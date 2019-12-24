package rest

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	"github.com/maticnetwork/heimdall/clerk"
	clerkTypes "github.com/maticnetwork/heimdall/clerk/types"
	hmRest "github.com/maticnetwork/heimdall/types/rest"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
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
		res, _, err := cliCtx.QueryStore(clerk.GetEventRecordKey(recordID), clerkTypes.StoreKey)
		if err != nil {
			hmRest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// the query will return empty if there is no data
		if len(res) == 0 {
			hmRest.WriteErrorResponse(w, http.StatusNoContent, errors.New("no content found for requested key").Error())
			return
		}

		var _record clerkTypes.EventRecord
		err = cdc.UnmarshalBinaryBare(res, &_record)
		if err != nil {
			RestLogger.Error("Error while marshalling state record data", "error", err)
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		result, err := json.Marshal(&_record)
		if err != nil {
			RestLogger.Error("Error while marshalling response to Json", "error", err)
			hmRest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		hmRest.PostProcessResponse(w, cliCtx, result)
	}
}
