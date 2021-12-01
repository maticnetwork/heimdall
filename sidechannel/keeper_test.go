package sidechannel_test

import (
	"errors"
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/app"
)

//
// Test suite
//

// KeeperTestSuite integrate test suite context object
type KeeperTestSuite struct {
	suite.Suite

	app *app.HeimdallApp
	ctx sdk.Context
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app, suite.ctx = createTestApp(false)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

//
// Tests
//

func (suite *KeeperTestSuite) TestTx() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	tx1 := tmTypes.Tx([]byte("transaction-1"))
	tx2 := tmTypes.Tx([]byte("transaction-2"))
	tx3 := tmTypes.Tx([]byte("transaction-3"))
	tx4 := tmTypes.Tx([]byte("transaction-4"))

	var height1 int64 = 10
	var height2 int64 = 24

	t.Run("SetTx", func(t *testing.T) {
		app.SidechannelKeeper.SetTx(ctx, height1, tx1)
		app.SidechannelKeeper.SetTx(ctx, height1, tx2)

		rtx1 := app.SidechannelKeeper.GetTx(ctx, height1, tx1.Hash())
		require.Equal(t, tx1, rtx1)

		rtx2 := app.SidechannelKeeper.GetTx(ctx, height1, tx2.Hash())
		require.Equal(t, tx2, rtx2)

		rtx3 := app.SidechannelKeeper.GetTx(ctx, height2, tx1.Hash())
		require.Nil(t, rtx3)
	})

	t.Run("HasTx", func(t *testing.T) {
		require.False(t, app.SidechannelKeeper.HasTx(ctx, height2, tx1.Hash()))

		app.SidechannelKeeper.SetTx(ctx, height2, tx1)
		app.SidechannelKeeper.SetTx(ctx, height2, tx2)
		app.SidechannelKeeper.SetTx(ctx, height2, tx3)
		app.SidechannelKeeper.SetTx(ctx, height2, tx4)

		require.True(t, app.SidechannelKeeper.HasTx(ctx, height2, tx1.Hash()))

		rtx1 := app.SidechannelKeeper.GetTx(ctx, height2, tx1.Hash())
		require.Equal(t, tx1, rtx1)

		rtx2 := app.SidechannelKeeper.GetTx(ctx, height2, tx2.Hash())
		require.Equal(t, tx2, rtx2)

		rtx3 := app.SidechannelKeeper.GetTx(ctx, int64(3), tx1.Hash())
		require.Nil(t, rtx3)

		require.False(t, app.SidechannelKeeper.HasTx(ctx, height1, tx3.Hash()))
		require.True(t, app.SidechannelKeeper.HasTx(ctx, height2, tx1.Hash()))
	})

	t.Run("GetTxs", func(t *testing.T) {
		txs := app.SidechannelKeeper.GetTxs(ctx, height2)
		require.Equal(t, 4, len(txs))
		require.Less(t, -1, txs.Index(tx1))
		require.Less(t, -1, txs.Index(tx2))
		require.Less(t, -1, txs.Index(tx3))
		require.Less(t, -1, txs.Index(tx4))

		txs = app.SidechannelKeeper.GetTxs(ctx, height1)
		require.Equal(t, 2, len(txs))
		require.Equal(t, -1, txs.Index(tx3))
		require.Less(t, -1, txs.Index(tx1))
		require.Less(t, -1, txs.Index(tx2))
	})

	t.Run("IterateTxAndApplyFn", func(t *testing.T) {
		txs := make(tmTypes.Txs, 0)
		app.SidechannelKeeper.IterateTxAndApplyFn(ctx, height2, func(tx tmTypes.Tx) error {
			txs = append(txs, tx)
			return errors.New("Retrieve only one")
		})
		require.Equal(t, 1, len(txs))
	})

	t.Run("IterateTxsAndApplyFn", func(t *testing.T) {
		txs := make(tmTypes.Txs, 0)
		heightMapping := make(map[int64]bool)
		app.SidechannelKeeper.IterateTxsAndApplyFn(ctx, func(height int64, tx tmTypes.Tx) error {
			txs = append(txs, tx)
			heightMapping[height] = true
			return nil
		})
		require.Equal(t, 6, len(txs))
		require.Equal(t, 2, len(heightMapping))
		require.True(t, heightMapping[height1])
		require.True(t, heightMapping[height2])
		require.False(t, heightMapping[17])

		txs = make(tmTypes.Txs, 0)
		app.SidechannelKeeper.IterateTxsAndApplyFn(ctx, func(height int64, tx tmTypes.Tx) error {
			txs = append(txs, tx)
			return errors.New("Only one tx")
		})
		require.Equal(t, 1, len(txs))
	})

	t.Run("RemoveTx", func(t *testing.T) {
		app.SidechannelKeeper.RemoveTx(ctx, height2, tx4.Hash())
		txs := app.SidechannelKeeper.GetTxs(ctx, height2)
		require.Equal(t, 3, len(txs))
		require.False(t, app.SidechannelKeeper.HasTx(ctx, height2, tx4.Hash()))
	})
}

func (suite *KeeperTestSuite) TestValidators() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	var height1 int64 = 10
	var height2 int64 = 24

	validators := make([]abci.Validator, 10)
	for i := 0; i < 10; i++ {
		validators[i] = abci.Validator{
			Address: []byte("address" + strconv.Itoa(i)),
			Power:   int64(i+1) * 1000,
		}
	}

	require.Equal(t, 10, len(validators))

	t.Run("SetValidators", func(t *testing.T) {
		err := app.SidechannelKeeper.SetValidators(ctx, height1, validators)
		require.NoError(t, err)

		err = app.SidechannelKeeper.SetValidators(ctx, height2, validators[:5])
		require.NoError(t, err)
	})

	t.Run("GetValidators", func(t *testing.T) {
		vals := app.SidechannelKeeper.GetValidators(ctx, height1)
		require.Equal(t, 10, len(vals))

		vals = app.SidechannelKeeper.GetValidators(ctx, height2)
		require.Equal(t, 5, len(vals))
		require.Equal(t, validators[0].Address, vals[0].Address)
		require.Equal(t, validators[0].Power, vals[0].Power)
	})

	t.Run("HasValidators", func(t *testing.T) {
		require.True(t, app.SidechannelKeeper.HasValidators(ctx, height1))
		require.True(t, app.SidechannelKeeper.HasValidators(ctx, height2))
		require.False(t, app.SidechannelKeeper.HasValidators(ctx, int64(5)))
	})

	t.Run("IterateValidatorsAndApplyFn", func(t *testing.T) {
		validators := make([]abci.Validator, 0)
		heightMapping := make(map[int64]bool)
		app.SidechannelKeeper.IterateValidatorsAndApplyFn(ctx, func(height int64, vs []abci.Validator) error {
			validators = append(validators, vs...)
			heightMapping[height] = true
			return nil
		})

		require.Equal(t, 15, len(validators))
		require.Equal(t, 2, len(heightMapping))
		require.True(t, heightMapping[height1])
		require.True(t, heightMapping[height2])
		require.False(t, heightMapping[17])

		validators = make([]abci.Validator, 0)
		app.SidechannelKeeper.IterateValidatorsAndApplyFn(ctx, func(height int64, vs []abci.Validator) error {
			validators = append(validators, vs...)
			return errors.New("Only with height")
		})
		require.Equal(t, 10, len(validators))
	})

	t.Run("RemoveValidators", func(t *testing.T) {
		app.SidechannelKeeper.RemoveValidators(ctx, height1)
		require.False(t, app.SidechannelKeeper.HasValidators(ctx, height1))
		require.True(t, app.SidechannelKeeper.HasValidators(ctx, height2))

		validators := app.SidechannelKeeper.GetValidators(ctx, height2)
		require.Equal(t, 5, len(validators))

		validators = app.SidechannelKeeper.GetValidators(ctx, height1)
		require.Equal(t, 0, len(validators))
	})
}

func (suite *KeeperTestSuite) TestLogger() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	logger := app.SidechannelKeeper.Logger(ctx)
	require.NotNil(t, logger)
}

func (suite *KeeperTestSuite) TestCodespace() {
	t, app, _ := suite.T(), suite.app, suite.ctx

	codespace := app.SidechannelKeeper.Codespace()
	require.NotEmpty(t, codespace)
}
