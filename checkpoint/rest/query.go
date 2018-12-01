package rest

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/gorilla/mux"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/types"
	"net/http"
)

func registerQueryRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	// Get all delegations from a delegator
	r.HandleFunc(
		"/checkpoint/buffer",
		CheckpointBufferHandlerFn(cdc, cliCtx),
	).Methods("GET")
}

// query accountREST Handler
func CheckpointBufferHandlerFn(
	cdc *codec.Codec,
	cliCtx context.CLIContext,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		res, err := cliCtx.QueryStore(common.BufferCheckpointKey, "checkpoint")
		if err != nil {
			utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		// the query will return empty if there is no data for this account
		if len(res) == 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// decode the value
		//account, err := decoder(res)
		//if err != nil {
		//	utils.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
		//	return
		//}

		//utils.PostProcessResponse(w, cdc, res, cliCtx.Indent)
		var _checkpoint types.CheckpointBlockHeader

		err = json.Unmarshal(res, &_checkpoint)
		if err != nil {
			RestLogger.Info("Unable to marshall")
		}
		result, err := json.Marshal(&_checkpoint)
		if err != nil {
			RestLogger.Error("Error while marshalling tendermint response", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write(result)
	}
}
