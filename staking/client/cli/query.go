package cli

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ethereum/go-ethereum/common"
	hmClient "github.com/maticnetwork/heimdall/client"
	"github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	// Group supply queries under a subcommand
	supplyQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the staking module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       hmClient.ValidateCmd,
	}

	// supply query command
	supplyQueryCmd.AddCommand(
		client.GetCommands(
			GetValidatorInfo(cdc),
			GetCurrentValSet(cdc),
			GetTotalStakingPower(cdc),
			GetValidatorStatus(cdc),
			GetProposer(cdc),
			GetCurrentProposer(cdc),
			IsOldTx(cdc),
		)...,
	)

	return supplyQueryCmd
}

// GetValidatorInfo validator information via id or address
func GetValidatorInfo(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-info",
		Short: "show validator information via validator id or validator address",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			validatorID := viper.GetInt64(FlagValidatorID)
			validatorAddressStr := viper.GetString(FlagValidatorAddress)
			if validatorID == 0 && validatorAddressStr == "" {
				return fmt.Errorf("validator ID or validator address required")
			}

			var queryParams []byte
			var err error
			var t string = ""
			if validatorAddressStr != "" {
				queryParams, err = cliCtx.Codec.MarshalJSON(types.NewQuerySignerParams(common.FromHex(validatorAddressStr)))
				if err != nil {
					return err
				}
				t = types.QuerySigner
			} else if validatorID != 0 {
				queryParams, err = cliCtx.Codec.MarshalJSON(types.NewQueryValidatorParams(hmTypes.ValidatorID(validatorID)))
				if err != nil {
					return err
				}
				t = types.QueryValidator
			}

			// get validator
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, t), queryParams)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().Int(FlagValidatorID, 0, "--id=<validator ID here>")
	cmd.Flags().String(FlagValidatorAddress, "", "--validator=<validator address here>")

	return cmd
}

// GetCurrentValSet validator information via address
func GetCurrentValSet(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current-validator-set",
		Short: "show current validator set",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// get validator set
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCurrentValidatorSet), nil)
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	return cmd
}

// Get total staking power
func GetTotalStakingPower(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "staking-power",
		Short: "show the current staking power",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			totalPowerBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryTotalValidatorPower), nil)
			if err != nil {
				return err
			}

			// check content
			if len(totalPowerBytes) == 0 {
				fmt.Printf("Total power not found")
				return nil
			}

			var totalPower uint64
			if err := json.Unmarshal(totalPowerBytes, &totalPower); err != nil {
				return err
			}

			res, err := json.Marshal(map[string]interface{}{"result": totalPower})
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	return cmd
}

// GetValidatorInfo validator status via address
func GetValidatorStatus(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-status",
		Short: "show validator status by validator address",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			validatorAddressStr := viper.GetString(FlagValidatorAddress)
			if validatorAddressStr == "" {
				return fmt.Errorf("validator address required")
			}

			queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQuerySignerParams(common.FromHex(validatorAddressStr)))
			if err != nil {
				return err
			}

			statusBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryValidatorStatus), queryParams)
			if err != nil {
				return err
			}

			// error if no checkpoint found
			if len(statusBytes) == 0 {
				fmt.Printf("Not Found")
				return nil
			}

			var status bool
			if err = json.Unmarshal(statusBytes, &status); err != nil {
				return err
			}

			res, err := json.Marshal(map[string]interface{}{"result": status})
			if err != nil {
				return err
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().String(FlagValidatorAddress, "", "--validator=<validator address here>")

	if err := cmd.MarkFlagRequired(FlagValidatorAddress); err != nil {
		logger.Error("GetValidatorStatus | MarkFlagRequired | FlagValidatorAddress", "Error", err)
	}

	return cmd
}

// GetValidatorInfo validator status via address
func GetProposer(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "proposer",
		Short: "show proposer info by times",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			times := viper.GetUint64(FlagTimes)

			// get query params
			queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryProposerParams(times))
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryProposer), queryParams)
			if err != nil {
				return err
			}

			// error if no checkpoint found
			if len(res) == 0 {
				fmt.Printf("No Proposer found")
				return nil
			}

			fmt.Println(string(res))
			return nil
		},
	}

	cmd.Flags().String(FlagTimes, "", "--times=<time here>")

	if err := cmd.MarkFlagRequired(FlagTimes); err != nil {
		logger.Error("GetProposer | MarkFlagRequired | FlagTimes", "Error", err)
	}

	return cmd
}

// Get Current proposer
func GetCurrentProposer(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current-proposer",
		Short: "show the current proposer",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCurrentProposer), nil)
			if err != nil {
				return err
			}

			// error if no checkpoint found
			if len(res) == 0 {
				fmt.Printf("Current Proposer not found")
				return nil
			}

			fmt.Println(string(res))
			return nil
		},
	}

	return cmd
}

// Check whether the transaction is old
func IsOldTx(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "is-old-tx",
		Short: "Check whether the transaction is old",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			// tx hash
			txHash := viper.GetString(FlagTxHash)
			if txHash == "" {
				return fmt.Errorf("tx hash cannot be empty")
			}

			// log index
			logIndex := viper.GetUint64(FlagLogIndex)

			// get query params
			queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryStakingSequenceParams(txHash, logIndex))
			if err != nil {
				return err
			}

			seqNo, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryStakingSequence), queryParams)
			if err != nil {
				return err
			}

			// error if no tx status found
			if len(seqNo) == 0 {
				fmt.Printf("false")
				return nil
			}

			res := true

			fmt.Println(res)
			return nil
		},
	}

	cmd.Flags().Uint64(FlagLogIndex, 0, "--log-index=<log index here>")
	cmd.Flags().Uint64(FlagTxHash, 0, "--tx-hash=<tx hash here>")

	if err := cmd.MarkFlagRequired(FlagLogIndex); err != nil {
		logger.Error("IsOldTx | MarkFlagRequired | FlagLogIndex", "Error", err)
	}

	if err := cmd.MarkFlagRequired(FlagTxHash); err != nil {
		logger.Error("IsOldTx | MarkFlagRequired | FlagTxHash", "Error", err)
	}

	return cmd
}
