package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethCommon "github.com/maticnetwork/bor/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	tmTypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/maticnetwork/heimdall/app"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	borTypes "github.com/maticnetwork/heimdall/bor/types"
	"github.com/maticnetwork/heimdall/helper"
	hmserver "github.com/maticnetwork/heimdall/server"
	stakingcli "github.com/maticnetwork/heimdall/staking/client/cli"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
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

const nodeDirPerm = 0755

// ValidatorAccountFormatter helps to print local validator account information
type ValidatorAccountFormatter struct {
	Address          string `json:"address" yaml:"address"`
	PrivKey          string `json:"priv_key" yaml:"priv_key"`
	PubKey           string `json:"pub_key" yaml:"pub_key"`
	AccountAddress   string `json:"account_address" yaml:"account_address"`
	AccountPubKey    string `json:"account_pubkey" yaml:"account_pubkey"`
	ValidatorAddress string `json:"validator_address" yaml:"validator_address"`
	ValidatorPubKey  string `json:"validator_pubkey" yaml:"validator_pubkey"`
	ConsensusAddress string `json:"consensus_address" yaml:"consensus_address"`
	ConsensusPubKey  string `json:"consensus_pubkey" yaml:"consensus_pubkey"`
}

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

	// just make pulp :)
	app.MakePulp()

	rootCmd := &cobra.Command{
		Use:               "heimdalld",
		Short:             "Heimdall Daemon (server)",
		PersistentPreRunE: server.PersistentPreRunEFn(ctx),
	}

	// add new persistent flag for heimdall-config
	rootCmd.PersistentFlags().String(
		helper.WithHeimdallConfigFlag,
		"",
		"Heimdall config file path (default <home>/config/heimdall-config.json)",
	)

	// bind with-heimdall-config config with root cmd
	viper.BindPFlag(
		helper.WithHeimdallConfigFlag,
		rootCmd.Flags().Lookup(helper.WithHeimdallConfigFlag),
	)
	server.AddCommands(ctx, cdc, rootCmd, newApp, exportAppStateAndTMValidators)
	rootCmd.AddCommand(newAccountCmd())
	rootCmd.AddCommand(hmserver.ServeCommands(cdc, hmserver.RegisterRoutes))
	rootCmd.AddCommand(InitCmd(ctx, cdc))
	rootCmd.AddCommand(TestnetCmd(ctx, cdc))
	rootCmd.AddCommand(VerifyGenesis(ctx, cdc))

	// prepare and add flags
	executor := cli.PrepareBaseCmd(rootCmd, "HD", os.ExpandEnv("$HOME/.heimdalld"))
	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}
}

func newApp(logger log.Logger, db dbm.DB, storeTracer io.Writer) abci.Application {
	// init heimdall config
	helper.InitHeimdallConfig("")
	// create new heimdall app
	return app.NewHeimdallApp(logger, db, baseapp.SetPruning(store.NewPruningOptionsFromString(viper.GetString("pruning"))))
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB, storeTracer io.Writer, height int64, forZeroHeight bool, jailWhiteList []string) (json.RawMessage, []tmTypes.GenesisValidator, error) {
	bapp := app.NewHeimdallApp(logger, db)
	return bapp.ExportAppStateAndValidators()
}

func newAccountCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show-account",
		Short: "Print the account's private key and public key",
		Run: func(cmd *cobra.Command, args []string) {
			// init heimdall config
			helper.InitHeimdallConfig("")

			// get private and public keys
			privObject := helper.GetPrivKey()
			pubObject := helper.GetPubKey()

			account := &ValidatorAccountFormatter{
				Address: ethCommon.BytesToAddress(pubObject.Address().Bytes()).String(),
				PrivKey: "0x" + hex.EncodeToString(privObject[:]),
				PubKey:  "0x" + hex.EncodeToString(pubObject[:]),
			}

			b, err := json.MarshalIndent(account, "", "    ")
			if err != nil {
				panic(err)
			}

			// prints json info
			fmt.Printf("%s", string(b))
		},
	}
}

// InitCmd initialises files required to start heimdall
func InitCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize genesis config, priv-validator file, and p2p-node file",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			// create chain id
			chainID := viper.GetString(client.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("heimdall-%v", common.RandStr(6))
			}

			validatorID := viper.GetInt64(stakingcli.FlagValidatorID)
			nodeID, valPubKey, _, err := InitializeNodeValidatorFiles(config)
			if err != nil {
				return err
			}

			// Heimdall config file
			heimdallConf := helper.GetDefaultHeimdallConfig()
			helper.WriteConfigFile(filepath.Join(config.RootDir, "config/heimdall-config.toml"), &heimdallConf)

			//
			// Genesis file
			//
			validatorPublicKey := helper.GetPubObjects(valPubKey)
			newPubkey := hmTypes.NewPubKey(validatorPublicKey[:])

			// create validator
			validator := hmTypes.Validator{
				ID:                   hmTypes.NewValidatorID(uint64(validatorID)),
				PubKey:               newPubkey,
				StartEpoch:           0,
				Signer:               hmTypes.BytesToHeimdallAddress(valPubKey.Address().Bytes()),
				VotingPower:          stakingTypes.DefaultValPower,
				DelegatedPower:       0,
				DelgatorRewardPool:   "",
				TotalDelegatorShares: "",
			}

			vals := []*hmTypes.Validator{&validator}
			validatorSet := hmTypes.NewValidatorSet(vals)

			// create genesis state
			appStateBytes := app.NewDefaultGenesisState()

			// auth state change
			appStateBytes, err = authTypes.SetGenesisStateToAppState(
				appStateBytes,
				[]authTypes.GenesisAccount{getGenesisAccount(validator.Signer.Bytes())},
			)
			if err != nil {
				return err
			}

			// staking state change
			appStateBytes, err = stakingTypes.SetGenesisStateToAppState(appStateBytes, vals, *validatorSet)
			if err != nil {
				return err
			}

			// bor state change
			appStateBytes, err = borTypes.SetGenesisStateToAppState(appStateBytes, *validatorSet)
			if err != nil {
				return err
			}

			// app state json
			appStateJSON, err := json.Marshal(appStateBytes)
			if err != nil {
				return err
			}

			toPrint := struct {
				ChainID string `json:"chain_id"`
				NodeID  string `json:"node_id"`
			}{
				chainID,
				nodeID,
			}

			out, err := codec.MarshalJSONIndent(cdc, toPrint)
			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stderr, "%s\n", string(out))
			return writeGenesisFile(config.GenesisFile(), chainID, appStateJSON)
		},
	}

	cmd.Flags().String(cli.HomeFlag, helper.DefaultNodeHome, "node's home directory")
	cmd.Flags().String(helper.FlagClientHome, helper.DefaultCLIHome, "client's home directory")
	cmd.Flags().String(client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().Int(stakingcli.FlagValidatorID, 1, "--id=<validator ID here>, if left blank will be assigned 1")
	return cmd
}

// VerifyGenesis verifies the genesis file and brings it in sync with on-chain contract
func VerifyGenesis(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "verify-genesis",
		Short: "Verify if the genesis matches",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))
			helper.InitHeimdallConfig("")

			// Loading genesis doc
			genDoc, err := tmTypes.GenesisDocFromFile(filepath.Join(config.RootDir, "config/genesis.json"))
			if err != nil {
				return err
			}

			// get genesis state
			var genesisState app.GenesisState
			err = json.Unmarshal(genDoc.AppState, &genesisState)
			if err != nil {
				return err
			}

			// verify genesis
			for _, b := range app.ModuleBasics {
				m := b.(hmTypes.HeimdallModuleBasic)
				if err := m.VerifyGenesis(genesisState); err != nil {
					return err
				}
			}

			return nil
		},
	}

	return cmd
}

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

				var privObject secp256k1.PrivKeySecp256k1
				cdc.MustUnmarshalBinaryBare(privKeys[i].Bytes(), &privObject)
				signers[i] = ValidatorAccountFormatter{
					Address:          ethCommon.BytesToAddress(valPubKeys[i].Address().Bytes()).String(),
					PubKey:           newPubkey.String(),
					PrivKey:          "0x" + hex.EncodeToString(privObject[:]),
					AccountAddress:   sdk.AccAddress(valPubKeys[i].Address().Bytes()).String(),
					AccountPubKey:    sdk.MustBech32ifyAccPub(valPubKeys[i]),
					ValidatorAddress: sdk.ValAddress(valPubKeys[i].Address().Bytes()).String(),
					ValidatorPubKey:  sdk.MustBech32ifyValPub(valPubKeys[i]),
					ConsensusAddress: sdk.ConsAddress(valPubKeys[i].Address().Bytes()).String(),
					ConsensusPubKey:  sdk.MustBech32ifyConsPub(valPubKeys[i]),
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
			appStateBytes, err = stakingTypes.SetGenesisStateToAppState(appStateBytes, validators, *validatorSet)
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

// Total Validators to be included in the testnet
func totalValidators() int {
	numValidators := viper.GetInt(flagNumValidators)
	numNonValidators := viper.GetInt(flagNumNonValidators)
	return numNonValidators + numValidators
}

// get node directory path
func nodeDir(i int) string {
	outDir := viper.GetString(flagOutputDir)
	nodeDirName := fmt.Sprintf("%s%d", viper.GetString(flagNodeDirPrefix), i)
	nodeDaemonHomeName := viper.GetString(flagNodeDaemonHome)
	return filepath.Join(outDir, nodeDirName, nodeDaemonHomeName)
}

// hostname of ip of nodes
func hostnameOrIP(i int) string {
	return fmt.Sprintf("%s%d", viper.GetString(flagNodeHostPrefix), i)
}

// populate persistent peers in config
func populatePersistentPeersInConfigAndWriteIt(config *cfg.Config) {
	persistentPeers := make([]string, totalValidators())
	for i := 0; i < totalValidators(); i++ {
		config.SetRoot(nodeDir(i))
		nodeKey, err := p2p.LoadNodeKey(config.NodeKeyFile())
		if err != nil {
			return
		}
		persistentPeers[i] = p2p.IDAddressString(nodeKey.ID(), fmt.Sprintf("%s:%d", hostnameOrIP(i), 26656))
	}

	persistentPeersList := strings.Join(persistentPeers, ",")
	for i := 0; i < totalValidators(); i++ {
		config.SetRoot(nodeDir(i))
		config.P2P.PersistentPeers = persistentPeersList
		config.P2P.AddrBookStrict = false

		// overwrite default config
		cfg.WriteConfigFile(filepath.Join(nodeDir(i), "config", "config.toml"), config)
	}
}

func getGenesisAccount(address []byte) authTypes.GenesisAccount {
	acc := authTypes.NewBaseAccountWithAddress(hmTypes.BytesToHeimdallAddress(address))
	acc.SetCoins(hmTypes.Coins{hmTypes.Coin{Denom: "vetic", Amount: hmTypes.NewInt(1000)}})
	return authTypes.BaseToGenesisAccount(acc)
}

// WriteGenesisFile creates and writes the genesis configuration to disk. An
// error is returned if building or writing the configuration to file fails.
// nolint: unparam
func writeGenesisFile(genesisFile, chainID string, appState json.RawMessage) error {
	genDoc := tmTypes.GenesisDoc{
		ChainID:  chainID,
		AppState: appState,
	}

	if err := genDoc.ValidateAndComplete(); err != nil {
		return err
	}

	return genDoc.SaveAs(genesisFile)
}

// InitializeNodeValidatorFiles initializes node and priv validator files
func InitializeNodeValidatorFiles(
	config *cfg.Config) (nodeID string, valPubKey crypto.PubKey, priv crypto.PrivKey, err error,
) {

	nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
	if err != nil {
		return nodeID, valPubKey, priv, err
	}

	nodeID = string(nodeKey.ID())
	server.UpgradeOldPrivValFile(config)

	pvKeyFile := config.PrivValidatorKeyFile()
	if err := common.EnsureDir(filepath.Dir(pvKeyFile), 0777); err != nil {
		return nodeID, valPubKey, priv, nil
	}

	pvStateFile := config.PrivValidatorStateFile()
	if err := common.EnsureDir(filepath.Dir(pvStateFile), 0777); err != nil {
		return nodeID, valPubKey, priv, nil
	}

	FilePv := privval.LoadOrGenFilePV(pvKeyFile, pvStateFile)
	valPubKey = FilePv.GetPubKey()
	return nodeID, valPubKey, FilePv.Key.PrivKey, nil
}
