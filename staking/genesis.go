package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	// get current val set
	var vals []*hmTypes.Validator
	if len(data.CurrentValSet.Validators) == 0 {
		vals = data.Validators
	} else {
		vals = data.CurrentValSet.Validators
	}

	if len(vals) != 0 {
		resultValSet := hmTypes.NewValidatorSet(vals)

		// add validators in store
		for _, validator := range resultValSet.Validators {
			// Add individual validator to state
			if err := keeper.AddValidator(ctx, *validator); err != nil {
				keeper.Logger(ctx).Error("Error InitGenesis", "error", err)
			}

			// update validator set in store
			if err := keeper.UpdateValidatorSetInStore(ctx, *resultValSet); err != nil {
				panic(err)
			}

			// increament accum if init validator set
			if len(data.CurrentValSet.Validators) == 0 {
				keeper.IncrementAccum(ctx, 1)
			}
		}
	}

	for _, sequence := range data.StakingSequences {
		keeper.SetStakingSequence(ctx, sequence)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	// return new genesis state
	return types.NewGenesisState(
		keeper.GetAllValidators(ctx),
		keeper.GetValidatorSet(ctx),
		keeper.GetStakingSequences(ctx),
	)
}
