package checkpoint_test

import (
	"math/rand"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/checkpoint/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/simulation"
	abci "github.com/tendermint/tendermint/abci/types"
)

//
// Create test app
//

// returns context and app
func createTestApp(isCheckTx bool) (*app.HeimdallApp, sdk.Context, context.CLIContext) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	genesisState := app.NewDefaultGenesisState()
	lastNoACK := simulation.RandIntBetween(r1, 1, 5)
	ackCount := simulation.RandIntBetween(r1, 1, 5)

	// create checkpoint BlockHeader
	startBlock := uint64(0)
	endBlock := uint64(256)
	rootHash := hmTypes.HexToHeimdallHash("123")
	proposerAddress := hmTypes.HexToHeimdallAddress("123")
	timestamp := uint64(time.Now().Unix())
	borChainId := "1234"
	checkpointBlockHeader := hmTypes.CreateBlock(
		startBlock,
		endBlock,
		rootHash,
		proposerAddress,
		borChainId,
		timestamp,
	)
	params := types.DefaultParams()

	checkpointBlockHeaders := make([]hmTypes.CheckpointBlockHeader, ackCount)

	for i := range checkpointBlockHeaders {
		checkpointBlockHeaders[i] = checkpointBlockHeader
	}

	checkpointGenesis := types.NewGenesisState(
		params,
		&checkpointBlockHeader,
		uint64(lastNoACK),
		uint64(ackCount),
		checkpointBlockHeaders,
	)

	app := app.Setup(isCheckTx)
	ctx := app.BaseApp.NewContext(isCheckTx, abci.Header{})
	cliCtx := context.NewCLIContext().WithCodec(app.Codec())

	genesisState[types.ModuleName] = app.Codec().MustMarshalJSON(checkpointGenesis)

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
