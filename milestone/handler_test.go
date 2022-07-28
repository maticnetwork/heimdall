package milestone_test

import (
	"math/big"
	"testing"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/checkpoint"
	chSim "github.com/maticnetwork/heimdall/checkpoint/simulation"
	"github.com/maticnetwork/heimdall/checkpoint/types"
	errs "github.com/maticnetwork/heimdall/common"

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
	suite.handler = checkpoint.NewHandler(suite.app.CheckpointKeeper, &suite.contractCaller)
	suite.sideHandler = checkpoint.NewSideTxHandler(suite.app.CheckpointKeeper, &suite.contractCaller)
	suite.postHandler = checkpoint.NewPostTxHandler(suite.app.CheckpointKeeper, &suite.contractCaller)
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
func (suite *HandlerTestSuite) TestHandleMsgCheckpoint() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper
	stakingKeeper := app.StakingKeeper
	topupKeeper := app.TopupKeeper
	start := uint64(0)
	maxSize := uint64(256)
	borChainId := "1234"
	params := keeper.GetParams(ctx)
	dividendAccount := hmTypes.DividendAccount{
		User:      hmTypes.HexToHeimdallAddress("123"),
		FeeAmount: big.NewInt(0).String(),
	}
	err := topupKeeper.AddDividendAccount(ctx, dividendAccount)
	require.NoError(t, err)

	// check valid checkpoint
	// generate proposer for validator set
	chSim.LoadValidatorSet(t, 2, stakingKeeper, ctx, false, 10)
	stakingKeeper.IncrementAccum(ctx, 1)

	lastCheckpoint, err := keeper.GetLastCheckpoint(ctx)
	if err == nil {
		start = start + lastCheckpoint.EndBlock + 1
	}

	header, err := chSim.GenRandCheckpoint(start, maxSize, params.MaxCheckpointLength)
	require.NoError(t, err)

	// add current proposer to header
	header.Proposer = stakingKeeper.GetValidatorSet(ctx).Proposer.Signer

	dividendAccounts := topupKeeper.GetAllDividendAccounts(ctx)
	accRootHash, err := types.GetAccountRootHash(dividendAccounts)
	require.NoError(t, err)

	accountRoot := hmTypes.BytesToHeimdallHash(accRootHash)

	suite.Run("Success", func() {
		msgCheckpoint := types.NewMsgCheckpointBlock(
			header.Proposer,
			header.StartBlock,
			header.EndBlock,
			header.RootHash,
			accountRoot,
			borChainId,
		)

		// send checkpoint to handler
		got := suite.handler(ctx, msgCheckpoint)
		require.True(t, got.IsOK(), "expected send-checkpoint to be ok, got %v", got)
		bufferedHeader, _ := keeper.GetCheckpointFromBuffer(ctx)
		require.Empty(t, bufferedHeader, "Should not store state")
	})

	suite.Run("Invalid Proposer", func() {
		header.Proposer = hmTypes.HexToHeimdallAddress("1234")
		msgCheckpoint := types.NewMsgCheckpointBlock(
			header.Proposer,
			header.StartBlock,
			header.EndBlock,
			header.RootHash,
			accountRoot,
			borChainId,
		)

		// send checkpoint to handler
		got := suite.handler(ctx, msgCheckpoint)
		require.True(t, !got.IsOK(), errs.CodeToDefaultMsg(got.Code))
	})

	suite.Run("Checkpoint not in countinuity", func() {
		headerId := uint64(10000)

		err = keeper.AddCheckpoint(ctx, headerId, header)
		require.NoError(t, err)

		_, err = keeper.GetCheckpointByNumber(ctx, headerId)
		require.NoError(t, err)

		keeper.UpdateACKCount(ctx)
		lastCheckpoint, err := keeper.GetLastCheckpoint(ctx)
		if err == nil {
			// pass wrong start
			start = start + lastCheckpoint.EndBlock + 2
		}

		msgCheckpoint := types.NewMsgCheckpointBlock(
			header.Proposer,
			start,
			start+256,
			header.RootHash,
			accountRoot,
			borChainId,
		)

		// send checkpoint to handler
		got := suite.handler(ctx, msgCheckpoint)
		require.True(t, !got.IsOK(), errs.CodeToDefaultMsg(got.Code))
	})
}

func (suite *HandlerTestSuite) TestHandleMsgCheckpointExistInBuffer() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper
	stakingKeeper := app.StakingKeeper
	topupKeeper := app.TopupKeeper
	start := uint64(0)
	maxSize := uint64(256)
	params := keeper.GetParams(ctx)
	dividendAccount := hmTypes.DividendAccount{
		User:      hmTypes.HexToHeimdallAddress("123"),
		FeeAmount: big.NewInt(0).String(),
	}

	err := topupKeeper.AddDividendAccount(ctx, dividendAccount)
	require.NoError(t, err)

	chSim.LoadValidatorSet(t, 2, stakingKeeper, ctx, false, 10)
	stakingKeeper.IncrementAccum(ctx, 1)

	lastCheckpoint, err := keeper.GetLastCheckpoint(ctx)
	if err == nil {
		start = start + lastCheckpoint.EndBlock + 1
	}

	header, err := chSim.GenRandCheckpoint(start, maxSize, params.MaxCheckpointLength)
	require.NoError(t, err)

	// add current proposer to header
	header.Proposer = stakingKeeper.GetValidatorSet(ctx).Proposer.Signer

	// send old checkpoint
	res := suite.SendCheckpoint(header)

	require.True(t, res.IsOK(), "expected send-checkpoint to be  ok, got %v", res)

	// send checkpoint to handler
	got := suite.SendCheckpoint(header)
	require.True(t, !got.IsOK(), errs.CodeToDefaultMsg(got.Code))
}

func (suite *HandlerTestSuite) SendCheckpoint(header hmTypes.Checkpoint) (res sdk.Result) {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	// keeper := app.CheckpointKeeper

	borChainId := "1234"
	// create checkpoint msg
	msgMilestone := types.NewMsgMilestoneBlock(
		header.Proposer,
		header.StartBlock,
		header.EndBlock,
		header.RootHash,
		borChainId,
	)

	suite.contractCaller.On("CheckIfBlocksExist", header.EndBlock).Return(true)
	suite.contractCaller.On("GetRootHash", header.StartBlock, header.EndBlock).Return(header.RootHash.Bytes(), nil)

	// send checkpoint to handler
	result := suite.handler(ctx, msgCheckpoint)
	sideResult := suite.sideHandler(ctx, msgCheckpoint)
	suite.postHandler(ctx, msgCheckpoint, sideResult.Result)

	return result
}
