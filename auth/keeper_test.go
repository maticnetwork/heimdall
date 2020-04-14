package auth_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/auth/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/simulation"
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
	require.NoError(t, acc.SetSequence(newSequence))
	app.AccountKeeper.SetAccount(ctx, acc)

	// check the new values
	acc = app.AccountKeeper.GetAccount(ctx, addr)
	require.NotNil(t, acc)
	require.Equal(t, newSequence, acc.GetSequence())

	// set coins values on the account and save it
	coins := simulation.RandomFeeCoins()
	require.NoError(t, acc.SetCoins(coins)) // set and check error
	app.AccountKeeper.SetAccount(ctx, acc)

	// check the new values
	acc = app.AccountKeeper.GetAccount(ctx, addr)
	require.NotNil(t, acc)
	require.True(t, coins.IsEqual(acc.GetCoins()))
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

func (suite *KeeperTestSuite) TestLogger() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	logger := app.AccountKeeper.Logger(ctx)
	require.NotNil(t, logger)
}

func (suite *KeeperTestSuite) TestGetSequence() {
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

	sequence, err := app.AccountKeeper.GetSequence(ctx, addr1)
	require.NoError(t, err)
	require.Equal(t, accSeq1, sequence)

	sequence, err = app.AccountKeeper.GetSequence(ctx, addr2)
	require.NoError(t, err)
	require.Equal(t, accSeq2, sequence)

	_, err = app.AccountKeeper.GetSequence(ctx, hmTypes.BytesToHeimdallAddress([]byte("addr3")))
	require.Error(t, err)
}

func (suite *KeeperTestSuite) TestGetPubKey() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	_, pubkey, addr := sdkAuth.KeyTestPubAddr()

	addr1 := hmTypes.AccAddressToHeimdallAddress(addr)
	acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, addr1)
	acc1.SetPubKey(pubkey)
	app.AccountKeeper.SetAccount(ctx, acc1)
	pubkey1, err := app.AccountKeeper.GetPubKey(ctx, addr1)
	require.Nil(t, err)
	require.NotEqual(t, 0, len(pubkey.Bytes()))
	require.Equal(t, pubkey, pubkey1)

	_, err = app.AccountKeeper.GetPubKey(ctx, hmTypes.BytesToHeimdallAddress([]byte("addr3")))
	require.Error(t, err)
}

func (suite *KeeperTestSuite) TestGetAllAccounts() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	accounts := app.AccountKeeper.GetAllAccounts(ctx) // module accounts
	require.True(t, len(accounts) > 0)
}

func (suite *KeeperTestSuite) TestIterateAccounts() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	newAccounts := 10
	beforeAccounts := app.AccountKeeper.GetAllAccounts(ctx) // current accounts
	for i := 0; i < newAccounts; i++ {
		addr := hmTypes.BytesToHeimdallAddress([]byte(fmt.Sprintf("address-%v", i)))
		acc := app.AccountKeeper.NewAccountWithAddress(ctx, addr)
		acc.SetCoins(simulation.RandomFeeCoins())
		app.AccountKeeper.SetAccount(ctx, acc)
	}
	afterAccounts := app.AccountKeeper.GetAllAccounts(ctx) // current accounts
	require.Equal(t, newAccounts, len(afterAccounts)-len(beforeAccounts))

	var filteredAccounts []types.Account
	app.AccountKeeper.IterateAccounts(ctx, func(acc types.Account) bool {
		filteredAccounts = append(filteredAccounts, acc)

		if acc.GetAccountNumber() > 5 {
			return true
		}
		return false
	})
	require.Equal(t, 5, len(filteredAccounts))
}

func (suite *KeeperTestSuite) TestProposer() {
	t, happ, ctx := suite.T(), suite.app, suite.ctx

	proposer, ok := happ.AccountKeeper.GetBlockProposer(ctx)
	require.False(t, ok)
	require.Equal(t, hmTypes.ZeroHeimdallAddress.Bytes(), proposer.Bytes())

	// set addr as proposer
	addr := hmTypes.BytesToHeimdallAddress([]byte("some-address"))
	happ.AccountKeeper.SetBlockProposer(ctx, addr)

	proposer, ok = happ.AccountKeeper.GetBlockProposer(ctx)
	require.True(t, ok)
	require.Equal(t, addr.Bytes(), proposer.Bytes())

	// remove block proposer
	happ.AccountKeeper.RemoveBlockProposer(ctx)

	proposer, ok = happ.AccountKeeper.GetBlockProposer(ctx)
	require.False(t, ok)
	require.Equal(t, hmTypes.ZeroHeimdallAddress.Bytes(), proposer.Bytes())
}
