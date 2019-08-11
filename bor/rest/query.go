package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/bor"
	"github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/rest"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// Get span details from start block
	r.HandleFunc("/bor/span", getSpanHandlerFn(cdc, cliCtx)).Methods("GET")
	r.HandleFunc("/bor/cache", getCacheFn(cdc, cliCtx)).Methods("GET")
}

func getSpanHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()

		// get start block
		startBlock, ok := rest.ParseUint64OrReturnBadRequest(w, params.Get("start-block"))
		if !ok {
			return
		}

		res, err := cliCtx.QueryStore(bor.GetSpanKey(startBlock), "bor")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// the query will return empty if there is no data in buffer
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNoContent, errors.New("No content found for requested key").Error())
			return
		}

		var span types.Span
		err = cdc.UnmarshalBinaryBare(res, &span)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		RestLogger.Debug("Span fetched", "Span", span.String())
		result, err := json.Marshal(&span)
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, result)
	}
}

func getCacheFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		res, err := cliCtx.QueryStore(bor.SpanCacheKey, "bor")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		fmt.Printf("Result is %v", res)
		// the query will return empty if there is no data in buffer
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNoContent, errors.New("no content found for requested key ").Error())
			return
		}
		result, err := json.Marshal(&res)
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, result)
	}
}
