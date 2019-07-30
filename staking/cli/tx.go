package cli

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/types"
)

// send validator join transaction
func GetValidatorJoinTx(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-join",
		Short: "Join Heimdall as a validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			validatorID := viper.GetInt64(FlagValidatorID)
			if validatorID == 0 {
				return fmt.Errorf("Validator ID cannot be zero")
			}

			pubkeyStr := viper.GetString(FlagSignerPubkey)
			if pubkeyStr == "" {
				return fmt.Errorf("pubkey has to be supplied")
			}

			txhash := viper.GetString(FlagTxHash)
			if txhash == "" {
				return fmt.Errorf("transaction hash has to be supplied")
			}

			pubkeyBytes, err := hex.DecodeString(pubkeyStr)
			if err != nil {
				return err
			}
			pubkey := types.NewPubKey(pubkeyBytes)

			msg := staking.NewMsgValidatorJoin(uint64(validatorID), pubkey, common.HexToHash(txhash))
			// utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
			txBldr := authTypes.NewTxBuilderFromCLI().WithTxEncoder(authTypes.RLPTxEncoder(authTypes.GetPulpInstance()))
			return helper.CreateAndSendTx(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}

	cmd.Flags().Int(FlagValidatorID, 0, "--id=<validator ID here>")
	cmd.Flags().String(FlagSignerPubkey, "", "--signer-pubkey=<signer pubkey here>")
	cmd.Flags().String(FlagTxHash, "", "--tx-hash=<transaction-hash>")

	cmd.MarkFlagRequired(FlagSignerPubkey)
	cmd.MarkFlagRequired(FlagTxHash)
	return cmd
}

// send validator exit transaction
func GetValidatorExitTx(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-exit",
		Short: "Exit heimdall as a validator ",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			validator := viper.GetInt64(FlagValidatorID)
			if validator == 0 {
				return fmt.Errorf("validator ID cannot be 0")
			}
			txhash := viper.GetString(FlagTxHash)
			if txhash == "" {
				return fmt.Errorf("transaction hash has to be supplied")
			}
			msg := staking.NewMsgValidatorExit(uint64(validator), common.HexToHash(txhash))

			return helper.CreateAndSendTx(msg, cliCtx)
		},
	}

	cmd.Flags().Int(FlagValidatorID, 0, "--id=<validator ID here>")
	cmd.Flags().String(FlagTxHash, "", "--tx-hash=<transaction-hash>")
	cmd.MarkFlagRequired(FlagTxHash)

	return cmd
}

// send validator update transaction
func GetValidatorUpdateTx(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signer-update",
		Short: "Update signer for a validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			validator := viper.GetInt64(FlagValidatorID)
			if validator == 0 {
				return fmt.Errorf("validator ID cannot be 0")
			}

			pubkeyStr := viper.GetString(FlagNewSignerPubkey)
			if pubkeyStr == "" {
				return fmt.Errorf("Pubkey has to be supplied")
			}

			amountStr := viper.GetString(FlagAmount)

			pubkeyBytes, err := hex.DecodeString(pubkeyStr)
			if err != nil {
				return err
			}
			pubkey := types.NewPubKey(pubkeyBytes)

			txhash := viper.GetString(FlagTxHash)
			if txhash == "" {
				return fmt.Errorf("transaction hash has to be supplied")
			}

			msg := staking.NewMsgValidatorUpdate(uint64(validator), pubkey, json.Number(amountStr), common.HexToHash(txhash))

			return helper.CreateAndSendTx(msg, cliCtx)
		},
	}
	cmd.Flags().Int(FlagValidatorID, 0, "--id=<validator ID here>")
	cmd.Flags().String(FlagNewSignerPubkey, "", "--new-pubkey=< new signer pubkey here>")
	cmd.Flags().String(FlagAmount, "", "--staked-amount=<staked amount>")
	cmd.Flags().String(FlagTxHash, "", "--tx-hash=<transaction-hash>")
	cmd.MarkFlagRequired(FlagTxHash)
	cmd.MarkFlagRequired(FlagNewSignerPubkey)
	cmd.MarkFlagRequired(FlagAmount)

	return cmd
}
