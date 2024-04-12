package checkpoint_test

import (
	"fmt"
	"math/big"
	"strconv"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/checkpoint"
	chSim "github.com/maticnetwork/heimdall/checkpoint/simulation"
	"github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/helper/mocks"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// QuerierTestSuite integrate test suite context object
type QuerierTestSuite struct {
	suite.Suite

	app            *app.HeimdallApp
	ctx            sdk.Context
	cliCtx         context.CLIContext
	querier        sdk.Querier
	contractCaller mocks.IContractCaller
}

// SetupTest setup all necessary things for querier testing
func (suite *QuerierTestSuite) SetupTest() {
	suite.app, suite.ctx, suite.cliCtx = createTestApp(false)
	suite.contractCaller = mocks.IContractCaller{}
	suite.querier = checkpoint.NewQuerier(suite.app.CheckpointKeeper, suite.app.StakingKeeper, suite.app.TopupKeeper, &suite.contractCaller)
}

// TestQuerierTestSuite
func TestQuerierTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(QuerierTestSuite))
}

// TestInvalidQuery checks request query
func (suite *QuerierTestSuite) TestInvalidQuery() {
	t, _, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	req := abci.RequestQuery{
		Path: "",
		Data: []byte{},
	}

	bz, err := querier(ctx, []string{"other"}, req)
	require.Error(t, err)
	require.Nil(t, bz)

	bz, err = querier(ctx, []string{types.QuerierRoute}, req)
	require.Error(t, err)
	require.Nil(t, bz)
}

func (suite *QuerierTestSuite) TestQueryParams() {
	t, _, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	var params types.Params

	defaultParams := types.DefaultParams()

	path := []string{types.QueryParams}

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParams)
	req := abci.RequestQuery{
		Path: route,
		Data: []byte{},
	}
	res, sdkErr := querier(ctx, path, req)
	require.NoError(t, sdkErr)
	require.NotNil(t, res)

	err := jsoniter.ConfigFastest.Unmarshal(res, &params)
	require.NoError(t, err)
	require.NotNil(t, params)
	require.Equal(t, defaultParams.AvgCheckpointLength, params.AvgCheckpointLength)
	require.Equal(t, defaultParams.MaxCheckpointLength, params.MaxCheckpointLength)
}

func (suite *QuerierTestSuite) TestQueryAckCount() {
	t, _, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier
	path := []string{types.QueryAckCount}

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAckCount)
	req := abci.RequestQuery{
		Path: route,
		Data: []byte{},
	}

	ackCount := uint64(1)
	suite.app.CheckpointKeeper.UpdateACKCountWithValue(ctx, ackCount)

	res, err := querier(ctx, path, req)
	require.NoError(t, err)
	require.NotNil(t, res)

	actualAckcount, _ := strconv.ParseUint(string(res), 0, 64)
	require.Equal(t, actualAckcount, ackCount)
}

func (suite *QuerierTestSuite) TestQueryCheckpoint() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	headerNumber := uint64(1)
	startBlock := uint64(0)
	endBlock := uint64(255)
	rootHash := hmTypes.HexToHeimdallHash("123")
	proposerAddress := hmTypes.HexToHeimdallAddress("123")
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
	err := app.CheckpointKeeper.AddCheckpoint(ctx, headerNumber, checkpointBlock)
	require.NoError(t, err)

	path := []string{types.QueryCheckpoint}
	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCheckpoint)

	req := abci.RequestQuery{
		Path: route,
		Data: app.Codec().MustMarshalJSON(types.NewQueryCheckpointParams(headerNumber)),
	}

	res, sdkErr := querier(ctx, path, req)
	require.NoError(t, sdkErr)
	require.NotNil(t, res)

	var checkpoint hmTypes.Checkpoint

	err = jsoniter.ConfigFastest.Unmarshal(res, &checkpoint)
	require.NoError(t, err)
	require.Equal(t, checkpoint, checkpointBlock)
}

func (suite *QuerierTestSuite) TestQueryCheckpointBuffer() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	path := []string{types.QueryCheckpointBuffer}
	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCheckpointBuffer)

	startBlock := uint64(0)
	endBlock := uint64(255)
	rootHash := hmTypes.HexToHeimdallHash("123")
	proposerAddress := hmTypes.HexToHeimdallAddress("123")
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
	err := app.CheckpointKeeper.SetCheckpointBuffer(ctx, checkpointBlock)
	require.NoError(t, err)

	req := abci.RequestQuery{
		Path: route,
		Data: []byte{},
	}

	res, err := querier(ctx, path, req)
	require.NoError(t, err)
	require.NotNil(t, res)

	var cp hmTypes.Checkpoint
	err = jsoniter.ConfigFastest.Unmarshal(res, &cp)
	require.NoError(t, err)
	require.Equal(t, cp, checkpointBlock)
}

func (suite *QuerierTestSuite) TestQueryLastNoAck() {
	t, _, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	path := []string{types.QueryLastNoAck}

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLastNoAck)
	req := abci.RequestQuery{
		Path: route,
		Data: []byte{},
	}

	noAck := uint64(time.Now().Unix())
	suite.app.CheckpointKeeper.SetLastNoAck(ctx, noAck)

	res, err := querier(ctx, path, req)
	require.NoError(t, err)
	require.NotNil(t, res)

	actualRes, _ := strconv.ParseUint(string(res), 10, 64)
	require.Equal(t, actualRes, noAck)
}

func (suite *QuerierTestSuite) TestQueryCheckpointList() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	keeper := app.CheckpointKeeper

	count := 5

	startBlock := uint64(0)
	endBlock := uint64(0)
	checkpoints := make([]hmTypes.Checkpoint, count)

	for i := 0; i < count; i++ {
		headerBlockNumber := uint64(i) + 1

		startBlock = startBlock + endBlock
		endBlock = endBlock + uint64(255)
		rootHash := hmTypes.HexToHeimdallHash("123")
		proposerAddress := hmTypes.HexToHeimdallAddress("123")
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
		err := keeper.AddCheckpoint(ctx, headerBlockNumber, checkpoint)
		require.NoError(t, err)
		keeper.UpdateACKCount(ctx)
	}

	path := []string{types.QueryCheckpointList}

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCheckpointList)
	req := abci.RequestQuery{
		Path: route,
		Data: app.Codec().MustMarshalJSON(hmTypes.NewQueryPaginationParams(uint64(1), uint64(10))),
	}
	res, sdkErr := querier(ctx, path, req)
	require.NoError(t, sdkErr)
	require.NotNil(t, res)

	var actualRes []hmTypes.Checkpoint

	err := jsoniter.ConfigFastest.Unmarshal(res, &actualRes)
	require.NoError(t, err)
	require.Equal(t, checkpoints, actualRes)
}

func (suite *QuerierTestSuite) TestQueryNextCheckpoint() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier
	chSim.LoadValidatorSet(t, 2, app.StakingKeeper, ctx, false, 10)

	dividendAccount := hmTypes.DividendAccount{
		User:      hmTypes.HexToHeimdallAddress("123"),
		FeeAmount: big.NewInt(0).String(),
	}
	err := app.TopupKeeper.AddDividendAccount(ctx, dividendAccount)
	require.NoError(t, err)

	headerNumber := uint64(1)
	startBlock := uint64(0)
	endBlock := uint64(256)
	rootHash := hmTypes.HexToHeimdallHash("123")
	proposerAddress := hmTypes.HexToHeimdallAddress("123")
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

	suite.contractCaller.On("GetRootHash", checkpointBlock.StartBlock, checkpointBlock.EndBlock, uint64(1024)).Return(checkpointBlock.RootHash.Bytes(), nil)
	err = app.CheckpointKeeper.AddCheckpoint(ctx, headerNumber, checkpointBlock)
	require.NoError(t, err)

	path := []string{types.QueryNextCheckpoint}

	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryNextCheckpoint)
	req := abci.RequestQuery{
		Path: route,
		Data: app.Codec().MustMarshalJSON(types.NewQueryBorChainID(borChainId)),
	}
	res, sdkErr := querier(ctx, path, req)
	require.NoError(t, sdkErr)
	require.NotNil(t, res)

	var actualRes types.MsgCheckpoint

	err = jsoniter.ConfigFastest.Unmarshal(res, &actualRes)
	require.NoError(t, err)
	require.Equal(t, checkpointBlock.StartBlock, actualRes.StartBlock)
	require.Equal(t, checkpointBlock.EndBlock, actualRes.EndBlock)
	require.Equal(t, checkpointBlock.RootHash, actualRes.RootHash)
	require.Equal(t, checkpointBlock.BorChainID, actualRes.BorChainID)
}
