package staking

import (
	"errors"
	"math/big"

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
		ID:          v.ID,
		PubKey:      v.PubKey,
		VotingPower: int64(v.Power),
		StartEpoch:  v.StartEpoch,
		EndEpoch:    v.EndEpoch,
		Signer:      v.Signer,
	}
}

// GenesisState is the checkpoint state that must be provided at genesis.
type GenesisState struct {
	Validators           []*hmTypes.Validator           `json:"validators" yaml:"validators"`
	CurrentValSet        hmTypes.ValidatorSet           `json:"current_val_set" yaml:"current_val_set"`
	ValidatorRewards     map[types.ValidatorID]*big.Int `json:"val_rewards" yaml:"val_rewards"`
	ProposerBonusPercent int64                          `json:"proposer_bonus_percent" yaml:"proposer_bonus_percent"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(
	validators []*hmTypes.Validator,
	currentValSet hmTypes.ValidatorSet,
	validatorRewards map[types.ValidatorID]*big.Int,
	proposerBonusPercent int64,

) GenesisState {
	return GenesisState{
		Validators:           validators,
		CurrentValSet:        currentValSet,
		ValidatorRewards:     validatorRewards,
		ProposerBonusPercent: proposerBonusPercent,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState(validators []*hmTypes.Validator, currentValSet hmTypes.ValidatorSet) GenesisState {
	validatorRewards := make(map[types.ValidatorID]*big.Int)
	for _, val := range validators {
		validatorRewards[val.ID] = big.NewInt(0)
	}
	return NewGenesisState(validators, currentValSet, validatorRewards, DefaultProposerBonusPercent)
}

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	// get current val set
	var vals []*hmTypes.Validator
	if len(data.CurrentValSet.Validators) == 0 {
		vals = data.Validators
	} else {
		vals = data.CurrentValSet.Validators
	}

	// result
	resultValSet := hmTypes.NewValidatorSet(vals)
	validatorRewards := make(map[types.ValidatorID]*big.Int)

	// add validators in store
	for _, validator := range resultValSet.Validators {
		// Add individual validator to state
		keeper.AddValidator(ctx, *validator)

	}

	// TODO match valSet and genesisState.CurrentValSet for difference in accum
	// update validator set in store
	if err := keeper.UpdateValidatorSetInStore(ctx, *resultValSet); err != nil {
		panic(err)
	}

	// Add rewards for initial validators
	for _, validator := range data.Validators {
		if _, ok := data.ValidatorRewards[validator.ID]; ok {
			validatorRewards[validator.ID] = data.ValidatorRewards[validator.ID]
		} else {
			validatorRewards[validator.ID] = big.NewInt(0)
		}
	}
	keeper.UpdateValidatorRewards(ctx, validatorRewards)

	keeper.SetProposerBonusPercent(ctx, data.ProposerBonusPercent)

}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	// return new genesis state
	return NewGenesisState(
		keeper.GetAllValidators(ctx),
		keeper.GetValidatorSet(ctx),
		keeper.GetAllValidatorRewards(ctx),
		keeper.GetProposerBonusPercent(ctx),
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
