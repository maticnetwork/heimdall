package app_test

import (
	"bytes"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmTypes "github.com/tendermint/tendermint/types"

	app "github.com/maticnetwork/heimdall/app"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

//
// Test suite
//

// SideTxProcessorTestSuite test suite for side-tx processor
type SideTxProcessorTestSuite struct {
	suite.Suite

	app     *app.HeimdallApp
	ctx     sdk.Context
	encoder sdk.TxEncoder
}

func (suite *SideTxProcessorTestSuite) SetupTest() {
	isCheckTx := false

	happ := app.Setup(isCheckTx)
	ctx := happ.NewContext(isCheckTx, abci.Header{})
	cdc := codec.New()
	registerTestCodec(cdc)
	encoder := authTypes.DefaultTxEncoder(cdc)

	suite.app, suite.ctx, suite.encoder = happ, ctx, encoder
}

func TestSideTxProcessorTestSuite(t *testing.T) {
	suite.Run(t, new(SideTxProcessorTestSuite))
}

//
// Test cases
//

func (suite *SideTxProcessorTestSuite) TestPostDeliverTxHandler() {
	t, happ, ctx, encoder := suite.T(), suite.app, suite.ctx, suite.encoder

	{
		// msg and signatures
		tx := hmTypes.BaseTx{
			Msg: msgCounter{Counter: 1},
		}
		txBytes, err := encoder(tx)
		require.Nil(t, err, "There should be no error while encoding tx")

		// test post deliver with height 2

		ctx = ctx.WithTxBytes(txBytes)
		happ.PostDeliverTxHandler(ctx, tx, sdk.Result{})
		resultTxBytes := happ.SidechannelKeeper.GetTx(ctx, 2, tmTypes.Tx(txBytes).Hash())
		require.Nil(t, resultTxBytes)

		// test post deliver with height 20 without side-tx
		ctx = ctx.WithBlockHeight(20).WithTxBytes(txBytes)
		happ.PostDeliverTxHandler(ctx, tx, sdk.Result{})
		resultTxBytes = happ.SidechannelKeeper.GetTx(ctx, 20, tmTypes.Tx(txBytes).Hash())
		require.Nil(t, resultTxBytes)
	}

	// test post deliver with height 20 with side-tx
	{
		// msg and signatures
		tx := hmTypes.BaseTx{
			Msg: msgSideCounter{Counter: 1},
		}
		txBytes, err := encoder(tx)
		require.Nil(t, err, "There should be no error while encoding tx")

		ctx = ctx.WithBlockHeight(20).WithTxBytes(txBytes)
		happ.PostDeliverTxHandler(ctx, tx, sdk.Result{})
		resultTxBytes := happ.SidechannelKeeper.GetTx(ctx, 20, tmTypes.Tx(txBytes).Hash())
		require.NotNil(t, resultTxBytes, "Stored tx bytes shouldn't be nil")
		require.Greater(t, len(resultTxBytes), 0)
		require.True(t, bytes.Equal(txBytes, resultTxBytes), "Stored tx bytes should same as actual tx bytes")
	}
}

func (suite *SideTxProcessorTestSuite) TestDeliverSideTxHandler() {
	// t, happ, ctx, encoder := suite.T(), suite.app, suite.ctx, suite.encoder
}

func (suite *SideTxProcessorTestSuite) TestBeginSideBlocker() {
	// t, happ, ctx, encoder := suite.T(), suite.app, suite.ctx, suite.encoder
}

//
// utils
//

func registerTestCodec(cdc *codec.Codec) {
	// register Tx, Msg
	sdk.RegisterCodec(cdc)

	// register test types
	cdc.RegisterConcrete(&msgCounter{}, "cosmos-sdk/baseapp/msgCounter", nil)
	cdc.RegisterConcrete(&msgSideCounter{}, "cosmos-sdk/baseapp/msgSideCounter", nil)
}

const (
	routeMsgCounter     = "msgCounter"
	routeMsgSideCounter = "msgSideCounter"
)

// ValidateBasic() fails on negative counters.
// Otherwise it's up to the handlers
type msgCounter struct {
	Counter       int64
	FailOnHandler bool
}

// Implements Msg
func (msg msgCounter) Route() string                { return routeMsgCounter }
func (msg msgCounter) Type() string                 { return "counter1" }
func (msg msgCounter) GetSignBytes() []byte         { return nil }
func (msg msgCounter) GetSigners() []sdk.AccAddress { return nil }
func (msg msgCounter) ValidateBasic() sdk.Error {
	if msg.Counter >= 0 {
		return nil
	}
	return sdk.ErrInvalidSequence("counter should be a non-negative integer.")
}

// ValidateBasic() fails on negative counters.
// Otherwise it's up to the handlers
type msgSideCounter struct {
	Counter int64
}

func (msg msgSideCounter) GetSideSignBytes() []byte {
	return nil
}

func (msg msgSideCounter) Route() string                { return routeMsgSideCounter }
func (msg msgSideCounter) Type() string                 { return "counter1" }
func (msg msgSideCounter) GetSignBytes() []byte         { return nil }
func (msg msgSideCounter) GetSigners() []sdk.AccAddress { return nil }
func (msg msgSideCounter) ValidateBasic() sdk.Error {
	if msg.Counter >= 0 {
		return nil
	}
	return sdk.ErrInvalidSequence("counter should be a non-negative integer.")
}
