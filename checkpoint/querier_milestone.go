package checkpoint

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/common"
)

// handleQueryLatestMilestone to get the latest milestone
func handleQueryLatestMilestone(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	res, err := keeper.GetLastMilestone(ctx)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not fetch milestone", err.Error()))
	}

	if res == nil {
		return nil, common.ErrNoMilestoneFound(keeper.Codespace())
	}

	bz, err := json.Marshal(res)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

// handleQueryMilestoneByNumber to get the milestone by number
func handleQueryMilestoneByNumber(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QueryMilestoneParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	res, err := keeper.GetMilestoneByNumber(ctx, params.Number)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not fetch milestone", err.Error()))
	}

	if res == nil {
		return nil, common.ErrNoMilestoneFound(keeper.Codespace())
	}

	bz, err := json.Marshal(res)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

// handleQueryCount to get the count
func handleQueryCount(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	bz, err := json.Marshal(keeper.GetMilestoneCount(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

// handleQueryLatestNoAckMilestone to get lasted no ack milestone id
func handleQueryLatestNoAckMilestone(ctx sdk.Context, keeper Keeper) ([]byte, sdk.Error) {
	res := keeper.GetLastNoAckMilestone(ctx)

	bz, err := json.Marshal(res)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

// handleQueryNoAckMilestoneByID to check whether the particular id exist in no-ack list(rejected milestone list)
func handleQueryNoAckMilestoneByID(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var ID types.QueryMilestoneID
	if err := keeper.cdc.UnmarshalJSON(req.Data, &ID); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse milestoneID: %s", err))
	}

	res := keeper.GetNoAckMilestone(ctx, ID.MilestoneID)

	bz, err := json.Marshal(res)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
