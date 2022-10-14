package main

import (
	"encoding/json"
	"errors"
	"math/big"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ethereum/go-ethereum/common"
	chainmanagerTypes "github.com/maticnetwork/heimdall/chainmanager/types"
	"github.com/maticnetwork/heimdall/helper"
	stakingcli "github.com/maticnetwork/heimdall/staking/client/cli"
)

var checkpointEndpoint = "/chainmanager/params"

// StakeCmd stakes for a validator
func StakeCmd(cliCtx cliContext.CLIContext) *cobra.Command {
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

			params, err := GetChainmanagerParams(cliCtx)
			if err != nil {
				return err
			}

			stakingManagerAddress := params.ChainParams.StakingManagerAddress.EthAddress()
			stakeManagerInstance, err := contractCaller.GetStakeManagerInstance(stakingManagerAddress)
			if err != nil {
				return err
			}

			return contractCaller.StakeFor(
				common.HexToAddress(validatorStr),
				stakeAmount,
				feeAmount,
				acceptDelegation,
				stakingManagerAddress,
				stakeManagerInstance,
			)
		},
	}

	cmd.Flags().String(stakingcli.FlagValidatorAddress, "", "--validator=<validator address here>")
	cmd.Flags().String(stakingcli.FlagAmount, "10000000000000000000", "--staked-amount=<stake amount>, if left blank it will be assigned as 10 matic tokens")
	cmd.Flags().String(stakingcli.FlagFeeAmount, "5000000000000000000", "--fee-amount=<heimdall fee amount>, if left blank will be assigned as 5 matic tokens")
	cmd.Flags().Bool(stakingcli.FlagAcceptDelegation, true, "--accept-delegation=<accept delegation>, if left blank will be assigned as true")
	return cmd
}

// ApproveCmd approves tokens for a validator
func ApproveCmd(cliCtx cliContext.CLIContext) *cobra.Command {
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

			params, err := GetChainmanagerParams(cliCtx)
			if err != nil {
				return err
			}

			stakingManagerAddress := params.ChainParams.StakingManagerAddress.EthAddress()
			maticTokenAddress := params.ChainParams.MaticTokenAddress.EthAddress()

			maticTokenInstance, err := contractCaller.GetMaticTokenInstance(maticTokenAddress)
			if err != nil {
				return err
			}

			return contractCaller.ApproveTokens(stakeAmount.Add(stakeAmount, feeAmount), stakingManagerAddress, maticTokenAddress, maticTokenInstance)
		},
	}

	cmd.Flags().String(stakingcli.FlagAmount, "10000000000000000000", "--staked-amount=<stake amount>, if left blank will be assigned as 10 matic tokens")
	cmd.Flags().String(stakingcli.FlagFeeAmount, "5000000000000000000", "--fee-amount=<heimdall fee amount>, if left blank will be assigned as 5 matic tokens")
	return cmd
}

// GetChainmanagerParams return configManager params
func GetChainmanagerParams(cliCtx cliContext.CLIContext) (*chainmanagerTypes.Params, error) {
	response, err := helper.FetchFromAPI(
		cliCtx,
		helper.GetHeimdallServerEndpoint(checkpointEndpoint),
	)

	if err != nil {
		return nil, err
	}

	var params chainmanagerTypes.Params
	if err := json.Unmarshal(response.Result, &params); err != nil {
		return nil, err
	}

	return &params, nil
}
