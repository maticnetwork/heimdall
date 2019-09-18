package bor

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the auth Querier
const (
	QueryParams = "params"

	ParamSpan          = "span"
	ParamSprint        = "sprint"
	ParamProducerCount = "producer-count"
	ParamLastEthBlock  = "last-eth-block"
)

// NewQuerier creates a querier for auth REST endpoints
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		case QueryParams:
			return queryParams(ctx, path[1:], req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown auth query endpoint")
		}
	}
}

func queryParams(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	switch path[0] {
	case ParamSpan:
		bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetSpanDuration(ctx))
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	case ParamSprint:
		bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetSprintDuration(ctx))
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	case ParamProducerCount:
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
	case ParamLastEthBlock:
		bz, err := codec.MarshalJSONIndent(keeper.cdc, keeper.GetLastEthBlock(ctx))
		if err != nil {
			return nil, sdk.ErrInternal(sdk.AppendMsgToErr("could not marshal result to JSON", err.Error()))
		}
		return bz, nil
	default:
		return nil, sdk.ErrUnknownRequest(fmt.Sprintf("%s is not a valid query request path", req.Path))
	}
}
