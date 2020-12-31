package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/x/checkpoint/types"
)

// Querier is used as Keeper will have duplicate methods if used directly, and gRPC names take precedence over keeper
type Querier struct {
	Keeper
	contractCaller helper.IContractCaller
}

// NewQueryServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewQueryServerImpl(keeper Keeper, contractCaller helper.IContractCaller) types.QueryServer {
	return &Querier{Keeper: keeper, contractCaller: contractCaller}
}

var _ types.QueryServer = Querier{}

// Params queries checkpoint params
func (k Querier) Params(c context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryParamsResponse{}, nil
}

// AckCount queries ack-count
func (k Querier) AckCount(c context.Context, req *types.QueryAckCountRequest) (*types.QueryAckCountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryAckCountResponse{}, nil
}

// Checkpoint queries checkpoint
func (k Querier) Checkpoint(c context.Context, req *types.QueryCheckpointRequest) (*types.QueryCheckpointResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryCheckpointResponse{}, nil
}

// CheckpointBuffer queries checkpoint buffer
func (k Querier) CheckpointBuffer(c context.Context, req *types.QueryCheckpointBufferRequest) (*types.QueryCheckpointBufferResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryCheckpointBufferResponse{}, nil
}

// LastNoAck queries last no-ack
func (k Querier) LastNoAck(c context.Context, req *types.QueryLastNoAckRequest) (*types.QueryLastNoAckResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryLastNoAckResponse{}, nil
}

// CheckpointList queries list of queries
func (k Querier) CheckpointList(c context.Context, req *types.QueryCheckpointListRequest) (*types.QueryCheckpointListResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryCheckpointListResponse{}, nil
}

// NextCheckpoint queries next checkpoint
func (k Querier) NextCheckpoint(c context.Context, req *types.QueryNextCheckpointRequest) (*types.QueryNextCheckpointResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	// ctx := sdk.UnwrapSDKContext(c)

	return &types.QueryNextCheckpointResponse{}, nil
}
