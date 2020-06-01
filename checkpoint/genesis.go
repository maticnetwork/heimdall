package checkpoint

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/checkpoint/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	keeper.SetParams(ctx, data.Params)

	// Set last no-ack
	if data.LastNoACK > 0 {
		keeper.SetLastNoAck(ctx, data.LastNoACK)
	}

	// Add finalised checkpoints to state
	if len(data.Checkpoints) != 0 {
		// check if we are provided all the headers
		if int(data.AckCount) != len(data.Checkpoints) {
			panic(errors.New("Incorrect state in state-dump , Please Check "))
		}
		// sort headers before loading to state
		data.Checkpoints = hmTypes.SortHeaders(data.Checkpoints)
		// load checkpoints to state
		for i, checkpoint := range data.Checkpoints {
			checkpointIndex := uint64(i) + 1
			if err := keeper.AddCheckpoint(ctx, checkpointIndex, checkpoint); err != nil {
				keeper.Logger(ctx).Error("InitGenesis | AddCheckpoint", "error", err)
			}
		}
	}

	// Add checkpoint in buffer
	if data.BufferedCheckpoint != nil {
		if err := keeper.SetCheckpointBuffer(ctx, *data.BufferedCheckpoint); err != nil {
			keeper.Logger(ctx).Error("InitGenesis | SetCheckpointBuffer", "error", err)
		}
	}

	// Set initial ack count
	keeper.UpdateACKCountWithValue(ctx, data.AckCount)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	params := keeper.GetParams(ctx)

	bufferedCheckpoint, _ := keeper.GetCheckpointFromBuffer(ctx)
	return types.NewGenesisState(
		params,
		bufferedCheckpoint,
		keeper.GetLastNoAck(ctx),
		keeper.GetACKCount(ctx),
		hmTypes.SortHeaders(keeper.GetCheckpoints(ctx)),
	)
}
