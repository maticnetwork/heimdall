package checkpoint_test

import (
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/types/simulation"
	"github.com/maticnetwork/heimdall/x/checkpoint"
	"github.com/maticnetwork/heimdall/x/checkpoint/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type GenesisTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
}

// SetupTest setup necessary things for genesis test
func (suite *GenesisTestSuite) SetupTest() {
	suite.app, suite.ctx, _ = createTestApp(true)
}

// TestGenesisTestSuite
func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) TestInitExportGenesis() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	lastNoACK := simulation.RandIntBetween(r1, 1, 5)
	ackCount := simulation.RandIntBetween(r1, 1, 5)
	startBlock := uint64(0)
	endBlock := uint64(256)
	rootHash := hmCommonTypes.HexToHeimdallHash("123")

	proposerAddress := hmCommonTypes.HexToHeimdallAddress("123")
	timestamp := uint64(time.Now().Unix())
	borChainId := "1234"

	bufferedCheckpoint := hmTypes.CreateBlock(
		startBlock,
		endBlock,
		rootHash,
		proposerAddress,
		borChainId,
		timestamp,
	)

	checkpoints := make([]*hmTypes.Checkpoint, ackCount)

	for i := range checkpoints {
		checkpoints[i] = bufferedCheckpoint
	}

	params := types.DefaultParams()
	genesisState := types.NewGenesisState(
		params,
		bufferedCheckpoint,
		uint64(lastNoACK),
		uint64(ackCount),
		checkpoints,
	)

	checkpoint.InitGenesis(ctx, initApp.CheckpointKeeper, genesisState)

	actualParams := checkpoint.ExportGenesis(ctx, initApp.CheckpointKeeper)

	require.Equal(t, genesisState.AckCount, actualParams.AckCount)
	require.Equal(t, genesisState.BufferedCheckpoint, actualParams.BufferedCheckpoint)
	require.Equal(t, genesisState.LastNoACK, actualParams.LastNoACK)
	require.Equal(t, genesisState.Params, actualParams.Params)
	require.LessOrEqual(t, len(actualParams.Checkpoints), len(genesisState.Checkpoints))

}
