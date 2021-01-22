package keeper_test

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
	"github.com/maticnetwork/heimdall/x/checkpoint/keeper"
	chSim "github.com/maticnetwork/heimdall/x/checkpoint/simulation"
	"github.com/maticnetwork/heimdall/x/checkpoint/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// GrpcQueryTestSuite integrate test suite context object
type GrpcQueryTestSuite struct {
	suite.Suite

	app            *app.HeimdallApp
	ctx            sdk.Context
	cliCtx         client.Context
	grpcQuery      types.QueryServer
	contractCaller mocks.IContractCaller
}

// SetupTest setup all necessary things for grpc query testing
func (suite *GrpcQueryTestSuite) SetupTest() {
	suite.app, suite.ctx, suite.cliCtx = createTestApp(false)
	suite.contractCaller = mocks.IContractCaller{}
	suite.grpcQuery = keeper.NewQueryServerImpl(suite.app.CheckpointKeeper, &suite.contractCaller)
}

// TestGrpcQueryTestSuite
func TestGrpcQueryTestSuite(t *testing.T) {
	suite.Run(t, new(GrpcQueryTestSuite))
}

func (suite *GrpcQueryTestSuite) TestQueryParams() {
	t, _, ctx, grpcQuery := suite.T(), suite.app, suite.ctx, suite.grpcQuery
	defaultParams := types.DefaultParams()

	result, err := grpcQuery.Params(sdk.WrapSDKContext(ctx), &types.QueryParamsRequest{})
	require.NotNil(t, result)
	require.Nil(t, err)
	require.Equal(t, defaultParams.AvgCheckpointLength, result.Params.AvgCheckpointLength)
	require.Equal(t, defaultParams.MaxCheckpointLength, result.Params.MaxCheckpointLength)
	require.Equal(t, defaultParams.ChildBlockInterval, result.Params.ChildBlockInterval)
}

func (suite *GrpcQueryTestSuite) TestQueryAckCount() {
	t, _, ctx, grpcQuery := suite.T(), suite.app, suite.ctx, suite.grpcQuery

	ackCount := uint64(1)
	suite.app.CheckpointKeeper.UpdateACKCountWithValue(ctx, ackCount)

	result, err := grpcQuery.AckCount(sdk.WrapSDKContext(ctx), &types.QueryAckCountRequest{})
	require.NotNil(t, result)
	require.NoError(t, err)
	require.Equal(t, result.AckCount, ackCount)
}

func (suite *GrpcQueryTestSuite) TestQueryInvalidCheckpoint() {
	t, _, ctx, grpcQuery := suite.T(), suite.app, suite.ctx, suite.grpcQuery

	result, err := grpcQuery.Checkpoint(sdk.WrapSDKContext(ctx), &types.QueryCheckpointRequest{})
	require.Nil(t, result)
	require.NotNil(t, err)
	require.Error(t, err)
	require.Equal(t, err.Error(), "rpc error: code = InvalidArgument desc = empty header param")
}

func (suite *GrpcQueryTestSuite) TestQueryCheckpoint() {
	t, initApp, ctx, grpcQuery := suite.T(), suite.app, suite.ctx, suite.grpcQuery

	headerNumber := uint64(1)
	startBlock := uint64(0)
	endBlock := uint64(255)
	rootHash := hmCommonTypes.HexToHeimdallHash("123")
	proposerAddress := hmCommonTypes.HexToHeimdallAddress("123")
	timestamp := uint64(time.Now().Unix())
	borChainId := "1234"

	checkpointBlock := hmTypes.CreateBlock(
		startBlock,
		endBlock,
		rootHash,
		proposerAddress,
		borChainId,
		timestamp,
	)

	err := initApp.CheckpointKeeper.AddCheckpoint(ctx, headerNumber, checkpointBlock)
	require.NoError(t, err)

	result, err := grpcQuery.Checkpoint(sdk.WrapSDKContext(ctx), &types.QueryCheckpointRequest{
		Number: uint64(1),
	})

	require.NotNil(t, result)
	require.NoError(t, err)
	require.Equal(t, result.Checkpoint, checkpointBlock)
}

func (suite *GrpcQueryTestSuite) TestQueryCheckpointBuffer() {
	t, initApp, ctx, grpcQuery := suite.T(), suite.app, suite.ctx, suite.grpcQuery

	startBlock := uint64(0)
	endBlock := uint64(255)
	rootHash := hmCommonTypes.HexToHeimdallHash("123")
	proposerAddress := hmCommonTypes.HexToHeimdallAddress("123")
	timestamp := uint64(time.Now().Unix())
	borChainId := "1234"

	checkpointBlock := hmTypes.CreateBlock(
		startBlock,
		endBlock,
		rootHash,
		proposerAddress,
		borChainId,
		timestamp,
	)
	err := initApp.CheckpointKeeper.SetCheckpointBuffer(ctx, checkpointBlock)
	require.NoError(t, err)

	result, err := grpcQuery.CheckpointBuffer(sdk.WrapSDKContext(ctx), &types.QueryCheckpointBufferRequest{})
	require.NotNil(t, result)
	require.NoError(t, err)

	require.Equal(t, result.CheckpointBuffer, checkpointBlock)
}

func (suite *GrpcQueryTestSuite) TestQueryLastNoAck() {
	t, initApp, ctx, grpcQuery := suite.T(), suite.app, suite.ctx, suite.grpcQuery

	noAck := uint64(time.Now().Unix())
	initApp.CheckpointKeeper.SetLastNoAck(ctx, noAck)

	result, err := grpcQuery.LastNoAck(sdk.WrapSDKContext(ctx), &types.QueryLastNoAckRequest{})
	require.NotNil(t, result)
	require.NoError(t, err)
	require.Equal(t, result.LastNoAck, noAck)
}

func (suite *GrpcQueryTestSuite) TestQueryCheckpointList() {
	t, initApp, ctx, grpcQuery := suite.T(), suite.app, suite.ctx, suite.grpcQuery

	checkPointKeeper := initApp.CheckpointKeeper

	count := 5
	startBlock := uint64(0)
	endBlock := uint64(0)
	checkpoints := make([]*hmTypes.Checkpoint, count)

	for i := 0; i < count; i++ {
		headerBlockNumber := uint64(i) + 1

		startBlock = startBlock + endBlock
		endBlock = endBlock + uint64(255)
		rootHash := hmCommonTypes.HexToHeimdallHash("123")
		proposerAddress := hmCommonTypes.HexToHeimdallAddress("123")
		timestamp := uint64(time.Now().Unix()) + uint64(i)
		borChainId := "1234"

		checkpoint := hmTypes.CreateBlock(
			startBlock,
			endBlock,
			rootHash,
			proposerAddress,
			borChainId,
			timestamp,
		)
		checkpoints[i] = checkpoint
		err := checkPointKeeper.AddCheckpoint(ctx, headerBlockNumber, checkpoint)
		require.NoError(t, err)
		checkPointKeeper.UpdateACKCount(ctx)
	}

	result, err := grpcQuery.CheckpointList(sdk.WrapSDKContext(ctx), &types.QueryCheckpointListRequest{
		Pagination: hmTypes.NewQueryPaginationParams(uint64(1), uint64(10)),
	})

	require.NotNil(t, result)
	require.NoError(t, err)
	require.Equal(t, checkpoints, result.CheckpointList)
}

func (suite *GrpcQueryTestSuite) TestQueryNextCheckpoint() {
	t, initApp, ctx, grpcQuery := suite.T(), suite.app, suite.ctx, suite.grpcQuery
	chSim.LoadValidatorSet(2, t, initApp.StakingKeeper, ctx, false, 10)

	dividendAccount := hmTypes.DividendAccount{
		User:      hmCommonTypes.HexToHeimdallAddress("123").String(),
		FeeAmount: big.NewInt(0).String(),
	}
	err := initApp.TopupKeeper.AddDividendAccount(ctx, dividendAccount)
	require.NoError(t, err)

	headerNumber := uint64(1)
	startBlock := uint64(0)
	endBlock := uint64(256)
	rootHash := hmCommonTypes.HexToHeimdallHash("123")
	proposerAddress := hmCommonTypes.HexToHeimdallAddress("123")
	timestamp := uint64(time.Now().Unix())
	borChainId := "1234"

	checkpointBlock := hmTypes.CreateBlock(
		startBlock,
		endBlock,
		rootHash,
		proposerAddress,
		borChainId,
		timestamp,
	)

	suite.contractCaller.On("GetRootHash", checkpointBlock.StartBlock, checkpointBlock.EndBlock, uint64(1024)).Return(hmCommonTypes.HexToHeimdallHash(checkpointBlock.RootHash).Bytes(), nil)
	err = initApp.CheckpointKeeper.AddCheckpoint(ctx, headerNumber, checkpointBlock)
	require.NoError(t, err)

	result, err := grpcQuery.NextCheckpoint(sdk.WrapSDKContext(ctx), &types.QueryNextCheckpointRequest{
		BorChainID: borChainId,
	})

	require.NotNil(t, result)
	require.NoError(t, err)
	require.Equal(t, checkpointBlock.StartBlock, result.NextCheckpoint.StartBlock)
	require.Equal(t, checkpointBlock.EndBlock, result.NextCheckpoint.EndBlock)
	require.Equal(t, checkpointBlock.RootHash, result.NextCheckpoint.RootHash)
	require.Equal(t, checkpointBlock.BorChainID, result.NextCheckpoint.BorChainID)
}
