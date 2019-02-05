package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// get validator information via address
func GetValidatorInfo(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-validator-info",
		Short: "show validator information via validator address",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			validatorAddress := viper.GetString(FlagValidatorAddress)

			res, err := cliCtx.QueryStore(common.GetValidatorKey([]byte(validatorAddress)), "staker")
			if err != nil {
				fmt.Printf("Error fetching validator information from store, Error: %v ValidatorAddr: %v", err, validatorAddress)
				return err
			}

			var _validator types.Validator
			err = cdc.UnmarshalBinary(res, &_validator)
			if err != nil {
				fmt.Printf("Error unmarshalling validator , Error: %v", err)
				return err
			}
			return nil
		},
	}

	cmd.MarkFlagRequired(FlagValidatorAddress)
	return cmd
}

// get validator information via address
func GetCurrentValSet(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-current-valset",
		Short: "show current validator set",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, err := cliCtx.QueryStore(common.CurrentValidatorSetKey, "staker")
			if err != nil {
				return err
			}
			var _validatorSet types.ValidatorSet
			err = cdc.UnmarshalBinary(res, &_validatorSet)
			if err != nil {
				return err
			}
			fmt.Printf("Current validator set : %v", _validatorSet)
			return nil
		},
	}

	return cmd
}
