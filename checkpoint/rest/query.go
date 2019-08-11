package rest

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/rest"
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
	r.HandleFunc("/checkpoint/latest-checkpoint",
		latestCheckpointHandlerFunc(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/checkpoint/{start}/{end}",
		checkpointHandlerFn(cdc, cliCtx)).Methods("GET")

	r.HandleFunc("/checkpoint/last-no-ack",
		noackHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/overview",
		overviewHandlerFunc(cdc, cliCtx)).Methods("GET")
	helper.InitHeimdallConfig("")
}

func checkpointBufferHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryStore(checkpoint.BufferCheckpointKey, "checkpoint")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// the query will return empty if there is no data in buffer
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNoContent, err.Error())
			return
		}

		var _checkpoint types.CheckpointBlockHeader
		err = cdc.UnmarshalBinaryBare(res, &_checkpoint)
		if err != nil {
			RestLogger.Error("Unable to unmarshall", "Error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		RestLogger.Debug("Checkpoint fetched", "Checkpoint", _checkpoint.String())

		result, err := json.Marshal(&_checkpoint)
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

func checkpointCountHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		RestLogger.Debug("Fetching number of checkpoints from state")
		res, err := cliCtx.QueryStore(staking.ACKCountKey, "checkpoint")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		// The query will return empty if there is no data
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNoContent, err.Error())
			return
		}
		ackCount, err := strconv.ParseInt(string(res), 10, 64)
		if err != nil {
			RestLogger.Error("Unable to parse int", "Response", res, "Error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		result, err := json.Marshal(map[string]interface{}{"result": ackCount})
		if err != nil {
			RestLogger.Error("Error while marshalling resposne to Json", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

func checkpointHeaderHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		// get header number
		headerNumber, ok := rest.ParseUint64OrReturnBadRequest(w, vars["headerBlockIndex"])
		if !ok {
			return
		}

		res, err := cliCtx.QueryStore(checkpoint.GetHeaderKey(uint64(headerNumber)), "checkpoint")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		// the query will return empty if there is no data
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		var _checkpoint types.CheckpointBlockHeader
		err = cdc.UnmarshalBinaryBare(res, &_checkpoint)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		result, err := json.Marshal(&_checkpoint)
		if err != nil {
			RestLogger.Error("Error while marshalling resposne to Json", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

func checkpointHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		start, err := strconv.Atoi(vars["start"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		end, err := strconv.Atoi(vars["end"])
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		roothash, err := checkpoint.GetHeaders(uint64(start), uint64(end))
		if err != nil {
			RestLogger.Error("Unable to get header", "Start", start, "End", end, "Error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var validatorSet types.ValidatorSet

		_validatorSet, err := cliCtx.QueryStore(staking.CurrentValidatorSetKey, "staking")
		if err == nil {
			err := cdc.UnmarshalBinaryBare(_validatorSet, &validatorSet)
			if err != nil {
				rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				RestLogger.Error("Unable to get validator set to form proposer", "Error", err)
				return
			}
		}

		checkpoint := HeaderBlock{
			Proposer:   validatorSet.Proposer.Signer,
			StartBlock: uint64(start),
			EndBlock:   uint64(end),
			RootHash:   ethcmn.BytesToHash(roothash),
		}

		result, err := json.Marshal(checkpoint)
		if err != nil {
			RestLogger.Error("Error while marshalling resposne to Json", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

func noackHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := cliCtx.QueryStore(checkpoint.CheckpointNoACKCacheKey, "checkpoint")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNoContent, err.Error())
			return
		}

		lastAckTime, err := strconv.ParseInt(string(res), 10, 64)
		if err != nil {
			RestLogger.Error("Unable to parse int", "Response", res, "Error", err)
			rest.WriteErrorResponse(w, http.StatusNoContent, err.Error())
			return
		}

		result, err := json.Marshal(map[string]interface{}{"result": lastAckTime})
		if err != nil {
			RestLogger.Error("Error while marshalling resposne to Json", "error", err)
			rest.WriteErrorResponse(w, http.StatusNoContent, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

type stateDump struct {
	ACKCount         int64                         `json:AckCount`
	CheckpointBuffer types.CheckpointBlockHeader   `json:HeaderBlock`
	ValidatorCount   int                           `json:ValidatorCount`
	ValidatorSet     types.ValidatorSet            `json:ValidatorSet`
	LastNoACK        time.Time                     `json:LastNoACKTime`
	Headers          []types.CheckpointBlockHeader `json:"headers"`
}

// get all state-dump of heimdall
func overviewHandlerFunc(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// ACk count
		var ackCountInt int64
		ackcount, err := cliCtx.QueryStore(staking.ACKCountKey, "staking")
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
		_checkpointBufferBytes, err := cliCtx.QueryStore(checkpoint.BufferCheckpointKey, "checkpoint")
		if err == nil {
			if len(_checkpointBufferBytes) != 0 {
				err = cdc.UnmarshalBinaryBare(_checkpointBufferBytes, &_checkpoint)
				if err != nil {
					RestLogger.Error("Unable to unmarshall checkpoint present in buffer", "Error", err, "CheckpointBuffer", _checkpointBufferBytes)
				}
			} else {
				RestLogger.Error("No checkpoint present in buffer")
			}
		} else {
			RestLogger.Error("Unable to fetch checkpoint from buffer", "Error", err)
		}

		// validator count
		var validatorCount int
		var validatorSet types.ValidatorSet

		_validatorSet, err := cliCtx.QueryStore(staking.CurrentValidatorSetKey, "staking")
		if err == nil {
			cdc.UnmarshalBinaryBare(_validatorSet, &validatorSet)
		}
		validatorCount = len(validatorSet.Validators)

		// last no ack
		var lastNoACKTime int64
		lastNoACK, err := cliCtx.QueryStore(checkpoint.CheckpointNoACKCacheKey, "checkpoint")
		if err == nil {
			lastNoACKTime, err = strconv.ParseInt(string(lastNoACK), 10, 64)
		}

		var headers []types.CheckpointBlockHeader
		storedHeaders, err := cliCtx.QuerySubspace(checkpoint.HeaderBlockKey, "checkpoint")
		if err != nil {
			RestLogger.Error("Unable to query subspace for headers", "Error", err)
		}
		for _, kv_pair := range storedHeaders {
			var checkpointHeader types.CheckpointBlockHeader
			if cdc.UnmarshalBinaryBare(kv_pair.Value, &checkpointHeader); err != nil {
				RestLogger.Error("Unable to unmarshall header", "Error", err, "Value", kv_pair.Value)
			}
			headers = append(headers, checkpointHeader)
		}

		state := stateDump{
			ACKCount:         ackCountInt,
			CheckpointBuffer: _checkpoint,
			ValidatorCount:   validatorCount,
			ValidatorSet:     validatorSet,
			LastNoACK:        time.Unix(lastNoACKTime, 0),
			Headers:          headers,
		}
		result, err := json.Marshal(map[string]interface{}{"result": state})
		if err != nil {
			RestLogger.Error("Error while marshalling resposne to Json", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

// get last checkpoint from store
func latestCheckpointHandlerFunc(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ackCount, err := cliCtx.QueryStore(staking.ACKCountKey, "staking")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		ackCountInt, err := strconv.ParseInt(string(ackCount), 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		RestLogger.Debug("ACK Count fetched", "ACKCount", ackCountInt)

		lastCheckpointKey := helper.GetConfig().ChildBlockInterval * uint64(ackCountInt)
		RestLogger.Debug("Last checkpoint key generated", "LastCheckpointKey", lastCheckpointKey, "min", helper.GetConfig().ChildBlockInterval)
		res, err := cliCtx.QueryStore(checkpoint.GetHeaderKey(lastCheckpointKey), "checkpoint")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		// the query will return empty if there is no data
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		var _checkpoint types.CheckpointBlockHeader
		err = cdc.UnmarshalBinaryBare(res, &_checkpoint)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		RestLogger.Debug("Fetched last checkpoint", "Checkpoint", _checkpoint)
		result, err := json.Marshal(&_checkpoint)
		if err != nil {
			RestLogger.Error("Error while marshalling resposne to Json", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, result)
	}
}
