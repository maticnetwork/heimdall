package types

import "errors"

// GenesisState is the bank state that must be provided at genesis.
type GenesisState struct {
	EventRecords    []*EventRecord `json:"event_records"`
	RecordSequences []string       `json:"record_sequences" yaml:"record_sequences"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(eventRecords []*EventRecord, recordSequences []string) GenesisState {
	return GenesisState{EventRecords: eventRecords, RecordSequences: recordSequences}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(make([]*EventRecord, 0), nil)
}

// ValidateGenesis performs basic validation of bank genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	for _, sq := range data.RecordSequences {
		if sq == "" {
			return errors.New("Invalid Sequence")
		}
	}
	return nil
}
