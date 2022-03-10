package clerk

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/clerk/types"
)

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	// add checkpoint headers
	if len(data.EventRecords) != 0 {
		for _, record := range data.EventRecords {
			if err := keeper.SetEventRecord(ctx, *record); err != nil {
				keeper.Logger(ctx).Error("InitGenesis | SetEventRecord", "error", err)
				panic(err)
			} else {
				fmt.Println("Imported record", "record", record.ID, "recordTime", record.RecordTime)
			}
		}
	}

	for _, sequence := range data.RecordSequences {
		keeper.Logger(ctx).Info("Imported record sequence", "sequence", sequence)
		keeper.SetRecordSequence(ctx, sequence)
	}

}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	return types.NewGenesisState(keeper.GetAllEventRecords(ctx), keeper.GetRecordSequences(ctx))
}
