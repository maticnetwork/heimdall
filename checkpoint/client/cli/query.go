package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/client"
	hmClient "github.com/maticnetwork/heimdall/client"
	"github.com/maticnetwork/heimdall/version"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
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
		GetQueryParams(),
		GetCheckpointBuffer(),
		GetLastNoACK(),
		GetHeaderFromIndex(),
		GetCheckpointCount(),
	)

	return supplyQueryCmd
}

// GetQueryParams implements the params query command.
func GetQueryParams() *cobra.Command {
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
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err = client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParams)
			bz, _, err := clientCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			if err := json.Unmarshal(bz, &params); err != nil {
				return nil
			}
			return clientCtx.PrintOutput(params)
		},
	}
}

// GetCheckpointBuffer get checkpoint present in buffer
func GetCheckpointBuffer() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checkpoint-buffer",
		Short: "show checkpoint present in buffer",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err = client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCheckpointBuffer), nil)
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
func GetLastNoACK() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "last-noack",
		Short: "get last no ack received time",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err = client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLastNoAck), nil)
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
func GetHeaderFromIndex() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "header",
		Short: "get checkpoint (header) from index",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err = client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}
			headerNumber := viper.GetUint64(FlagHeaderNumber)

			// get query params
			queryParams, err := clientCtx.Codec.MarshalJSON(types.NewQueryCheckpointParams(headerNumber))
			if err != nil {
				return err
			}

			// fetch checkpoint
			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCheckpoint), queryParams)
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
func GetCheckpointCount() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checkpoint-count",
		Short: "get checkpoint counts",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err = client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAckCount), nil)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return errors.New("No ack count found")
			}

			var ackCount uint64
			if err := clientCtx.Codec.UnmarshalJSON(res, &ackCount); err != nil {
				return err
			}

			fmt.Printf("Total number of checkpoint so far : %v", ackCount)
			return nil
		},
	}

	return cmd
}
