package checkpoint_test

import (
	"math/big"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/helper/mocks"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/x/checkpoint"
	chSim "github.com/maticnetwork/heimdall/x/checkpoint/simulation"
	"github.com/maticnetwork/heimdall/x/checkpoint/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type HandlerTestSuite struct {
	suite.Suite

	app    *app.HeimdallApp
	ctx    sdk.Context
	cliCtx client.Context

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
	result, err := suite.handler(ctx, nil)
	require.Nil(t, result)
	require.NotNil(t, err)
	require.Error(t, err, "handler should fail")
}

// test handler for message
func (suite *HandlerTestSuite) TestHandleMsgCheckpoint() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := initApp.CheckpointKeeper
	stakingKeeper := initApp.StakingKeeper
	topupKeeper := initApp.TopupKeeper
	start := uint64(0)
	maxSize := uint64(256)
	borChainId := "1234"
	params := keeper.GetParams(ctx)
	dividendAccount := hmTypes.DividendAccount{
		User:      hmCommonTypes.HexToHeimdallAddress("123").String(),
		FeeAmount: big.NewInt(0).String(),
	}
	err := topupKeeper.AddDividendAccount(ctx, dividendAccount)
	require.NoError(t, err)
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
	accountRoot := hmCommonTypes.BytesToHeimdallHash(accRootHash)
	proposer, err := sdk.AccAddressFromHex(header.Proposer)
	require.NoError(t, err)
	suite.Run("Success", func() {
		msgCheckpoint := types.NewMsgCheckpointBlock(
			proposer,
			header.StartBlock,
			header.EndBlock,
			hmCommonTypes.HexToHeimdallHash(header.RootHash),
			accountRoot,
			borChainId,
		)

		// send checkpoint to handler
		result, err := suite.handler(ctx, &msgCheckpoint)
		require.NotNil(t, result)
		require.NoError(t, err)
		bufferedHeader, _ := keeper.GetCheckpointFromBuffer(ctx)
		require.Empty(t, bufferedHeader, "Should not store state")
	})

	suite.Run("Invalid Proposer", func() {
		invalidProposer, err := sdk.AccAddressFromHex(hmCommonTypes.HexToHeimdallAddress("1234").String())
		require.NoError(t, err)
		msgCheckpoint := types.NewMsgCheckpointBlock(
			invalidProposer,
			header.StartBlock,
			header.EndBlock,
			hmCommonTypes.HexToHeimdallHash(header.RootHash),
			accountRoot,
			borChainId,
		)

		// send checkpoint to handler
		_, err = suite.handler(ctx, &msgCheckpoint)
		require.Error(t, err)
	})

	suite.Run("Checkpoint not in continuity", func() {
		headerId := uint64(10000)

		err := keeper.AddCheckpoint(ctx, headerId, header)
		require.NoError(t, err)
		checkpointInfo, err := keeper.GetCheckpointByNumber(ctx, headerId)
		require.NotNil(t, checkpointInfo)
		require.NoError(t, err)
		keeper.UpdateACKCount(ctx)
		lastCheckpoint, err := keeper.GetLastCheckpoint(ctx)
		if err == nil {
			// pass wrong start
			start = start + lastCheckpoint.EndBlock + 2
		}

		msgCheckpoint := types.NewMsgCheckpointBlock(
			proposer,
			start,
			start+256,
			hmCommonTypes.HexToHeimdallHash(header.RootHash),
			accountRoot,
			borChainId,
		)

		// send checkpoint to handler
		_, err = suite.handler(ctx, &msgCheckpoint)
		require.NoError(t, err)
	})
}

func (suite *HandlerTestSuite) TestHandleMsgCheckpointAck() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := initApp.CheckpointKeeper
	stakingKeeper := initApp.StakingKeeper
	topupKeeper := initApp.TopupKeeper
	start := uint64(0)
	maxSize := uint64(256)
	params := keeper.GetParams(ctx)
	dividendAccount := hmTypes.DividendAccount{
		User:      hmCommonTypes.HexToHeimdallAddress("123").String(),
		FeeAmount: big.NewInt(0).String(),
	}
	err := topupKeeper.AddDividendAccount(ctx, dividendAccount)
	require.NoError(t, err)

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

	got, err := suite.SendCheckpoint(header)
	require.NotNil(t, got)
	require.Nil(t, err)

	bufferedHeader, err := keeper.GetCheckpointFromBuffer(ctx)
	require.NotNil(t, bufferedHeader)
	require.Nil(t, err)

	// send ack
	headerId := uint64(1)
	suite.Run("success", func() {
		fromAccAddr, err := sdk.AccAddressFromHex("123")
		require.NoError(t, err)
		proposer, err := sdk.AccAddressFromHex(header.Proposer)
		require.NoError(t, err)

		msgCheckpointAck := types.NewMsgCheckpointAck(
			fromAccAddr,
			headerId,
			proposer,
			header.StartBlock,
			header.EndBlock,
			hmCommonTypes.HexToHeimdallHash(header.RootHash),
			hmCommonTypes.HexToHeimdallHash("123123"),
			uint64(1),
		)
		result, err := suite.handler(ctx, &msgCheckpointAck)
		require.NotNil(t, result)
		require.NoError(t, err)

		afterAckBufferedCheckpoint, _ := keeper.GetCheckpointFromBuffer(ctx)
		require.NotNil(t, afterAckBufferedCheckpoint, "should not remove from buffer")
	})

	suite.Run("Invalid start", func() {
		msgCheckpointAck := types.NewMsgCheckpointAck(
			hmCommonTypes.HexToHeimdallAddress("123").Bytes(),
			headerId,
			sdk.AccAddress(header.Proposer),
			uint64(123),
			header.EndBlock,
			hmCommonTypes.HexToHeimdallHash(header.RootHash),
			hmCommonTypes.HexToHeimdallHash("123123"),
			uint64(1),
		)

		result, err := suite.handler(ctx, &msgCheckpointAck)
		require.Error(t, err)
		require.Nil(t, result)
	})

	suite.Run("Invalid Roothash", func() {
		msgCheckpointAck := types.NewMsgCheckpointAck(
			hmCommonTypes.HexToHeimdallAddress("123").Bytes(),
			headerId,
			sdk.AccAddress(header.Proposer),
			header.StartBlock,
			header.EndBlock,
			hmCommonTypes.HexToHeimdallHash("9887"),
			hmCommonTypes.HexToHeimdallHash("123123"),
			uint64(1),
		)

		result, err := suite.handler(ctx, &msgCheckpointAck)
		require.NotNil(t, err)
		require.Nil(t, result)
		require.Error(t, err)
	})
}

func (suite *HandlerTestSuite) TestHandleMsgCheckpointNoAck() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := initApp.CheckpointKeeper
	stakingKeeper := initApp.StakingKeeper
	topupKeeper := initApp.TopupKeeper
	start := uint64(0)
	maxSize := uint64(256)
	params := keeper.GetParams(ctx)
	checkpointBufferTime := params.CheckpointBufferTime

	dividendAccount := hmTypes.DividendAccount{
		User:      hmCommonTypes.HexToHeimdallAddress("123").String(),
		FeeAmount: big.NewInt(0).String(),
	}
	err := topupKeeper.AddDividendAccount(ctx, dividendAccount)
	require.NoError(t, err)

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

	got, err := suite.SendCheckpoint(header)
	require.NotNil(t, got)
	require.Nil(t, err)

	// set time lastCheckpoint timestamp + checkpointBufferTime
	newTime := lastCheckpoint.TimeStamp + uint64(checkpointBufferTime)
	suite.ctx = ctx.WithBlockTime(time.Unix(0, int64(newTime)))
	result, err := suite.SendNoAck()
	require.NotNil(t, result)
	require.Nil(t, err)
	ackCount := keeper.GetACKCount(ctx)
	require.Equal(t, uint64(0), uint64(ackCount), "Should not update state")
}

func (suite *HandlerTestSuite) TestHandleMsgCheckpointNoAckBeforeBufferTimeout() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := initApp.CheckpointKeeper
	stakingKeeper := initApp.StakingKeeper
	topupKeeper := initApp.TopupKeeper
	start := uint64(0)
	maxSize := uint64(256)
	params := keeper.GetParams(ctx)
	dividendAccount := hmTypes.DividendAccount{
		User:      hmCommonTypes.HexToHeimdallAddress("123").String(),
		FeeAmount: big.NewInt(0).String(),
	}
	err := topupKeeper.AddDividendAccount(ctx, dividendAccount)
	require.NoError(t, err)

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

	got, err := suite.SendCheckpoint(header)
	require.NotNil(t, got)
	require.NoError(t, err)

	result, err := suite.SendNoAck()
	require.Nil(t, result)
	require.NotNil(t, err)
}

func (suite *HandlerTestSuite) SendCheckpoint(header *hmTypes.Checkpoint) (res *sdk.Result, err error) {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	// keeper := app.CheckpointKeeper
	topupKeeper := initApp.TopupKeeper

	dividendAccounts := topupKeeper.GetAllDividendAccounts(ctx)
	accRootHash, err := types.GetAccountRootHash(dividendAccounts)
	require.NoError(t, err)
	accountRoot := hmCommonTypes.BytesToHeimdallHash(accRootHash)

	borChainId := "1234"
	// create checkpoint msg
	proposer, err := sdk.AccAddressFromHex(header.Proposer)
	require.NoError(t, err)
	msgCheckpoint := types.NewMsgCheckpointBlock(
		proposer,
		header.StartBlock,
		header.EndBlock,
		hmCommonTypes.HexToHeimdallHash(header.RootHash),
		accountRoot,
		borChainId,
	)

	suite.contractCaller.On("CheckIfBlocksExist", header.EndBlock).Return(true)
	suite.contractCaller.On("GetRootHash", header.StartBlock, header.EndBlock, uint64(1024)).Return(hmCommonTypes.HexToHeimdallHash(header.RootHash).Bytes(), nil)

	// send checkpoint to handler
	result, err := suite.handler(ctx, &msgCheckpoint)
	sideResult := suite.sideHandler(ctx, &msgCheckpoint)
	_, _ = suite.postHandler(ctx, &msgCheckpoint, sideResult.Result)
	return result, err
}

func (suite *HandlerTestSuite) SendNoAck() (res *sdk.Result, err error) {
	_, _, ctx := suite.T(), suite.app, suite.ctx
	msgNoAck := types.NewMsgCheckpointNoAck(hmCommonTypes.HexToHeimdallAddress("123").Bytes())

	result, err := suite.handler(ctx, &msgNoAck)
	sideResult := suite.sideHandler(ctx, &msgNoAck)
	_, _ = suite.postHandler(ctx, &msgNoAck, sideResult.Result)
	return result, err
}
