package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// Get all delegations from a delegator
	r.HandleFunc(
		"/checkpoint/buffer",
		CheckpointBufferHandlerFn(cdc, cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/checkpoint/{headerBlockIndex}",
		CheckpointHandlerFb(cdc, cliCtx),
	).Methods("GET")

	r.HandleFunc("checkpoint/count",
		CheckpointCountHandlerFn(cdc, cliCtx),
	).Methods("GET")

}

func CheckpointBufferHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//cliCtx.TrustNode = true
		res, err := cliCtx.QueryStore(common.BufferCheckpointKey, "checkpoint")
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// the query will return empty if there is no data in buffer
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		var _checkpoint types.CheckpointBlockHeader
		err = json.Unmarshal(res, &_checkpoint)
		if err != nil {
			RestLogger.Info("Unable to marshall")
		}

		result, err := json.Marshal(&_checkpoint)
		if err != nil {
			RestLogger.Error("Error while marshalling resposne to Json", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	}
}

func CheckpointHandlerFb(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		headerint, err := strconv.Atoi(vars["headerBlockIndex"])

		res, err := cliCtx.QueryStore(common.GetHeaderKey(uint64(headerint)), "checkpoint")
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// the query will return empty if there is no data
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		var _checkpoint types.CheckpointBlockHeader
		err = json.Unmarshal(res, &_checkpoint)
		if err != nil {
			RestLogger.Info("Unable to marshall")
		}

		result, err := json.Marshal(&_checkpoint)
		if err != nil {
			RestLogger.Error("Error while marshalling resposne to Json", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	}
}

func CheckpointCountHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, err := cliCtx.QueryStore(common.ACKCountKey, "checkpoint")
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		ackCount, err := strconv.Atoi(string(res))
		//the query will return empty if there is no data
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		result, err := json.Marshal(&ackCount)
		if err != nil {
			RestLogger.Error("Error while marshalling resposne to Json", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)

	}
}
