package types

import "errors"

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return NewGenesisState(make([]*EventRecord, 0), nil)
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	for _, sq := range gs.RecordSequences {
		if sq == "" {
			return errors.New("Invalid Sequence")
		}
	}
	return nil
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(eventRecords []*EventRecord, recordSequences []string) *GenesisState {
	return &GenesisState{EventRecords: eventRecords, RecordSequences: recordSequences}
}
