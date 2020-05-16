package checkpoint_test

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/helper/mocks"
	cmn "github.com/maticnetwork/heimdall/test"
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
func (suite *SideHandlerTestSuite) TestSideHandleMsgCheckpoint() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper
	// stakingKeeper := app.StakingKeeper
	start := uint64(0)
	maxSize := uint64(256)
	params := keeper.GetParams(ctx)

	header, err := suite.GenRandCheckpointHeader(start, maxSize, params.MaxCheckpointLength)
	require.NoError(t, err)
	borChainId := "1234"
	// create checkpoint msg
	msgCheckpoint := types.NewMsgCheckpointBlock(
		header.Proposer,
		header.StartBlock,
		header.EndBlock,
		header.RootHash,
		header.RootHash,
		borChainId,
	)

	suite.contractCaller.On("CheckIfBlocksExist", header.EndBlock).Return(true)
	suite.contractCaller.On("GetRootHash", header.StartBlock, header.EndBlock, uint64(1024)).Return(header.RootHash.Bytes(), nil)

	result := suite.sideHandler(ctx, msgCheckpoint)
	require.Equal(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should be success")

	bufferedHeader, _ := keeper.GetCheckpointFromBuffer(ctx)
	require.Nil(t, bufferedHeader, "Should not store state")
}

func (suite *SideHandlerTestSuite) TestHandleMsgCheckpointAck() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper
	start := uint64(0)
	maxSize := uint64(256)
	params := keeper.GetParams(ctx)

	header, _ := suite.GenRandCheckpointHeader(start, maxSize, params.MaxCheckpointLength)
	headerId := uint64(1)

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
	suite.contractCaller.On("GetHeaderInfo", headerId, rootchainInstance).Return(header.RootHash.EthHash(), header.StartBlock, header.EndBlock, header.TimeStamp, header.Proposer, nil)

	result := suite.sideHandler(ctx, msgCheckpointAck)
	require.Equal(t, uint32(sdk.CodeOK), result.Code, "Side tx handler should be success")
	require.Equal(t, abci.SideTxResultType_Yes, result.Result, "Result should be `yes`")
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
	cmn.LoadValidatorSet(2, t, stakingKeeper, ctx, false, 10)
	stakingKeeper.IncrementAccum(ctx, 1)

	lastCheckpoint, err := keeper.GetLastCheckpoint(ctx)
	if err == nil {
		start = start + lastCheckpoint.EndBlock + 1
	}

	header, err := suite.GenRandCheckpointHeader(start, maxSize, params.MaxCheckpointLength)

	// add current proposer to header
	header.Proposer = stakingKeeper.GetValidatorSet(ctx).Proposer.Signer

	// make sure proposer has min ether
	suite.contractCaller.On("GetBalance", stakingKeeper.GetValidatorSet(ctx).Proposer.Signer).Return(helper.MinBalance, nil)

	borChainId := "1234"
	// create checkpoint msg
	msgCheckpoint := types.NewMsgCheckpointBlock(
		header.Proposer,
		header.StartBlock,
		header.EndBlock,
		header.RootHash,
		header.RootHash,
		borChainId,
	)

	suite.contractCaller.On("CheckIfBlocksExist", header.EndBlock).Return(true)
	suite.contractCaller.On("GetRootHash", header.StartBlock, header.EndBlock, uint64(1024)).Return(header.RootHash, nil)
	result := suite.postHandler(ctx, msgCheckpoint, abci.SideTxResultType_Yes)

	require.True(t, result.IsOK(), "expected send-checkpoint to be ok, got %v", result)
	bufferedHeader, err := keeper.GetCheckpointFromBuffer(ctx)
	require.Equal(t, bufferedHeader.StartBlock, header.StartBlock)
	require.Equal(t, bufferedHeader.EndBlock, header.EndBlock)
	require.Equal(t, bufferedHeader.RootHash, header.RootHash)
	require.Equal(t, bufferedHeader.Proposer, header.Proposer)
	require.Equal(t, bufferedHeader.BorChainID, header.BorChainID)
	require.Empty(t, err, "Unable to set checkpoint from buffer, Error: %v", err)
}

func (suite *SideHandlerTestSuite) TestPostHandleMsgCheckpointAck() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.CheckpointKeeper

	start := uint64(0)
	maxSize := uint64(256)
	params := keeper.GetParams(ctx)
	header, _ := suite.GenRandCheckpointHeader(start, maxSize, params.MaxCheckpointLength)

	// add current proposer to header
	// header.Proposer = stakingKeeper.GetValidatorSet(ctx).Proposer.Signer

	// send ack
	headerId := uint64(1)

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
	rootchainInstance := &rootchain.Rootchain{}

	suite.contractCaller.On("GetRootChainInstance", mock.Anything).Return(rootchainInstance, nil)
	suite.contractCaller.On("GetHeaderInfo", headerId, rootchainInstance).Return(header.RootHash.EthHash(), header.StartBlock, header.EndBlock, header.TimeStamp, header.Proposer, nil)

	suite.postHandler(ctx, msgCheckpointAck, abci.SideTxResultType_Yes)
	afterAckBufferedHeader, _ := keeper.GetCheckpointFromBuffer(ctx)
	require.Nil(t, afterAckBufferedHeader)
}

//TODO: make it reusable
// GenRandCheckpointHeader create random header block
func (suite *SideHandlerTestSuite) GenRandCheckpointHeader(start uint64, headerSize uint64, maxCheckpointLenght uint64) (headerBlock hmTypes.CheckpointBlockHeader, err error) {
	app, ctx := suite.app, suite.ctx

	topupKeeper := app.TopupKeeper
	end := start + headerSize
	borChainID := "1234"

	dividendAccounts := topupKeeper.GetAllDividendAccounts(ctx)
	accRootHash, err := types.GetAccountRootHash(dividendAccounts)
	rootHash := hmTypes.BytesToHeimdallHash(accRootHash)
	proposer := hmTypes.HeimdallAddress{}
	headerBlock = hmTypes.CreateBlock(
		start,
		end,
		rootHash,
		proposer,
		borChainID,
		uint64(time.Now().UTC().Unix()))

	return headerBlock, nil
}
