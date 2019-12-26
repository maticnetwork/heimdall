package checkpoint

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	// Set last no-ack
	if data.LastNoACK > 0 {
		keeper.SetLastNoAck(ctx, data.LastNoACK)
	}

	// Add finalised checkpoints to state
	if len(data.Headers) != 0 {
		// check if we are provided all the headers
		if int(data.AckCount) != len(data.Headers) {
			panic(errors.New("Incorrect state in state-dump , Please Check "))
		}
		// sort headers before loading to state
		data.Headers = hmTypes.SortHeaders(data.Headers)

		// load checkpoints to state
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
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	bufferedCheckpoint, _ := keeper.GetCheckpointFromBuffer(ctx)
	return types.NewGenesisState(
		bufferedCheckpoint,
		keeper.GetLastNoAck(ctx),
		keeper.GetACKCount(ctx),
		hmTypes.SortHeaders(keeper.GetCheckpointHeaders(ctx)),
	)
}
