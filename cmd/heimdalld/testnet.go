package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	ethCommon "github.com/maticnetwork/bor/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/maticnetwork/heimdall/app"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/helper"
	stakingcli "github.com/maticnetwork/heimdall/staking/client/cli"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// TestnetCmd initialises files required to start heimdall testnet
func TestnetCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
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
			numValidators := viper.GetInt(flagNumValidators)
			numNonValidators := viper.GetInt(flagNumNonValidators)
			startID := viper.GetInt64(stakingcli.FlagValidatorID)
			if startID == 0 {
				startID = 1
			}
			totalValidators := totalValidators()
			signers := make([]ValidatorAccountFormatter, totalValidators)
			nodeIDs := make([]string, totalValidators)
			valPubKeys := make([]crypto.PubKey, totalValidators)
			privKeys := make([]crypto.PrivKey, totalValidators)
			validators := make([]*hmTypes.Validator, totalValidators)
			dividendAccounts := make([]hmTypes.DividendAccount, totalValidators)
			genFiles := make([]string, totalValidators)
			var err error
			// create chain id
			chainID := viper.GetString(client.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("heimdall-%v", common.RandStr(6))
			}

			for i := 0; i < numValidators+numNonValidators; i++ {
				nodeDirName := fmt.Sprintf("%s%d", viper.GetString(flagNodeDirPrefix), i)
				nodeDaemonHomeName := viper.GetString(flagNodeDaemonHome)
				nodeCliHomeName := viper.GetString(flagNodeCliHome)
				nodeDir := filepath.Join(outDir, nodeDirName, nodeDaemonHomeName)
				clientDir := filepath.Join(outDir, nodeDirName, nodeCliHomeName)
				config.SetRoot(nodeDir)

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
				//
				// Genesis file
				//
				validatorPublicKey := helper.GetPubObjects(valPubKeys[i])
				newPubkey := hmTypes.NewPubKey(validatorPublicKey[:])

				// create validator
				validators[i] = &hmTypes.Validator{
					ID:          hmTypes.NewValidatorID(uint64(startID + int64(i))),
					PubKey:      newPubkey,
					StartEpoch:  0,
					Signer:      hmTypes.BytesToHeimdallAddress(valPubKeys[i].Address().Bytes()),
					VotingPower: 1,
				}

				dividendAccounts[i] = hmTypes.DividendAccount{
					ID:            hmTypes.NewDividendAccountID(uint64(validators[i].ID)),
					FeeAmount:     big.NewInt(0).String(),
					SlashedAmount: big.NewInt(0).String(),
				}

				var privObject secp256k1.PrivKeySecp256k1
				cdc.MustUnmarshalBinaryBare(privKeys[i].Bytes(), &privObject)
				signers[i] = ValidatorAccountFormatter{
					Address: ethCommon.BytesToAddress(valPubKeys[i].Address().Bytes()).String(),
					PubKey:  newPubkey.String(),
					PrivKey: "0x" + hex.EncodeToString(privObject[:]),
				}

				heimdallConf := helper.GetDefaultHeimdallConfig()
				helper.WriteConfigFile(filepath.Join(config.RootDir, "config/heimdall-config.toml"), &heimdallConf)
			}

			// other data
			accounts := make([]authTypes.GenesisAccount, totalValidators)
			for i := 0; i < totalValidators; i++ {
				populatePersistentPeersInConfigAndWriteIt(config)
				// genesis account
				accounts[i] = getGenesisAccount(validators[i].Signer.Bytes())
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
			appStateBytes, err = stakingTypes.SetGenesisStateToAppState(appStateBytes, validators, *validatorSet, dividendAccounts)
			if err != nil {
				return err
			}

			appStateJSON, err := json.Marshal(appStateBytes)
			if err != nil {
				return err
			}

			for i := 0; i < len(validators); i++ {
				writeGenesisFile(genFiles[i], chainID, appStateJSON)
			}

			// dump signer information in a json file
			// TODO move to const string flag
			dump := viper.GetBool("signer-dump")
			if dump {
				signerJSON, err := json.MarshalIndent(signers, "", "  ")
				if err != nil {
					return err
				}

				if err := common.WriteFileAtomic(filepath.Join(outDir, "signer-dump.json"), signerJSON, 0600); err != nil {
					fmt.Println("Error writing signer-dump", err)
					return err
				}
			}

			fmt.Printf("Successfully initialized %d node directories\n", numValidators+numNonValidators)
			return nil
		},
	}

	cmd.Flags().Int(flagNumValidators, 4,
		"Number of validators to initialize the testnet with",
	)

	cmd.Flags().Int(flagNumNonValidators, 8,
		"Number of validators to initialize the testnet with",
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

	cmd.Flags().String(client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().Bool("signer-dump", true, "dumps all signer information in a json file")
	return cmd
}
