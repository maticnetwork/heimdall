package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/ethereum/go-ethereum/consensus/bor"
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

var spanOverrides map[uint64]*HeimdallSpanResultWithHeight = nil

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc("/bor/span/list", spanListHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/bor/span/{id}", spanHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/bor/latest-span", latestSpanHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/bor/prepare-next-span", prepareNextSpanHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/bor/next-span-seed", fetchNextSpanSeedHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/bor/params", paramsHandlerFn(cliCtx)).Methods("GET")
}

func fetchNextSpanSeedHandlerFn(
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryNextSpanSeed), nil)
		RestLogger.Debug("nextSpanSeed querier response", "res", res)

		if err != nil {
			RestLogger.Error("Error while fetching next span seed  ", "Error", err.Error())
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// error if span seed found
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res, "NextSpanSeed not found"); !ok {
			RestLogger.Error("NextSpanSeed not found ", "Error", err.Error())
			return
		}

		// return result
		rest.PostProcessResponse(w, cliCtx, res)

	}
}

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
		if err := json.Unmarshal(spanDurationBytes, &spanDuration); err != nil {
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
		if err := json.Unmarshal(ackCountBytes, &ackCount); err != nil {
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
		if ok := hmRest.ReturnNotFoundIfNoContent(w, validatorSetBytes, "No current validator set found"); !ok {
			return
		}

		var _validatorSet hmTypes.ValidatorSet
		err = json.Unmarshal(validatorSetBytes, &_validatorSet)
		if err != nil {
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
		if err := json.Unmarshal(nextProducerBytes, &selectedProducers); err != nil {
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

		result, err := json.Marshal(&msg)
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		hmRest.PostProcessResponse(w, cliCtx, result)
	}
}

// HTTP request handler to query the bor params values
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

func loadSpanOverrides() {
	spanOverrides = map[uint64]*HeimdallSpanResultWithHeight{}

	j, ok := SPAN_OVERRIDES[helper.GenesisDoc.ChainID]
	if !ok {
		return
	}

	var spans []*bor.ResponseWithHeight
	if err := json.Unmarshal(j, &spans); err != nil {
		return
	}

	for _, span := range spans {
		var heimdallSpan bor.HeimdallSpan
		if err := json.Unmarshal(span.Result, &heimdallSpan); err != nil {
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
