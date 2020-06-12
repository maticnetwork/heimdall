package staking_test

import (
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/app"
)

//
// Create test app
//

// returns context and app with params set on staking keeper
func createTestApp(isCheckTx bool) (*app.HeimdallApp, sdk.Context, context.CLIContext) {
	genesisState := app.NewDefaultGenesisState()
	stakingGenesis := stakingTypes.NewGenesisState(
		stakingTypes.DefaultGenesisState().Validators,
		stakingTypes.DefaultGenesisState().CurrentValSet,
		stakingTypes.DefaultGenesisState().StakingSequences)

	app := app.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	cliCtx := context.NewCLIContext().WithCodec(app.Codec())

	genesisState[stakingTypes.ModuleName] = app.Codec().MustMarshalJSON(stakingGenesis)
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

	return app, ctx, cliCtx
}
