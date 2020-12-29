package cli

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/maticnetwork/heimdall/x/clerk/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group clerk queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetStateRecord(),
	)

	return cmd
}

// GetStateRecord get state record
func GetStateRecord() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record",
		Short: "show state record",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			recordIDStr, err := cmd.Flags().GetString(FlagRecordID)
			if err != nil {
				return err
			}
			if recordIDStr == "" {
				return fmt.Errorf("record id cannot be empty")
			}
			recordID, err := strconv.ParseUint(recordIDStr, 10, 64)
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryRecordParams{RecordId: recordID}
			res, err := queryClient.Record(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res.EventRecord)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	cmd.Flags().Uint64(FlagRecordID, 0, "--id=<record ID here>")
	cmd.MarkFlagRequired(FlagRecordID) //nolint

	return cmd
}
