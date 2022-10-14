package checkpoint_test

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	borCommon "github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/heimdall/app"
	cmTypes "github.com/maticnetwork/heimdall/chainmanager/types"
	"github.com/maticnetwork/heimdall/checkpoint"
	chSim "github.com/maticnetwork/heimdall/checkpoint/simulation"
	"github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/common"
	errs "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/helper/mocks"
	hmTypes "github.com/maticnetwork/heimdall/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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
	require.Equal(t, uint32(sdk.CodeUnknownRequest), result.Code)
	require.Equal(t, abci.SideTxResultType_Skip, result.Result)
}

// test handler for message
func (suite *HandlerTestSuite) TestHandleMsgCheckpointAdjustSuccess() {
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
		Proposer:    hmTypes.HexToHeimdallAddress("456"),
		StartBlock:  0,
		EndBlock:    512,
		RootHash:    hmTypes.HexToHeimdallHash("456"),
	}
	rootchainInstance := &rootchain.Rootchain{}
	suite.contractCaller.On("GetRootChainInstance", mock.Anything).Return(rootchainInstance, nil)
	suite.contractCaller.On("GetHeaderInfo", mock.Anything, mock.Anything, mock.Anything).Return(borCommon.HexToHash("456"), uint64(0), uint64(512), uint64(1), hmTypes.HexToHeimdallAddress("456"), nil)

	suite.handler(ctx, checkpointAdjust)
	sideResult := suite.sideHandler(ctx, checkpointAdjust)
	suite.postHandler(ctx, checkpointAdjust, sideResult.Result)

	responseCheckpoint, _ := keeper.GetCheckpointByNumber(ctx, 1)
	require.Equal(t, responseCheckpoint.EndBlock, uint64(512))
	require.Equal(t, responseCheckpoint.Proposer, hmTypes.HexToHeimdallAddress("456"))
	require.Equal(t, responseCheckpoint.RootHash, hmTypes.HexToHeimdallHash("456"))
}

func (suite *HandlerTestSuite) TestHandleMsgCheckpointAdjustSameCheckpointAsRootChain() {
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
		RootHash:    hmTypes.HexToHeimdallHash("456"),
	}
	rootchainInstance := &rootchain.Rootchain{}
	suite.contractCaller.On("GetRootChainInstance", mock.Anything).Return(rootchainInstance, nil)
	suite.contractCaller.On("GetHeaderInfo", mock.Anything, mock.Anything, mock.Anything).Return(borCommon.HexToHash("123"), uint64(0), uint64(256), uint64(1), hmTypes.HexToHeimdallAddress("123"), nil)

	suite.handler(ctx, checkpointAdjust)
	sideResult := suite.sideHandler(ctx, checkpointAdjust)
	require.Equal(t, sideResult.Code, uint32(common.CodeCheckpointAlreadyExists))
}

func (suite *HandlerTestSuite) TestHandleMsgCheckpointAdjustNotSameCheckpointAsRootChain() {
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

	rootchainInstance := &rootchain.Rootchain{}
	suite.contractCaller.On("GetRootChainInstance", mock.Anything).Return(rootchainInstance, nil)
	suite.contractCaller.On("GetHeaderInfo", mock.Anything, mock.Anything, mock.Anything).Return(borCommon.HexToHash("222"), uint64(0), uint64(256), uint64(1), hmTypes.HexToHeimdallAddress("123"), nil)

	result := suite.sideHandler(ctx, checkpointAdjust)
	require.Equal(t, result.Code, uint32(common.CodeCheckpointAlreadyExists))
}

func (suite *SideHandlerTestSuite) TestSideHandleMsgCheckpoint() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper

	start := uint64(0)
	maxSize := uint64(256)
	params := keeper.GetParams(ctx)

	header, err := chSim.GenRandCheckpoint(start, maxSize, params.MaxCheckpointLength)
	require.NoError(t, err)
	borChainId := "1234"
	suite.Run("Success", func() {
		suite.contractCaller = mocks.IContractCaller{}

		// create checkpoint msg
		msgCheckpoint := types.NewMsgCheckpointBlock(
			header.Proposer,
			header.StartBlock,
			header.EndBlock,
			header.RootHash,
			header.RootHash,
			borChainId,
		)

		suite.contractCaller.On("CheckIfBlocksExist", header.EndBlock+cmTypes.DefaultMaticchainTxConfirmations).Return(true)
		suite.contractCaller.On("GetRootHash", header.StartBlock, header.EndBlock, uint64(1024)).Return(header.RootHash.Bytes(), nil)

		result := suite.sideHandler(ctx, msgCheckpoint)
		require.Equal(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should be success")

		bufferedHeader, _ := keeper.GetCheckpointFromBuffer(ctx)
		require.Nil(t, bufferedHeader, "Should not store state")
	})

	suite.Run("No Roothash", func() {
		suite.contractCaller = mocks.IContractCaller{}

		// create checkpoint msg
		msgCheckpoint := types.NewMsgCheckpointBlock(
			header.Proposer,
			header.StartBlock,
			header.EndBlock,
			header.RootHash,
			header.RootHash,
			borChainId,
		)

		suite.contractCaller.On("CheckIfBlocksExist", header.EndBlock+cmTypes.DefaultMaticchainTxConfirmations).Return(true)
		suite.contractCaller.On("GetRootHash", header.StartBlock, header.EndBlock, uint64(1024)).Return(nil, nil)

		result := suite.sideHandler(ctx, msgCheckpoint)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should Fail")
		require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should be `skip`")
		require.Equal(t, uint32(common.CodeInvalidBlockInput), result.Code)

		bufferedHeader, err := keeper.GetCheckpointFromBuffer(ctx)
		require.Error(t, err)
		require.Nil(t, bufferedHeader, "Should not store state")
	})

	suite.Run("invalid checkpoint", func() {
		suite.contractCaller = mocks.IContractCaller{}

		// create checkpoint msg
		msgCheckpoint := types.NewMsgCheckpointBlock(
			header.Proposer,
			header.StartBlock,
			header.EndBlock,
			header.RootHash,
			header.RootHash,
			borChainId,
		)

		suite.contractCaller.On("CheckIfBlocksExist", header.EndBlock+cmTypes.DefaultMaticchainTxConfirmations).Return(true)
		suite.contractCaller.On("GetRootHash", header.StartBlock, header.EndBlock, uint64(1024)).Return([]byte{1}, nil)

		result := suite.sideHandler(ctx, msgCheckpoint)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should fail")
		require.Equal(t, uint32(common.CodeInvalidBlockInput), result.Code)
	})
}

func (suite *SideHandlerTestSuite) TestSideHandleMsgCheckpointAck() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper
	start := uint64(0)
	maxSize := uint64(256)
	params := keeper.GetParams(ctx)

	header, _ := chSim.GenRandCheckpoint(start, maxSize, params.MaxCheckpointLength)
	headerId := uint64(1)
	suite.Run("Success", func() {
		suite.contractCaller = mocks.IContractCaller{}

		// prepare ack msg
		msgCheckpointAck := types.NewMsgCheckpointAck(
			hmTypes.HexToHeimdallAddress("123"),
			uint64(1),
			header.Proposer,
			header.StartBlock,
			header.EndBlock,
			header.RootHash,
			hmTypes.HexToHeimdallHash("123123"),
			uint64(1),
		)
		rootchainInstance := &rootchain.Rootchain{}

		suite.contractCaller.On("GetRootChainInstance", mock.Anything).Return(rootchainInstance, nil)
		suite.contractCaller.On("GetHeaderInfo", headerId, rootchainInstance, params.ChildBlockInterval).Return(header.RootHash.EthHash(), header.StartBlock, header.EndBlock, header.TimeStamp, header.Proposer, nil)

		result := suite.sideHandler(ctx, msgCheckpointAck)
		require.Equal(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should be success")
		require.Equal(t, abci.SideTxResultType_Yes, result.Result, "Result should be `yes`")
	})

	suite.Run("No HeaderInfo", func() {
		suite.contractCaller = mocks.IContractCaller{}

		// prepare ack msg
		msgCheckpointAck := types.NewMsgCheckpointAck(
			hmTypes.HexToHeimdallAddress("123"),
			uint64(1),
			header.Proposer,
			header.StartBlock,
			header.EndBlock,
			hmTypes.HexToHeimdallHash("123"),
			hmTypes.HexToHeimdallHash("123123"),
			uint64(1),
		)
		rootchainInstance := &rootchain.Rootchain{}

		suite.contractCaller.On("GetRootChainInstance", mock.Anything).Return(rootchainInstance, nil)
		suite.contractCaller.On("GetHeaderInfo", headerId, rootchainInstance, params.ChildBlockInterval).Return(nil, header.StartBlock, header.EndBlock, header.TimeStamp, header.Proposer, nil)

		result := suite.sideHandler(ctx, msgCheckpointAck)
		require.NotEqual(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should fail")
		require.Equal(t, abci.SideTxResultType_Skip, result.Result, "Result should skip")
	})
}

func (suite *SideHandlerTestSuite) TestPostHandler() {
	t, ctx := suite.T(), suite.ctx

	// side handler
	result := suite.postHandler(ctx, nil, abci.SideTxResultType_Yes)
	require.False(t, result.IsOK(), "Post handler should fail")
	require.Equal(t, sdk.CodeUnknownRequest, result.Code)
}

func (suite *SideHandlerTestSuite) TestPostHandleMsgCheckpoint() {
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

	borChainId := "1234"
	suite.Run("Failure", func() {
		// create checkpoint msg
		msgCheckpoint := types.NewMsgCheckpointBlock(
			header.Proposer,
			header.StartBlock,
			header.EndBlock,
			header.RootHash,
			header.RootHash,
			borChainId,
		)

		result := suite.postHandler(ctx, msgCheckpoint, abci.SideTxResultType_No)
		require.True(t, !result.IsOK(), errs.CodeToDefaultMsg(result.Code))

		bufferedHeader, err := keeper.GetCheckpointFromBuffer(ctx)
		require.Nil(t, bufferedHeader)
		require.Error(t, err)
	})

	suite.Run("Success", func() {
		// create checkpoint msg
		msgCheckpoint := types.NewMsgCheckpointBlock(
			header.Proposer,
			header.StartBlock,
			header.EndBlock,
			header.RootHash,
			header.RootHash,
			borChainId,
		)

		result := suite.postHandler(ctx, msgCheckpoint, abci.SideTxResultType_Yes)
		require.True(t, result.IsOK(), "expected send-checkpoint to be ok, got %v", result)

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
		msgCheckpoint := types.NewMsgCheckpointBlock(
			header.Proposer,
			header.StartBlock,
			header.EndBlock,
			header.RootHash,
			header.RootHash,
			borChainId,
		)

		result := suite.postHandler(ctx, msgCheckpoint, abci.SideTxResultType_Yes)
		require.False(t, result.IsOK(), "expected send-checkpoint to be ok, got %v", result)
		require.Equal(t, common.CodeNoACK, result.Code)
	})
}

func (suite *SideHandlerTestSuite) TestPostHandleMsgCheckpointAck() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper

	start := uint64(0)
	maxSize := uint64(256)
	params := keeper.GetParams(ctx)
	header, _ := chSim.GenRandCheckpoint(start, maxSize, params.MaxCheckpointLength)
	// generate proposer for validator set
	chSim.LoadValidatorSet(2, t, app.StakingKeeper, ctx, false, 10)
	app.StakingKeeper.IncrementAccum(ctx, 1)

	// send ack
	checkpointNumber := uint64(1)

	suite.Run("Failure", func() {
		msgCheckpointAck := types.NewMsgCheckpointAck(
			hmTypes.HexToHeimdallAddress("123"),
			checkpointNumber,
			header.Proposer,
			header.StartBlock,
			header.EndBlock,
			header.RootHash,
			hmTypes.HexToHeimdallHash("123123"),
			uint64(1),
		)

		result := suite.postHandler(ctx, msgCheckpointAck, abci.SideTxResultType_No)
		require.True(t, !result.IsOK(), errs.CodeToDefaultMsg(result.Code))

		afterAckBufferedCheckpoint, _ := keeper.GetCheckpointFromBuffer(ctx)
		require.Nil(t, afterAckBufferedCheckpoint)
	})

	suite.Run("Success", func() {
		msgCheckpoint := types.NewMsgCheckpointBlock(
			header.Proposer,
			header.StartBlock,
			header.EndBlock,
			header.RootHash,
			header.RootHash,
			"1234",
		)

		result := suite.postHandler(ctx, msgCheckpoint, abci.SideTxResultType_Yes)
		require.True(t, result.IsOK(), "expected send-checkpoint to be ok, got %v", result)

		msgCheckpointAck := types.NewMsgCheckpointAck(
			hmTypes.HexToHeimdallAddress("123"),
			checkpointNumber,
			header.Proposer,
			header.StartBlock,
			header.EndBlock,
			header.RootHash,
			hmTypes.HexToHeimdallHash("123123"),
			uint64(1),
		)

		result = suite.postHandler(ctx, msgCheckpointAck, abci.SideTxResultType_Yes)
		require.True(t, result.IsOK(), "expected send-ack to be ok, got %v", result)

		afterAckBufferedCheckpoint, _ := keeper.GetCheckpointFromBuffer(ctx)
		require.Nil(t, afterAckBufferedCheckpoint)
	})

	suite.Run("Replay", func() {
		msgCheckpointAck := types.NewMsgCheckpointAck(
			hmTypes.HexToHeimdallAddress("123"),
			checkpointNumber,
			header.Proposer,
			header.StartBlock,
			header.EndBlock,
			header.RootHash,
			hmTypes.HexToHeimdallHash("123123"),
			uint64(1),
		)

		result := suite.postHandler(ctx, msgCheckpointAck, abci.SideTxResultType_Yes)
		require.False(t, result.IsOK())
		require.Equal(t, common.CodeInvalidACK, result.Code)

		afterAckBufferedCheckpoint, _ := keeper.GetCheckpointFromBuffer(ctx)
		require.Nil(t, afterAckBufferedCheckpoint)
	})

	suite.Run("InvalidEndBlock", func() {
		suite.contractCaller = mocks.IContractCaller{}
		header2, _ := chSim.GenRandCheckpoint(header.EndBlock+1, maxSize, params.MaxCheckpointLength)
		msgCheckpoint := types.NewMsgCheckpointBlock(
			header2.Proposer,
			header2.StartBlock,
			header2.EndBlock,
			header2.RootHash,
			header2.RootHash,
			"1234",
		)

		result := suite.postHandler(ctx, msgCheckpoint, abci.SideTxResultType_Yes)
		require.True(t, result.IsOK(), "expected send-checkpoint to be ok, got %v", result)

		msgCheckpointAck := types.NewMsgCheckpointAck(
			hmTypes.HexToHeimdallAddress("123"),
			checkpointNumber,
			header2.Proposer,
			header2.StartBlock,
			header2.EndBlock,
			header2.RootHash,
			hmTypes.HexToHeimdallHash("123123"),
			uint64(1),
		)

		result = suite.postHandler(ctx, msgCheckpointAck, abci.SideTxResultType_Yes)
		require.True(t, result.IsOK(), "expected send-ack to be ok, got %v", result)

		afterAckBufferedCheckpoint, _ := keeper.GetCheckpointFromBuffer(ctx)
		require.Nil(t, afterAckBufferedCheckpoint)
	})
}
