package cli

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/maticnetwork/heimdall/checkpoint/types"
	hmClient "github.com/maticnetwork/heimdall/client"
	"github.com/maticnetwork/heimdall/helper"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/version"
)

var cliLogger = helper.Logger.With("module", "checkpoint/client/cli")

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
			GetCheckpointByNumber(cdc),
			GetCheckpointCount(cdc),
			GetCheckpointLatest(cdc),
			GetCheckpointList(cdc),
			GetOverview(cdc),
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
			if err := jsoniter.ConfigFastest.Unmarshal(bz, &params); err != nil {
				// nolint: nilerr
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

			fmt.Println(string(res))
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
			if err := jsoniter.ConfigFastest.Unmarshal(res, &lastNoAck); err != nil {
				return err
			}

			fmt.Printf("LastNoACK received at %v", time.Unix(int64(lastNoAck), 0))
			return nil
		},
	}

	return cmd
}

// GetHeaderFromIndex get checkpoint given header index
func GetCheckpointByNumber(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "header",
		Short: "get checkpoint (header) by number",
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

			fmt.Println(string(res))
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
			if err := jsoniter.ConfigFastest.Unmarshal(res, &ackCount); err != nil {
				return err
			}

			fmt.Printf("Total number of checkpoint so far : %d\n", ackCount)
			return nil
		},
	}

	return cmd
}

// Temporary Checkpoint struct to store the Checkpoint ID
type CheckpointWithID struct {
	ID         uint64                  `json:"id"`
	Proposer   hmTypes.HeimdallAddress `json:"proposer"`
	StartBlock uint64                  `json:"start_block"`
	EndBlock   uint64                  `json:"end_block"`
	RootHash   hmTypes.HeimdallHash    `json:"root_hash"`
	BorChainID string                  `json:"bor_chain_id"`
	TimeStamp  uint64                  `json:"timestamp"`
}

// GetCheckpointLatest get the latest checkpoint
func GetCheckpointLatest(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checkpoint-latest",
		Short: "show the latest checkpoint",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			//
			// Get ack count
			//
			ackcountBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAckCount), nil)
			if err != nil {
				return err
			}

			if len(ackcountBytes) == 0 {
				fmt.Printf("Not found")
				return nil
			}

			var ackCount uint64
			if err := jsoniter.Unmarshal(ackcountBytes, &ackCount); err != nil {
				return err
			}

			//
			// Last checkpoint key
			//

			lastCheckpointKey := ackCount

			// get query params
			queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryCheckpointParams(lastCheckpointKey))
			if err != nil {
				return err
			}

			//
			// Get checkpoint
			//

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCheckpoint), queryParams)
			if err != nil {
				return err
			}

			var checkpointUnmarshal hmTypes.Checkpoint
			if err = jsoniter.Unmarshal(res, &checkpointUnmarshal); err != nil {
				return err
			}

			checkpointWithID := &CheckpointWithID{
				ID:         ackCount,
				Proposer:   checkpointUnmarshal.Proposer,
				StartBlock: checkpointUnmarshal.StartBlock,
				EndBlock:   checkpointUnmarshal.EndBlock,
				RootHash:   checkpointUnmarshal.RootHash,
				BorChainID: checkpointUnmarshal.BorChainID,
				TimeStamp:  checkpointUnmarshal.TimeStamp,
			}

			resWithID, err := jsoniter.Marshal(checkpointWithID)
			if err != nil {
				return err
			}

			//error if checkpoint not found
			if len(resWithID) == 0 {
				fmt.Printf("No checkpoint found")
				return nil
			}

			fmt.Println(string(resWithID))
			return nil

		},
	}

	return cmd
}

// GetLastNoACK get last no ack time
func GetCheckpointList(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checkpoint-list",
		Short: "get checkpoint list",
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

			fmt.Println(string(res))
			return nil

		},
	}

	cmd.Flags().Uint64(FlagPage, 0, "--page=<page number here>")
	cmd.Flags().Uint64(FlagLimit, 0, "--id=<limit here>")

	if err := cmd.MarkFlagRequired(FlagPage); err != nil {
		cliLogger.Error("GetCheckpointList | MarkFlagRequired | FlagPage", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagLimit); err != nil {
		cliLogger.Error("GetCheckpointList | MarkFlagRequired | FlagLimit", "Error", err)
	}

	return cmd
}

type stateDump struct {
	ACKCount         uint64               `json:"ack_count"`
	CheckpointBuffer *hmTypes.Checkpoint  `json:"checkpoint_buffer"`
	ValidatorCount   int                  `json:"validator_count"`
	ValidatorSet     hmTypes.ValidatorSet `json:"validator_set"`
	LastNoACK        time.Time            `json:"last_noack_time"`
}

// GetOverview gives the complete state dump of heimdall
func GetOverview(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "overview",
		Short: "get overview",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var ackCountInt uint64

			ackCountBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAckCount), nil)
			if err == nil {
				// check content
				if len(ackCountBytes) == 0 {
					if err = jsoniter.Unmarshal(ackCountBytes, &ackCountInt); err != nil {
						// log and ignore
						cliLogger.Error("Error while unmarshing no-ack count", "error", err.Error())
					}
				}
			}

			//
			// Checkpoint buffer
			//

			var _checkpoint *hmTypes.Checkpoint

			checkpointBufferBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCheckpointBuffer), nil)
			if err == nil {
				if len(checkpointBufferBytes) != 0 {
					_checkpoint = new(hmTypes.Checkpoint)
					if err = jsoniter.Unmarshal(checkpointBufferBytes, _checkpoint); err != nil {
						// log and ignore
						cliLogger.Error("Error while unmarshing checkpoint header", "error", err.Error())
					}
				}
			}

			//
			// Current validator set
			//

			var validatorSet hmTypes.ValidatorSet

			validatorSetBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", stakingTypes.QuerierRoute, stakingTypes.QueryCurrentValidatorSet), nil)
			if err == nil {
				if err := jsoniter.Unmarshal(validatorSetBytes, &validatorSet); err != nil {
					// log and ignore
					cliLogger.Error("Error while unmarshing validator set", "error", err.Error())
				}
			}

			// validator count
			validatorCount := len(validatorSet.Validators)

			//
			// Last no-ack
			//

			// last no ack
			var lastNoACKTime uint64

			lastNoACKBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLastNoAck), nil)
			if err == nil {
				// check content
				if len(lastNoACKBytes) == 0 {
					if err = jsoniter.Unmarshal(lastNoACKBytes, &lastNoACKTime); err != nil {
						// log and ignore
						cliLogger.Error("Error while unmarshing last no-ack time", "error", err.Error())
					}
				}
			}

			//
			// State dump
			//

			state := stateDump{
				ACKCount:         ackCountInt,
				CheckpointBuffer: _checkpoint,
				ValidatorCount:   validatorCount,
				ValidatorSet:     validatorSet,
				LastNoACK:        time.Unix(int64(lastNoACKTime), 0),
			}

			result, err := jsoniter.Marshal(state)
			if err != nil {
				return err
			}

			fmt.Println(string(result))
			return nil

		},
	}

	return cmd
}
