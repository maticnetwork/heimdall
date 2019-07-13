package rest

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/helper"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc(
		"/bor/proposeSpan",
		postProposeSpanHandlerFn(cdc, cliCtx),
	).Methods("POST")
}

type (
	// ProposeSpan struct for proposing new span
	ProposeSpan struct {
		StartBlock uint64 `json:"startBlock"`
	}
)

func postProposeSpanHandlerFn(cdc *codec.Codec, cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var m ProposeSpan
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		err = json.Unmarshal(body, &m)
		if err != nil {
			RestLogger.Error("Error unmarshalling json epoch checkpoint", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
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
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		resp, err := helper.SendTendermintRequest(cliCtx, txBytes, helper.BroadcastAsync)
		if err != nil {
			RestLogger.Error("Error while sending request to Tendermint", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := json.Marshal(&resp)
		if err != nil {
			RestLogger.Error("Error while marshalling tendermint response", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cdc, result, cliCtx.Indent)
	}
}
