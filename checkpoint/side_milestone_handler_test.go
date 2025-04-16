package checkpoint_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	cmTypes "github.com/maticnetwork/heimdall/chainmanager/types"
	chSim "github.com/maticnetwork/heimdall/checkpoint/simulation"
	"github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/helper/mocks"
)

func TestMilestoneSideHandlerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(SideHandlerTestSuite))
}

func (suite *SideHandlerTestSuite) TestSideHandleMsgMilestone() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper

	start := uint64(0)
	milestoneLength := helper.MilestoneLength

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
			milestone.Hash,
			borChainId,
			milestone.MilestoneID,
		)

		suite.contractCaller.On("CheckIfBlocksExist", milestone.EndBlock+cmTypes.DefaultMaticchainMilestoneTxConfirmations).Return(true)
		suite.contractCaller.On("GetVoteOnHash", milestone.StartBlock, milestone.EndBlock, milestoneLength, milestone.Hash.String(), milestone.MilestoneID).Return(true, nil)

		result := suite.sideHandler(ctx, msgMilestone)
		require.Equal(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should be success")

		milestoneReceived, _ := keeper.GetLastMilestone(ctx)
		require.Nil(t, milestoneReceived, "Should not store state")

	})

	suite.Run("No Hash", func() {
		suite.contractCaller = mocks.IContractCaller{}

		// create milestone msg
		msgMilestone := types.NewMsgMilestoneBlock(
			milestone.Proposer,
			milestone.StartBlock,
			milestone.EndBlock,
			milestone.Hash,
			borChainId,
			milestone.MilestoneID,
		)

		suite.contractCaller.On("CheckIfBlocksExist", milestone.EndBlock+cmTypes.DefaultMaticchainMilestoneTxConfirmations).Return(true)
		suite.contractCaller.On("GetVoteOnHash", milestone.StartBlock, milestone.EndBlock, milestoneLength, milestone.Hash.String(), milestone.MilestoneID).Return(false, nil)

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
			milestone.EndBlock-1,
			milestone.Hash,
			borChainId,
			milestone.MilestoneID,
		)

		suite.contractCaller.On("CheckIfBlocksExist", milestone.EndBlock+cmTypes.DefaultMaticchainMilestoneTxConfirmations).Return(true)
		suite.contractCaller.On("GetVoteOnHash", milestone.StartBlock, milestone.EndBlock, milestoneLength, milestone.Hash.String(), milestone.MilestoneID).Return(true, nil)

		result := suite.sideHandler(ctx, msgMilestone)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should fail")
		require.Equal(t, uint32(common.CodeInvalidBlockInput), result.Code)
	})

	suite.Run("invalid milestone uuid", func() {
		suite.contractCaller = mocks.IContractCaller{}

		milestone.MilestoneID = "0-0a18-41a8-ab7e-59d8002f027b - 0x901a64406d97a3fa9b87b320cbeb86b3c62328f5"

		// create milestone msg
		msgMilestone := types.NewMsgMilestoneBlock(
			milestone.Proposer,
			milestone.StartBlock,
			milestone.EndBlock-1,
			milestone.Hash,
			borChainId,
			milestone.MilestoneID,
		)

		msgMilestone.BorChainID = msgMilestone.MilestoneID

		suite.contractCaller.On("CheckIfBlocksExist", milestone.EndBlock+cmTypes.DefaultMaticchainMilestoneTxConfirmations).Return(true)
		suite.contractCaller.On("GetVoteOnHash", milestone.StartBlock, milestone.EndBlock, milestoneLength, milestone.Hash.String(), milestone.MilestoneID).Return(true, nil)

		result := suite.sideHandler(ctx, msgMilestone)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should fail")
		require.Equal(t, uint32(common.CodeInvalidBlockInput), result.Code)
	})

	suite.Run("invalid milestone proposer", func() {
		suite.contractCaller = mocks.IContractCaller{}

		milestone.MilestoneID = "17ce48fe-0a18-41a8-ab7e-59d8002f027b - 0xz01a64406d97a3fa9b87b320cbeb86b3c62328f5"

		// create milestone msg
		msgMilestone := types.NewMsgMilestoneBlock(
			milestone.Proposer,
			milestone.StartBlock,
			milestone.EndBlock-1,
			milestone.Hash,
			borChainId,
			milestone.MilestoneID,
		)

		msgMilestone.BorChainID = msgMilestone.MilestoneID

		suite.contractCaller.On("CheckIfBlocksExist", milestone.EndBlock+cmTypes.DefaultMaticchainMilestoneTxConfirmations).Return(true)
		suite.contractCaller.On("GetVoteOnHash", milestone.StartBlock, milestone.EndBlock, milestoneLength, milestone.Hash.String(), milestone.MilestoneID).Return(true, nil)

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
			milestone.Hash,
			borChainId,
			milestone.MilestoneID,
		)

		suite.contractCaller.On("CheckIfBlocksExist", milestone.EndBlock+cmTypes.DefaultMaticchainMilestoneTxConfirmations).Return(true)
		suite.contractCaller.On("GetVoteOnHash", milestone.StartBlock, milestone.EndBlock, milestoneLength, milestone.Hash.String(), milestone.MilestoneID).Return(true, nil)

		result := suite.sideHandler(ctx, msgMilestone)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should fail as milestone is not in continuity to latest stored milestone ")

	})
}

func (suite *SideHandlerTestSuite) TestPostHandleMsgMilestone() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper
	stakingKeeper := app.StakingKeeper
	start := uint64(0)
	milestoneLength := helper.MilestoneLength

	// check valid milestone
	// generate proposer for validator set
	chSim.LoadValidatorSet(t, 2, stakingKeeper, ctx, false, 10, 0)
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
			milestone.Hash,
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

	suite.Run("Failure-Invalid Start Block", func() {
		// create milestone msg
		msgMilestone := types.NewMsgMilestoneBlock(
			milestone.Proposer,
			milestone.StartBlock+1,
			milestone.EndBlock+1,
			milestone.Hash,
			borChainId,
			"00000",
		)

		_ = suite.postHandler(ctx, msgMilestone, abci.SideTxResultType_Yes)

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
			milestone.Hash,
			borChainId,
			"00001",
		)
		_ = suite.postHandler(ctx, msgMilestone, abci.SideTxResultType_Yes)

		bufferedHeader, err := keeper.GetLastMilestone(ctx)
		require.Equal(t, bufferedHeader.StartBlock, milestone.StartBlock)
		require.Equal(t, bufferedHeader.EndBlock, milestone.EndBlock)
		require.Equal(t, bufferedHeader.Hash, milestone.Hash)
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
			milestone.Hash,
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
			milestone.Hash,
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
			milestone.Hash,
			borChainId,
			"00004",
		)
		_ = suite.postHandler(ctx, msgMilestone, abci.SideTxResultType_No)
		lastNoAckMilestone := keeper.GetLastNoAckMilestone(ctx)
		require.Equal(t, lastNoAckMilestone, "00004")
	})
}
