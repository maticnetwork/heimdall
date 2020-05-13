package clerk_test

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/clerk/types"
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

// setupClerkGenesis initializes a new Heimdall with the default genesis data.
func setupClerkGenesis() *app.HeimdallApp {
	happ := app.Setup(true)

	// initialize the chain with the default genesis state
	genesisState := app.NewDefaultGenesisState()

	clerkGenesis := types.NewGenesisState(types.DefaultGenesisState().EventRecords, types.DefaultGenesisState().RecordSequences)
	genesisState[types.ModuleName] = happ.Codec().MustMarshalJSON(clerkGenesis)

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
