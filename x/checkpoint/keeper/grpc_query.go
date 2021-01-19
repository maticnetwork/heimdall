package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"
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

	ctx := sdk.UnwrapSDKContext(c)

	params := k.GetParams(ctx)

	return &types.QueryParamsResponse{Params: params}, nil
}

// AckCount queries ack-count
func (k Querier) AckCount(c context.Context, req *types.QueryAckCountRequest) (*types.QueryAckCountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	ackCount := k.GetACKCount(ctx)

	return &types.QueryAckCountResponse{AckCount: ackCount}, nil
}

// Checkpoint queries checkpoint
func (k Querier) Checkpoint(c context.Context, req *types.QueryCheckpointRequest) (*types.QueryCheckpointResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	if req.Number == 0 {
		return nil, status.Error(codes.InvalidArgument, "empty header param")
	}

	ctx := sdk.UnwrapSDKContext(c)

	res, err := k.GetCheckpointByNumber(ctx, req.Number)
	if err != nil {
		return nil, types.ErrNoCheckpointFound
	}

	return &types.QueryCheckpointResponse{Checkpoint: &res}, nil
}

// CheckpointBuffer queries checkpoint buffer
func (k Querier) CheckpointBuffer(c context.Context, req *types.QueryCheckpointBufferRequest) (*types.QueryCheckpointBufferResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	res, err := k.GetCheckpointFromBuffer(ctx)
	if err != nil {
		return nil, types.ErrNoCheckpointBufferFound
	}

	if res == nil {
		return nil, types.ErrNoCheckpointBufferFound
	}

	return &types.QueryCheckpointBufferResponse{CheckpointBuffer: res}, nil
}

// LastNoAck queries last no-ack
func (k Querier) LastNoAck(c context.Context, req *types.QueryLastNoAckRequest) (*types.QueryLastNoAckResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	res := k.GetLastNoAck(ctx)

	return &types.QueryLastNoAckResponse{LastNoAck: res}, nil
}

// CheckpointList queries list of queries
func (k Querier) CheckpointList(c context.Context, req *types.QueryCheckpointListRequest) (*types.QueryCheckpointListResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	if req.Pagination == nil {
		return nil, status.Error(codes.InvalidArgument, "empty pagination limit, page params")
	}

	res, err := k.GetCheckpointList(ctx, req.Pagination.Page, req.Pagination.Limit)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "nocheckpoint list error")
	}
	return &types.QueryCheckpointListResponse{CheckpointList: res}, nil
}

// NextCheckpoint queries next checkpoint
func (k Querier) NextCheckpoint(c context.Context, req *types.QueryNextCheckpointRequest) (*types.QueryNextCheckpointResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	if req.BorChainID == "" {
		return nil, status.Error(codes.InvalidArgument, "empty borchain id")
	}

	borChainID := req.BorChainID
	ctx := sdk.UnwrapSDKContext(c)

	// get validator set
	validatorSet := k.Sk.GetValidatorSet(ctx)
	proposer := validatorSet.GetProposer()
	ackCount := k.GetACKCount(ctx)
	params := k.GetParams(ctx)

	var start uint64

	if ackCount != 0 {
		checkpointNumber := ackCount
		lastCheckpoint, err := k.GetCheckpointByNumber(ctx, checkpointNumber)
		if err != nil {
			return nil, err
		}
		start = lastCheckpoint.EndBlock + 1
	}

	end := start + params.AvgCheckpointLength

	rootHash, err := k.contractCaller.GetRootHash(start, end, params.MaxCheckpointLength)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "root has error")
	}

	// accs := k.tk.GetAllDividendAccounts(ctx)
	// accRootHash, err := types.GetAccountRootHash(accs)
	// if err != nil {
	// 	return nil, sdk.ErrInternal(sdk.AppendMsgToErr(fmt.Sprintf("could not get generate account root hash. Error:%v", err), err.Error()))
	// }

	checkpointMsg := types.NewMsgCheckpointBlock(
		sdk.AccAddress([]byte(proposer.Signer)),
		start,
		start+params.AvgCheckpointLength,
		hmCommonTypes.BytesToHeimdallHash(rootHash),
		hmCommonTypes.BytesToHeimdallHash(rootHash), //hmTypes.BytesToHeimdallHash(accRootHash),
		borChainID,
	)

	return &types.QueryNextCheckpointResponse{NextCheckpoint: &checkpointMsg}, nil
}
