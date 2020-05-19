package rest

import (
	"encoding/json"
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
	r.HandleFunc(
		"/clerk/deposit-count",
		depositCountHandlerFn(cliCtx),
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
		var queryParams []byte
		var err error
		var query string

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		page := uint64(1) // default page
		if vars.Get("page") != "" {
			_page, ok := rest.ParseUint64OrReturnBadRequest(w, vars.Get("page"))
			if !ok {
				return
			}

			page = _page
		}

		limit := uint64(50) // default limit
		if vars.Get("limit") != "" {
			_limit, ok := rest.ParseUint64OrReturnBadRequest(w, vars.Get("limit"))
			if !ok {
				return
			}

			limit = _limit
		}

		if vars.Get("from-time") != "" && vars.Get("to-time") != "" {
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
			queryParams, err = cliCtx.Codec.MarshalJSON(types.NewQueryTimeRangePaginationParams(time.Unix(fromTime, 0), time.Unix(toTime, 0), page, limit))
			if err != nil {
				return
			}

			query = types.QueryRecordListWithTime
		} else {
			// get query params
			queryParams, err = cliCtx.Codec.MarshalJSON(hmTypes.NewQueryPaginationParams(page, limit))
			if err != nil {
				return
			}

			query = types.QueryRecordList
		}

		// query records
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, query), queryParams)
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

// DepositTxStatusHandlerFn returns deposit tx status information
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

func depositCountHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		RestLogger.Debug("Fetching number of deposits from state")
		depositCountBytes, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryDepositCount), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, depositCountBytes, "No deposit count found"); !ok {
			return
		}

		var depositCount uint64
		if err := json.Unmarshal(depositCountBytes, &depositCount); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := json.Marshal(&depositCount)
		if err != nil {
			RestLogger.Error("Error while marshalling resposne to Json", "error", err)
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, result)
	}
}
