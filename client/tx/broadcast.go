// nolint
package tx

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/spf13/cobra"
	amino "github.com/tendermint/go-amino"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types/rest"
)

var logger = helper.Logger.With("module", "client/tx")

// BroadcastReq defines a tx broadcasting request.
type BroadcastReq struct {
	Tx   authTypes.StdTx `json:"tx"`
	Mode string          `json:"mode"`
}

//swagger:parameters txsBroadcast
type txsBroadcast struct {

	//Body
	//required:true
	//in:body
	Input txsBroadcastInput `json:"input"`
}

type txsBroadcastInput struct {
	Tx   StdTx  `json:"tx"`
	Mode string `json:"mode"`
}

type StdTx struct {
	Msg       interface{} `json:"msg"`
	Signature string      `json:"signature"`
	Memo      string      `json:"memo"`
}

// swagger:route POST /txs  txs txsBroadcast
// It broadcast the signed transaction to the network.

// BroadcastTxRequest implements a tx broadcasting handler that is responsible
// for broadcasting a valid and signed tx to a full node. The tx can be
// broadcasted via a sync|async|block mechanism.
func BroadcastTxRequest(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BroadcastReq

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = cliCtx.Codec.UnmarshalJSON(body, &req)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// check if msg is not nil
		if req.Tx.Msg == nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, errors.New("Invalid msg input").Error())
			return
		}

		// broadcast tx
		res, err := helper.BroadcastTx(cliCtx, req.Tx, req.Mode)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// GetBroadcastCommand returns the tx broadcast command.
func GetBroadcastCommand(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "broadcast [file_path]",
		Short: "Broadcast transactions generated offline",
		Long: strings.TrimSpace(`Broadcast transactions created with the --generate-only
flag and signed with the sign command. Read a transaction from [file_path] and
broadcast it to a node. If you supply a dash (-) argument in place of an input
filename, the command reads from standard input.

$ gaiacli tx broadcast ./mytxn.json
`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			stdTx, err := helper.ReadStdTxFromFile(cliCtx.Codec, args[0])
			if err != nil {
				return err
			}

			// broadcast tx
			res, err := helper.BroadcastTx(cliCtx, stdTx, "")
			if err != nil {
				return err
			}
			if err := cliCtx.PrintOutput(res); err != nil {
				return err
			}
			return nil
		},
	}

	return client.PostCommands(cmd)[0]
}
