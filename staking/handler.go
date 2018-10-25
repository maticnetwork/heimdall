package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
	conf "github.com/maticnetwork/heimdall/helper"
	abci "github.com/tendermint/tendermint/abci/types"
)

func EndBlocker(ctx sdk.Context, k Keeper) (validators []abci.Validator) {
	var StakingLogger = conf.Logger.With("module", "staking")

	// validator := k.GetValidatorInfo(ctx, _address)
	StakingLogger.Info("Current Validators Fetched \n", "Validators", k.GetAllValidators(ctx))

	validatorSet := helper.GetValidators()
	k.FlushValidatorSet(ctx)
	k.SetValidatorSet(ctx, validatorSet)
	StakingLogger.Info("New Validators Are \n", "Validators", k.GetAllValidators(ctx))

	return validatorSet
}
