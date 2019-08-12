package cli

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// send checkpoint transaction
func GetSendCheckpointTx(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-checkpoint",
		Short: "send checkpoint to tendermint and ethereum chain ",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			ProposerStr := viper.GetString(FlagProposerAddress)
			if viper.GetString(FlagProposerAddress) == "" {
				return fmt.Errorf("proposer address cannot be empty")
			}

			//if common.IsHexAddress(ProposerStr) {
			//	return fmt.Errorf("Not valid validator address")
			//}

			StartBlockStr := viper.GetString(FlagStartBlock)
			if StartBlockStr == "" {
				return fmt.Errorf("start block cannot be empty")
			}

			EndBlockStr := viper.GetString(FlagEndBlock)
			if EndBlockStr == "" {
				return fmt.Errorf("end block cannot be empty")
			}

			RootHashStr := viper.GetString(FlagRootHash)
			if RootHashStr == "" {
				return fmt.Errorf("root hash cannot be empty")
			}

			Proposer := common.HexToAddress(ProposerStr)

			StartBlock, err := strconv.ParseUint(StartBlockStr, 10, 64)
			if err != nil {
				return err
			}

			EndBlock, err := strconv.ParseUint(EndBlockStr, 10, 64)
			if err != nil {
				return err
			}

			RootHash := common.HexToHash(RootHashStr)

			msg := checkpoint.NewMsgCheckpointBlock(
				Proposer,
				StartBlock,
				EndBlock,
				RootHash,
				uint64(time.Now().Unix()),
			)
			if resp, err := helper.CreateAndSendTx(msg, cliCtx); err != nil {
				return err
			} else {
				fmt.Printf("Sent transaction successfully %v", resp.TxHash)
				return nil
			}
		},
	}
	cmd.Flags().StringP(FlagProposerAddress, "p", helper.GetPubKey().Address().String(), "--proposer=<proposer-address>")
	cmd.Flags().String(FlagStartBlock, "", "--start-block=<start-block-number>")
	cmd.Flags().String(FlagEndBlock, "", "--end-block=<end-block-number>")
	cmd.Flags().StringP(FlagRootHash, "r", "", "--root-hash=<root-hash>")
	cmd.MarkFlagRequired(FlagStartBlock)
	cmd.MarkFlagRequired(FlagEndBlock)
	cmd.MarkFlagRequired(FlagRootHash)

	return cmd
}

// send checkpoint ack transaction
func GetCheckpointACKTx(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-ack",
		Short: "send acknowledgement for checkpoint in buffer",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			HeaderBlockStr := viper.GetString(FlagHeaderNumber)
			if HeaderBlockStr == "" {
				return fmt.Errorf("header number cannot be empty")
			}

			HeaderBlock, err := strconv.ParseUint(HeaderBlockStr, 10, 64)
			if err != nil {
				return err
			}

			msg := checkpoint.NewMsgCheckpointAck(HeaderBlock, uint64(time.Now().Unix()))
			if resp, err := helper.CreateAndSendTx(msg, cliCtx); err != nil {
				return err
			} else {
				fmt.Printf("Sent transaction successfully %v", resp.TxHash)
				return nil
			}
		},
	}

	cmd.Flags().String(FlagHeaderNumber, "", "--header=<header-index>")
	cmd.MarkFlagRequired(FlagHeaderNumber)
	return cmd
}

// send no-ack transaction
func GetCheckpointNoACKTx(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send-noack",
		Short: "send no-acknowledgement for last proposer",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			msg := checkpoint.NewMsgCheckpointNoAck(uint64(time.Now().Unix()))
			if resp, err := helper.CreateAndSendTx(msg, cliCtx); err != nil {
				return err
			} else {
				fmt.Printf("Sent transaction successfully %v", resp.TxHash)
				return nil
			}
		},
	}
	return cmd
}
