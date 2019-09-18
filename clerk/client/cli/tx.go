package cli

import (
	"encoding/hex"
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

			// msg state record
			data, err := hex.DecodeString(viper.GetString(FlagRecordData))
			if err != nil {
				return err
			}

			// create new state record
			msg := clerkTypes.NewMsgEventRecord(
				proposer,
				types.HexToHeimdallHash(txHashStr),
				logIndex,
				recordID,
				types.HexToHeimdallAddress(viper.GetString(FlagRecordContract)),
				data,
			)

			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagTxHash, "", "--tx-hash=<tx-hash>")
	cmd.Flags().String(FlagLogIndex, "", "--log-index=<log-index>")
	cmd.Flags().String(FlagRecordID, "", "--id=<record-id>")
	cmd.Flags().String(FlagRecordData, "", "--data=<record-data>")
	cmd.Flags().String(FlagRecordContract, "", "--contract=<record-contract>")
	cmd.MarkFlagRequired(FlagProposerAddress)
	cmd.MarkFlagRequired(FlagRecordID)
	cmd.MarkFlagRequired(FlagRecordContract)
	cmd.MarkFlagRequired(FlagRecordData)
	cmd.MarkFlagRequired(FlagTxHash)
	cmd.MarkFlagRequired(FlagLogIndex)

	return cmd
}
