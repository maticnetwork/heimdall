package rest

import (
"github.com/cosmos/cosmos-sdk/client/context"
"github.com/gorilla/mux"
"github.com/cosmos/cosmos-sdk/wire"
"github.com/cosmos/cosmos-sdk/crypto/keys"
"net/http"
"fmt"
"io/ioutil"
"encoding/json"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	r.HandleFunc(
		"/checkpoint/submitCheckpoint",
		submitCheckpointRequestHandlerFn(cdc, kb, cliCtx),
	).Methods("POST")

}
type CheckpointBody struct {
	Password         string `json:"password"`
	LocalAccountName string `json:"local_account_name"`
	ChainID          string `json:"chain_id"`
	AccountNumber    int64  `json:"account_number"`
	Sequence         int64  `json:"sequence"`
	Gas              int64  `json:"gas"`
	RootHash    	 string `json:"root_hash"`
	StartBlock 		 int64 `json:"start_block"`
	EndBlock 		 int64 `json:"end_block"`
	ProposerAddress  string `json:"proposer"`
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

		fmt.Printf("account name is %v \n %v" ,m.LocalAccountName,m)
		fmt.Println("--------")
		fmt.Println(kb.List())
		fmt.Println("--------")
		fmt.Printf("the body is %v",r)
		fmt.Println("--------")
		info, err := kb.Get(m.LocalAccountName)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		//txCtx := authctx.TxContext{
		//	Codec:         cdc,
		//	ChainID:       m.ChainID,
		//	AccountNumber: m.AccountNumber,
		//	Sequence:      m.Sequence,
		//	Gas:           m.Gas,
		//}

		//variableAddress, err := sdk.AccAddressFromBech32(m.VariableAddr)
		//if err != nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	w.Write([]byte(fmt.Sprintf("Couldn't decode address. Error: %s", err.Error())))
		//	return
		//}
		//fmt.Println(variableAddress)
		//fmt.Println(txCtx)

		//msg := checkpoint.NewMsgCheckpointBlock	(sdk.AccAddress(m.ProposerAddress),int(m.StartBlock),int(m.EndBlock),string(m.RootHash))
		//txBytes, err := txCtx.BuildAndSign(m.LocalAccountName, m.Password, []sdk.Msg{msg})
		//if err != nil {
		//	w.WriteHeader(http.StatusUnauthorized)
		//	w.Write([]byte(err.Error()))
		//	return
		//}

		//
		//res, err := cliCtx.BroadcastTx(txBytes)
		//if err != nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	w.Write([]byte(err.Error()))
		//	return
		//}
		//
		//output, err := json.MarshalIndent(res, "", "  ")
		fmt.Printf("info is %v%v%v%v%v\n%v\n%v",body,m,info)
		//if err != nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	w.Write([]byte(err.Error()))
		//	return
		//}
		//
		//w.Write(output)

	}
}
