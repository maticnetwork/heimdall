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

var testTxStateData1 = []byte("test-tx-state1")
var testTxStateData2 = []byte("test-tx-state2")

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
	happ.SetCodec(cdc) // set to app
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
	t, happ, ctx, encoder := suite.T(), suite.app, suite.ctx, suite.encoder

	// msg and signatures
	msg := msgSideCounter{Counter: 1}
	tx := hmTypes.BaseTx{
		Msg: msg,
	}
	txBytes, err := encoder(tx)
	require.Nil(t, err, "There should be no error while encoding tx")

	// deliver side thandler
	require.Panics(t, func() {
		happ.DeliverSideTxHandler(ctx, tx, abci.RequestDeliverSideTx{
			Tx: tmTypes.Tx(txBytes),
		})
	}, "Tx without side-tx handler should panic")

	// Set router (testing only)
	router := hmTypes.NewSideRouter()
	router.AddRoute(routeMsgSideCounter, &hmTypes.SideHandlers{
		SideTxHandler: func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
			return abci.ResponseDeliverSideTx{
				Result: abci.SideTxResultType_Yes,
			}
		},
		PostTxHandler: func(ctx sdk.Context, msg sdk.Msg, sideTxResult abci.SideTxResultType) sdk.Result {
			return sdk.Result{}
		},
	})
	happ.SetSideRouter(router)

	t.Run("YesVote", func(t *testing.T) {
		var res abci.ResponseDeliverSideTx
		// test route
		require.NotPanics(t, func() {
			res = happ.DeliverSideTxHandler(ctx, tx, abci.RequestDeliverSideTx{
				Tx: tmTypes.Tx(txBytes),
			})
		})

		require.Equal(t, abci.SideTxResultType_Yes, res.GetResult(), "Result from deliver side-tx should be vote `Yes`")
		require.Equal(t, len(msg.GetSideSignBytes()), len(res.Data), "Result from deliver side-tx data should be same as expected")
	})

	t.Run("NoVote", func(t *testing.T) {
		resultData := []byte("hello")
		router := hmTypes.NewSideRouter()
		router.AddRoute(routeMsgSideCounter, &hmTypes.SideHandlers{
			SideTxHandler: func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
				return abci.ResponseDeliverSideTx{
					Result: abci.SideTxResultType_No,
					Data:   resultData,
				}
			},
			PostTxHandler: func(ctx sdk.Context, msg sdk.Msg, sideTxResult abci.SideTxResultType) sdk.Result {
				return sdk.Result{}
			},
		})
		happ.SetSideRouter(router)

		res := happ.DeliverSideTxHandler(ctx, tx, abci.RequestDeliverSideTx{
			Tx: tmTypes.Tx(txBytes),
		})

		require.Equal(t, abci.SideTxResultType_No, res.GetResult(), "Result from deliver side-tx should be vote `No`")
		require.Equal(t, resultData, res.Data, "Result from deliver side-tx data should be same as expected")
	})

	t.Run("SkipVote", func(t *testing.T) {
		router := hmTypes.NewSideRouter()
		router.AddRoute(routeMsgSideCounter, &hmTypes.SideHandlers{
			SideTxHandler: func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
				return abci.ResponseDeliverSideTx{
					Code: uint32(sdk.CodeInternal),
				}
			},
			PostTxHandler: func(ctx sdk.Context, msg sdk.Msg, sideTxResult abci.SideTxResultType) sdk.Result {
				return sdk.Result{}
			},
		})
		happ.SetSideRouter(router)
		res := happ.DeliverSideTxHandler(ctx, tx, abci.RequestDeliverSideTx{
			Tx: tmTypes.Tx(txBytes),
		})

		require.Equal(t, abci.SideTxResultType_Skip, res.GetResult(), "Result from deliver side-tx should be vote `Skip` if result is not OK")
	})

	t.Run("State", func(t *testing.T) {
		// testing by storing random txs to store
		happ.SidechannelKeeper.SetTx(ctx, 800, testTxStateData1)
		require.Equal(t, 1, len(happ.SidechannelKeeper.GetTxs(ctx, 800)), "It should set state in storage in normal case")

		router := hmTypes.NewSideRouter()
		router.AddRoute(routeMsgSideCounter, &hmTypes.SideHandlers{
			SideTxHandler: func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
				// set random tx in deliver side-tx
				// it should reset the storage (shouldn't change the state)
				happ.SidechannelKeeper.SetTx(ctx, 800, testTxStateData2)
				require.Equal(t, 2, len(happ.SidechannelKeeper.GetTxs(ctx, 800)), "It should set state in storage temperory in side-tx handler")

				return abci.ResponseDeliverSideTx{
					Result: abci.SideTxResultType_Yes,
				}
			},
			PostTxHandler: func(ctx sdk.Context, msg sdk.Msg, sideTxResult abci.SideTxResultType) sdk.Result {
				return sdk.Result{}
			},
		})
		happ.SetSideRouter(router)
		res := happ.DeliverSideTxHandler(ctx, tx, abci.RequestDeliverSideTx{
			Tx: tmTypes.Tx(txBytes),
		})

		require.Equal(t, abci.SideTxResultType_Yes, res.GetResult(), "It should return `yes` vote from deliver side-tx")
		// testing by fetching txs from store
		require.Equal(t, 1, len(happ.SidechannelKeeper.GetTxs(ctx, 800)), "It shouldn't change state in deliver side-tx")
	})

	t.Run("Panic", func(t *testing.T) {
		router := hmTypes.NewSideRouter()
		router.AddRoute(routeMsgSideCounter, &hmTypes.SideHandlers{
			SideTxHandler: func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
				panic("TestSideTxPanic")
			},
			PostTxHandler: func(ctx sdk.Context, msg sdk.Msg, sideTxResult abci.SideTxResultType) sdk.Result {
				return sdk.Result{}
			},
		})
		happ.SetSideRouter(router)
		require.Panics(t, func() {
			happ.DeliverSideTxHandler(ctx, tx, abci.RequestDeliverSideTx{
				Tx: tmTypes.Tx(txBytes),
			})
		}, "It should panic because of side-tx handler panic")

		var res abci.ResponseDeliverSideTx
		require.NotPanics(t, func() {
			res = happ.DeliverSideTx(abci.RequestDeliverSideTx{
				Tx: tmTypes.Tx(txBytes),
			})
		}, "It shouldn't panic due to recover in app.DeliverSideTx")

		require.Equal(t, abci.SideTxResultType_Skip, res.GetResult(), "It should return `skip` vote due to panic")
	})
}

func (suite *SideTxProcessorTestSuite) TestBeginSideBlocker() {
	t, happ, ctx, encoder := suite.T(), suite.app, suite.ctx, suite.encoder

	// msg and signatures
	msg := msgSideCounter{Counter: 1}
	tx := hmTypes.BaseTx{
		Msg: msg,
	}
	txBytes, err := encoder(tx)
	require.NotEmpty(t, txBytes)
	require.Nil(t, err, "There should be no error while encoding tx")

	txHash := tmTypes.Tx(txBytes).Hash()

	t.Run("Height_2", func(t *testing.T) {
		res := happ.BeginSideBlocker(ctx, abci.RequestBeginSideBlock{})
		require.Equal(t, 0, len(res.Events), "Events from begin side blocker result should be empty for <= 2 blocks")
	})

	t.Run("Height_20", func(t *testing.T) {
		ctx = ctx.WithBlockHeight(20)
		res := happ.BeginSideBlocker(ctx, abci.RequestBeginSideBlock{})
		require.Equal(t, 0, len(res.Events), "Events from begin side blocker result should be empty for empty validators")
	})

	t.Run("Votes", func(t *testing.T) {
		var height int64 = 20
		ctx = ctx.WithBlockHeight(height)

		addr1 := []byte("hello-1")
		addr2 := []byte("hello-2")
		addr3 := []byte("hello-3")
		addr4 := []byte("hello-4")
		// set validators
		happ.SidechannelKeeper.SetValidators(ctx, height, []abci.Validator{
			{Address: addr1, Power: 10},
			{Address: addr2, Power: 20},
			{Address: addr3, Power: 30},
			{Address: addr4, Power: 40},
		})

		res := happ.BeginSideBlocker(ctx, abci.RequestBeginSideBlock{})
		require.Equal(t, 0, len(res.Events), "Events from begin side blocker result should be empty with validators but no sigs")

		// test all types of result
		for key, value := range abci.SideTxResultType_value {
			t.Run(key, func(t *testing.T) {
				// setup router and handler
				router := hmTypes.NewSideRouter()
				handler := &hmTypes.SideHandlers{
					SideTxHandler: func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
						return abci.ResponseDeliverSideTx{}
					},
					PostTxHandler: func(ctx sdk.Context, msg sdk.Msg, sideTxResult abci.SideTxResultType) sdk.Result {
						return sdk.Result{}
					},
				}
				router.AddRoute(routeMsgSideCounter, handler)
				happ.SetSideRouter(router)

				happ.SidechannelKeeper.SetTx(ctx, height-2, txBytes) // set tx in the store for process
				res = happ.BeginSideBlocker(ctx, abci.RequestBeginSideBlock{
					SideTxResults: []abci.SideTxResult{
						{
							TxHash: txHash,
							Sigs: []abci.SideTxSig{
								{
									Result:  abci.SideTxResultType(value),
									Address: addr1,
								},
								{
									Result:  abci.SideTxResultType(value),
									Address: addr2,
								},
								{
									Result:  abci.SideTxResultType(value),
									Address: addr3,
								},
								{
									Result:  abci.SideTxResultType(value),
									Address: addr4,
								},
							},
						},
					},
				})
				require.Equal(t, 0, len(res.Events), "It should have no event")
				require.Nil(t, happ.SidechannelKeeper.GetTx(ctx, height-2, txHash), "Tx should not be present in store after begin block")
			})
		}

		t.Run("SkipWithoutPowerCalculation", func(t *testing.T) {
			// skip without any power calculation
			happ.SidechannelKeeper.SetTx(ctx, height-2, txBytes)
			res = happ.BeginSideBlocker(ctx, abci.RequestBeginSideBlock{})
			require.Equal(t, 0, len(res.Events), "It should have no event with validators")
			require.Nil(t, happ.SidechannelKeeper.GetTx(ctx, height-2, txHash), "Tx should not be present in store after begin block")
		})
	})

	t.Run("State", func(t *testing.T) {
		var height int64 = 20
		ctx = ctx.WithBlockHeight(height)

		addr1 := []byte("hello-1")
		addr2 := []byte("hello-2")
		addr3 := []byte("hello-3")
		addr4 := []byte("hello-4")
		// set validators
		happ.SidechannelKeeper.SetValidators(ctx, height, []abci.Validator{
			{Address: addr1, Power: 10},
			{Address: addr2, Power: 20},
			{Address: addr3, Power: 30},
			{Address: addr4, Power: 40},
		})

		req := abci.RequestBeginSideBlock{
			SideTxResults: []abci.SideTxResult{
				{
					TxHash: txHash,
					Sigs: []abci.SideTxSig{
						{
							Result:  abci.SideTxResultType_No,
							Address: addr1,
						},
						{
							Result:  abci.SideTxResultType_No,
							Address: addr2,
						},
						{
							Result:  abci.SideTxResultType_Yes,
							Address: addr3,
						},
						{
							Result:  abci.SideTxResultType_Yes,
							Address: addr4,
						},
					},
				},
			},
		}

		// should save state on successful execution of  post-tx handler
		{
			// context with new event manager
			ctx = ctx.WithEventManager(sdk.NewEventManager())

			// setup router and handler
			router := hmTypes.NewSideRouter()
			handler := &hmTypes.SideHandlers{
				SideTxHandler: func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
					return abci.ResponseDeliverSideTx{}
				},
				PostTxHandler: func(ctx sdk.Context, msg sdk.Msg, sideTxResult abci.SideTxResultType) sdk.Result {
					require.Equal(t, abci.SideTxResultType_Yes, sideTxResult, "Result `sideTxResult` should be `yes`")

					// try to set random state data to store
					happ.SidechannelKeeper.SetTx(ctx, 700, testTxStateData1)

					// events
					ctx.EventManager().EmitEvents(sdk.Events{
						sdk.NewEvent(
							"property1",
							sdk.NewAttribute("key1", "value1"),
						),
						sdk.NewEvent(
							"property2",
							sdk.NewAttribute("key2", "value2"),
						),
					})

					return sdk.Result{
						Events: ctx.EventManager().Events(),
					}
				},
			}
			router.AddRoute(routeMsgSideCounter, handler)
			happ.SetSideRouter(router)

			happ.SidechannelKeeper.SetTx(ctx, height-2, txBytes) // set tx in the store for process
			res := happ.BeginSideBlocker(ctx, req)
			require.Equal(t, 2, len(res.Events), "It should include correct emitted events")
			require.Nil(t, happ.SidechannelKeeper.GetTx(ctx, height-2, txHash), "Tx should not be present in store after begin block")

			// check if it saved the data
			require.Equal(t, 1, len(happ.SidechannelKeeper.GetTxs(ctx, 700)), "It should save state correctly after successful post-tx execution")
		}

		// shouldn't save state on failed execution of post-tx handler
		{
			// context with new event manager
			ctx = ctx.WithEventManager(sdk.NewEventManager())

			// setup router and handler
			router := hmTypes.NewSideRouter()
			handler := &hmTypes.SideHandlers{
				SideTxHandler: func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
					return abci.ResponseDeliverSideTx{}
				},
				PostTxHandler: func(ctx sdk.Context, msg sdk.Msg, sideTxResult abci.SideTxResultType) sdk.Result {
					require.Equal(t, abci.SideTxResultType_Yes, sideTxResult, "Result `sideTxResult` should be `yes`")

					// try to set random state data to store
					happ.SidechannelKeeper.SetTx(ctx, 900, testTxStateData1)

					// error in post-tx handler
					return sdk.Result{
						Code:      sdk.CodeInternal,
						Codespace: sdk.CodespaceRoot,
						Events:    ctx.EventManager().Events(),
					}
				},
			}
			router.AddRoute(routeMsgSideCounter, handler)
			happ.SetSideRouter(router)

			happ.SidechannelKeeper.SetTx(ctx, height-2, txBytes) // set tx in the store for process
			res := happ.BeginSideBlocker(ctx, req)
			require.Equal(t, 0, len(res.Events), "It should have 0 events")
			require.Nil(t, happ.SidechannelKeeper.GetTx(ctx, height-2, txHash), "Tx should not be present in store after begin block")

			// check if it saved the data
			require.Equal(t, 0, len(happ.SidechannelKeeper.GetTxs(ctx, 900)), "It shouldn't save state after failed post-tx execution")
		}

		// shouldn't save state on failed execution of post-tx handler
		{
			// context with new event manager
			ctx = ctx.WithEventManager(sdk.NewEventManager())

			// setup router and handler
			router := hmTypes.NewSideRouter()
			handler := &hmTypes.SideHandlers{
				SideTxHandler: func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
					return abci.ResponseDeliverSideTx{}
				},
				PostTxHandler: func(ctx sdk.Context, msg sdk.Msg, sideTxResult abci.SideTxResultType) sdk.Result {
					require.Equal(t, abci.SideTxResultType_Yes, sideTxResult, "Result `sideTxResult` should be `yes`")

					// try to set random state data to store
					happ.SidechannelKeeper.SetTx(ctx, 900, testTxStateData1)

					panic("Paniced in handler")
				},
			}
			router.AddRoute(routeMsgSideCounter, handler)
			happ.SetSideRouter(router)

			happ.SidechannelKeeper.SetTx(ctx, height-2, txBytes) // set tx in the store for process
			res := happ.BeginSideBlocker(ctx, req)
			require.Equal(t, 0, len(res.Events), "It should have 0 events")
			require.Nil(t, happ.SidechannelKeeper.GetTx(ctx, height-2, txHash), "Tx should not be present in store after begin block")

			// check if it saved the data
			require.Equal(t, 0, len(happ.SidechannelKeeper.GetTxs(ctx, 900)), "It shouldn't save state after failed post-tx execution")
		}
	})
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
