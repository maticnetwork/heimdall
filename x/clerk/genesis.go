package clerk

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/x/clerk/keeper"
	"github.com/maticnetwork/heimdall/x/clerk/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// add checkpoint headers
	if len(genState.EventRecords) != 0 {
		for _, record := range genState.EventRecords {
			if err := k.SetEventRecord(ctx, *record); err != nil {
				k.Logger(ctx).Error("InitGenesis | SetEventRecord", "error", err)
			}
		}
	}

	for _, sequence := range genState.RecordSequences {
		k.SetRecordSequence(ctx, sequence)
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genState := types.NewGenesisState(k.GetAllEventRecords(ctx), k.GetRecordSequences(ctx))
	return &genState
}