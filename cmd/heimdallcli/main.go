package main

import (
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/rpc"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/cli"

	"encoding/json"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	"github.com/maticnetwork/heimdall/app"
	checkpoint "github.com/maticnetwork/heimdall/checkpoint/cli"
	hmcmn "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	staking "github.com/maticnetwork/heimdall/staking/cli"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/common"
	tmTypes "github.com/tendermint/tendermint/types"

	"bytes"
	"path/filepath"
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
	ctx := server.NewDefaultContext()

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
			checkpoint.GetCheckpointBuffer(cdc),
			checkpoint.GetLastNoACK(cdc),
			checkpoint.GetHeaderFromIndex(cdc),
			checkpoint.GetCheckpointCount(cdc),
			staking.GetValidatorInfo(cdc),
			staking.GetCurrentValSet(cdc),
		)...,
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
	rootCmd.AddCommand(ExportCmd(ctx, cdc))

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

// Exports a state dump file
func ExportCmd(ctx *server.Context, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export-heimdall",
		Short: "Export genesis file with state-dump",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {

			cliCtx := context.NewCLIContext().WithCodec(cdc)
			config := ctx.Config
			config.SetRoot(viper.GetString(cli.HomeFlag))

			// create chain id
			chainID := viper.GetString(client.FlagChainID)
			if chainID == "" {
				chainID = fmt.Sprintf("heimdall-%v", common.RandStr(6))
			}
			//
			// ack count
			//
			stored_ackcount, err := cliCtx.QueryStore(hmcmn.ACKCountKey, "checkpoint")
			if err != nil {
				fmt.Printf("Error retriving query")
				return err
			}

			ackCount, err := strconv.ParseInt(string(stored_ackcount), 10, 64)
			if err != nil {
				fmt.Printf("Unable to parse int", "Response", stored_ackcount, "Error", err)
				return err
			}
			//
			// buffered checkpoint
			//
			var buffer_checkpoint hmTypes.CheckpointBlockHeader

			_checkpointBuffer, err := cliCtx.QueryStore(hmcmn.BufferCheckpointKey, "checkpoint")
			if err == nil {
				if len(_checkpointBuffer) != 0 {
					err = cdc.UnmarshalBinary(_checkpointBuffer, &buffer_checkpoint)
					if err != nil {
						fmt.Printf("Unable to unmarshall checkpoint present in buffer", "Error", err, "CheckpointBuffer", _checkpointBuffer)
					}
				}
			} else {
				fmt.Printf("Unable to fetch checkpoint from buffer", "Error", err)
			}
			////
			//// Caches
			////
			storedCheckpointCache, err := cliCtx.QueryStore(hmcmn.CheckpointCacheKey, "checkpoint")
			if err != nil {
				return err
			}
			var checkpointCache bool
			if bytes.Compare(storedCheckpointCache, hmcmn.DefaultValue) == 0 {
				checkpointCache = true
			} else {
				checkpointCache = false
			}

			storedCheckpointACK, err := cliCtx.QueryStore(hmcmn.CheckpointACKCacheKey, "checkpoint")
			if err != nil {
				return err
			}
			var checkpointACKCache bool
			if bytes.Compare(storedCheckpointACK, hmcmn.DefaultValue) == 0 {
				checkpointACKCache = true
			} else {
				checkpointACKCache = false
			}
			////
			//// last no ack time
			////
			var lastNoACKTime int64
			lastNoACK, err := cliCtx.QueryStore(hmcmn.CheckpointNoACKCacheKey, "checkpoint")
			if err == nil && len(lastNoACK) != 0 {
				lastNoACKTime, err = strconv.ParseInt(string(lastNoACK), 10, 64)
				if err != nil {
					return err
				}
			}
			////
			//// Headers
			////
			var headers []hmTypes.CheckpointBlockHeader
			storedHeaders, err := cliCtx.QuerySubspace(hmcmn.HeaderBlockKey, "checkpoint")
			if err != nil {
				return err
			}
			for _, kv_pair := range storedHeaders {
				var checkpointHeader hmTypes.CheckpointBlockHeader
				if cdc.UnmarshalBinary(kv_pair.Value, &checkpointHeader); err != nil {
					return err
				}
				headers = append(headers, checkpointHeader)
			}
			////
			//// validators
			////
			var validators []hmTypes.Validator
			storedVals, err := cliCtx.QuerySubspace(hmcmn.ValidatorsKey, "staker")
			if err != nil {
				return err
			}
			for _, kv_pair := range storedVals {
				var hmVal hmTypes.Validator
				if cdc.UnmarshalBinary(kv_pair.Value, &hmVal); err != nil {
					return err
				}
				validators = append(validators, hmVal)
			}
			////
			//// Current val set
			////
			var currentValSet hmTypes.ValidatorSet
			storedCurrValSet, err := cliCtx.QueryStore(hmcmn.CurrentValidatorSetKey, "staker")
			if err != nil {
				return err
			}
			if err := cdc.UnmarshalBinary(storedCurrValSet, &currentValSet); err != nil {
				return err
			}

			// create genesis state
			appState := &app.GenesisState{
				Validators:         validators,
				AckCount:           uint64(ackCount),
				BufferedCheckpoint: buffer_checkpoint,
				CheckpointCache:    checkpointCache,
				CheckpointACKCache: checkpointACKCache,
				LastNoACK:          uint64(lastNoACKTime),
				Headers:            headers,
				CurrentValSet:      currentValSet,
			}

			appStateJSON, err := json.Marshal(appState)
			if err != nil {
				return err
			}

			toPrint := struct {
				ChainID string `json:"chain_id"`
			}{
				chainID,
			}

			out, err := codec.MarshalJSONIndent(cdc, toPrint)
			if err != nil {
				return err
			}

			fmt.Fprintf(os.Stderr, "%s\n", string(out))
			return writeGenesisFile(rootify("config/dump-genesis.json", config.RootDir), chainID, appStateJSON)

			return nil
		},
	}
	cmd.Flags().String(cli.HomeFlag, helper.DefaultNodeHome, "node's home directory")
	cmd.Flags().String(helper.FlagClientHome, helper.DefaultCLIHome, "client's home directory")
	cmd.Flags().String(client.FlagChainID, "", "genesis file chain-id, if left blank will be randomly created")
	return cmd
}

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

func rootify(path, root string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, path)
}
