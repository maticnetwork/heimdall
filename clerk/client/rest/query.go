package rest

import (
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/clerk/types"
	hmRest "github.com/maticnetwork/heimdall/types/rest"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(
		"/clerk/event-record/list",
		recordListHandlerFn(cliCtx),
	).Methods("GET")

	r.HandleFunc(
		"/clerk/event-record/{recordId}",
		recordHandlerFn(cliCtx),
	).Methods("GET")
}

// recordHandlerFn returns record by record id
func recordHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// record id
		recordID, ok := rest.ParseUint64OrReturnBadRequest(w, vars["recordId"])
		if !ok {
			return
		}

		// get query params
		queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryRecordParams(recordID))
		if err != nil {
			return
		}

		// get record from store
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRecord), queryParams)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No record found"); !ok {
			return
		}

		hmRest.PostProcessResponse(w, cliCtx, res)
	}
}

func recordListHandlerFn(
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
		queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryRecordListParams(page, limit))
		if err != nil {
			return
		}

		// query records
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRecordList), queryParams)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No records found"); !ok {
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
