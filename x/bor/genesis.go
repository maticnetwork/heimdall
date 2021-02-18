package bor

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	borKeeper "github.com/maticnetwork/heimdall/x/bor/keeper"
	"github.com/maticnetwork/heimdall/x/bor/types"
)

func InitGenesis(ctx sdk.Context, keeper borKeeper.Keeper, data types.GenesisState) {
	keeper.SetParams(ctx, data.Params)

	if len(data.Spans) > 0 {
		// sort data spans before inserting to ensure lastspanId fetched is correct
		hmTypes.SortSpanByID(data.Spans)
		// add new span
		for _, span := range data.Spans {
			if err := keeper.AddNewRawSpan(ctx, *span); err != nil {
				keeper.Logger(ctx).Error("Error AddNewRawSpan", "error", err)
			}
		}

		// update last span
		keeper.UpdateLastSpan(ctx, data.Spans[len(data.Spans)-1].ID)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper borKeeper.Keeper) *types.GenesisState {
	params := keeper.GetParams(ctx)

	allSpans, err := keeper.GetAllSpans(ctx)
	if err != nil {
		keeper.Logger(ctx).Error("Error ExportGenesis", "error", err)
	}
	hmTypes.SortSpanByID(allSpans)
	return types.NewGenesisState(
		params,
		allSpans,
	)
}
