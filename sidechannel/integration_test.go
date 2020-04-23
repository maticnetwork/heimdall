package sidechannel_test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/sidechannel/types"
)

//
// Create test app
//

// returns context and app with params set on account keeper
func createTestApp(isCheckTx bool) (*app.HeimdallApp, sdk.Context) {
	app := app.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})

	return app, ctx
}

func setupWithGenesis() *app.HeimdallApp {
	// setup with isCheckTx
	happ := app.Setup(true)

	// initialize the chain with the passed in genesis accounts
	genesisState := app.NewDefaultGenesisState()

	// past commits
	sidechannelGenesis := types.DefaultGenesisState()
	genesisState[types.ModuleName] = happ.Codec().MustMarshalJSON(sidechannelGenesis)

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
