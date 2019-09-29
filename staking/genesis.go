package staking

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// GenesisValidator genesis validator
type GenesisValidator struct {
	ID         hmTypes.ValidatorID   `json:"id"`
	StartEpoch uint64                `json:"start_epoch"`
	EndEpoch   uint64                `json:"end_epoch"`
	Power      uint64                `json:"power"` // aka Amount
	PubKey     hmTypes.PubKey        `json:"pub_key"`
	Signer     types.HeimdallAddress `json:"signer"`
}

// HeimdallValidator converts genesis validator validator to Heimdall validator
func (v *GenesisValidator) HeimdallValidator() hmTypes.Validator {
	return hmTypes.Validator{
		ID:         v.ID,
		PubKey:     v.PubKey,
		Power:      v.Power,
		StartEpoch: v.StartEpoch,
		EndEpoch:   v.EndEpoch,
		Signer:     v.Signer,
	}
}

// GenesisState is the checkpoint state that must be provided at genesis.
type GenesisState struct {
	Validators    []hmTypes.Validator  `json:"validators" yaml:"validators"`
	CurrentValSet hmTypes.ValidatorSet `json:"current_val_set" yaml:"current_val_set"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(
	validators []hmTypes.Validator,
	currentValSet hmTypes.ValidatorSet,
) GenesisState {
	return GenesisState{
		Validators:    validators,
		CurrentValSet: currentValSet,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState(validators []hmTypes.Validator, currentValSet hmTypes.ValidatorSet) GenesisState {
	return NewGenesisState(validators, currentValSet)
}

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	// update validators from genesis
	var newValSet hmTypes.ValidatorSet
	for _, validator := range data.Validators {
		if ok := newValSet.Add(&validator); !ok {
			panic(errors.New("Error while adding new validator"))
		} else {
			// Add individual validator to state
			keeper.AddValidator(ctx, validator)
		}
	}

	// get current val set
	var currentValSet hmTypes.ValidatorSet
	if len(data.CurrentValSet.Validators) == 0 {
		currentValSet = newValSet
	} else {
		currentValSet = data.CurrentValSet
	}

	// result
	resultValSet := hmTypes.NewValidatorSet(currentValSet.Validators)

	// TODO match valSet and genesisState.CurrentValSet for difference in accum
	// update validator set in store
	if err := keeper.UpdateValidatorSetInStore(ctx, *resultValSet); err != nil {
		panic(err)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	var validators []hmTypes.Validator
	for _, v := range keeper.GetAllValidators(ctx) {
		validators = append(validators, *v)
	}

	// return new genesis state
	return NewGenesisState(
		validators,
		keeper.GetValidatorSet(ctx),
	)
}

// ValidateGenesis performs basic validation of bor genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	for _, validator := range data.Validators {
		if !validator.ValidateBasic() {
			return errors.New("Invalid validator")
		}
	}

	return nil
}
