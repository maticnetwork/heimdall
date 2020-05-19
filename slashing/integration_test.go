package slashing_test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/slashing/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

//
// Create test app
//

// returns context and app on clerk keeper
func createTestApp(isCheckTx bool) (*app.HeimdallApp, sdk.Context) {
	app := app.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	return app, ctx
}

// setupSlashingGenesis initializes a new Heimdall with the default genesis data.
func setupSlashingGenesis() *app.HeimdallApp {
	happ := app.Setup(true)

	// initialize the chain with the default genesis state
	genesisState := app.NewDefaultGenesisState()
	slashingInfoList := []*hmTypes.ValidatorSlashingInfo{}
	slashingGenesis := types.NewGenesisState(types.DefaultGenesisState().Params, types.DefaultGenesisState().SigningInfos, types.DefaultGenesisState().MissedBlocks, slashingInfoList, slashingInfoList, uint64(0))
	genesisState[types.ModuleName] = happ.Codec().MustMarshalJSON(slashingGenesis)

	stateBytes, err := codec.MarshalJSONIndent(happ.Codec(), genesisState)
	if err != nil {
		panic(err)
	}

	happ.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)

	happ.Commit()
	happ.BeginBlock(abci.RequestBeginBlock{Header: abci.Header{Height: happ.LastBlockHeight() + 1}})

	return happ
}
