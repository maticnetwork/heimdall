package checkpoint_test

import (
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/checkpoint/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

//
// Create test app
//

// createTestApp returns context and app
func createTestApp(isCheckTx bool) (*app.HeimdallApp, sdk.Context, context.CLIContext) {
	genesisState := app.NewDefaultGenesisState()

	app := app.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	cliCtx := context.NewCLIContext().WithCodec(app.Codec())

	params := types.NewParams(5*time.Second, 256, 1024)

	checkpointBlockHeaders := make([]hmTypes.CheckpointBlockHeader, 0)

	for i := range checkpointBlockHeaders {
		checkpointBlockHeaders[i] = hmTypes.CheckpointBlockHeader{}
	}

	checkpointGenesis := types.NewGenesisState(
		params,
		nil,
		uint64(0),
		uint64(0),
		nil,
	)

	genesisState[types.ModuleName] = app.Codec().MustMarshalJSON(checkpointGenesis)

	stateBytes, err := codec.MarshalJSONIndent(app.Codec(), genesisState)
	if err != nil {
		panic(err)
	}

	app.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	app.Commit()
	app.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: app.LastBlockHeight() + 1}})
	app.CheckpointKeeper.SetParams(ctx, params)
	return app, ctx, cliCtx
}
