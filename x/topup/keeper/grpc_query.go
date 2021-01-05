package keeper

import (
	"context"

	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/x/topup/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

	chainParams := k.chainKeeper.GetParams(ctx)
	receipt, err := k.contractCaller.GetConfirmedTxReceipt(hmTypes.HexToHeimdallHash(txHash).EthHash(), chainParams.MainchainTxConfirmations)

	if err != nil || receipt == nil {
		return nil, status.Errorf(codes.NotFound, "Transaction is not confirmed yet. Please wait for sometime and try again")
	}

	sequence := new(big.Int).Mul(receipt.BlockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(logIndex))

	if !k.HasTopupSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("No sequence exists: %s %s", txHash, logIndex)
		return nil, nil
	}

	return &types.QuerySequenceResponse{Sequence: sequence.Uint64()}, nil
}
