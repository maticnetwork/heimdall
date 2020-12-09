package keeper

import (
	"bytes"
	"context"
	"math/big"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/x/staking/types"
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

func (k msgServer) ValidatorJoin(goCtx context.Context, msg *types.MsgValidatorJoin) (*types.MsgValidatorJoinResponse, error) {
	// k.Logger(ctx).Info("Handling new validator join", "msg", msg)

	ctx := sdk.UnwrapSDKContext(goCtx)

	// Generate PubKey from Pubkey in message and signer
	pubkey := msg.SignerPubKey
	signer := pubkey.Address()

	// Check if validator has been validator before
	if _, ok := k.GetSignerFromValidatorID(ctx, msg.ID); ok {
		k.Logger(ctx).Error("Validator has been validator before, cannot join with same ID", "validatorId", msg.ID)
		return nil, hmCommon.ErrValidatorAlreadyJoined
	}

	// get validator by signer
	checkVal, err := k.GetValidatorInfo(ctx, signer.Bytes())
	if err == nil || bytes.Equal([]byte(checkVal.Signer), signer.Bytes()) {
		return nil, hmCommon.ErrValidatorAlreadyJoined
	}

	// get voting power from amount
	_, err = helper.GetPowerFromAmount(msg.Amount.BigInt())
	if err != nil {
		return nil, hmCommon.ErrInvalidMsg
	}

	// sequence id
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// check if incoming tx is older
	if k.HasStakingSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return nil, hmCommon.ErrOldTx
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeValidatorJoin,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(msg.ID.Uint64(), 10)),
			sdk.NewAttribute(types.AttributeKeySigner, strconv.FormatUint(msg.Nonce, 10)),
		),
	})

	return &types.MsgValidatorJoinResponse{}, nil

}

func (k msgServer) StakeUpdate(goCtx context.Context, msg *types.MsgStakeUpdate) (*types.MsgStakeUpdateResponse, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.MsgStakeUpdateResponse{}, nil
}

func (k msgServer) SignerUpdate(goCtx context.Context, msg *types.MsgSignerUpdate) (*types.MsgSignerUpdateResponse, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.MsgSignerUpdateResponse{}, nil

}

func (k msgServer) ValidatorExit(goCtx context.Context, msg *types.MsgValidatorExit) (*types.MsgValidatorExitResponse, error) {
	// ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.MsgValidatorExitResponse{}, nil
}
