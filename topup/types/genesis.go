package types

import "errors"

// GenesisState is the bank state that must be provided at genesis.
type GenesisState struct {
	TopupSequences []*uint64 `json:"tx_sequences" yaml:"tx_sequences"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(topupSequence []*uint64) GenesisState {
	return GenesisState{
		TopupSequences: topupSequence,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(nil)
}

// ValidateGenesis performs basic validation of topup genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	for _, sq := range data.TopupSequences {
		if sq == nil {
			return errors.New("Invalid Sequence")
		}
	}
	return nil
}
