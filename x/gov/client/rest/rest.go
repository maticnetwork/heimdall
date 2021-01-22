package rest

import (
	"github.com/gorilla/mux"
	"net/http"

	"github.com/cosmos/cosmos-sdk/client"
	// this line is used by starport scaffolding # 1
)

const (
	MethodGet = "GET"
)

// RegisterRoutes registers gov-related REST handlers to a router
func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
}

// ProposalRESTHandler defines a REST handler implemented in another module. The
// sub-route is mounted on the governance REST handler.
type ProposalRESTHandler struct {
	SubRoute string
	Handler  func(http.ResponseWriter, *http.Request)
}
