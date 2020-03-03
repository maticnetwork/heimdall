package types

// GenesisState is the bank state that must be provided at genesis.
type GenesisState struct {
	// TopupSequences map[uint256]bool `json:"topup_sequence" yaml:"topup_sequence"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState() GenesisState {
	return GenesisState{
		// TopupSequence: topupSequence,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState { return NewGenesisState() }

// ValidateGenesis performs basic validation of bank genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error { return nil }
