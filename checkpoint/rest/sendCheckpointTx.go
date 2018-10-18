package rest

import (
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	authctx "github.com/cosmos/cosmos-sdk/x/auth/client/context"
	"github.com/gorilla/mux"

	"encoding/json"
	"fmt"
	"io/ioutil"

	"encoding/hex"
	"github.com/basecoin/checkpoint"
	"github.com/cosmos/cosmos-sdk/x/stake"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	r.HandleFunc(
		"/checkpoint/submitCheckpoint",
		submitCheckpointFromBridgeRequestHandlerFn(cdc, kb, cliCtx),
	).Methods("POST")
	r.HandleFunc("/stake/createValidator",
		createNewValidatorRequestHandlerFn(cdc, kb, cliCtx),
	).Methods("POST")
	// TODO add editValidator etc

}

type AddValidatorBody struct {
	Password           string `json:"password"`
	Local_account_name string `json:"local_account_name"`
	Chain_id           string `json:"chain_id"`
	Account_number     int64  `json:"account_number"`
	Sequence           int64  `json:"sequence"`
	Gas                int64  `json:"gas"`
	Moniker            string `json:"moniker"`
	Public_key         string `json:"public_key"`
	Amount             string `json:"amount"`
	Identity           string `json:"identity"`
	Website            string `json:"website"`
	Details            string `json:"details"`
	ValidatorAddress   string `json:"validator_address"`
}

func createNewValidatorRequestHandlerFn(cdc *wire.Codec, kb keys.Keybase, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m AddValidatorBody
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
		txCtx := authctx.TxContext{
			Codec:         cdc,
			ChainID:       m.Chain_id,
			AccountNumber: m.Account_number,
			Sequence:      m.Sequence,
			Gas:           m.Gas,
		}

		validatorAddress, err := sdk.AccAddressFromBech32(m.ValidatorAddress)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Couldn't decode address. Error: %s", err.Error())))
			return
		}
		description := stake.Description{
			Moniker:  m.Moniker,
			Identity: m.Identity,
			Website:  m.Website,
			Details:  m.Details,
		}
		pk, err := sdk.GetValPubKeyBech32(m.Public_key)
		if err != nil {
			fmt.Printf("Error decoding public key ")
		}
		amountInt, err := sdk.ParseCoin(m.Amount)
		if err != nil {
			fmt.Printf("Error decoding amount ")
		}

		msg := stake.NewMsgCreateValidator(validatorAddress, pk, amountInt, description)

		txBytes, err := txCtx.BuildAndSign(m.Local_account_name, m.Password, []sdk.Msg{msg})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(err.Error()))
			return
		}
		res, err := cliCtx.BroadcastTx(txBytes)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		output, err := json.MarshalIndent(res, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(output)

	}
}

//
//type CheckpointBody struct {
//	Password           string `json:"password"`
//	Local_account_name string `json:"local_account_name"`
//	Chain_id           string `json:"chain_id"`
//	Account_number     int64  `json:"account_number"`
//	Sequence           int64  `json:"sequence"`
//	Gas                int64  `json:"gas"`
//	Root_hash          string `json:"root_hash"`
//	Start_block        int64  `json:"start_block"`
//	End_block          int64  `json:"end_block"`
//	Proposer_address   string `json:"proposer_address"`
//}
//
//func submitCheckpointRequestHandlerFn(cdc *wire.Codec, kb keys.Keybase, cliCtx context.CLIContext) http.HandlerFunc {
//	return func(w http.ResponseWriter, r *http.Request) {
//
//		var m CheckpointBody
//		body, err := ioutil.ReadAll(r.Body)
//		if err != nil {
//			w.WriteHeader(http.StatusBadRequest)
//			w.Write([]byte(err.Error()))
//			return
//		}
//
//		err = json.Unmarshal(body, &m)
//		if err != nil {
//			fmt.Printf("we have error")
//			w.WriteHeader(http.StatusBadRequest)
//			w.Write([]byte(err.Error()))
//			return
//		}
//
//		//info, err := kb.Get(m.Local_account_name)
//		//if err != nil {
//		//	w.WriteHeader(http.StatusUnauthorized)
//		//	w.Write([]byte(err.Error()))
//		//	return
//		//}
//
//		txCtx := authctx.TxContext{
//			Codec:         cdc,
//			ChainID:       m.Chain_id,
//			AccountNumber: m.Account_number,
//			Sequence:      m.Sequence,
//			Gas:           m.Gas,
//		}
//
//		//proposerAddress, err := sdk.AccAddressFromBech32(m.Proposer_address)
//		//if err != nil {
//		//	w.WriteHeader(http.StatusInternalServerError)
//		//	w.Write([]byte(fmt.Sprintf("Couldn't decode address. Error: %s", err.Error())))
//		//	return
//		//}
//		msg := checkpoint.NewMsgCheckpointBlock(uint64(m.Start_block), uint64(m.End_block), common.BytesToHash([]byte(m.Root_hash)))
//		txBytes, err := txCtx.BuildAndSign(m.Local_account_name, m.Password, []sdk.Msg{msg})
//		if err != nil {
//			w.WriteHeader(http.StatusUnauthorized)
//			w.Write([]byte(err.Error()))
//			return
//		}
//		//
//
//		//TODO uncomment to send transaction
//		res, err := cliCtx.BroadcastTx(txBytes)
//		if err != nil {
//			w.WriteHeader(http.StatusInternalServerError)
//			w.Write([]byte(err.Error()))
//			return
//		}
//
//		//fmt.Printf("response is are %v \n",res)
//
//		//fmt.Printf("info is %v /n",info)
//
//		// TODO uncomment
//		output, err := json.MarshalIndent(res, "", "  ")
//		if err != nil {
//			w.WriteHeader(http.StatusInternalServerError)
//			w.Write([]byte(err.Error()))
//			return
//		}
//		//fmt.Printf("yay output is %v",output)
//		//TODO this too
//		w.Write(output)
//
//	}
//}

type CheckpointFromBridge struct {
	Root_hash        string `json:"root_hash"`
	Start_block      int64  `json:"start_block"`
	End_block        int64  `json:"end_block"`
	Proposer_address string `json:"proposer_address"`
}

func submitCheckpointFromBridgeRequestHandlerFn(cdc *wire.Codec, kb keys.Keybase, cliCtx context.CLIContext) http.HandlerFunc {
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
		//TODO replace with our own AUTH context
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
