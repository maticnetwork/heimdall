package checkpoint_test

import (
	"testing"
	"time"

	cmn "github.com/maticnetwork/heimdall/common"
	hmTypes "github.com/maticnetwork/heimdall/types"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestKeeperTestSuiteMilestone(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(KeeperTestSuite))
}

func (suite *KeeperTestSuite) TestAddMilestone() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper

	startBlock := uint64(0)
	endBlock := uint64(63)
	hash := hmTypes.HexToHeimdallHash("123")
	proposerAddress := hmTypes.HexToHeimdallAddress("123")
	timestamp := uint64(time.Now().Unix())
	borChainId := "1234"
	milestoneID := "0000"

	milestone := hmTypes.CreateMilestone(
		startBlock,
		endBlock,
		hash,
		proposerAddress,
		borChainId,
		milestoneID,
		timestamp,
	)
	err := keeper.AddMilestone(ctx, milestone)
	require.NoError(t, err)

	result, err := keeper.GetLastMilestone(ctx)
	require.NoError(t, err)
	require.Equal(t, startBlock, result.StartBlock)
	require.Equal(t, endBlock, result.EndBlock)
	require.Equal(t, hash, result.Hash)
	require.Equal(t, borChainId, result.BorChainID)
	require.Equal(t, proposerAddress, result.Proposer)
	require.Equal(t, timestamp, result.TimeStamp)

	result, err = keeper.GetMilestoneByNumber(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, startBlock, result.StartBlock)
	require.Equal(t, endBlock, result.EndBlock)
	require.Equal(t, hash, result.Hash)
	require.Equal(t, borChainId, result.BorChainID)
	require.Equal(t, proposerAddress, result.Proposer)
	require.Equal(t, timestamp, result.TimeStamp)

	result, err = keeper.GetMilestoneByNumber(ctx, 2)
	require.Nil(t, result)
	require.Equal(t, err, cmn.ErrInvalidMilestoneIndex(keeper.Codespace()))
}

func (suite *KeeperTestSuite) TestPruneMilestone() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper

	startBlock := uint64(0)
	endBlock := uint64(63)
	hash := hmTypes.HexToHeimdallHash("123")
	proposerAddress := hmTypes.HexToHeimdallAddress("123")
	timestamp := uint64(time.Now().Unix())
	borChainId := "1234"
	milestoneID := "0000"

	milestone := hmTypes.CreateMilestone(
		startBlock,
		endBlock,
		hash,
		proposerAddress,
		borChainId,
		milestoneID,
		timestamp,
	)
	err := keeper.AddMilestone(ctx, milestone)
	require.NoError(t, err)

	keeper.PruneMilestone(ctx, 2)
	result, err := keeper.GetMilestoneByNumber(ctx, 1)
	require.NotNil(t, result)
	require.Nil(t, err)

	keeper.PruneMilestone(ctx, 1)
	result, err = keeper.GetMilestoneByNumber(ctx, 1)
	require.Nil(t, result)
	require.Equal(t, err, cmn.ErrInvalidMilestoneIndex(keeper.Codespace()))
}

func (suite *KeeperTestSuite) TestGetCount() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper

	result := keeper.GetMilestoneCount(ctx)
	require.Equal(t, uint64(0), result)

	startBlock := uint64(0)
	endBlock := uint64(63)
	hash := hmTypes.HexToHeimdallHash("123")
	proposerAddress := hmTypes.HexToHeimdallAddress("123")
	timestamp := uint64(time.Now().Unix())
	borChainId := "1234"
	milestoneID := "0000"

	milestone := hmTypes.CreateMilestone(
		startBlock,
		endBlock,
		hash,
		proposerAddress,
		borChainId,
		milestoneID,
		timestamp,
	)
	err := keeper.AddMilestone(ctx, milestone)
	require.NoError(t, err)

	result = keeper.GetMilestoneCount(ctx)
	require.Equal(t, uint64(1), result)
}

func (suite *KeeperTestSuite) TestGetNoAckMilestone() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper

	result := keeper.GetMilestoneCount(ctx)
	require.Equal(t, uint64(0), result)

	milestoneID := "0000"

	keeper.SetNoAckMilestone(ctx, milestoneID)

	val := keeper.GetNoAckMilestone(ctx, "0000")
	require.True(t, val)

	val = keeper.GetNoAckMilestone(ctx, "00001")
	require.False(t, val)

	val = keeper.GetNoAckMilestone(ctx, "")
	require.False(t, val)

	milestoneID = "0001"
	keeper.SetNoAckMilestone(ctx, milestoneID)

	val = keeper.GetNoAckMilestone(ctx, "0001")
	require.True(t, val)

	val = keeper.GetNoAckMilestone(ctx, "0000")
	require.True(t, val)
}

func (suite *KeeperTestSuite) TestLastNoAckMilestone() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper

	result := keeper.GetMilestoneCount(ctx)
	require.Equal(t, uint64(0), result)

	milestoneID := "0000"

	val := keeper.GetLastNoAckMilestone(ctx)
	require.NotEqual(t, val, milestoneID)

	keeper.SetNoAckMilestone(ctx, milestoneID)

	val = keeper.GetLastNoAckMilestone(ctx)
	require.Equal(t, val, milestoneID)

	milestoneID = "0001"

	keeper.SetNoAckMilestone(ctx, milestoneID)

	val = keeper.GetLastNoAckMilestone(ctx)
	require.Equal(t, val, milestoneID)
}

func (suite *KeeperTestSuite) TestGetMilestoneTimeout() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper

	val := keeper.GetLastMilestoneTimeout(ctx)
	require.Zero(t, val)

	keeper.SetLastMilestoneTimeout(ctx, uint64(21))

	val = keeper.GetLastMilestoneTimeout(ctx)
	require.Equal(t, uint64(21), val)
}
