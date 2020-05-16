package checkpoint_test

import (
	"math/rand"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/checkpoint/types"
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
