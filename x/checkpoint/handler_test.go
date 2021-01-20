package checkpoint_test

import (
	"fmt"
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
	"math/big"
	"testing"
	"time"
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
	accountRoot := hmCommonTypes.BytesToHeimdallHash(accRootHash)
	suite.Run("Success", func() {
		msgCheckpoint := types.NewMsgCheckpointBlock(
			sdk.AccAddress(header.Proposer),
			header.StartBlock,
			header.EndBlock,
			hmCommonTypes.BytesToHeimdallHash(header.RootHash),
			accountRoot,
			borChainId,
		)

		// send checkpoint to handler
		result, err := suite.handler(ctx, &msgCheckpoint)
		fmt.Printf("error %+v\n", err)
		fmt.Printf("resukt %+v\n", result)
		//require.True(t, got.IsOK(), "expected send-checkpoint to be ok, got %v", got)
		bufferedHeader, _ := keeper.GetCheckpointFromBuffer(ctx)
		require.Empty(t, bufferedHeader, "Should not store state")
	})

	suite.Run("Invalid Proposer", func() {
		header.Proposer = hmCommonTypes.HexToHeimdallAddress("1234").String()
		msgCheckpoint := types.NewMsgCheckpointBlock(
			sdk.AccAddress(header.Proposer),
			header.StartBlock,
			header.EndBlock,
			hmCommonTypes.BytesToHeimdallHash(header.RootHash),
			accountRoot,
			borChainId,
		)

		// send checkpoint to handler
		result, err := suite.handler(ctx, &msgCheckpoint)
		fmt.Printf("Result is %+v\n", result)
		fmt.Printf("error is %+v\n", err)
		//require.True(t, !got.IsOK(), errs.CodeToDefaultMsg(got.Code))
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
			sdk.AccAddress(header.Proposer),
			start,
			start+256,
			hmCommonTypes.BytesToHeimdallHash(header.RootHash),
			accountRoot,
			borChainId,
		)

		// send checkpoint to handler
		result, err := suite.handler(ctx, &msgCheckpoint)
		fmt.Printf("error %+v\n", err)
		fmt.Printf("resukt %+v\n", result)

		//require.True(t, !got.IsOK(), errs.CodeToDefaultMsg(got.Code))
	})
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
		User:      hmCommonTypes.HexToHeimdallAddress("123").String(),
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

	got, err := suite.SendCheckpoint(header)
	//require.NotNil(t, got)
	//require.Nil(t, err)
	fmt.Printf("got %+v\n err %+vn", got, err)

	bufferedHeader, err := keeper.GetCheckpointFromBuffer(ctx)
	require.NotNil(t, bufferedHeader)
	require.Nil(t, err)

	// send ack
	headerId := uint64(1)
	suite.Run("success", func() {
		msgCheckpointAck := types.NewMsgCheckpointAck(
			hmCommonTypes.HexToHeimdallAddress("123").Bytes(),
			headerId,
			sdk.AccAddress(header.Proposer),
			header.StartBlock,
			header.EndBlock,
			hmCommonTypes.BytesToHeimdallHash(header.RootHash),
			hmCommonTypes.HexToHeimdallHash("123123"),
			uint64(1),
		)
		result, err := suite.handler(ctx, &msgCheckpointAck)
		fmt.Printf("error %+v\n", err)
		fmt.Printf("resukt %+v\n", result)
		//require.NotNil(t, result)
		//require.Nil(t, err)

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
			hmCommonTypes.BytesToHeimdallHash(header.RootHash),
			hmCommonTypes.HexToHeimdallHash("123123"),
			uint64(1),
		)

		result, err := suite.handler(ctx, &msgCheckpointAck)
		fmt.Printf("error %+v\n", err)
		fmt.Printf("resukt %+v\n", result)
		//require.NotNil(t, result)
		//require.Nil(t, err)

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
		fmt.Printf("error %+v\n", err)
		fmt.Printf("resukt %+v\n", result)
		//require.NotNil(t, result)
		//require.Nil(t, err)

		//require.True(t, !got.IsOK(), errs.CodeToDefaultMsg(got.Code))
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
		User:      hmCommonTypes.HexToHeimdallAddress("123").String(),
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

	got, err := suite.SendCheckpoint(header)
	fmt.Printf("error %+v\n", err)
	fmt.Printf("resukt %+v\n", got)
	//require.NotNil(t, got)
	//require.Nil(t, err)

	fmt.Printf("result no ack %+v\n %+v\n", got, err)

	//require.True(t, got.IsOK(), "expected send-NoAck to be ok, got %v", got)

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
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper
	stakingKeeper := app.StakingKeeper
	topupKeeper := app.TopupKeeper
	start := uint64(0)
	maxSize := uint64(256)
	params := keeper.GetParams(ctx)
	dividendAccount := hmTypes.DividendAccount{
		User:      hmCommonTypes.HexToHeimdallAddress("123").String(),
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

	got, err := suite.SendCheckpoint(header)
	fmt.Printf("got %+v\n and err %+v\n", got, err)
	//require.Nil(t, got)
	//require.NotNil(t, err)

	result, err := suite.SendNoAck()
	fmt.Printf("result %+v\n and err %+v\n", result, err)
	//require.Nil(t, result)
	//require.NotNil(t, err)
}

func (suite *HandlerTestSuite) SendCheckpoint(header *hmTypes.Checkpoint) (res *sdk.Result, err error) {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	// keeper := app.CheckpointKeeper
	topupKeeper := app.TopupKeeper

	dividendAccounts := topupKeeper.GetAllDividendAccounts(ctx)
	accRootHash, err := types.GetAccountRootHash(dividendAccounts)
	require.NoError(t, err)
	accountRoot := hmCommonTypes.BytesToHeimdallHash(accRootHash)

	borChainId := "1234"
	// create checkpoint msg
	msgCheckpoint := types.NewMsgCheckpointBlock(
		sdk.AccAddress(header.Proposer),
		header.StartBlock,
		header.EndBlock,
		hmCommonTypes.BytesToHeimdallHash(header.RootHash),
		accountRoot,
		borChainId,
	)

	suite.contractCaller.On("CheckIfBlocksExist", header.EndBlock).Return(true)
	suite.contractCaller.On("GetRootHash", header.StartBlock, header.EndBlock, uint64(1024)).Return(header.RootHash, nil)

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
