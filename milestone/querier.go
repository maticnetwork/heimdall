package milestone

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/milestone/types"
	"github.com/maticnetwork/heimdall/staking"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier creates a querier for auth REST endpoints
func NewQuerier(keeper Keeper, stakingKeeper staking.Keeper, contractCaller helper.IContractCaller) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryParams:
			return handleQueryParams(ctx, req, keeper)

		case types.QueryMilestone:
			return handleQueryMilestone(ctx, req, keeper)

		case types.QueryCount:
			return handleQueryCount(ctx, req, keeper)

		default:
			return nil, sdk.ErrUnknownRequest("unknown auth query endpoint")
		}
	}
}

func handleQueryParams(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	bz, err := json.Marshal(keeper.GetParams(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}

func handleQueryMilestone(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	res, err := keeper.GetMilestone(ctx)
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

func handleQueryCount(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	bz, err := json.Marshal(keeper.GetCount(ctx))
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}

	return bz, nil
}
