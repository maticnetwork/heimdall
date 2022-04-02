package bor_test

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	abci "github.com/tendermint/tendermint/abci/types"
)

//
// Create test app
//

// createTestApp returns context and app
func createTestApp(isCheckTx bool) (*app.HeimdallApp, sdk.Context, context.CLIContext) {
	app := app.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	cliCtx := context.NewCLIContext().WithCodec(app.Codec())
	return app, ctx, cliCtx
}
