package bor

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/bor/types"
)

// NewQuerier creates a querier for auth REST endpoints
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case types.QueryParams:
			return queryParams(ctx, path[1:], req, keeper)
		case types.QuerySpan:
			return handleQuerySpan(ctx, req, keeper)
		case types.QueryLatestSpan:
			return handleQueryLatestSpan(ctx, req, keeper)
		case types.QueryNextProducers:
			return handleQueryNextProducers(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown auth query endpoint")
		}
	}
}

func queryParams(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	switch path[0] {
	case types.ParamSpan:
		bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetSpanDuration(ctx))
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	case types.ParamSprint:
		bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetSprintDuration(ctx))
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	case types.ParamProducerCount:
		var bz []byte
		count, err := keeper.GetProducerCount(ctx)
		if err != nil {
			return bz, sdk.ErrInternal(sdk.AppendMsgToErr("cannot fetch producer count from keeper", err.Error()))
		}
		bz, err = codec.MarshalJSONIndent(keeper.cdc, count)
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	case types.ParamLastEthBlock:
		bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetLastEthBlock(ctx))
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
	bz, err := codec.MarshalJSONIndent(keeper.cdc, span)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func handleQueryLatestSpan(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	span, err := keeper.GetLastSpan(ctx)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not get span", err.Error()))
	}

	// return error if span doesn't exist
	if span == nil {
		return nil, sdk.ErrInternal(fmt.Sprintf("latest span does not exist"))
	}

	// json record
	bz, err := codec.MarshalJSONIndent(keeper.cdc, span)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}

func handleQueryNextProducers(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	nextProducers, err := keeper.SelectNextProducers(ctx)
	if err != nil {
		return nil, sdk.ErrInternal((sdk.AppendMsgToErr("cannot fetch next producers from keeper", err.Error())))
	}

	bz, err := codec.MarshalJSONIndent(keeper.cdc, nextProducers)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
	}
	return bz, nil
}
