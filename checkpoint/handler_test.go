package checkpoint_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	cmTypes "github.com/maticnetwork/heimdall/chainmanager/types"
	"github.com/maticnetwork/heimdall/checkpoint/types"
	errs "github.com/maticnetwork/heimdall/common"

	"github.com/maticnetwork/heimdall/checkpoint"
	chSim "github.com/maticnetwork/heimdall/checkpoint/simulation"

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
	topupKeeper.AddDividendAccount(ctx, dividendAccount)

	// check valid checkpoint
	// generate proposer for validator set
	chSim.LoadValidatorSet(2, t, stakingKeeper, ctx, false, 10)
	stakingKeeper.IncrementAccum(ctx, 1)

	lastCheckpoint, err := keeper.GetLastCheckpoint(ctx)
	if err == nil {
		start = start + lastCheckpoint.EndBlock + 1
	}

	header, err := chSim.GenRandCheckpoint(start, maxSize, params.MaxCheckpointLength)

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

		keeper.AddCheckpoint(ctx, headerId, header)
		keeper.GetCheckpointByNumber(ctx, headerId)
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

func (suite *HandlerTestSuite) TestHandleMsgCheckpointAdjustCheckpointBuffer() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper

	checkpoint := hmTypes.Checkpoint{
		Proposer:   hmTypes.HexToHeimdallAddress("123"),
		StartBlock: 0,
		EndBlock:   256,
		RootHash:   hmTypes.HexToHeimdallHash("123"),
		BorChainID: "testchainid",
		TimeStamp:  1,
	}
	keeper.SetCheckpointBuffer(ctx, checkpoint)

	checkpointAdjust := types.MsgCheckpointAdjust{
		HeaderIndex: 1,
		Proposer:    hmTypes.HexToHeimdallAddress("456"),
		StartBlock:  0,
		EndBlock:    512,
		RootHash:    hmTypes.HexToHeimdallHash("456"),
	}

	result := suite.handler(ctx, checkpointAdjust)
	require.False(t, result.IsOK())
}

func (suite *HandlerTestSuite) TestHandleMsgCheckpointAdjustSameCheckpointAsMsg() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper

	checkpoint := hmTypes.Checkpoint{
		Proposer:   hmTypes.HexToHeimdallAddress("123"),
		StartBlock: 0,
		EndBlock:   256,
		RootHash:   hmTypes.HexToHeimdallHash("123"),
		BorChainID: "testchainid",
		TimeStamp:  1,
	}
	keeper.AddCheckpoint(ctx, 1, checkpoint)

	checkpointAdjust := types.MsgCheckpointAdjust{
		HeaderIndex: 1,
		Proposer:    hmTypes.HexToHeimdallAddress("123"),
		StartBlock:  0,
		EndBlock:    256,
		RootHash:    hmTypes.HexToHeimdallHash("123"),
	}

	result := suite.handler(ctx, checkpointAdjust)
	require.False(t, result.IsOK())
}

func (suite *HandlerTestSuite) TestHandleMsgCheckpointAfterBufferTimeOut() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper
	stakingKeeper := app.StakingKeeper
	topupKeeper := app.TopupKeeper
	start := uint64(0)
	maxSize := uint64(256)
	params := keeper.GetParams(ctx)
	checkpointBufferTime := params.CheckpointBufferTime
	dividendAccount := hmTypes.DividendAccount{
		User:      hmTypes.HexToHeimdallAddress("123"),
		FeeAmount: big.NewInt(0).String(),
	}
	topupKeeper.AddDividendAccount(ctx, dividendAccount)

	// generate proposer for validator set
	chSim.LoadValidatorSet(2, t, stakingKeeper, ctx, false, 10)
	stakingKeeper.IncrementAccum(ctx, 1)
	lastCheckpoint, err := keeper.GetLastCheckpoint(ctx)
	if err == nil {
		start = start + lastCheckpoint.EndBlock + 1
	}

	header, err := chSim.GenRandCheckpoint(start, maxSize, params.MaxCheckpointLength)

	// add current proposer to header
	header.Proposer = stakingKeeper.GetValidatorSet(ctx).Proposer.Signer

	// send old checkpoint
	res := suite.SendCheckpoint(header)
	require.True(t, res.IsOK(), "expected send-checkpoint to be  ok, got %v", res)

	checkpointBuffer, err := keeper.GetCheckpointFromBuffer(ctx)

	// set time buffered checkpoint timestamp + checkpointBufferTime
	newTime := checkpointBuffer.TimeStamp + uint64(checkpointBufferTime)
	suite.ctx = ctx.WithBlockTime(time.Unix(0, int64(newTime)))

	// send new checkpoint which should replace old one
	got := suite.SendCheckpoint(header)
	require.True(t, got.IsOK(), "expected send-checkpoint to be  ok, got %v", got)
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
	topupKeeper.AddDividendAccount(ctx, dividendAccount)

	chSim.LoadValidatorSet(2, t, stakingKeeper, ctx, false, 10)
	stakingKeeper.IncrementAccum(ctx, 1)
	lastCheckpoint, err := keeper.GetLastCheckpoint(ctx)
	if err == nil {
		start = start + lastCheckpoint.EndBlock + 1
	}

	header, err := chSim.GenRandCheckpoint(start, maxSize, params.MaxCheckpointLength)

	// add current proposer to header
	header.Proposer = stakingKeeper.GetValidatorSet(ctx).Proposer.Signer

	// send old checkpoint
	res := suite.SendCheckpoint(header)

	require.True(t, res.IsOK(), "expected send-checkpoint to be  ok, got %v", res)

	// send checkpoint to handler
	got := suite.SendCheckpoint(header)
	require.True(t, !got.IsOK(), errs.CodeToDefaultMsg(got.Code))
}

func (suite *HandlerTestSuite) TestHandleMsgCheckpointAck() {
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
	topupKeeper.AddDividendAccount(ctx, dividendAccount)

	// check valid checkpoint
	// generate proposer for validator set
	chSim.LoadValidatorSet(2, t, stakingKeeper, ctx, false, 10)
	stakingKeeper.IncrementAccum(ctx, 1)

	lastCheckpoint, err := keeper.GetLastCheckpoint(ctx)
	if err == nil {
		start = start + lastCheckpoint.EndBlock + 1
	}

	header, err := chSim.GenRandCheckpoint(start, maxSize, params.MaxCheckpointLength)

	// add current proposer to header
	header.Proposer = stakingKeeper.GetValidatorSet(ctx).Proposer.Signer

	got := suite.SendCheckpoint(header)
	require.True(t, got.IsOK(), "expected send-checkpoint to be ok, got %v", got)

	bufferedHeader, err := keeper.GetCheckpointFromBuffer(ctx)
	require.NotNil(t, bufferedHeader)

	// send ack
	headerId := uint64(1)
	suite.Run("success", func() {
		msgCheckpointAck := types.NewMsgCheckpointAck(
			hmTypes.HexToHeimdallAddress("123"),
			headerId,
			header.Proposer,
			header.StartBlock,
			header.EndBlock,
			header.RootHash,
			hmTypes.HexToHeimdallHash("123123"),
			uint64(1),
		)
		result := suite.handler(ctx, msgCheckpointAck)
		require.True(t, result.IsOK(), "expected send-ack to be ok, got %v", result)
		afterAckBufferedCheckpoint, _ := keeper.GetCheckpointFromBuffer(ctx)
		require.NotNil(t, afterAckBufferedCheckpoint, "should not remove from buffer")
	})

	suite.Run("Invalid start", func() {
		msgCheckpointAck := types.NewMsgCheckpointAck(
			hmTypes.HexToHeimdallAddress("123"),
			headerId,
			header.Proposer,
			uint64(123),
			header.EndBlock,
			header.RootHash,
			hmTypes.HexToHeimdallHash("123123"),
			uint64(1),
		)

		got := suite.handler(ctx, msgCheckpointAck)
		require.True(t, !got.IsOK(), errs.CodeToDefaultMsg(got.Code))
	})

	suite.Run("Invalid Roothash", func() {
		msgCheckpointAck := types.NewMsgCheckpointAck(
			hmTypes.HexToHeimdallAddress("123"),
			headerId,
			header.Proposer,
			header.StartBlock,
			header.EndBlock,
			hmTypes.HexToHeimdallHash("9887"),
			hmTypes.HexToHeimdallHash("123123"),
			uint64(1),
		)

		got := suite.handler(ctx, msgCheckpointAck)
		require.True(t, !got.IsOK(), errs.CodeToDefaultMsg(got.Code))
	})
}

func (suite *HandlerTestSuite) TestHandleMsgCheckpointNoAck() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper
	stakingKeeper := app.StakingKeeper
	topupKeeper := app.TopupKeeper
	start := uint64(0)
	maxSize := uint64(256)
	params := keeper.GetParams(ctx)
	checkpointBufferTime := params.CheckpointBufferTime

	dividendAccount := hmTypes.DividendAccount{
		User:      hmTypes.HexToHeimdallAddress("123"),
		FeeAmount: big.NewInt(0).String(),
	}
	topupKeeper.AddDividendAccount(ctx, dividendAccount)

	// check valid checkpoint
	// generate proposer for validator set
	chSim.LoadValidatorSet(2, t, stakingKeeper, ctx, false, 10)
	stakingKeeper.IncrementAccum(ctx, 1)

	lastCheckpoint, err := keeper.GetLastCheckpoint(ctx)
	if err == nil {
		start = start + lastCheckpoint.EndBlock + 1
	}

	header, err := chSim.GenRandCheckpoint(start, maxSize, params.MaxCheckpointLength)

	// add current proposer to header
	header.Proposer = stakingKeeper.GetValidatorSet(ctx).Proposer.Signer

	got := suite.SendCheckpoint(header)
	require.True(t, got.IsOK(), "expected send-NoAck to be ok, got %v", got)

	// set time lastCheckpoint timestamp + checkpointBufferTime
	newTime := lastCheckpoint.TimeStamp + uint64(checkpointBufferTime)
	suite.ctx = ctx.WithBlockTime(time.Unix(0, int64(newTime)))
	result := suite.SendNoAck()
	require.True(t, result.IsOK(), "expected send-NoAck to be ok, got %v", got)
	ackCount := keeper.GetACKCount(ctx)
	require.Equal(t, uint64(0), uint64(ackCount), "Should not update state")
}

func (suite *HandlerTestSuite) TestHandleMsgCheckpointNoAckBeforeBufferTimeout() {
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
	topupKeeper.AddDividendAccount(ctx, dividendAccount)

	// check valid checkpoint
	// generate proposer for validator set
	chSim.LoadValidatorSet(2, t, stakingKeeper, ctx, false, 10)
	stakingKeeper.IncrementAccum(ctx, 1)

	lastCheckpoint, err := keeper.GetLastCheckpoint(ctx)
	if err == nil {
		start = start + lastCheckpoint.EndBlock + 1
	}

	header, err := chSim.GenRandCheckpoint(start, maxSize, params.MaxCheckpointLength)

	// add current proposer to header
	header.Proposer = stakingKeeper.GetValidatorSet(ctx).Proposer.Signer

	got := suite.SendCheckpoint(header)
	require.True(t, got.IsOK(), "expected send-checkpoint to be ok, got %v", got)

	result := suite.SendNoAck()
	require.True(t, !result.IsOK(), errs.CodeToDefaultMsg(result.Code))
}

func (suite *HandlerTestSuite) SendCheckpoint(header hmTypes.Checkpoint) (res sdk.Result) {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	// keeper := app.CheckpointKeeper
	topupKeeper := app.TopupKeeper

	dividendAccounts := topupKeeper.GetAllDividendAccounts(ctx)
	accRootHash, err := types.GetAccountRootHash(dividendAccounts)
	require.NoError(t, err)
	accountRoot := hmTypes.BytesToHeimdallHash(accRootHash)

	borChainId := "1234"
	// create checkpoint msg
	msgCheckpoint := types.NewMsgCheckpointBlock(
		header.Proposer,
		header.StartBlock,
		header.EndBlock,
		header.RootHash,
		accountRoot,
		borChainId,
	)

	suite.contractCaller.On("CheckIfBlocksExist", header.EndBlock+cmTypes.DefaultMaticchainTxConfirmations).Return(true)
	suite.contractCaller.On("GetRootHash", header.StartBlock, header.EndBlock, uint64(1024)).Return(header.RootHash.Bytes(), nil)

	// send checkpoint to handler
	result := suite.handler(ctx, msgCheckpoint)
	sideResult := suite.sideHandler(ctx, msgCheckpoint)
	suite.postHandler(ctx, msgCheckpoint, sideResult.Result)
	return result
}

func (suite *HandlerTestSuite) SendNoAck() (res sdk.Result) {
	_, _, ctx := suite.T(), suite.app, suite.ctx
	msgNoAck := types.NewMsgCheckpointNoAck(hmTypes.HexToHeimdallAddress("123"))

	result := suite.handler(ctx, msgNoAck)
	sideResult := suite.sideHandler(ctx, msgNoAck)
	suite.postHandler(ctx, msgNoAck, sideResult.Result)
	return result
}
