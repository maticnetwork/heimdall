package checkpoint_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app, suite.ctx, _ = createTestApp(false)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestAddCheckpoint() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper

	headerBlockNumber := uint64(2000)
	startBlock := uint64(0)
	endBlock := uint64(256)
	rootHash := hmTypes.HexToHeimdallHash("123")
	accountRootHash := hmTypes.HexToHeimdallHash("456")
	proposerAddress := hmTypes.HexToHeimdallAddress("123")
	timestamp := uint64(time.Now().Unix())

	checkpointBlockHeader := hmTypes.CreateBlock(
		startBlock,
		endBlock,
		rootHash,
		accountRootHash,
		proposerAddress,
		timestamp,
	)
	err := keeper.AddCheckpoint(ctx, headerBlockNumber, checkpointBlockHeader)
	require.NoError(t, err)

	result, err := keeper.GetCheckpointByIndex(ctx, headerBlockNumber)
	require.NoError(t, err)
	require.Equal(t, startBlock, result.StartBlock)
	require.Equal(t, endBlock, result.EndBlock)
	require.Equal(t, rootHash, result.RootHash)
	require.Equal(t, accountRootHash, result.AccountRootHash)
	require.Equal(t, proposerAddress, result.Proposer)
	require.Equal(t, timestamp, result.TimeStamp)
}
