package cli

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// send validator join transaction
func GetValidatorJoinTx(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-join",
		Short: "Join Heimdall as a validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			validatorStr := viper.GetString(FlagValidatorAddress)
			if validatorStr != "" {
				return fmt.Errorf("Validator address has to be supplied")
			}
			if common.IsHexAddress(validatorStr) {
				return fmt.Errorf("Invalid validator address")
			}

			pubkeyStr := viper.GetString(FlagSignerPubkey)
			if pubkeyStr != "" {
				return fmt.Errorf("Pubkey has to be supplied")
			}

			startEpoch := viper.GetInt64(FlagStartEpoch)

			endEpoch := viper.GetInt64(FlagEndEpoch)

			amountStr := viper.GetString(FlagAmount)

			validatorAddr := common.HexToAddress(validatorStr)

			pubkeyBytes, err := hex.DecodeString(pubkeyStr)
			if err != nil {
				return err
			}
			pubkey := types.NewPubKey(pubkeyBytes)

			msg := staking.NewMsgValidatorJoin(validatorAddr, pubkey, uint64(startEpoch), uint64(endEpoch), json.Number(amountStr))

			return helper.CreateAndSendTx(msg, cliCtx)
		},
	}

	cmd.Flags().String(FlagValidatorAddress, helper.GetPubKey().Address().String(), "--validator=<validator address here>")
	cmd.Flags().String(FlagSignerPubkey, "", "--signer-pubkey=<signer pubkey here>")
	cmd.Flags().String(FlagStartEpoch, "0", "--start-epoch=<start epoch of validator here>")
	cmd.Flags().String(FlagEndEpoch, "0", "--end-epoch=<end epoch of validator here>")
	cmd.Flags().String(FlagAmount, "", "--staked-amount=<staked amount>")

	cmd.MarkFlagRequired(FlagSignerPubkey)
	cmd.MarkFlagRequired(FlagStartEpoch)
	cmd.MarkFlagRequired(FlagEndEpoch)
	cmd.MarkFlagRequired(FlagAmount)
	return cmd
}

// send validator exit transaction
func GetValidatorExitTx(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-exit",
		Short: "Exit heimdall as a valdiator ",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			validatorStr := viper.GetString(FlagValidatorAddress)
			if validatorStr != "" {
				return fmt.Errorf("Validator address has to be supplied")
			}
			if common.IsHexAddress(validatorStr) {
				return fmt.Errorf("Invalid validator address")
			}
			validatorAddr := common.HexToAddress(validatorStr)
			msg := staking.NewMsgValidatorExit(validatorAddr)

			return helper.CreateAndSendTx(msg, cliCtx)
		},
	}

	cmd.Flags().String(FlagValidatorAddress, helper.GetPubKey().Address().String(), "--validator=<validator address here>")
	return cmd
}

// send validator update transaction
func GetValidatorUpdateTx(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "signer-update",
		Short: "Update signer for a validator",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			validatorStr := viper.GetString(FlagValidatorAddress)
			if validatorStr != "" {
				return fmt.Errorf("Validator address has to be supplied")
			}
			if common.IsHexAddress(validatorStr) {
				return fmt.Errorf("Invalid validator address")
			}
			pubkeyStr := viper.GetString(FlagNewSignerPubkey)
			if pubkeyStr != "" {
				return fmt.Errorf("Pubkey has to be supplied")
			}

			amountStr := viper.GetString(FlagAmount)
			validatorAddr := common.HexToAddress(validatorStr)

			pubkeyBytes, err := hex.DecodeString(pubkeyStr)
			if err != nil {
				return err
			}
			pubkey := types.NewPubKey(pubkeyBytes)

			msg := staking.NewMsgValidatorUpdate(validatorAddr, pubkey, json.Number(amountStr))

			return helper.CreateAndSendTx(msg, cliCtx)
		},
	}
	cmd.Flags().String(FlagValidatorAddress, helper.GetPubKey().Address().String(), "--validator=<validator address here>")
	cmd.Flags().String(FlagNewSignerPubkey, "", "--new-pubkey=< new signer pubkey here>")
	cmd.Flags().String(FlagAmount, "", "--staked-amount=<staked amount>")

	cmd.MarkFlagRequired(FlagNewSignerPubkey)
	cmd.MarkFlagRequired(FlagAmount)

	return cmd
}
