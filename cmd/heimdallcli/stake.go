package main

import (
	"errors"
	"math/big"

	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/heimdall/helper"
	stakingcli "github.com/maticnetwork/heimdall/staking/client/cli"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// StakeCmd stakes for a validator
func StakeCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "stake",
		Short: "Stake matic tokens for your account",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			helper.InitHeimdallConfig("")

			validatorStr := viper.GetString(stakingcli.FlagValidatorAddress)
			stakeAmountStr := viper.GetString(stakingcli.FlagAmount)
			feeAmountStr := viper.GetString(stakingcli.FlagFeeAmount)
			acceptDelegation := viper.GetBool(stakingcli.FlagAcceptDelegation)

			// validator str
			if validatorStr == "" {
				return errors.New("Validator address is required")
			}

			// stake amount
			stakeAmount, ok := big.NewInt(0).SetString(stakeAmountStr, 10)
			if !ok {
				return errors.New("Invalid stake amount")
			}

			// fee amount
			feeAmount, ok := big.NewInt(0).SetString(feeAmountStr, 10)
			if !ok {
				return errors.New("Invalid fee amount")
			}

			// contract caller
			contractCaller, err := helper.NewContractCaller()
			if err != nil {
				return err
			}

			return contractCaller.StakeFor(
				common.HexToAddress(validatorStr),
				stakeAmount,
				feeAmount,
				acceptDelegation,
			)
		},
	}

	cmd.Flags().String(stakingcli.FlagValidatorAddress, "", "--validator=<validator address here>")
	cmd.Flags().String(stakingcli.FlagAmount, "10000000000000000000", "--staked-amount=<stake amount>, if left blank it will be assigned as 10 matic tokens")
	cmd.Flags().String(stakingcli.FlagFeeAmount, "5000000000000000000", "--fee-amount=<heimdall fee amount>, if left blank will be assigned as 5 matic tokens")
	cmd.Flags().Bool(stakingcli.FlagAcceptDelegation, false, "--accept-delegation=<accept delegation>, if left blank will be assigned as false")
	return cmd
}

// ApproveCmd approves tokens for a validator
func ApproveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve",
		Short: "Approve the tokens to stake",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			helper.InitHeimdallConfig("")

			stakeAmountStr := viper.GetString(stakingcli.FlagAmount)
			feeAmountStr := viper.GetString(stakingcli.FlagFeeAmount)

			// stake amount
			stakeAmount, ok := big.NewInt(0).SetString(stakeAmountStr, 10)
			if !ok {
				return errors.New("Invalid stake amount")
			}

			// fee amount
			feeAmount, ok := big.NewInt(0).SetString(feeAmountStr, 10)
			if !ok {
				return errors.New("Invalid fee amount")
			}

			contractCaller, err := helper.NewContractCaller()
			if err != nil {
				return err
			}

			return contractCaller.ApproveTokens(stakeAmount.Add(stakeAmount, feeAmount))
		},
	}

	cmd.Flags().String(stakingcli.FlagAmount, "10000000000000000000", "--staked-amount=<stake amount>, if left blank will be assigned as 10 matic tokens")
	cmd.Flags().String(stakingcli.FlagFeeAmount, "5000000000000000000", "--fee-amount=<heimdall fee amount>, if left blank will be assigned as 5 matic tokens")
	return cmd
}
