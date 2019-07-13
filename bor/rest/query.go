package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/types"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// Get span details from start block
	r.HandleFunc(
		"/bor/buffer",
		getSpanHandlerFn(cdc, cliCtx),
	).Methods("GET")
}

func getSpanHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		startBlock, err := strconv.ParseUint(vars["startBlock"], 10, 64)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		res, err := cliCtx.QueryStore(common.GetSpanKey(startBlock), "bor")
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		// the query will return empty if there is no data in buffer
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNoContent, err.Error())
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
		rest.PostProcessResponse(w, cdc, result, cliCtx.Indent)
	}
}
