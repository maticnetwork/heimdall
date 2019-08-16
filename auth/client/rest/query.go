package rest

import (
	"errors"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/rest"
)

// QueryAccountRequestHandlerFn query account REST Handler
func QueryAccountRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)

		// key
		key := types.HexToHeimdallAddress(vars["address"])
		if key.Empty() {
			rest.WriteErrorResponse(w, http.StatusNotFound, errors.New("Invalid address").Error())
			return
		}

		// account getter
		accGetter := authTypes.NewAccountRetriever(cliCtx)

		if err := accGetter.EnsureExists(key); err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		account, err := accGetter.GetAccount(key)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, account)
	}
}

// QueryAccountSequenceRequestHandlerFn query accoun sequence REST Handler
func QueryAccountSequenceRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)

		// key
		key := types.HexToHeimdallAddress(vars["address"])
		if key.Empty() {
			rest.WriteErrorResponse(w, http.StatusNotFound, errors.New("Invalid address").Error())
			return
		}

		// account getter
		accGetter := authTypes.NewAccountRetriever(cliCtx)

		if err := accGetter.EnsureExists(key); err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		account, err := accGetter.GetAccount(key)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		// get result
		result := authTypes.LightBaseAccount{
			Address:       account.GetAddress(),
			Sequence:      account.GetSequence(),
			AccountNumber: account.GetAccountNumber(),
		}

		rest.PostProcessResponse(w, cliCtx, result)
	}
}
