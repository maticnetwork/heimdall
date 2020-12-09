package keeper_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/maticnetwork/heimdall/testutil"
	"github.com/maticnetwork/heimdall/x/sidechannel/keeper"
	"github.com/maticnetwork/heimdall/x/sidechannel/types"
)

//
// Test suite
//

// KeeperTestSuite integrate test suite context object
type KeeperTestSuite struct {
	suite.Suite

	keeper keeper.Keeper
	ctx    sdk.Context
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.ctx, suite.keeper = setupKeeper(suite.T())
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

//
// Tests
//

func (suite *KeeperTestSuite) TestTx() {
	t, k, ctx := suite.T(), suite.keeper, suite.ctx

	tx1 := tmtypes.Tx([]byte("transaction-1"))
	tx2 := tmtypes.Tx([]byte("transaction-2"))
	tx3 := tmtypes.Tx([]byte("transaction-3"))
	tx4 := tmtypes.Tx([]byte("transaction-4"))

	var height1 uint64 = 10
	var height2 uint64 = 24

	t.Run("SetTx", func(t *testing.T) {
		k.SetTx(ctx, height1, tx1)
		k.SetTx(ctx, height1, tx2)

		rtx1 := k.GetTx(ctx, height1, tx1.Hash())
		require.Equal(t, tx1, rtx1)

		rtx2 := k.GetTx(ctx, height1, tx2.Hash())
		require.Equal(t, tx2, rtx2)

		rtx3 := k.GetTx(ctx, height2, tx1.Hash())
		require.Nil(t, rtx3)
	})

	t.Run("HasTx", func(t *testing.T) {
		require.False(t, k.HasTx(ctx, height2, tx1.Hash()))

		k.SetTx(ctx, height2, tx1)
		k.SetTx(ctx, height2, tx2)
		k.SetTx(ctx, height2, tx3)
		k.SetTx(ctx, height2, tx4)

		require.True(t, k.HasTx(ctx, height2, tx1.Hash()))

		rtx1 := k.GetTx(ctx, height2, tx1.Hash())
		require.Equal(t, tx1, rtx1)

		rtx2 := k.GetTx(ctx, height2, tx2.Hash())
		require.Equal(t, tx2, rtx2)

		rtx3 := k.GetTx(ctx, uint64(3), tx1.Hash())
		require.Nil(t, rtx3)

		require.False(t, k.HasTx(ctx, height1, tx3.Hash()))
		require.True(t, k.HasTx(ctx, height2, tx1.Hash()))
	})

	t.Run("GetTxs", func(t *testing.T) {
		txs := k.GetTxs(ctx, height2)
		require.Equal(t, 4, len(txs))
		require.Less(t, -1, txs.Index(tx1))
		require.Less(t, -1, txs.Index(tx2))
		require.Less(t, -1, txs.Index(tx3))
		require.Less(t, -1, txs.Index(tx4))

		txs = k.GetTxs(ctx, height1)
		require.Equal(t, 2, len(txs))
		require.Equal(t, -1, txs.Index(tx3))
		require.Less(t, -1, txs.Index(tx1))
		require.Less(t, -1, txs.Index(tx2))
	})

	t.Run("IterateTxAndApplyFn", func(t *testing.T) {
		txs := make(tmtypes.Txs, 0)
		k.IterateTxAndApplyFn(ctx, height2, func(tx tmtypes.Tx) error {
			txs = append(txs, tx)
			return errors.New("Retrieve only one")
		})
		require.Equal(t, 1, len(txs))
	})

	t.Run("IterateTxsAndApplyFn", func(t *testing.T) {
		txs := make(tmtypes.Txs, 0)
		heightMapping := make(map[uint64]bool)
		k.IterateTxsAndApplyFn(ctx, func(height uint64, tx tmtypes.Tx) error {
			txs = append(txs, tx)
			heightMapping[height] = true
			return nil
		})
		require.Equal(t, 6, len(txs))
		require.Equal(t, 2, len(heightMapping))
		require.True(t, heightMapping[height1])
		require.True(t, heightMapping[height2])
		require.False(t, heightMapping[17])

		txs = make(tmtypes.Txs, 0)
		k.IterateTxsAndApplyFn(ctx, func(height uint64, tx tmtypes.Tx) error {
			txs = append(txs, tx)
			return errors.New("Only one tx")
		})
		require.Equal(t, 1, len(txs))
	})

	t.Run("RemoveTx", func(t *testing.T) {
		k.RemoveTx(ctx, height2, tx4.Hash())
		txs := k.GetTxs(ctx, height2)
		require.Equal(t, 3, len(txs))
		require.False(t, k.HasTx(ctx, height2, tx4.Hash()))
	})
}

func (suite *KeeperTestSuite) TestValidators() {
	t, k, ctx := suite.T(), suite.keeper, suite.ctx

	var height1 uint64 = 10
	var height2 uint64 = 24

	validators := make([]*abci.Validator, 10)
	for i := 0; i < 10; i++ {
		validators[i] = &abci.Validator{
			Address: []byte(fmt.Sprintf("address %d", i)),
			Power:   int64(i+1) * 1000,
		}
	}

	require.Equal(t, 10, len(validators))

	t.Run("SetValidators", func(t *testing.T) {
		err := k.SetValidators(ctx, height1, validators)
		require.NoError(t, err)

		err = k.SetValidators(ctx, height2, validators[:5])
		require.NoError(t, err)
	})

	t.Run("GetValidators", func(t *testing.T) {
		vals := k.GetValidators(ctx, height1)
		require.Equal(t, 10, len(vals))

		vals = k.GetValidators(ctx, height2)
		require.Equal(t, 5, len(vals))
		require.Equal(t, validators[0].Address, vals[0].Address)
		require.Equal(t, validators[0].Power, vals[0].Power)
	})

	t.Run("HasValidators", func(t *testing.T) {
		require.True(t, k.HasValidators(ctx, height1))
		require.True(t, k.HasValidators(ctx, height2))
		require.False(t, k.HasValidators(ctx, uint64(5)))
	})

	t.Run("IterateValidatorsAndApplyFn", func(t *testing.T) {
		validators := make([]*abci.Validator, 0)
		heightMapping := make(map[uint64]bool)
		k.IterateValidatorsAndApplyFn(ctx, func(height uint64, vs []*abci.Validator) error {
			validators = append(validators, vs...)
			heightMapping[height] = true
			return nil
		})

		require.Equal(t, 15, len(validators))
		require.Equal(t, 2, len(heightMapping))
		require.True(t, heightMapping[height1])
		require.True(t, heightMapping[height2])
		require.False(t, heightMapping[17])

		validators = make([]*abci.Validator, 0)
		k.IterateValidatorsAndApplyFn(ctx, func(height uint64, vs []*abci.Validator) error {
			validators = append(validators, vs...)
			return errors.New("Only with height")
		})
		require.Equal(t, 10, len(validators))
	})

	t.Run("RemoveValidators", func(t *testing.T) {
		k.RemoveValidators(ctx, height1)
		require.False(t, k.HasValidators(ctx, height1))
		require.True(t, k.HasValidators(ctx, height2))

		validators := k.GetValidators(ctx, height2)
		require.Equal(t, 5, len(validators))

		validators = k.GetValidators(ctx, height1)
		require.Equal(t, 0, len(validators))
	})
}

func (suite *KeeperTestSuite) TestLogger() {
	t, k, ctx := suite.T(), suite.keeper, suite.ctx

	logger := k.Logger(ctx)
	require.NotNil(t, logger)
}

//
// Internal setup keeper
//

func setupKeeper(t *testing.T) (sdk.Context, keeper.Keeper) {
	t.Helper()
	key := sdk.NewKVStoreKey(types.StoreKey)
	db := dbm.NewMemDB()
	ms := store.NewCommitMultiStore(db)
	ms.MountStoreWithDB(key, sdk.StoreTypeIAVL, db)
	err := ms.LoadLatestVersion()
	require.NoError(t, err)
	ctx := sdk.NewContext(ms, tmproto.Header{Time: time.Unix(0, 0)}, false, testutil.Logger(t))
	return ctx, keeper.NewKeeper(types.ModuleCdc, key)
}
