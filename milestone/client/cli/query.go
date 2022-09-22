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
	"github.com/spf13/viper"

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
			GetLatestMilestone(cdc),
			GetMilestoneByNumber(cdc),
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
func GetLatestMilestone(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "latest",
		Short: "show latest milestone",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLatestMilestone), nil)
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

// GetHeaderFromIndex get checkpoint given header index
func GetMilestoneByNumber(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "",
		Short: "get milestone by number",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			number := viper.GetUint64(FlagMilestoneNumber)

			// get query params
			queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryMilestoneParams(number))
			if err != nil {
				return err
			}

			// fetch checkpoint
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryMilestoneByNumber), queryParams)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().Uint64(FlagMilestoneNumber, 0, "--number=<milesstone-number>")

	if err := cmd.MarkFlagRequired(FlagMilestoneNumber); err != nil {
		logger.Error("GetMilestoneByNumber | MarkFlagRequired | FlagMilestoneNumber", "Error", err)
	}

	return cmd
}
