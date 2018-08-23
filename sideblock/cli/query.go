package cli

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client/context"
	"fmt"
)

func GetSideBlockGetterCmd(storeName string, cdc *wire.Codec)  *cobra.Command {
	cmd := &cobra.Command{
		Use:   "GetBlock",
		Short: "Query for matic block",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			key := args[0]
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			res, err := cliCtx.QueryStore([]byte(key), storeName)
			if err != nil {
				fmt.Printf("error is %v",err)
				return err
			}
			//sideblockObj := new(sideBlock.MsgSideBlock)
			//cdc.MustUnmarshalBinary(res,sideblockObj)

			fmt.Printf("response from query is %v",string(res))
			//validator := new(stake.Validator)
			//cdc.MustUnmarshalBinary(res, validator)
			//
			//switch viper.Get(cli.OutputFlag) {
			//case "text":
			//	human, err := validator.HumanReadableString()
			//	if err != nil {
			//		return err
			//	}
			//	fmt.Println(human)
			//
			//case "json":
			//	// parse out the validator
			//	output, err := wire.MarshalJSONIndent(cdc, validator)
			//	if err != nil {
			//		return err
			//	}
			//	fmt.Println(string(output))
			//}
			// TODO output with proofs / machine parseable etc.
			return nil
		},
	}

	return cmd
}