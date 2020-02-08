package main

import (
	"os"

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
			validatorStr := viper.GetString(stakingcli.FlagValidatorAddress)
			stakeAmount := viper.GetInt(stakingcli.FlagAmount)
			feeAmount := viper.GetInt(stakingcli.FlagFeeAmount)

			helper.InitHeimdallConfig(os.ExpandEnv("$HOME/.heimdalld"))
			contractCaller, err := helper.NewContractCaller()
			if err != nil {
				return err
			}

			return contractCaller.StakeFor(common.HexToAddress(validatorStr), int64(stakeAmount), int64(feeAmount))
		},
	}

	cmd.Flags().Int(stakingcli.FlagValidatorAddress, 1, "--id=<validator ID here>, if left blank will be assigned 1")
	return cmd
}

func ApproveCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "approve",
		Short: "Approve the tokens to stake",
		Args:  cobra.NoArgs,
		RunE: func(_ *cobra.Command, _ []string) error {
			stakeAmount := viper.GetInt(stakingcli.FlagAmount)
			feeAmount := viper.GetInt(stakingcli.FlagFeeAmount)
			helper.InitHeimdallConfig(os.ExpandEnv("$HOME/.heimdalld"))
			contractCaller, err := helper.NewContractCaller()
			if err != nil {
				return err
			}

			return contractCaller.ApproveTokens(int64(stakeAmount + feeAmount))
		},
	}
	return cmd
}
