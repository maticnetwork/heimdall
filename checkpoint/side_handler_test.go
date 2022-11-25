package checkpoint_test

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	borCommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/app"
	cmTypes "github.com/maticnetwork/heimdall/chainmanager/types"
	"github.com/maticnetwork/heimdall/checkpoint"
	chSim "github.com/maticnetwork/heimdall/checkpoint/simulation"
	"github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/common"
	errs "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/helper/mocks"
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
	suite.sideHandler = checkpoint.NewSideTxHandler(suite.app.CheckpointKeeper, &suite.contractCaller)
	suite.postHandler = checkpoint.NewPostTxHandler(suite.app.CheckpointKeeper, &suite.contractCaller)

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
	err := keeper.AddCheckpoint(ctx, 1, checkpoint)
	require.NoError(t, err)

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
	err := keeper.AddCheckpoint(ctx, 1, checkpoint)
	require.NoError(t, err)

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
	err := keeper.AddCheckpoint(ctx, 1, checkpoint)
	require.NoError(t, err)

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
	chSim.LoadValidatorSet(t, 2, app.StakingKeeper, ctx, false, 10)
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

func (suite *SideHandlerTestSuite) TestPostHandleMsgMilestone() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper
	stakingKeeper := app.StakingKeeper
	start := uint64(0)
	milestoneLength := helper.MilestoneLength

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
