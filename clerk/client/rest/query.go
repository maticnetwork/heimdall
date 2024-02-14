// nolint
package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"

	"github.com/maticnetwork/heimdall/clerk/types"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmRest "github.com/maticnetwork/heimdall/types/rest"
)

//swagger:response clerkEventListResponse
type clerkEventListResponse struct {
	//in:body
	Output clerkEventList `json:"output"`
}

type clerkEventList struct {
	Height string  `json:"height"`
	Result []event `json:"result"`
}

//swagger:response clerkEventByIdResponse
type clerkEventByIdResponse struct {
	//in:body
	Output clerkEventById `json:"output"`
}

type clerkEventById struct {
	Height string `json:"height"`
	Result event  `json:"result"`
}

type event struct {
	Id         int64  `json:"id"`
	Contract   string `json:"contract"`
	Data       string `json:"data"`
	TxHash     string `json:"tx_hash"`
	LogIndex   int64  `json:"log_index"`
	BorChainId string `json:"bor_chain_id"`
	RecoedTime string `json:"record_time"`
}

//swagger:response clerkIsOldTxResponse
type clerkIdOldTxResponse struct {
	//in:body
	Output isOldTx `json:"output"`
}

type isOldTx struct {
	Height string `json:"height"`
	Result bool   `json:"result"`
}

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

//swagger:parameters clerkEventById
type clerkEventID struct {

	//ID of the checkpoint
	//required:true
	//in:path
	Id int64 `json:"recordID"`
}

// swagger:route GET /clerk/event-record/{recordID} clerk clerkEventById
// It returns the clerk event based on ID
// responses:
//
//	200: clerkEventByIdResponse
//
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

//swagger:parameters clerkEventList
type clerkEventListParams struct {

	//Page number
	//required:true
	//in:query
	Page int64 `json:"page"`

	//Limit per page
	//required:true
	//in:query
	Limit int64 `json:"limit"`
}

// swagger:route GET /clerk/event-record/list clerk clerkEventList
// It returns the clerk events list
// responses:
//
//	200: clerkEventListResponse
func recordListHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var logger = helper.Logger.With("module", "/clerk/event-record/list")
		logger.Info("Serving event record list")
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

		logger.Info("Serving event record list", "page", page)

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
		logger.Info("Serving event record list", "limit", limit)

		var (
			res []byte
			err error
		)

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

			logger.Info("Serving event record list", "from-time", fromTime, "to-time", toTime)

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

			logger.Info("Serving event record list", "from-id", fromID, "to-time", toTime)

			// get result by till time-range query
			res, err = tillTimeRangeQuery(cliCtx, fromID, toTime, limit)
		} else {
			// get result by range query
			res, err = rangeQuery(cliCtx, page, limit)
		}

		// send internal server error if error occurred during the query
		if err != nil {
			logger.Error("Error while querying event record list", "error", err.Error())
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No records found"); !ok {
			logger.Error("Error while querying event record list", "error", "No records found")
			return
		}

		logger.Info("Served event record list successfully")
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters clerkIsOldTx
type clerkTxParams struct {

	//Log Index of the transaction
	//required:true
	//in:query
	LogIndex int64 `json:"logindex"`

	//Hash of the transaction
	//required:true
	//in:query
	Txhash string `json:"txhash"`
}

// swagger:route GET /clerk/isoldtx clerk clerkIsOldTx
// It checks for whether the transaction is old or new.
// responses:
//
//	200: clerkIsOldTxResponse
//
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
		return jsoniter.ConfigFastest.Marshal(result)
	}

	var fromRecord types.EventRecord
	if err = jsoniter.ConfigFastest.Unmarshal(fromData, &fromRecord); err != nil {
		return nil, err
	}

	fromTime := fromRecord.RecordTime.Unix()

	rangeData, err := timeRangeQuery(cliCtx, fromTime, toTime, 1, limit)
	if err != nil {
		return nil, err
	}

	rangeRecords := make([]*types.EventRecord, 0)
	if err = jsoniter.ConfigFastest.Unmarshal(rangeData, &rangeRecords); err != nil {
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
			err = jsoniter.ConfigFastest.Unmarshal(recordData, &record)
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
	return jsoniter.ConfigFastest.Marshal(result)
}

//swagger:parameters clerkIsOldTx clerkEventList clerkEventById
type Height struct {

	//Block Height
	//in:query
	Height string `json:"height"`
}
