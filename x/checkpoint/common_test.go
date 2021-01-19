package checkpoint_test

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client"
	abci "github.com/tendermint/tendermint/abci/types"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/x/checkpoint/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

//
// Create test app
//

// createTestApp returns context and app
func createTestApp(isCheckTx bool) (*app.HeimdallApp, sdk.Context, client.Context) {
	genesisState := app.NewDefaultGenesisState()

	app := app.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	cliCtx := client.Context{}.WithJSONMarshaler(app.AppCodec())

	helper.SetTestConfig(helper.GetDefaultHeimdallConfig())

	params := types.NewParams(5*time.Second, 256, 1024, 10000)

	Checkpoints := make([]hmTypes.Checkpoint, 0)

	for i := range Checkpoints {
		Checkpoints[i] = hmTypes.Checkpoint{}
	}

	checkpointGenesis := types.NewGenesisState(
		types.DefaultGenesis().Params,
		types.DefaultGenesis().BufferedCheckpoint,
		types.DefaultGenesis().LastNoACK,
		types.DefaultGenesis().AckCount,
		types.DefaultGenesis().Checkpoints,
	)

	genesisState[types.ModuleName] = app.AppCodec().MustMarshalJSON(&checkpointGenesis)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
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
	app.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: app.LastBlockHeight() + 1}})
	app.CheckpointKeeper.SetParams(ctx, params)
	return app, ctx, cliCtx
}
