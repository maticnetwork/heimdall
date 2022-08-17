package app

import (
	"encoding/json"
)

// GenesisState the genesis state of the blockchain is represented here as a map of raw json messages keyed by an identifier string
type GenesisState map[string]json.RawMessage

// NewDefaultGenesisState generates the default state for the application.
func NewDefaultGenesisState() GenesisState {
	return ModuleBasics.DefaultGenesis()
}
