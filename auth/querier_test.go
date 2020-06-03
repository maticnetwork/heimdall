package auth_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/auth"
	"github.com/maticnetwork/heimdall/auth/exported"
	"github.com/maticnetwork/heimdall/auth/types"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
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
	suite.querier = auth.NewQuerier(suite.app.AccountKeeper)
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

func (suite *QuerierTestSuite) TestQueryAccount() {
	t, happ, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier
	cdc := happ.Codec()

	// account path
	path := []string{types.QueryAccount}

	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryAccount),
		Data: []byte{},
	}
	res, err := querier(ctx, path, req)
	require.Error(t, err)
	require.Nil(t, res)

	req.Data = cdc.MustMarshalJSON(types.NewQueryAccountParams(hmTypes.BytesToHeimdallAddress([]byte(""))))
	res, err = querier(ctx, path, req)
	require.Error(t, err)
	require.Nil(t, res)

	_, _, addr := sdkAuth.KeyTestPubAddr()
	req.Data = cdc.MustMarshalJSON(types.NewQueryAccountParams(hmTypes.AccAddressToHeimdallAddress(addr)))
	res, err = querier(ctx, path, req)
	require.Error(t, err)
	require.Nil(t, res)

	happ.AccountKeeper.SetAccount(ctx, happ.AccountKeeper.NewAccountWithAddress(ctx, hmTypes.AccAddressToHeimdallAddress(addr)))
	res, err = querier(ctx, path, req)
	require.NoError(t, err)
	require.NotNil(t, res)

	res, err = querier(ctx, path, req)
	require.NoError(t, err)
	require.NotNil(t, res)

	var account exported.Account
	err2 := cdc.UnmarshalJSON(res, &account)
	require.Nil(t, err2)
	require.Equal(t, account.GetAddress().Bytes(), addr.Bytes())

	{
		// setting tnil to account
		require.Panics(t, func() {
			happ.AccountKeeper.SetAccount(ctx, nil)
		})

		// store invalid/empty account
		store := ctx.KVStore(happ.GetKey(authTypes.StoreKey))
		store.Set(types.AddressStoreKey(hmTypes.AccAddressToHeimdallAddress(addr)), []byte(""))
		require.Panics(t, func() {
			querier(ctx, path, req)
		})
	}
}

func (suite *QuerierTestSuite) TestQueryParams() {
	t, happ, ctx, querier := suite.T(), suite.app, suite.ctx, suite.querier

	path := []string{types.QueryParams}
	req := abci.RequestQuery{
		Path: fmt.Sprintf("custom/%s/%s", types.QuerierRoute, types.QueryParams),
		Data: []byte{},
	}
	res, err := querier(ctx, path, req)
	require.NoError(t, err)
	require.NotNil(t, res)

	// default params
	defaultParams := authTypes.DefaultParams()

	var params types.Params
	err2 := json.Unmarshal(res, &params)
	require.Nil(t, err2)
	require.Equal(t, defaultParams.MaxMemoCharacters, params.MaxMemoCharacters)
	require.Equal(t, defaultParams.TxSigLimit, params.TxSigLimit)
	require.Equal(t, defaultParams.TxSizeCostPerByte, params.TxSizeCostPerByte)
	require.Equal(t, defaultParams.SigVerifyCostED25519, params.SigVerifyCostED25519)
	require.Equal(t, defaultParams.SigVerifyCostSecp256k1, params.SigVerifyCostSecp256k1)
	require.Equal(t, defaultParams.MaxTxGas, params.MaxTxGas)
	require.Equal(t, defaultParams.TxFees, params.TxFees)

	// set max characters
	params.MaxMemoCharacters = 10
	params.TxSizeCostPerByte = 8
	happ.AccountKeeper.SetParams(ctx, params)
	res, err = querier(ctx, path, req)
	require.NoError(t, err)
	require.NotEmpty(t, string(res))

	var params3 types.Params
	err3 := json.Unmarshal(res, &params3)
	require.NoError(t, err3)
	require.Equal(t, uint64(10), params.MaxMemoCharacters)
	require.Equal(t, uint64(8), params.TxSizeCostPerByte)

	{
		happ := app.Setup(true)
		ctx := happ.BaseApp.NewContext(true, abci.Header{})
		querier := auth.NewQuerier(happ.AccountKeeper)
		require.Panics(t, func() {
			querier(ctx, path, req)
		})
	}
}

func (suite *QuerierTestSuite) TestSimulation() {
	t, app, _, _ := suite.T(), suite.app, suite.ctx, suite.querier

	txBytesStr := "57f0625dee0a4b8ac198bc080112144293d4f71cf3f6908a5002562ccd77f4a953354d18800220ff332a0531353030313220f34a0e6b1df0c95720b84b7f540b0259e083e1667f3538dc02ded5e4042be1d9120410c09a0c"

	// simulate by calling Query with encoded tx

	txBytes, err := hex.DecodeString(txBytesStr)
	if err != nil {
		t.Log("Error tx str", err)
		return
	}

	// decoder := authTypes.DefaultTxDecoder(app.Codec())
	// tx, err := decoder(txBytes)
	// if err != nil {
	// 	t.Log("Error Decoding tx", err)
	// 	return
	// }

	// t.Log("Decoded tx", tx)

	query := abci.RequestQuery{
		Path: "/app/simulate",
		Data: txBytes,
	}

	queryResult := app.Query(query)
	require.True(t, queryResult.IsOK(), queryResult.Log)
	// t.Log("queryResult", queryResult)
	// t.Log("queryResult", queryResult.GetValue())

	var result sdk.Result
	err = codec.Cdc.UnmarshalBinaryLengthPrefixed(queryResult.Value, &result)

	estimate, err := parseQueryResponse(codec.Cdc, queryResult.Value)
	t.Log("parsedResponse", "Estimate", estimate, "error", err)
	if err != nil {
		return
	}
	// t.Log("parsedResponse", "Value", queryResult.Value, "error", err, "GasUsed", result.GasUsed)
}

func parseQueryResponse(cdc *codec.Codec, rawRes []byte) (uint64, error) {
	fmt.Println("Venky - parseQueryResponse - rawRes", rawRes)
	var simulationResult sdk.Result
	if err := cdc.UnmarshalBinaryLengthPrefixed(rawRes, &simulationResult); err != nil {
		fmt.Println("Venky - parseQueryResponse - error", err)
		return 0, err
	}

	return simulationResult.GasUsed, nil
}
