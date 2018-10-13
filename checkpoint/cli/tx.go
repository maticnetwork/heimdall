package cli

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/wire"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authctx "github.com/cosmos/cosmos-sdk/x/auth/client/context"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/spf13/cobra"

	"strconv"

	"github.com/basecoin/checkpoint"
	"github.com/ethereum/go-ethereum/common"
	"github.com/pkg/errors"
)

func SubmitCheckpointCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "submitCheckpoint",
		Args:  cobra.ExactArgs(3),
		Short: "submit checkpoint from matic chain",
		RunE: func(cmd *cobra.Command, args []string) error {

			//-- new beginning

			// --- end

			txCtx := authctx.NewTxContextFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithLogger(os.Stdout).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))

			roothash := args[0]
			// TODO replace these with flags
			start := args[1]
			end := args[2]
			fmt.Printf("the data we recieved is %v", roothash)

			from, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			startInt, err := strconv.Atoi(start)
			endInt, err := strconv.Atoi(end)
			fmt.Printf("from address is %v with txctx as %v", from, txCtx)
			msg := checkpoint.NewMsgCheckpointBlock(uint64(startInt), uint64(endInt), common.BytesToHash([]byte(roothash)))

			tx := checkpoint.NewBaseTx(msg)

			//
			////txBytes, err := rlp.EncodeToBytes(tx)
			////if err != nil {
			////	fmt.Printf("Error generating TXBYtes %v", err)
			////}
			//msgs := []sdk.Msg{msg}
			//
			////result, err := cliCtx.BroadcastTx(txBytes)
			////fmt.Printf("The Result is %v", result)
			//return utils.SendTx(txCtx, cliCtx, []sdk.Msg{msg})
			if err := cliCtx.EnsureAccountExists(); err != nil {
				return err
			}

			//from, err := cliCtx.GetFromAddress()
			//if err != nil {
			//	return err
			//}

			// TODO: (ref #1903) Allow for user supplied account number without
			// automatically doing a manual lookup.
			if txCtx.AccountNumber == 0 {
				accNum, err := cliCtx.GetAccountNumber(from)
				if err != nil {
					return err
				}

				txCtx = txCtx.WithAccountNumber(accNum)
			}

			// TODO: (ref #1903) Allow for user supplied account sequence without
			// automatically doing a manual lookup.
			if txCtx.Sequence == 0 {
				accSeq, err := cliCtx.GetAccountSequence(from)
				if err != nil {
					return err
				}

				txCtx = txCtx.WithSequence(accSeq)
			}

			//passphrase, err := keys.GetPassphrase(cliCtx.FromAddressName)
			//if err != nil {
			//	return err
			//}

			//msgAfterBuild, err := txCtx.Build(msgs)
			//if err != nil {
			//	return nil, err
			//}
			//keybase, err := keys.GetKeyBase()
			//if err != nil {
			//	return nil, err
			//}
			//
			//sig, pubkey, err := keybase.Sign(cliCtx.FromAddressName, passphrase, msgAfterBuild.Bytes())
			//if err != nil {
			//	return nil, err
			//}
			//
			//sigs := []auth.StdSignature{{
			//	AccountNumber: msgAfterBuild.AccountNumber,
			//	Sequence:      msgAfterBuild.Sequence,
			//	PubKey:        pubkey,
			//	Signature:     sig,
			//}}

			txBytes, err := rlp.EncodeToBytes(tx)
			if err != nil {
				fmt.Printf("Error generating TXBYtes %v", err)
			}
			fmt.Println("The txbytes are %v", txBytes)

			//-------
			node, err := cliCtx.GetNode()
			if err != nil {
				return err
			}

			res, err := node.BroadcastTxCommit(txBytes)
			if err != nil {
				return err
			}

			if res.CheckTx.Code != uint32(0) {
				return errors.Errorf("CheckTx failed: (%d) %s",
					res.CheckTx.Code, res.CheckTx.Log)
			}
			if res.DeliverTx.Code != uint32(0) {
				return errors.Errorf("DeliverTx failed: (%d) %s",
					res.DeliverTx.Code, res.DeliverTx.Log)
			}

			//--
			fmt.Printf("Res is %v", res)

			return nil

		},
	}
	return cmd
}
