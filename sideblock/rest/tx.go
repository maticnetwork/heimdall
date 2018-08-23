package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
	"github.com/cosmos/cosmos-sdk/wire"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	"net/http"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router, cdc *wire.Codec, kb keys.Keybase) {
	r.HandleFunc(
		"/sideblock/submitBlock",
		unrevokeRequestHandlerFn(cdc, kb, cliCtx),
	).Methods("POST")
}
func unrevokeRequestHandlerFn(codec *wire.Codec, keybase keys.Keybase, cliContext context.CLIContext) func(http.ResponseWriter, *http.Request) {
	
}
