package clerk

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/clerk/types"
)

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	// add checkpoint headers
	if len(data.EventRecords) != 0 {
		for _, record := range data.EventRecords {
			keeper.SetEventRecord(ctx, *record)
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	return types.NewGenesisState(keeper.GetAllEventRecords(ctx))
}