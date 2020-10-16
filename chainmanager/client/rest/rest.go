package rest

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/gorilla/mux"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
// RegisterRoutes registers the auth module REST routes.
func RegisterRoutes(cliCtx client.Context, r *mux.Router) {
	r.HandleFunc("/chainmanager/params", paramsHandlerFn(cliCtx)).Methods("GET")
}
