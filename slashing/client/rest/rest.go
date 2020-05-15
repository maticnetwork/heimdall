package rest

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/gorilla/mux"
	tmLog "github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/helper"
)

// RestLogger for slashing module logger
var RestLogger tmLog.Logger

func init() {
	RestLogger = helper.Logger.With("module", "slashing/rest")
}

// func RegisterHandlers(ctx context.CLIContext, m codec.Marshaler, txg tx.Generator, r *mux.Router) {
// 	registerQueryRoutes(ctx, r)
// 	registerTxHandlers(ctx, m, txg, r)
// }

// RegisterRoutes registers slashing-related REST handlers to a router
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}
