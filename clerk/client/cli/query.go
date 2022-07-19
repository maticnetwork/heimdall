package cli

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/maticnetwork/heimdall/clerk/types"
	clerkTypes "github.com/maticnetwork/heimdall/clerk/types"
	hmClient "github.com/maticnetwork/heimdall/client"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

var logger = helper.Logger.With("module", "clerk/client/cli")

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	// Group supply queries under a subcommand
	queryCmds := &cobra.Command{
		Use:                        clerkTypes.ModuleName,
		Short:                      "Querying commands for the clerk module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       hmClient.ValidateCmd,
	}

	// clerk query command
	queryCmds.AddCommand(
		client.GetCommands(
			GetStateRecord(cdc),
		)...,
	)

	return queryCmds
}

// GetStateRecord get state record
func GetStateRecord(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record",
		Short: "show state record",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			recordIDStr := viper.GetString(FlagRecordID)
			if recordIDStr == "" {
				return fmt.Errorf("record id cannot be empty")
			}

			recordID, err := strconv.ParseUint(recordIDStr, 10, 64)
			if err != nil {
				return err
			}

			// get query params
			queryParams, err := cliCtx.Codec.MarshalJSON(clerkTypes.NewQueryRecordParams(recordID))
			if err != nil {
				return err
			}

			// fetch state reocrd
			res, _, err := cliCtx.QueryWithData(
				fmt.Sprintf("custom/%s/%s", clerkTypes.QuerierRoute, clerkTypes.QueryRecord),
				queryParams,
			)

			if err != nil {
				return err
			}

			if len(res) == 0 {
				return errors.New("Record not found")
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().Uint64(FlagRecordID, 0, "--id=<record ID here>")

	if err := cmd.MarkFlagRequired(FlagRecordID); err != nil {
		logger.Error("GetStateRecord | MarkFlagRequired | FlagRecordID", "Error", err)
	}

	return cmd
}

// GetStateRecord get state record
func IsOldTx(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "is-old-tx",
		Short: "Check whether the transaction is old",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// tx hash
			txHash := viper.GetString(FlagTxHash)
			if txHash == "" {
				return fmt.Errorf("tx hash cannot be empty")
			}

			// log index
			logIndexStr := viper.GetString(FlagLogIndex)
			if logIndexStr == "" {
				return fmt.Errorf("log index cannot be empty")
			}

			logIndex, err := strconv.ParseUint(logIndexStr, 10, 64)
			if err != nil {
				return fmt.Errorf("log index cannot be parsed")
			}

			// get query params
			queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryRecordSequenceParams(txHash, logIndex))
			if err != nil {
				return err
			}

			seqNo, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRecordSequence), queryParams)
			if err != nil {
				return err
			}

			// error if no tx status found
			if len(seqNo) == 0 {
				fmt.Printf("false")
				return nil
			}

			res := true

			fmt.Println(res)
			return nil
		},
	}

	cmd.Flags().Uint64(FlagLogIndex, 0, "--log-index=<log index here>")
	cmd.Flags().Uint64(FlagTxHash, 0, "--tx-hash=<tx hash here>")

	if err := cmd.MarkFlagRequired(FlagLogIndex); err != nil {
		logger.Error("IsOldTx | MarkFlagRequired | FlagLogIndex", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagTxHash); err != nil {
		logger.Error("IsOldTx | MarkFlagRequired | FlagTxHash", "Error", err)
	}

	return cmd
}

// GetStateRecordList g
func GetStateRecordList(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record-list",
		Short: "get records list",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			pageStr := viper.GetString(FlagPage)
			if pageStr == "" {
				return fmt.Errorf("page can't be empty")
			}

			limitStr := viper.GetString(FlagLimit)
			if limitStr == "" {
				return fmt.Errorf("limit can't be empty")
			}

			page, err := strconv.ParseUint(pageStr, 10, 64)
			if err != nil {
				return err
			}

			limit, err := strconv.ParseUint(limitStr, 10, 64)
			if err != nil {
				return err
			}

			// get query params
			queryParams, err := cliCtx.Codec.MarshalJSON(hmTypes.NewQueryPaginationParams(page, limit))
			if err != nil {
				return err
			}

			// query checkpoint
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCheckpointList), queryParams)
			if err != nil {
				return err
			}

			// check content
			if len(res) == 0 {
				fmt.Printf("checkpoint list not found")
				return nil
			}

			fmt.Printf(string(res))
			return nil

		},
	}

	cmd.Flags().Uint64(FlagPage, 0, "--page=<page number here>")
	cmd.Flags().Uint64(FlagLimit, 0, "--id=<limit here>")

	if err := cmd.MarkFlagRequired(FlagPage); err != nil {
		logger.Error("GetCheckpointList | MarkFlagRequired | FlagPage", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagLimit); err != nil {
		logger.Error("GetCheckpointList | MarkFlagRequired | FlagLimit", "Error", err)
	}

	return cmd
}
