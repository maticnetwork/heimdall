package bor

import (
	"encoding/json"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/bor/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// NewQuerier creates a querier for auth REST endpoints
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryParams:
			if len(path) == 1 {
				return queryParams(ctx, nil, req, keeper)
			}
			return queryParams(ctx, path[1:], req, keeper)
		case types.QuerySpan:
			return handleQuerySpan(ctx, req, keeper)
		case types.QuerySpanList:
			return handleQuerySpanList(ctx, req, keeper)
		case types.QueryLatestSpan:
			return handleQueryLatestSpan(ctx, req, keeper)
		case types.QueryNextProducers:
			return handleQueryNextProducers(ctx, req, keeper)
		case types.QueryNextSpanSeed:
			return handlerQueryNextSpanSeed(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown auth query endpoint")
		}
	}
}

func queryParams(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	if len(path) == 0 {
		bz, err := json.Marshal(keeper.GetParams(ctx))
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	}

	switch path[0] {
	case types.ParamSpan:
		bz, err := json.Marshal(keeper.GetParams(ctx).SpanDuration)
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	case types.ParamSprint:
		bz, err := json.Marshal(keeper.GetParams(ctx).SprintDuration)
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	case types.ParamProducerCount:
		bz, err := json.Marshal(keeper.GetParams(ctx).ProducerCount)
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	case types.ParamLastEthBlock:
		bz, err := json.Marshal(keeper.GetLastEthBlock(ctx))
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	default:
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("%s is not a valid query request path", req.Path))
	}
}

func handleQuerySpan(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params types.QuerySpanParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	span, err := keeper.GetSpan(ctx, params.RecordID)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not get span", err.Error()))
	}

	// return error if span doesn't exist
	if span == nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("span %v does not exist", params.RecordID))
	}

	// json record
	bz, err := json.Marshal(span)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func handleQuerySpanList(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var params hmTypes.QueryPaginationParams
	if err := keeper.cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("failed to parse params: %s", err))
	}

	res, err := keeper.GetSpanList(ctx, params.Page, params.Limit)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr(fmt.Sprintf("could not fetch span list with page %v and limit %v", params.Page, params.Limit), err.Error()))
	}

	bz, err := json.Marshal(res)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func handleQueryLatestSpan(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var defaultSpan hmTypes.Span
	spans := keeper.GetAllSpans(ctx)
	// if this is the first span return empty span
	if len(spans) == 0 {
		// json record
		bz, err := json.Marshal(defaultSpan)
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	}

	// explcitly fetch the last span
	span, err := keeper.GetLastSpan(ctx)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not get span", err.Error()))
	}

	// return error if span doesn't exist
	if span == nil {
		return nil, sdk.ErrInternal("latest span does not exist")
	}

	// json record
	bz, err := json.Marshal(span)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func handleQueryNextProducers(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	nextSpanSeed, err := keeper.GetNextSpanSeed(ctx)
	if err != nil {
		return nil, sdk.ErrInternal((sdk.AppendMsgToErr("cannot fetch next span seed from keeper", err.Error())))
	}

	nextProducers, err := keeper.SelectNextProducers(ctx, nextSpanSeed)
	if err != nil {
		return nil, sdk.ErrInternal((sdk.AppendMsgToErr("cannot fetch next producers from keeper", err.Error())))
	}

	bz, err := json.Marshal(nextProducers)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func handlerQueryNextSpanSeed(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	nextSpanSeed, err := keeper.GetNextSpanSeed(ctx)

	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("Error fetching next span seed", err.Error()))
	}

	// json record
	bz, err := json.Marshal(nextSpanSeed)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
