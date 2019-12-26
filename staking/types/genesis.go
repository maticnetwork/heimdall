package types

import (
	"encoding/json"
	"errors"
	"math/big"

	"github.com/cosmos/cosmos-sdk/codec"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// GenesisValidator genesis validator
type GenesisValidator struct {
	ID         hmTypes.ValidatorID     `json:"id"`
	StartEpoch uint64                  `json:"start_epoch"`
	EndEpoch   uint64                  `json:"end_epoch"`
	Power      uint64                  `json:"power"` // aka Amount
	PubKey     hmTypes.PubKey          `json:"pub_key"`
	Signer     hmTypes.HeimdallAddress `json:"signer"`
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
	Validators           []*hmTypes.Validator `json:"validators" yaml:"validators"`
	CurrentValSet        hmTypes.ValidatorSet `json:"current_val_set" yaml:"current_val_set"`
	ValidatorRewards     map[string]*big.Int  `json:"val_rewards" yaml:"val_rewards"`
	ProposerBonusPercent int64                `json:"proposer_bonus_percent" yaml:"proposer_bonus_percent"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(
	validators []*hmTypes.Validator,
	currentValSet hmTypes.ValidatorSet,
	validatorRewards map[string]*big.Int,
	proposerBonusPercent int64,

) GenesisState {
	return GenesisState{
		Validators:           validators,
		CurrentValSet:        currentValSet,
		ValidatorRewards:     validatorRewards,
		ProposerBonusPercent: proposerBonusPercent,
	}
}

// // DefaultGenesisState returns a default genesis state
// func DefaultGenesisState(validators []*hmTypes.Validator, currentValSet hmTypes.ValidatorSet) GenesisState {
// 	validatorRewards := make(map[hmTypes.ValidatorID]*big.Int)
// 	for _, val := range validators {
// 		validatorRewards[val.ID] = big.NewInt(0)
// 	}
// 	return NewGenesisState(nil, currentValSet, validatorRewards, DefaultProposerBonusPercent)
// }

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(nil, hmTypes.ValidatorSet{}, nil, DefaultProposerBonusPercent)
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

// GetGenesisStateFromAppState returns staking GenesisState given raw application genesis state
func GetGenesisStateFromAppState(cdc *codec.Codec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		err := json.Unmarshal(appState[ModuleName], &genesisState)
		if err != nil {
			panic(err)
		}
	}

	return genesisState
}
