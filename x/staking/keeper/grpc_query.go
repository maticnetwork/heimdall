package keeper

import (
	"context"
	"fmt"
	"math/big"

	"github.com/maticnetwork/bor/common"

	hmTypes "github.com/maticnetwork/heimdall/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/x/staking/types"
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
func (k Querier) Validator(c context.Context, req *types.QueryValidatorRequest) (*types.QueryValidatorResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	validatorID := hmTypes.ValidatorID(req.ValidatorId)
	if req.ValidatorId == 0 {
		return nil, status.Error(codes.InvalidArgument, "validator ID cannot be empty")
	}

	ctx := sdk.UnwrapSDKContext(c)
	validator, found := k.GetValidatorFromValID(ctx, hmTypes.ValidatorID(validatorID))
	if !found {
		return nil, status.Errorf(codes.NotFound, "validator %d not found", req.ValidatorId)
	}

	return &types.QueryValidatorResponse{Validator: &validator}, nil
}

// ValidatorSet queries validatorSet info
func (k Querier) ValidatorSet(c context.Context, req *types.QueryValidatorSetRequest) (*types.QueryValidatorSetResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	validatorSet := k.GetValidatorSet(ctx)

	return &types.QueryValidatorSetResponse{ValidatorSet: validatorSet}, nil
}

// StakingOldTx returns the tx is old or not with given txhash and logindex
func (k Querier) StakingOldTx(c context.Context, req *types.QueryStakingOldTxRequest) (*types.QueryStakingOldTxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	params := k.Keeper.ChainKeeper.GetParams(ctx)

	// get main tx receipt
	receipt, err := k.contractCaller.GetConfirmedTxReceipt(common.HexToHash(req.GetTxHash()), params.MainchainTxConfirmations)
	if err != nil || receipt == nil {
		return nil, status.Error(codes.Internal, "Transaction is not confirmed yet. Please wait for sometime and try again")
	}

	// sequence id
	sequence := new(big.Int).Mul(receipt.BlockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(req.GetLogIndex()))

	if !k.Keeper.HasStakingSequence(ctx, sequence.String()) {
		return nil, status.Error(codes.NotFound, fmt.Sprintf("No staking sequence exist: %s %d", req.GetTxHash(), req.GetLogIndex()))
	}

	return &types.QueryStakingOldTxResponse{
		Status: true,
	}, nil
}

// QueryProposer will the proposers list
func (k Querier) QueryProposer(c context.Context, req *types.QueryProposerRequest) (*types.QueryProposerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	// get validator set
	validatorSet := k.GetValidatorSet(ctx)

	times := int(req.GetTimes())
	if times > len(validatorSet.Validators) {
		times = len(validatorSet.Validators)
	}

	// init proposers
	var proposers []*hmTypes.Validator

	// get proposers
	for index := 0; index < times; index++ {
		proposers = append(proposers, validatorSet.GetProposer())
		validatorSet.IncrementProposerPriority(1)
	}

	return &types.QueryProposerResponse{
		Proposers: proposers,
	}, nil
}
