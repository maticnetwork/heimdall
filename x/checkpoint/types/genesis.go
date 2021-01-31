package types

import (
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// DefaultIndex is the default capability global index
const DefaultIndex uint64 = 1

// NewGenesisState creates a new genesis state.
func NewGenesisState(
	params Params,
	bufferedCheckpoint *hmTypes.Checkpoint,
	lastNoACK uint64,
	ackCount uint64,
	checkpoints []*hmTypes.Checkpoint,
) *GenesisState {
	return &GenesisState{
		Params:             params,
		BufferedCheckpoint: bufferedCheckpoint,
		LastNoACK:          lastNoACK,
		AckCount:           ackCount,
		Checkpoints:        checkpoints,
	}
}

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		Params: DefaultParams(),
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	return nil
}
