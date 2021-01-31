package types

import (
	"fmt"
)

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	for _, pastCommit := range gs.PastCommits {
		if pastCommit.Height <= 2 {
			return fmt.Errorf("past commit height must be greater 2")
		}

		if len(pastCommit.Txs) == 0 {
			return fmt.Errorf("txs must be present")
		}
	}
	return nil
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(pastCommits []*PastCommit) *GenesisState {
	return &GenesisState{
		PastCommits: pastCommits,
		Params: Params{
			Enabled: true,
		},
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() *GenesisState {
	return NewGenesisState(make([]*PastCommit, 0))
}
