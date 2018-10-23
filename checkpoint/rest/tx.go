package rest

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/checkpoint"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	r.HandleFunc(
		"/checkpoint/new",
		newCheckpointHandler(cdc, kb, cliCtx),
	).Methods("POST")
}

type CheckpointFromBridge struct {
	Root_hash        string `json:"root_hash"`
	Start_block      int64  `json:"start_block"`
	End_block        int64  `json:"end_block"`
	Proposer_address string `json:"proposer_address"`
}

func newCheckpointHandler(cdc *wire.Codec, kb keys.Keybase, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m CheckpointFromBridge
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
		//TODO add proposer address
		msg := checkpoint.NewMsgCheckpointBlock(uint64(m.Start_block), uint64(m.End_block), common.HexToHash(m.Root_hash), m.Proposer_address)

		tx := checkpoint.NewBaseTx(msg)
		txBytes, err := rlp.EncodeToBytes(tx)
		if err != nil {
			fmt.Printf("Error generating TXBYtes %v", err)
		}
		fmt.Printf("The tx bytes are %v ", hex.EncodeToString(txBytes))

		client := &http.Client{}
		// TODO replace with our own AUTH context
		req, _ := http.NewRequest("GET", "http://localhost:26657/broadcast_tx_commit", nil)
		q := req.URL.Query()
		q.Add("tx", "0x"+hex.EncodeToString(txBytes))
		req.URL.RawQuery = q.Encode()

		resp, err := client.Do(req)
		fmt.Printf("The result is %v", resp)
		var bodyString string
		if resp.StatusCode == http.StatusOK {
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			bodyString = string(bodyBytes)
		}
		w.Write([]byte(bodyString))
	}
}
