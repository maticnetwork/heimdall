package checkpoint_test

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/maticnetwork/heimdall/x/checkpoint/test_helper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/helper/mocks"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/x/checkpoint"
	chSim "github.com/maticnetwork/heimdall/x/checkpoint/simulation"
	"github.com/maticnetwork/heimdall/x/checkpoint/types"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/proto/tendermint/types"
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
	suite.app, suite.ctx, _ = test_helper.CreateTestApp(false)
	suite.contractCaller = mocks.IContractCaller{}
	suite.sideHandler = checkpoint.NewSideTxHandler(suite.app.CheckpointKeeper, &suite.contractCaller)
	suite.postHandler = checkpoint.NewPostTxHandler(suite.app.CheckpointKeeper, &suite.contractCaller)

	// random generator
	s1 := rand.NewSource(time.Now().UnixNano())
	suite.r = rand.New(s1)
}

func TestSideHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(SideHandlerTestSuite))
}

//
// Test cases
//

func (suite *SideHandlerTestSuite) TestSideHandler() {
	t, ctx := suite.T(), suite.ctx

	// side handler
	result := suite.sideHandler(ctx, nil)
	require.Equal(t, sdkerrors.ErrUnknownRequest.ABCICode(), result.Code)
	require.Equal(t, abci.SideTxResultType_SKIP, result.Result)
}

// test handler for message
func (suite *SideHandlerTestSuite) TestSideHandleMsgCheckpoint() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := initApp.CheckpointKeeper

	start := uint64(0)
	maxSize := uint64(256)
	params := keeper.GetParams(ctx)

	header, err := chSim.GenRandCheckpoint(start, maxSize, params.MaxCheckpointLength)
	require.NoError(t, err)
	borChainId := "1234"
	suite.Run("Success", func() {
		suite.contractCaller = mocks.IContractCaller{}

		proposer, err := sdk.AccAddressFromHex(header.Proposer)
		require.NoError(t, err)

		// create checkpoint msg
		msgCheckpoint := types.NewMsgCheckpointBlock(
			proposer,
			header.StartBlock,
			header.EndBlock,
			hmCommonTypes.HexToHeimdallHash(header.RootHash),
			hmCommonTypes.HexToHeimdallHash(header.RootHash),
			borChainId,
		)

		suite.contractCaller.On("CheckIfBlocksExist", header.EndBlock).Return(true)
		suite.contractCaller.On("GetRootHash", header.StartBlock, header.EndBlock, uint64(1024)).Return(hmCommonTypes.HexToHeimdallHash(header.RootHash).Bytes(), nil)

		result := suite.sideHandler(ctx, &msgCheckpoint)
		// todo:   uint32(0) needs convert sdk.CodeOK
		require.Equal(t, uint32(0), result.Code, "Side tx handler should be success")
		require.Equal(t, abci.SideTxResultType_YES, result.Result, "Result should be `yes`")

		bufferedHeader, _ := keeper.GetCheckpointFromBuffer(ctx)
		require.Nil(t, bufferedHeader, "Should not store state")
	})

	suite.Run("No root hash", func() {
		suite.contractCaller = mocks.IContractCaller{}
		proposer, err := sdk.AccAddressFromHex(header.Proposer)
		require.NoError(t, err)
		// create checkpoint msg
		msgCheckpoint := types.NewMsgCheckpointBlock(
			proposer,
			header.StartBlock,
			header.EndBlock,
			hmCommonTypes.HexToHeimdallHash(header.RootHash),
			hmCommonTypes.HexToHeimdallHash(header.RootHash),
			borChainId,
		)

		suite.contractCaller.On("CheckIfBlocksExist", header.EndBlock).Return(true)
		suite.contractCaller.On("GetRootHash", header.StartBlock, header.EndBlock, uint64(1024)).Return(nil, nil)

		result := suite.sideHandler(ctx, &msgCheckpoint)
		// 0 needs convert sdk.CodeOK
		require.Equal(t, uint32(0), result.Code, "Side tx handler should succeed")
		require.Equal(t, abci.SideTxResultType_SKIP, result.Result, "Result should be `skip`")

		bufferedHeader, err := keeper.GetCheckpointFromBuffer(ctx)
		require.Error(t, err)
		require.Nil(t, bufferedHeader, "Should not store state")
	})
	//

	suite.Run("invalid checkpoint", func() {
		suite.contractCaller = mocks.IContractCaller{}
		proposer, err := sdk.AccAddressFromHex(header.Proposer)
		require.NoError(t, err)
		// create checkpoint msg
		msgCheckpoint := types.NewMsgCheckpointBlock(
			proposer,
			header.StartBlock,
			header.EndBlock,
			hmCommonTypes.HexToHeimdallHash(header.RootHash),
			hmCommonTypes.HexToHeimdallHash(header.RootHash),
			borChainId,
		)

		suite.contractCaller.On("CheckIfBlocksExist", header.EndBlock).Return(true)
		suite.contractCaller.On("GetRootHash", header.StartBlock, header.EndBlock, uint64(1024)).Return([]byte{1}, nil)

		result := suite.sideHandler(ctx, &msgCheckpoint)
		require.NotNil(t, result)
		require.Equal(t, uint32(0), result.Code, "Side tx handler should not succeed")
		require.Equal(t, abci.SideTxResultType_SKIP, result.Result, "Result should be `skip`")
	})
}

func (suite *SideHandlerTestSuite) TestSideHandleMsgCheckpointAck() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := initApp.CheckpointKeeper
	start := uint64(0)
	maxSize := uint64(256)
	params := keeper.GetParams(ctx)

	header, _ := chSim.GenRandCheckpoint(start, maxSize, params.MaxCheckpointLength)
	headerId := uint64(1)
	suite.Run("Success", func() {
		suite.contractCaller = mocks.IContractCaller{}

		// prepare ack msg
		accAddr, err := sdk.AccAddressFromHex("123")
		require.NoError(t, err)
		proposer, err := sdk.AccAddressFromHex(header.Proposer)
		require.NoError(t, err)
		msgCheckpointAck := types.NewMsgCheckpointAck(
			accAddr,
			uint64(1),
			proposer,
			header.StartBlock,
			header.EndBlock,
			hmCommonTypes.HexToHeimdallHash(header.RootHash),
			hmCommonTypes.HexToHeimdallHash("123123"),
			uint64(1),
		)
		rootChainInstance := &rootchain.Rootchain{}

		suite.contractCaller.On("GetRootChainInstance", mock.Anything).Return(rootChainInstance, nil)
		suite.contractCaller.On("GetHeaderInfo", headerId, rootChainInstance, params.ChildBlockInterval).Return(common.HexToHash(header.RootHash), header.StartBlock, header.EndBlock, header.TimeStamp, accAddr, nil)

		result := suite.sideHandler(ctx, &msgCheckpointAck)
		require.Equal(t, uint32(0), result.Code, "Side tx handler should be success")
		require.Equal(t, abci.SideTxResultType_SKIP, result.Result, "Result should be `yes`")
	})

	suite.Run("No HeaderInfo", func() {
		suite.contractCaller = mocks.IContractCaller{}
		accAddr, err := sdk.AccAddressFromHex("123")
		require.NoError(t, err)
		proposer, err := sdk.AccAddressFromHex(header.Proposer)
		require.NoError(t, err)
		// prepare ack msg
		msgCheckpointAck := types.NewMsgCheckpointAck(
			accAddr,
			uint64(1),
			proposer,
			header.StartBlock,
			header.EndBlock,
			hmCommonTypes.HexToHeimdallHash("123"),
			hmCommonTypes.HexToHeimdallHash("123123"),
			uint64(1),
		)
		rootChainInstance := &rootchain.Rootchain{}

		suite.contractCaller.On("GetRootChainInstance", mock.Anything).Return(rootChainInstance, nil)
		suite.contractCaller.On("GetHeaderInfo", headerId, rootChainInstance, params.ChildBlockInterval).Return(nil, header.StartBlock, header.EndBlock, header.TimeStamp, accAddr, nil)

		result := suite.sideHandler(ctx, &msgCheckpointAck)
		require.Equal(t, uint32(0), result.Code, "Side tx handler should succeed")
		require.Equal(t, abci.SideTxResultType_SKIP, result.Result, "Result should skip")
	})
}

func (suite *SideHandlerTestSuite) TestPostHandler() {
	t, ctx := suite.T(), suite.ctx

	// side handler
	_, err := suite.postHandler(ctx, nil, abci.SideTxResultType_YES)
	require.Error(t, err, "Post handler should fail")
}

func (suite *SideHandlerTestSuite) TestPostHandleMsgCheckpoint() {
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
	require.NoError(t, err, "failed to generate random checkpoint")
	// add current proposer to header
	header.Proposer = stakingKeeper.GetValidatorSet(ctx).Proposer.Signer

	borChainId := "1234"
	suite.Run("Failure", func() {
		proposer, err := sdk.AccAddressFromHex(header.Proposer)
		require.NoError(t, err)
		// create checkpoint msg
		msgCheckpoint := types.NewMsgCheckpointBlock(
			proposer,
			header.StartBlock,
			header.EndBlock,
			hmCommonTypes.HexToHeimdallHash(header.RootHash),
			hmCommonTypes.HexToHeimdallHash(header.RootHash),
			borChainId,
		)

		_, err = suite.postHandler(ctx, &msgCheckpoint, abci.SideTxResultType_NO)
		require.Error(t, err)

		bufferedHeader, err := keeper.GetCheckpointFromBuffer(ctx)
		require.Nil(t, bufferedHeader)
		require.Error(t, err)
	})

	suite.Run("Success", func() {
		proposer, err := sdk.AccAddressFromHex(header.Proposer)
		require.NoError(t, err)
		// create checkpoint msg
		msgCheckpoint := types.NewMsgCheckpointBlock(
			proposer,
			header.StartBlock,
			header.EndBlock,
			hmCommonTypes.HexToHeimdallHash(header.RootHash),
			hmCommonTypes.HexToHeimdallHash(header.RootHash),
			borChainId,
		)

		_, err = suite.postHandler(ctx, &msgCheckpoint, abci.SideTxResultType_YES)
		require.NoError(t, err)

		bufferedHeader, err := keeper.GetCheckpointFromBuffer(ctx)
		require.Equal(t, bufferedHeader.StartBlock, header.StartBlock)
		require.Equal(t, bufferedHeader.EndBlock, header.EndBlock)
		require.Equal(t, bufferedHeader.RootHash, header.RootHash)
		require.Equal(t, bufferedHeader.Proposer, header.Proposer)
		require.Equal(t, bufferedHeader.BorChainID, header.BorChainID)
		require.Empty(t, err, "Unable to set checkpoint from buffer, Error: %v", err)
	})

	suite.Run("Replay", func() {
		// create checkpoint msg
		proposer, err := sdk.AccAddressFromHex(header.Proposer)
		require.NoError(t, err)
		msgCheckpoint := types.NewMsgCheckpointBlock(
			proposer,
			header.StartBlock,
			header.EndBlock,
			hmCommonTypes.HexToHeimdallHash(header.RootHash),
			hmCommonTypes.HexToHeimdallHash(header.RootHash),
			borChainId,
		)

		result, err := suite.postHandler(ctx, &msgCheckpoint, abci.SideTxResultType_YES)
		require.Error(t, err)
		require.Nil(t, result)
	})
}

func (suite *SideHandlerTestSuite) TestPostHandleMsgCheckpointAck() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	keeper := initApp.CheckpointKeeper

	start := uint64(0)
	maxSize := uint64(256)
	params := keeper.GetParams(ctx)
	header, _ := chSim.GenRandCheckpoint(start, maxSize, params.MaxCheckpointLength)
	// generate proposer for validator set
	chSim.LoadValidatorSet(2, t, initApp.StakingKeeper, ctx, false, 10)
	initApp.StakingKeeper.IncrementAccum(ctx, 1)

	// send ack
	checkpointNumber := uint64(1)

	suite.Run("Failure", func() {
		accAddr, err := sdk.AccAddressFromHex("123")
		require.NoError(t, err)
		proposer, err := sdk.AccAddressFromHex(header.Proposer)
		require.NoError(t, err)

		msgCheckpointAck := types.NewMsgCheckpointAck(
			accAddr,
			checkpointNumber,
			proposer,
			header.StartBlock,
			header.EndBlock,
			hmCommonTypes.HexToHeimdallHash(header.RootHash),
			hmCommonTypes.HexToHeimdallHash("123123"),
			uint64(1),
		)

		result, err := suite.postHandler(ctx, &msgCheckpointAck, abci.SideTxResultType_NO)
		require.Error(t, err)
		require.Nil(t, result)

		afterAckBufferedCheckpoint, _ := keeper.GetCheckpointFromBuffer(ctx)
		require.Nil(t, afterAckBufferedCheckpoint)
	})

	suite.Run("Success", func() {
		accAddr, err := sdk.AccAddressFromHex("123")
		require.NoError(t, err)
		proposer, err := sdk.AccAddressFromHex(header.Proposer)
		require.NoError(t, err)

		msgCheckpoint := types.NewMsgCheckpointBlock(
			proposer,
			header.StartBlock,
			header.EndBlock,
			hmCommonTypes.HexToHeimdallHash(header.RootHash),
			hmCommonTypes.HexToHeimdallHash(header.RootHash),
			"1234",
		)

		_, err = suite.postHandler(ctx, &msgCheckpoint, abci.SideTxResultType_YES)
		require.NoError(t, err)

		msgCheckpointAck := types.NewMsgCheckpointAck(
			accAddr,
			checkpointNumber,
			proposer,
			header.StartBlock,
			header.EndBlock,
			hmCommonTypes.HexToHeimdallHash(header.RootHash),
			hmCommonTypes.HexToHeimdallHash("123123"),
			uint64(1),
		)

		_, err = suite.postHandler(ctx, &msgCheckpointAck, abci.SideTxResultType_YES)
		require.NoError(t, err)

		afterAckBufferedCheckpoint, _ := keeper.GetCheckpointFromBuffer(ctx)
		require.Nil(t, afterAckBufferedCheckpoint)
	})

	suite.Run("Replay", func() {
		accAddr, err := sdk.AccAddressFromHex("123")
		require.NoError(t, err)
		proposer, err := sdk.AccAddressFromHex(header.Proposer)
		require.NoError(t, err)
		msgCheckpointAck := types.NewMsgCheckpointAck(
			accAddr,
			checkpointNumber,
			proposer,
			header.StartBlock,
			header.EndBlock,
			hmCommonTypes.HexToHeimdallHash(header.RootHash),
			hmCommonTypes.HexToHeimdallHash("123123"),
			uint64(1),
		)

		_, err = suite.postHandler(ctx, &msgCheckpointAck, abci.SideTxResultType_YES)
		require.Error(t, err)

		afterAckBufferedCheckpoint, _ := keeper.GetCheckpointFromBuffer(ctx)
		require.Nil(t, afterAckBufferedCheckpoint)
	})

	suite.Run("InvalidEndBlock", func() {
		suite.contractCaller = mocks.IContractCaller{}
		header2, _ := chSim.GenRandCheckpoint(header.EndBlock+1, maxSize, params.MaxCheckpointLength)
		proposer2, err := sdk.AccAddressFromHex(header2.Proposer)
		require.NoError(t, err)
		msgCheckpoint := types.NewMsgCheckpointBlock(
			proposer2,
			header2.StartBlock,
			header2.EndBlock,
			hmCommonTypes.HexToHeimdallHash(header2.RootHash),
			hmCommonTypes.HexToHeimdallHash(header2.RootHash),
			"1234",
		)

		_, err = suite.postHandler(ctx, &msgCheckpoint, abci.SideTxResultType_YES)
		require.NoError(t, err)

		accAddr, err := sdk.AccAddressFromHex("123123")
		require.NoError(t, err)
		msgCheckpointAck := types.NewMsgCheckpointAck(
			accAddr,
			checkpointNumber,
			proposer2,
			header2.StartBlock,
			header2.EndBlock,
			hmCommonTypes.HexToHeimdallHash(header2.RootHash),
			hmCommonTypes.HexToHeimdallHash("123123"),
			uint64(1),
		)

		_, err = suite.postHandler(ctx, &msgCheckpointAck, abci.SideTxResultType_YES)
		require.NoError(t, err)

		afterAckBufferedCheckpoint, _ := keeper.GetCheckpointFromBuffer(ctx)
		require.Nil(t, afterAckBufferedCheckpoint)
	})
}
