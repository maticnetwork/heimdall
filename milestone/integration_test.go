package milestone_test

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/milestone/types"
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

	helper.SetTestConfig(helper.GetDefaultHeimdallConfig())

	params := types.NewParams(64)

	milestoneGenesis := types.NewGenesisState(
		types.DefaultGenesisState().Params,
		types.DefaultGenesisState().Milestones,
		types.DefaultGenesisState().NoAckMilestones,
	)

	genesisState[types.ModuleName] = app.Codec().MustMarshalJSON(milestoneGenesis)

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
	app.MilestoneKeeper.SetParams(ctx, params)

	return app, ctx, cliCtx
}
