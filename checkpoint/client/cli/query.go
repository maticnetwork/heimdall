package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/maticnetwork/heimdall/checkpoint/types"
	hmClient "github.com/maticnetwork/heimdall/client"
	"github.com/maticnetwork/heimdall/version"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	// Group supply queries under a subcommand
	supplyQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the checkpoint module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       hmClient.ValidateCmd,
	}

	// supply query command
	supplyQueryCmd.AddCommand(
		client.GetCommands(
			GetQueryParams(cdc),
			GetCheckpointBuffer(cdc),
			GetLastNoACK(cdc),
			GetHeaderFromIndex(cdc),
			GetCheckpointCount(cdc),
		)...,
	)

	return supplyQueryCmd
}

// GetQueryParams implements the params query command.
func GetQueryParams(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Args:  cobra.NoArgs,
		Short: "show the current checkpoint parameters information",
		Long: strings.TrimSpace(
			fmt.Sprintf(`Query values set as checkpoint parameters.

Example:
$ %s query checkpoint params
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

// GetCheckpointBuffer get checkpoint present in buffer
func GetCheckpointBuffer(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checkpoint-buffer",
		Short: "show checkpoint present in buffer",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCheckpointBuffer), nil)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return errors.New("No checkpoint buffer found")
			}

			fmt.Printf(string(res))
			return nil
		},
	}

	return cmd
}

// GetLastNoACK get last no ack time
func GetLastNoACK(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last-noack",
		Short: "get last no ack received time",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLastNoAck), nil)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return errors.New("No last-no-ack count found")
			}

			var lastNoAck uint64
			if err := json.Unmarshal(res, &lastNoAck); err != nil {
				return err
			}

			fmt.Printf("LastNoACK received at %v", time.Unix(int64(lastNoAck), 0))
			return nil
		},
	}

	return cmd
}

// GetHeaderFromIndex get checkpoint given header index
func GetHeaderFromIndex(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "header",
		Short: "get checkpoint (header) from index",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			headerNumber := viper.GetUint64(FlagHeaderNumber)

			// get query params
			queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryCheckpointParams(headerNumber))
			if err != nil {
				return err
			}

			// fetch checkpoint
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCheckpoint), queryParams)
			if err != nil {
				return err
			}

			fmt.Printf(string(res))
			return nil
		},
	}

	cmd.Flags().Uint64(FlagHeaderNumber, 0, "--header=<header-number>")
	if err := cmd.MarkFlagRequired(FlagHeaderNumber); err != nil {
		logger.Error("GetHeaderFromIndex | MarkFlagRequired | FlagHeaderNumber", "Error", err)
	}

	return cmd
}

// GetCheckpointCount get number of checkpoint received count
func GetCheckpointCount(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checkpoint-count",
		Short: "get checkpoint counts",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAckCount), nil)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return errors.New("No ack count found")
			}

			var ackCount uint64
			if err := json.Unmarshal(res, &ackCount); err != nil {
				return err
			}

			fmt.Printf("Total number of checkpoint so far : %d\n", ackCount)
			return nil
		},
	}

	return cmd
}
