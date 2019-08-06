package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/server"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/cli"
	"github.com/tendermint/tendermint/libs/common"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/app"
	bor "github.com/maticnetwork/heimdall/bor/cli"
	ck "github.com/maticnetwork/heimdall/checkpoint"
	checkpoint "github.com/maticnetwork/heimdall/checkpoint/cli"
	"github.com/maticnetwork/heimdall/helper"
	sk "github.com/maticnetwork/heimdall/staking"
	staking "github.com/maticnetwork/heimdall/staking/cli"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// rootCmd is the entry point for this binary
var (
	rootCmd = &cobra.Command{
		Use:   "heimdallcli",
		Short: "Heimdall light-client",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// initialise config
			initTendermintViperConfig(cmd)
			return nil
		},
	}
)

func initTendermintViperConfig(cmd *cobra.Command) {
	tendermintNode, _ := cmd.Flags().GetString(helper.NodeFlag)
	homeValue, _ := cmd.Flags().GetString(helper.HomeFlag)
	withHeimdallConfigValue, _ := cmd.Flags().GetString(helper.WithHeimdallConfigFlag)

	// set to viper
	viper.Set(helper.NodeFlag, tendermintNode)
	viper.Set(helper.HomeFlag, homeValue)
	viper.Set(helper.WithHeimdallConfigFlag, withHeimdallConfigValue)

	// start heimdall config
	helper.InitHeimdallConfig("")

	// set prefix
	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(hmTypes.PrefixAccAddr, hmTypes.PrefixAccPub)
	config.SetBech32PrefixForValidator(hmTypes.PrefixValAddr, hmTypes.PrefixValPub)
	config.SetBech32PrefixForConsensusNode(hmTypes.PrefixConsAddr, hmTypes.PrefixConsPub)
	config.Seal()
}

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
	//rpc.AddCommands(rootCmd)
	rootCmd.AddCommand(client.LineBreak)
	//tx.AddCommands(rootCmd, cdc)
	rootCmd.AddCommand(client.LineBreak)

	// add query/post commands (custom to binary)
	rootCmd.AddCommand(
		client.GetCommands(
			// checkpoint related cli get commands
			checkpoint.GetCheckpointBuffer(cdc),
			checkpoint.GetLastNoACK(cdc),
			checkpoint.GetHeaderFromIndex(cdc),
			checkpoint.GetCheckpointCount(cdc),

			// staking related cli get commands
			staking.GetValidatorInfo(cdc),
			staking.GetCurrentValSet(cdc),
		)...,
	)
	rootCmd.AddCommand(
		client.PostCommands(
			// checkpoint related cli post commands
			checkpoint.GetSendCheckpointTx(cdc),
			checkpoint.GetCheckpointACKTx(cdc),
			checkpoint.GetCheckpointNoACKTx(cdc),

			// staking related cli post commands
			staking.GetValidatorExitTx(cdc),
			staking.GetValidatorJoinTx(cdc),
			staking.GetValidatorUpdateTx(cdc),

			// bor related cli post commands
			bor.PostSendProposeSpanTx(cdc),
		)...,
	)

	// export cmds
	rootCmd.AddCommand(
		client.LineBreak,
		client.LineBreak,
		ExportCmd(ctx, cdc),
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
	executor := cli.PrepareMainCmd(rootCmd, "HD", os.ExpandEnv("$HOME/.heimdalld"))
	err := executor.Execute()
	if err != nil {
		// Note: Handle with #870
		panic(err)
	}
}

// Exportcmd a state dump file
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
			stored_ackcount, err := cliCtx.QueryStore(sk.ACKCountKey, "staking")
			if err != nil {
				fmt.Printf("Error retriving query")
				return err
			}

			ackCount, err := strconv.ParseInt(string(stored_ackcount), 10, 64)
			if err != nil {
				fmt.Printf("Unable to parse int. Response: %v Error: %v", stored_ackcount, err)
				return err
			}
			//
			// buffered checkpoint
			//
			var buffer_checkpoint hmTypes.CheckpointBlockHeader

			_checkpointBuffer, err := cliCtx.QueryStore(ck.BufferCheckpointKey, "checkpoint")
			if err == nil {
				if len(_checkpointBuffer) != 0 {
					err = cdc.UnmarshalBinaryBare(_checkpointBuffer, &buffer_checkpoint)
					if err != nil {
						fmt.Printf("Unable to unmarshall checkpoint present in buffer. Error: %v CheckpointBuffer: %v", err, _checkpointBuffer)
					}
				}
			} else {
				fmt.Printf("Unable to fetch checkpoint from buffer. Error: %v", err)
			}
			////
			//// Caches
			////
			storedCheckpointCache, err := cliCtx.QueryStore(ck.CheckpointCacheKey, "checkpoint")
			if err != nil {
				return err
			}
			var checkpointCache bool
			if bytes.Compare(storedCheckpointCache, ck.DefaultValue) == 0 {
				checkpointCache = true
			} else {
				checkpointCache = false
			}

			storedCheckpointACK, err := cliCtx.QueryStore(ck.CheckpointACKCacheKey, "checkpoint")
			if err != nil {
				return err
			}
			var checkpointACKCache bool
			if bytes.Compare(storedCheckpointACK, ck.DefaultValue) == 0 {
				checkpointACKCache = true
			} else {
				checkpointACKCache = false
			}
			////
			//// last no ack time
			////
			var lastNoACKTime int64
			lastNoACK, err := cliCtx.QueryStore(ck.CheckpointNoACKCacheKey, "checkpoint")
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
			storedHeaders, err := cliCtx.QuerySubspace(ck.HeaderBlockKey, "checkpoint")
			if err != nil {
				return err
			}
			for _, kv_pair := range storedHeaders {
				var checkpointHeader hmTypes.CheckpointBlockHeader
				if cdc.UnmarshalBinaryBare(kv_pair.Value, &checkpointHeader); err != nil {
					return err
				}
				headers = append(headers, checkpointHeader)
			}
			////
			//// validators
			////
			var validators []hmTypes.Validator
			storedVals, err := cliCtx.QuerySubspace(sk.ValidatorsKey, "staking")
			if err != nil {
				return err
			}
			for _, kv_pair := range storedVals {
				var hmVal hmTypes.Validator
				if cdc.UnmarshalBinaryBare(kv_pair.Value, &hmVal); err != nil {
					return err
				}
				validators = append(validators, hmVal)
			}
			////
			//// Current val set
			////
			var currentValSet hmTypes.ValidatorSet
			storedCurrValSet, err := cliCtx.QueryStore(sk.CurrentValidatorSetKey, "staking")
			if err != nil {
				return err
			}
			if err := cdc.UnmarshalBinaryBare(storedCurrValSet, &currentValSet); err != nil {
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
