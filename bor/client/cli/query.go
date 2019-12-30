package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/maticnetwork/heimdall/bor/types"
	hmClient "github.com/maticnetwork/heimdall/client"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
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
		client.GetCommands(
			GetSpan(cdc),
			GetLatestSpan(cdc),
		)...,
	)

	return queryCmds
}

// GetSpan get state record
func GetSpan(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "span",
		Short: "show span",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			spanIDStr := viper.GetString(FlagSpanId)
			if spanIDStr == "" {
				return fmt.Errorf("span id cannot be empty")
			}

			spanID, err := strconv.ParseUint(spanIDStr, 10, 64)
			if err != nil {
				return err
			}

			// get query params
			queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQuerySpanParams(spanID))
			if err != nil {
				return err
			}

			// fetch span
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySpan), queryParams)
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
	cmd.MarkFlagRequired(FlagSpanId)

	return cmd
}

// GetLatestSpan get state record
func GetLatestSpan(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "latest-span",
		Short: "show latest span",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// fetch latest span
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLatestSpan), nil)

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
