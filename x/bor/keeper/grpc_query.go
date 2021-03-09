package keeper

import (
	"context"
	"fmt"

	"github.com/maticnetwork/heimdall/helper"

	hmTypes "github.com/maticnetwork/heimdall/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/maticnetwork/heimdall/x/bor/types"
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

const (
	ParamSpan          = "span"
	ParamSprint        = "sprint"
	ParamProducerCount = "producer-count"
	ParamLastEthBlock  = "last-eth-block"
)

// Params returns all bor params info
func (k Querier) Params(context context.Context, req *types.QueryParamsRequest) (*types.QueryParamsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(context)
	getParams := k.GetParams(ctx)
	latestEthBlock := k.GetLastEthBlock(ctx)
	return &types.QueryParamsResponse{
		SpanDuration:   getParams.GetSpanDuration(),
		LatestEthBlock: latestEthBlock.Uint64(),
		ProducerCount:  getParams.GetProducerCount(),
		Sprint:         getParams.GetSprintDuration(),
	}, nil
}

// Param returns bor parameters info
func (k Querier) Param(context context.Context, req *types.QueryParamRequest) (*types.QueryParamResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(context)
	switch req.GetParamsType() {
	case ParamSpan:
		params := k.GetParams(ctx)
		return &types.QueryParamResponse{
			Params: &types.QueryParamResponse_SpanDuration{
				SpanDuration: params.SpanDuration,
			},
		}, nil
	case ParamSprint:
		params := k.GetParams(ctx)
		return &types.QueryParamResponse{
			Params: &types.QueryParamResponse_Sprint{
				Sprint: params.SprintDuration,
			},
		}, nil
	case ParamProducerCount:
		params := k.GetParams(ctx)
		return &types.QueryParamResponse{
			Params: &types.QueryParamResponse_ProducerCount{
				ProducerCount: params.ProducerCount,
			},
		}, nil
	case ParamLastEthBlock:
		latestEthBlock := k.GetLastEthBlock(ctx)
		return &types.QueryParamResponse{
			Params: &types.QueryParamResponse_LatestEthBlock{
				LatestEthBlock: latestEthBlock.Uint64(),
			},
		}, nil
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid param type ")

	}
}

// Span returns span info with span-id
func (k Querier) Span(goCtx context.Context, req *types.QuerySpanRequest) (*types.QuerySpanResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	resp, err := k.GetSpan(ctx, req.GetSpanId())
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, status.Error(codes.NotFound, "span not found for id")
	}
	return &types.QuerySpanResponse{
		Span: resp,
	}, nil
}

func (k Querier) SpanList(context context.Context, req *types.QuerySpanListRequest) (*types.QuerySpanListResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(context)
	resp, err := k.GetSpanList(ctx, req.Page, req.Limit)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if resp == nil {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("could not fetch span list with page %v and limit %v", req.Page, req.Limit))
	}
	return &types.QuerySpanListResponse{
		Spans: resp,
	}, nil
}

func (k Querier) LatestSpan(context context.Context, req *types.QueryLatestSpanRequest) (*types.QueryLatestSpanResponse, error) {
	var defaultSpan *hmTypes.Span
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(context)
	spans, err := k.GetAllSpans(ctx)
	if err != nil {
		return nil, err
	}
	// if this is the first span return empty span
	if len(spans) == 0 {
		return &types.QueryLatestSpanResponse{Span: defaultSpan}, nil
	}
	// explicitly fetch the last span
	span, err := k.GetLastSpan(ctx)
	if err != nil {
		return nil, err
	}
	if span == nil {
		return nil, status.Error(codes.NotFound, "latest span does not exist")
	}
	return &types.QueryLatestSpanResponse{Span: span}, nil
}

func (k Querier) NextSpanSeed(context context.Context, req *types.QueryNextSpanSeedRequest) (*types.QueryNextSpanSeedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(context)
	nextSpanSeed, err := k.GetNextSpanSeed(ctx, k.contractCaller)
	if err != nil {
		return nil, err
	}
	return &types.QueryNextSpanSeedResponse{
		NextSpanSeed: nextSpanSeed.String(),
	}, nil
}

func (k Querier) PrepareNextSpan(context context.Context, req *types.PrepareNextSpanRequest) (*types.PrepareNextSpanResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	// get the query params
	spanId := req.GetSpanId()
	startBlock := req.GetStartBlock()
	chainId := req.GetBorChainId()

	if len(chainId) == 0 {
		return nil, status.Error(codes.InvalidArgument, "invalid request, chain_id required")
	}

	ctx := sdk.UnwrapSDKContext(context)
	// Get the span duration
	params := k.GetParams(ctx)
	spanDuration := params.SpanDuration

	//ackCount := k.checkpointKeeper.GetACKCount(ctx)
	//if ackCount == 0 {
	//	return nil, status.Errorf(codes.NotFound, "Ack not found")
	//}

	currentValidatorSet := k.sk.GetValidatorSet(ctx)
	if currentValidatorSet == nil {
		return nil, status.Errorf(codes.NotFound, "validator set not found")
	}

	nextSpanSeed, err := k.GetNextSpanSeed(ctx, k.contractCaller)
	if err != nil {
		return nil, err
	}
	nextProducers, err := k.SelectNextProducers(ctx, nextSpanSeed)
	if err != nil {
		return nil, err
	}

	selectedProducers := hmTypes.SortValidatorByAddress(nextProducers)

	// creat new span
	newSpan := hmTypes.NewSpan(
		spanId,
		startBlock,
		startBlock+spanDuration-1,
		hmTypes.ValidatorSet{
			Validators:       currentValidatorSet.Validators,
			Proposer:         currentValidatorSet.Proposer,
			TotalVotingPower: currentValidatorSet.TotalVotingPower,
		},
		selectedProducers,
		chainId,
	)

	return &types.PrepareNextSpanResponse{
		Span: &newSpan,
	}, nil
}
