package checkpoint_test

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	cmTypes "github.com/maticnetwork/heimdall/chainmanager/types"
	chSim "github.com/maticnetwork/heimdall/checkpoint/simulation"
	"github.com/maticnetwork/heimdall/checkpoint/types"
	errs "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	hmTypes "github.com/maticnetwork/heimdall/types"
)

func TestMilestoneHandlerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestHandleMsgMilestone() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper
	stakingKeeper := app.StakingKeeper
	start := uint64(0)
	borChainId := "1234"
	milestoneID := "0000"
	milestoneLength := helper.MilestoneLength

	// check valid milestone
	// generate proposer for validator set
	chSim.LoadValidatorSet(t, 2, stakingKeeper, ctx, false, 10)
	stakingKeeper.MilestoneIncrementAccum(ctx, 1)

	lastMilestone, err := keeper.GetLastMilestone(ctx)
	if err == nil {
		start = start + lastMilestone.EndBlock + 1
	}

	header, err := chSim.GenRandMilestone(start, milestoneLength)
	require.NoError(t, err)

	ctx = ctx.WithBlockHeight(3)

	//Test1- When milestone is proposed by wrong proposer
	suite.Run("Invalid Proposer", func() {
		header.Proposer = hmTypes.HexToHeimdallAddress("1234")
		msgMilestone := types.NewMsgMilestoneBlock(
			header.Proposer,
			header.StartBlock,
			header.EndBlock,
			header.Hash,
			borChainId,
			milestoneID,
		)

		// send milestone to handler
		got := suite.handler(ctx, msgMilestone)
		require.True(t, !got.IsOK(), errs.CodeToDefaultMsg(got.Code))
		require.Equal(t, errs.CodeInvalidMsg, got.Code)
	})

	// add current proposer to header
	header.Proposer = stakingKeeper.GetMilestoneValidatorSet(ctx).Proposer.Signer

	//Test2- When milestone is proposed of length shorter than configured minimum length
	suite.Run("Invalid msg based on milestone length", func() {
		msgMilestone := types.NewMsgMilestoneBlock(
			header.Proposer,
			header.StartBlock,
			header.EndBlock-1,
			header.Hash,
			borChainId,
			milestoneID,
		)

		// send milestone to handler
		got := suite.handler(ctx, msgMilestone)
		require.True(t, !got.IsOK(), errs.CodeToDefaultMsg(got.Code))
		require.Equal(t, errs.CodeMilestoneInvalid, got.Code)
	})

	// add current proposer to header
	header.Proposer = stakingKeeper.GetMilestoneValidatorSet(ctx).Proposer.Signer

	//Test3- When the first milestone is composed of incorrect start number
	suite.Run("Failure-Invalid Start Block Number", func() {
		msgMilestone := types.NewMsgMilestoneBlock(
			header.Proposer,
			uint64(1),
			header.EndBlock+1,
			header.Hash,
			borChainId,
			milestoneID,
		)

		// send milestone to handler
		got := suite.handler(ctx, msgMilestone)
		require.False(t, got.IsOK(), "expected send-milstone to be ok, got %v", got)
		require.Equal(t, errs.CodeNoMilestone, got.Code)
	})

	header.Proposer = stakingKeeper.GetMilestoneValidatorSet(ctx).Proposer.Signer

	//Test4- When the correct milestone is proposed
	suite.Run("Success", func() {
		msgMilestone := types.NewMsgMilestoneBlock(
			header.Proposer,
			header.StartBlock,
			header.EndBlock,
			header.Hash,
			borChainId,
			milestoneID,
		)

		// send milestone to handler
		got := suite.handler(ctx, msgMilestone)
		require.True(t, got.IsOK(), "expected send-milstone to be ok, got %v", got)
		bufferedHeader, _ := keeper.GetLastMilestone(ctx)
		require.Empty(t, bufferedHeader, "Should not store state")
		milestoneBlockNumber := keeper.GetMilestoneBlockNumber(ctx)
		require.Equal(t, int64(3), milestoneBlockNumber, "Mismatch in milestoneBlockNumber")
	})

	header.Proposer = stakingKeeper.GetMilestoneValidatorSet(ctx).Proposer.Signer

	ctx = ctx.WithBlockHeight(int64(4))

	//Test5- When previous milestone is still in voting phase
	suite.Run("Previous milestone is still in voting phase", func() {

		msgMilestone := types.NewMsgMilestoneBlock(
			header.Proposer,
			start,
			start+milestoneLength-1,
			header.Hash,
			borChainId,
			milestoneID,
		)

		// send milestone to handler
		got := suite.handler(ctx, msgMilestone)
		require.True(t, !got.IsOK(), errs.CodeToDefaultMsg(got.Code))
		require.Equal(t, errs.CodePrevMilestoneInVoting, got.Code)
	})

	header.Proposer = stakingKeeper.GetMilestoneValidatorSet(ctx).Proposer.Signer

	ctx = ctx.WithBlockHeight(int64(6))

	//Test5- When milestone is not in continuity
	suite.Run("Milestone not in countinuity", func() {

		err := keeper.AddMilestone(ctx, header)
		require.NoError(t, err)

		_, err = keeper.GetLastMilestone(ctx)
		require.NoError(t, err)

		lastMilestone, err := keeper.GetLastMilestone(ctx)
		if err == nil {
			// pass wrong start
			start = start + lastMilestone.EndBlock + 2 //Start block is 2 more than last milestone's end block
		}

		msgMilestone := types.NewMsgMilestoneBlock(
			header.Proposer,
			start,
			start+milestoneLength-1,
			header.Hash,
			borChainId,
			milestoneID,
		)

		// send milestone to handler
		got := suite.handler(ctx, msgMilestone)
		require.True(t, !got.IsOK(), errs.CodeToDefaultMsg(got.Code))
		require.Equal(t, errs.CodeMilestoneNotInContinuity, got.Code)
	})

	header.Proposer = stakingKeeper.GetMilestoneValidatorSet(ctx).Proposer.Signer

	//Test6- When milestone is not in continuity
	suite.Run("Milestone not in countinuity", func() {

		_, err = keeper.GetLastMilestone(ctx)
		require.NoError(t, err)

		lastMilestone, err := keeper.GetLastMilestone(ctx)
		if err == nil {
			// pass wrong start
			start = start + lastMilestone.EndBlock - 2 //Start block is 2 less than last milestone's end block
		}

		msgMilestone := types.NewMsgMilestoneBlock(
			header.Proposer,
			start,
			start+milestoneLength-1,
			header.Hash,
			borChainId,
			milestoneID,
		)

		// send milestone to handler
		got := suite.handler(ctx, msgMilestone)
		require.True(t, !got.IsOK(), errs.CodeToDefaultMsg(got.Code))
		require.Equal(t, errs.CodeMilestoneNotInContinuity, got.Code)
	})
}

// Test to check that passed milestone should be in the store
func (suite *HandlerTestSuite) TestHandleMsgMilestoneExistInStore() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper
	stakingKeeper := app.StakingKeeper
	start := uint64(0)
	milestoneLength := helper.MilestoneLength

	chSim.LoadValidatorSet(t, 2, stakingKeeper, ctx, false, 10)
	stakingKeeper.IncrementAccum(ctx, 1)

	lastMilestone, err := keeper.GetLastMilestone(ctx)
	if err == nil {
		start = start + lastMilestone.EndBlock + 1
	}

	header, err := chSim.GenRandMilestone(start, milestoneLength)
	require.NoError(t, err)

	// add current proposer to header
	header.Proposer = stakingKeeper.GetValidatorSet(ctx).Proposer.Signer

	// send old milestone
	res := suite.SendMilestone(header)

	require.True(t, res.IsOK(), "expected send-milestone to be  ok, got %v", res)

	// send milestone to handler
	got := suite.SendMilestone(header)
	require.False(t, got.IsOK(), errs.CodeToDefaultMsg(got.Code))
}

func (suite *HandlerTestSuite) SendMilestone(header hmTypes.Milestone) (res sdk.Result) {
	_, ctx := suite.app, suite.ctx

	milestoneLength := helper.MilestoneLength

	// keeper := app.MilestoneKeeper

	borChainId := "1234"
	milestoneID := "00000"
	// create milestone msg
	msgMilestone := types.NewMsgMilestoneBlock(
		header.Proposer,
		header.StartBlock,
		header.EndBlock,
		header.Hash,
		borChainId,
		milestoneID,
	)

	suite.contractCaller.On("CheckIfBlocksExist", header.EndBlock).Return(true)
	suite.contractCaller.On("CheckIfBlocksExist", header.EndBlock+cmTypes.DefaultMaticchainMilestoneTxConfirmations).Return(true)
	suite.contractCaller.On("GetRootHash", header.StartBlock, header.EndBlock, milestoneLength).Return(header.Hash.Bytes(), nil)
	suite.contractCaller.On("GetVoteOnHash", header.StartBlock, header.EndBlock, milestoneLength, header.Hash.String(), header.MilestoneID).Return(true, nil)

	ctx = ctx.WithBlockHeight(int64(3))

	// send milestone to handler
	result := suite.handler(ctx, msgMilestone)
	sideResult := suite.sideHandler(ctx, msgMilestone)
	suite.postHandler(ctx, msgMilestone, sideResult.Result)

	return result
}

func (suite *HandlerTestSuite) TestHandleMsgMilestoneTimeout() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper
	stakingKeeper := app.StakingKeeper

	startBlock := uint64(0)
	endBlock := uint64(63)
	hash := hmTypes.HexToHeimdallHash("123")
	proposerAddress := hmTypes.HexToHeimdallAddress("123")
	timestamp := uint64(0)
	borChainId := "1234"
	milestoneID := "0000"

	proposer := hmTypes.HeimdallAddress{}

	chSim.LoadValidatorSet(t, 2, stakingKeeper, ctx, false, 10)

	suite.Run("Last milestone not found", func() {
		msgMilestoneTimeout := types.NewMsgMilestoneTimeout(
			proposer,
		)

		// send milestone to handler
		got := suite.handler(ctx, msgMilestoneTimeout)
		require.True(t, !got.IsOK(), errs.CodeToDefaultMsg(got.Code))
		require.Equal(t, errs.CodeNoMilestone, got.Code)
	})

	milestone := hmTypes.CreateMilestone(
		startBlock,
		endBlock,
		hash,
		proposerAddress,
		borChainId,
		milestoneID,
		timestamp,
	)
	_ = keeper.AddMilestone(ctx, milestone)

	newTime := milestone.TimeStamp + uint64(helper.MilestoneBufferTime) - 1
	suite.ctx = ctx.WithBlockTime(time.Unix(0, int64(newTime)))

	msgMilestoneTimeout := types.NewMsgMilestoneTimeout(
		proposer,
	)

	// send milestone to handler
	got := suite.handler(ctx, msgMilestoneTimeout)
	require.True(t, !got.IsOK(), errs.CodeToDefaultMsg(got.Code))
	require.Equal(t, errs.CodeInvalidMilestoneTimeout, got.Code)

	newTime = milestone.TimeStamp + 2*uint64(helper.MilestoneBufferTime) + 10000000
	suite.ctx = ctx.WithBlockTime(time.Unix(0, int64(newTime)))

	msgMilestoneTimeout = types.NewMsgMilestoneTimeout(
		proposer,
	)

	// send milestone to handler
	got = suite.handler(suite.ctx, msgMilestoneTimeout)
	require.True(t, got.IsOK(), errs.CodeToDefaultMsg(got.Code))
}
