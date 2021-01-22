package keeper

import (
	"context"
	"encoding/hex"
	"math/big"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/x/clerk/types"
)

type msgServer struct {
	Keeper
	contractCaller helper.IContractCaller
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper, contractCaller helper.IContractCaller) types.MsgServer {
	return &msgServer{Keeper: keeper, contractCaller: contractCaller}
}

var _ types.MsgServer = msgServer{}

func (k msgServer) MsgEventRecord(goCtx context.Context, msg *types.MsgEventRecordRequest) (*types.MsgEventRecordResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	k.Logger(ctx).Debug("✅ Validating clerk msg",
		"id", msg.Id,
		"contract", msg.ContractAddress,
		"data", hex.EncodeToString(msg.Data),
		"txHash", msg.TxHash,
		"logIndex", uint64(msg.LogIndex),
		"blockNumber", msg.BlockNumber,
	)

	// check if event record exists
	if exists := k.HasEventRecord(ctx, msg.Id); exists {
		return nil, hmCommon.ErrEventRecordAlreadySynced
	}

	// TODO - Check this
	// chainManager params
	// params := k.ChainKeeper.GetParams(ctx)
	// chainParams := params.ChainParams

	// check chain id
	// if chainParams.BorChainID != msg.ChainId {
	// 	k.Logger(ctx).Error("Invalid Bor chain id", "msgChainID", msg.ChainId)
	// 	return nil, hmCommon.ErrInvalidBorChainID
	// }

	// sequence id
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// check if incoming tx is older
	if k.HasRecordSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return nil, hmCommon.ErrOldTx
	}

	// add events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeRecord,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyRecordID, strconv.FormatUint(msg.Id, 10)),
			sdk.NewAttribute(types.AttributeKeyRecordContract, msg.ContractAddress),
			sdk.NewAttribute(types.AttributeKeyRecordTxHash, msg.TxHash),
			sdk.NewAttribute(types.AttributeKeyRecordTxLogIndex, strconv.FormatUint(msg.LogIndex, 10)),
		),
	})

	return &types.MsgEventRecordResponse{}, nil
}
