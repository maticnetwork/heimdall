package milestone_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/helper/mocks"
	"github.com/maticnetwork/heimdall/milestone"
	"github.com/maticnetwork/heimdall/milestone/types"
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

// SetupTest setup all necessary things for querier tesing
func (suite *QuerierTestSuite) SetupTest() {
	suite.app, suite.ctx, suite.cliCtx = createTestApp(false)
	suite.contractCaller = mocks.IContractCaller{}
	suite.querier = milestone.NewQuerier(suite.app.MilestoneKeeper, suite.app.StakingKeeper, &suite.contractCaller)
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

	err := json.Unmarshal(res, &params)
	require.NoError(t, err)
	require.NotNil(t, params)
	require.Equal(t, defaultParams.SprintLength, params.SprintLength)
}

func (suite *QuerierTestSuite) TestQueryLatestMilestone() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	path := []string{types.QueryLatestMilestone}
	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLatestMilestone)

	startBlock := uint64(0)
	endBlock := uint64(255)
	rootHash := hmTypes.HexToHeimdallHash("123")
	proposerAddress := hmTypes.HexToHeimdallAddress("123")
	timestamp := uint64(time.Now().Unix())
	borChainId := "1234"
	milestoneID := "00000"

	milestoneBlock := hmTypes.CreateMilestone(
		startBlock,
		endBlock,
		rootHash,
		proposerAddress,
		borChainId,
		milestoneID,
		timestamp,
	)
	err := app.MilestoneKeeper.AddMilestone(ctx, milestoneBlock)
	require.NoError(t, err)

	req := abci.RequestQuery{
		Path: route,
		Data: []byte{},
	}

	res, err := querier(ctx, path, req)
	require.NoError(t, err)
	require.NotNil(t, res)

	var milestone hmTypes.Milestone
	err = json.Unmarshal(res, &milestone)
	require.NoError(t, err)
	require.Equal(t, milestone, milestoneBlock)
}

func (suite *QuerierTestSuite) TestQueryLastNoAckMilestone() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	path := []string{types.QueryLatestNoAckMilestone}
	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLatestNoAckMilestone)

	milestoneID := "00000"

	app.MilestoneKeeper.SetNoAckMilestone(ctx, milestoneID)

	req := abci.RequestQuery{
		Path: route,
		Data: []byte{},
	}

	res, err1 := querier(ctx, path, req)
	require.NoError(t, err1)
	require.NotNil(t, res)

	var _milestoneID string

	err2 := json.Unmarshal(res, &_milestoneID)

	require.NoError(t, err2)
	require.Equal(t, _milestoneID, milestoneID)

	milestoneID = "00001"
	app.MilestoneKeeper.SetNoAckMilestone(ctx, milestoneID)

	res, err1 = querier(ctx, path, req)
	require.NoError(t, err1)
	require.NotNil(t, res)

	err2 = json.Unmarshal(res, &_milestoneID)

	require.NoError(t, err2)
	require.Equal(t, _milestoneID, milestoneID)

}

func (suite *QuerierTestSuite) TestQueryNoAckMilestoneByID() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	path := []string{types.QueryNoAckMilestoneByID}
	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryNoAckMilestoneByID)

	milestoneID := "00000"

	req := abci.RequestQuery{
		Path: route,
		Data: app.Codec().MustMarshalJSON(types.NewQueryMilestoneID(milestoneID)),
	}

	res, err1 := querier(ctx, path, req)
	require.NoError(t, err1)
	require.NotNil(t, res)

	var val bool

	err2 := json.Unmarshal(res, &val)

	require.NoError(t, err2)
	require.Equal(t, val, false)

	app.MilestoneKeeper.SetNoAckMilestone(ctx, milestoneID)

	res, err1 = querier(ctx, path, req)
	require.NoError(t, err1)
	require.NotNil(t, res)

	err2 = json.Unmarshal(res, &val)

	require.NoError(t, err2)
	require.Equal(t, val, true)

	milestoneID = "00001"
	app.MilestoneKeeper.SetNoAckMilestone(ctx, milestoneID)

	res, err1 = querier(ctx, path, req)
	require.NoError(t, err1)
	require.NotNil(t, res)

	err2 = json.Unmarshal(res, &val)

	require.NoError(t, err2)
	require.Equal(t, val, true)

}

func (suite *QuerierTestSuite) TestSetNoAckMilestone() {
	t, app, ctx, _ := suite.T(), suite.app, suite.ctx, suite.querier

	milestoneID := "00000"

	startBlock := uint64(0)
	endBlock := uint64(255)
	rootHash := hmTypes.HexToHeimdallHash("123")
	proposerAddress := hmTypes.HexToHeimdallAddress("123")
	timestamp := uint64(time.Now().Unix())
	borChainId := "1234"

	milestoneBlock := hmTypes.CreateMilestone(
		startBlock,
		endBlock,
		rootHash,
		proposerAddress,
		borChainId,
		milestoneID,
		timestamp,
	)
	app.MilestoneKeeper.AddMilestone(ctx, milestoneBlock)

	milstoneBlockTemp, _ := app.MilestoneKeeper.GetLastMilestone(ctx)
	require.Equal(t, milstoneBlockTemp.MilestoneID, milestoneBlock.MilestoneID)

	app.MilestoneKeeper.SetNoAckMilestone(ctx, milestoneID)

	res := app.MilestoneKeeper.GetLastNoAckMilestone(ctx)

	require.Equal(t, res, "00000")

	var ct sdk.Context
	ct = app.BaseApp.NewContext(false, abci.Header{})

	//app.MilestoneKeeper.AddMilestone(ct, milestoneBlock)

	milstoneBlockTemp, _ = app.MilestoneKeeper.GetLastMilestone(ct)
	require.Equal(t, milstoneBlockTemp.MilestoneID, milestoneBlock.MilestoneID)

	res = app.MilestoneKeeper.GetLastNoAckMilestone(ct)
	require.Equal(t, res, "00000")

}
