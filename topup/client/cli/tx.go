package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	hmClient "github.com/maticnetwork/heimdall/client"
	"github.com/maticnetwork/heimdall/helper"
	topupTypes "github.com/maticnetwork/heimdall/topup/types"
	"github.com/maticnetwork/heimdall/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	txCmd := &cobra.Command{
		Use:                        topupTypes.ModuleName,
		Short:                      "Topup transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       hmClient.ValidateCmd,
	}

	txCmd.AddCommand(
		client.PostCommands(
			TopupTxCmd(cdc),
			WithdrawFeeTxCmd(cdc),
		)...,
	)
	return txCmd
}

// TopupTxCmd will create a topup tx
func TopupTxCmd(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fee",
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
			msg := topupTypes.NewMsgTopup(
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
		Use:   "withdraw",
		Short: "Fee token withdrawal for validators",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// get proposer
			proposer := types.HexToHeimdallAddress(viper.GetString(FlagProposerAddress))
			if proposer.Empty() {
				proposer = helper.GetFromAddress(cliCtx)
			}

			// get msg
			msg := topupTypes.NewMsgWithdrawFee(
				proposer,
			)
			// broadcast msg with cli
			return helper.BroadcastMsgsWithCLI(cliCtx, []sdk.Msg{msg})
		},
	}

	return cmd
}
