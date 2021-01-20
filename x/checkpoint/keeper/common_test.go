package keeper_test

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

	initApp := app.Setup(isCheckTx)
	ctx := initApp.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	cliCtx := client.Context{}.WithJSONMarshaler(initApp.AppCodec())

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

	genesisState[types.ModuleName] = initApp.AppCodec().MustMarshalJSON(checkpointGenesis)

	stateBytes, err := json.MarshalIndent(genesisState, "", " ")
	if err != nil {
		panic(err)
	}

	initApp.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)
	initApp.Commit()
	initApp.BeginBlock(abci.RequestBeginBlock{Header: tmproto.Header{Height: initApp.LastBlockHeight() + 1}})
	initApp.CheckpointKeeper.SetParams(ctx, params)
	return initApp, ctx, cliCtx
}
