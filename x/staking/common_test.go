package staking_test

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingTypes "github.com/maticnetwork/heimdall/x/staking/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	"github.com/maticnetwork/heimdall/app"
)

//
// Create test app
//

// returns context and app with params set on staking keeper
func createTestApp(isCheckTx bool) (*app.HeimdallApp, sdk.Context, client.Context) {
	genesisState := app.NewDefaultGenesisState()
	stakingGenesis := stakingTypes.NewGenesisState(
		stakingTypes.DefaultGenesis().Validators,
		stakingTypes.DefaultGenesis().CurrentValSet,
		stakingTypes.DefaultGenesis().StakingSequences)

	initApp := app.Setup(isCheckTx)
	ctx := initApp.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	cliCtx := client.Context{}.WithJSONMarshaler(initApp.AppCodec())

	genesisState[stakingTypes.ModuleName] = initApp.AppCodec().MustMarshalJSON(stakingGenesis)
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
