package topup_test

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/types/simulation"
)

//
// Create test app
//

// returns context and app with params set on chainmanager keeper
func createTestApp(isCheckTx bool) (*app.HeimdallApp, sdk.Context, context.CLIContext) {
	app := app.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	cliCtx := context.NewCLIContext().WithCodec(app.Codec())

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	topupSequence := strconv.Itoa(simulation.RandIntBetween(r1, 1000, 100000))
	app.TopupKeeper.SetTopupSequence(ctx, topupSequence)

	return app, ctx, cliCtx
}
