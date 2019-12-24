package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/bor"
	"github.com/maticnetwork/heimdall/bor/types"
	"github.com/maticnetwork/heimdall/checkpoint"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/staking"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmRest "github.com/maticnetwork/heimdall/types/rest"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// Get span details from start block
	r.HandleFunc("/bor/span/{id}", getSpanHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/bor/latest-span", getLatestSpanHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/bor/span-proposer", getSpanProposersHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/bor/prepare-next-span", prepareNextSpanHandlerFn(cliCtx)).Methods("GET")
}

func getSpanHandlerFn(
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		// get to address
		spanID, ok := rest.ParseUint64OrReturnBadRequest(w, vars["id"])
		if !ok {
			return
		}

		res, _, err := cliCtx.QueryStore(bor.GetSpanKey(spanID), types.StoreKey)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// the query will return empty if there is no data in buffer
		if len(res) == 0 {
			hmRest.WriteErrorResponse(w, http.StatusNotFound, errors.New("No content found for requested span").Error())
			return
		}

		var span hmTypes.Span
		err = cliCtx.Codec.UnmarshalBinaryBare(res, &span)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := json.Marshal(&span)
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		hmRest.PostProcessResponse(w, cliCtx, result)
	}
}

func getLatestSpanHandlerFn(
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//
		// Get latest span start block
		//

		res, _, err := cliCtx.QueryStore(bor.LastSpanIDKey, types.StoreKey)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res); !ok {
			return
		}

		lastestSpanStart, ok := rest.ParseUint64OrReturnBadRequest(w, string(res))
		if !ok {
			return
		}

		//
		// Get latest span
		//

		res, _, err = cliCtx.QueryStore(bor.GetSpanKey(lastestSpanStart), types.StoreKey)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res); !ok {
			return
		}

		var span hmTypes.Span
		err = cliCtx.Codec.UnmarshalBinaryBare(res, &span)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := json.Marshal(&span)
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// return result
		hmRest.PostProcessResponse(w, cliCtx, result)
	}
}

func getSpanProposersHandlerFn(
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// TODO add span proposer logic selection here
		// for now its set as last producer sorted by address in current span

		//
		// Get latest span start block
		//

		res, _, err := cliCtx.QueryStore(bor.LastSpanIDKey, types.StoreKey)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res); !ok {
			return
		}

		lastestSpanID, ok := rest.ParseUint64OrReturnBadRequest(w, string(res))
		if !ok {
			return
		}

		//
		// Get latest span
		//

		res, _, err = cliCtx.QueryStore(bor.GetSpanKey(lastestSpanID), types.StoreKey)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// check content
		if ok := hmRest.ReturnNotFoundIfNoContent(w, res); !ok {
			return
		}

		var span hmTypes.Span
		err = cliCtx.Codec.UnmarshalBinaryBare(res, &span)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// selected producers
		selectedProducers := hmTypes.SortValidatorByAddress(span.SelectedProducers)

		//
		// Get span proposers
		//

		result, err := json.Marshal(&selectedProducers)
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		hmRest.PostProcessResponse(w, cliCtx, result)
	}
}

func prepareNextSpanHandlerFn(
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, bor.QueryParams, bor.ParamSpan), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if len(res) == 0 {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, errors.New("Span duration not found ").Error())
			return
		}
		var spanDuration uint64
		if err := cliCtx.Codec.UnmarshalJSON(res, &spanDuration); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		//
		// Get ack count
		//

		// fetch ack count
		res, _, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s", checkpointTypes.QuerierRoute, checkpoint.QueryAckCount), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if len(res) == 0 {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, errors.New("Ack not found").Error())
			return
		}

		var ackCount uint64
		if err := cliCtx.Codec.UnmarshalJSON(res, &ackCount); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		//
		// Validators
		//

		res, _, err = cliCtx.QueryStore(staking.CurrentValidatorSetKey, "staking")
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		// the query will return empty if there is no data
		if len(res) == 0 {
			hmRest.WriteErrorResponse(w, http.StatusNoContent, errors.New("no content found for requested key").Error())
			return
		}

		var _validatorSet hmTypes.ValidatorSet
		err = cliCtx.Codec.UnmarshalBinaryBare(res, &_validatorSet)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusNoContent, errors.New("unable to unmarshall binary bare").Error())
			return
		}

		// Fetching SelectedProducers

		res, _, err = cliCtx.QueryWithData(fmt.Sprintf("custom/%s/%s/%s", types.QuerierRoute, bor.QueryParams, bor.ParamNextProducers), nil)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		if len(res) == 0 {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, errors.New("Next Producers not found").Error())
			return
		}

		var selectedProducers []hmTypes.Validator
		if err := cliCtx.Codec.UnmarshalJSON(res, &selectedProducers); err != nil {
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

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
