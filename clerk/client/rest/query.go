package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/clerk"
	clerkTypes "github.com/maticnetwork/heimdall/clerk/types"
	"github.com/maticnetwork/heimdall/types/rest"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// Get all delegations from a delegator
	r.HandleFunc(
		"/clerk/event-record/{recordId}",
		handlerRecordFn(cdc, cliCtx),
	).Methods("GET")
}

// handlerRecordFn returns record by record id
func handlerRecordFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		// record id
		recordID, ok := rest.ParseUint64OrReturnBadRequest(w, vars["recordId"])
		if !ok {
			return
		}

		// get record from store
		res, err := cliCtx.QueryStore(clerk.GetEventRecordKey(recordID), clerkTypes.StoreKey)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// the query will return empty if there is no data
		if len(res) == 0 {
			rest.WriteErrorResponse(w, http.StatusNoContent, errors.New("no content found for requested key").Error())
			return
		}

		var _record clerkTypes.EventRecord
		err = cdc.UnmarshalBinaryBare(res, &_record)
		if err != nil {
			RestLogger.Error("Error while marshalling state record data", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := json.Marshal(&_record)
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, result)
	}
}
