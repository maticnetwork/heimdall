package checkpoint

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
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

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	// Set last no-ack
	if data.LastNoACK > 0 {
		keeper.SetLastNoAck(ctx, data.LastNoACK)
	}

	// add checkpoint headers
	if len(data.Headers) != 0 {
		if int(data.AckCount) != len(data.Headers) {
			panic(errors.New("Incorrect state in state-dump , Please Check "))
		}

		for i, header := range data.Headers {
			checkpointHeaderIndex := helper.GetConfig().ChildBlockInterval * (uint64(i) + 1)
			keeper.AddCheckpoint(ctx, checkpointHeaderIndex, header)
		}
	}

	// Add checkpoint in buffer
	if data.BufferedCheckpoint != nil {
		keeper.SetCheckpointBuffer(ctx, *data.BufferedCheckpoint)
	}

	// Set initial ack count
	keeper.UpdateACKCountWithValue(ctx, data.AckCount)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	bufferedCheckpoint, _ := keeper.GetCheckpointFromBuffer(ctx)
	return NewGenesisState(
		bufferedCheckpoint,
		keeper.GetLastNoAck(ctx),
		keeper.GetACKCount(ctx),
		keeper.GetCheckpointHeaders(ctx),
	)
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
