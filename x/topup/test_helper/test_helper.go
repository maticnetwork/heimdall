package test_helper

import (
	"encoding/json"
	"math/rand"
	"strconv"
	"time"

	chainManagerTypes "github.com/maticnetwork/heimdall/x/chainmanager/types"

	"github.com/maticnetwork/heimdall/types/simulation"
	abci "github.com/tendermint/tendermint/abci/types"

	topupTypes "github.com/maticnetwork/heimdall/x/topup/types"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

//
// Create test app
//

// returns context and app with params set on chainmanager keeper
func CreateTestApp(isCheckTx bool) (*app.HeimdallApp, sdk.Context, client.Context) {
	genesisState := app.NewDefaultGenesisState()
	topUpGenesis := topupTypes.NewGenesisState(
		topupTypes.DefaultGenesis().TopupSequences,
		topupTypes.DefaultGenesis().DividendAccounts)

	initApp := app.Setup(false)
	ctx := initApp.BaseApp.NewContext(isCheckTx, tmproto.Header{})
	cliCtx := client.Context{}.WithJSONMarshaler(initApp.AppCodec())

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	genesisState[topupTypes.ModuleName] = initApp.AppCodec().MustMarshalJSON(&topUpGenesis)
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
	//
	topupSequence := strconv.Itoa(simulation.RandIntBetween(r1, 1000, 100000))
	initApp.TopupKeeper.SetTopupSequence(ctx, topupSequence)
	initApp.ChainKeeper.SetParams(ctx, chainManagerTypes.DefaultParams())

	return initApp, ctx, cliCtx
}
