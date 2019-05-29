package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	cfg "github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/p2p"

	"github.com/tendermint/tendermint/privval"

	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/helper"
	hmserver "github.com/maticnetwork/heimdall/server"
	"github.com/tendermint/tendermint/crypto"

	stakingcli "github.com/maticnetwork/heimdall/staking/cli"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/cosmos/cosmos-sdk/store"
)

// ValidatorAccountFormatter helps to print local validator account information
type ValidatorAccountFormatter struct {
	Address string `json:"address"`
	PrivKey string `json:"priv_key"`
	PubKey  string `json:"pub_key"`
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

	//rootCmd.PersistentFlags().String("log_level", ctx.Config.LogLevel, "Log level")
	tendermintCmd := &cobra.Command{
		Use:   "tendermint",
		Short: "Tendermint subcommands",
	}

	tendermintCmd.AddCommand(
		server.ShowNodeIDCmd(ctx),
		server.ShowValidatorCmd(ctx),
		server.ShowAddressCmd(ctx),
		server.VersionCmd(ctx),
	)
	rootCmd.AddCommand(
		server.StartCmd(ctx, server.AppCreator(newApp)),
		server.UnsafeResetAllCmd(ctx),
		client.LineBreak,
		client.LineBreak,
		tendermintCmd,
		version.VersionCmd,
	)
	rootCmd.AddCommand(newAccountCmd())
	rootCmd.AddCommand(hmserver.ServeCommands(cdc))
	rootCmd.AddCommand(InitCmd(ctx, cdc))
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

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB, storeTracer io.Writer, height int64, forZeroHeight bool, jailWhiteList []string,) (json.RawMessage, []tmTypes.GenesisValidator, error) {
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
				Address: "0x" + hex.EncodeToString(pubObject.Address().Bytes()),
				PrivKey: "0x" + hex.EncodeToString(privObject[:]),
				PubKey:  hex.EncodeToString(pubObject[:]),
			}

			b, err := json.Marshal(&account)
			if err != nil {
				panic(err)
			}

			// prints json info
			fmt.Printf("%s", string(b))
		},
	}
}

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
			nodeID, valPubKey, err :=  InitializeNodeValidatorFiles(config)
			if err != nil {
				return err
			}

			//
			// Heimdall config file
			//

			heimdallConf := helper.Configuration{
				MainRPCUrl:          helper.MainRPCUrl,
				MaticRPCUrl:         helper.MaticRPCUrl,
				StakeManagerAddress: (ethCommon.Address{}).Hex(),
				RootchainAddress:    (ethCommon.Address{}).Hex(),
				ChildBlockInterval:  helper.DefaultChildBlockInterval,

				CheckpointerPollInterval: helper.DefaultCheckpointerPollInterval,
				SyncerPollInterval:       helper.DefaultSyncerPollInterval,
				NoACKPollInterval:        helper.DefaultNoACKPollInterval,
				AvgCheckpointLength:      helper.DefaultCheckpointLength,
				MaxCheckpointLength:      helper.MaxCheckpointLength,
				NoACKWaitTime:            helper.NoACKWaitTime,
				CheckpointBufferTime:     helper.CheckpointBufferTime,
				ConfirmationBlocks: helper.ConfirmationBlocks,
			}

			heimdallConfBytes, err := json.MarshalIndent(heimdallConf, "", "  ")
			if err != nil {
				return err
			}

			if err := common.WriteFileAtomic(filepath.Join(config.RootDir, "config/heimdall-config.json"), heimdallConfBytes, 0600); err != nil {
				fmt.Println("Error writing heimdall-config", err)
				return err
			}

			//
			// Genesis file
			//
			validatorPublicKey := helper.GetPubObjects(valPubKey)
			newPubkey:=hmTypes.NewPubKey(validatorPublicKey[:])

			// create validator
			validator := app.GenesisValidator{
				ID:         hmTypes.NewValidatorID(uint64(validatorID)),
				PubKey:     newPubkey,
				StartEpoch: 0,
				Signer:     ethCommon.BytesToAddress(valPubKey.Address().Bytes()),
				Power:      1,
			}

			// create genesis state
			appState := &app.GenesisState{
				GenValidators: []app.GenesisValidator{validator},
			}

			appStateJSON, err := json.Marshal(appState)
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



func InitializeNodeValidatorFiles(
	config *cfg.Config) (nodeID string, valPubKey crypto.PubKey, err error,
) {

	nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
	if err != nil {
		return nodeID, valPubKey, err
	}

	nodeID = string(nodeKey.ID())
	server.UpgradeOldPrivValFile(config)

	pvKeyFile := config.PrivValidatorKeyFile()
	if err := common.EnsureDir(filepath.Dir(pvKeyFile), 0777); err != nil {
		return nodeID, valPubKey, nil
	}

	pvStateFile := config.PrivValidatorStateFile()
	if err := common.EnsureDir(filepath.Dir(pvStateFile), 0777); err != nil {
		return nodeID, valPubKey, nil
	}

	valPubKey = privval.LoadOrGenFilePV(pvKeyFile, pvStateFile).GetPubKey()
	return nodeID, valPubKey, nil
}

