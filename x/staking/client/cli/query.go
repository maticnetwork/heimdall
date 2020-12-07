package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"

	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/x/staking/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group staking queries under a subcommand
	stakingQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	stakingQueryCmd.AddCommand(
	// GetCmdQueryValidators(),
	)

	return stakingQueryCmd
}

// GetCmdQueryValidators implements the query all validators command.
// func GetCmdQueryValidators() *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "validators",
// 		Short: "Query for all validators",
// 		Args:  cobra.NoArgs,
// 		Long: strings.TrimSpace(
// 			fmt.Sprintf(`Query details about all validators on a network.

// Example:
// $ %s query staking validators
// `,
// 				version.AppName,
// 			),
// 		),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			clientCtx := client.GetClientContextFromCmd(cmd)
// 			clientCtx, err := client.ReadQueryCommandFlags(clientCtx, cmd.Flags())
// 			if err != nil {
// 				return err
// 			}

// 			queryClient := types.NewQueryClient(clientCtx)
// 			pageReq, err := client.ReadPageRequest(cmd.Flags())
// 			if err != nil {
// 				return err
// 			}

// 			result, err := queryClient.Validators(context.Background(), &types.QueryValidatorsRequest{
// 				// Leaving status empty on purpose to query all validators.
// 				Pagination: pageReq,
// 			})
// 			if err != nil {
// 				return err
// 			}

// 			return clientCtx.PrintOutput(result)
// 		},
// 	}

// 	flags.AddQueryFlagsToCmd(cmd)

// 	return cmd
// }
