package auth_test

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/crypto"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/auth"
	"github.com/maticnetwork/heimdall/auth/types"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
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

func (suite *AnteTestSuite) TestAnteValidation() {
	t, happ, ctx, anteHandler := suite.T(), suite.app, suite.ctx, suite.anteHandler

	// module account validation
	// removing module account
	acc := happ.SupplyKeeper.GetModuleAccount(ctx, authTypes.FeeCollectorName)
	happ.AccountKeeper.RemoveAccount(ctx, acc)
	happ.SupplyKeeper.RemoveModuleAddress(authTypes.FeeCollectorName)

	// keys and addresses
	priv1, _, addr1 := sdkAuth.KeyTestPubAddr()
	msg1 := sdkAuth.NewTestMsg(addr1)
	tx1 := types.NewTestTx(ctx, msg1, priv1, uint64(0), uint64(0)) // use sdk's auth module for msg

	_, result1, _ := checkInvalidTx(t, anteHandler, ctx, tx1, false, sdk.CodeInternal)
	require.Contains(t, result1.Log, "fee_collector module account has not been set")
}

func (suite *AnteTestSuite) TestGasLimit() {
	t, happ, ctx, anteHandler := suite.T(), suite.app, suite.ctx, suite.anteHandler
	ctx = ctx.WithBlockHeight(1)

	// keys and addresses
	priv1, _, addr1 := sdkAuth.KeyTestPubAddr()

	// set the accounts
	acc1 := happ.AccountKeeper.NewAccountWithAddress(ctx, hmTypes.AccAddressToHeimdallAddress(addr1))
	happ.AccountKeeper.SetAccount(ctx, acc1)

	// set default amount for one tx
	amt, _ := sdk.NewIntFromString(authTypes.DefaultTxFees)
	acc1.SetCoins(sdk.NewCoins(sdk.NewCoin(authTypes.FeeToken, amt)))
	happ.AccountKeeper.SetAccount(ctx, acc1)

	// get stored account
	acc1 = happ.AccountKeeper.GetAccount(ctx, acc1.GetAddress())

	// msg and signatures
	var tx sdk.Tx
	msg := sdkAuth.NewTestMsg(addr1)

	// get params
	params := happ.AccountKeeper.GetParams(ctx)
	params.MaxTxGas = params.SigVerifyCostSecp256k1 - 1
	happ.AccountKeeper.SetParams(ctx, params)

	// test good tx from one signer
	tx = types.NewTestTx(ctx, msg, priv1, acc1.GetAccountNumber(), uint64(0))
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeOutOfGas)
}

func (suite *AnteTestSuite) TestCheckpointGasLimit() {
	t, happ, ctx, anteHandler := suite.T(), suite.app, suite.ctx, suite.anteHandler
	ctx = ctx.WithBlockHeight(1)

	// keys and addresses
	priv1, _, addr1 := sdkAuth.KeyTestPubAddr()
	priv2, _, addr2 := sdkAuth.KeyTestPubAddr()

	// set the accounts
	acc1 := happ.AccountKeeper.NewAccountWithAddress(ctx, hmTypes.AccAddressToHeimdallAddress(addr1))
	amt1, _ := sdk.NewIntFromString(authTypes.DefaultTxFees)
	acc1.SetCoins(sdk.NewCoins(sdk.NewCoin(authTypes.FeeToken, amt1)))
	happ.AccountKeeper.SetAccount(ctx, acc1)
	acc1 = happ.AccountKeeper.GetAccount(ctx, acc1.GetAddress()) // get stored account

	acc2 := happ.AccountKeeper.NewAccountWithAddress(ctx, hmTypes.AccAddressToHeimdallAddress(addr2))
	amt2, _ := sdk.NewIntFromString(authTypes.DefaultTxFees)
	acc2.SetCoins(sdk.NewCoins(sdk.NewCoin(authTypes.FeeToken, amt2)))
	happ.AccountKeeper.SetAccount(ctx, acc2)
	acc2 = happ.AccountKeeper.GetAccount(ctx, acc2.GetAddress()) // get stored account

	// msg and signatures
	var tx sdk.Tx
	msg := sdkAuth.NewTestMsg(addr1)

	// test good tx from one signer
	tx = types.NewTestTx(ctx, msg, priv1, acc1.GetAccountNumber(), uint64(0))
	_, result, _ := checkValidTx(t, anteHandler, ctx, tx, false)

	// get params
	params := happ.AccountKeeper.GetParams(ctx)
	require.Equal(t, params.MaxTxGas, result.GasWanted)

	// checkpoint msg

	cmsg := TestCheckpointMsg{*sdkAuth.NewTestMsg(addr2)}
	// test good tx from one signer
	tx = types.NewTestTx(ctx, sdk.Msg(&cmsg), priv2, acc2.GetAccountNumber(), uint64(0))
	_, result, _ = checkValidTx(t, anteHandler, ctx, tx, false)
	// check gas wanted for checkpoint msg
	// require.Equal(t, uint64(10000000), uint64(result.GasWanted))
}

func (suite *AnteTestSuite) TestStdTx() {
	t, happ, ctx, anteHandler := suite.T(), suite.app, suite.ctx, suite.anteHandler

	// keys and addresses
	priv1, _, addr1 := sdkAuth.KeyTestPubAddr()

	// test stdTx
	fee := sdkAuth.NewTestStdFee()
	msg := sdkAuth.NewTestMsg(addr1)
	msgs := []sdk.Msg{msg}

	// test no signatures
	privs, accNums, seqs := []crypto.PrivKey{}, []uint64{}, []uint64{}
	tx := sdkAuth.NewTestTx(ctx, msgs, privs, accNums, seqs, fee)

	// Check no signatures fails
	_, result1, _ := checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeInternal)
	require.Contains(t, result1.Log, "tx must be StdTx")

	// test memo
	params := happ.AccountKeeper.GetParams(ctx)
	params.MaxMemoCharacters = 5 // setting 5 as max length temporary
	happ.AccountKeeper.SetParams(ctx, params)

	msg2 := sdkAuth.NewTestMsg(addr1)
	memo := "more than 5 length memo"
	tx2 := types.NewTestTxWithMemo(ctx, msg2, priv1, uint64(0), uint64(0), memo) // use sdk's auth module for msg

	checkInvalidTx(t, anteHandler, ctx, tx2, false, sdk.CodeMemoTooLarge)

	// test tx fees
	params.TxFees = "non integer" // setting non integer
	happ.AccountKeeper.SetParams(ctx, params)
	tx2 = types.NewTestTx(ctx, msg2, priv1, uint64(0), uint64(0)) // use sdk's auth module for msg

	_, result2, _ := checkInvalidTx(t, anteHandler, ctx, tx2, false, sdk.CodeInternal)
	require.Contains(t, result2.Log, "Invalid param tx fees")
}

func (suite *AnteTestSuite) TestSigErrors() {
	t, _, ctx, anteHandler := suite.T(), suite.app, suite.ctx, suite.anteHandler

	// keys and addresses
	priv1, _, addr1 := sdkAuth.KeyTestPubAddr()
	priv2, _, addr2 := sdkAuth.KeyTestPubAddr()

	// test no signers
	msg1 := sdkAuth.NewTestMsg()
	tx1 := types.NewTestTx(ctx, msg1, priv1, uint64(0), uint64(0)) // use sdk's auth module for msg

	// Check no signatures fails
	require.Equal(t, 0, len(msg1.GetSigners()))
	checkInvalidTx(t, anteHandler, ctx, tx1, false, sdk.CodeUnauthorized)

	// unknown address error
	msg2 := sdkAuth.NewTestMsg(addr1) // using first address
	tx2 := types.NewTestTx(ctx, msg2, priv2, uint64(0), uint64(0))

	// Check no signatures fails
	checkInvalidTx(t, anteHandler, ctx, tx2, false, sdk.CodeUnknownAddress)

	// multi signers
	msg3 := sdkAuth.NewTestMsg(addr1, addr2) // using first address
	tx3 := types.NewTestTx(ctx, msg3, priv1, uint64(0), uint64(0))

	// Check no signatures fails
	checkInvalidTx(t, anteHandler, ctx, tx3, false, sdk.CodeUnauthorized)
}

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

// Test logic around account number checking with many signers when BlockHeight is 0.
func (suite *AnteTestSuite) TestAccountNumbersAtBlockHeightZero() {
	t, happ, ctx, anteHandler := suite.T(), suite.app, suite.ctx, suite.anteHandler

	// keys and addresses
	priv1, _, addr1 := sdkAuth.KeyTestPubAddr()
	priv2, _, addr2 := sdkAuth.KeyTestPubAddr()

	// set the accounts, we don't need the acc numbers as it is in the genesis block
	acc1 := happ.AccountKeeper.NewAccountWithAddress(ctx, hmTypes.AccAddressToHeimdallAddress(addr1))
	acc1.SetCoins(simulation.RandomFeeCoins())
	happ.AccountKeeper.SetAccount(ctx, acc1)
	acc2 := happ.AccountKeeper.NewAccountWithAddress(ctx, hmTypes.AccAddressToHeimdallAddress(addr2))
	acc2.SetCoins(simulation.RandomFeeCoins())
	require.NoError(t, acc2.SetAccountNumber(100))
	happ.AccountKeeper.SetAccount(ctx, acc2)

	// msg and signatures
	var tx sdk.Tx
	msg1 := sdkAuth.NewTestMsg(addr1)
	msg2 := sdkAuth.NewTestMsg(addr2)

	acc1 = happ.AccountKeeper.GetAccount(ctx, hmTypes.AccAddressToHeimdallAddress(addr1))
	acc2 = happ.AccountKeeper.GetAccount(ctx, hmTypes.AccAddressToHeimdallAddress(addr2))
	// accNumber1 := acc1.GetAccountNumber()
	// accNumber2 := acc2.GetAccountNumber()

	// test good tx from one signer
	tx = types.NewTestTx(ctx, msg1, priv1, uint64(0), uint64(0))
	checkValidTx(t, anteHandler, ctx, tx, false)

	// test good tx from one signer
	tx = types.NewTestTx(ctx, msg1, priv1, uint64(0), uint64(1))
	checkValidTx(t, anteHandler, ctx, tx, false)

	// // new tx from wrong account number
	tx = types.NewTestTx(ctx, msg2, priv2, uint64(1), uint64(1))
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeUnauthorized)

	// from correct account number but wrong private key
	tx = types.NewTestTx(ctx, msg2, priv1, uint64(1), uint64(0)) // with private key 1
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeUnauthorized)

	// from correct account number but wrong private key
	tx = types.NewTestTx(ctx, msg2, priv2, uint64(0), uint64(0)) // with private key 2 (account 2)
	checkValidTx(t, anteHandler, ctx, tx, false)
}

func (suite *AnteTestSuite) TestSequences() {
	t, happ, ctx, anteHandler := suite.T(), suite.app, suite.ctx, suite.anteHandler
	ctx = ctx.WithBlockHeight(1)

	// keys and addresses
	priv1, _, addr1 := sdkAuth.KeyTestPubAddr()
	priv2, _, addr2 := sdkAuth.KeyTestPubAddr()

	// set the accounts, we don't need the acc numbers as it is in the genesis block
	acc1 := happ.AccountKeeper.NewAccountWithAddress(ctx, hmTypes.AccAddressToHeimdallAddress(addr1))
	acc1.SetCoins(simulation.RandomFeeCoins())
	happ.AccountKeeper.SetAccount(ctx, acc1)
	acc2 := happ.AccountKeeper.NewAccountWithAddress(ctx, hmTypes.AccAddressToHeimdallAddress(addr2))
	acc2.SetCoins(simulation.RandomFeeCoins())
	require.NoError(t, acc2.SetAccountNumber(100))
	happ.AccountKeeper.SetAccount(ctx, acc2)

	// msg and signatures
	var tx sdk.Tx
	msg1 := sdkAuth.NewTestMsg(addr1)
	msg2 := sdkAuth.NewTestMsg(addr2)

	acc1 = happ.AccountKeeper.GetAccount(ctx, hmTypes.AccAddressToHeimdallAddress(addr1))
	acc2 = happ.AccountKeeper.GetAccount(ctx, hmTypes.AccAddressToHeimdallAddress(addr2))
	accNumber1 := acc1.GetAccountNumber()
	accNumber2 := acc2.GetAccountNumber()

	// test good tx from one signer
	tx = types.NewTestTx(ctx, msg1, priv1, accNumber1, uint64(0))
	checkValidTx(t, anteHandler, ctx, tx, false)

	// test good tx from one signer
	tx = types.NewTestTx(ctx, msg1, priv1, accNumber1, uint64(1))
	checkValidTx(t, anteHandler, ctx, tx, false)

	// // new tx from wrong account number
	tx = types.NewTestTx(ctx, msg2, priv2, accNumber2, uint64(1))
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeUnauthorized)

	// from correct account number but wrong private key
	tx = types.NewTestTx(ctx, msg2, priv1, accNumber2, uint64(0)) // with private key 1
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeUnauthorized)

	// from correct account number but wrong private key
	tx = types.NewTestTx(ctx, msg2, priv2, accNumber2, uint64(0)) // with private key 2 (account 2)
	checkValidTx(t, anteHandler, ctx, tx, false)
}

// Test logic around fee deduction.
func (suite *AnteTestSuite) TestFees() {
	t, happ, ctx, anteHandler := suite.T(), suite.app, suite.ctx, suite.anteHandler

	// keys and addresses
	priv1, _, addr1 := sdkAuth.KeyTestPubAddr()

	// set the accounts
	acc1 := happ.AccountKeeper.NewAccountWithAddress(ctx, hmTypes.AccAddressToHeimdallAddress(addr1))
	happ.AccountKeeper.SetAccount(ctx, acc1)

	// msg and signatures
	var tx sdk.Tx
	msg1 := sdkAuth.NewTestMsg(addr1)
	acc1 = happ.AccountKeeper.GetAccount(ctx, hmTypes.AccAddressToHeimdallAddress(addr1))
	tx = types.NewTestTx(ctx, msg1, priv1, uint64(0), uint64(0))
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeInsufficientFunds)

	// set some coins
	acc1.SetCoins(sdk.NewCoins(sdk.NewInt64Coin(authTypes.FeeToken, 149)))
	happ.AccountKeeper.SetAccount(ctx, acc1)
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeInsufficientFunds)

	require.True(t, happ.SupplyKeeper.GetModuleAccount(ctx, authTypes.FeeCollectorName).GetCoins().Empty())
	require.True(sdk.IntEq(t, happ.AccountKeeper.GetAccount(ctx, hmTypes.AccAddressToHeimdallAddress(addr1)).GetCoins().AmountOf(authTypes.FeeToken), sdk.NewInt(149)))

	amt, _ := sdk.NewIntFromString(authTypes.DefaultTxFees)
	acc1.SetCoins(sdk.NewCoins(sdk.NewCoin(authTypes.FeeToken, amt)))
	happ.AccountKeeper.SetAccount(ctx, acc1)
	checkValidTx(t, anteHandler, ctx, tx, false)

	require.True(sdk.IntEq(t, happ.SupplyKeeper.GetModuleAccount(ctx, types.FeeCollectorName).GetCoins().AmountOf(authTypes.FeeToken), amt))
	require.True(sdk.IntEq(t, happ.AccountKeeper.GetAccount(ctx, hmTypes.AccAddressToHeimdallAddress(addr1)).GetCoins().AmountOf(authTypes.FeeToken), sdk.NewInt(0)))

	// try to send tx again
	checkInvalidTx(t, anteHandler, ctx, tx, false, sdk.CodeInsufficientFunds)
}

//
// utils
//

// run the tx through the anteHandler and ensure its valid
func checkValidTx(t *testing.T, anteHandler sdk.AnteHandler, ctx sdk.Context, tx sdk.Tx, simulate bool) (sdk.Context, sdk.Result, bool) {
	newCtx, result, abort := anteHandler(ctx, tx, simulate)
	require.Equal(t, "", result.Log)
	require.False(t, abort)
	require.Equal(t, sdk.CodeOK, result.Code)
	require.True(t, result.IsOK())
	return newCtx, result, abort
}

// run the tx through the anteHandler and ensure it fails with the given code
func checkInvalidTx(t *testing.T, anteHandler sdk.AnteHandler, ctx sdk.Context, tx sdk.Tx, simulate bool, code sdk.CodeType) (sdk.Context, sdk.Result, bool) {
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

	return newCtx, result, abort
}

//
// Test checkpoint
//

var _ sdk.Msg = (*TestCheckpointMsg)(nil)

// msg type for testing
type TestCheckpointMsg struct {
	sdk.TestMsg
}

func (msg *TestCheckpointMsg) Route() string { return "checkpoint" }
func (msg *TestCheckpointMsg) Type() string  { return "checkpoint" }
