package milestone

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/milestone/types"
)

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	// Add milestone
	if data.Milestone != nil {

		if err := keeper.AddMilestone(ctx, *data.Milestone); err != nil {
			keeper.Logger(ctx).Error("InitGenesis | SetMilestone", "error", err)
		}
	}

}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {

	milestone, _ := keeper.GetLastMilestone(ctx)

	return types.NewGenesisState(
		milestone,
	)
}
