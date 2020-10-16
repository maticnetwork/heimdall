package cli

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	hmClient "github.com/maticnetwork/heimdall/client"
	"github.com/maticnetwork/heimdall/helper"
	topupTypes "github.com/maticnetwork/heimdall/topup/types"
	"github.com/maticnetwork/heimdall/types"
)

var cliLogger = helper.Logger.With("module", "topup/client/cli")

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        topupTypes.ModuleName,
		Short:                      "Topup transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       hmClient.ValidateCmd,
	}

	txCmd.AddCommand(
		TopupTxCmd(),
		WithdrawFeeTxCmd(),
	)
	return txCmd
}

// TopupTxCmd will create a topup tx
func TopupTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fee",
		Short: "Topup tokens for validators",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			// get proposer
			proposer := types.HexToHeimdallAddress(viper.GetString(FlagProposerAddress))
			if proposer.Empty() {
				proposer = helper.GetFromAddress(clientCtx)
			}

			validatorID := viper.GetUint64(FlagValidatorID)
			if validatorID == 0 {
				return fmt.Errorf("Validator ID cannot be zero")
			}

			// get user
			user := types.HexToHeimdallAddress(viper.GetString(FlagUserAddress))
			if user.Empty() {
				return fmt.Errorf("user address cannot be zero")
			}

			// fee amount
			fee, ok := sdk.NewIntFromString(viper.GetString(FlagFeeAmount))
			if !ok {
				return errors.New("Invalid fee amount")
			}

			txhash := viper.GetString(FlagTxHash)
			if txhash == "" {
				return fmt.Errorf("transaction hash has to be supplied")
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := topupTypes.NewMsgTopup(
				proposer,
				user,
				fee,
				types.HexToHeimdallHash(txhash),
				viper.GetUint64(FlagLogIndex),
				viper.GetUint64(FlagBlockNumber),
			)

			// broadcast msg with cli
			return helper.BroadcastMsgsWithCLI(clientCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagTxHash, "", "--tx-hash=<transaction-hash>")
	cmd.Flags().String(FlagUserAddress, "", "--user=<user>")
	cmd.Flags().String(FlagFeeAmount, "", "--topup-amount=<topup-amount>")
	cmd.Flags().Uint64(FlagLogIndex, 0, "--log-index=<log-index>")
	cmd.Flags().Uint64(FlagBlockNumber, 0, "--block-number=<block-number>")

	if err := cmd.MarkFlagRequired(FlagTxHash); err != nil {
		cliLogger.Error("TopupTxCmd | MarkFlagRequired | FlagTxHash", "Error", err)
	}
	if err := cmd.MarkFlagRequired(FlagLogIndex); err != nil {
		cliLogger.Error("TopupTxCmd | MarkFlagRequired | FlagLogIndex", "Error", err)
	}
	if err := cmd.MarkFlagRequired(FlagUserAddress); err != nil {
		cliLogger.Error("TopupTxCmd | MarkFlagRequired | FlagUserAddress", "Error", err)
	}
	if err := cmd.MarkFlagRequired(FlagFeeAmount); err != nil {
		cliLogger.Error("TopupTxCmd | MarkFlagRequired | FlagFeeAmount", "Error", err)
	}
	if err := cmd.MarkFlagRequired(FlagBlockNumber); err != nil {
		cliLogger.Error("TopupTxCmd | MarkFlagRequired | FlagBlockNumber", "Error", err)
	}

	return cmd
}

// WithdrawFeeTxCmd will create a fee withdraw tx
func WithdrawFeeTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw",
		Short: "Fee token withdrawal for validators",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			// get proposer
			proposer := types.HexToHeimdallAddress(viper.GetString(FlagProposerAddress))
			if proposer.Empty() {
				proposer = helper.GetFromAddress(clientCtx)
			}

			// withdraw amount
			amountStr := viper.GetString(FlagAmount)

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
			return helper.BroadcastMsgsWithCLI(clientCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().String(FlagAmount, "0", "--amount=<withdraw-amount>")

	return cmd
}
