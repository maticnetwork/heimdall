package cli

import (
	"context"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/cosmos/cosmos-sdk/version"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/maticnetwork/heimdall/x/bor/types"
	"github.com/spf13/cobra"
)

func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group gov queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(
		GetQueryParams(),
		GetQueryParam(),
		GetQuerySpan(),
		GetQuerySpanList(),
		GetQueryLatestSpan(),
		GetQueryNextSpanSeed(),
		PrepareNextSpan(),
	)
	return cmd
}

func GetQueryNextSpanSeed() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "next-span-seed",
		Short: "Query the next span seed ",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the next span seed.
Example:
$ %s query bor next-span-seed
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCmd := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(cliCmd, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(cliCmd)
			resp, err := queryClient.NextSpanSeed(cmd.Context(), &types.QueryNextSpanSeedRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetQueryLatestSpan() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "latest-span",
		Short: "Query the latest span ",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the latest span.
Example:
$ %s query bor latest-span
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCmd := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(cliCmd, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.LatestSpan(context.Background(), &types.QueryLatestSpanRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

func GetQuerySpanList() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "span-list [page] [limit]",
		Short: "Query the span list",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the span-list.
Example:
$ %s query|q bor span-list --page 1 --limit 10
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			page, err := cmd.Flags().GetUint64(FlagPage)
			if err != nil {
				return err
			}
			limit, err := cmd.Flags().GetUint64(FlagLimit)
			if err != nil {
				return err
			}
			cmdCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(cmdCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.SpanList(context.Background(), &types.QuerySpanListRequest{
				Page:  page,
				Limit: limit,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	cmd.Flags().Uint64(FlagPage, 0, "--page=1")
	cmd.Flags().Uint64(FlagLimit, 10, "--limit=10  maximum 20")

	_ = cmd.MarkFlagRequired(FlagPage)
	_ = cmd.MarkFlagRequired(FlagLimit)
	return cmd
}

func GetQuerySpan() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "span [span-id]",
		Short: "show span info with span-id",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Get the span info with span-id.
Example:
$ %s query bor span --span-id 1
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			spanId, err := cmd.Flags().GetUint64(FlagSpanId)
			if err != nil {
				return err
			}

			cliCmd := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(cliCmd, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Span(context.Background(), &types.QuerySpanRequest{
				SpanId: spanId,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	cmd.Flags().Uint64(FlagSpanId, 0, "span-id")
	return cmd
}

func GetQueryParam() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "param [param-type]",
		Short: "Query the parameters (span|sprint|producer-count|last-eth-block) of the bor process",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query the all the parameters for the bor.
Example:
$ %s query bor param --param-type span
$ %s query bor param --param-type sprint
$ %s query bor param --param-type producer-count
$ %s query bor param --param-type last-eth-block

`,
				version.AppName, version.AppName, version.AppName, version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			paramType, err := cmd.Flags().GetString(FlagParamTypes)
			if err != nil {
				return err
			}
			cmdCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(cmdCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Param(context.Background(), &types.QueryParamRequest{
				ParamsType: paramType,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	cmd.Flags().String(FlagParamTypes, "", "--param-type=<param type span|sprint|producer-count|last-eth-block >")
	_ = cmd.MarkFlagRequired(FlagParamTypes)
	return cmd
}

func GetQueryParams() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "show the current bor parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as bor parameters.

Example:
$ %s query bor params
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(cmdCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.Params(context.Background(), &types.QueryParamsRequest{})
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(resp)
		},
	}
	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}

// PrepareNextSpan query it will fetch next span
func PrepareNextSpan() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prepare-next-span",
		Args:  cobra.NoArgs,
		Short: "prepare next span ",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Prepare next span with startBlock .

Example:
$ %s query bor prepare-next-span --bor-chain-id 1 --start-block 254 --span-id 1
`,
				version.AppName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			borChainId, err := cmd.Flags().GetString(FlagBorChainId)
			if err != nil {
				return err
			}

			startBlock, err := cmd.Flags().GetUint64(FlagStartBlock)
			if err != nil {
				return err
			}

			spanId, err := cmd.Flags().GetUint64(FlagSpanId)
			if err != nil {
				return err
			}

			cmdCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadQueryCommandFlags(cmdCtx, cmd.Flags())
			if err != nil {
				return err
			}
			queryClient := types.NewQueryClient(clientCtx)
			resp, err := queryClient.PrepareNextSpan(context.Background(), &types.PrepareNextSpanRequest{
				SpanId:     spanId,
				BorChainId: borChainId,
				StartBlock: startBlock,
			})
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(resp)
		},
	}
	cmd.Flags().String(FlagBorChainId, "", "--bor-chain-id=15001")
	cmd.Flags().Uint64(FlagSpanId, 10, "--span-id=1")
	cmd.Flags().Uint64(FlagStartBlock, 10, "--start-block=10 ")

	_ = cmd.MarkFlagRequired(FlagBorChainId)
	_ = cmd.MarkFlagRequired(FlagSpanId)
	_ = cmd.MarkFlagRequired(FlagStartBlock)

	flags.AddQueryFlagsToCmd(cmd)
	return cmd
}
