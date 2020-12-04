package staking

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/x/staking/keeper"
	"github.com/maticnetwork/heimdall/x/staking/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, genState types.GenesisState) {
	// get current val set
	var vals []*hmTypes.Validator
	if len(genState.CurrentValSet.Validators) == 0 {
		vals = genState.Validators
	} else {
		vals = genState.CurrentValSet.Validators
	}

	fmt.Println("genState", genState.Validators)

	if len(vals) != 0 {
		resultValSet := hmTypes.NewValidatorSet(vals)

		// add validators in store
		for _, validator := range resultValSet.Validators {
			// Add individual validator to state
			if err := keeper.AddValidator(ctx, *validator); err != nil {
				keeper.Logger(ctx).Error("Error InitGenesis", "error", err)
			}

			// update validator set in store
			if err := keeper.UpdateValidatorSetInStore(ctx, resultValSet); err != nil {
				panic(err)
			}

			// increament accum if init validator set
			if len(genState.CurrentValSet.Validators) == 0 {
				keeper.IncrementAccum(ctx, 1)
			}
		}
	}

	for _, sequence := range genState.StakingSequences {
		keeper.SetStakingSequence(ctx, sequence)
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, keeper keeper.Keeper) *types.GenesisState {
	// return new genesis state
	return types.NewGenesisState(
		keeper.GetAllValidators(ctx),
		keeper.GetValidatorSet(ctx),
		keeper.GetStakingSequences(ctx),
	)
}
