package staking

import (
	"bytes"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking/tags"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
	"github.com/maticnetwork/heimdall/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// NewHandler new handler
func NewHandler(k Keeper, contractCaller helper.IContractCaller) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgValidatorJoin:
			return HandleMsgValidatorJoin(ctx, msg, k, contractCaller)
		case MsgValidatorExit:
			return HandleMsgValidatorExit(ctx, msg, k, contractCaller)
		case MsgSignerUpdate:
			return HandleMsgSignerUpdate(ctx, msg, k, contractCaller)
		case MsgStakeUpdate:
			return HandleMsgStakeUpdate(ctx, msg, k, contractCaller)
		default:
			return sdk.ErrTxDecode("Invalid message in checkpoint module").Result()
		}
	}
}

// handleMsgValidatorJoin msg validator join
func HandleMsgValidatorJoin(ctx sdk.Context, msg MsgValidatorJoin, k Keeper, contractCaller helper.IContractCaller) sdk.Result {
	k.Logger(ctx).Debug("Handing new validator join", "msg", msg)

	if confirmed := contractCaller.IsTxConfirmed(msg.TxHash.EthHash()); !confirmed {
		return hmCommon.ErrWaitForConfirmation(k.Codespace()).Result()
	}

	// fetch validator from mainchain
	validator, err := contractCaller.GetValidatorInfo(msg.ID)
	if err != nil {
		k.Logger(ctx).Error(
			"Unable to fetch validator from rootchain",
			"error", err,
		)
		return hmCommon.ErrNoValidator(k.Codespace()).Result()
	}

	if bytes.Equal(validator.Signer.Bytes(), helper.ZeroAddress.Bytes()) {
		k.Logger(ctx).Error(
			"No validator signer found",
			"msgValidator", msg.ID,
		)
		return hmCommon.ErrNoValidator(k.Codespace()).Result()
	}

	k.Logger(ctx).Debug("Fetched validator from rootchain successfully", "validator", validator.String())

	// Generate PubKey from Pubkey in message and signer
	pubkey := msg.SignerPubKey
	signer := pubkey.Address()

	// check signer in message corresponds
	if !bytes.Equal(signer.Bytes(), validator.Signer.Bytes()) {
		k.Logger(ctx).Error(
			"Signer Address does not match",
			"msgValidator", signer.String(),
			"mainchainValidator", validator.Signer.String(),
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

	// create new validator
	newValidator := hmTypes.Validator{
		ID:          validator.ID,
		StartEpoch:  validator.StartEpoch,
		EndEpoch:    validator.EndEpoch,
		VotingPower: validator.VotingPower,
		PubKey:      pubkey,
		Signer:      validator.Signer,
		LastUpdated: 0,
	}

	// add validator to store
	k.Logger(ctx).Debug("Adding new validator to state", "validator", newValidator.String())
	err = k.AddValidator(ctx, newValidator)
	if err != nil {
		k.Logger(ctx).Error("Unable to add validator to state", "error", err, "validator", newValidator.String())
		return hmCommon.ErrValidatorSave(k.Codespace()).Result()
	}

	resTags := sdk.NewTags(
		tags.ValidatorJoin, []byte(newValidator.Signer.String()),
		tags.ValidatorID, []byte(strconv.FormatUint(newValidator.ID.Uint64(), 10)),
	)

	return sdk.Result{Tags: resTags}
}

// handleMsgStakeUpdate handles stake update message
func HandleMsgStakeUpdate(ctx sdk.Context, msg MsgStakeUpdate, k Keeper, contractCaller helper.IContractCaller) sdk.Result {
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

	// last updated
	lastUpdated := (receipt.BlockNumber.Uint64() * stakingTypes.DefaultLogIndexUnit) + msg.LogIndex

	// check if incoming tx is older
	if lastUpdated <= validator.LastUpdated {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	// update last udpated
	validator.LastUpdated = lastUpdated

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

	resTags := sdk.NewTags(
		tags.UpdatedAt, []byte(strconv.FormatUint(validator.LastUpdated, 10)),
		tags.ValidatorID, []byte(strconv.FormatUint(validator.ID.Uint64(), 10)),
	)

	return sdk.Result{Tags: resTags}
}

// handleMsgSignerUpdate handles signer update message
func HandleMsgSignerUpdate(ctx sdk.Context, msg MsgSignerUpdate, k Keeper, contractCaller helper.IContractCaller) sdk.Result {
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

	// last updated
	lastUpdated := (receipt.BlockNumber.Uint64() * stakingTypes.DefaultLogIndexUnit) + msg.LogIndex

	// check if incoming tx is older
	if lastUpdated <= validator.LastUpdated {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	// update last udpated
	validator.LastUpdated = lastUpdated

	// check if we are actually updating signer
	if !bytes.Equal(newSigner.Bytes(), validator.Signer.Bytes()) {
		// Update signer in prev Validator
		validator.Signer = types.HeimdallAddress(newSigner)
		validator.PubKey = newPubKey
		k.Logger(ctx).Debug("Updating new signer", "signer", newSigner.String(), "oldSigner", oldValidator.Signer.String(), "validatorID", msg.ID)
	}

	k.Logger(ctx).Debug("Removing old validator", "validator", oldValidator.String())

	// remove old validator from HM
	oldValidator.EndEpoch = k.ackRetriever.GetACKCount(ctx)

	// remove old validator from TM
	oldValidator.VotingPower = 0
	// updated last
	oldValidator.LastUpdated = lastUpdated

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

	resTags := sdk.NewTags(
		tags.SignerUpdate, []byte(newSigner.String()),
		tags.UpdatedAt, []byte(strconv.FormatUint(validator.LastUpdated, 10)),
		tags.ValidatorID, []byte(strconv.FormatUint(validator.ID.Uint64(), 10)),
	)

	return sdk.Result{Tags: resTags}
}

// handleMsgValidatorExit handle msg validator exit
func HandleMsgValidatorExit(ctx sdk.Context, msg MsgValidatorExit, k Keeper, contractCaller helper.IContractCaller) sdk.Result {
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

	resTags := sdk.NewTags(
		tags.ValidatorExit, []byte(validator.Signer.String()),
		tags.ValidatorID, []byte(strconv.FormatUint(uint64(validator.ID), 10)),
	)

	return sdk.Result{Tags: resTags}
}
