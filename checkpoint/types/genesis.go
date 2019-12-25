package types

import (
	"errors"

	hmTypes "github.com/maticnetwork/heimdall/types"
)

// GenesisState is the checkpoint state that must be provided at genesis.
type GenesisState struct {
	BufferedCheckpoint *hmTypes.CheckpointBlockHeader  `json:"buffered_checkpoint" yaml:"buffered_checkpoint"`
	LastNoACK          uint64                          `json:"last_no_ack" yaml:"last_no_ack"`
	AckCount           uint64                          `json:"ack_count" yaml:"ack_count"`
	Headers            []hmTypes.CheckpointBlockHeader `json:"headers" yaml:"headers"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(
	bufferedCheckpoint *hmTypes.CheckpointBlockHeader,
	lastNoACK uint64,
	ackCount uint64,
	headers []hmTypes.CheckpointBlockHeader,
) GenesisState {
	return GenesisState{
		BufferedCheckpoint: bufferedCheckpoint,
		LastNoACK:          lastNoACK,
		AckCount:           ackCount,
		Headers:            headers,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return GenesisState{}
}

// ValidateGenesis performs basic validation of bor genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	if len(data.Headers) != 0 {
		if int(data.AckCount) != len(data.Headers) {
			return errors.New("Incorrect state in state-dump , Please Check")
		}
	}

	return nil
}
