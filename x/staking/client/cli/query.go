package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/x/staking/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group staking queries under a subcommand
	stakingQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	stakingQueryCmd.AddCommand(
		GetValidatorInfoCmd(),
		GetCurrentValSetCmd(),
	)

	return stakingQueryCmd
}

// GetValidatorInfo validator information via id or address
func GetValidatorInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validator-info",
		Short: "show validator information via validator id",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			validatorID := hmTypes.ValidatorID(viper.GetInt64(FlagValidatorID))
			validatorAddressStr := viper.GetString(FlagValidatorAddress)
			if validatorID == 0 && validatorAddressStr == "" {
				return fmt.Errorf("validator ID or validator address required")
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryValidatorRequest{ValidatorId: validatorID}
			res, err := queryClient.Validator(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res.Validator)
		},
	}

	cmd.Flags().Int(FlagValidatorID, 0, "--id=<validator ID here>")
	cmd.Flags().String(FlagValidatorAddress, "", "--validator=<validator address here>")
	return cmd
}

// GetCurrentValSet Queries Current ValidatorSet information
func GetCurrentValSetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "current-validator-set",
		Short: "show current validator set",
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)
			clientCtx, err := client.ReadTxCommandFlags(clientCtx, cmd.Flags())
			if err != nil {
				return err
			}

			queryClient := types.NewQueryClient(clientCtx)

			params := &types.QueryValidatorSetRequest{}
			res, err := queryClient.ValidatorSet(cmd.Context(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintOutput(res.ValidatorSet)
		},
	}

	return cmd
}
