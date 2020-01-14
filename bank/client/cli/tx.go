package cli

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	bankTypes "github.com/maticnetwork/heimdall/bank/types"
	hmClient "github.com/maticnetwork/heimdall/client"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        bankTypes.ModuleName,
		Short:                      "Bank transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       hmClient.ValidateCmd,
	}

	txCmd.AddCommand(
		client.PostCommands(
			SendTxCmd(cdc),
			TopupTxCmd(cdc),
			WithdrawFeeTxCmd(cdc),
		)...,
	)
	return txCmd
}

// SendTxCmd will create a send tx and sign it with the given key.
func SendTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send [to_address] [amount]",
		Short: "Send coin transfer tx",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// get account getter
			accGetter := authTypes.NewAccountRetriever(cliCtx)

			// get from account
			from := helper.GetFromAddress(cliCtx)

			// to key
			to := types.HexToHeimdallAddress(args[0])
			if to.Empty() {
				return errors.New("Invalid to address")
			}

			if err := accGetter.EnsureExists(from); err != nil {
				return err
			}

			account, err := accGetter.GetAccount(from)
			if err != nil {
				return err
			}

			// parse coins trying to be sent
			coins, err := types.ParseCoins(args[1])
			if err != nil {
				return err
			}

			// ensure account has enough coins
			if !account.GetCoins().IsAllGTE(coins) {
				return fmt.Errorf("address %s doesn't have enough coins to pay for this transaction", from)
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := bankTypes.NewMsgSend(from, to, coins)
			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}

// TopupTxCmd will create a topup tx
func TopupTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "topup",
		Short: "Topup tokens for validators",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// get proposer
			proposer := types.HexToHeimdallAddress(viper.GetString(FlagProposerAddress))
			if proposer.Empty() {
				proposer = helper.GetFromAddress(cliCtx)
			}

			validatorID := viper.GetInt64(FlagValidatorID)
			if validatorID == 0 {
				return fmt.Errorf("Validator ID cannot be zero")
			}

			txhash := viper.GetString(FlagTxHash)
			if txhash == "" {
				return fmt.Errorf("transaction hash has to be supplied")
			}

			// build and sign the transaction, then broadcast to Tendermint
			msg := bankTypes.NewMsgTopup(
				proposer,
				uint64(validatorID),
				types.HexToHeimdallHash(txhash),
				uint64(viper.GetInt64(FlagLogIndex)),
			)

			// broadcast msg with cli
			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().Int(FlagValidatorID, 0, "--validator-id=<validator ID here>")
	cmd.Flags().String(FlagTxHash, "", "--tx-hash=<transaction-hash>")
	cmd.Flags().String(FlagLogIndex, "", "--log-index=<log-index>")
	cmd.MarkFlagRequired(FlagValidatorID)
	cmd.MarkFlagRequired(FlagTxHash)
	cmd.MarkFlagRequired(FlagLogIndex)
	return cmd
}

// WithdrawFeeTxCmd will create a fee withdraw tx
func WithdrawFeeTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "withdraw-fee",
		Short: "Fee token withdrawal for validators",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// get proposer
			proposer := types.HexToHeimdallAddress(viper.GetString(FlagProposerAddress))
			if proposer.Empty() {
				proposer = helper.GetFromAddress(cliCtx)
			}

			validatorID := viper.GetInt64(FlagValidatorID)
			if validatorID == 0 {
				return fmt.Errorf("Validator ID cannot be zero")
			}

			// get msg
			msg := bankTypes.NewMsgWithdrawFee(
				proposer,
				uint64(validatorID),
			)

			// broadcast msg with cli
			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}

	cmd.Flags().StringP(FlagProposerAddress, "p", "", "--proposer=<proposer-address>")
	cmd.Flags().Int(FlagValidatorID, 0, "--validator-id=<validator ID here>")
	cmd.MarkFlagRequired(FlagValidatorID)
	cmd.MarkFlagRequired(FlagProposerAddress)
	return cmd
}
