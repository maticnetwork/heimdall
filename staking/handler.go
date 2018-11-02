package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
	abci "github.com/tendermint/tendermint/abci/types"
)

// EndBlocker refreshes validator set after block commit
func EndBlocker(ctx sdk.Context, k Keeper) (validators []abci.ValidatorUpdate) {
	StakingLogger.Info("Current validators fetched", "validators", (k.GetAllValidators(ctx)))

	// flush exiting validator set
	k.FlushValidatorSet(ctx)
	// fetch current validator set
	validatorSet := helper.GetValidators()
	// update
	k.SetValidatorSet(ctx, validatorSet)

	StakingLogger.Info("New validators set", "validators", (k.GetAllValidators(ctx)))
	return validatorSet
}
