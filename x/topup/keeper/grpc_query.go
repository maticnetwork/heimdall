package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/x/topup/types"
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

// Validator queries validator info for given validator addr
func (k Querier) Sequence(c context.Context, req *types.QuerySequenceRequest) (*types.QuerySequenceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	txHash := req.TxHash
	logIndex := req.LogIndex
	ctx := sdk.UnwrapSDKContext(c)

	seq, found := k.GetTopupSequenceFromTxHashLogIndex(ctx, txHash, logIndex)
	if !found {
		return nil, status.Errorf(codes.NotFound, "Sequence with tx hash: %s and log index: %d not found", txHash, logIndex)
	}
	return &types.QuerySequenceResponse{Sequence: &seq}, nil
}

// // ValidatorSet queries validatorSet info
// func (k Querier) ValidatorSet(c context.Context, req *types.QueryValidatorSetRequest) (*types.QueryValidatorSetResponse, error) {
// 	if req == nil {
// 		return nil, status.Error(codes.InvalidArgument, "empty request")
// 	}

// 	ctx := sdk.UnwrapSDKContext(c)
// 	validatorSet := k.GetValidatorSet(ctx)

// 	return &types.QueryValidatorSetResponse{ValidatorSet: validatorSet}, nil
// }
