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
	"strconv"
)

func SubmitCheckpointCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submitCheckpoint",
		Args:  cobra.ExactArgs(3),
		Short: "submit checkpoint from matic chain",
		RunE: func(cmd *cobra.Command, args []string) error {
			//fmt.Printf("the account thingy is this %v",authcmd.GetAccountDecoder(cdc))
			txCtx := authctx.NewTxContextFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			roothash:=args[0]
			// TODO replace these with flags
			start :=args[1]
			end := args[2]
			fmt.Printf("the data we recieved is %v",roothash)

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			startInt,err:=strconv.Atoi(start)
			endInt,err:=strconv.Atoi(end)
			fmt.Printf("from address is %v with txctx as %v",from,txCtx)
			msg := checkpoint.NewMsgCheckpointBlock(from,startInt,endInt,roothash)


			return utils.SendTx(txCtx,cliCtx,[]sdk.Msg{msg})

		},
	}
	return cmd
}
