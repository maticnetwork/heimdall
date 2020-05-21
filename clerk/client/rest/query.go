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

		// get record from store
		res, err := recordQuery(cliCtx, recordID)
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

			// truncate limit to default limit
			if _limit < limit {
				limit = _limit
			}
		}

		var res []byte
		var err error

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

			// get result by time-range query
			res, err = timeRangeQuery(cliCtx, fromTime, toTime, page, limit)

		} else if vars.Get("from-id") != "" && vars.Get("to-time") != "" {
			// get from id
			fromID, ok := rest.ParseUint64OrReturnBadRequest(w, vars.Get("from-id"))
			if !ok {
				return
			}

			// get to time (epoch)
			toTime, ok := rest.ParseInt64OrReturnBadRequest(w, vars.Get("to-time"))
			if !ok {
				return
			}

			// get result by till time-range query
			res, err = tillTimeRangeQuery(cliCtx, fromID, toTime, limit)
		} else {
			// get result by range query
			res, err = rangeQuery(cliCtx, page, limit)
		}

		// send internal server error if error occured during the query
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
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

//
// Internal helpers
//

func recordQuery(cliCtx context.CLIContext, recordID uint64) ([]byte, error) {
	// get query params
	queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryRecordParams(recordID))
	if err != nil {
		return nil, err
	}

	// get record from store
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryRecord), queryParams)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func timeRangeQuery(cliCtx context.CLIContext, fromTime int64, toTime int64, page uint64, limit uint64) ([]byte, error) {
	// get query params
	queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQueryTimeRangePaginationParams(time.Unix(fromTime, 0), time.Unix(toTime, 0), page, limit))
	if err != nil {
		return nil, err
	}

	// set query as record list with time
	query := types.QueryRecordListWithTime

	// query records
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, query), queryParams)
	if err != nil {
		return nil, err
	}

	// return result
	return res, nil
}

func rangeQuery(cliCtx context.CLIContext, page uint64, limit uint64) ([]byte, error) {
	// get query params
	queryParams, err := cliCtx.Codec.MarshalJSON(hmTypes.NewQueryPaginationParams(page, limit))
	if err != nil {
		return nil, err
	}

	// set query as record list
	query := types.QueryRecordList

	// query records
	res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, query), queryParams)
	if err != nil {
		return nil, err
	}

	// return result
	return res, nil
}

func tillTimeRangeQuery(cliCtx context.CLIContext, fromID uint64, toTime int64, limit uint64) ([]byte, error) {
	result := make([]*types.EventRecord, 0, limit)

	// if from id not found, return empty result
	fromData, err := recordQuery(cliCtx, fromID)
	if err != nil {
		return json.Marshal(result)
	}

	var fromRecord types.EventRecord
	err = json.Unmarshal(fromData, &fromRecord)
	if err != nil {
		return nil, err
	}

	fromTime := fromRecord.RecordTime.Unix()
	rangeData, err := timeRangeQuery(cliCtx, fromTime, toTime, 1, limit)
	if err != nil {
		return nil, err
	}

	rangeRecords := make([]*types.EventRecord, 0)
	err = json.Unmarshal(rangeData, &rangeRecords)
	if err != nil {
		return nil, err
	}

	rangeMapping := make(map[uint64]*types.EventRecord)
	for _, r := range rangeRecords {
		rangeMapping[r.ID] = r
	}

	nextID := fromID
	toTimeObj := time.Unix(toTime, 0)
	for nextID-fromID < limit {
		if found, ok := rangeMapping[nextID]; ok {
			result = append(result, found)
		} else {
			// fetch record for nextID and unmarshal to record
			recordData, err := recordQuery(cliCtx, nextID)
			if err != nil {
				break
			}

			var record types.EventRecord
			err = json.Unmarshal(recordData, &record)
			if err != nil {
				return nil, err
			}

			// checks if record time < to time
			if !record.RecordTime.Before(toTimeObj) {
				break
			}

			// add into result
			result = append(result, &record)
		}

		nextID++
	}

	// return result in json
	return json.Marshal(result)
}
