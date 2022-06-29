package types

import (
	"fmt"

	tmTypes "github.com/tendermint/tendermint/types"
)

// PastCommit represent past commit for the record and process side-txs
type PastCommit struct {
	Height int64       `json:"height" yaml:"height"`
	Txs    tmTypes.Txs `json:"txs" yaml:"txs"`
}

// GenesisState is the sidechannel state that must be provided at genesis.
type GenesisState struct {
	PastCommits []PastCommit `json:"past_commits" yaml:"past_commits"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(pastCommits []PastCommit) GenesisState {
	return GenesisState{
		PastCommits: pastCommits,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(make([]PastCommit, 0))
}

// ValidateGenesis performs basic validation of topup genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	for _, pastCommit := range data.PastCommits {
		if pastCommit.Height <= 2 {
			return fmt.Errorf("Past commit height must be greater 2")
		}

		if len(pastCommit.Txs) == 0 {
			return fmt.Errorf("Txs must be present")
		}
	}

	return nil
}
