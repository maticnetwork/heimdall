package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	clerkTypes "github.com/maticnetwork/heimdall/clerk/types"
	hmClient "github.com/maticnetwork/heimdall/client"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        clerkTypes.ModuleName,
		Short:                      "Checkpoint transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       hmClient.ValidateCmd,
	}

	txCmd.AddCommand(
		client.PostCommands(
			CreateNewStateRecord(cdc),
		)...,
	)
	return txCmd
}

// CreateNewStateRecord send checkpoint transaction
func CreateNewStateRecord(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record",
		Short: "new state record",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// bor chain id
			borChainIDStr := viper.GetString(FlagBorChainId)
			if borChainIDStr == "" {
				return fmt.Errorf("BorChainID cannot be empty")
			}

			borChainID, err := strconv.ParseUint(borChainIDStr, 10, 64)
			if err != nil {
				return fmt.Errorf("BorChainID cannot be empty")
			}

			// get proposer
			proposer := types.HexToHeimdallAddress(viper.GetString(FlagProposerAddress))
			if proposer.Empty() {
				proposer = helper.GetFromAddress(cliCtx)
			}

			// tx hash
			txHashStr := viper.GetString(FlagTxHash)
			if txHashStr == "" {
				return fmt.Errorf("tx hash cannot be empty")
			}

			// tx hash
			recordIDStr := viper.GetString(FlagRecordID)
			if recordIDStr == "" {
				return fmt.Errorf("record id cannot be empty")
			}

			recordID, err := strconv.ParseUint(recordIDStr, 10, 64)
			if err != nil {
				return fmt.Errorf("record id cannot be empty")
			}

			// log index
			logIndexStr := viper.GetString(FlagLogIndex)
			if logIndexStr == "" {
				return fmt.Errorf("log index cannot be empty")
			}

			logIndex, err := strconv.ParseUint(logIndexStr, 10, 64)
			if err != nil {
				return fmt.Errorf("log index cannot be empty")
			}

			// create new state record
			msg := clerkTypes.NewMsgEventRecord(
				proposer,
				types.HexToHeimdallHash(txHashStr),
				logIndex,
				recordID,
				borChainID,
			)

			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(FlagTxHash, "", "--tx-hash=<tx-hash>")
	cmd.Flags().String(FlagLogIndex, "", "--log-index=<log-index>")
	cmd.Flags().String(FlagRecordID, "", "--id=<record-id>")
	cmd.Flags().String(FlagBorChainId, "", "--bor-chain-id=<bor-chain-id>")

	if err := cmd.MarkFlagRequired(FlagProposerAddress); err != nil {
		logger.Error("CreateNewStateRecord | MarkFlagRequired | FlagProposerAddress", "Error", err)
	}
	if err := cmd.MarkFlagRequired(FlagRecordID); err != nil {
		logger.Error("CreateNewStateRecord | MarkFlagRequired | FlagRecordID", "Error", err)
	}
	if err := cmd.MarkFlagRequired(FlagTxHash); err != nil {
		logger.Error("CreateNewStateRecord | MarkFlagRequired | FlagTxHash", "Error", err)
	}
	if err := cmd.MarkFlagRequired(FlagLogIndex); err != nil {
		logger.Error("CreateNewStateRecord | MarkFlagRequired | FlagLogIndex", "Error", err)
	}
	if err := cmd.MarkFlagRequired(FlagBorChainId); err != nil {
		logger.Error("CreateNewStateRecord | MarkFlagRequired | FlagBorChainId", "Error", err)
	}

	return cmd
}
