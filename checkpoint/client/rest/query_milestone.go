package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"

	"github.com/maticnetwork/heimdall/checkpoint/types"
	hmRest "github.com/maticnetwork/heimdall/types/rest"
)

func registerQueryMilestoneRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/milestone/latest", milestoneLatestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/milestone/count", milestoneCountHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/milestone/lastNoAck", latestNoAckMilestoneHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/milestone/{number}", milestoneByNumberHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/milestone/noAck/{id}", noAckMilestoneByIDHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/milestone/ID/{id}", noAckMilestoneByIDHandlerFn(cliCtx)).Methods("GET")
}

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

		countBytes, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCount), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, countBytes, "No milestone count found"); !ok {
			return
		}

		var count uint64
		if err := jsoniter.Unmarshal(countBytes, &count); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := jsoniter.Marshal(map[string]interface{}{"count": count})
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

func latestNoAckMilestoneHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// fetch checkpoint
		result, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLatestNoAckMilestone), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, result, "No last ack milestone found"); !ok {
			return
		}

		var milestoneID string
		if err := jsoniter.Unmarshal(result, &milestoneID); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := jsoniter.Marshal(map[string]interface{}{"result": milestoneID})
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func noAckMilestoneByIDHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// get milestone number
		id := vars["id"]

		// get query params
		queryID, err := cliCtx.Codec.MarshalJSON(types.NewQueryMilestoneID(id))
		if err != nil {
			return
		}

		result, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryNoAckMilestoneByID), queryID)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var val bool
		if err := jsoniter.Unmarshal(result, &val); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := jsoniter.Marshal(map[string]interface{}{"result": val})
		if err != nil {
			RestLogger.Error("Error while marshalling resposne to Json", "error", err)
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func milestoneByIDHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// get milestone number
		id := vars["id"]

		// get query params
		queryID, err := cliCtx.Codec.MarshalJSON(types.NewQueryMilestoneID(id))
		if err != nil {
			return
		}

		result, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryMilestoneByID), queryID)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var val bool
		if err := jsoniter.Unmarshal(result, &val); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		res, err := jsoniter.Marshal(map[string]interface{}{"result": val})
		if err != nil {
			RestLogger.Error("Error while marshalling resposne to Json", "error", err)
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		cliCtx = cliCtx.WithHeight(height)

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
