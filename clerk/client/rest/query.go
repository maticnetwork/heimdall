package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/clerk/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
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
	r.HandleFunc(
		"/clerk/isoldtx",
		DepositTxStatusHandlerFn(cliCtx),
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

		if vars.Get("page") != "" && vars.Get("limit") != "" {
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
		} else if vars.Get("from-time") != "" && vars.Get("to-time") != "" {
			// get from time (epoch)
			fromTime, ok := rest.ParseInt64OrReturnBadRequest(w, vars.Get("from-time"))
			if !ok {
				return
			}

			// get to time (epoch)
			toTime, ok := rest.ParseInt64OrReturnBadRequest(w, vars.Get("to-time"))
			if !ok {
				return
			}

			// get query params
			queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryTimeRangeParams(time.Unix(fromTime, 0), time.Unix(toTime, 0)))
			if err != nil {
				return
			}

			// query records
			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRecordListWithTime), queryParams)
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
		return
	}
}

// Returns deposit tx status information
func DepositTxStatusHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := r.URL.Query()
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// get logIndex
		logindex, ok := rest.ParseUint64OrReturnBadRequest(w, vars.Get("logindex"))
		if !ok {
			return
		}

		txHash := vars.Get("txhash")
		if txHash == "" {
			return
		}

		// get query params
		queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryRecordSequenceParams(txHash, logindex))
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		seqNo, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRecordSequence), queryParams)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// error if no tx status found
		if ok := hmRest.ReturnNotFoundIfNoContent(w, seqNo, "No sequence found"); !ok {
			return
		}

		res := true

		// return result
		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
