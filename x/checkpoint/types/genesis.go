package types

import (
	"encoding/json"
	"errors"
	"github.com/cosmos/cosmos-sdk/codec"
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
	if err := gs.Params.Validate(); err != nil {
		return err
	}

	if len(gs.Checkpoints) != 0 {
		if int(gs.AckCount) != len(gs.Checkpoints) {
			return errors.New("Incorrect state in state-dump , Please Check")
		}
	}

	return nil
}

// GetGenesisStateFromAppState returns staking GenesisState given raw application genesis state
func GetGenesisStateFromAppState(cdc codec.Marshaler, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}
	return genesisState
}
