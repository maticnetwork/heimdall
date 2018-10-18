package main

import (
	"github.com/basecoin/app"
	checkpointcmd "github.com/basecoin/checkpoint/cli"
	restCmds "github.com/basecoin/rest_client"
	"github.com/basecoin/staker/client/cli"
	//stakecmd "github.com/basecoin/staking/client/cli"
	"github.com/basecoin/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	bankcmd "github.com/cosmos/cosmos-sdk/x/bank/client/cli"
	ibccmd "github.com/cosmos/cosmos-sdk/x/ibc/client/cli"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"
	"os"
)

// rootCmd is the entry point for this binary
var (
	rootCmd = &cobra.Command{
		Use:   "basecli",
		Short: "Basecoin light-client",
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
		client.GetCommands(
			//stakecmd.GetCmdQueryValidator("stake", cdc),
			//stakecmd.GetCmdQueryValidators("stake", cdc),
			//stakecmd.GetCmdQueryDelegation("stake", cdc),
			//stakecmd.GetCmdQueryDelegations("stake", cdc),
			authcmd.GetAccountCmd("acc", cdc, types.GetAccountDecoder(cdc)),
		)...)

	rootCmd.AddCommand(
		client.PostCommands(
			bankcmd.SendTxCmd(cdc),
			ibccmd.IBCTransferCmd(cdc),
			ibccmd.IBCRelayCmd(cdc),
			//stakecmd.GetCmdCreateValidator(cdc),
			//stakecmd.GetCmdEditValidator(cdc),
			staker.GetCmdCreateMaticValidator(cdc),
			checkpointcmd.SubmitCheckpointCmd(cdc),
			//stakecmd.GetCmdDelegate(cdc),
			//stakecmd.GetCmdUnbond("stake", cdc),
		)...)

	// add proxy, version and key info
	rootCmd.AddCommand(
		client.LineBreak,
		//lcd.ServeCommand(cdc),
		restCmds.ServeMaticCommands(cdc),
		//TODO insert rest client
		keys.Commands(),
		client.LineBreak,
		version.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(rootCmd, "BC", os.ExpandEnv("$HOME/.basecli"))
	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}
}
