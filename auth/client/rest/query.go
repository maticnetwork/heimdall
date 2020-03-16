package rest

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/gorilla/mux"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/types"
	hmRest "github.com/maticnetwork/heimdall/types/rest"
)

// QueryAccountRequestHandlerFn query account REST Handler
func QueryAccountRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)

		// key
		key := types.HexToHeimdallAddress(vars["address"])
		if key.Empty() {
			hmRest.WriteErrorResponse(w, http.StatusNotFound, errors.New("Invalid address").Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// account getter
		accGetter := authTypes.NewAccountRetriever(cliCtx)

		account, height, err := accGetter.GetAccountWithHeight(key)
		if err != nil {
			if err := accGetter.EnsureExists(key); err != nil {
				cliCtx = cliCtx.WithHeight(height)
				hmRest.PostProcessResponse(w, cliCtx, authTypes.BaseAccount{})
				return
			}
			hmRest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		hmRest.PostProcessResponse(w, cliCtx, account)
	}
}

// QueryAccountSequenceRequestHandlerFn query account sequence REST Handler
func QueryAccountSequenceRequestHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)

		// key
		key := types.HexToHeimdallAddress(vars["address"])
		if key.Empty() {
			hmRest.WriteErrorResponse(w, http.StatusNotFound, errors.New("Invalid address").Error())
			return
		}

		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		// account getter
		accGetter := authTypes.NewAccountRetriever(cliCtx)

		account, height, err := accGetter.GetAccountWithHeight(key)
		if err != nil {
			if err := accGetter.EnsureExists(key); err != nil {
				cliCtx = cliCtx.WithHeight(height)
				hmRest.PostProcessResponse(w, cliCtx, authTypes.BaseAccount{})
				return
			}
			hmRest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		// get result
		result := authTypes.LightBaseAccount{
			Address:       account.GetAddress(),
			Sequence:      account.GetSequence(),
			AccountNumber: account.GetAccountNumber(),
		}

		cliCtx = cliCtx.WithHeight(height)
		hmRest.PostProcessResponse(w, cliCtx, result)
	}
}

// HTTP request handler to query the auth params values
func paramsHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cliCtx, ok := rest.ParseQueryHeightOrReturnBadRequest(w, cliCtx, r)
		if !ok {
			return
		}

		route := fmt.Sprintf("custom/%s/%s", authTypes.QuerierRoute, authTypes.QueryParams)
		res, height, err := cliCtx.QueryWithData(route, nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}

		cliCtx = cliCtx.WithHeight(height)
		rest.PostProcessResponse(w, cliCtx, res)
	}
}
