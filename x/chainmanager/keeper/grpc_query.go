package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/x/chainmanager/types"
)

// Querier is used as Keeper will have duplicate methods if used directly,
// and gRPC names take precedence over keeper
type Querier struct {
	Keeper
}

var _ types.QueryServer = Querier{}

// Params queries params info
func (k Querier) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	params := k.GetParams(sdk.UnwrapSDKContext(c))

	return &types.QueryParamsResponse{Params: &params}, nil
}
