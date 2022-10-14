package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/ethereum/go-ethereum/common"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/helper"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmRest "github.com/maticnetwork/heimdall/types/rest"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/checkpoints/params", paramsHandlerFn(cliCtx)).Methods("GET")

	r.HandleFunc("/overview", overviewHandlerFn(cliCtx)).Methods("GET")

	r.HandleFunc("/checkpoints/buffer", checkpointBufferHandlerFn(cliCtx)).Methods("GET")

	r.HandleFunc("/checkpoints/count", checkpointCountHandlerFn(cliCtx)).Methods("GET")

	r.HandleFunc("/checkpoints/prepare", prepareCheckpointHandlerFn(cliCtx)).Methods("GET")

	r.HandleFunc("/checkpoints/latest", latestCheckpointHandlerFunc(cliCtx)).Methods("GET")

	r.HandleFunc("/checkpoints/last-no-ack", noackHandlerFn(cliCtx)).Methods("GET")

	r.HandleFunc("/checkpoints/list", checkpointListhandlerFn(cliCtx)).Methods("GET")

	r.HandleFunc("/checkpoints/{number}", checkpointByNumberHandlerFunc(cliCtx)).Methods("GET")

}

// HTTP request handler to query the auth params values
func paramsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParams)
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func checkpointBufferHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// fetch checkpoint
		result, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCheckpointBuffer), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

func checkpointCountHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		RestLogger.Debug("Fetching number of checkpoints from state")
		ackCountBytes, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAckCount), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, ackCountBytes, "No ack count found"); !ok {
			return
		}

		var ackCount uint64
		if err := json.Unmarshal(ackCountBytes, &ackCount); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := json.Marshal(map[string]interface{}{"result": ackCount})
		if err != nil {
			RestLogger.Error("Error while marshalling resposne to Json", "error", err)
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

func prepareCheckpointHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// Get params
		params := r.URL.Query()

		var result []byte
		var height int64
		var validatorSetBytes []byte

		// get start and start
		if params.Get("start") != "" && params.Get("end") != "" {
			start, err := strconv.ParseUint(params.Get("start"), 10, 64)
			if err != nil {
				hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			end, err := strconv.ParseUint(params.Get("end"), 10, 64)
			if err != nil {
				hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParams), nil)
			if err != nil {
				hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				RestLogger.Error("Unable to get checkpoint params", "Error", err)
				return
			}

			var params types.Params
			json.Unmarshal(res, &params)
			contractCallerObj, err := helper.NewContractCaller()

			// get headers
			roothash, err := contractCallerObj.GetRootHash(uint64(start), uint64(end), params.MaxCheckpointLength)
			if err != nil {
				RestLogger.Error("Unable to get roothash", "Start", start, "End", end, "Error", err)
				hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}

			//
			// Get current validator set
			//

			var validatorSet hmTypes.ValidatorSet
			validatorSetBytes, height, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", stakingTypes.QuerierRoute, stakingTypes.QueryCurrentValidatorSet), nil)
			if err == nil {
				err := json.Unmarshal(validatorSetBytes, &validatorSet)
				if err != nil {
					hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
					RestLogger.Error("Unable to get validator set to form proposer", "Error", err)
					return
				}
			}

			// header block -- checkpoint
			checkpoint := HeaderBlockResult{
				Proposer:   validatorSet.Proposer.Signer,
				StartBlock: uint64(start),
				EndBlock:   uint64(end),
				RootHash:   ethcmn.BytesToHash(roothash),
			}

			result, err = json.Marshal(checkpoint)
			if err != nil {
				RestLogger.Error("Error while marshalling resposne to Json", "error", err)
				hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
				return
			}
		} else {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, "`start` and `end` query params required")
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

// HeaderBlockResult represents header block result
type HeaderBlockResult struct {
	Proposer   hmTypes.HeimdallAddress `json:"proposer"`
	RootHash   common.Hash             `json:"rootHash"`
	StartBlock uint64                  `json:"startBlock"`
	EndBlock   uint64                  `json:"endBlock"`
}

func noackHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLastNoAck), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "Last NoAck not found"); !ok {
			return
		}

		var lastAckTime uint64
		if err := json.Unmarshal(res, &lastAckTime); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := json.Marshal(map[string]interface{}{"result": lastAckTime})
		if err != nil {
			RestLogger.Error("Error while marshalling resposne to Json", "error", err)
			hmRest.WriteErrorResponse(w, http.StatusNoContent, errors.New("Error while sending last ack time").Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

type stateDump struct {
	ACKCount         uint64               `json:"ack_count"`
	CheckpointBuffer *hmTypes.Checkpoint  `json:"checkpoint_buffer"`
	ValidatorCount   int                  `json:"validator_count"`
	ValidatorSet     hmTypes.ValidatorSet `json:"validator_set"`
	LastNoACK        time.Time            `json:"last_noack_time"`
}

// get all state-dump of heimdall
func overviewHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		//
		// Ack acount
		//

		var ackCountInt uint64
		ackCountBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAckCount), nil)
		if err == nil {
			// check content
			if ok := hmRest.ReturnNotFoundIfNoContent(w, ackCountBytes, "No ack count found"); ok {
				if err := json.Unmarshal(ackCountBytes, &ackCountInt); err != nil {
					// log and ignore
					RestLogger.Error("Error while unmarshing no-ack count", "error", err)
				}
			}
		}

		//
		// Checkpoint buffer
		//

		var _checkpoint *hmTypes.Checkpoint
		checkpointBufferBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCheckpointBuffer), nil)
		if err == nil {
			if len(checkpointBufferBytes) != 0 {
				_checkpoint = new(hmTypes.Checkpoint)
				if err = json.Unmarshal(checkpointBufferBytes, _checkpoint); err != nil {
					// log and ignore
					RestLogger.Error("Error while unmarshing checkpoint header", "error", err)
				}
			}
		}

		//
		// Current validator set
		//

		var validatorSet hmTypes.ValidatorSet
		validatorSetBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", stakingTypes.QuerierRoute, stakingTypes.QueryCurrentValidatorSet), nil)
		if err == nil {
			if err := json.Unmarshal(validatorSetBytes, &validatorSet); err != nil {
				// log and ignore
				RestLogger.Error("Error while unmarshing validator set", "error", err)
			}
		}

		// validator count
		validatorCount := len(validatorSet.Validators)

		//
		// Last no-ack
		//

		// last no ack
		var lastNoACKTime uint64
		lastNoACKBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLastNoAck), nil)
		if err == nil {
			// check content
			if ok := hmRest.ReturnNotFoundIfNoContent(w, lastNoACKBytes, "No last-no-ack count found"); ok {
				if err := json.Unmarshal(lastNoACKBytes, &lastNoACKTime); err != nil {
					// log and ignore
					RestLogger.Error("Error while unmarshing last no-ack time", "error", err)
				}
			}
		}

		//
		// State dump
		//

		state := stateDump{
			ACKCount:         ackCountInt,
			CheckpointBuffer: _checkpoint,
			ValidatorCount:   validatorCount,
			ValidatorSet:     validatorSet,
			LastNoACK:        time.Unix(int64(lastNoACKTime), 0),
		}

		result, err := json.Marshal(state)
		if err != nil {
			RestLogger.Error("Error while marshalling resposne to Json", "error", err)
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

// get last checkpoint from store
func latestCheckpointHandlerFunc(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		//
		// Get ack count
		//

		ackcountBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAckCount), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, ackcountBytes, "No ack count found"); !ok {
			return
		}

		var ackCount uint64
		if err := json.Unmarshal(ackcountBytes, &ackCount); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		//
		// Last checkpoint key
		//

		RestLogger.Debug("ACK Count fetched", "ackCount", ackCount)
		lastCheckpointKey := ackCount
		RestLogger.Debug("Last checkpoint key generated", "lastCheckpointKey", lastCheckpointKey)

		// get query params
		queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryCheckpointParams(lastCheckpointKey))
		if err != nil {
			return
		}

		//
		// Get checkpoint
		//

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCheckpoint), queryParams)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// error if no checkpoint found
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No checkpoint found"); !ok {
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

// get checkpoint by checkppint number from store
func checkpointByNumberHandlerFunc(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// get checkpoint number
		number, ok := rest.ParseUint64OrReturnBadRequest(w, vars["number"])
		if !ok {
			return
		}

		// get query params
		queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryCheckpointParams(number))
		if err != nil {
			return
		}

		// query checkpoint
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCheckpoint), queryParams)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No checkpoint found"); !ok {
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func checkpointListhandlerFn(
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := r.URL.Query()

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// get page
		page, ok := rest.ParseUint64OrReturnBadRequest(w, vars.Get("page"))
		if !ok {
			return
		}

		// get limit
		limit, ok := rest.ParseUint64OrReturnBadRequest(w, vars.Get("limit"))
		if !ok {
			return
		}

		// get query params
		queryParams, err := cliCtx.Codec.MarshalJSON(hmTypes.NewQueryPaginationParams(page, limit))
		if err != nil {
			return
		}

		// query checkpoint
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCheckpointList), queryParams)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No checkpoints found"); !ok {
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
