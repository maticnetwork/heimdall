package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	hmClient "github.com/maticnetwork/heimdall/client"
	"github.com/maticnetwork/heimdall/milestone/types"
	"github.com/maticnetwork/heimdall/version"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	// Group supply queries under a subcommand
	supplyQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the milestone module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       hmClient.ValidateCmd,
	}

	// supply query command
	supplyQueryCmd.AddCommand(
		client.GetCommands(
			GetQueryParams(cdc),
			GetMilestone(cdc),
		)...,
	)

	return supplyQueryCmd
}

// GetQueryParams implements the params query command.
func GetQueryParams(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "show the current milestone parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as milestone parameters.

Example:
$ %s query milestone params
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParams)
			bz, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			if err := json.Unmarshal(bz, &params); err != nil {
				return nil
			}
			return cliCtx.PrintOutput(params)
		},
	}
}

// GetMilestone get milestone
func GetMilestone(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "milestone",
		Short: "show milestone present",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryMilestone), nil)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return errors.New("No Milestone")
			}

			fmt.Println(string(res))
			return nil
		},
	}

	return cmd
}
