// nolint
package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"

	"github.com/ethereum/go-ethereum/common"
	ethcmn "github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/helper"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmRest "github.com/maticnetwork/heimdall/types/rest"
)

// It represents the checkpoint parameters
//
//swagger:response checkpointParamsResponse
type checkpointParamsResponse struct {
	//in:body
	Output checkpointParams `json:"output"`
}

type checkpointParams struct {
	Height string `json:"height"`
	Result params `json:"result"`
}

type params struct {
	CheckpointBufferTime    int `json:"checkpoint_buffer_time"`
	AvgCheckpointLength     int `json:"avg_checkpoint_length"`
	MaxCheckPoint           int `json:"max_checkpoint_length"`
	ChildChainBlockInterval int `json:"child_chain_block_interval"`
}

// It represents the checkpoint
//
//swagger:response checkpointResponse
type checkpointResponse struct {
	//in:body
	Output checkpointStructure `json:"output"`
}
type checkpointStructure struct {
	Height string     `json:"height"`
	Result checkpoint `json:"result"`
}

type checkpoint struct {
	Proposer   string `json:"proposer"`
	StartBlock int64  `json:"start_block"`
	EndBlock   int64  `json:"end_block"`
	RootHash   string `json:"root_hash"`
	BorChainId string `json:"bor_chain_id"`
	Timestamp  int64  `json:"timestamp"`
}

// It represents the checkpoint prepare
//
//swagger:response checkpointPrepareResponse
type checkpointPrepareResponse struct {
	//in:body
	Output checkpointParams `json:"output"`
}

type checkpointPrepareStructure struct {
	Height string            `json:"height"`
	Result prepareCheckpoint `json:"result"`
}

type prepareCheckpoint struct {
	Proposer   string `json:"proposer"`
	StartBlock int64  `json:"start_block"`
	EndBlock   int64  `json:"end_block"`
	RootHash   string `json:"root_hash"`
}

// It represents the checkpoint list
//
//swagger:response checkpointListResponse
type checkpointListResponse struct {
	//in:body
	Output checkpointListStructure `json:"output"`
}

type checkpointListStructure struct {
	Height string       `json:"height"`
	Result []checkpoint `json:"result"`
}

// It represents the last-no-ack
//
//swagger:response lastNoAckResponse
type lastNoAckResponse struct {
	//in:body
	Output lastNoAckResponseStructure `json:"output"`
}

type lastNoAckResponseStructure struct {
	Height string    `json:"height"`
	Result lastNoAck `json:"result"`
}

type lastNoAck struct {
	Result int64 `json:"result"`
}

// It represents the checkpoint count
//
//swagger:response checkpointCountResponse
type checkpointCountResponse struct {
	//in:body
	Output checkpointCountResponseStructure `json:"output"`
}

type checkpointCountResponseStructure struct {
	Height string          `json:"height"`
	Result checkpointCount `json:"result"`
}

type checkpointCount struct {
	Result int64 `json:"result"`
}

// It represents the overview
//
//swagger:response overviewResponse
type overviewResponse struct {
	//in:body
	Output overviewResponseStructure `json:"output"`
}

type overviewResponseStructure struct {
	Height string   `json:"height"`
	Result overview `json:"result"`
}

type overview struct {
	AckCount         int64        `json:"ack_count"`
	CheckpointBuffer checkpoint   `json:"checkpoint_buffer"`
	ValidatorCount   int64        `json:"validator_Count"`
	ValidatorSet     validatorSet `json:"validator_set"`
	LastNoAckTime    string       `json:"last_noack_time"`
}

type validatorSet struct {
	Validators []validator `json:"validators"`
	Proposer   validator   `json:"proposer"`
}

type validator struct {
	ID           int    `json:"ID"`
	StartEpoch   int    `json:"startEpoch"`
	EndEpoch     int    `json:"endEpoch"`
	Nonce        int    `json:"nonce"`
	Power        int    `json:"power"`
	PubKey       string `json:"pubKey"`
	Signer       string `json:"signer"`
	Last_Updated string `json:"last_updated"`
	Jailed       bool   `json:"jailed"`
	Accum        int    `json:"accum"`
}

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

// swagger:route GET /checkpoints/params checkpoint checkpointParams
// It returns the checkpoint parameters
// responses:
//
//	200: checkpointParamsResponse
//
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

// swagger:route GET /checkpoints/buffer checkpoint checkpointBuffer
// It returns the checkpoint buffer
// responses:
//
//	200: checkpointResponse
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

// swagger:route GET /checkpoints/count checkpoint checkpointCount
// It returns the checkpoint counts
// responses:
//
//	200: checkpointCountResponse
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
		if err := jsoniter.ConfigFastest.Unmarshal(ackCountBytes, &ackCount); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := jsoniter.ConfigFastest.Marshal(map[string]interface{}{"result": ackCount})
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

//swagger:parameters checkpointPrepare
type checkpointPrepareParams struct {

	//Start Block
	//required:true
	//in:query
	Start int64 `json:"start"`

	//End Block
	//required:true
	//in:query
	End int64 `json:"end"`
}

// swagger:route GET /checkpoints/prepare checkpoint checkpointPrepare
// It returns the prepared checkpoint
// responses:
//
//	200: checkpointPrepareResponse
func prepareCheckpointHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// Get params
		params := r.URL.Query()

		var (
			result            []byte
			height            int64
			validatorSetBytes []byte
		)

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
			if err = jsoniter.ConfigFastest.Unmarshal(res, &params); err != nil {
				RestLogger.Error("Unable to unmarshal params", "Start", start, "End", end, "Error", err)
				hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

				return
			}

			contractCallerObj, err := helper.NewContractCaller()
			if err != nil {
				RestLogger.Error("Unable to create contract caller", "Start", start, "End", end, "Error", err)
				hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

				return
			}

			// get headers
			roothash, err := contractCallerObj.GetRootHash(start, end, params.MaxCheckpointLength)
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
				if err = jsoniter.ConfigFastest.Unmarshal(validatorSetBytes, &validatorSet); err != nil {
					hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
					RestLogger.Error("Unable to get validator set to form proposer", "Error", err)

					return
				}
			}

			// header block -- checkpoint
			checkpoint := HeaderBlockResult{
				Proposer:   validatorSet.Proposer.Signer,
				StartBlock: start,
				EndBlock:   end,
				RootHash:   ethcmn.BytesToHash(roothash),
			}

			result, err = jsoniter.ConfigFastest.Marshal(checkpoint)
			if err != nil {
				RestLogger.Error("Error while marshalling response to Json", "error", err)
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

// swagger:route GET /checkpoints/last-no-ack checkpoint checkpointLastNoAck
// It returns the last no ack
// responses:
//
//	200: lastNoAckResponse
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
		if err := jsoniter.ConfigFastest.Unmarshal(res, &lastAckTime); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := jsoniter.ConfigFastest.Marshal(map[string]interface{}{"result": lastAckTime})
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
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

// swagger:route GET /overview checkpoint overview
// It returns the complete overview
// responses:
//
//	200: overviewResponse
//
// get all state-dump of heimdall
func overviewHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		//
		// Ack account
		//

		var ackCountInt uint64

		ackCountBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAckCount), nil)
		if err == nil {
			// check content
			if hmRest.ReturnNotFoundIfNoContent(w, ackCountBytes, "No ack count found") {
				if err = jsoniter.ConfigFastest.Unmarshal(ackCountBytes, &ackCountInt); err != nil {
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
				if err = jsoniter.ConfigFastest.Unmarshal(checkpointBufferBytes, _checkpoint); err != nil {
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
			if err := jsoniter.ConfigFastest.Unmarshal(validatorSetBytes, &validatorSet); err != nil {
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
			if hmRest.ReturnNotFoundIfNoContent(w, lastNoACKBytes, "No last-no-ack count found") {
				if err = jsoniter.ConfigFastest.Unmarshal(lastNoACKBytes, &lastNoACKTime); err != nil {
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

		result, err := jsoniter.ConfigFastest.Marshal(state)
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		rest.PostProcessResponse(w, cliCtx, result)
	}
}

// swagger:route GET /checkpoints/latest checkpoint checkpointLatest
// It returns the last checkpoint from the store
// responses:
//
//	200: checkpointResponse
//
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

		ackcountBytes, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAckCount), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, ackcountBytes, "No ack count found"); !ok {
			return
		}

		var ackCount uint64
		if err := jsoniter.ConfigFastest.Unmarshal(ackcountBytes, &ackCount); err != nil {
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

		var checkpointUnmarshal hmTypes.Checkpoint
		if err = jsoniter.ConfigFastest.Unmarshal(res, &checkpointUnmarshal); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		checkpointWithID := &CheckpointWithID{
			ID:         ackCount,
			Proposer:   checkpointUnmarshal.Proposer,
			StartBlock: checkpointUnmarshal.StartBlock,
			EndBlock:   checkpointUnmarshal.EndBlock,
			RootHash:   checkpointUnmarshal.RootHash,
			BorChainID: checkpointUnmarshal.BorChainID,
			TimeStamp:  checkpointUnmarshal.TimeStamp,
		}

		resWithID, err := jsoniter.ConfigFastest.Marshal(checkpointWithID)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		// error if no checkpoint found
		if !hmRest.ReturnNotFoundIfNoContent(w, resWithID, "No checkpoint found") {
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, resWithID)
	}
}

// Temporary Checkpoint struct to store the Checkpoint ID
type CheckpointWithID struct {
	ID         uint64                  `json:"id"`
	Proposer   hmTypes.HeimdallAddress `json:"proposer"`
	StartBlock uint64                  `json:"start_block"`
	EndBlock   uint64                  `json:"end_block"`
	RootHash   hmTypes.HeimdallHash    `json:"root_hash"`
	BorChainID string                  `json:"bor_chain_id"`
	TimeStamp  uint64                  `json:"timestamp"`
}

//swagger:parameters checkpointById
type checkpointID struct {

	//ID of the checkpoint
	//required:true
	//in:path
	Id int64 `json:"id"`
}

// swagger:route GET /checkpoints/{id} checkpoint checkpointById
// It returns the checkpoint by ID
// responses:
//
//	200: checkpointResponse
//
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
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCheckpoint), queryParams)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var checkpointUnmarshal hmTypes.Checkpoint
		if err = jsoniter.ConfigFastest.Unmarshal(res, &checkpointUnmarshal); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		checkpointWithID := &CheckpointWithID{
			ID:         number,
			Proposer:   checkpointUnmarshal.Proposer,
			StartBlock: checkpointUnmarshal.StartBlock,
			EndBlock:   checkpointUnmarshal.EndBlock,
			RootHash:   checkpointUnmarshal.RootHash,
			BorChainID: checkpointUnmarshal.BorChainID,
			TimeStamp:  checkpointUnmarshal.TimeStamp,
		}

		resWithID, err := jsoniter.ConfigFastest.Marshal(checkpointWithID)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, resWithID, "No checkpoint found"); !ok {
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, resWithID)
	}
}

//swagger:parameters checkpointList
type checkpointListParams struct {

	//Page number
	//required:true
	//in:query
	Page int64 `json:"page"`

	//Limit per page
	//required:true
	//in:query
	Limit int64 `json:"limit"`
}

// swagger:route GET /checkpoints/list checkpoint checkpointList
// It returns the checkpoints list
// responses:
//
//	200: checkpointListResponse
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
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCheckpointList), queryParams)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No checkpoints found"); !ok {
			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters checkpointList checkpointById checkpointLatest overview checkpointLastNoAck checkpointPrepare checkpointCount checkpointParams checkpointBuffer
type Height struct {

	//Block Height
	//in:query
	Height string `json:"height"`
}
