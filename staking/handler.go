package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
	abci "github.com/tendermint/tendermint/abci/types"
)

// EndBlocker refreshes validator set after block commit
func EndBlocker(ctx sdk.Context, k Keeper) (validators []abci.Validator) {
	StakingLogger.Info("Current validators fetched", "validators", abci.ValidatorsString(k.GetAllValidators(ctx)))

	// flush exiting validator set
	k.FlushValidatorSet(ctx)
	// fetch current validator set
	validatorSet := helper.GetValidators()
	// update
	k.SetValidatorSet(ctx, validatorSet)

	StakingLogger.Info("New validators set", "validators", abci.ValidatorsString(k.GetAllValidators(ctx)))
	return validatorSet
}
