package staking

import "github.com/maticnetwork/heimdall/types"

// GenesisState - all staking state that must be provided at genesis
type GenesisState struct {
	Validators []types.Validator `json:"validators"`
}

func NewGenesisState(validators []types.Validator) GenesisState {
	return GenesisState{
		Validators: validators,
	}
}

// get raw genesis raw message for testing
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}
