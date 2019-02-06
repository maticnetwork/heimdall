package rest

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"

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

	r.HandleFunc("/checkpoint/{start}/{end}", checkpointHandlerFn(cdc, cliCtx)).Methods("GET")

	r.HandleFunc("/checkpoint/last-no-ack", noackHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/state-dump", stateDumpHandlerFunc(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("checkpoint/latest-checkpoint", latestCheckpointHandlerFunc(cdc, cliCtx)).Methods("GET")
}

func checkpointBufferHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		err = cdc.UnmarshalBinary(res, &_checkpoint)
		if err != nil {
			RestLogger.Error("Unable to unmarshall", "Error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		RestLogger.Debug("Checkpoint fetched", "Checkpoint", _checkpoint.String())

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
		RestLogger.Debug("Fetching number of checkpoints from state")
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

		ackCount, err := strconv.ParseInt(string(res), 10, 64)
		if err != nil {
			RestLogger.Error("Unable to parse int", "Response", res, "Error", err)
			w.Write([]byte(err.Error()))
			return
		}
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
		err = cdc.UnmarshalBinary(res, &_checkpoint)
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

func checkpointHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		helper.InitHeimdallConfig("")
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
		checkpoint := HeaderBlock{
			Proposer:   helper.ZeroAddress,
			StartBlock: uint64(start),
			EndBlock:   uint64(end),
			RootHash:   ethcmn.BytesToHash(roothash),
		}
		result, err := json.Marshal(checkpoint)
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

func noackHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryStore(common.CheckpointNoACKCacheKey, "checkpoint")
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		lastAckTime, err := strconv.ParseInt(string(res), 10, 64)
		if err != nil {
			RestLogger.Error("Unable to parse int", "Response", res, "Error", err)
			w.Write([]byte(err.Error()))
			return
		}

		result, err := json.Marshal(map[string]interface{}{"result": lastAckTime})
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

type stateDump struct {
	ACKCount         int64                       `json:AckCount`
	CheckpointBuffer types.CheckpointBlockHeader `json:HeaderBlock`
	ValidatorCount   int                         `json:ValidatorCount`
	ValidatorSet     types.ValidatorSet          `json:ValidatorSet`
	LastNoACK        time.Time                   `json:LastNoACKTime`
}

// get all state-dump of heimdall
func stateDumpHandlerFunc(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// ACk count
		var ackCountInt int64
		ackcount, err := cliCtx.QueryStore(common.ACKCountKey, "checkpoint")
		if err == nil {
			ackCountInt, err = strconv.ParseInt(string(ackcount), 10, 64)
			if err != nil {
				RestLogger.Error("Unable to parse int for getting ack count", "Response", ackcount, "Error", err)
			}
		} else {
			RestLogger.Error("Unable to fetch ack count from store", "Error", err)
		}

		// checkpoint buffer
		var _checkpoint types.CheckpointBlockHeader
		_checkpointBufferBytes, err := cliCtx.QueryStore(common.BufferCheckpointKey, "checkpoint")
		if err == nil {
			err = cdc.UnmarshalBinary(_checkpointBufferBytes, &_checkpoint)
			if err != nil {
				RestLogger.Error("Unable to unmarshall checkpoint present in buffer", "Error", err, "CheckpointBuffer", _checkpointBufferBytes)
			}
		} else {
			RestLogger.Error("Unable to fetch checkpoint from buffer", "Error", err)
		}

		// validator count
		var validatorCount int
		var validatorSet types.ValidatorSet

		_validatorSet, err := cliCtx.QueryStore(common.CurrentValidatorSetKey, "staker")
		if err == nil {
			cdc.UnmarshalBinary(_validatorSet, &validatorSet)
		}
		validatorCount = len(validatorSet.Validators)

		// last no ack
		var lastACKTime int64
		lastNoACK, err := cliCtx.QueryStore(common.CheckpointNoACKCacheKey, "checkpoint")
		if err == nil {
			lastACKTime, err = strconv.ParseInt(string(lastNoACK), 10, 64)
		}

		state := stateDump{
			ACKCount:         ackCountInt,
			CheckpointBuffer: _checkpoint,
			ValidatorCount:   validatorCount,
			ValidatorSet:     validatorSet,
			LastNoACK:        time.Unix(lastACKTime, 0),
		}
		result, err := json.Marshal(map[string]interface{}{"result": state})
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

// get last checkpoint from store
func latestCheckpointHandlerFunc(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ackCount, err := cliCtx.QueryStore(common.ACKCountKey, "checkpoint")
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		ackCountInt, err := strconv.ParseInt(string(ackCount), 10, 64)
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		lastCheckpointKey := helper.GetConfig().ChildBlockInterval * uint64(ackCountInt)
		res, err := cliCtx.QueryStore(common.GetHeaderKey(lastCheckpointKey), "checkpoint")
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
		err = cdc.UnmarshalBinary(res, &_checkpoint)
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
