package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// GetValidatorInfo validator information via address
func GetValidatorInfo(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-info",
		Short: "show validator information via validator address",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			validatorID := viper.GetInt64(FlagValidatorID)
			if validatorID == 0 {
				return fmt.Errorf("validator ID cannot be 0")
			}
			signerAddr, err := cliCtx.QueryStore(staking.GetValidatorMapKey(types.NewValidatorID(uint64(validatorID)).Bytes()), "staking")
			if err != nil {
				fmt.Printf("Error fetching signer address from validator ID")
				return err
			}
			res, err := cliCtx.QueryStore(staking.GetValidatorKey(signerAddr), "staking")
			if err != nil {
				fmt.Printf("Error fetching validator information from store, Error: %v ValidatorID: %v", err, validatorID)
				return err
			}

			var _validator types.Validator
			err = cdc.UnmarshalBinaryBare(res, &_validator)
			if err != nil {
				fmt.Printf("Error unmarshalling validator , Error: %v", err)
				return err
			}
			return nil
		},
	}
	cmd.Flags().Int(FlagValidatorID, 0, "--id=<validator ID here>")
	cmd.MarkFlagRequired(FlagValidatorID)
	return cmd
}

// GetCurrentValSet validator information via address
func GetCurrentValSet(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current-validator-set",
		Short: "show current validator set",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, err := cliCtx.QueryStore(staking.CurrentValidatorSetKey, "staking")
			if err != nil {
				return err
			}
			var _validatorSet types.ValidatorSet
			err = cdc.UnmarshalBinaryBare(res, &_validatorSet)
			if err != nil {
				return err
			}
			fmt.Printf("Current validator set : %v", _validatorSet)
			return nil
		},
	}

	return cmd
}
