package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
	conf "github.com/maticnetwork/heimdall/helper"
	abci "github.com/tendermint/tendermint/abci/types"
)

func EndBlocker(ctx sdk.Context, k Keeper) (validators []abci.Validator) {
	var StakingLogger = conf.Logger.With("module", "staking")
	StakingLogger.Info("Current Validators Fetched ", "Validators", abci.ValidatorsString(k.GetAllValidators(ctx)))

	// flush exiting validator set
	k.FlushValidatorSet(ctx)
	// fetch current validator set
	validatorSet := helper.GetValidators()
	// update
	k.SetValidatorSet(ctx, validatorSet)

	StakingLogger.Info("New Validators ", "Validators", abci.ValidatorsString(k.GetAllValidators(ctx)))

	return validatorSet
}
