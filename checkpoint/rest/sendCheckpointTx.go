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

		info, err := kb.Get("alice")
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}

		txCtx := authctx.TxContext{
			Codec:         cdc,
			ChainID:       "test-chain-ZLAVft",
			AccountNumber: 0,
			Sequence:      1,
			Gas:           1000000,
		}

		proposerAddress, err := sdk.AccAddressFromBech32("cosmosaccaddr15epdnu350uzfjau26scvazux4yj5mpag2ag6us")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't decode address. Error: %s", err.Error())))
			return
		}
		fmt.Printf( "vairable address is %v  \n",proposerAddress)
		fmt.Printf( "tx ctx is  is %v  \n",txCtx)

		msg := checkpoint.NewMsgCheckpointBlock	(sdk.AccAddress(proposerAddress),12,34,"vaibhavroothash")
		txBytes, err := txCtx.BuildAndSign("alice", "password", []sdk.Msg{msg})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		fmt.Printf("txbytes are %v \n",txBytes)
		//
		res, err := cliCtx.BroadcastTx(txBytes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		fmt.Printf("response is are %v \n",res)

		fmt.Printf("info is %v /n",info)
		output, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		fmt.Printf("yay output is %v",output)
		//
		w.Write(output)

	}
}
