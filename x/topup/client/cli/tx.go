package cli

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/x/topup/types"
	topupTypes "github.com/maticnetwork/heimdall/x/topup/types"
	"github.com/spf13/cobra"
)

var cliLogger = helper.Logger.With("module", "topup/client/cli")

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
		client.PostCommands(
			TopupTxCmd(cdc),
			WithdrawFeeTxCmd(cdc),
		)...,
	)

	// this line is used by starport scaffolding # 1

	return txCmd
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
				return fmt.Errorf("Invalid proposer address", err)
			}
			if proposer.Empty() {
				proposer = helper.GetFromAddress(cliCtx)
			}

			validatorID, _ := cmd.Flags().GetUint64(FlagValidatorID)
			if validatorID == 0 {
				return fmt.Errorf("Validator ID cannot be zero")
			}

			// get user
			userAddrStr, _ := cmd.Flags().GetString(FlagUserAddress)
			user, err := sdk.AccAddressFromHex(userAddrStr)
			if err != nil {
				return fmt.Errorf("Invalid user address", err)
			}
			if user.Empty() {
				return fmt.Errorf("user address cannot be zero")
			}

			// fee amount
			amountStr, _ := cmd.Flags().GetString(FlagFeeAmount)
			amount, ok := sdk.NewIntFromString(amountStr)
			if !ok {
				return errors.New("Invalid fee amount")
			}

			txhash, _ := cmd.Flags().GetString(FlagTxHash)
			if txhash == "" {
				return fmt.Errorf("transaction hash has to be supplied")
			}

			logIndex, _ := cmd.Flags().GetUint64(FlagLogIndex)
			blockNumber, _ := cmd.Flags().GetUint64(FlagBlockNumber)
			// build and sign the transaction, then broadcast to Tendermint
			msg := topupTypes.NewMsgTopup(
				proposer,
				user,
				fee,
				types.HexToHeimdallHash(txhash),
				logIndex,
				blockNumber,
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

	cmd.MarkFlagRequired(FlagTxHash)
	cmd.MarkFlagRequired(FlagLogIndex)
	cmd.MarkFlagRequired(FlagUserAddress)
	cmd.MarkFlagRequired(FlagFeeAmount)
	cmd.MarkFlagRequired(FlagBlockNumber)

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

	return cmd
}
