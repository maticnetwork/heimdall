package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/maticnetwork/bor/common"
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
		)...,
	)

	return supplyQueryCmd
}

// GetValidatorInfo validator information via id or address
func GetValidatorInfo(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-info",
		Short: "show validator information via validator id",
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
