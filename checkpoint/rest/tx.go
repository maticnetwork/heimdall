package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/helper"
	"time"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(
		"/checkpoint/new",
		newCheckpointHandler(cliCtx),
	).Methods("POST")
	r.HandleFunc("/checkpoint/ack", NewCheckpointACKHandler(cliCtx)).Methods("POST")
	r.HandleFunc("/checkpoint/no-ack", NewCheckpointNoACKHandler(cliCtx)).Methods("POST")
}

// HeaderBlock struct for incoming checkpoint
type HeaderBlock struct {
	Proposer   common.Address `json:"proposer"`
	RootHash   common.Hash    `json:"rootHash"`
	StartBlock uint64         `json:"startBlock"`
	EndBlock   uint64         `json:"endBlock"`
}

func newCheckpointHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m HeaderBlock

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		err = json.Unmarshal(body, &m)
		if err != nil {
			RestLogger.Error("Error unmarshalling json epoch checkpoint", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		msg := checkpoint.NewMsgCheckpointBlock(
			m.Proposer,
			m.StartBlock,
			m.EndBlock,
			m.RootHash,
			uint64(time.Now().Unix()),
		)

		txBytes, err := helper.CreateTxBytes(msg)
		if err != nil {
			RestLogger.Error("Unable to create txBytes", "proposer", m.Proposer.Hex(), "endBlock", m.EndBlock, "startBlock", m.StartBlock, "rootHash", m.RootHash.Hex())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		resp, err := helper.SendTendermintRequest(cliCtx, txBytes)
		if err != nil {
			RestLogger.Error("Error while sending request to Tendermint", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		result, err := json.Marshal(&resp)
		if err != nil {
			RestLogger.Error("Error while marshalling tendermint response", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(result)
	}
}

type HeaderACK struct {
	HeaderBlock uint64 `json:"headerBlock"`
}

func NewCheckpointACKHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		var m HeaderACK
		err = json.Unmarshal(body, &m)
		if err != nil {
			RestLogger.Error("Error unmarshalling Header ACK", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		// create new msg checkpoint ack
		msg := checkpoint.NewMsgCheckpointAck(m.HeaderBlock)

		txBytes, err := helper.CreateTxBytes(msg)
		if err != nil {
			RestLogger.Error("Unable to create txBytes", "error", err, "headerBlock", m.HeaderBlock)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		resp, err := helper.SendTendermintRequest(cliCtx, txBytes)
		if err != nil {
			RestLogger.Error("Error while sending request to Tendermint", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		result, err := json.Marshal(&resp)
		if err != nil {
			RestLogger.Error("Error while marshalling tendermint response", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(result)
	}
}

//
//type HeaderNoACK struct {
//	TimeStamp uint64 `json:"timestamp"`
//}

func NewCheckpointNoACKHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//body, err := ioutil.ReadAll(r.Body)
		//if err != nil {
		//	w.WriteHeader(http.StatusBadRequest)
		//	w.Write([]byte(err.Error()))
		//	return
		//}

		//var m HeaderNoACK
		//err = json.Unmarshal(body, &m)
		//if err != nil {
		//	RestLogger.Error("Error unmarshalling Header No-ACK", "error", err)
		//	w.WriteHeader(http.StatusBadRequest)
		//	w.Write([]byte(err.Error()))
		//	return
		//}

		// create new msg checkpoint ack
		msg := checkpoint.NewMsgCheckpointNoAck(uint64(time.Now().Unix()))

		txBytes, err := helper.CreateTxBytes(msg)
		if err != nil {
			RestLogger.Error("Unable to create txBytes", "error", err, "timestamp", time.Now().Unix())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		resp, err := helper.SendTendermintRequest(cliCtx, txBytes)
		if err != nil {
			RestLogger.Error("Error while sending request to Tendermint", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		result, err := json.Marshal(&resp)
		if err != nil {
			RestLogger.Error("Error while marshalling tendermint response", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(result)
	}
}
