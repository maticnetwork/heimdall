package staking

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// NewHandler new handler
func NewHandler(k Keeper, contractCaller helper.IContractCaller) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgValidatorJoin:
			return HandleMsgValidatorJoin(ctx, msg, k, contractCaller)
		case types.MsgValidatorExit:
			return HandleMsgValidatorExit(ctx, msg, k, contractCaller)
		case types.MsgSignerUpdate:
			return HandleMsgSignerUpdate(ctx, msg, k, contractCaller)
		case types.MsgStakeUpdate:
			return HandleMsgStakeUpdate(ctx, msg, k, contractCaller)
		default:
			return sdk.ErrTxDecode("Invalid message in staking module").Result()
		}
	}
}

// HandleMsgValidatorJoin msg validator join
func HandleMsgValidatorJoin(ctx sdk.Context, msg types.MsgValidatorJoin, k Keeper, contractCaller helper.IContractCaller) sdk.Result {

	k.Logger(ctx).Debug("✅ Validating validator join msg",
		"validatorId", msg.ID,
		"activationEpoch", msg.ActivationEpoch,
		"amount", msg.Amount,
		"SignerPubkey", msg.SignerPubKey.String(),
		"txHash", msg.TxHash,
		"logIndex", msg.LogIndex,
		"blockNumber", msg.BlockNumber,
	)

	// Generate PubKey from Pubkey in message and signer
	pubkey := msg.SignerPubKey
	signer := pubkey.Address()

	// Check if validator has been validator before
	if _, ok := k.GetSignerFromValidatorID(ctx, msg.ID); ok {
		k.Logger(ctx).Error("Validator has been validator before, cannot join with same ID", "validatorId", msg.ID)
		return hmCommon.ErrValidatorAlreadyJoined(k.Codespace()).Result()
	}

	// get validator by signer
	checkVal, err := k.GetValidatorInfo(ctx, signer.Bytes())
	if err == nil || bytes.Equal(checkVal.Signer.Bytes(), signer.Bytes()) {
		return hmCommon.ErrValidatorAlreadyJoined(k.Codespace()).Result()
	}

	// validate voting power
	_, err = helper.GetPowerFromAmount(msg.Amount.BigInt())
	if err != nil {
		return hmCommon.ErrInvalidMsg(k.Codespace(), fmt.Sprintf("Invalid amount %v for validator %v", msg.Amount, msg.ID)).Result()
	}

	// sequence id
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// check if incoming tx is older
	if k.HasStakingSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	// Emit event join
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeValidatorJoin,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(msg.ID.Uint64(), 10)),
			sdk.NewAttribute(types.AttributeKeyValidatorNonce, strconv.FormatUint(msg.Nonce, 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// HandleMsgStakeUpdate handles stake update message
func HandleMsgStakeUpdate(ctx sdk.Context, msg types.MsgStakeUpdate, k Keeper, contractCaller helper.IContractCaller) sdk.Result {

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
		return hmCommon.ErrNoValidator(k.Codespace()).Result()
	}

	// sequence id
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// check if incoming tx is older
	if k.HasStakingSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	// pull validator from store
	validator, ok := k.GetValidatorFromValID(ctx, msg.ID)
	if !ok {
		k.Logger(ctx).Error("Fetching of validator from store failed", "validatorId", msg.ID)
		return hmCommon.ErrNoValidator(k.Codespace()).Result()
	}

	if msg.Nonce != validator.Nonce+1 {
		k.Logger(ctx).Error("Incorrect validator nonce")
		return hmCommon.ErrNonce(k.Codespace()).Result()
	}

	// set validator amount
	_, err := helper.GetPowerFromAmount(msg.NewAmount.BigInt())
	if err != nil {
		return hmCommon.ErrInvalidMsg(k.Codespace(), fmt.Sprintf("Invalid newamount %v for validator %v", msg.NewAmount, msg.ID)).Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeStakeUpdate,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(validator.ID.Uint64(), 10)),
			sdk.NewAttribute(types.AttributeKeyValidatorNonce, strconv.FormatUint(msg.Nonce, 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// HandleMsgSignerUpdate handles signer update message
func HandleMsgSignerUpdate(ctx sdk.Context, msg types.MsgSignerUpdate, k Keeper, contractCaller helper.IContractCaller) sdk.Result {

	k.Logger(ctx).Debug("✅ Validating signer update msg",
		"validatorID", msg.ID,
		"NewSignerPubkey", msg.NewSignerPubKey.String(),
		"txHash", msg.TxHash,
		"logIndex", msg.LogIndex,
		"blockNumber", msg.BlockNumber,
	)

	newPubKey := msg.NewSignerPubKey
	newSigner := newPubKey.Address()

	// pull validator from store
	validator, ok := k.GetValidatorFromValID(ctx, msg.ID)
	if !ok {
		k.Logger(ctx).Error("Fetching of validator from store failed", "validatorId", msg.ID)
		return hmCommon.ErrNoValidator(k.Codespace()).Result()
	}

	// sequence id
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// check if incoming tx is older
	if k.HasStakingSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	// check if new signer address is same as existing signer
	if bytes.Equal(newSigner.Bytes(), validator.Signer.Bytes()) {
		// No signer change
		k.Logger(ctx).Error("NewSigner same as OldSigner.")
		return hmCommon.ErrNoSignerChange(k.Codespace()).Result()
	}

	// check nonce validity
	if msg.Nonce != validator.Nonce+1 {
		k.Logger(ctx).Error("Incorrect validator nonce")
		return hmCommon.ErrNonce(k.Codespace()).Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSignerUpdate,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(validator.ID.Uint64(), 10)),
			sdk.NewAttribute(types.AttributeKeyValidatorNonce, strconv.FormatUint(msg.Nonce, 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// HandleMsgValidatorExit handle msg validator exit
func HandleMsgValidatorExit(ctx sdk.Context, msg types.MsgValidatorExit, k Keeper, contractCaller helper.IContractCaller) sdk.Result {

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
		return hmCommon.ErrNoValidator(k.Codespace()).Result()
	}

	k.Logger(ctx).Debug("validator in store", "validator", validator)
	// check if validator deactivation period is set
	if validator.EndEpoch != 0 {
		k.Logger(ctx).Error("Validator already unbonded")
		return hmCommon.ErrValUnbonded(k.Codespace()).Result()
	}

	// sequence id
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// check if incoming tx is older
	if k.HasStakingSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	// check nonce validity
	if msg.Nonce != validator.Nonce+1 {
		k.Logger(ctx).Error("Incorrect validator nonce")
		return hmCommon.ErrNonce(k.Codespace()).Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeValidatorExit,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(validator.ID.Uint64(), 10)),
			sdk.NewAttribute(types.AttributeKeyValidatorNonce, strconv.FormatUint(msg.Nonce, 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
