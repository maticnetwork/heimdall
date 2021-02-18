package keeper_test

import (
	"testing"
	"time"

	"github.com/maticnetwork/heimdall/x/checkpoint/test_helper"

	checkpointKeeper "github.com/maticnetwork/heimdall/x/checkpoint/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app, suite.ctx, _ = test_helper.CreateTestApp(false)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestAddCheckpoint() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := initApp.CheckpointKeeper

	headerBlockNumber := uint64(2000)
	startBlock := uint64(0)
	endBlock := uint64(256)
	rootHash := hmCommonTypes.HexToHeimdallHash("123")
	proposerAddress := hmCommonTypes.HexToHeimdallAddress("123")
	timestamp := uint64(time.Now().Unix())
	borChainId := "1234"

	Checkpoint := hmTypes.CreateBlock(
		startBlock,
		endBlock,
		rootHash,
		proposerAddress,
		borChainId,
		timestamp,
	)
	err := keeper.AddCheckpoint(ctx, headerBlockNumber, Checkpoint)
	require.NoError(t, err)

	result, err := keeper.GetCheckpointByNumber(ctx, headerBlockNumber)
	require.NoError(t, err)
	require.Equal(t, startBlock, result.StartBlock)
	require.Equal(t, endBlock, result.EndBlock)
	require.Equal(t, rootHash.String(), result.RootHash)
	require.Equal(t, borChainId, result.BorChainID)
	require.Equal(t, proposerAddress, hmCommonTypes.HexToHeimdallAddress(result.Proposer))
	require.Equal(t, timestamp, result.TimeStamp)
}

func (suite *KeeperTestSuite) TestGetCheckpointList() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := initApp.CheckpointKeeper

	count := 5

	startBlock := uint64(0)
	endBlock := uint64(0)

	for i := 0; i < count; i++ {
		headerBlockNumber := uint64(i) + 1

		startBlock = startBlock + endBlock
		endBlock = endBlock + uint64(255)
		rootHash := hmCommonTypes.HexToHeimdallHash("123")
		proposerAddress := hmCommonTypes.HexToHeimdallAddress("123")
		timestamp := uint64(time.Now().Unix()) + uint64(i)
		borChainId := "1234"

		Checkpoint := hmTypes.CreateBlock(
			startBlock,
			endBlock,
			rootHash,
			proposerAddress,
			borChainId,
			timestamp,
		)

		err := keeper.AddCheckpoint(ctx, headerBlockNumber, Checkpoint)
		require.NoError(t, err)
		keeper.UpdateACKCount(ctx)
	}

	result, err := keeper.GetCheckpointList(ctx, uint64(1), uint64(20))
	require.NoError(t, err)
	require.LessOrEqual(t, count, len(result))
}

func (suite *KeeperTestSuite) TestHasStoreValue() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := initApp.CheckpointKeeper
	key := checkpointKeeper.ACKCountKey
	result := keeper.HasStoreValue(ctx, key)
	require.True(t, result)
}

func (suite *KeeperTestSuite) TestFlushCheckpointBuffer() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := initApp.CheckpointKeeper
	key := checkpointKeeper.BufferCheckpointKey
	keeper.FlushCheckpointBuffer(ctx)
	result := keeper.HasStoreValue(ctx, key)
	require.False(t, result)
}
