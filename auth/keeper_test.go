package auth_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/auth/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

//
// Test suite
//

// KeeperTestSuite integrate test suite context object
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

//
// Tests
//

func (suite *KeeperTestSuite) TestAccountMapperGetSet() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	addr := hmTypes.BytesToHeimdallAddress([]byte("some-address"))

	// no account before its created
	acc := app.AccountKeeper.GetAccount(ctx, addr)
	require.Nil(t, acc)

	// create account and check default values
	acc = app.AccountKeeper.NewAccountWithAddress(ctx, addr)
	require.NotNil(t, acc)
	require.Equal(t, addr, acc.GetAddress())
	require.EqualValues(t, nil, acc.GetPubKey())
	require.EqualValues(t, 0, acc.GetSequence())

	// NewAccount doesn't call Set, so it's still nil
	require.Nil(t, app.AccountKeeper.GetAccount(ctx, addr))

	// set some values on the account and save it
	newSequence := uint64(20)
	err := acc.SetSequence(newSequence)
	require.NoError(t, err)
	app.AccountKeeper.SetAccount(ctx, acc)

	// check the new values
	acc = app.AccountKeeper.GetAccount(ctx, addr)
	require.NotNil(t, acc)
	require.Equal(t, newSequence, acc.GetSequence())
}

func (suite *KeeperTestSuite) TestAccountMapperRemoveAccount() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	addr1 := hmTypes.BytesToHeimdallAddress([]byte("addr1"))
	addr2 := hmTypes.BytesToHeimdallAddress([]byte("addr2"))

	// create accounts
	acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, addr1)
	acc2 := app.AccountKeeper.NewAccountWithAddress(ctx, addr2)

	accSeq1 := uint64(20)
	accSeq2 := uint64(40)

	err := acc1.SetSequence(accSeq1)
	require.NoError(t, err)
	err = acc2.SetSequence(accSeq2)
	require.NoError(t, err)
	app.AccountKeeper.SetAccount(ctx, acc1)
	app.AccountKeeper.SetAccount(ctx, acc2)

	acc1 = app.AccountKeeper.GetAccount(ctx, addr1)
	require.NotNil(t, acc1)
	require.Equal(t, accSeq1, acc1.GetSequence())

	// remove one account
	app.AccountKeeper.RemoveAccount(ctx, acc1)
	acc1 = app.AccountKeeper.GetAccount(ctx, addr1)
	require.Nil(t, acc1)

	acc2 = app.AccountKeeper.GetAccount(ctx, addr2)
	require.NotNil(t, acc2)
	require.Equal(t, accSeq2, acc2.GetSequence())
}

func (suite *KeeperTestSuite) TestGetSetParams() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	params := types.DefaultParams()

	app.AccountKeeper.SetParams(ctx, params)

	actualParams := app.AccountKeeper.GetParams(ctx)
	require.Equal(t, params, actualParams)
}
