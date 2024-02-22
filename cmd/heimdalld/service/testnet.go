package service

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/common"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/maticnetwork/heimdall/app"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	borTypes "github.com/maticnetwork/heimdall/bor/types"
	"github.com/maticnetwork/heimdall/helper"
	slashingTypes "github.com/maticnetwork/heimdall/slashing/types"
	stakingcli "github.com/maticnetwork/heimdall/staking/client/cli"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
	topupTypes "github.com/maticnetwork/heimdall/topup/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// TestnetCmd initialises files required to start heimdall testnet
func testnetCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-testnet",
		Short: "Initialize files for a Heimdall testnet",
		Long: `testnet will create "v" + "n" number of directories and populate each with
necessary files (private validator, genesis, config, etc.).

Note, strict routability for addresses is turned off in the config file.
Optionally, it will fill in persistent_peers list in config file using either hostnames or IPs.

Example:
testnet --v 4 --n 8 --output-dir ./output --starting-ip-address 192.168.10.2
`,
		Args: cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			outDir := viper.GetString(flagOutputDir)

			// create chain id
			chainID := viper.GetString(client.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("heimdall-%v", common.RandStr(6))
			}

			// num of validators = validators in genesis files
			numValidators := viper.GetInt(flagNumValidators)

			// get total number of validators to be generated
			totalValidators := totalValidators()

			// first validators start ID
			startID := viper.GetInt64(stakingcli.FlagValidatorID)
			if startID == 0 {
				startID = 1
			}

			// signers data to dump in the signer-dump file
			signers := make([]ValidatorAccountFormatter, totalValidators)

			// Initialise variables for all validators
			nodeIDs := make([]string, totalValidators)
			valPubKeys := make([]crypto.PubKey, totalValidators)
			privKeys := make([]crypto.PrivKey, totalValidators)
			validators := make([]*hmTypes.Validator, numValidators)
			dividendAccounts := make([]hmTypes.DividendAccount, numValidators)

			// slashing
			valSigningInfoMap := make(map[string]hmTypes.ValidatorSigningInfo)

			genFiles := make([]string, totalValidators)
			var err error

			nodeDaemonHomeName := viper.GetString(flagNodeDaemonHome)
			nodeCliHomeName := viper.GetString(flagNodeCliHome)

			// get genesis time
			genesisTime := tmtime.Now()

			for i := 0; i < totalValidators; i++ {
				// get node dir name = PREFIX+INDEX
				nodeDirName := fmt.Sprintf("%s%d", viper.GetString(flagNodeDirPrefix), i)

				// generate node and client dir
				nodeDir := filepath.Join(outDir, nodeDirName, nodeDaemonHomeName)
				clientDir := filepath.Join(outDir, nodeDirName, nodeCliHomeName)

				// set root in config
				config.SetRoot(nodeDir)

				// create config folder
				err := os.MkdirAll(filepath.Join(nodeDir, "config"), nodeDirPerm)
				if err != nil {
					_ = os.RemoveAll(outDir)
					return err
				}

				err = os.MkdirAll(clientDir, nodeDirPerm)
				if err != nil {
					_ = os.RemoveAll(outDir)
					return err
				}

				nodeIDs[i], valPubKeys[i], privKeys[i], err = InitializeNodeValidatorFiles(config)
				if err != nil {
					return err
				}

				genFiles[i] = config.GenesisFile()
				newPubkey := CryptoKeyToPubkey(valPubKeys[i])

				if i < numValidators {
					// create validator account
					validators[i] = hmTypes.NewValidator(
						hmTypes.NewValidatorID(uint64(startID+int64(i))),
						0,
						0,
						1,
						10000,
						newPubkey,
						hmTypes.BytesToHeimdallAddress(valPubKeys[i].Address().Bytes()),
					)

					// create dividend account for validator
					dividendAccounts[i] = hmTypes.NewDividendAccount(validators[i].Signer, ZeroIntString)
					valSigningInfoMap[validators[i].ID.String()] = hmTypes.NewValidatorSigningInfo(validators[i].ID, 0, 0, 0)
				}

				signers[i] = GetSignerInfo(valPubKeys[i], privKeys[i].Bytes(), cdc)

				WriteDefaultHeimdallConfig(filepath.Join(config.RootDir, "config/heimdall-config.toml"), helper.GetDefaultHeimdallConfig())
			}

			// other data
			accounts := make([]authTypes.GenesisAccount, totalValidators)
			for i := 0; i < totalValidators; i++ {
				populatePersistentPeersInConfigAndWriteIt(config)
				// genesis account
				accounts[i] = getGenesisAccount(valPubKeys[i].Address().Bytes())
			}
			validatorSet := hmTypes.NewValidatorSet(validators)

			// new app state
			appStateBytes := app.NewDefaultGenesisState()

			// auth state change
			appStateBytes, err = authTypes.SetGenesisStateToAppState(appStateBytes, accounts)
			if err != nil {
				return err
			}

			// staking state change
			appStateBytes, err = stakingTypes.SetGenesisStateToAppState(appStateBytes, validators, *validatorSet)
			if err != nil {
				return err
			}

			// bor state change
			appStateBytes, err = borTypes.SetGenesisStateToAppState(appStateBytes, *validatorSet)
			if err != nil {
				return err
			}

			// topup state change
			appStateBytes, err = topupTypes.SetGenesisStateToAppState(appStateBytes, dividendAccounts)
			if err != nil {
				return err
			}

			// slashing state change
			appStateBytes, err = slashingTypes.SetGenesisStateToAppState(appStateBytes, valSigningInfoMap)
			if err != nil {
				return err
			}

			appStateJSON, err := jsoniter.ConfigCompatibleWithStandardLibrary.Marshal(appStateBytes)
			if err != nil {
				return err
			}

			for i := 0; i < totalValidators; i++ {
				if err = writeGenesisFile(genesisTime, genFiles[i], chainID, appStateJSON); err != nil {
					return err
				}
			}

			// dump signer information in a json file
			// TODO move to const string flag
			dump := viper.GetBool("signer-dump")
			if dump {
				signerJSON, err := jsoniter.ConfigFastest.MarshalIndent(signers, "", "  ")
				if err != nil {
					return err
				}

				if err := common.WriteFileAtomic(filepath.Join(outDir, "signer-dump.json"), signerJSON, 0600); err != nil {
					fmt.Println("Error writing signer-dump", err)
					return err
				}
			}

			fmt.Printf("Successfully initialized %d node directories\n", totalValidators)
			return nil
		},
	}

	cmd.Flags().Int(flagNumValidators, 4,
		"Number of validators to initialize the testnet with",
	)

	cmd.Flags().Int(flagNumNonValidators, 8,
		"Number of non-validators to initialize the testnet with",
	)

	cmd.Flags().StringP(flagOutputDir, "o", "./mytestnet",
		"Directory to store initialization data for the testnet",
	)

	cmd.Flags().String(flagNodeDirPrefix, "node",
		"Prefix the directory name for each node with (node results in node0, node1, ...)",
	)

	cmd.Flags().String(flagNodeDaemonHome, "heimdalld",
		"Home directory of the node's daemon configuration",
	)

	cmd.Flags().String(flagNodeCliHome, "heimdallcli",
		"Home directory of the node's cli configuration",
	)

	cmd.Flags().String(flagNodeHostPrefix, "node",
		"Hostname prefix (node results in persistent peers list ID0@node0:26656, ID1@node1:26656, ...)")

	cmd.Flags().String(client.FlagChainID, "", "Genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().Bool("signer-dump", true, "Dumps all signer information in a json file")

	return cmd
}
