package keeper

import (
	"context"

	checkpointTypes "github.com/maticnetwork/heimdall/x/checkpoint/types"

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

	txHash := req.GetTxHash()
	logIndex := req.GetLogIndex()
	ctx := sdk.UnwrapSDKContext(c)

	chainParams := k.ChainKeeper.GetParams(ctx)
	receipt, err := k.contractCaller.GetConfirmedTxReceipt(hmTypes.HexToHeimdallHash(txHash).EthHash(), chainParams.MainchainTxConfirmations)

	if err != nil || receipt == nil {
		return nil, status.Errorf(codes.NotFound, "Transaction is not confirmed yet. Please wait for sometime and try again")
	}

	sequence := new(big.Int).Mul(receipt.BlockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(logIndex))

	if !k.HasTopupSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("No sequence exists: %s %s", txHash, logIndex)
		return nil, status.Errorf(codes.NotFound, "Sequence not found")
	}

	return &types.QuerySequenceResponse{Sequence: sequence.Uint64()}, nil
}

// IsOldTx will returns boolean if tx is old or not
func (k Querier) IsOldTx(c context.Context, req *types.QueryIsOldTxSequenceRequest) (*types.QueryIsOldTxSequenceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	_, err := k.Sequence(c, &types.QuerySequenceRequest{
		TxHash:   req.GetTxHash(),
		LogIndex: req.GetLogIndex(),
	})

	if err != nil {
		return nil, err
	}

	return &types.QueryIsOldTxSequenceResponse{Status: true}, nil
}

// QueryDividendAccountRoot will return dividend account root hash
func (k Querier) QueryDividendAccountRoot(c context.Context, req *types.QueryDividendAccountRootRequest) (*types.QueryDividendAccountRootResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)

	dividendAccounts := k.GetAllDividendAccounts(ctx)
	accountRoot, err := checkpointTypes.GetAccountRootHash(dividendAccounts)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "could not fetch accountroothash")
	}

	return &types.QueryDividendAccountRootResponse{
		AccountRootHash: hmTypes.BytesToHeimdallHash(accountRoot).String(),
	}, nil
}

// QueryDividendAccounts will return all dividend account
func (k Querier) QueryDividendAccounts(c context.Context, req *types.QueryDividendAccountsRequest) (*types.QueryDividendAccountsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}
	ctx := sdk.UnwrapSDKContext(c)
	dividendAccounts := k.GetAllDividendAccounts(ctx)
	return &types.QueryDividendAccountsResponse{DividendAccounts: dividendAccounts}, nil
}

// QueryDividendAccount will return dividend account info with given addr
func (k Querier) QueryDividendAccount(c context.Context, req *types.QueryDividendAccountRequest) (*types.QueryDividendAccountResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	reqAddr := req.GetAddress()
	if reqAddr == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid address format")
	}

	addr, _ := sdk.AccAddressFromHex(reqAddr)
	dividendAccount, err := k.GetDividendAccountByAddress(ctx, addr)
	if err != nil {
		return nil, status.Error(codes.NotFound, "dividend account not found")
	}

	return &types.QueryDividendAccountResponse{
		DividendAccount: &dividendAccount,
	}, nil
}
