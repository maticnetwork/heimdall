// nolint
package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/milestone/types"

	hmRest "github.com/maticnetwork/heimdall/types/rest"
)

// It represents the milestone parameters
//
//swagger:response milestoneParamsResponse
type milestoneParamsResponse struct {
	//in:body
	Output milestoneParams `json:"output"`
}

type milestoneParams struct {
	Height string `json:"height"`
	Result params `json:"result"`
}

type params struct {
	SprintLength int `json:"sprint_length"`
}

// It represents the milestone
//
//swagger:response milestoneResponse
type milestoneResponse struct {
	//in:body
	Output milestoneStructure `json:"output"`
}
type milestoneStructure struct {
	Height string    `json:"height"`
	Result milestone `json:"result"`
}

type milestone struct {
	Proposer   string `json:"proposer"`
	StartBlock int64  `json:"start_block"`
	EndBlock   int64  `json:"end_block"`
	RootHash   string `json:"root_hash"`
	BorChainId string `json:"bor_chain_id"`
	Timestamp  int64  `json:"timestamp"`
}

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/milestone/params", paramsHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/milestone", milestoneLatestHandlerFn(cliCtx)).Methods("GET")
	//r.HandleFunc("/milestone/count", milestoneCountHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/milestone/{number}", milestoneByNumberHandlerFn(cliCtx)).Methods("GET")
}

// swagger:route GET /milestone/params milestone milistoneParams
// It returns the milestone parameters
// responses:
//
//	200: milestoneParamsResponse
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

// swagger:route GET /milestone milestone milestone
// It returns the milestone
// responses:
//
//	200: milestoneResponse
func milestoneLatestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// fetch checkpoint
		result, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLatestMilestone), nil)
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
func milestoneCountHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		RestLogger.Error("Fetching number of milestone from state")

		countBytes, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCount), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, countBytes, "No milestone count found"); !ok {
			return
		}

		RestLogger.Error("Fetching number of milestone from state")
		var count uint64
		if err := json.Unmarshal(countBytes, &count); err != nil {
			//hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		RestLogger.Error("Fetching number of milestone from state")

		result, err := json.Marshal(map[string]interface{}{"count": count})
		if err != nil {
			RestLogger.Error("Error while marshalling resposne to Json", "error", err)
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

func milestoneByNumberHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// get milestone number
		number, ok := rest.ParseUint64OrReturnBadRequest(w, vars["number"])
		if !ok {
			return
		}

		// get query params
		queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryMilestoneParams(number))
		if err != nil {
			return
		}

		// query checkpoint
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryMilestoneByNumber), queryParams)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
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
