package types

// GenesisState is the bank state that must be provided at genesis.
type GenesisState struct {
	EventRecords []*EventRecord `json:"event_records"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(eventRecords []*EventRecord) GenesisState {
	return GenesisState{EventRecords: eventRecords}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(make([]*EventRecord, 0))
}

// ValidateGenesis performs basic validation of bank genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	return nil
}
