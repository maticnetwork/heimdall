package cli

import (
	"errors"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	bankTypes "github.com/maticnetwork/heimdall/bank/types"
	hmClient "github.com/maticnetwork/heimdall/client"

	"github.com/maticnetwork/heimdall/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	// Group supply queries under a subcommand
	supplyQueryCmd := &cobra.Command{
		Use:                        bankTypes.ModuleName,
		Short:                      "Querying commands for the checkpoint module",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       hmClient.ValidateCmd,
	}

	// supply query command
	supplyQueryCmd.AddCommand(
		client.GetCommands(
			GetBalanceByAccountNumber(cdc),
		)...,
	)

	return supplyQueryCmd
}

func GetBalanceByAccountNumber(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "balance [address]",
		Short: "get the bank balance by account number",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			//
			addr := types.HexToHeimdallAddress(args[0])
			if addr.Empty() {
				return errors.New("Invalid account address")
			}

			params := bankTypes.NewQueryBalanceParams(addr)

			bz, err := cliCtx.Codec.MarshalJSON(params)
			if err != nil {
				return err
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", bankTypes.QuerierRoute, bankTypes.QueryBalance), bz)
			if err != nil {
				return err
			}

			// the query will return empty if there is no data for this account
			if len(res) == 0 {
				fmt.Println("No data available for this account")
				return nil
			}

			fmt.Println(string(res))
			return nil

		},
	}

	return cmd
}
