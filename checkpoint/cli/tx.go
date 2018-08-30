package cli


import (
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/spf13/cobra"
	"github.com/cosmos/cosmos-sdk/client/context"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authctx "github.com/cosmos/cosmos-sdk/x/auth/client/context"
	"os"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/basecoin/checkpoint"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"encoding/json"
)

func SubmitCheckpointCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submitCheckpoint",
		Args:  cobra.ExactArgs(1),
		Short: "submit checkpoint from matic chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			//fmt.Printf("the account thingy is this %v",authcmd.GetAccountDecoder(cdc))
			txCtx := authctx.NewTxContextFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			data:=args[0]
			fmt.Printf("the data we recieved is %v",data)
			out, err := json.Marshal(data)
			if err != nil {
				panic (err)
			}
			fmt.Printf("output is %v",out)

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			fmt.Printf("from address is %v with txctx as %v",from,txCtx)
			//submit :=[
			//	checkpoint.BlockHeader{BlockHash:"dsds",TxRoot:"dsdsdsds",ReceiptRoot:"dsdsdsdsdsdsdsdsd"},
			//	checkpoint.BlockHeader{BlockHash:"dsds",TxRoot:"dsdsdsds",ReceiptRoot:"dsdsdsdsdsdsdsdsd"}
			//]
			submit:=[]checkpoint.BlockHeader{}
			submit= append(submit, checkpoint.BlockHeader{BlockHash:"dsds",TxRoot:"dsdsdsds",ReceiptRoot:"dsdsdsdsdsdsdsdsd"})
			submit= append(submit, checkpoint.BlockHeader{BlockHash:"dsds",TxRoot:"dsdsdsds",ReceiptRoot:"dsdsdsdsdsdsdsdsd"})
			msg := checkpoint.NewMsgSideBlock(from,submit)


			return utils.SendTx(txCtx,cliCtx,[]sdk.Msg{msg})

		},
	}
	return cmd
}
