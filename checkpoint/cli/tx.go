package cli

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/maticnetwork/heimdall/checkpoint"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	hmClient "github.com/maticnetwork/heimdall/client"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        checkpointTypes.ModuleName,
		Short:                      "Checkpoint transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       hmClient.ValidateCmd,
	}

	txCmd.AddCommand(
		client.PostCommands(
			SendCheckpointTx(cdc),
			SendCheckpointACKTx(cdc),
			SendCheckpointNoACKTx(cdc),
		)...,
	)
	return txCmd
}

// SendCheckpointTx send checkpoint transaction
func SendCheckpointTx(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-checkpoint",
		Short: "send checkpoint to tendermint and ethereum chain ",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// get proposer
			proposer := types.HexToHeimdallAddress(viper.GetString(FlagProposerAddress))
			if proposer.Empty() {
				proposer = helper.GetFromAddress(cliCtx)
			}

			//	start block

			startBlockStr := viper.GetString(FlagStartBlock)
			if startBlockStr == "" {
				return fmt.Errorf("start block cannot be empty")
			}

			startBlock, err := strconv.ParseUint(startBlockStr, 10, 64)
			if err != nil {
				return err
			}

			//	end block

			endBlockStr := viper.GetString(FlagEndBlock)
			if endBlockStr == "" {
				return fmt.Errorf("end block cannot be empty")
			}

			endBlock, err := strconv.ParseUint(endBlockStr, 10, 64)
			if err != nil {
				return err
			}

			// root hash

			rootHashStr := viper.GetString(FlagRootHash)
			if rootHashStr == "" {
				return fmt.Errorf("root hash cannot be empty")
			}

			msg := checkpoint.NewMsgCheckpointBlock(
				proposer,
				startBlock,
				endBlock,
				common.HexToHash(rootHashStr),
				uint64(time.Now().Unix()),
			)

			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagStartBlock, "", "--start-block=<start-block-number>")
	cmd.Flags().String(FlagEndBlock, "", "--end-block=<end-block-number>")
	cmd.Flags().StringP(FlagRootHash, "r", "", "--root-hash=<root-hash>")
	cmd.MarkFlagRequired(FlagStartBlock)
	cmd.MarkFlagRequired(FlagEndBlock)
	cmd.MarkFlagRequired(FlagRootHash)

	return cmd
}

// SendCheckpointACKTx send checkpoint ack transaction
func SendCheckpointACKTx(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-ack",
		Short: "send acknowledgement for checkpoint in buffer",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// get proposer
			proposer := types.HexToHeimdallAddress(viper.GetString(FlagProposerAddress))
			if proposer.Empty() {
				proposer = helper.GetFromAddress(cliCtx)
			}

			headerBlockStr := viper.GetString(FlagHeaderNumber)
			if headerBlockStr == "" {
				return fmt.Errorf("header number cannot be empty")
			}

			headerBlock, err := strconv.ParseUint(headerBlockStr, 10, 64)
			if err != nil {
				return err
			}

			// new checkpoint
			msg := checkpoint.NewMsgCheckpointAck(proposer, headerBlock)

			// msg
			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagHeaderNumber, "", "--header=<header-index>")
	cmd.MarkFlagRequired(FlagHeaderNumber)
	return cmd
}

// SendCheckpointNoACKTx send no-ack transaction
func SendCheckpointNoACKTx(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-noack",
		Short: "send no-acknowledgement for last proposer",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// get proposer
			proposer := types.HexToHeimdallAddress(viper.GetString(FlagProposerAddress))
			if proposer.Empty() {
				proposer = helper.GetFromAddress(cliCtx)
			}

			// create new checkpoint no-ack
			msg := checkpoint.NewMsgCheckpointNoAck(
				proposer,
				uint64(time.Now().Unix()),
			)

			// broadcast messages
			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	return cmd
}
