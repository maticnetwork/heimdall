package cmd

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethCommon "github.com/maticnetwork/bor/common"
	"github.com/spf13/cobra"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/p2p"
	tmtime "github.com/tendermint/tendermint/types/time"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommon "github.com/maticnetwork/heimdall/types/common"
	bortypes "github.com/maticnetwork/heimdall/x/bor/types"
	stakingcli "github.com/maticnetwork/heimdall/x/staking/client/cli"
	stakingtypes "github.com/maticnetwork/heimdall/x/staking/types"
	topuptypes "github.com/maticnetwork/heimdall/x/topup/types"
)

var (
	flagNodeDirPrefix    = "node-dir-prefix"
	flagNumValidators    = "v"
	flagNumNonValidators = "n"
	flagOutputDir        = "output-dir"
	flagNodeDaemonHome   = "node-daemon-home"
	flagNodeCliHome      = "node-cli-home"
	flagNodeHostPrefix   = "node-host-prefix"
)

const (
	nodeDirPerm = 0755
)

// ValidatorAccountFormatter helps to print local validator account information
type ValidatorAccountFormatter struct {
	Address string `json:"address,omitempty" yaml:"address"`
	PrivKey string `json:"priv_key,omitempty" yaml:"priv_key"`
	PubKey  string `json:"pub_key,omitempty" yaml:"pub_key"`
}

// TestnetCmd initialises files required to start heimdall testnet
func testnetCmd(ctx *server.Context) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, _ []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			depCdc := clientCtx.JSONMarshaler
			cdc := depCdc.(codec.Marshaler)

			serverCtx := server.GetServerContextFromCmd(cmd)
			config := serverCtx.Config

			config.SetRoot(clientCtx.HomeDir)

			outDir, err := cmd.Flags().GetString(flagOutputDir)
			if err != nil {
				return err
			}

			// create chain id
			chainID, err := cmd.Flags().GetString(flags.FlagChainID)
			if err != nil {
				return err
			}
			// TODO - check this
			// if chainID == "" {
			// 	chainID = fmt.Sprintf("heimdall-%v", common.RandStr(6))
			// }

			// num of validators = validators in genesis files
			numValidators, err := cmd.Flags().GetInt(flagNumValidators)
			if err != nil {
				return err
			}

			// get total number of validators to be generated
			totalValidators := totalValidators(cmd)

			// first validators start ID
			startID, err := cmd.Flags().GetInt64(stakingcli.FlagValidatorID)
			if err != nil {
				return err
			}
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
			dividendAccounts := make([]*hmTypes.DividendAccount, numValidators)

			// slashing
			valSigningInfoMap := make(map[string]hmTypes.ValidatorSigningInfo)

			genFiles := make([]string, totalValidators)

			nodeDaemonHomeName, err := cmd.Flags().GetString(flagNodeDaemonHome)
			if err != nil {
				return err
			}

			nodeCliHomeName, err := cmd.Flags().GetString(flagNodeCliHome)
			if err != nil {
				return err
			}

			// get genesis time
			genesisTime := tmtime.Now()

			for i := 0; i < totalValidators; i++ {
				// get node dir name = PREFIX+INDEX
				nodeDirName, err := cmd.Flags().GetString(flagNodeDirPrefix)
				if err != nil {
					return err
				}
				nodeDirName = fmt.Sprintf("%s%d", nodeDirName, i)

				// generate node and client dir
				nodeDir := filepath.Join(outDir, nodeDirName, nodeDaemonHomeName)
				clientDir := filepath.Join(outDir, nodeDirName, nodeCliHomeName)

				// set root in config
				config.SetRoot(nodeDir)

				// create config folder
				err = os.MkdirAll(filepath.Join(nodeDir, "config"), nodeDirPerm)
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
					sdkAddress := hmCommon.HeimdallAddressToAccAddress(hmCommon.BytesToHeimdallAddress(valPubKeys[i].Address().Bytes()))
					// create validator account
					validators[i] = hmTypes.NewValidator(
						hmTypes.NewValidatorID(uint64(startID+int64(i))),
						0,
						0,
						1,
						10000,
						newPubkey,
						sdkAddress,
					)

					signer, _ := sdk.AccAddressFromHex(validators[i].Signer)
					// create dividend account for validator
					dividendAcc := hmTypes.NewDividendAccount(signer, ZeroIntString)
					dividendAccounts[i] = &dividendAcc
					valSigningInfoMap[validators[i].ID.String()] = hmTypes.NewValidatorSigningInfo(validators[i].ID, 0, 0, 0)
				}

				signers[i] = GetSignerInfo(valPubKeys[i], privKeys[i].Bytes(), cdc)

				WriteDefaultHeimdallConfig(filepath.Join(config.RootDir, "config/heimdall-config.toml"), helper.GetDefaultHeimdallConfig())
			}

			// other data
			accounts := make([]authtypes.GenesisAccount, totalValidators)
			for i := 0; i < totalValidators; i++ {
				populatePersistentPeersInConfigAndWriteIt(config, cmd)
				// genesis account
				accounts[i] = getGenesisAccount([]byte(validators[i].Signer), valPubKeys[i].Address().Bytes())
			}
			validatorSet := hmTypes.NewValidatorSet(validators)

			// new app state
			appStateBytes := app.NewDefaultGenesisState()

			authGenState := authtypes.GetGenesisStateFromAppState(cdc, appStateBytes)
			anyGenAccounts, err := authtypes.PackAccounts(accounts)
			if err != nil {
				return err
			}
			authGenState.Accounts = anyGenAccounts
			appStateBytes[authtypes.ModuleName] = authtypes.ModuleCdc.MustMarshalJSON(&authGenState)

			// staking state change
			appStateBytes, err = stakingtypes.SetGenesisStateToAppState(cdc, appStateBytes, validators, validatorSet)
			if err != nil {
				return err
			}

			// bor state change
			appStateBytes, err = bortypes.SetGenesisStateToAppState(appStateBytes, *validatorSet)
			if err != nil {
				return err
			}

			// topup state change
			appStateBytes, err = topuptypes.SetGenesisStateToAppState(appStateBytes, dividendAccounts)
			if err != nil {
				return err
			}

			// TODO - Uncomment when slashing is added
			// slashing state change
			// appStateBytes, err = slashingtypes.SetGenesisStateToAppState(appStateBytes, valSigningInfoMap)
			// if err != nil {
			// 	return err
			// }

			appStateJSON, err := json.Marshal(appStateBytes)
			if err != nil {
				return err
			}

			for i := 0; i < totalValidators; i++ {
				if err = writeGenesisFile(genesisTime, genFiles[i], chainID, appStateJSON); err != nil {
					return err
				}
			}

			// TODO - check this
			// dump signer information in a json file
			// TODO move to const string flag
			// dump := cmd.Flags().GetBool("signer-dump")
			// if dump {
			// 	signerJSON, err := json.MarshalIndent(signers, "", "  ")
			// 	if err != nil {
			// 		return err
			// 	}

			// 	if err := common.WriteFileAtomic(filepath.Join(outDir, "signer-dump.json"), signerJSON, 0600); err != nil {
			// 		fmt.Println("Error writing signer-dump", err)
			// 		return err
			// 	}
			// }

			fmt.Printf("Successfully initialized %d node directories\n", totalValidators)
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

	cmd.Flags().String(flags.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().Bool("signer-dump", true, "dumps all signer information in a json file")
	return cmd
}

// Total Validators to be included in the testnet
func totalValidators(cmd *cobra.Command) int {
	numValidators, _ := cmd.Flags().GetInt(flagNumValidators)
	numNonValidators, _ := cmd.Flags().GetInt(flagNumNonValidators)
	return numNonValidators + numValidators
}

// WriteDefaultHeimdallConfig writes default heimdall config to the given path
func WriteDefaultHeimdallConfig(path string, conf helper.Configuration) {
	heimdallConf := helper.GetDefaultHeimdallConfig()
	helper.WriteConfigFile(path, &heimdallConf)
}

// populate persistent peers in config
func populatePersistentPeersInConfigAndWriteIt(config *cfg.Config, cmd *cobra.Command) {
	persistentPeers := make([]string, totalValidators(cmd))
	for i := 0; i < totalValidators(cmd); i++ {
		config.SetRoot(nodeDir(i, cmd))
		nodeKey, err := p2p.LoadNodeKey(config.NodeKeyFile())
		if err != nil {
			return
		}
		persistentPeers[i] = p2p.IDAddressString(nodeKey.ID(), fmt.Sprintf("%s:%d", hostnameOrIP(i, cmd), 26656))
	}

	persistentPeersList := strings.Join(persistentPeers, ",")
	for i := 0; i < totalValidators(cmd); i++ {
		config.SetRoot(nodeDir(i, cmd))
		config.P2P.PersistentPeers = persistentPeersList
		config.P2P.AddrBookStrict = false

		// overwrite default config
		cfg.WriteConfigFile(filepath.Join(nodeDir(i, cmd), "config", "config.toml"), config)
	}
}

// get node directory path
func nodeDir(i int, cmd *cobra.Command) string {
	outDir, _ := cmd.Flags().GetString(flagOutputDir)
	nodeDirName, _ := cmd.Flags().GetString(flagNodeDirPrefix)
	nodeDirName = fmt.Sprintf("%s%d", nodeDirName, i)
	nodeDaemonHomeName, _ := cmd.Flags().GetString(flagNodeDaemonHome)
	return filepath.Join(outDir, nodeDirName, nodeDaemonHomeName)
}

// hostname of ip of nodes
func hostnameOrIP(i int, cmd *cobra.Command) string {
	hOrIP, _ := cmd.Flags().GetString(flagNodeHostPrefix)
	return fmt.Sprintf("%s%d", hOrIP, i)
}

// GetSignerInfo returns signer information
func GetSignerInfo(pub crypto.PubKey, priv []byte, cdc codec.Marshaler) ValidatorAccountFormatter {
	privObject := secp256k1.GenPrivKeyFromSecret(priv)
	return ValidatorAccountFormatter{
		Address: ethCommon.BytesToAddress(pub.Address().Bytes()).String(),
		PubKey:  CryptoKeyToPubkey(pub).String(),
		PrivKey: "0x" + hex.EncodeToString(privObject.Bytes()),
	}
}
