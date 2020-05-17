package cli

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	"github.com/maticnetwork/heimdall/slashing/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	// Group slashing queries under a subcommand
	slashingQueryCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Querying commands for the slashing module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// slashingQueryCmd query command
	slashingQueryCmd.AddCommand(
		client.GetCommands(
			// GetCmdQuerySigningInfo(cdc),
			GetCmdQueryParams(cdc),
		)...,
	)
	return slashingQueryCmd

}

/* // GetCmdQuerySigningInfo implements the command to query signing info.
func GetCmdQuerySigningInfo(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "signing-info [validator-id]",
		Short: "Query a validator's signing information",
		Long: strings.TrimSpace(`Use a validators' id to find the signing-info for that validator:

$ <appcli> query slashing signing-info {valID}
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			validatorID := viper.GetUint64(FlagValidatorID)

			if validatorID == 0 {
				return fmt.Errorf("validator ID is required")
			}

			key := types.GetValidatorSigningInfoKey(hmTypes.NewValidatorID(validatorID).Bytes())

			res, _, err := cliCtx.QueryStore(key, storeName)
			if err != nil {
				return err
			}

			if len(res) == 0 {
				return fmt.Errorf("validator %s not found in slashing store", validatorID)
			}

			var signingInfo hmTypes.ValidatorSigningInfo
			signingInfo, err = hmTypes.UnmarshallValSigningInfo(types.ModuleCdc, res)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(signingInfo)
		},
	}
} */

// GetCmdQueryParams implements a command to fetch slashing parameters.
func GetCmdQueryParams(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "params",
		Short: "Query the current slashing parameters",
		Args:  cobra.NoArgs,
		Long: strings.TrimSpace(`Query genesis parameters for the slashing module:

$ <appcli> query slashing params
`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			route := fmt.Sprintf("custom/%s/parameters", types.QuerierRoute)
			res, _, err := cliCtx.QueryWithData(route, nil)
			if err != nil {
				return err
			}

			var params types.Params
			cdc.MustUnmarshalJSON(res, &params)
			return cliCtx.PrintOutput(params)
		},
	}
}
