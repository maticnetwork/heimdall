package cli

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"

	"github.com/maticnetwork/heimdall/x/topup/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group topup queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetSequenceCmd(),
	)

	return cmd
}

// GetSequence validator information via id or address
func GetSequenceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "get sequence from txhash and logindex",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadTxCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}

			logIndex, err := cmd.Flags().GetUint64(FlagLogIndex)
			if err != nil {
				return err
			}

			txHashStr, err := cmd.Flags().GetString(FlagTxHash)
			if err != nil {
				return err
			}

			if txHashStr == "" {
				return fmt.Errorf("LogIndex and transaction hash required")
			}

			queryClient := types.NewQueryClient(cliCtx)
			fmt.Println("************************txHashStr*********")
			fmt.Println(txHashStr)
			fmt.Println("************************logIndex*********")
			fmt.Println(logIndex)
			fmt.Println("****************************************")

			params := &types.QuerySequenceRequest{TxHash: txHashStr, LogIndex: logIndex}

			res, err := queryClient.Sequence(context.Background(), params)
			if err != nil {
				return err
			}

			return cliCtx.PrintString(fmt.Sprint(res.Sequence))
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	cmd.Flags().String(FlagTxHash, "", "--tx-hash=<transaction-hash>")
	cmd.Flags().Uint64(FlagLogIndex, 0, "--log-index=<log-index>")

	_ = cmd.MarkFlagRequired(FlagTxHash)
	_ = cmd.MarkFlagRequired(FlagLogIndex)

	return cmd
}
