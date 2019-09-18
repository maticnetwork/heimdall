package tx

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
)

// register REST routes
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *codec.Codec) {
	r.HandleFunc("/txs/{hash}", QueryTxRequestHandlerFn(cliCtx, cdc)).Methods("GET")
	r.HandleFunc("/txs/{hash}/commit-proof", QueryCommitTxRequestHandlerFn(cliCtx, cdc)).Methods("GET")
	r.HandleFunc("/txs", QueryTxsByTagsRequestHandlerFn(cliCtx, cdc)).Methods("GET")
	r.HandleFunc("/txs", BroadcastTxRequest(cliCtx, cdc)).Methods("POST")
	r.HandleFunc("/txs/encode", EncodeTxRequestHandlerFn(cliCtx, cdc)).Methods("POST")
}
