package bank_test

import (
	"encoding/json"
	"fmt"
	"testing"

	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/bank"
	"github.com/maticnetwork/heimdall/bank/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/simulation"
)

//
// Test suite
//

// QuerierTestSuite integrate test suite context object
type QuerierTestSuite struct {
	suite.Suite

	app     *app.HeimdallApp
	ctx     sdk.Context
	querier sdk.Querier
}

func (suite *QuerierTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(false)
	suite.querier = bank.NewQuerier(suite.app.BankKeeper)
}

func TestQuerierTestSuite(t *testing.T) {
	suite.Run(t, new(QuerierTestSuite))
}

//
// Tests
//

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

func (suite *QuerierTestSuite) TestQueryBalance() {
	t, happ, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier
	cdc := happ.Codec()

	// account path
	path := []string{types.QueryBalance}

	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryBalance),
		Data: []byte{},
	}
	res, err := querier(ctx, path, req)
	require.Error(t, err)
	require.Nil(t, res)

	// balance for non-existing address (empty address)

	req.Data = cdc.MustMarshalJSON(types.NewQueryBalanceParams(hmTypes.BytesToHeimdallAddress([]byte(""))))
	res, err = querier(ctx, path, req)
	require.NoError(t, err)
	require.NotNil(t, res)

	// fetch balance
	var balance sdk.Coins
	require.NoError(t, json.Unmarshal(res, &balance))
	require.True(t, balance.IsZero())

	// balance for non-existing address
	_, _, addr := sdkAuth.KeyTestPubAddr()
	req.Data = cdc.MustMarshalJSON(types.NewQueryBalanceParams(hmTypes.AccAddressToHeimdallAddress(addr)))
	res, err = querier(ctx, path, req)
	require.NoError(t, err)
	require.NotNil(t, res)

	require.NoError(t, json.Unmarshal(res, &balance))
	require.True(t, balance.IsZero())

	// set account
	acc1 := happ.AccountKeeper.NewAccountWithAddress(ctx, hmTypes.AccAddressToHeimdallAddress(addr))
	amt := simulation.RandomFeeCoins()
	acc1.SetCoins(amt)
	happ.AccountKeeper.SetAccount(ctx, acc1)

	res, err = querier(ctx, path, req)
	require.NoError(t, err)
	require.NotNil(t, res)

	require.NoError(t, json.Unmarshal(res, &balance))
	require.True(t, balance.IsEqual(amt), "address coins stored in the store should be equal to amt")

	{
		// setting nil to account
		require.Panics(t, func() {
			happ.AccountKeeper.SetAccount(ctx, nil)
		})

		// store invalid/empty account
		store := ctx.KVStore(happ.GetKey(authTypes.StoreKey))
		store.Set(authTypes.AddressStoreKey(hmTypes.AccAddressToHeimdallAddress(addr)), []byte(""))
		require.Panics(t, func() {
			querier(ctx, path, req)
		})
	}
}
