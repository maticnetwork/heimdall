package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/x/staking/keeper"
	stakingSim "github.com/maticnetwork/heimdall/x/staking/simulation"
	"github.com/maticnetwork/heimdall/x/staking/types"
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestQueryValidator() {
	t, app, ctx := suite.T(), suite.app, suite.ctx

	k := keeper.Querier{
		Keeper: app.StakingKeeper,
	}

	// loading the validators
	stakingSim.LoadValidatorSet(4, t, k.Keeper, ctx, false, 10)

	// getting the validators set
	validators := app.StakingKeeper.GetValidatorSet(ctx)
	validator := validators.Validators[0]

	queryValidatorInfo, err := k.Validator(sdk.WrapSDKContext(ctx), &types.QueryValidatorRequest{ValidatorId: validator.ID})

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
	stakingSim.LoadValidatorSet(4, t, k.Keeper, ctx, false, 10)

	vaSet, err := k.ValidatorSet(sdk.WrapSDKContext(ctx), &types.QueryValidatorSetRequest{})
	require.NoError(t, err)
	require.NotNil(t, vaSet)
}
