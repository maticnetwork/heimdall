package milestone_test

import (
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/helper/mocks"
	"github.com/maticnetwork/heimdall/milestone"
	chSim "github.com/maticnetwork/heimdall/milestone/simulation"
	"github.com/maticnetwork/heimdall/milestone/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// SideHandlerTestSuite integrate test suite context object
type SideHandlerTestSuite struct {
	suite.Suite

	app            *app.HeimdallApp
	ctx            sdk.Context
	sideHandler    hmTypes.SideTxHandler
	postHandler    hmTypes.PostTxHandler
	contractCaller mocks.IContractCaller
	r              *rand.Rand
}

func (suite *SideHandlerTestSuite) SetupTest() {
	suite.app, suite.ctx, _ = createTestApp(false)
	suite.contractCaller = mocks.IContractCaller{}
	suite.sideHandler = milestone.NewSideTxHandler(suite.app.MilestoneKeeper, &suite.contractCaller)
	suite.postHandler = milestone.NewPostTxHandler(suite.app.MilestoneKeeper, &suite.contractCaller)

	// random generator
	s1 := rand.NewSource(time.Now().UnixNano())
	suite.r = rand.New(s1)
}

func TestSideHandlerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(SideHandlerTestSuite))
}

//
// Test cases
//

func (suite *SideHandlerTestSuite) TestSideHandler() {
	t, ctx := suite.T(), suite.ctx

	// side handler
	result := suite.sideHandler(ctx, nil)
	require.Equal(t, uint32(sdk.CodeUnknownRequest), result.Code)
	require.Equal(t, abci.SideTxResultType_Skip, result.Result)
}

func (suite *SideHandlerTestSuite) TestSideHandleMsgMilestone() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.MilestoneKeeper

	start := uint64(0)
	milestoneLength := helper.GetConfig().MilestoneLength

	milestone, err := chSim.GenRandMilestone(start, milestoneLength)
	require.NoError(t, err)

	borChainId := "1234"

	suite.Run("Success", func() {
		suite.contractCaller = mocks.IContractCaller{}

		// create milestone msg
		msgMilestone := types.NewMsgMilestoneBlock(
			milestone.Proposer,
			milestone.StartBlock,
			milestone.EndBlock,
			milestone.RootHash,
			borChainId,
			milestone.MilestoneID,
		)

		suite.contractCaller.On("CheckIfBlocksExist", milestone.EndBlock).Return(true)
		suite.contractCaller.On("GetVoteOnRootHash", milestone.StartBlock, milestone.EndBlock, milestoneLength, milestone.RootHash.String(), milestone.MilestoneID).Return(true, nil)

		result := suite.sideHandler(ctx, msgMilestone)
		require.Equal(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should be success")

		milestoneReceived, _ := keeper.GetLastMilestone(ctx)
		require.Nil(t, milestoneReceived, "Should not store state")

	})

	suite.Run("No Roothash", func() {
		suite.contractCaller = mocks.IContractCaller{}

		// create milestone msg
		msgMilestone := types.NewMsgMilestoneBlock(
			milestone.Proposer,
			milestone.StartBlock,
			milestone.EndBlock,
			milestone.RootHash,
			borChainId,
			milestone.MilestoneID,
		)

		suite.contractCaller.On("CheckIfBlocksExist", milestone.EndBlock).Return(true)
		suite.contractCaller.On("GetVoteOnRootHash", milestone.StartBlock, milestone.EndBlock, milestoneLength, milestone.RootHash.String(), milestone.MilestoneID).Return(false, nil)

		result := suite.sideHandler(ctx, msgMilestone)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should Fail")
		require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should be `skip`")
		require.Equal(t, uint32(common.CodeInvalidBlockInput), result.Code)

		Header, err := keeper.GetLastMilestone(ctx)
		require.Error(t, err)
		require.Nil(t, Header, "Should not store state")
	})

	suite.Run("invalid milestone", func() {
		suite.contractCaller = mocks.IContractCaller{}

		// create milestone msg
		msgMilestone := types.NewMsgMilestoneBlock(
			milestone.Proposer,
			milestone.StartBlock,
			milestone.EndBlock+1,
			milestone.RootHash,
			borChainId,
			milestone.MilestoneID,
		)

		suite.contractCaller.On("CheckIfBlocksExist", milestone.EndBlock).Return(true)
		suite.contractCaller.On("GetVoteOnRootHash", milestone.StartBlock, milestone.EndBlock, milestoneLength, milestone.RootHash.String(), milestone.MilestoneID).Return(true, nil)

		result := suite.sideHandler(ctx, msgMilestone)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should fail")
		require.Equal(t, uint32(common.CodeInvalidBlockInput), result.Code)
	})

	suite.Run("Not in continuity", func() {
		suite.contractCaller = mocks.IContractCaller{}
		err := keeper.AddMilestone(ctx, milestone)

		if err != nil {
			t.Error("Could add the milestone")
		}

		// create milestone msg
		msgMilestone := types.NewMsgMilestoneBlock(
			milestone.Proposer,
			milestone.StartBlock,
			milestone.EndBlock,
			milestone.RootHash,
			borChainId,
			milestone.MilestoneID,
		)

		suite.contractCaller.On("CheckIfBlocksExist", milestone.EndBlock).Return(true)
		suite.contractCaller.On("GetVoteOnRootHash", milestone.StartBlock, milestone.EndBlock, milestoneLength, milestone.RootHash.String(), milestone.MilestoneID).Return(true, nil)

		result := suite.sideHandler(ctx, msgMilestone)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should fail as milestone is not in continuity to latest stored milestone ")

	})

}

func (suite *SideHandlerTestSuite) TestPostHandler() {
	t, ctx := suite.T(), suite.ctx

	// side handler
	result := suite.postHandler(ctx, nil, abci.SideTxResultType_Yes)
	require.False(t, result.IsOK(), "Post handler should fail")
	require.Equal(t, sdk.CodeUnknownRequest, result.Code)
}

func (suite *SideHandlerTestSuite) TestPostHandleMsgMilestone() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.MilestoneKeeper
	stakingKeeper := app.StakingKeeper
	start := uint64(0)
	milestoneLength := helper.GetConfig().MilestoneLength

	// check valid milestone
	// generate proposer for validator set
	chSim.LoadValidatorSet(t, 2, stakingKeeper, ctx, false, 10)
	stakingKeeper.IncrementAccum(ctx, 1)

	lastMilestone, err := keeper.GetLastMilestone(ctx)
	if err == nil {
		start = start + lastMilestone.EndBlock + 1
	}

	milestone, err := chSim.GenRandMilestone(start, milestoneLength)
	require.NoError(t, err)

	// add current proposer to header
	milestone.Proposer = stakingKeeper.GetValidatorSet(ctx).Proposer.Signer

	borChainId := "1234"

	suite.Run("Failure", func() {
		// create milestone msg
		msgMilestone := types.NewMsgMilestoneBlock(
			milestone.Proposer,
			milestone.StartBlock,
			milestone.EndBlock,
			milestone.RootHash,
			borChainId,
			"00000",
		)

		_ = suite.postHandler(ctx, msgMilestone, abci.SideTxResultType_No)

		lastMilestone, err = keeper.GetLastMilestone(ctx)
		require.Nil(t, lastMilestone)
		require.Error(t, err)

		lastNoAckMilestone := keeper.GetLastNoAckMilestone(ctx)
		require.Equal(t, lastNoAckMilestone, "00000")

		IsNoAckMilestone := keeper.GetNoAckMilestone(ctx, "00000")
		require.True(t, IsNoAckMilestone)

		IsNoAckMilestone = keeper.GetNoAckMilestone(ctx, "WrongID")
		require.False(t, IsNoAckMilestone)

	})

	suite.Run("Success", func() {
		// create milestone msg
		msgMilestone := types.NewMsgMilestoneBlock(
			milestone.Proposer,
			milestone.StartBlock,
			milestone.EndBlock,
			milestone.RootHash,
			borChainId,
			"00001",
		)
		_ = suite.postHandler(ctx, msgMilestone, abci.SideTxResultType_Yes)

		bufferedHeader, err := keeper.GetLastMilestone(ctx)
		require.Equal(t, bufferedHeader.StartBlock, milestone.StartBlock)
		require.Equal(t, bufferedHeader.EndBlock, milestone.EndBlock)
		require.Equal(t, bufferedHeader.RootHash, milestone.RootHash)
		require.Equal(t, bufferedHeader.Proposer, milestone.Proposer)
		require.Equal(t, bufferedHeader.BorChainID, milestone.BorChainID)
		require.Empty(t, err, "Unable to set milestone, Error: %v", err)

		lastNoAckMilestone := keeper.GetLastNoAckMilestone(ctx)
		require.NotEqual(t, lastNoAckMilestone, "00001")

		lastNoAckMilestone = keeper.GetLastNoAckMilestone(ctx)
		require.Equal(t, lastNoAckMilestone, "00000")

		IsNoAckMilestone := keeper.GetNoAckMilestone(ctx, "00001")
		require.False(t, IsNoAckMilestone)

	})

	suite.Run("Pre Exist", func() {
		// create milestone msg
		msgMilestone := types.NewMsgMilestoneBlock(
			milestone.Proposer,
			milestone.StartBlock,
			milestone.EndBlock,
			milestone.RootHash,
			borChainId,
			"00002",
		)
		_ = suite.postHandler(ctx, msgMilestone, abci.SideTxResultType_Yes)

		lastNoAckMilestone := keeper.GetLastNoAckMilestone(ctx)
		require.Equal(t, lastNoAckMilestone, "00002")

		IsNoAckMilestone := keeper.GetNoAckMilestone(ctx, "00002")
		require.True(t, IsNoAckMilestone)

	})

	suite.Run("Not in continuity", func() {
		// create milestone msg
		msgMilestone := types.NewMsgMilestoneBlock(
			milestone.Proposer,
			milestone.StartBlock+64+1,
			milestone.EndBlock+64+1,
			milestone.RootHash,
			borChainId,
			"00003",
		)
		_ = suite.postHandler(ctx, msgMilestone, abci.SideTxResultType_Yes)

		lastNoAckMilestone := keeper.GetLastNoAckMilestone(ctx)
		require.Equal(t, lastNoAckMilestone, "00003")

		IsNoAckMilestone := keeper.GetNoAckMilestone(ctx, "00003")
		require.True(t, IsNoAckMilestone)

	})

	suite.Run("Replay", func() {
		// create milestone msg
		msgMilestone := types.NewMsgMilestoneBlock(
			milestone.Proposer,
			milestone.StartBlock,
			milestone.EndBlock,
			milestone.RootHash,
			borChainId,
			"00004",
		)
		_ = suite.postHandler(ctx, msgMilestone, abci.SideTxResultType_No)
		lastNoAckMilestone := keeper.GetLastNoAckMilestone(ctx)
		require.Equal(t, lastNoAckMilestone, "00004")
	})
}
