package cli

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/spf13/cobra"

	hmClient "github.com/maticnetwork/heimdall/client"
	supplyTypes "github.com/maticnetwork/heimdall/supply/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	// Group supply queries under a subcommand
	supplyQueryCmd := &cobra.Command{
		Use:                        supplyTypes.ModuleName,
		Short:                      "Querying commands for the supply module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       hmClient.ValidateCmd,
	}

	// supply query command
	supplyQueryCmd.AddCommand(
		client.GetCommands(
			GetCmdQueryTotalSupply(cdc),
		)...,
	)

	return supplyQueryCmd
}

// GetCmdQueryTotalSupply implements the query total supply command.
func GetCmdQueryTotalSupply(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "total [denom]",
		Args:  cobra.MaximumNArgs(1),
		Short: "Query the total supply of coins of the chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			if len(args) == 0 {
				return queryTotalSupply(cliCtx, cdc)
			}
			return querySupplyOf(cliCtx, cdc, args[0])
		},
	}
}

func queryTotalSupply(cliCtx context.CLIContext, cdc *codec.Codec) error {
	params := supplyTypes.NewQueryTotalSupplyParams(1, 0) // no pagination
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return err
	}

	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", supplyTypes.QuerierRoute, supplyTypes.QueryTotalSupply), bz)
	if err != nil {
		return err
	}

	var totalSupply sdk.Coins
	err = cdc.UnmarshalJSON(res, &totalSupply)
	if err != nil {
		return err
	}

	return cliCtx.PrintOutput(totalSupply)
}

func querySupplyOf(cliCtx context.CLIContext, cdc *codec.Codec, denom string) error {
	params := supplyTypes.NewQuerySupplyOfParams(denom)
	bz, err := cdc.MarshalJSON(params)
	if err != nil {
		return err
	}

	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", supplyTypes.QuerierRoute, supplyTypes.QuerySupplyOf), bz)
	if err != nil {
		return err
	}

	var supply sdk.Int
	err = cdc.UnmarshalJSON(res, &supply)
	if err != nil {
		return err
	}

	return cliCtx.PrintOutput(supply)
}
