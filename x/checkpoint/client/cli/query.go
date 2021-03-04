package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/version"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/maticnetwork/heimdall/x/checkpoint/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group checkpoint queries under a subcommand
	checkpointQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	checkpointQueryCmd.AddCommand(
		GetCmdQueryParams(),
		GetCmdQueryCheckpointBuffer(),
		GetCmdQueryLastNoACK(),
		GetCmdQueryHeaderFromIndex(),
		GetCmdQueryCheckpointCount(),
	)

	return checkpointQueryCmd
}

// GetCmdQueryParams implements the params query command.
func GetCmdQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "show the current checkpoint parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as checkpoint parameters.

Example:
$ %s query checkpoint params
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(&res.Params)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

// GetCmdQueryCheckpointBuffer implements the checkpoint buffer query command.
func GetCmdQueryCheckpointBuffer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checkpoint-buffer",
		Args:  cobra.NoArgs,
		Short: "show checkpoint present in buffer",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.CheckpointBuffer(context.Background(), &types.QueryCheckpointBufferRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res.CheckpointBuffer)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryLastNoACK get last no ack time
func GetCmdQueryLastNoACK() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last-noack",
		Short: "get last no ack received time",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			res, err := queryClient.LastNoAck(context.Background(), &types.QueryLastNoAckRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintString(fmt.Sprint(res.LastNoAck))
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryHeaderFromIndex get checkpoint given header index
func GetCmdQueryHeaderFromIndex() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "header",
		Short: "get checkpoint (header) from index",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			headerNumber, err := cmd.Flags().GetUint64(FlagHeaderNumber)
			if err != nil {
				return err
			}

			res, err := queryClient.Checkpoint(context.Background(), &types.QueryCheckpointRequest{Number: headerNumber})
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res.Checkpoint)

		},
	}

	cmd.Flags().Uint64(FlagHeaderNumber, 0, "--header=<header-number>")
	_ = cmd.MarkFlagRequired(FlagHeaderNumber)

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// GetCmdQueryCheckpointCount get number of checkpoint received count
func GetCmdQueryCheckpointCount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checkpoint-count",
		Short: "get checkpoint counts",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			res, err := queryClient.AckCount(context.Background(), &types.QueryAckCountRequest{})
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res)

		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
