package bank_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
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

func (suite *KeeperTestSuite) TestSetCoins() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.BankKeeper
	amount := int64(10000000)
	address := hmTypes.HexToHeimdallAddress("123")

	coins := sdk.NewCoins(sdk.NewCoin(authTypes.FeeToken, sdk.NewInt(amount*10)))
	// keeper.SendCoins(ctx, hmTypes.HexToHeimdallAddress("123"), hmTypes.HexToHeimdallAddress("456"), coins)
	keeper.SetCoins(ctx, address, coins)

	res := keeper.GetCoins(ctx, address)
	require.Equal(t, res, coins)
}

func (suite *KeeperTestSuite) TestAddCoins() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.BankKeeper
	amount := int64(10000000)
	address := hmTypes.HexToHeimdallAddress("123")

	coins := sdk.NewCoins(sdk.NewCoin(authTypes.FeeToken, sdk.NewInt(amount*10)))
	// keeper.SendCoins(ctx, hmTypes.HexToHeimdallAddress("123"), hmTypes.HexToHeimdallAddress("456"), coins)
	result, err := keeper.AddCoins(ctx, address, coins)
	require.NotNil(t, result)
	require.NoError(t, err)

	res := keeper.GetCoins(ctx, address)
	require.Equal(t, res, coins)
}

func (suite *KeeperTestSuite) TestSendCoins() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.BankKeeper
	amount := int64(10000000)
	address := hmTypes.HexToHeimdallAddress("123")
	to := hmTypes.HexToHeimdallAddress("456")
	coins := sdk.NewCoins(sdk.NewCoin(authTypes.FeeToken, sdk.NewInt(amount)))
	keeper.SetCoins(ctx, address, coins)

	err := keeper.SendCoins(ctx, address, to, coins)
	require.NoError(t, err)

	fromAcc := keeper.GetCoins(ctx, address)
	require.Equal(t, sdk.NewInt(0), fromAcc.AmountOf(authTypes.FeeToken))

	toAcc := app.BankKeeper.GetCoins(ctx, to)
	require.LessOrEqual(t, sdk.NewInt(amount).Int64(), toAcc.AmountOf(authTypes.FeeToken).Int64())
}

func (suite *KeeperTestSuite) TestHasCoins() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.BankKeeper
	amount := int64(10000000)
	address := hmTypes.HexToHeimdallAddress("123")

	coins := sdk.NewCoins(sdk.NewCoin(authTypes.FeeToken, sdk.NewInt(amount*10)))
	keeper.SetCoins(ctx, address, coins)

	res := app.BankKeeper.HasCoins(ctx, address, coins)
	require.True(t, res)
}

func (suite *KeeperTestSuite) TestGetSendEnabled() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	keeper := app.BankKeeper
	res := keeper.GetSendEnabled(ctx)
	require.True(t, res)
}
