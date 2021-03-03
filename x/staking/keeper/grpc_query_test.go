package keeper_test

import (
	"math/big"
	"math/rand"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethTypes "github.com/maticnetwork/bor/core/types"
	hmTypes "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/types/simulation"
	checkPointSim "github.com/maticnetwork/heimdall/x/checkpoint/simulation"
	"github.com/maticnetwork/heimdall/x/staking/keeper"
	"github.com/maticnetwork/heimdall/x/staking/types"
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestQueryValidator() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	k := keeper.Querier{
		Keeper: app.StakingKeeper,
	}

	// loading the validators
	checkPointSim.LoadValidatorSet(4, t, k.Keeper, ctx, false, 10)

	// getting the validators set
	validators := app.StakingKeeper.GetValidatorSet(ctx)
	validator := validators.Validators[0]

	queryValidatorInfo, err := k.Validator(sdk.WrapSDKContext(ctx), &types.QueryValidatorRequest{ValidatorId: int32(validator.ID)})

	require.NoError(t, err)
	require.NotNil(t, queryValidatorInfo)
	require.Equal(t, queryValidatorInfo.Validator.ID, validator.ID)
}

func (suite *KeeperTestSuite) TestQueryValidatorSet() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	k := keeper.Querier{
		Keeper: app.StakingKeeper,
	}

	// loading the validators
	checkPointSim.LoadValidatorSet(4, t, k.Keeper, ctx, false, 10)

	vaSet, err := k.ValidatorSet(sdk.WrapSDKContext(ctx), &types.QueryValidatorSetRequest{})
	require.NoError(t, err)
	require.NotNil(t, vaSet)
}

func (suite *KeeperTestSuite) TestQueryIsOlxTxStaking() {
	t, app, ctx := suite.T(), suite.app, suite.ctx
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	k := keeper.NewQueryServerImpl(app.StakingKeeper, suite.contractCaller)

	txHash := hmTypes.HexToHeimdallHash("123")
	chainParams := app.ChainKeeper.GetParams(ctx)
	txReceipt := &ethTypes.Receipt{BlockNumber: big.NewInt(10)}
	logIndex := uint64(simulation.RandIntBetween(r1, 0, 100))

	tc := []struct {
		status string
		error  bool
		msg    *types.QueryStakingOldTxRequest
		resp   *types.QueryStakingOldTxResponse
		addSeq bool
	}{
		{
			status: "no sequence exists",
			error:  true,
			msg:    nil,
			addSeq: false,
		},
		{
			status: "success",
			error:  false,
			msg:    &types.QueryStakingOldTxRequest{LogIndex: logIndex, TxHash: txHash.String()},
			resp:   &types.QueryStakingOldTxResponse{Status: true},
			addSeq: true,
		},
	}

	for _, c := range tc {
		suite.contractCaller.On("GetConfirmedTxReceipt", txHash.EthHash(), chainParams.MainchainTxConfirmations).Return(txReceipt, nil)
		if c.addSeq {
			sequence := new(big.Int).Mul(txReceipt.BlockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
			sequence.Add(sequence, new(big.Int).SetUint64(logIndex))
			app.StakingKeeper.SetStakingSequence(ctx, sequence.String())
		}

		resp, err := k.StakingOldTx(sdk.WrapSDKContext(ctx), c.msg)
		if c.error {
			require.Error(t, err)
			require.Nil(t, resp)
		} else {
			require.NoError(t, err)
			require.NotNil(t, resp)
			require.True(t, resp.Status)
		}
	}
}
