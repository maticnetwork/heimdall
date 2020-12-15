package app_test

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/cosmos/cosmos-sdk/store"
	testdata "github.com/cosmos/cosmos-sdk/testutil/testdata"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmprototypes "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/maticnetwork/heimdall/app"
	"github.com/maticnetwork/heimdall/testutil"
	hmtestdata "github.com/maticnetwork/heimdall/testutil/testdata"
	hmtypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/x/gov/types"
	sidechannelkeeper "github.com/maticnetwork/heimdall/x/sidechannel/keeper"
)

var testTxStateData1 = []byte("test-tx-state1")
var testTxStateData2 = []byte("test-tx-state2")

//
// Test suite
//

// SideTxProcessorTestSuite test suite for side-tx processor
type SideTxProcessorTestSuite struct {
	suite.Suite

	happ           *app.HeimdallApp
	keeper         sidechannelkeeper.Keeper
	ctx            sdk.Context
	encodingConfig client.TxConfig
}

func (suite *SideTxProcessorTestSuite) SetupTest() {
	encodingConfig := simapp.MakeTestEncodingConfig()
	encodingConfig.InterfaceRegistry.RegisterImplementations(&hmtestdata.Dog{})
	testdata.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	hmtestdata.RegisterInterfaces(encodingConfig.InterfaceRegistry)

	// set encodingConfig
	suite.encodingConfig = encodingConfig.TxConfig
	// set keeper and context
	suite.ctx, suite.keeper = setupKeeper(suite.T())

	// create app wrapper
	suite.happ = &app.HeimdallApp{
		SidechannelKeeper: suite.keeper,
		BaseApp:           baseapp.NewBaseApp("test", testutil.Logger(suite.T()), nil, encodingConfig.TxConfig.TxDecoder()),
	}
	suite.happ.SetTxDecoder(encodingConfig.TxConfig.TxDecoder())
}

func TestSideTxProcessorTestSuite(t *testing.T) {
	suite.Run(t, new(SideTxProcessorTestSuite))
}

func (suite *SideTxProcessorTestSuite) getTx() (tmtypes.Tx, sdk.Tx) {
	t, encodingConfig := suite.T(), suite.encodingConfig

	txBuilder := encodingConfig.NewTxBuilder()
	msg := hmtestdata.NewServiceSideMsgCreateDog(&hmtestdata.SideMsgCreateDog{Dog: &hmtestdata.Dog{Name: "Spot"}})
	err := txBuilder.SetMsgs(msg)
	require.Nil(t, err, "It should be no error while setting msg")
	txBytes, err := encodingConfig.TxEncoder()(txBuilder.GetTx())
	require.Nil(t, err, "It should be no error while encoding tx")

	tx, err := encodingConfig.TxDecoder()(txBytes)
	require.Nil(t, err, "It should throw no error while decoding tx")

	return txBytes, tx
}

//
// Test cases
//

func (suite *SideTxProcessorTestSuite) TestPostDeliverTxHandler() {
	t, keeper, ctx, happ, encodingConfig := suite.T(), suite.keeper, suite.ctx, suite.happ, suite.encodingConfig

	{
		txBuilder := encodingConfig.NewTxBuilder()
		msg := testdata.NewServiceMsgCreateDog(&testdata.MsgCreateDog{Dog: &testdata.Dog{Name: "Spot"}})
		err := txBuilder.SetMsgs(msg)
		require.Nil(t, err, "It should throw no error while setting msg")
		txBytes, err := encodingConfig.TxEncoder()(txBuilder.GetTx())
		require.Nil(t, err, "It should throw no error while encoding tx")

		tx, err := encodingConfig.TxDecoder()(txBytes)
		require.Nil(t, err, "It should throw no error while decoding tx")

		// test post deliver with height 3
		ctx = ctx.WithBlockHeight(3).WithTxBytes(txBytes)
		happ.PostDeliverTxHandler(ctx, tx, &sdk.Result{})
		resultTxBytes := keeper.GetTx(ctx, 3, tmtypes.Tx(txBytes).Hash())
		require.Nil(t, resultTxBytes, "Stored tx should be nil")
		require.Zero(t, len(resultTxBytes), "Tx bytes should be stored in keeper")

		// test post deliver with height 20 without side-tx
		ctx = ctx.WithBlockHeight(20).WithTxBytes(txBytes)
		happ.PostDeliverTxHandler(ctx, tx, &sdk.Result{})
		resultTxBytes = keeper.GetTx(ctx, 20, tmtypes.Tx(txBytes).Hash())
		require.Nil(t, resultTxBytes, "Stored tx should be nil")
		require.Zero(t, len(resultTxBytes), "Tx bytes should be stored in keeper")
	}

	// test post deliver with height 20 with side-tx
	{
		txBytes, tx := suite.getTx()

		ctx = ctx.WithBlockHeight(20).WithTxBytes(txBytes)
		happ.PostDeliverTxHandler(ctx, tx, &sdk.Result{})
		resultTxBytes := keeper.GetTx(ctx, 20, tmtypes.Tx(txBytes).Hash())
		require.NotNil(t, resultTxBytes, "Stored tx bytes shouldn't be nil")
		require.Greater(t, len(resultTxBytes), 0)
		require.True(t, bytes.Equal(txBytes, resultTxBytes), "Stored tx bytes should same as actual tx bytes")
	}
}

func (suite *SideTxProcessorTestSuite) TestDeliverSideTxHandler() {
	t, keeper, ctx, happ := suite.T(), suite.keeper, suite.ctx, suite.happ

	txBytes, tx := suite.getTx()

	// deliver side thandler
	require.Panics(t, func() {
		happ.DeliverSideTxHandler(ctx, tx, abci.RequestDeliverSideTx{
			Tx: txBytes,
		})
	}, "Tx without side-tx handler should panic")

	// Set router (testing only)
	router := hmtypes.NewSideRouter()
	msg := tx.GetMsgs()[0]
	router.AddRoute(msg.Route(), &hmtypes.SideHandlers{
		SideTxHandler: func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
			return abci.ResponseDeliverSideTx{
				Result: tmprototypes.SideTxResultType_YES,
			}
		},
		PostTxHandler: func(ctx sdk.Context, msg sdk.Msg, sideTxResult tmprototypes.SideTxResultType) (*sdk.Result, error) {
			return &sdk.Result{}, nil
		},
	})
	happ.SetSideRouter(router)

	t.Run("YesVote", func(t *testing.T) {
		var res abci.ResponseDeliverSideTx
		// test route
		require.NotPanics(t, func() {
			res = happ.DeliverSideTxHandler(ctx, tx, abci.RequestDeliverSideTx{
				Tx: txBytes,
			})
		})

		require.Equal(t, tmprototypes.SideTxResultType_YES, res.GetResult(), "Result from deliver side-tx should be vote `Yes`")
		if m, ok := app.IsSideMsg(msg); ok {
			require.Equal(t, len(m.GetSideSignBytes()), len(res.Data), "Result from deliver side-tx data should be same as expected")
		}
	})

	t.Run("NoVote", func(t *testing.T) {
		resultData := []byte("hello")
		router := hmtypes.NewSideRouter()
		router.AddRoute(msg.Route(), &hmtypes.SideHandlers{
			SideTxHandler: func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
				return abci.ResponseDeliverSideTx{
					Result: tmprototypes.SideTxResultType_NO,
					Data:   resultData,
				}
			},
			PostTxHandler: func(ctx sdk.Context, msg sdk.Msg, sideTxResult tmprototypes.SideTxResultType) (*sdk.Result, error) {
				return &sdk.Result{}, nil
			},
		})
		happ.SetSideRouter(router)

		res := happ.DeliverSideTxHandler(ctx, tx, abci.RequestDeliverSideTx{
			Tx: txBytes,
		})

		require.Equal(t, tmprototypes.SideTxResultType_NO, res.GetResult(), "Result from deliver side-tx should be vote `No`")
		require.Equal(t, resultData, res.Data, "Result from deliver side-tx data should be same as expected")
	})

	t.Run("SkipVote", func(t *testing.T) {
		router := hmtypes.NewSideRouter()
		router.AddRoute(msg.Route(), &hmtypes.SideHandlers{
			SideTxHandler: func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
				return abci.ResponseDeliverSideTx{
					Code: uint32(1),
				}
			},
			PostTxHandler: func(ctx sdk.Context, msg sdk.Msg, sideTxResult tmprototypes.SideTxResultType) (*sdk.Result, error) {
				return &sdk.Result{}, nil
			},
		})
		happ.SetSideRouter(router)
		res := happ.DeliverSideTxHandler(ctx, tx, abci.RequestDeliverSideTx{
			Tx: tmtypes.Tx(txBytes),
		})

		require.Equal(t, tmprototypes.SideTxResultType_SKIP, res.GetResult(), "Result from deliver side-tx should be vote `Skip` if result is not OK")
	})

	t.Run("State", func(t *testing.T) {
		// testing by storing random txs to store
		keeper.SetTx(ctx, 800, testTxStateData1)
		require.Equal(t, 1, len(keeper.GetTxs(ctx, 800)), "It should set state in storage in normal case")

		router := hmtypes.NewSideRouter()
		router.AddRoute(msg.Route(), &hmtypes.SideHandlers{
			SideTxHandler: func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
				// set random tx in deliver side-tx
				// it should reset the storage (shouldn't change the state)
				keeper.SetTx(ctx, 800, testTxStateData2)
				require.Equal(t, 2, len(keeper.GetTxs(ctx, 800)), "It should set state in storage temperory in side-tx handler")

				return abci.ResponseDeliverSideTx{
					Result: tmprototypes.SideTxResultType_YES,
				}
			},
			PostTxHandler: func(ctx sdk.Context, msg sdk.Msg, sideTxResult tmprototypes.SideTxResultType) (*sdk.Result, error) {
				return &sdk.Result{}, nil
			},
		})
		happ.SetSideRouter(router)
		res := happ.DeliverSideTxHandler(ctx, tx, abci.RequestDeliverSideTx{
			Tx: tmtypes.Tx(txBytes),
		})

		require.Equal(t, tmprototypes.SideTxResultType_YES, res.GetResult(), "It should return `yes` vote from deliver side-tx")
		// testing by fetching txs from store
		require.Equal(t, 1, len(keeper.GetTxs(ctx, 800)), "It shouldn't change state in deliver side-tx")
	})

	t.Run("Panic", func(t *testing.T) {
		router := hmtypes.NewSideRouter()
		router.AddRoute(msg.Route(), &hmtypes.SideHandlers{
			SideTxHandler: func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
				panic("TestSideTxPanic")
			},
			PostTxHandler: func(ctx sdk.Context, msg sdk.Msg, sideTxResult tmprototypes.SideTxResultType) (*sdk.Result, error) {
				return &sdk.Result{}, nil
			},
		})
		happ.SetSideRouter(router)
		require.Panics(t, func() {
			happ.DeliverSideTxHandler(ctx, tx, abci.RequestDeliverSideTx{
				Tx: txBytes,
			})
		}, "It should panic because of side-tx handler panic")

		var res abci.ResponseDeliverSideTx
		require.NotPanics(t, func() {
			res = happ.DeliverSideTx(abci.RequestDeliverSideTx{
				Tx: txBytes,
			})
		}, "It shouldn't panic due to recover in app.DeliverSideTx")

		require.Equal(t, tmprototypes.SideTxResultType_SKIP, res.GetResult(), "It should return `skip` vote due to panic")
	})
}

func (suite *SideTxProcessorTestSuite) TestBeginSideBlocker() {
	t, keeper, ctx, happ := suite.T(), suite.keeper, suite.ctx, suite.happ

	// get tx and txbytes
	txBytes, tx := suite.getTx()
	msg := tx.GetMsgs()[0]
	txHash := tmtypes.Tx(txBytes).Hash()

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
		var height uint64 = 20
		ctx = ctx.WithBlockHeight(int64(height))

		addr1 := []byte("hello-1")
		addr2 := []byte("hello-2")
		addr3 := []byte("hello-3")
		addr4 := []byte("hello-4")
		// set validators
		err := keeper.SetValidators(ctx, height, []*abci.Validator{
			{Address: addr1, Power: 10},
			{Address: addr2, Power: 20},
			{Address: addr3, Power: 30},
			{Address: addr4, Power: 40},
		})
		require.Nil(t, err, "It should throw no error while setting validators")

		res := happ.BeginSideBlocker(ctx, abci.RequestBeginSideBlock{})
		require.Equal(t, 0, len(res.Events), "Events from begin side blocker result should be empty with validators but no sigs")

		// test all types of result
		for key, value := range tmprototypes.SideTxResultType_value {
			t.Run(key, func(t *testing.T) {
				// setup router and handler
				router := hmtypes.NewSideRouter()
				handler := &hmtypes.SideHandlers{
					SideTxHandler: func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
						return abci.ResponseDeliverSideTx{}
					},
					PostTxHandler: func(ctx sdk.Context, msg sdk.Msg, sideTxResult tmprototypes.SideTxResultType) (*sdk.Result, error) {
						return &sdk.Result{}, nil
					},
				}
				router.AddRoute(msg.Route(), handler)
				happ.SetSideRouter(router)

				keeper.SetTx(ctx, height-2, txBytes) // set tx in the store for process
				res = happ.BeginSideBlocker(ctx, abci.RequestBeginSideBlock{
					SideTxResults: []tmprototypes.SideTxResult{
						{
							TxHash: txHash,
							Sigs: []tmprototypes.SideTxSig{
								{
									Result:  tmprototypes.SideTxResultType(value),
									Address: addr1,
								},
								{
									Result:  tmprototypes.SideTxResultType(value),
									Address: addr2,
								},
								{
									Result:  tmprototypes.SideTxResultType(value),
									Address: addr3,
								},
								{
									Result:  tmprototypes.SideTxResultType(value),
									Address: addr4,
								},
							},
						},
					},
				})
				require.Equal(t, 0, len(res.Events), "It should have no event")
				require.Nil(t, keeper.GetTx(ctx, height-2, txHash), "Tx should not be present in store after begin block")
			})
		}

		t.Run("SkipWithoutPowerCalculation", func(t *testing.T) {
			// skip without any power calculation
			keeper.SetTx(ctx, height-2, txBytes)
			res = happ.BeginSideBlocker(ctx, abci.RequestBeginSideBlock{})
			require.Equal(t, 0, len(res.Events), "It should have no event with validators")
			require.Nil(t, keeper.GetTx(ctx, height-2, txHash), "Tx should not be present in store after begin block")
		})
	})

	t.Run("State", func(t *testing.T) {
		var height uint64 = 20
		ctx = ctx.WithBlockHeight(int64(height))

		addr1 := []byte("hello-1")
		addr2 := []byte("hello-2")
		addr3 := []byte("hello-3")
		addr4 := []byte("hello-4")
		// set validators
		err := keeper.SetValidators(ctx, height, []*abci.Validator{
			{Address: addr1, Power: 10},
			{Address: addr2, Power: 20},
			{Address: addr3, Power: 30},
			{Address: addr4, Power: 40},
		})
		require.Nil(t, err, "It should throw no error while setting validators")

		req := abci.RequestBeginSideBlock{
			SideTxResults: []tmprototypes.SideTxResult{
				{
					TxHash: txHash,
					Sigs: []tmprototypes.SideTxSig{
						{
							Result:  tmprototypes.SideTxResultType_NO,
							Address: addr1,
						},
						{
							Result:  tmprototypes.SideTxResultType_NO,
							Address: addr2,
						},
						{
							Result:  tmprototypes.SideTxResultType_YES,
							Address: addr3,
						},
						{
							Result:  tmprototypes.SideTxResultType_YES,
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
			router := hmtypes.NewSideRouter()
			handler := &hmtypes.SideHandlers{
				SideTxHandler: func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
					return abci.ResponseDeliverSideTx{}
				},
				PostTxHandler: func(ctx sdk.Context, msg sdk.Msg, sideTxResult tmprototypes.SideTxResultType) (*sdk.Result, error) {
					require.Equal(t, tmprototypes.SideTxResultType_YES, sideTxResult, "Result `sideTxResult` should be `yes`")

					// try to set random state data to store
					keeper.SetTx(ctx, 700, testTxStateData1)

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

					return &sdk.Result{
						Events: ctx.EventManager().ABCIEvents(),
					}, nil
				},
			}
			router.AddRoute(msg.Route(), handler)
			happ.SetSideRouter(router)

			keeper.SetTx(ctx, height-2, txBytes) // set tx in the store for process
			res := happ.BeginSideBlocker(ctx, req)
			require.Equal(t, 2, len(res.Events), "It should include correct emitted events")
			require.Nil(t, keeper.GetTx(ctx, height-2, txHash), "Tx should not be present in store after begin block")

			// check if it saved the data
			require.Equal(t, 1, len(keeper.GetTxs(ctx, 700)), "It should save state correctly after successful post-tx execution")
		}

		// shouldn't save state on failed execution of post-tx handler
		{
			// context with new event manager
			ctx = ctx.WithEventManager(sdk.NewEventManager())

			// setup router and handler
			router := hmtypes.NewSideRouter()
			handler := &hmtypes.SideHandlers{
				SideTxHandler: func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
					return abci.ResponseDeliverSideTx{}
				},
				PostTxHandler: func(ctx sdk.Context, msg sdk.Msg, sideTxResult tmprototypes.SideTxResultType) (*sdk.Result, error) {
					require.Equal(t, tmprototypes.SideTxResultType_YES, sideTxResult, "Result `sideTxResult` should be `yes`")

					// try to set random state data to store
					keeper.SetTx(ctx, 900, testTxStateData1)

					// error in post-tx handler
					return nil, errors.New("Error example")
				},
			}
			router.AddRoute(msg.Route(), handler)
			happ.SetSideRouter(router)

			keeper.SetTx(ctx, height-2, txBytes) // set tx in the store for process
			res := happ.BeginSideBlocker(ctx, req)
			require.Equal(t, 0, len(res.Events), "It should have 0 events")
			require.Nil(t, keeper.GetTx(ctx, height-2, txHash), "Tx should not be present in store after begin block")

			// check if it saved the data
			require.Equal(t, 0, len(keeper.GetTxs(ctx, 900)), "It shouldn't save state after failed post-tx execution")
		}

		// shouldn't save state on failed execution of post-tx handler
		{
			// context with new event manager
			ctx = ctx.WithEventManager(sdk.NewEventManager())

			// setup router and handler
			router := hmtypes.NewSideRouter()
			handler := &hmtypes.SideHandlers{
				SideTxHandler: func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
					return abci.ResponseDeliverSideTx{}
				},
				PostTxHandler: func(ctx sdk.Context, msg sdk.Msg, sideTxResult tmprototypes.SideTxResultType) (*sdk.Result, error) {
					require.Equal(t, tmprototypes.SideTxResultType_YES, sideTxResult, "Result `sideTxResult` should be `yes`")

					// try to set random state data to store
					keeper.SetTx(ctx, 900, testTxStateData1)

					panic("Paniced in handler")
				},
			}
			router.AddRoute(msg.Route(), handler)
			happ.SetSideRouter(router)

			keeper.SetTx(ctx, height-2, txBytes) // set tx in the store for process
			res := happ.BeginSideBlocker(ctx, req)
			require.Equal(t, 0, len(res.Events), "It should have 0 events")
			require.Nil(t, keeper.GetTx(ctx, height-2, txHash), "Tx should not be present in store after begin block")

			// check if it saved the data
			require.Equal(t, 0, len(keeper.GetTxs(ctx, 900)), "It shouldn't save state after failed post-tx execution")
		}
	})
}

//
// Internal setup keeper
//

func setupKeeper(t *testing.T) (sdk.Context, sidechannelkeeper.Keeper) {
	t.Helper()
	key := sdk.NewKVStoreKey(types.StoreKey)
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.NoError(t, err)
	ctx := sdk.NewContext(ms, tmproto.Header{Time: time.Unix(0, 0)}, false, testutil.Logger(t))
	return ctx, sidechannelkeeper.NewKeeper(types.ModuleCdc, key)
}
