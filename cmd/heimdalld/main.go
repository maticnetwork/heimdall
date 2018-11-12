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
	gaiaInit "github.com/cosmos/cosmos-sdk/cmd/gaia/init"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/p2p"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/app"
	checkpointRestCmds "github.com/maticnetwork/heimdall/checkpoint/rest"
	"github.com/maticnetwork/heimdall/helper"
)

func main() {
	cdc := app.MakeCodec()
	ctx := server.NewDefaultContext()

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

	// cosmos server commands
	server.AddCommands(
		ctx,
		cdc,
		rootCmd,
		server.DefaultAppInit,
		server.AppCreator(newApp),
		server.AppExporter(exportAppStateAndTMValidators),
	)

	rootCmd.AddCommand(newAccountCmd())

	rootCmd.AddCommand(checkpointRestCmds.ServeCommands(cdc))
	rootCmd.AddCommand(InitCmd(ctx, cdc, server.DefaultAppInit))
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
	helper.InitHeimdallConfig()

	// create new heimdall app
	return app.NewHeimdallApp(logger, db, baseapp.SetPruning(viper.GetString("pruning")))
}

func exportAppStateAndTMValidators(logger log.Logger, db dbm.DB, storeTracer io.Writer) (json.RawMessage, []tmtypes.GenesisValidator, error) {
	bapp := app.NewHeimdallApp(logger, db)
	return bapp.ExportAppStateAndValidators()
}

// InitCmd get cmd to initialize all files for tendermint and application
// nolint: errcheck
func InitCmd(ctx *server.Context, cdc *codec.Codec, appInit server.AppInit) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize genesis config, priv-validator file, and p2p-node file",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {

			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))
			chainID := viper.GetString(client.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("heimdall-%v", common.RandStr(6))
			}

			nodeKey, err := p2p.LoadOrGenNodeKey(config.NodeKeyFile())
			if err != nil {
				return err
			}
			nodeID := string(nodeKey.ID())

			pk := gaiaInit.ReadOrCreatePrivValidator(config.PrivValidatorFile())

			// TODO pull validator from main chain and add to genesis
			genTx, appMessage, validator, err := server.SimpleAppGenTx(cdc, pk)
			if err != nil {
				return err
			}

			appState, err := appInit.AppGenState(cdc, []json.RawMessage{genTx})
			if err != nil {
				return err
			}
			appStateJSON, err := cdc.MarshalJSON(appState)
			if err != nil {
				return err
			}

			heimdallConf := helper.Configuration{
				MainRPCUrl:          helper.MainRPCUrl,
				MaticRPCUrl:         helper.MaticRPCUrl,
				StakeManagerAddress: "",
				RootchainAddress:    "",
			}
			heimdallConfBytes, err := cdc.MarshalJSONIndent(heimdallConf, "", "  ")
			if err != nil {
				return err
			}

			if err := common.WriteFileAtomic(filepath.Join(config.RootDir, "config/heimdall-config.json"), heimdallConfBytes, 0600); err != nil {
				fmt.Errorf("Error writing heimdall-config %s\n", err)
				return err
			}

			toPrint := struct {
				ChainID    string          `json:"chain_id"`
				NodeID     string          `json:"node_id"`
				AppMessage json.RawMessage `json:"app_message"`
			}{
				chainID,
				nodeID,
				appMessage,
			}

			out, err := codec.MarshalJSONIndent(cdc, toPrint)
			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stderr, "%s\n", string(out))
			return gaiaInit.WriteGenesisFile(config.GenesisFile(), chainID, []tmtypes.GenesisValidator{validator}, appStateJSON)
		},
	}

	cmd.Flags().String(cli.HomeFlag, helper.DefaultNodeHome, "node's home directory")
	cmd.Flags().String(helper.FlagClientHome, helper.DefaultCLIHome, "client's home directory")
	cmd.Flags().String(client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	cmd.Flags().String(client.FlagName, "", "validator's moniker")
	return cmd
}

func newAccountCmd() *cobra.Command {
	type Account struct {
		Address string `json:"address"`
		PrivKey string `json:"private_key"`
		PubKey  string `json:"public_key"`
	}

	return &cobra.Command{
		Use:   "show-account",
		Short: "Print the account's private key and public key",
		Run: func(cmd *cobra.Command, args []string) {
			// init heimdall config
			helper.InitHeimdallConfig()

			// get private and public keys
			privObject := helper.GetPrivKey()
			pubObject := helper.GetPubKey()

			account := &Account{
				Address: "0x" + hex.EncodeToString(pubObject.Address().Bytes()),
				PrivKey: "0x" + hex.EncodeToString(privObject[:]),
				PubKey:  "0x" + hex.EncodeToString(pubObject[:]),
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
