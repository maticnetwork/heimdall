package rest

import (
"github.com/cosmos/cosmos-sdk/client/context"
"github.com/gorilla/mux"
"github.com/cosmos/cosmos-sdk/wire"
"github.com/cosmos/cosmos-sdk/crypto/keys"
"net/http"
	authctx "github.com/cosmos/cosmos-sdk/x/auth/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"fmt"
"io/ioutil"
"encoding/json"
	"github.com/basecoin/checkpoint"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	r.HandleFunc(
		"/checkpoint/submitCheckpoint",
		submitCheckpointRequestHandlerFn(cdc, kb, cliCtx),
	).Methods("POST")

}
type CheckpointBody struct {
	Password         string `json:"password"`
	Local_account_name string `json:"local_account_name"`
	Chain_id          string `json:"chain_id"`
	Account_number    int64  `json:"account_number"`
	Sequence         int64  `json:"sequence"`
	Gas              int64  `json:"gas"`
	Root_hash    	 string `json:"root_hash"`
	Start_block 		 int64 `json:"start_block"`
	End_block 		 int64 `json:"end_block"`
	Proposer_address  string `json:"proposer_address"`
}
func submitCheckpointRequestHandlerFn(cdc *wire.Codec, kb keys.Keybase, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var m CheckpointBody
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		err = json.Unmarshal(body, &m)
		if err != nil {
			fmt.Printf("we have error")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		fmt.Println("--------")
		fmt.Printf("the body is \n and %v",m)
		fmt.Println("--------")

		//info, err := kb.Get(m.Local_account_name)
		//if err != nil {
		//	w.WriteHeader(http.StatusUnauthorized)
		//	w.Write([]byte(err.Error()))
		//	return
		//}

		txCtx := authctx.TxContext{
			Codec:         cdc,
			ChainID:       m.Chain_id,
			AccountNumber: m.Account_number,
			Sequence:      m.Sequence,
			Gas:           m.Gas,
		}

		proposerAddress, err := sdk.AccAddressFromBech32(m.Proposer_address)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't decode address. Error: %s", err.Error())))
			return
		}
		//fmt.Printf( "vairable address is %v  \n",proposerAddress)
		//fmt.Printf( "tx ctx is  is %v  \n",txCtx)

		msg := checkpoint.NewMsgCheckpointBlock	(sdk.AccAddress(proposerAddress),int(m.Start_block),int(m.End_block),m.Root_hash)
		txBytes, err := txCtx.BuildAndSign(m.Local_account_name, m.Password, []sdk.Msg{msg})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		fmt.Printf("txbytes are %v \n",txBytes)
		//


		//TODO uncomment to send transaction
		res, err := cliCtx.BroadcastTx(txBytes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}



		//fmt.Printf("response is are %v \n",res)

		//fmt.Printf("info is %v /n",info)

		// TODO uncomment
		output, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		//fmt.Printf("yay output is %v",output)
		//TODO this too
		w.Write(output)

	}
}
