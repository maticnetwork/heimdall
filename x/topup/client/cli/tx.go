package cli

import (
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types/common"
	chainmanagerTypes "github.com/maticnetwork/heimdall/x/chainmanager/types"
	"github.com/maticnetwork/heimdall/x/topup/types"
	topupTypes "github.com/maticnetwork/heimdall/x/topup/types"
	"github.com/spf13/cobra"
)

// var cliLogger = helper.Logger.With("module", "topup/client/cli")

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	txCmd.AddCommand(
		TopupTxCmd(),
		WithdrawFeeTxCmd(),
	)

	// this line is used by starport scaffolding # 1

	return txCmd
}

// Fetch chain manager params
func getChainmanagerParams(clientCtx client.Context) (*chainmanagerTypes.Params, error) {
	// create query client
	queryClient := chainmanagerTypes.NewQueryClient(clientCtx)
	req := &chainmanagerTypes.QueryParamsRequest{}
	res, err := queryClient.Params(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return res.GetParams(), nil
}

// TopupTxCmd will create a topup tx
func TopupTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fee",
		Short: "Topup tokens for validators",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadTxCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}

			// get proposer
			proposerAddrStr, _ := cmd.Flags().GetString(FlagProposerAddress)
			proposer, err := sdk.AccAddressFromHex(proposerAddrStr)
			if err != nil {
				return fmt.Errorf("Invalid proposer address: %s", err)
			}
			if proposer.Empty() {
				proposer = helper.GetFromAddress(cliCtx)
			}

			// get user
			userAddrStr, _ := cmd.Flags().GetString(FlagUserAddress)
			user, err := sdk.AccAddressFromHex(userAddrStr)
			if err != nil {
				return fmt.Errorf("Invalid user address: %s", err)
			}
			if user.Empty() {
				return fmt.Errorf("user address cannot be zero")
			}

			txhash, _ := cmd.Flags().GetString(FlagTxHash)
			if txhash == "" {
				return fmt.Errorf("transaction hash has to be supplied")
			}

			logIndex, _ := cmd.Flags().GetUint64(FlagLogIndex)

			// Get contractCaller ref
			contractCallerObj, err := helper.NewContractCaller()
			if err != nil {
				return err
			}

			chainmanagerParams, err := getChainmanagerParams(cliCtx)
			if err != nil {
				return err
			}

			stakingManagerAddress, _ := sdk.AccAddressFromHex(chainmanagerParams.ChainParams.StakingManagerAddress)

			// get main tx receipt
			receipt, err := contractCallerObj.GetConfirmedTxReceipt(
				hmTypes.HexToHeimdallHash(txhash).EthHash(),
				chainmanagerParams.MainchainTxConfirmations,
			)
			if err != nil || receipt == nil {
				return errors.New("Transaction is not confirmed yet. Please wait for sometime and try again")
			}

			event, err := contractCallerObj.DecodeValidatorTopupFeesEvent(
				stakingManagerAddress,
				receipt,
				logIndex,
			)
			if err != nil {
				return err
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := topupTypes.NewMsgTopup(
				proposer,
				user,
				sdk.NewIntFromBigInt(event.Fee),
				hmTypes.HexToHeimdallHash(txhash),
				logIndex,
				receipt.BlockNumber.Uint64(),
			)

			// broadcast msg with cli
			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagTxHash, "", "--tx-hash=<transaction-hash>")
	cmd.Flags().String(FlagUserAddress, "", "--user=<user>")
	cmd.Flags().String(FlagFeeAmount, "", "--topup-amount=<topup-amount>")
	cmd.Flags().Uint64(FlagLogIndex, 0, "--log-index=<log-index>")
	cmd.Flags().Uint64(FlagBlockNumber, 0, "--block-number=<block-number>")

	_ = cmd.MarkFlagRequired(FlagTxHash)
	_ = cmd.MarkFlagRequired(FlagLogIndex)
	_ = cmd.MarkFlagRequired(FlagUserAddress)
	_ = cmd.MarkFlagRequired(FlagFeeAmount)
	_ = cmd.MarkFlagRequired(FlagBlockNumber)

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}

// WithdrawFeeTxCmd will create a fee withdraw tx
func WithdrawFeeTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw",
		Short: "Fee token withdrawal for validators",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := client.GetClientContextFromCmd(cmd)
			cliCtx, err := client.ReadTxCommandFlags(cliCtx, cmd.Flags())
			if err != nil {
				return err
			}

			// get proposer
			proposerAddrStr, _ := cmd.Flags().GetString(FlagProposerAddress)
			proposer, err := sdk.AccAddressFromHex(proposerAddrStr)
			if err != nil {
				return err
			}
			if proposer.Empty() {
				proposer = helper.GetFromAddress(cliCtx)
			}

			// withdraw amount
			amountStr, _ := cmd.Flags().GetString(FlagAmount)
			amount, ok := big.NewInt(0).SetString(amountStr, 10)
			if !ok {
				return errors.New("Invalid withdraw amount")
			}

			// get msg
			msg := topupTypes.NewMsgWithdrawFee(
				proposer,
				sdk.NewIntFromBigInt(amount),
			)
			// broadcast msg with cli

			return tx.GenerateOrBroadcastTxCLI(cliCtx, cmd.Flags(), &msg)
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagAmount, "0", "--amount=<withdraw-amount>")

	flags.AddTxFlagsToCmd(cmd)
	return cmd
}
