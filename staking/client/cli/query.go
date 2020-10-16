package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/heimdall/client"
	hmClient "github.com/maticnetwork/heimdall/client"
	"github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
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
		GetValidatorInfo(),
		GetCurrentValSet(),
	)

	return supplyQueryCmd
}

// GetValidatorInfo validator information via id or address
func GetValidatorInfo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-info",
		Short: "show validator information via validator id",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			validatorID := viper.GetInt64(FlagValidatorID)
			validatorAddressStr := viper.GetString(FlagValidatorAddress)
			if validatorID == 0 && validatorAddressStr == "" {
				return fmt.Errorf("validator ID or validator address required")
			}

			var queryParams []byte
			var err error
			var t string = ""
			if validatorAddressStr != "" {
				queryParams, err = clientCtx.Codec.MarshalJSON(types.NewQuerySignerParams(common.FromHex(validatorAddressStr)))
				if err != nil {
					return err
				}
				t = types.QuerySigner
			} else if validatorID != 0 {
				queryParams, err = clientCtx.Codec.MarshalJSON(types.NewQueryValidatorParams(hmTypes.ValidatorID(validatorID)))
				if err != nil {
					return err
				}
				t = types.QueryValidator
			}

			// get validator
			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, t), queryParams)
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
func GetCurrentValSet() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current-validator-set",
		Short: "show current validator set",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())

			// get validator set
			res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCurrentValidatorSet), nil)
			if err != nil {
				return err
			}

			fmt.Println(res)
			return nil
		},
	}

	return cmd
}
