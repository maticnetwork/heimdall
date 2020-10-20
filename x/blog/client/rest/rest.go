package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/maticnetwork/heimdall/x/blog/types"
)

const (
	MethodGet = "GET"
)

// RegisterRoutes registers blog-related REST handlers to a router
func RegisterRoutes(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 2
	registerQueryRoutes(clientCtx, r)
	registerTxHandlers(clientCtx, r)

	registerQueryRoutes(clientCtx, r)
	registerTxHandlers(clientCtx, r)

}

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 3
	r.HandleFunc("custom/heimdall/"+types.QueryListComment, listCommentHandler(clientCtx)).Methods("GET")

	r.HandleFunc("custom/heimdall/"+types.QueryListPost, listPostHandler(clientCtx)).Methods("GET")

}

func registerTxHandlers(clientCtx client.Context, r *mux.Router) {
	// this line is used by starport scaffolding # 4
	r.HandleFunc("/heimdall/comment", createCommentHandler(clientCtx)).Methods("POST")

	r.HandleFunc("/heimdall/post", createPostHandler(clientCtx)).Methods("POST")

}
