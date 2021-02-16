package test_helper

import (
	"encoding/json"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	borTypes "github.com/maticnetwork/heimdall/x/bor/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func CreateTestApp(isCheckTx bool) (*app.HeimdallApp, sdk.Context, client.Context) {
	genesisState := app.NewDefaultGenesisState()
	stakingGenesis := borTypes.NewGenesisState(
		borTypes.DefaultParams(),
		borTypes.DefaultGenesisState().Spans,
	)

	initApp := app.Setup(isCheckTx)
	ctx := initApp.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	cliCtx := client.Context{}.WithJSONMarshaler(initApp.AppCodec())

	genesisState[borTypes.ModuleName] = initApp.AppCodec().MustMarshalJSON(stakingGenesis)
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
