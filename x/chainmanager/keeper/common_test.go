package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	chainManagerTypes "github.com/maticnetwork/heimdall/x/chainmanager/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/maticnetwork/heimdall/app"
)

//
// Create test app
//

// returns context and app with params set on chainmanager keeper
func createTestApp(isCheckTx bool) (*app.HeimdallApp, sdk.Context) {

	initApp := app.Setup(isCheckTx)
	ctx := initApp.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	initApp.ChainKeeper.SetParams(ctx, chainManagerTypes.DefaultParams())
	return initApp, ctx
}
