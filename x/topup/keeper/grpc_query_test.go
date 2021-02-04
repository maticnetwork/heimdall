package keeper_test

import (
	"math/big"
	"testing"

	ethTypes "github.com/maticnetwork/bor/core/types"

	"github.com/stretchr/testify/require"

	hmTypes "github.com/maticnetwork/heimdall/types/common"

	"github.com/maticnetwork/heimdall/x/topup/types"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/x/topup/keeper"
)

func (suite *KeeperTestSuite) TestSequence() {
	t, initApp, ctx := suite.T(), suite.app, suite.ctx
	k := keeper.NewQueryServerImpl(initApp.TopupKeeper, &suite.contractCaller)

	logIndex := uint64(10)
	txHash := hmTypes.HexToHeimdallHash("123123123")
	chainParams := initApp.ChainKeeper.GetParams(ctx)
	blockNumber := big.NewInt(10)
	txReceipt := &ethTypes.Receipt{
		BlockNumber: blockNumber,
	}

	suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)

	// check if incoming tx is older
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(logIndex))

	t.Run("Without Sequence Setup ", func(t *testing.T) {
		resp, err := k.Sequence(sdk.WrapSDKContext(ctx), &types.QuerySequenceRequest{
			TxHash:   txHash.EthHash().String(),
			LogIndex: logIndex,
		})

		require.Nil(t, resp)
		require.NoError(t, err)
	})

	t.Run("With Sequence Setup ", func(t *testing.T) {
		initApp.TopupKeeper.SetTopupSequence(ctx, sequence.String())
		resp, err := k.Sequence(sdk.WrapSDKContext(ctx), &types.QuerySequenceRequest{
			TxHash:   txHash.EthHash().String(),
			LogIndex: logIndex,
		})

		require.NotNil(t, resp)
		require.NoError(t, err)
		require.Equal(t, resp.Sequence, sequence.Uint64())
	})
}
