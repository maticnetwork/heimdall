package checkpoint_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/maticnetwork/heimdall/checkpoint/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestMilestoneQuerierTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(QuerierTestSuite))
}

func (suite *QuerierTestSuite) TestQueryLatestMilestone() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	path_latest := []string{types.QueryLatestMilestone}
	path_byNumber := []string{types.QueryMilestoneByNumber}
	path_count := []string{types.QueryCount}
	route_latest := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLatestMilestone)
	route_byNumber := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryMilestoneByNumber)
	route_count := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryCount)
	startBlock := uint64(0)
	endBlock := uint64(255)
	hash := hmTypes.HexToHeimdallHash("123")
	proposerAddress := hmTypes.HexToHeimdallAddress("123")
	timestamp := uint64(time.Now().Unix())
	borChainId := "1234"
	milestoneID := "00000"
	milestoneBlock := hmTypes.CreateMilestone(
		startBlock,
		endBlock,
		hash,
		proposerAddress,
		borChainId,
		milestoneID,
		timestamp,
	)

	req_latest := abci.RequestQuery{
		Path: route_latest,
		Data: []byte{},
	}

	res, err := querier(ctx, path_latest, req_latest)

	require.Error(t, err)
	require.Nil(t, res)

	req_byNumber := abci.RequestQuery{
		Path: route_byNumber,
		Data: app.Codec().MustMarshalJSON(types.NewQueryMilestoneID(milestoneID)),
	}

	res, err = querier(ctx, path_byNumber, req_byNumber)

	require.Error(t, err)
	require.Nil(t, res)

	req_byNumber = abci.RequestQuery{
		Path: route_byNumber,
		Data: app.Codec().MustMarshalJSON(types.NewQueryMilestoneParams(1)),
	}

	res, err = querier(ctx, path_latest, req_byNumber)

	require.Error(t, err)
	require.Nil(t, res)

	errNew := app.CheckpointKeeper.AddMilestone(ctx, milestoneBlock)
	require.NoError(t, errNew)

	res, err = querier(ctx, path_latest, req_latest)

	require.NoError(t, err)
	require.NotNil(t, res)

	var milestone hmTypes.Milestone

	errNew = json.Unmarshal(res, &milestone)
	require.NoError(t, errNew)
	require.Equal(t, milestone, milestoneBlock)

	res, err = querier(ctx, path_latest, req_byNumber)

	require.NoError(t, err)
	require.NotNil(t, res)

	req_count := abci.RequestQuery{
		Path: route_count,
		Data: []byte{},
	}

	res, err = querier(ctx, path_count, req_count)

	require.NoError(t, err)
	require.NotNil(t, res)

	var count uint64
	errNew = json.Unmarshal(res, &count)

	require.NoError(t, err)
	require.Equal(t, count, uint64(1))
}
func (suite *QuerierTestSuite) TestQueryLastNoAckMilestone() {
	t, app, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier
	path := []string{types.QueryLatestNoAckMilestone}
	route := fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryLatestNoAckMilestone)
	milestoneID := "00000"

	app.CheckpointKeeper.SetNoAckMilestone(ctx, milestoneID)

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
	app.CheckpointKeeper.SetNoAckMilestone(ctx, milestoneID)
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

	app.CheckpointKeeper.SetNoAckMilestone(ctx, milestoneID)

	res, err1 = querier(ctx, path, req)
	require.NoError(t, err1)
	require.NotNil(t, res)

	err2 = json.Unmarshal(res, &val)
	require.NoError(t, err2)
	require.Equal(t, val, true)

	milestoneID = "00001"

	app.CheckpointKeeper.SetNoAckMilestone(ctx, milestoneID)
	app.CheckpointKeeper.SetNoAckMilestone(ctx, milestoneID)

	res, err1 = querier(ctx, path, req)
	require.NoError(t, err1)
	require.NotNil(t, res)

	err2 = json.Unmarshal(res, &val)
	require.NoError(t, err2)
	require.Equal(t, val, true)
}
