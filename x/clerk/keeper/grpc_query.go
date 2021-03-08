package keeper

import (
	"context"
	"math/big"
	"time"

	"github.com/jinzhu/copier"

	hmTypes "github.com/maticnetwork/heimdall/types/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/x/clerk/types"
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

func (k Querier) Record(c context.Context, req *types.QueryRecordParams) (*types.QueryRecordResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	// get state record by record id
	record, err := k.GetEventRecord(ctx, req.RecordId)
	if err != nil {
		return nil, err
	}

	return &types.QueryRecordResponse{EventRecord: record}, nil
}

// QueryIsOldTxClerk will returns bool if isoldtx or not
func (k Querier) QueryIsOldTxClerk(c context.Context, req *types.QueryIsOldTxRequest) (*types.QueryIsOldTxResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)

	txHash := req.GetTxHash()
	logIndex := req.GetLogIndex()
	chainParams := k.ChainKeeper.GetParams(ctx)
	receipt, err := k.contractCaller.GetConfirmedTxReceipt(hmTypes.HexToHeimdallHash(txHash).EthHash(), chainParams.MainchainTxConfirmations)

	if err != nil || receipt == nil {
		return nil, status.Errorf(codes.NotFound, "Transaction is not confirmed yet. Please wait for sometime and try again")
	}

	// sequence id

	sequence := new(big.Int).Mul(receipt.BlockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(logIndex))

	// check if incoming tx already exists
	if !k.HasRecordSequence(ctx, sequence.String()) {
		return nil, status.Errorf(codes.NotFound, "Sequence not found")
	}

	return &types.QueryIsOldTxResponse{Status: true}, nil
}

// Event Records List
func (k Querier) Records(c context.Context, req *types.QueryRecordListRequest) (*types.QueryRecordListResponse, error) {
	var records []types.EventRecord
	var err error

	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	page := uint64(1)
	if req.Page != 0 {
		page = req.Page
	}

	limit := uint64(50)
	if req.Limit != 0 {
		limit = req.Limit
	}

	ctx := sdk.UnwrapSDKContext(c)

	if req.FromTime != 0 && req.ToTime != 0 {
		records, err = k.GetEventRecordListWithTime(ctx, time.Unix(int64(req.FromTime), 0), time.Unix(int64(req.ToTime), 0), page, limit)
		if err != nil {
		return nil, err
	}
	} else if (req.FromId != 0) && (req.ToTime != 0) {
		fromRecord, err := k.GetEventRecord(ctx, req.FromId)
		if err != nil && err.Error() != "No record found" {
			return nil, err
		}
		if fromRecord != nil {
			fromTime := fromRecord.RecordTime.Unix()
		records, err = k.GetEventRecordListWithTime(ctx, time.Unix(fromTime, 0), time.Unix(int64(req.ToTime), 0), 1, limit)
		if err != nil {
			return nil, err
		}
		}

	} else {
		records, err = k.GetEventRecordList(ctx, page, limit)
		if err != nil {
		return nil, err
	}
	}

	var ptrRecords []*types.EventRecord
	for _, record := range records {
		newRecord := types.EventRecord{}
		err := copier.Copy(&newRecord, &record)
		if err != nil {
			return nil, err
		}
		ptrRecords = append(ptrRecords, &newRecord)
	}
	return &types.QueryRecordListResponse{
		EventRecords: ptrRecords,
	}, nil
}
