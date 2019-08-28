package cli

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/cobra"

	"github.com/maticnetwork/heimdall/auth/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// GetAccountCmd returns a query account that will display the state of the
// account at a given address.
// nolint: unparam
func GetAccountCmd(storeName string, cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account [address]",
		Short: "Query account balance",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc).WithAccountDecoder(cdc)
			accGetter := types.NewAccountRetriever(cliCtx)

			// key
			key := hmTypes.HexToHeimdallAddress(args[0])

			if err := accGetter.EnsureExists(key); err != nil {
				return err
			}

			acc, err := accGetter.GetAccount(key)
			if err != nil {
				return err
			}

			return cliCtx.PrintOutput(acc)
		},
	}
	return client.GetCommands(cmd)[0]
}
