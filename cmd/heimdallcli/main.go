package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/maticnetwork/heimdall/app"
	checkpoint "github.com/maticnetwork/heimdall/checkpoint/cli"
	"github.com/maticnetwork/heimdall/helper"
	staking "github.com/maticnetwork/heimdall/staking/cli"
	"github.com/spf13/viper"
)

// rootCmd is the entry point for this binary
var (
	rootCmd = &cobra.Command{
		Use:   "heimdallcli",
		Short: "Heimdall light-client",
	}
)

func main() {
	// disable sorting
	cobra.EnableCommandSorting = false

	// get the codec
	cdc := app.MakeCodec()

	// TODO: Setup keybase, viper object, etc. to be passed into
	// the below functions and eliminate global vars, like we do
	// with the cdc.

	// add standard rpc, and tx commands
	rpc.AddCommands(rootCmd)
	rootCmd.AddCommand(client.LineBreak)
	tx.AddCommands(rootCmd, cdc)
	rootCmd.AddCommand(client.LineBreak)

	// add query/post commands (custom to binary)
	rootCmd.AddCommand(
		client.GetCommands()...,
	)
	rootCmd.AddCommand(
		client.PostCommands(
			checkpoint.GetSendCheckpointTx(cdc),
			checkpoint.GetCheckpointACKTx(cdc),
			checkpoint.GetCheckpointNoACKTx(cdc),
			staking.GetValidatorExitTx(cdc),
			staking.GetValidatorJoinTx(cdc),
			staking.GetValidatorUpdateTx(cdc),
		)...,
	)

	// add proxy, version and key info
	rootCmd.AddCommand(
		client.LineBreak,
		client.LineBreak,
		version.VersionCmd,
	)

	// bind with-heimdall-config config with root cmd
	viper.BindPFlag(
		helper.WithHeimdallConfigFlag,
		rootCmd.Flags().Lookup(helper.WithHeimdallConfigFlag),
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(rootCmd, "HD", os.ExpandEnv("$HOME/.heimdallcli"))
	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}
}
