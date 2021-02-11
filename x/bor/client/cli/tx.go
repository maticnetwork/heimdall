package cli

import (
	"context"
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/flags"

	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/maticnetwork/heimdall/x/bor/types"
	"github.com/spf13/cobra"
)

const (
	ParamSpan          = "span"
	ParamSprint        = "sprint"
	ParamProducerCount = "producer-count"
	ParamLastEthBlock  = "last-eth-block"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Bor transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		PostSendProposeSpanTx(),
	)
	return txCmd
}

// PostSendProposeSpanTx send propose span transaction
func PostSendProposeSpanTx() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "propose-span",
		Short: "send propose span tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmdCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadTxCommandFlags(cmdCtx, cmd.Flags())
			if err != nil {
				return err
			}
			borChainID, err := cmd.Flags().GetString(FlagBorChainId)
			if err != nil {
				return err
			}
			if borChainID == "" {
				return fmt.Errorf("BorChainID cannot be empty")
			}
			//
			//// get proposer
			proposerAddrStr, err := cmd.Flags().GetString(FlagProposerAddress)
			if err != nil {
				return err
			}
			proposer, err := sdk.AccAddressFromHex(proposerAddrStr)
			if err != nil {
				return fmt.Errorf("invalid proposer address: %v", err)
			}
			if proposer.Empty() {
				proposer = helper.GetFromAddress(cliCtx)
			}

			// start block
			startBlock, err := cmd.Flags().GetUint64(FlagStartBlock)
			if err != nil {
				return err
			}
			// span
			spanId, err := cmd.Flags().GetUint64(FlagSpanId)
			if err != nil {
				return err
			}
			//
			// Query data
			//
			queryClient := types.NewQueryClient(cliCtx)

			// fetch duration
			resp, err := queryClient.Param(context.Background(), &types.QueryParamRequest{
				ParamsType: ParamSpan,
			})
			if err != nil {
				return errors.New("span duration not found")
			}

			if len(resp.String()) == 0 {
				return errors.New("span duration not found")
			}

			spanDuration := resp.GetSpanDuration()

			// span seed
			nextSpanResp, errs := queryClient.NextSpanSeed(context.Background(), &types.QueryNextSpanSeedRequest{})
			if errs != nil {
				return err
			}
			seed := nextSpanResp.GetNextSpanSeed()
			//
			msg := types.NewMsgProposeSpan(
				spanId,
				proposer.String(),
				startBlock,
				startBlock+spanDuration-1,
				borChainID,
				seed,
			)
			//broadcast message
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}
	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().Uint64(FlagSpanId, 0, "--span-id=<span-id>")
	cmd.Flags().String(FlagBorChainId, "", "--bor-chain-id=<bor-chain-id>")
	cmd.Flags().Uint64(FlagStartBlock, 0, "--start-block=<start-block-number>")
	_ = cmd.MarkFlagRequired(FlagBorChainId)
	_ = cmd.MarkFlagRequired(FlagStartBlock)

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
