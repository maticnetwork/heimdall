package topup_test

import (
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkAuth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/maticnetwork/heimdall/app"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	chainTypes "github.com/maticnetwork/heimdall/chainmanager/types"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper/mocks"
	"github.com/maticnetwork/heimdall/topup"
	"github.com/maticnetwork/heimdall/topup/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/simulation"
)

// HandlerTestSuite integrate test suite context object
type HandlerTestSuite struct {
	suite.Suite

	app            *app.HeimdallApp
	ctx            sdk.Context
	cliCtx         context.CLIContext
	handler        sdk.Handler
	contractCaller mocks.IContractCaller
	chainParams    chainTypes.Params
}

// SetupTest setup all necessary things for querier testing
func (suite *HandlerTestSuite) SetupTest() {
	suite.app, suite.ctx, suite.cliCtx = createTestApp(false)

	suite.contractCaller = mocks.IContractCaller{}
	suite.handler = topup.NewHandler(suite.app.TopupKeeper, &suite.contractCaller)
	suite.chainParams = suite.app.ChainKeeper.GetParams(suite.ctx)
}

// TestHandlerTestSuite
func TestHandlerTestSuite(t *testing.T) {
	t.Parallel()
	suite.Run(t, new(HandlerTestSuite))
}

func (suite *HandlerTestSuite) TestHandleMsgUnknown() {
	t, _, ctx := suite.T(), suite.app, suite.ctx

	result := suite.handler(ctx, nil)
	require.False(t, result.IsOK())
}

func (suite *HandlerTestSuite) TestHandleMsgTopup() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	txHash := hmTypes.BytesToHeimdallHash([]byte("topup hash"))
	logIndex := r1.Uint64()
	blockNumber := r1.Uint64()

	_, _, addr := sdkAuth.KeyTestPubAddr()
	fee := sdk.NewInt(100000000000000000)

	t.Run("Success", func(t *testing.T) {
		msg := types.NewMsgTopup(
			hmTypes.BytesToHeimdallAddress(addr.Bytes()),
			hmTypes.BytesToHeimdallAddress(addr.Bytes()),
			fee,
			txHash,
			logIndex,
			blockNumber,
		)

		// handler
		result := suite.handler(ctx, msg)
		require.True(t, result.IsOK(), "Expected topup to be done, but failed")
	})

	t.Run("OlderTx", func(t *testing.T) {
		msg := types.NewMsgTopup(
			hmTypes.BytesToHeimdallAddress(addr.Bytes()),
			hmTypes.BytesToHeimdallAddress(addr.Bytes()),
			fee,
			txHash,
			logIndex,
			blockNumber,
		)

		// sequence id
		blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
		sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
		sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

		// set sequence
		app.TopupKeeper.SetTopupSequence(ctx, sequence.String())

		// handler
		result := suite.handler(ctx, msg)
		require.False(t, result.IsOK(), "Expected topup to be failed, but succeeded")
		require.Equal(t, common.CodeOldTx, result.Code)
	})
}

func (suite *HandlerTestSuite) TestHandleMsgWithdrawFee() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	t.Run("FullAmount", func(t *testing.T) {
		_, _, addr := sdkAuth.KeyTestPubAddr()

		msg := types.NewMsgWithdrawFee(
			hmTypes.BytesToHeimdallAddress(addr.Bytes()),
			sdk.NewInt(0),
		)

		// execute handler
		result := suite.handler(ctx, msg)
		require.False(t, result.IsOK(), "Expected topup to be failed without fee tokens, but succeeded")
		require.Equal(t, types.CodeNoBalanceToWithdraw, result.Code)

		// set coins
		coins := simulation.RandomFeeCoins()
		acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, hmTypes.AccAddressToHeimdallAddress(addr))
		err := acc1.SetCoins(coins)
		require.NoError(t, err)
		app.AccountKeeper.SetAccount(ctx, acc1)

		// check if coins > 0
		require.True(t, acc1.GetCoins().AmountOf(authTypes.FeeToken).GT(sdk.NewInt(0)))

		// execute handler
		result = suite.handler(ctx, msg)
		require.True(t, result.IsOK(), "Expected topup to be succeed with fee tokens, but failed")
		require.Greater(t, len(result.Events), 0)

		// check if account has zero
		acc1 = app.AccountKeeper.GetAccount(ctx, hmTypes.AccAddressToHeimdallAddress(addr))
		require.True(t, acc1.GetCoins().AmountOf(authTypes.FeeToken).IsZero())
	})

	t.Run("PartialAmount", func(t *testing.T) {
		_, _, addr := sdkAuth.KeyTestPubAddr()

		// set coins
		coins := simulation.RandomFeeCoins()
		acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, hmTypes.AccAddressToHeimdallAddress(addr))
		err := acc1.SetCoins(coins)
		require.NoError(t, err)
		app.AccountKeeper.SetAccount(ctx, acc1)

		// check if coins > 0
		require.True(t, acc1.GetCoins().AmountOf(authTypes.FeeToken).GT(sdk.NewInt(0)))

		m, _ := sdk.NewIntFromString("2")
		coins = coins.Sub(sdk.Coins{sdk.Coin{Denom: authTypes.FeeToken, Amount: m}})
		msg := types.NewMsgWithdrawFee(
			hmTypes.BytesToHeimdallAddress(addr.Bytes()),
			coins.AmountOf(authTypes.FeeToken),
		)

		// execute handler
		result := suite.handler(ctx, msg)
		require.True(t, result.IsOK(), "Expected topup to be succeed with fee tokens (partial amount), but failed")
		require.Greater(t, len(result.Events), 0)

		// check if account has 1 tok
		acc1 = app.AccountKeeper.GetAccount(ctx, hmTypes.AccAddressToHeimdallAddress(addr))
		require.True(t, acc1.GetCoins().AmountOf(authTypes.FeeToken).Equal(m))
	})

	t.Run("NotEnoughAmount", func(t *testing.T) {
		_, _, addr := sdkAuth.KeyTestPubAddr()

		// set coins
		coins := simulation.RandomFeeCoins()
		acc1 := app.AccountKeeper.NewAccountWithAddress(ctx, hmTypes.AccAddressToHeimdallAddress(addr))
		err := acc1.SetCoins(coins)
		require.NoError(t, err)
		app.AccountKeeper.SetAccount(ctx, acc1)

		m, _ := sdk.NewIntFromString("1")
		coins = coins.Add(sdk.Coins{sdk.Coin{Denom: authTypes.FeeToken, Amount: m}})
		msg := types.NewMsgWithdrawFee(
			hmTypes.BytesToHeimdallAddress(addr.Bytes()),
			coins.AmountOf(authTypes.FeeToken),
		)

		result := suite.handler(ctx, msg)
		require.False(t, result.IsOK(), "Expected withdraw to be failed while withdrawing more than account's coins")
	})
}
