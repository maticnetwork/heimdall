package milestone

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/milestone/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	keeper.SetParams(ctx, data.Params)

	// Add milestone
	// Add finalised milestones to state
	if len(data.Milestones) != 0 {
		// check if we are provided all the headers
		// sort headers before loading to state based on Timestamp
		data.Milestones = hmTypes.SortMilestone(data.Milestones)
		// load checkpoints to state
		for _, milestone := range data.Milestones {
			if err := keeper.AddMilestone(ctx, milestone); err != nil {
				keeper.Logger(ctx).Error("InitGenesis | AddMilestone",
					"milestone", milestone.String(),
					"error", err)
			}
		}
	}

	// Add No Ack Milestone
	// Add no ack milestone to state
	if len(data.NoAckMilestones) != 0 {
		// sort no ack milestone before loading to state based on the timestamps
		data.NoAckMilestones = hmTypes.SortMilestone(data.NoAckMilestones)
		// load milestones to state
		for _, milestone := range data.NoAckMilestones {
			keeper.SetNoAckMilestone(ctx, milestone.MilestoneID)
			keeper.Logger(ctx).Error("InitGenesis | AddMilestone",
				"noAckMilestone", milestone.String(),
				"milestoneID", milestone.MilestoneID,
			)
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	params := keeper.GetParams(ctx)

	milestone, _ := keeper.GetLastMilestone(ctx)

	return types.NewGenesisState(
		params,
		milestone,
	)
}
