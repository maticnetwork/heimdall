package chainmanager_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	jsoniter "github.com/json-iterator/go"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/chainmanager"
	"github.com/maticnetwork/heimdall/chainmanager/types"
)

// QuerierTestSuite integrate test suite context object
type QuerierTestSuite struct {
	suite.Suite

	app     *app.HeimdallApp
	ctx     sdk.Context
	querier sdk.Querier
}

// SetupTest setup all necessary things for querier testing
func (suite *QuerierTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(false)
	suite.querier = chainmanager.NewQuerier(suite.app.ChainKeeper)
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

// TestQueryParams queries params
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

	// match response params
	require.Equal(t, defaultParams.MainchainTxConfirmations, params.MainchainTxConfirmations)
	require.Equal(t, defaultParams.MaticchainTxConfirmations, params.MaticchainTxConfirmations)
	require.Equal(t, defaultParams.ChainParams, params.ChainParams)

	{
		rapp := app.Setup(true)
		ctx := rapp.BaseApp.NewContext(true, abci.Header{})
		querier := chainmanager.NewQuerier(rapp.ChainKeeper)
		require.Panics(t, func() {
			_, err = querier(ctx, path, req)
			require.NoError(t, err)
		})
	}
}
