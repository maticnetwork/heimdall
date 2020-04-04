package staking

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
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
	receipt, err := contractCaller.GetConfirmedTxReceipt(ctx.BlockTime(), msg.TxHash.EthHash())
	if err != nil || receipt == nil {
		return hmCommon.ErrWaitForConfirmation(k.Codespace()).Result()
	}

	// decode validator join event
	eventLog, err := contractCaller.DecodeValidatorJoinEvent(helper.GetStakingInfoAddress(), receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Unable to fetch logs for txHash").Result()
	}

	// Generate PubKey from Pubkey in message and signer
	pubkey := msg.SignerPubKey
	signer := pubkey.Address()

	// check signer pubkey in message corresponds
	if !bytes.Equal(pubkey.Bytes()[1:], eventLog.SignerPubkey) {
		k.Logger(ctx).Error(
			"Signer Pubkey does not match",
			"msgValidator", pubkey.String(),
			"mainchainValidator", hmTypes.BytesToHexBytes(eventLog.SignerPubkey),
		)
		return hmCommon.ErrValSignerPubKeyMismatch(k.Codespace()).Result()
	}

	// check signer corresponding to pubkey matches signer from event
	if !bytes.Equal(signer.Bytes(), eventLog.Signer.Bytes()) {
		k.Logger(ctx).Error(
			"Signer Address from Pubkey does not match",
			"Validator", signer.String(),
			"mainchainValidator", eventLog.Signer.Hex(),
		)
		return hmCommon.ErrValSignerMismatch(k.Codespace()).Result()
	}

	// check msg id
	if eventLog.ValidatorId.Uint64() != msg.ID.Uint64() {
		k.Logger(ctx).Error("ID in message doesn't match with id in log", "msgId", msg.ID, "validatorIdFromTx", eventLog.ValidatorId)
		return hmCommon.ErrInvalidMsg(k.Codespace(), "ID in message doesn't match with id in log. msgId %v validatorIdFromTx %v", msg.ID, eventLog.ValidatorId).Result()
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
		return hmCommon.ErrInvalidMsg(k.Codespace(), fmt.Sprintf("Invalid amount %v for validator %v", eventLog.Amount, msg.ID)).Result()
	}

	// create new validator
	newValidator := hmTypes.Validator{
		ID:          msg.ID,
		StartEpoch:  eventLog.ActivationEpoch.Uint64(),
		EndEpoch:    0,
		VotingPower: votingPower.Int64(),
		PubKey:      pubkey,
		Signer:      hmTypes.BytesToHeimdallAddress(signer.Bytes()),
		LastUpdated: "",
	}

	// sequence id

	sequence := new(big.Int).Mul(receipt.BlockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// check if incoming tx is older
	if k.HasStakingSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	// update last updated
	newValidator.LastUpdated = sequence.String()

	// add validator to store
	k.Logger(ctx).Debug("Adding new validator to state", "validator", newValidator.String())
	err = k.AddValidator(ctx, newValidator)
	if err != nil {
		k.Logger(ctx).Error("Unable to add validator to state", "error", err, "validator", newValidator.String())
		return hmCommon.ErrValidatorSave(k.Codespace()).Result()
	}

	// save staking sequence
	k.SetStakingSequence(ctx, sequence.String())

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
	receipt, err := contractCaller.GetConfirmedTxReceipt(ctx.BlockTime(), msg.TxHash.EthHash())
	if err != nil || receipt == nil {
		return hmCommon.ErrWaitForConfirmation(k.Codespace()).Result()
	}

	eventLog, err := contractCaller.DecodeValidatorStakeUpdateEvent(helper.GetStakingInfoAddress(), receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Unable to fetch logs for txHash").Result()
	}

	if eventLog.ValidatorId.Uint64() != msg.ID.Uint64() {
		k.Logger(ctx).Error("ID in message doesn't match with id in log", "msgId", msg.ID, "validatorIdFromTx", eventLog.ValidatorId)
		return hmCommon.ErrInvalidMsg(k.Codespace(), "ID in message doesn't match with id in log. msgId %v validatorIdFromTx %v", msg.ID, eventLog.ValidatorId).Result()
	}

	// pull validator from store
	validator, ok := k.GetValidatorFromValID(ctx, msg.ID)
	if !ok {
		k.Logger(ctx).Error("Fetching of validator from store failed", "validatorId", msg.ID)
		return hmCommon.ErrNoValidator(k.Codespace()).Result()
	}

	// sequence id

	sequence := new(big.Int).Mul(receipt.BlockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// check if incoming tx is older
	if k.HasStakingSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	// update last updated
	validator.LastUpdated = sequence.String()

	// set validator amount
	p, err := helper.GetPowerFromAmount(eventLog.NewAmount)
	if err != nil {
		return hmCommon.ErrInvalidMsg(k.Codespace(), fmt.Sprintf("Invalid amount %v for validator %v", eventLog.NewAmount, msg.ID)).Result()
	}
	validator.VotingPower = p.Int64()

	// save validator
	err = k.AddValidator(ctx, validator)
	if err != nil {
		k.Logger(ctx).Error("Unable to update signer", "error", err, "ValidatorID", validator.ID)
		return hmCommon.ErrSignerUpdateError(k.Codespace()).Result()
	}

	// save staking sequence

	k.SetStakingSequence(ctx, sequence.String())

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeStakeUpdate,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(validator.ID.Uint64(), 10)),
			sdk.NewAttribute(types.AttributeKeyUpdatedAt, validator.LastUpdated),
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
	receipt, err := contractCaller.GetConfirmedTxReceipt(ctx.BlockTime(), msg.TxHash.EthHash())
	if err != nil || receipt == nil {
		return hmCommon.ErrWaitForConfirmation(k.Codespace()).Result()
	}

	newPubKey := msg.NewSignerPubKey
	newSigner := newPubKey.Address()

	eventLog, err := contractCaller.DecodeSignerUpdateEvent(helper.GetStakingInfoAddress(), receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Unable to fetch signer update log for txHash").Result()
	}

	if eventLog.ValidatorId.Uint64() != msg.ID.Uint64() {
		k.Logger(ctx).Error("ID in message doesn't match with id in log", "msgId", msg.ID, "validatorIdFromTx", eventLog.ValidatorId)
		return hmCommon.ErrInvalidMsg(k.Codespace(), "ID in message doesn't match with id in log. msgId %v validatorIdFromTx %v", msg.ID, eventLog.ValidatorId).Result()
	}

	if bytes.Compare(eventLog.SignerPubkey, newPubKey.Bytes()[1:]) != 0 {
		k.Logger(ctx).Error("Newsigner pubkey in txhash and msg dont match", "msgPubKey", newPubKey.String(), "pubkeyTx", hmTypes.NewPubKey(eventLog.SignerPubkey[:]).String())
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Newsigner pubkey in txhash and msg dont match").Result()
	}

	// check signer corresponding to pubkey matches signer from event
	if !bytes.Equal(newSigner.Bytes(), eventLog.NewSigner.Bytes()) {
		k.Logger(ctx).Error("Signer Address from Pubkey does not match", "Validator", newSigner.String(), "mainchainValidator", eventLog.NewSigner.Hex())
		return hmCommon.ErrValSignerMismatch(k.Codespace()).Result()
	}

	// pull validator from store
	validator, ok := k.GetValidatorFromValID(ctx, msg.ID)
	if !ok {
		k.Logger(ctx).Error("Fetching of validator from store failed", "validatorId", msg.ID)
		return hmCommon.ErrNoValidator(k.Codespace()).Result()
	}
	oldValidator := validator.Copy()

	// sequence id

	sequence := new(big.Int).Mul(receipt.BlockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// check if incoming tx is older
	if k.HasStakingSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	// update last udpated
	validator.LastUpdated = sequence.String()

	// check if we are actually updating signer
	if !bytes.Equal(newSigner.Bytes(), validator.Signer.Bytes()) {
		// Update signer in prev Validator
		validator.Signer = hmTypes.HeimdallAddress(newSigner)
		validator.PubKey = newPubKey
		k.Logger(ctx).Debug("Updating new signer", "signer", newSigner.String(), "oldSigner", oldValidator.Signer.String(), "validatorID", msg.ID)
	}

	k.Logger(ctx).Debug("Removing old validator", "validator", oldValidator.String())

	// remove old validator from HM
	oldValidator.EndEpoch = k.moduleCommunicator.GetACKCount(ctx)

	// remove old validator from TM
	oldValidator.VotingPower = 0
	// updated last
	oldValidator.LastUpdated = sequence.String()

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
	k.SetStakingSequence(ctx, sequence.String())

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeSignerUpdate,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyValidatorID, validator.ID.String()),
			sdk.NewAttribute(types.AttributeKeyUpdatedAt, validator.LastUpdated),
		),
	})

	//
	// Move heimdall fee to new signer
	//

	// check if fee is already withdrawn
	coins := k.moduleCommunicator.GetCoins(ctx, oldValidator.Signer)
	maticBalance := coins.AmountOf(authTypes.FeeToken)
	if !maticBalance.IsZero() {
		k.Logger(ctx).Info("Transferring fee", "from", oldValidator.Signer.String(), "to", validator.Signer.String(), "balance", maticBalance.String())
		maticCoins := hmTypes.Coins{hmTypes.Coin{Denom: authTypes.FeeToken, Amount: maticBalance}}
		if err := k.moduleCommunicator.SendCoins(ctx, oldValidator.Signer, validator.Signer, maticCoins); err != nil {
			k.Logger(ctx).Info("Error while transferring fee", "from", oldValidator.Signer.String(), "to", validator.Signer.String(), "balance", maticBalance.String())
			return err.Result()
		}
	}

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// HandleMsgValidatorExit handle msg validator exit
func HandleMsgValidatorExit(ctx sdk.Context, msg types.MsgValidatorExit, k Keeper, contractCaller helper.IContractCaller) sdk.Result {
	k.Logger(ctx).Info("Handling validator exit", "ValidatorID", msg.ID)

	// get main tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(ctx.BlockTime(), msg.TxHash.EthHash())
	if err != nil || receipt == nil {
		return hmCommon.ErrWaitForConfirmation(k.Codespace()).Result()
	}

	// decode validator exit
	eventLog, err := contractCaller.DecodeValidatorExitEvent(helper.GetStakingInfoAddress(), receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Unable to fetch unstake log for txHash").Result()
	}

	if eventLog.ValidatorId.Uint64() != msg.ID.Uint64() {
		k.Logger(ctx).Error("ID in message doesn't match with id in log", "msgId", msg.ID, "validatorIdFromTx", eventLog.ValidatorId)
		return hmCommon.ErrInvalidMsg(k.Codespace(), "ID in message doesn't match with id in log. msgId %v validatorIdFromTx %v", msg.ID, eventLog.ValidatorId).Result()
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

	// set end epoch
	validator.EndEpoch = eventLog.DeactivationEpoch.Uint64()

	// sequence id
	sequence := new(big.Int).Mul(receipt.BlockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// check if incoming tx is older
	if k.HasStakingSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	// update last updated
	validator.LastUpdated = sequence.String()

	// Add deactivation time for validator
	if err := k.AddValidator(ctx, validator); err != nil {
		k.Logger(ctx).Error("Error while setting deactivation epoch to validator", "error", err, "validatorID", validator.ID.String())
		return hmCommon.ErrValidatorNotDeactivated(k.Codespace()).Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeValidatorExit,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyValidatorID, validator.ID.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
