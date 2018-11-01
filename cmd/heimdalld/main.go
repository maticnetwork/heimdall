package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/cli"
	dbm "github.com/tendermint/tendermint/libs/db"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/app"
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

	// add custom root command
	rootCmd.AddCommand(newAccountCmd())

	// cosmos server commands
	server.AddCommands(
		ctx,
		cdc,
		rootCmd,
		server.DefaultAppInit,
		server.ConstructAppCreator(newApp, app.AppName),
		server.ConstructAppExporter(exportAppStateAndTMValidators, app.AppName),
	)

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
