package keeper_test

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/x/gov/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/maticnetwork/heimdall/app"
)

// TODO - Merge this with keeper/common_test.go
//
// Create test app
//

// returns context and app with params set on gov keeper
func createTestApp(isCheckTx bool) (*app.HeimdallApp, sdk.Context, client.Context) {
	genesisState := app.NewDefaultGenesisState()
	govGenesis := types.NewGenesisState(types.DefaultGenesis().StartingProposalId, types.DefaultGenesis().DepositParams, types.DefaultGenesis().VotingParams, types.DefaultGenesis().TallyParams)

	app := app.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	cliCtx := client.Context{}.WithJSONMarshaler(app.AppCodec())

	genesisState[types.ModuleName] = app.AppCodec().MustMarshalJSON(&govGenesis)
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

	return app, ctx, cliCtx
}
