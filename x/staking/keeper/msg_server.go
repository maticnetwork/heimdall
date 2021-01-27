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
	ctx := sdk.UnwrapSDKContext(goCtx)

	k.Logger(ctx).Debug("✅ Validating validator join msg",
		"validatorId", msg.ID,
		"activationEpoch", msg.ActivationEpoch,
		"amount", msg.Amount,
		"SignerPubkey", msg.SignerPubKey,
		"txHash", msg.TxHash,
		"logIndex", msg.LogIndex,
		"blockNumber", msg.BlockNumber,
	)

	// Generate PubKey from Pubkey in message and signer
	pubkey := msg.GetSignerPubKey()
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
		k.Logger(ctx).Error("Error occurred while converting amount to power", "error", err)
		return nil, hmCommon.ErrInvalidPower
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
	ctx := sdk.UnwrapSDKContext(goCtx)
	k.Logger(ctx).Debug("✅ Validating stake update msg",
		"validatorID", msg.ID,
		"newAmount", msg.NewAmount,
		"txHash", msg.TxHash,
		"logIndex", msg.LogIndex,
		"blockNumber", msg.BlockNumber,
	)

	// pull validator from store
	_, ok := k.GetValidatorFromValID(ctx, msg.ID)
	if !ok {
		k.Logger(ctx).Error("Fetching of validator from store failed", "validatorId", msg.ID)
		return nil, hmCommon.ErrNoValidator
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

	// pull validator from store
	validator, ok := k.GetValidatorFromValID(ctx, msg.ID)
	if !ok {
		k.Logger(ctx).Error("Fetching of validator from store failed", "validatorId", msg.ID)
		return nil, hmCommon.ErrNoValidator
	}

	if msg.Nonce != validator.Nonce+1 {
		k.Logger(ctx).Error("Incorrect validator nonce")
		return nil, hmCommon.ErrNonce
	}

	// set validator amount
	_, err := helper.GetPowerFromAmount(msg.NewAmount.BigInt())
	if err != nil {
		k.Logger(ctx).Error("Invalid newamount", msg.NewAmount, "for validator %v", msg.ID)
		return nil, hmCommon.ErrInvalidMsg
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeStakeUpdate,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(validator.ID.Uint64(), 10)),
			sdk.NewAttribute(types.AttributeKeyValidatorNonce, strconv.FormatUint(msg.Nonce, 10)),
		),
	})

	return &types.MsgStakeUpdateResponse{}, nil
}

func (k msgServer) SignerUpdate(goCtx context.Context, msg *types.MsgSignerUpdate) (*types.MsgSignerUpdateResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	k.Logger(ctx).Debug("✅ Validating signer update msg",
		"validatorID", msg.ID,
		"NewSignerPubkey", msg.NewSignerPubKey,
		"txHash", msg.TxHash,
		"logIndex", msg.LogIndex,
		"blockNumber", msg.BlockNumber,
	)

	pubkey := msg.GetNewSignerPubKey()
	newSigner := pubkey.Address()

	// pull validator from store
	validator, ok := k.GetValidatorFromValID(ctx, msg.ID)
	if !ok {
		k.Logger(ctx).Error("Fetching of validator from store failed", "validatorId", msg.ID)
		return nil, hmCommon.ErrNoValidator
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

	// check if new signer address is same as existing signer
	if bytes.Equal(newSigner.Bytes(), []byte(validator.Signer)) {
		// No signer change
		k.Logger(ctx).Error("NewSigner same as OldSigner.")
		return nil, hmCommon.ErrNoSignerChange
	}

	// check nonce validity
	if msg.Nonce != validator.Nonce+1 {
		k.Logger(ctx).Error("Incorrect validator nonce")
		return nil, hmCommon.ErrNonce
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSignerUpdate,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(validator.ID.Uint64(), 10)),
			sdk.NewAttribute(types.AttributeKeyValidatorNonce, strconv.FormatUint(msg.Nonce, 10)),
		),
	})

	return &types.MsgSignerUpdateResponse{}, nil

}

func (k msgServer) ValidatorExit(goCtx context.Context, msg *types.MsgValidatorExit) (*types.MsgValidatorExitResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	k.Logger(ctx).Debug("✅ Validating validator exit msg",
		"validatorID", msg.ID,
		"deactivatonEpoch", msg.DeactivationEpoch,
		"txHash", msg.TxHash,
		"logIndex", msg.LogIndex,
		"blockNumber", msg.BlockNumber,
	)

	validator, ok := k.GetValidatorFromValID(ctx, msg.ID)
	if !ok {
		k.Logger(ctx).Error("Fetching of validator from store failed", "validatorID", msg.ID)
		return nil, hmCommon.ErrNoValidator
	}

	k.Logger(ctx).Debug("validator in store", "validator", validator)
	// check if validator deactivation period is set
	if validator.EndEpoch != 0 {
		k.Logger(ctx).Error("Validator already unbonded")
		return nil, hmCommon.ErrValUnbonded
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

	// check nonce validity
	if msg.Nonce != validator.Nonce+1 {
		k.Logger(ctx).Error("Incorrect validator nonce")
		return nil, hmCommon.ErrNonce
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeValidatorExit,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(validator.ID.Uint64(), 10)),
			sdk.NewAttribute(types.AttributeKeyValidatorNonce, strconv.FormatUint(msg.Nonce, 10)),
		),
	})

	return &types.MsgValidatorExitResponse{}, nil
}
