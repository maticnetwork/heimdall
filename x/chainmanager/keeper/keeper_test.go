package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/x/chainmanager/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(false)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

// Tests

func (suite *KeeperTestSuite) TestParamsGetterSetter() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	params := types.DefaultParams()

	initApp.ChainKeeper.SetParams(ctx, params)

	actualParams := initApp.ChainKeeper.GetParams(ctx)
	require.Equal(t, params, actualParams)
}
