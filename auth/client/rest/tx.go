package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
// RegisterRoutes registers the auth module REST routes.
func RegisterRoutes(cliCtx client.Context, r *mux.Router, storeName string) {
	r.HandleFunc("/auth/accounts/{address}", QueryAccountRequestHandlerFn(storeName, cliCtx)).Methods("GET")
	r.HandleFunc("/auth/accounts/{address}/sequence", QueryAccountSequenceRequestHandlerFn(cliCtx)).Methods("GET")
	r.HandleFunc("/auth/params", queryParamsHandler(cliCtx)).Methods("GET")
}
