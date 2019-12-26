package bor

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/bor/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	keeper.SetSprintDuration(ctx, data.SprintDuration)
	keeper.SetSpanDuration(ctx, data.SpanDuration)
	keeper.SetProducerCount(ctx, data.ProducerCount)
	if len(data.Spans) > 0 {
		// sort data spans before inserting to ensure lastspanId fetched is correct
		hmTypes.SortSpanByID(data.Spans)
		// add new span
		for _, span := range data.Spans {
			keeper.AddNewRawSpan(ctx, *span)
		}

		// update last span
		keeper.UpdateLastSpan(ctx, data.Spans[len(data.Spans)-1].ID)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	producerCount, _ := keeper.GetProducerCount(ctx)
	allSpans := keeper.GetAllSpans(ctx)
	hmTypes.SortSpanByID(allSpans)
	return types.NewGenesisState(
		keeper.GetSprintDuration(ctx),
		keeper.GetSpanDuration(ctx),
		producerCount,
		// TODO think better way to export all spans
		allSpans,
	)
}
