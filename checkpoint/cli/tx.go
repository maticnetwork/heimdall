package cli

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/helper"
	"time"
	"github.com/spf13/viper"
	"github.com/ethereum/go-ethereum/common"
	"strconv"
	"fmt"
)

func GetSendCheckpointTx(cdc *codec.Codec) *cobra.Command  {
	cmd:=&cobra.Command{
		Use:   "send-checkpoint",
		Short: "send checkpoint to tendermint and ethereum chain ",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)


			ProposerStr := viper.GetString(FlagProposerAddress)
			StartBlockStr := viper.GetString(FlagStartBlock)
			EndBlockStr:= viper.GetString(FlagEndBlock)
			RootHashStr := viper.GetString(FlagRootHash)

			Proposer:= common.HexToAddress(ProposerStr)
			StartBlock,err := strconv.ParseUint(StartBlockStr,10,64)
			if err!=nil{
				return err
			}

			EndBlock,err:= strconv.ParseUint(EndBlockStr,10,64)
			if err!=nil{
				return err
			}

			RootHash:=common.HexToHash(RootHashStr)

			msg := checkpoint.NewMsgCheckpointBlock(
				Proposer,
				StartBlock,
				EndBlock,
				RootHash,
				uint64(time.Now().Unix()),
			)

			txBytes, err := helper.CreateTxBytes(msg)
			if err != nil {
				return err
			}

			resp, err := helper.SendTendermintRequest(cliCtx, txBytes)
			if err != nil {
				return err
			}
			fmt.Printf("Transaction sent %v",resp.Hash)
			return err
		},
	}
	return cmd
}

func GetSendCheckpointACK(cdc *codec.Codec) *cobra.Command  {
	cmd:=&cobra.Command{
		Use:   "send-ack",
		Short: "send acknowledgement for checkpoint in buffer",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			HeaderBlockStr:= viper.GetString(FlagHeaderNumber)

			HeaderBlock,err := strconv.ParseUint(HeaderBlockStr,10,64)
			if err!=nil{
				return err
			}

			msg := checkpoint.NewMsgCheckpointAck(HeaderBlock)

			txBytes, err := helper.CreateTxBytes(msg)
			if err != nil {
				return err
			}

			resp, err := helper.SendTendermintRequest(cliCtx, txBytes)
			if err != nil {
				return err
			}
			fmt.Printf("Transaction sent %v",resp.Hash)
			return err
		},
	}
	return cmd
}

func GetSendCheckpointNoACK(cdc *codec.Codec) *cobra.Command{
	cmd:=&cobra.Command{
		Use:   "send-NoACK",
		Short: "send no-acknowledgement for last proposer",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			msg := checkpoint.NewMsgCheckpointNoAck(uint64(time.Now().Unix()))

			txBytes, err := helper.CreateTxBytes(msg)
			if err != nil {
				return err
			}

			resp, err := helper.SendTendermintRequest(cliCtx, txBytes)
			if err != nil {
				return err
			}
			fmt.Printf("Transaction sent %v",resp.Hash)
			return err
		},
	}
	return cmd
}

