// Package classification HiemdallRest API
//
//	    Schemes: http
//	    BasePath: /
//	    Version: 0.0.1
//	    title: Heimdall APIs
//	    Consumes:
//	    - application/json
//		   Host:localhost:1317
//	    - application/json
//
// nolint
//
//swagger:meta
package rest

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"

	"github.com/maticnetwork/heimdall/bor/types"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/helper"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmRest "github.com/maticnetwork/heimdall/types/rest"
)

type HeimdallSpanResultWithHeight struct {
	Height int64
	Result []byte
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

type span struct {
	SpanID     int `json:"span_id"`
	StartBlock int `json:"start_block"`
	EndBlock   int `json:"end_block"`
	//in:body
	ValidatorSet      validatorSet `json:"validator_set"`
	SelectedProducers []validator  `json:"selected_producer"`
	BorChainId        string       `json:"bor_chain_id"`
}

type validatorSet struct {
	Validators []validator `json:"validators"`
	Proposer   validator   `json:"Proposer"`
}

// It represents the list of spans
//
//swagger:response borSpanListResponse
type borSpanListResponse struct {
	//in:body
	Output borSpanList `json:"output"`
}

type borSpanList struct {
	Height string `json:"height"`
	Result []span `json:"result"`
}

// It represents the span
//
//swagger:response borSpanResponse
type borSpanResponse struct {
	//in:body
	Output borSpan `json:"output"`
}

type borSpan struct {
	Height string `json:"height"`
	Result span   `json:"result"`
}

// It represents the bor span parameters
//
//swagger:response borSpanParamsResponse
type borSpanParamsResponse struct {
	//in:body
	Output borSpanParams `json:"output"`
}

type borSpanParams struct {
	Height string     `json:"height"`
	Result spanParams `json:"result"`
}

type spanParams struct {

	//type:integer
	SprintDuration int64 `json:"sprint_duration"`
	//type:integer
	SpanDuration int64 `json:"span_duration"`
	//type:integer
	ProducerCount int64 `json:"producer_count"`
}

// It represents the next span seed
//
//swagger:response borNextSpanSeedResponse
type borNextSpanSeedResponse struct {
	//in:body
	Output spanSeed `json:"output"`
}

type spanSeed struct {
	Height string `json:"height"`
	Result string `json:"result"`
}

var spanOverrides map[uint64]*HeimdallSpanResultWithHeight = nil

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/bor/span/list", spanListHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/bor/span/{id}", spanHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/bor/latest-span", latestSpanHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/bor/prepare-next-span", prepareNextSpanHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/bor/next-span-seed", fetchNextSpanSeedHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/bor/params", paramsHandlerFn(cliCtx)).Methods("GET")
}

// swagger:route GET /bor/next-span-seed bor borNextSpanSeed
// It returns the seed for the next span
// responses:
//   200: borNextSpanSeedResponse

func fetchNextSpanSeedHandlerFn(
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryNextSpanSeed), nil)
		if err != nil {
			RestLogger.Error("Error while fetching next span seed  ", "Error", err.Error())
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		RestLogger.Debug("nextSpanSeed querier response", "res", res)

		// error if span seed found
		if !hmRest.ReturnNotFoundIfNoContent(w, res, "NextSpanSeed not found") {
			RestLogger.Error("NextSpanSeed not found ", "Error", err.Error())
			return
		}

		// return result
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters borSpanList
type borSpanListParam struct {

	//Page Number
	//required:true
	//type:integer
	//in:query
	Page int `json:"page"`

	//Limit
	//required:true
	//type:integer
	//in:query
	Limit int `json:"limit"`
}

// swagger:route GET /bor/span/list bor borSpanList
// It returns the list of Bor Span
// responses:
//
//	200: borSpanListResponse
func spanListHandlerFn(
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

		// query spans
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySpanList), queryParams)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No spans found"); !ok {
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters borSpanById
type borSpanById struct {

	//Id number of the span
	//required:true
	//type:integer
	//in:path
	Id int `json:"id"`
}

// swagger:route GET /bor/span/{id} bor borSpanById
// It returns the span based on ID
// responses:
//
//	200: borSpanResponse
func spanHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		vars := mux.Vars(r)

		// get to address
		spanID, ok := rest.ParseUint64OrReturnBadRequest(w, vars["id"])
		if !ok {
			return
		}

		var (
			res            []byte
			height         int64
			spanOverridden bool
		)

		if spanOverrides == nil {
			loadSpanOverrides()
		}

		if span, ok := spanOverrides[spanID]; ok {
			res = span.Result
			height = span.Height
			spanOverridden = true
		}

		if !spanOverridden {
			// get query params
			queryParams, err := cliCtx.Codec.MarshalJSON(types.NewQuerySpanParams(spanID))
			if err != nil {
				return
			}

			// fetch span
			res, height, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QuerySpan), queryParams)
			if err != nil {
				hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
				return
			}
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No span found"); !ok {
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		hmRest.PostProcessResponse(w, cliCtx, res)
	}
}

// swagger:route GET /bor/latest-span bor borSpanLatest
// It returns the latest-span
// responses:
//
//	200: borSpanResponse
func latestSpanHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// fetch latest span
		res, height, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLatestSpan), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "No latest span found"); !ok {
			return
		}

		// return result
		cliCtx = cliCtx.WithHeight(height)
		hmRest.PostProcessResponse(w, cliCtx, res)
	}
}

//swagger:parameters borPrepareNextSpan
type borPrepareNextSpanParam struct {

	//Start Block
	//required:true
	//type:integer
	//in:query
	StartBlock int `json:"start_block"`

	//Span ID of the span
	//required:true
	//type:integer
	//in:query
	SpanId int `json:"span_id"`

	//Chain ID of the network
	//required:true
	//type:integer
	//in:query
	ChainId int `json:"chain_id"`
}

// swagger:route GET /bor/prepare-next-span bor borPrepareNextSpan
// It returns the prepared next span
// responses:
//
//	200: borSpanResponse
func prepareNextSpanHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		params := r.URL.Query()

		spanID, ok := rest.ParseUint64OrReturnBadRequest(w, params.Get("span_id"))
		if !ok {
			return
		}

		startBlock, ok := rest.ParseUint64OrReturnBadRequest(w, params.Get("start_block"))
		if !ok {
			return
		}

		chainID := params.Get("chain_id")

		//
		// Get span duration
		//

		// fetch duration
		spanDurationBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, types.QueryParams, types.ParamSpan), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, spanDurationBytes, "No span duration"); !ok {
			return
		}

		var spanDuration uint64
		if err := jsoniter.ConfigFastest.Unmarshal(spanDurationBytes, &spanDuration); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		//
		// Get ack count
		//

		// fetch ack count
		ackCountBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", checkpointTypes.QuerierRoute, checkpointTypes.QueryAckCount), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, ackCountBytes, "Ack not found"); !ok {
			return
		}

		var ackCount uint64
		if err := jsoniter.ConfigFastest.Unmarshal(ackCountBytes, &ackCount); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		//
		// Validators
		//

		validatorSetBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", stakingTypes.QuerierRoute, stakingTypes.QueryCurrentValidatorSet), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check content
		if !hmRest.ReturnNotFoundIfNoContent(w, validatorSetBytes, "No current validator set found") {
			return
		}

		var _validatorSet hmTypes.ValidatorSet
		if err = jsoniter.ConfigFastest.Unmarshal(validatorSetBytes, &_validatorSet); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusNoContent, errors.New("unable to unmarshall JSON").Error())
			return
		}

		//
		// Fetching SelectedProducers
		//

		nextProducerBytes, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryNextProducers), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, nextProducerBytes, "Next Producers not found"); !ok {
			return
		}

		var selectedProducers []hmTypes.Validator
		if err := jsoniter.ConfigFastest.Unmarshal(nextProducerBytes, &selectedProducers); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		selectedProducers = hmTypes.SortValidatorByAddress(selectedProducers)

		// draft a propose span message
		msg := hmTypes.NewSpan(
			spanID,
			startBlock,
			startBlock+spanDuration-1,
			_validatorSet,
			selectedProducers,
			chainID,
		)

		result, err := jsoniter.ConfigFastest.Marshal(&msg)
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())

			return
		}

		hmRest.PostProcessResponse(w, cliCtx, result)
	}
}

// swagger:route GET /bor/params bor borSpanParams
// It returns the span parameters
// responses:
//
//	200: borSpanParamsResponse
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

// ResponseWithHeight defines a response object type that wraps an original
// response with a height.
// TODO:Link it with bor
type ResponseWithHeight struct {
	Height string              `json:"height"`
	Result jsoniter.RawMessage `json:"result"`
}

func loadSpanOverrides() {
	spanOverrides = map[uint64]*HeimdallSpanResultWithHeight{}

	j, ok := SPAN_OVERRIDES[helper.GenesisDoc.ChainID]
	if !ok {
		return
	}

	var spans []*types.ResponseWithHeight
	if err := jsoniter.ConfigFastest.Unmarshal(j, &spans); err != nil {
		return
	}

	for _, span := range spans {
		var heimdallSpan types.HeimdallSpan
		if err := jsoniter.ConfigFastest.Unmarshal(span.Result, &heimdallSpan); err != nil {
			continue
		}

		height, err := strconv.ParseInt(span.Height, 10, 64)
		if err != nil {
			continue
		}

		spanOverrides[heimdallSpan.ID] = &HeimdallSpanResultWithHeight{
			Height: height,
			Result: span.Result,
		}
	}
}

//swagger:parameters borSpanList borSpanById borPrepareNextSpan borSpanLatest borSpanParams borNextSpanSeed
type Height struct {

	//Block Height
	//in:query
	Height string `json:"height"`
}
