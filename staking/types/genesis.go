package types

import (
	"encoding/json"
	"errors"

	hmTypes "github.com/maticnetwork/heimdall/types"
)

// GenesisValidator genesis validator
type GenesisValidator struct {
	ID                   hmTypes.ValidatorID     `json:"id"`
	StartEpoch           uint64                  `json:"start_epoch"`
	EndEpoch             uint64                  `json:"end_epoch"`
	Power                uint64                  `json:"power"` // aka Amount
	DelegatedPower       int64                   `json:"delegatedpower"`
	DelgatorRewardPool   string                  `json:delegatorRewardPool`
	TotalDelegatorShares string                  `json:totalDelegatorShares`
	PubKey               hmTypes.PubKey          `json:"pub_key"`
	Signer               hmTypes.HeimdallAddress `json:"signer"`
}

// HeimdallValidator converts genesis validator validator to Heimdall validator
func (v *GenesisValidator) HeimdallValidator() hmTypes.Validator {
	return hmTypes.Validator{
		ID:                   v.ID,
		PubKey:               v.PubKey,
		VotingPower:          int64(v.Power),
		DelegatedPower:       int64(v.DelegatedPower),
		DelgatorRewardPool:   v.DelgatorRewardPool,
		TotalDelegatorShares: v.TotalDelegatorShares,
		StartEpoch:           v.StartEpoch,
		EndEpoch:             v.EndEpoch,
		Signer:               v.Signer,
	}
}

// GenesisState is the checkpoint state that must be provided at genesis.
type GenesisState struct {
	Validators           []*hmTypes.Validator      `json:"validators" yaml:"validators"`
	CurrentValSet        hmTypes.ValidatorSet      `json:"current_val_set" yaml:"current_val_set"`
	DividentAccounts     []hmTypes.DividendAccount `json:"dividend_accounts" yaml:"dividend_accounts"`
	ProposerBonusPercent int64                     `json:"proposer_bonus_percent" yaml:"proposer_bonus_percent"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(
	validators []*hmTypes.Validator,
	currentValSet hmTypes.ValidatorSet,
	dividentAccounts []hmTypes.DividendAccount,

) GenesisState {
	return GenesisState{
		Validators:       validators,
		CurrentValSet:    currentValSet,
		DividentAccounts: dividentAccounts,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(nil, hmTypes.ValidatorSet{}, nil)
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
func GetGenesisStateFromAppState(appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		err := json.Unmarshal(appState[ModuleName], &genesisState)
		if err != nil {
			panic(err)
		}
	}

	return genesisState
}

// SetGenesisStateToAppState sets state into app state
func SetGenesisStateToAppState(appState map[string]json.RawMessage, validators []*hmTypes.Validator, currentValSet hmTypes.ValidatorSet) (map[string]json.RawMessage, error) {
	// set state to staking state
	stakingState := GetGenesisStateFromAppState(appState)
	stakingState.Validators = validators
	stakingState.CurrentValSet = currentValSet
	stakingState.DividentAccounts = make([]hmTypes.DividendAccount, 0)

	var err error
	appState[ModuleName], err = json.Marshal(stakingState)
	if err != nil {
		return appState, err
	}
	return appState, nil
}
