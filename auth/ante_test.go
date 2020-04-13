package auth_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/auth"
	"github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/simulation"
)

//
// Test suite
//

// AnteTestSuite integrate test suite context object
type AnteTestSuite struct {
	suite.Suite

	app         *app.HeimdallApp
	ctx         sdk.Context
	anteHandler sdk.AnteHandler
}

func (suite *AnteTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(false)

	caller, err := helper.NewContractCaller()
	require.NoError(suite.T(), err)

	suite.anteHandler = auth.NewAnteHandler(
		suite.app.AccountKeeper,
		suite.app.ChainKeeper,
		suite.app.SupplyKeeper,
		&caller,
		auth.DefaultSigVerificationGasConsumer,
	)
}

func TestAnteTestSuite(t *testing.T) {
	suite.Run(t, new(AnteTestSuite))
}

//
// Tests
//

func (suite *AnteTestSuite) TestAccountNumbers() {
	t, happ, ctx, anteHandler := suite.T(), suite.app, suite.ctx, suite.anteHandler

	// keys and addresses
	priv1, _, addr1 := sdkAuth.KeyTestPubAddr()
	_, _, addr2 := sdkAuth.KeyTestPubAddr()

	// set the accounts
	acc1 := happ.AccountKeeper.NewAccountWithAddress(ctx, hmTypes.AccAddressToHeimdallAddress(addr1))
	acc1.SetCoins(simulation.RandomFeeCoins())
	require.NoError(t, acc1.SetAccountNumber(0))
	happ.AccountKeeper.SetAccount(ctx, acc1)
	acc2 := happ.AccountKeeper.NewAccountWithAddress(ctx, hmTypes.AccAddressToHeimdallAddress(addr2))
	acc2.SetCoins(simulation.RandomFeeCoins())
	require.NoError(t, acc2.SetAccountNumber(1))
	happ.AccountKeeper.SetAccount(ctx, acc2)

	// msg and signatures
	var tx sdk.Tx
	msg := sdkAuth.NewTestMsg(addr1)

	// test good tx from one signer
	tx = types.NewTestTx(ctx, msg, priv1, uint64(0), uint64(0))
	checkValidTx(t, anteHandler, ctx, tx, false)

	// new tx from wrong account number
	tx = types.NewTestTx(ctx, msg, priv1, uint64(1), uint64(0))
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeUnauthorized)

	// from correct account number
	tx = types.NewTestTx(ctx, msg, priv1, uint64(0), uint64(1))
	checkValidTx(t, anteHandler, ctx, tx, false)
}

//
// utils
//

// run the tx through the anteHandler and ensure its valid
func checkValidTx(t *testing.T, anteHandler sdk.AnteHandler, ctx sdk.Context, tx sdk.Tx, simulate bool) {
	_, result, abort := anteHandler(ctx, tx, simulate)
	require.Equal(t, "", result.Log)
	require.False(t, abort)
	require.Equal(t, sdk.CodeOK, result.Code)
	require.True(t, result.IsOK())
}

// run the tx through the anteHandler and ensure it fails with the given code
func checkInvalidTx(t *testing.T, anteHandler sdk.AnteHandler, ctx sdk.Context, tx sdk.Tx, simulate bool, code sdk.CodeType) {
	newCtx, result, abort := anteHandler(ctx, tx, simulate)
	require.True(t, abort)

	require.Equal(t, code, result.Code, fmt.Sprintf("Expected %v, got %v", code, result))
	require.Equal(t, sdk.CodespaceRoot, result.Codespace)

	if code == sdk.CodeOutOfGas {
		_, ok := tx.(types.StdTx)
		require.True(t, ok, "tx must be in form auth.types.StdTx")
		// GasWanted set correctly
		require.True(t, result.GasUsed > result.GasWanted, "GasUsed not greated than GasWanted")
		// Check that context is set correctly
		require.Equal(t, result.GasUsed, newCtx.GasMeter().GasConsumed(), "Context not updated correctly")
	}
}
