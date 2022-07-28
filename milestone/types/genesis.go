package types

import (
	"encoding/json"

	"github.com/maticnetwork/heimdall/bor/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// GenesisState is the milestone state that must be provided at genesis.
type GenesisState struct {
	Params Params `json:"params" yaml:"params"`

	Milestone *hmTypes.Milestone `json:"milestone" yaml:"milestone"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(
	params Params,
	milestone *hmTypes.Milestone,
) GenesisState {
	return GenesisState{
		Params: params,

		Milestone: milestone,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{
		Params: DefaultParams(),
	}
}

// ValidateGenesis performs basic validation of bor genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	if err := data.Params.Validate(); err != nil {
		return err
	}

	return nil
}

// GetGenesisStateFromAppState returns staking GenesisState given raw application genesis state
func GetGenesisStateFromAppState(appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		types.ModuleCdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return genesisState
}
