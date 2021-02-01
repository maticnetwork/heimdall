package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/x/chainmanager/keeper"
	"github.com/maticnetwork/heimdall/x/chainmanager/types"
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestQuery() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	k := keeper.Querier{
		Keeper: app.ChainKeeper,
	}

	result, err := k.Params(sdk.WrapSDKContext(ctx), &types.QueryParamsRequest{})
	require.NotNil(t, result)
	require.Nil(t, err)
}
