package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"

	"encoding/hex"
	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// Get all delegations from a delegator
	r.HandleFunc(
		"/checkpoint/buffer",
		checkpointBufferHandlerFn(cdc, cliCtx),
	).Methods("GET")

	r.HandleFunc("/checkpoint/count",
		checkpointCountHandlerFn(cdc, cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/checkpoint/headers/{headerBlockIndex}",
		checkpointHeaderHandlerFn(cdc, cliCtx),
	).Methods("GET")
	r.HandleFunc("/checkpoint/{start}/{end}", checkpointHandlerFn()).Methods("GET")
}

func checkpointBufferHandlerFn(
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
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(result)
	}
}

func checkpointCountHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryStore(common.ACKCountKey, "checkpoint")
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// The query will return empty if there is no data
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		ackCount, err := strconv.ParseUint(string(res), 10, 64)
		result, err := json.Marshal(map[string]interface{}{"result": ackCount})
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

func checkpointHeaderHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		headerNumber, err := strconv.ParseUint(vars["headerBlockIndex"], 10, 64)
		res, err := cliCtx.QueryStore(common.GetHeaderKey(uint64(headerNumber)), "checkpoint")
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// the query will return empty if there is no data
		if len(res) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		var _checkpoint types.CheckpointBlockHeader
		err = json.Unmarshal(res, &_checkpoint)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
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

func checkpointHandlerFn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		helper.InitHeimdallConfig()
		start, err := strconv.Atoi(vars["start"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		end, err := strconv.Atoi(vars["end"])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		roothash, err := checkpoint.GetHeaders(uint64(start), uint64(end))
		if err != nil {
			RestLogger.Error("Unable to get header", "Start", start, "End", end, "Error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		result, err := json.Marshal(hex.EncodeToString(roothash))
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
