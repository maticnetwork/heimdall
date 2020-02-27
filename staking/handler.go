package staking

import (
	"bytes"
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
			return sdk.ErrTxDecode("Invalid message in checkpoint module").Result()
		}
	}
}

// HandleMsgValidatorJoin msg validator join
func HandleMsgValidatorJoin(ctx sdk.Context, msg types.MsgValidatorJoin, k Keeper, contractCaller helper.IContractCaller) sdk.Result {
	k.Logger(ctx).Info("Handling new validator join", "msg", msg)

	// get main tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(msg.TxHash.EthHash())
	if err != nil || receipt == nil {
		return hmCommon.ErrWaitForConfirmation(k.Codespace()).Result()
	}

	// TODO remove this block in next release and use msg.LogIndex -- start
	logIndex, found := contractCaller.FindStakedEventLogIndex(receipt)
	if !found {
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Unable to find validator join log from txHash").Result()
	}
	// -- END

	// decode validator join event
	eventLog, err := contractCaller.DecodeValidatorJoinEvent(receipt, logIndex)
	if err != nil || eventLog == nil {
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Unable to fetch logs for txHash").Result()
	}

	// Generate PubKey from Pubkey in message and signer
	pubkey := msg.SignerPubKey
	signer := pubkey.Address()

	// check signer in message corresponds
	if !bytes.Equal(signer.Bytes(), eventLog.Signer.Bytes()) {
		k.Logger(ctx).Error(
			"Signer Address does not match",
			"msgValidator", signer.String(),
			"mainchainValidator", eventLog.Signer.Hex(),
		)
		return hmCommon.ErrNoValidator(k.Codespace()).Result()
	}

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

	// get voting power from amount
	votingPower, err := helper.GetPowerFromAmount(eventLog.Amount)
	if err != nil {
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Invalid amount for validator: %v", msg.ID).Result()
	}

	// create new validator
	newValidator := hmTypes.Validator{
		ID:          msg.ID,
		StartEpoch:  eventLog.ActivationEpoch.Uint64(),
		EndEpoch:    0,
		VotingPower: votingPower.Int64(),
		PubKey:      pubkey,
		Signer:      hmTypes.BytesToHeimdallAddress(signer.Bytes()),
		LastUpdated: 0,
	}

	// add validator to store
	k.Logger(ctx).Debug("Adding new validator to state", "validator", newValidator.String())
	err = k.AddValidator(ctx, newValidator)
	if err != nil {
		k.Logger(ctx).Error("Unable to add validator to state", "error", err, "validator", newValidator.String())
		return hmCommon.ErrValidatorSave(k.Codespace()).Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeValidatorJoin,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(newValidator.ID.Uint64(), 10)),
			sdk.NewAttribute(types.AttributeKeySigner, newValidator.Signer.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// HandleMsgStakeUpdate handles stake update message
func HandleMsgStakeUpdate(ctx sdk.Context, msg types.MsgStakeUpdate, k Keeper, contractCaller helper.IContractCaller) sdk.Result {
	k.Logger(ctx).Debug("Handling stake update", "Validator", msg.ID)

	// get main tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(msg.TxHash.EthHash())
	if err != nil || receipt == nil {
		return hmCommon.ErrWaitForConfirmation(k.Codespace()).Result()
	}

	eventLog, err := contractCaller.DecodeValidatorStakeUpdateEvent(receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Unable to fetch logs for txHash").Result()
	}

	if eventLog.ValidatorId.Uint64() != msg.ID.Uint64() {
		k.Logger(ctx).Error("ID in message doesnt match id in logs", "MsgID", msg.ID, "IdFromTx", eventLog.ValidatorId)
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Invalid txhash, id's dont match. Id from tx hash is %v", eventLog.ValidatorId.Uint64()).Result()
	}

	// pull validator from store
	validator, ok := k.GetValidatorFromValID(ctx, msg.ID)
	if !ok {
		k.Logger(ctx).Error("Fetching of validator from store failed", "validatorId", msg.ID)
		return hmCommon.ErrNoValidator(k.Codespace()).Result()
	}

	// sequence id
	sequence := (receipt.BlockNumber.Uint64() * hmTypes.DefaultLogIndexUnit) + msg.LogIndex

	// check if incoming tx is older
	if k.HasStakingSequence(ctx, sequence) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	// update last updated
	validator.LastUpdated = sequence

	// set validator amount
	p, err := helper.GetPowerFromAmount(eventLog.NewAmount)
	if err != nil {
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Invalid amount for validator: %v", msg.ID).Result()
	}
	validator.VotingPower = p.Int64()

	// save validator
	err = k.AddValidator(ctx, validator)
	if err != nil {
		k.Logger(ctx).Error("Unable to update signer", "error", err, "ValidatorID", validator.ID)
		return hmCommon.ErrSignerUpdateError(k.Codespace()).Result()
	}

	// save staking sequence
	k.SetStakingSequence(ctx, sequence)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeStakeUpdate,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(validator.ID.Uint64(), 10)),
			sdk.NewAttribute(types.AttributeKeyUpdatedAt, strconv.FormatUint(validator.LastUpdated, 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// HandleMsgSignerUpdate handles signer update message
func HandleMsgSignerUpdate(ctx sdk.Context, msg types.MsgSignerUpdate, k Keeper, contractCaller helper.IContractCaller) sdk.Result {
	k.Logger(ctx).Debug("Handling signer update", "Validator", msg.ID, "Signer", msg.NewSignerPubKey.Address())

	// get main tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(msg.TxHash.EthHash())
	if err != nil || receipt == nil {
		return hmCommon.ErrWaitForConfirmation(k.Codespace()).Result()
	}

	newPubKey := msg.NewSignerPubKey
	newSigner := newPubKey.Address()

	eventLog, err := contractCaller.DecodeSignerUpdateEvent(receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Unable to fetch signer update log for txHash").Result()
	}

	if eventLog.ValidatorId.Uint64() != msg.ID.Uint64() {
		k.Logger(ctx).Error("ID in message doesnt match id in logs", "MsgID", msg.ID, "IdFromTx", eventLog.ValidatorId.Uint64())
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Invalid txhash, id's dont match. Id from tx hash is %v", eventLog.ValidatorId.Uint64()).Result()
	}

	if bytes.Compare(eventLog.NewSigner.Bytes(), newSigner.Bytes()) != 0 {
		k.Logger(ctx).Error("Signer in txhash and msg dont match", "MsgSigner", newSigner.String(), "SignerTx", eventLog.NewSigner.String())
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Signer in txhash and msg dont match").Result()
	}

	// pull validator from store
	validator, ok := k.GetValidatorFromValID(ctx, msg.ID)
	if !ok {
		k.Logger(ctx).Error("Fetching of validator from store failed", "validatorId", msg.ID)
		return hmCommon.ErrNoValidator(k.Codespace()).Result()
	}
	oldValidator := validator.Copy()

	// sequence id
	sequence := (receipt.BlockNumber.Uint64() * hmTypes.DefaultLogIndexUnit) + msg.LogIndex

	// check if incoming tx is older
	if k.HasStakingSequence(ctx, sequence) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	// update last udpated
	validator.LastUpdated = sequence

	// check if we are actually updating signer
	if !bytes.Equal(newSigner.Bytes(), validator.Signer.Bytes()) {
		// Update signer in prev Validator
		validator.Signer = hmTypes.HeimdallAddress(newSigner)
		validator.PubKey = newPubKey
		k.Logger(ctx).Debug("Updating new signer", "signer", newSigner.String(), "oldSigner", oldValidator.Signer.String(), "validatorID", msg.ID)
	}

	k.Logger(ctx).Debug("Removing old validator", "validator", oldValidator.String())

	// remove old validator from HM
	oldValidator.EndEpoch = k.ackRetriever.GetACKCount(ctx)

	// remove old validator from TM
	oldValidator.VotingPower = 0
	// updated last
	oldValidator.LastUpdated = sequence

	// save old validator
	if err := k.AddValidator(ctx, *oldValidator); err != nil {
		k.Logger(ctx).Error("Unable to update signer", "error", err, "validatorId", validator.ID)
		return hmCommon.ErrSignerUpdateError(k.Codespace()).Result()
	}

	// adding new validator
	k.Logger(ctx).Debug("Adding new validator", "validator", validator.String())

	// save validator
	err = k.AddValidator(ctx, validator)
	if err != nil {
		k.Logger(ctx).Error("Unable to update signer", "error", err, "ValidatorID", validator.ID)
		return hmCommon.ErrSignerUpdateError(k.Codespace()).Result()
	}
	// save staking sequence
	k.SetStakingSequence(ctx, sequence)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSignerUpdate,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(validator.ID.Uint64(), 10)),
			sdk.NewAttribute(types.AttributeKeyUpdatedAt, strconv.FormatUint(validator.LastUpdated, 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// HandleMsgValidatorExit handle msg validator exit
func HandleMsgValidatorExit(ctx sdk.Context, msg types.MsgValidatorExit, k Keeper, contractCaller helper.IContractCaller) sdk.Result {
	k.Logger(ctx).Info("Handling validator exit", "ValidatorID", msg.ID)

	if confirmed := contractCaller.IsTxConfirmed(msg.TxHash.EthHash()); !confirmed {
		return hmCommon.ErrWaitForConfirmation(k.Codespace()).Result()
	}
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

	// get validator from mainchain
	updatedVal, err := contractCaller.GetValidatorInfo(validator.ID)
	if err != nil {
		k.Logger(ctx).Error("Cannot fetch validator info while unstaking", "Error", err, "validatorID", validator.ID)
		return hmCommon.ErrNoValidator(k.Codespace()).Result()
	}

	// Add deactivation time for validator
	if err := k.AddDeactivationEpoch(ctx, validator, updatedVal); err != nil {
		k.Logger(ctx).Error("Error while setting deactivation epoch to validator", "error", err, "validatorID", validator.ID)
		return hmCommon.ErrValidatorNotDeactivated(k.Codespace()).Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeValidatorExit,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(validator.ID.Uint64(), 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
