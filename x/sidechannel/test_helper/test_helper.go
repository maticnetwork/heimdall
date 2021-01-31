package test_helper

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/x/sidechannel/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

//
// Create test app
//

// returns context and app with params set on account keeper
func CreateTestApp(isCheckTx bool) (*app.HeimdallApp, sdk.Context, client.Context) {
	genesisState := app.NewDefaultGenesisState()
	sideChannelGenesis := types.NewGenesisState(types.DefaultGenesisState().PastCommits)

	// setup with isCheckTx
	initApp := app.Setup(false)
	ctx := initApp.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	cliCtx := client.Context{}.WithJSONMarshaler(initApp.AppCodec())

	genesisState[types.ModuleName] = initApp.AppCodec().MustMarshalJSON(sideChannelGenesis)
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

	return initApp, ctx, cliCtx
}
