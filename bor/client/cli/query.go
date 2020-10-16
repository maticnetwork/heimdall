package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/maticnetwork/heimdall/bor/types"
	hmClient "github.com/maticnetwork/heimdall/client"
	"github.com/maticnetwork/heimdall/version"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group supply queries under a subcommand
	queryCmds := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the bor module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       hmClient.ValidateCmd,
	}

	// clerk query command
	queryCmds.AddCommand(
		GetSpan(),
		GetLatestSpan(),
		GetQueryParams(),
	)

	return queryCmds
}

// GetSpan get state record
func GetSpan() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "span",
		Short: "show span",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())

			spanIDStr := viper.GetString(FlagSpanId)
			if spanIDStr == "" {
				return fmt.Errorf("span id cannot be empty")
			}

			spanID, err := strconv.ParseUint(spanIDStr, 10, 64)
			if err != nil {
				return err
			}

			// get query params
			queryParams, err := clientCtx.Codec.MarshalJSON(types.NewQuerySpanParams(spanID))
			if err != nil {
				return err
			}

			// fetch span
			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySpan), queryParams)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return errors.New("Span not found")
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().Uint64(FlagSpanId, 0, "--id=<span ID here>")
	if err := cmd.MarkFlagRequired(FlagSpanId); err != nil {
		cliLogger.Error("GetSpan | MarkFlagRequired | FlagSpanId", "Error", err)
	}

	return cmd
}

// GetLatestSpan get state record
func GetLatestSpan() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "latest-span",
		Short: "show latest span",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())

			// fetch latest span
			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLatestSpan), nil)

			// fetch span
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return errors.New("Latest span not found")
			}

			fmt.Println(string(res))
			return nil
		},
	}

	return cmd
}

// GetQueryParams implements the params query command.
func GetQueryParams() *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "show the current bor parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as bor parameters.

Example:
$ %s query bor params
`,
				version.ClientName,
			),
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParams)
			bz, _, err := clientCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			err = json.Unmarshal(bz, &params)
			if err != nil {
				return err
			}
			return clientCtx.PrintOutput(params)
		},
	}
}
