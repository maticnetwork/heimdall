package topup_test

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/app"
	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper/mocks"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/types/simulation"
	chainTypes "github.com/maticnetwork/heimdall/x/chainmanager/types"
	"github.com/maticnetwork/heimdall/x/topup"
	"github.com/maticnetwork/heimdall/x/topup/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// HandlerTestSuite integrate test suite context object
type HandlerTestSuite struct {
	suite.Suite

	app            *app.HeimdallApp
	ctx            sdk.Context
	cliCtx         client.Context
	querier        sdk.Querier
	handler        sdk.Handler
	contractCaller mocks.IContractCaller
	chainParams    chainTypes.Params
}

// SetupTest setup all necessary things for querier tesing
func (suite *HandlerTestSuite) SetupTest() {
	suite.app, suite.ctx, suite.cliCtx = createTestApp(false)

	suite.contractCaller = mocks.IContractCaller{}
	suite.handler = topup.NewHandler(suite.app.TopupKeeper, &suite.contractCaller)
	suite.chainParams = suite.app.ChainKeeper.GetParams(suite.ctx)
}

// TestHandlerTestSuite
func TestHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestHandleMsgUnknown() {
	t, _, ctx := suite.T(), suite.app, suite.ctx

	result, err := suite.handler(ctx, nil)
	require.NotNil(t, err)
	require.Nil(t, result)
	// require.False(t, result.IsOK())
}

func (suite *HandlerTestSuite) TestHandleMsgTopup() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	txHash := hmCommonTypes.BytesToHeimdallHash([]byte("topup hash"))
	logIndex := r1.Uint64()
	blockNumber := r1.Uint64()

	_, _, addr := testdata.KeyTestPubAddr()
	fee := sdk.NewInt(100000000000000000)
	generated_address, _ := sdk.AccAddressFromHex(addr.String())

	t.Run("Success", func(t *testing.T) {
		msg := types.NewMsgTopup(
			generated_address,
			generated_address,
			fee,
			txHash,
			uint64(logIndex),
			uint64(blockNumber),
		)

		// handler
		result, err := suite.handler(ctx, &msg)
		require.Nil(t, err)
		require.NotNil(t, result)
		// require.True(t, result.IsOK(), "Expected topup to be done, but failed")
	})

	t.Run("OlderTx", func(t *testing.T) {
		msg := types.NewMsgTopup(
			generated_address,
			generated_address,
			fee,
			txHash,
			uint64(logIndex),
			uint64(blockNumber),
		)

		// sequence id
		blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
		sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
		sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

		// set sequence
		app.TopupKeeper.SetTopupSequence(ctx, sequence.String())

		// handler
		result, err := suite.handler(ctx, &msg)
		//TODO: check if the error code is the same as CodeOldTx
		require.Error(t, err)
		require.Equal(t, hmCommon.ErrOldTx, err)
		require.Nil(t, result)
		// require.False(t, result.IsOK(), "Expected topup to be failed, but succeeded")
		// require.Equal(t, common.CodeOldTx, result.Code)
	})
}

// TODO: use bankKeeper for set and get coins (coins -> balances)
func (suite *HandlerTestSuite) TestHandleMsgWithdrawFee() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	t.Run("FullAmount", func(t *testing.T) {
		_, _, addr := testdata.KeyTestPubAddr()
		tAddr, err := sdk.AccAddressFromHex(addr.String())
		msg := types.NewMsgWithdrawFee(
			tAddr,
			sdk.NewInt(0),
		)

		// execute handler
		result, err := suite.handler(ctx, &msg)
		// require.False(t, result.IsOK(), "Expected topup to be failed without fee tokens, but succeeded")
		// require.Equal(t, types.CodeNoBalanceToWithdraw, result.Code)
		require.Error(t, err)
		require.Equal(t, types.ErrNoBalanceToWithdraw, err)

		// set coins
		coins := simulation.RandomFeeCoins()
		acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, tAddr)
		acc1.BankKeeper.SetBalance(coins)
		// acc1.SetCoins(coins)
		app.AccountKeeper.SetAccount(ctx, acc1)

		// check if coins > 0
		require.True(t, acc1.GetCoins().AmountOf(hmTypes.FeeToken).GT(sdk.NewInt(0)))

		// execute handler
		result, err = suite.handler(ctx, &msg)
		// require.True(t, result.IsOK(), "Expected topup to be succeed with fee tokens, but failed")
		require.Greater(t, len(result.Events), 0)

		// check if account has zero
		acc1 = app.AccountKeeper.GetAccount(ctx, tAddr)
		require.True(t, acc1.GetCoins().AmountOf(hmTypes.FeeToken).IsZero())
	})

	t.Run("PartialAmount", func(t *testing.T) {
		_, _, addr := testdata.KeyTestPubAddr()
		tAddr, err := sdk.AccAddressFromHex(addr.String())
		// set coins
		coins := simulation.RandomFeeCoins()
		acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, tAddr)
		acc1.SetCoins(coins)
		app.AccountKeeper.SetAccount(ctx, acc1)

		// check if coins > 0
		require.True(t, acc1.GetCoins().AmountOf(hmTypes.FeeToken).GT(sdk.NewInt(0)))

		m, _ := sdk.NewIntFromString("2")
		coins = coins.Sub(sdk.Coins{sdk.Coin{Denom: hmTypes.FeeToken, Amount: m}})
		msg := types.NewMsgWithdrawFee(
			tAddr,
			coins.AmountOf(hmTypes.FeeToken),
		)

		// execute handler
		result, err := suite.handler(ctx, &msg)
		// require.True(t, result.IsOK(), "Expected topup to be succeed with fee tokens (partial amount), but failed")
		require.Nil(t, err)
		require.NotNil(t, result)
		require.Greater(t, len(result.Events), 0)

		// check if account has 1 tok
		acc1 = app.AccountKeeper.GetAccount(ctx, tAddr)
		require.True(t, acc1.GetCoins().AmountOf(hmTypes.FeeToken).Equal(m))
	})

	t.Run("NotEnoughAmount", func(t *testing.T) {
		_, _, addr := testdata.KeyTestPubAddr()
		tAddr, err := sdk.AccAddressFromHex(addr.String())
		// set coins
		coins := simulation.RandomFeeCoins()
		acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, tAddr)
		acc1.SetCoins(coins)
		app.AccountKeeper.SetAccount(ctx, acc1)

		m, _ := sdk.NewIntFromString("1")
		coins = coins.Add(sdk.Coins{sdk.Coin{Denom: hmTypes.FeeToken, Amount: m}})
		msg := types.NewMsgWithdrawFee(
			tAddr,
			coins.AmountOf(hmTypes.FeeToken),
		)

		result, err := suite.handler(ctx, &msg)
		require.Error(t, err)

		// require.False(t, result.IsOK(), "Expected withdraw to be failed while withdrawing more than account's coins")
	})
}
