//nolint
package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/milestone/types"

	hmRest "github.com/maticnetwork/heimdall/types/rest"
)

//It represents the milestone parameters
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

//It represents the milestone
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
	r.HandleFunc("/milestone", milestoneHandlerFn(cliCtx)).Methods("GET")
}

// swagger:route GET /milestone/params milestone milistoneParams
// It returns the milestone parameters
// responses:
//   200: milestoneParamsResponse
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
//   200: milestoneResponse
func milestoneHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// fetch checkpoint
		result, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryMilestone), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

//swagger:parameters checkpointList checkpointById checkpointLatest overview checkpointLastNoAck checkpointPrepare checkpointCount checkpointParams checkpointBuffer
type Height struct {

	//Block Height
	//in:query
	Height string `json:"height"`
}
