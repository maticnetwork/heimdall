package cli

import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client/context"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	"github.com/cosmos/cosmos-sdk/client/utils"
	authctx "github.com/cosmos/cosmos-sdk/x/auth/client/context"
	"github.com/basecoin/sideblock"
	sdk "github.com/cosmos/cosmos-sdk/types"


	"os"
	"fmt"
)

func GetSideBlockSetterCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submitBlock",
		Args:  cobra.ExactArgs(3),
		Short: "submit block from matic chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			//fmt.Printf("the account thingy is this %v",authcmd.GetAccountDecoder(cdc))
			txCtx := authctx.NewTxContextFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			//validatorAddr, err := sdk.GetAccAddressBech32(args[0])
			//if err != nil {
			//	return err
			//}
			//
			//msg := slashing.NewMsgUnrevoke(validatorAddr)
			//blockNumber := args[3]
			//blockNumberInt := new(big.Int)
			//blockNumberInt, ok := blockNumberInt.SetString(blockNumber, 10)
			//if !ok {
			//	fmt.Println("SetString: error")
			//}

			blockHash := args[0]
			txroot := args[1]
			rroot := args[2]

			//fromStr := viper.GetString("from")
			//senderAddress,err := sdk.GetAccAddressBech32(fromStr)
			//if err!= nil {
			//	return err
			//}

			// TODO i dont know if this is neeeded ,check
			//fmt.Printf("ctx is %v",ctx)

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			msg := sideBlock.NewMsgSideBlock(from,blockHash,txroot,rroot)
			//fmt.Printf("the message is %v",msg)
			// build and sign the transaction, then broadcast to Tendermint
			//res, err := ctx.EnsureSignBuildBroadcast(ctx.FromAddressName, msg, cdc)
			//res, err := ctx.EnsureSignBuildBroadcast(ctx.FromAddressName, msg, cdc)
			//if err != nil {
			//	fmt.Printf("from the last err")
			//	return err
			//}

			fmt.Printf("yoyoyo")
			return utils.SendTx(txCtx,cliCtx,[]sdk.Msg{msg})
		},
	}
	//cmd.Flags().String("to", "", "Address to send coins")
	//cmd.Flags().String("from", "", "from address")
	return cmd
}
