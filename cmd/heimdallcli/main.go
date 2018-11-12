package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/keys"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/maticnetwork/heimdall/app"
	checkpointRestCmds "github.com/maticnetwork/heimdall/checkpoint/rest"
	"github.com/maticnetwork/heimdall/helper"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
		client.GetCommands(
		//authcmd.GetAccountCmd("acc", cdc, types.GetAccountDecoder(cdc)),
		)...,
	)
	rootCmd.AddCommand(newAccountCmd())
	rootCmd.AddCommand(
		client.PostCommands()...,
	)

	// add proxy, version and key info
	rootCmd.AddCommand(
		client.LineBreak,
		checkpointRestCmds.ServeCommands(cdc),
		// TODO insert rest client
		keys.Commands(),
		client.LineBreak,
		version.VersionCmd,
	)

	// prepare and add flags
	executor := cli.PrepareMainCmd(rootCmd, "HD", os.ExpandEnv("$HOME/.heimdallcli"))
	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}
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