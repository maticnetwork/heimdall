package topup_test

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/types/simulation"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

//
// Create test app
//

// returns context and app with params set on chainmanager keeper
func createTestApp(isCheckTx bool) (*app.HeimdallApp, sdk.Context, client.Context) {
	app := app.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	cliCtx := client.Context{}.WithJSONMarshaler(app.AppCodec())

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	topupSequence := strconv.Itoa(simulation.RandIntBetween(r1, 1000, 100000))
	app.TopupKeeper.SetTopupSequence(ctx, topupSequence)
	return app, ctx, cliCtx
}
