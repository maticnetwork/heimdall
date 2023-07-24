package bank_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/app"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	bankTypes "github.com/maticnetwork/heimdall/bank/types"
)

//
// Create test app
//

// returns context and app with params set on account keeper
func createTestApp(isCheckTx bool) (*app.HeimdallApp, sdk.Context) {
	app := app.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	app.AccountKeeper.SetParams(ctx, authTypes.DefaultParams())
	app.BankKeeper.SetSendEnabled(ctx, bankTypes.DefaultSendEnabled)

	return app, ctx
}
