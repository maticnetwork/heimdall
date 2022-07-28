package milestone_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	errs "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/milestone"
	chSim "github.com/maticnetwork/heimdall/milestone/simulation"
	"github.com/maticnetwork/heimdall/milestone/types"

	"github.com/maticnetwork/heimdall/helper/mocks"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite

	app    *app.HeimdallApp
	ctx    sdk.Context
	cliCtx context.CLIContext

	handler        sdk.Handler
	sideHandler    hmTypes.SideTxHandler
	postHandler    hmTypes.PostTxHandler
	contractCaller mocks.IContractCaller
}

func (suite *HandlerTestSuite) SetupTest() {
	suite.app, suite.ctx, suite.cliCtx = createTestApp(false)
	suite.contractCaller = mocks.IContractCaller{}
	suite.handler = milestone.NewHandler(suite.app.MilestoneKeeper, &suite.contractCaller)
	suite.sideHandler = milestone.NewSideTxHandler(suite.app.MilestoneKeeper, &suite.contractCaller)
	suite.postHandler = milestone.NewPostTxHandler(suite.app.MilestoneKeeper, &suite.contractCaller)
}

func TestHandlerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestHandler() {
	t, ctx := suite.T(), suite.ctx

	// side handler
	result := suite.handler(ctx, nil)
	require.False(t, result.IsOK(), "Handler should fail")
}

// test handler for message
func (suite *HandlerTestSuite) TestHandleMsgMilestone() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.MilestoneKeeper
	stakingKeeper := app.StakingKeeper
	start := uint64(0)
	borChainId := "1234"
	params := keeper.GetParams(ctx)
	sprintLength := params.SprintLength

	// check valid milestone
	// generate proposer for validator set
	chSim.LoadValidatorSet(t, 2, stakingKeeper, ctx, false, 10)
	stakingKeeper.IncrementAccum(ctx, 1)

	lastMilestone, err := keeper.GetMilestone(ctx)
	if err == nil {
		start = start + lastMilestone.EndBlock + 1
	}

	header, err := chSim.GenRandMilestone(start, sprintLength)
	require.NoError(t, err)

	// add current proposer to header
	header.Proposer = stakingKeeper.GetValidatorSet(ctx).Proposer.Signer

	suite.Run("Success", func() {
		msgMilestone := types.NewMsgMilestoneBlock(
			header.Proposer,
			header.StartBlock,
			header.EndBlock,
			header.RootHash,
			borChainId,
		)

		// send milestone to handler
		got := suite.handler(ctx, msgMilestone)
		require.True(t, got.IsOK(), "expected send-milstone to be ok, got %v", got)
		bufferedHeader, _ := keeper.GetMilestone(ctx)
		require.Empty(t, bufferedHeader, "Should not store state")
	})

	suite.Run("Invalid Proposer", func() {
		header.Proposer = hmTypes.HexToHeimdallAddress("1234")
		msgMilestone := types.NewMsgMilestoneBlock(
			header.Proposer,
			header.StartBlock,
			header.EndBlock,
			header.RootHash,
			borChainId,
		)

		// send milestone to handler
		got := suite.handler(ctx, msgMilestone)
		require.True(t, !got.IsOK(), errs.CodeToDefaultMsg(got.Code))
	})

	suite.Run("Milestone not in countinuity", func() {

		err = keeper.AddMilestone(ctx, header)
		require.NoError(t, err)

		_, err = keeper.GetMilestone(ctx)
		require.NoError(t, err)

		lastMilestone, err := keeper.GetMilestone(ctx)
		if err == nil {
			// pass wrong start
			start = start + lastMilestone.EndBlock + 2
		}

		msgMilestone := types.NewMsgMilestoneBlock(
			header.Proposer,
			start,
			start+sprintLength,
			header.RootHash,
			borChainId,
		)

		// send milestone to handler
		got := suite.handler(ctx, msgMilestone)
		require.True(t, !got.IsOK(), errs.CodeToDefaultMsg(got.Code))
	})
}

func (suite *HandlerTestSuite) TestHandleMsgMilestoneExistInStore() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.MilestoneKeeper
	stakingKeeper := app.StakingKeeper
	start := uint64(0)
	params := keeper.GetParams(ctx)

	chSim.LoadValidatorSet(t, 2, stakingKeeper, ctx, false, 10)
	stakingKeeper.IncrementAccum(ctx, 1)

	lastMilestone, err := keeper.GetMilestone(ctx)
	if err == nil {
		start = start + lastMilestone.EndBlock + 1
	}

	header, err := chSim.GenRandMilestone(start, params.SprintLength)
	require.NoError(t, err)

	// add current proposer to header
	header.Proposer = stakingKeeper.GetValidatorSet(ctx).Proposer.Signer

	// send old milestone
	res := suite.SendMilestone(header)

	require.True(t, res.IsOK(), "expected send-milestone to be  ok, got %v", res)

	// send milestone to handler
	got := suite.SendMilestone(header)
	require.True(t, !got.IsOK(), errs.CodeToDefaultMsg(got.Code))
}

func (suite *HandlerTestSuite) SendMilestone(header hmTypes.Milestone) (res sdk.Result) {
	app, ctx := suite.app, suite.ctx
	keeper := app.MilestoneKeeper
	params := keeper.GetParams(ctx)
	sprintLength := params.SprintLength
	// keeper := app.MilestoneKeeper

	borChainId := "1234"
	// create milestone msg
	msgMilestone := types.NewMsgMilestoneBlock(
		header.Proposer,
		header.StartBlock,
		header.EndBlock,
		header.RootHash,
		borChainId,
	)

	suite.contractCaller.On("CheckIfBlocksExist", header.EndBlock).Return(true)
	suite.contractCaller.On("GetRootHash", header.StartBlock, header.EndBlock, sprintLength).Return(header.RootHash.Bytes(), nil)

	// send milestone to handler
	result := suite.handler(ctx, msgMilestone)
	sideResult := suite.sideHandler(ctx, msgMilestone)
	suite.postHandler(ctx, msgMilestone, sideResult.Result)

	return result
}
