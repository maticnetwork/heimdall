package cli

import (
	"encoding/binary"
	"fmt"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// get checkpoint present in buffer
func GetCheckpointBuffer(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-checkpoint-buffer",
		Short: "show checkpoint present in buffer",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryStore(common.BufferCheckpointKey, "checkpoint")
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

// get last no ack time
func GetLastNoACK(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-last-noack",
		Short: "get last no ack received time",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, err := cliCtx.QueryStore(common.CheckpointNoACKCacheKey, "checkpoint")
			if err != nil {
				return err
			}

			fmt.Printf("LastNoACK received at %v", time.Unix(int64(binary.BigEndian.Uint64(res)), 0))
			return nil
		},
	}

	return cmd
}

// get checkpoint given header index
func GetHeaderFromIndex(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-last-noack",
		Short: "get last no ack received time",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			headerNumber := viper.GetInt(FlagHeaderNumber)
			res, err := cliCtx.QueryStore(common.GetHeaderKey(uint64(headerNumber)), "checkpoint")
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

// get number of checkpoint received count
func GetCheckpointCount(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-last-noack",
		Short: "get last no ack received time",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, err := cliCtx.QueryStore(common.ACKCountKey, "checkpoint")
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
