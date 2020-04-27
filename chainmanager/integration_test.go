package chainmanager_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/app"
	chainManagerTypes "github.com/maticnetwork/heimdall/chainmanager/types"
)

//
// Create test app
//

// returns context and app with params set on chainmanager keeper
func createTestApp(isCheckTx bool) (*app.HeimdallApp, sdk.Context) {

	app := app.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	app.ChainKeeper.SetParams(ctx, chainManagerTypes.DefaultParams())
	return app, ctx
}
