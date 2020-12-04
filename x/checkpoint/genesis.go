package checkpoint

import (
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/x/checkpoint/keeper"
	"github.com/maticnetwork/heimdall/x/checkpoint/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, keeper keeper.Keeper, genState types.GenesisState) {
	keeper.SetParams(ctx, genState.Params)

	// Set last no-ack
	if genState.LastNoACK > 0 {
		keeper.SetLastNoAck(ctx, genState.LastNoACK)
	}

	// Add finalised checkpoints to state
	if len(genState.Checkpoints) != 0 {
		// check if we are provided all the headers
		if int(genState.AckCount) != len(genState.Checkpoints) {
			panic(errors.New("Incorrect state in state-dump , Please Check "))
		}
		// sort headers before loading to state
		genState.Checkpoints = hmTypes.SortHeaders(genState.Checkpoints)
		// load checkpoints to state
		for i, checkpoint := range genState.Checkpoints {
			checkpointIndex := uint64(i) + 1
			if err := keeper.AddCheckpoint(ctx, checkpointIndex, checkpoint); err != nil {
				keeper.Logger(ctx).Error("InitGenesis | AddCheckpoint", "error", err)
			}
		}
	}

	// Add checkpoint in buffer
	if genState.BufferedCheckpoint != nil {
		if err := keeper.SetCheckpointBuffer(ctx, genState.BufferedCheckpoint); err != nil {
			keeper.Logger(ctx).Error("InitGenesis | SetCheckpointBuffer", "error", err)
		}
	}

	// Set initial ack count
	keeper.UpdateACKCountWithValue(ctx, genState.AckCount)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.DefaultGenesis()
}
