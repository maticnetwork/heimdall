package cli

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/maticnetwork/heimdall/checkpoint"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	hmClient "github.com/maticnetwork/heimdall/client"
	"github.com/maticnetwork/heimdall/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	// Group supply queries under a subcommand
	supplyQueryCmd := &cobra.Command{
		Use:                        checkpointTypes.ModuleName,
		Short:                      "Querying commands for the checkpoint module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       hmClient.ValidateCmd,
	}

	// supply query command
	supplyQueryCmd.AddCommand(
		client.GetCommands(
			GetCheckpointBuffer(cdc),
			GetLastNoACK(cdc),
			GetHeaderFromIndex(cdc),
			GetCheckpointCount(cdc),
		)...,
	)

	return supplyQueryCmd
}

// GetCheckpointBuffer get checkpoint present in buffer
func GetCheckpointBuffer(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checkpoint-buffer",
		Short: "show checkpoint present in buffer",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryStore(checkpoint.BufferCheckpointKey, "checkpoint")
			if err != nil {
				return err
			}
			var _checkpoint types.CheckpointBlockHeader
			err = cdc.UnmarshalBinaryBare(res, &_checkpoint)
			if err != nil {
				fmt.Printf("Unable to unmarshall Error: %v", err)
				return err
			}
			fmt.Printf("Proposer: %v , StartBlock: %v , EndBlock: %v, Roothash: %v", _checkpoint.Proposer.String(), _checkpoint.StartBlock, _checkpoint.EndBlock, _checkpoint.RootHash.String())
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
			res, err := cliCtx.QueryStore(checkpoint.CheckpointNoACKCacheKey, "checkpoint")
			if err != nil {
				return err
			}

			fmt.Printf("LastNoACK received at %v", time.Unix(int64(binary.BigEndian.Uint64(res)), 0))
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
			headerNumber := viper.GetInt(FlagHeaderNumber)
			res, err := cliCtx.QueryStore(checkpoint.GetHeaderKey(uint64(headerNumber)), "checkpoint")
			if err != nil {
				fmt.Printf("Unable to fetch header block , Error:%v HeaderIndex:%v", err, headerNumber)
				return err
			}
			var _checkpoint types.CheckpointBlockHeader
			err = cdc.UnmarshalBinaryBare(res, &_checkpoint)
			if err != nil {
				fmt.Printf("Unable to unmarshall header block , Error:%v HeaderIndex:%v", err, headerNumber)
				return err
			}
			fmt.Printf("Proposer: %v , StartBlock: %v , EndBlock: %v, Roothash: %v", _checkpoint.Proposer.String(), _checkpoint.StartBlock, _checkpoint.EndBlock, _checkpoint.RootHash.String())

			return nil
		},
	}
	cmd.MarkFlagRequired(FlagHeaderNumber)

	return cmd
}

// GetCheckpointCount get number of checkpoint received count
func GetCheckpointCount(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "checkpoint-count",
		Short: "get checkpoint counts",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryStore(checkpoint.ACKCountKey, "staking")
			if err != nil {
				return err
			}

			ackCount, err := strconv.ParseInt(string(res), 10, 64)
			if err != nil {
				return err
			}
			fmt.Printf("Total number of checkpoint so far : %v", ackCount)
			return nil
		},
	}

	return cmd
}
