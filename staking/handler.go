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
		case types.MsgDelegatorBond:
			return HandleMsgDelegatorBond(ctx, msg, k, contractCaller)
		case types.MsgDelegatorUnBond:
			return HandleMsgDelegatorUnBond(ctx, msg, k, contractCaller)
		case types.MsgDelegatorReBond:
			return HandleMsgDelegatorReBond(ctx, msg, k, contractCaller)
		case types.MsgDelStakeUpdate:
			return HandleMsgDelStakeUpdate(ctx, msg, k, contractCaller)
		default:
			return sdk.ErrTxDecode("Invalid message in checkpoint module").Result()
		}
	}
}

// HandleMsgValidatorJoin msg validator join
func HandleMsgValidatorJoin(ctx sdk.Context, msg types.MsgValidatorJoin, k Keeper, contractCaller helper.IContractCaller) sdk.Result {
	k.Logger(ctx).Info("Handing new validator join", "msg", msg)

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

// HandleMsgCommissionRateUpdate handles commission rate update message
func HandleMsgCommissionRateUpdate(ctx sdk.Context, msg types.MsgStakeUpdate, k Keeper, contractCaller helper.IContractCaller) sdk.Result {
	k.Logger(ctx).Debug("Handling commission rate update", "Validator", msg.ID)

	// get main tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(msg.TxHash.EthHash())
	if err != nil || receipt == nil {
		return hmCommon.ErrWaitForConfirmation(k.Codespace()).Result()
	}

	eventLog, err := contractCaller.DecodeCommissionRateUpdateEvent(receipt, msg.LogIndex)
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
	validator.CommissionRate = eventLog.Rate.Uint64()

	// save validator
	err = k.AddValidator(ctx, validator)
	if err != nil {
		k.Logger(ctx).Error("Unable to update commission rate", "error", err, "ValidatorID", validator.ID)
		return hmCommon.ErrSignerUpdateError(k.Codespace()).Result()
	}

	// save staking sequence
	k.SetStakingSequence(ctx, sequence)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCommissionUpdate,
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

// HandleMsgDelegatorBond msg delegator Bond with Validator
// 1. On Bonding event, Validator to whom delegator is bonded will send `MsgDelegatorBond` transaction to Heimdall.
// 2. Delegator is updated with Validator ID.
// 3. VotingPower of the bonded validator is updated.
// 4. shares are added to Delegator proportional to his stake and exchange rate. // delegatorshares = (delegatorstake / exchangeRate)
// 5. Exchange rate is calculated instantly.  //   ExchangeRate = (delegatedpower + delegatorRewardPool) / totaldelegatorshares
// 6. TotalDelegatorShares of bonded validator is updated.
// 7. DelegatedPower of bonded validator is updated.
func HandleMsgDelegatorBond(ctx sdk.Context, msg types.MsgDelegatorBond, k Keeper, contractCaller helper.IContractCaller) sdk.Result {

	k.Logger(ctx).Info("Handling delegator bond with validator", "Delegator", msg.ID)

	// get main tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(msg.TxHash.EthHash())
	if err != nil || receipt == nil {
		return hmCommon.ErrWaitForConfirmation(k.Codespace()).Result()
	}

	eventLog, err := contractCaller.DecodeDelegatorBondEvent(receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Unable to fetch delegator bond log for txHash").Result()
	}

	if eventLog.DelegatorId.Uint64() != msg.ID.Uint64() {
		k.Logger(ctx).Error("Delegator ID in message doesnt match delegator id in logs", "MsgID", msg.ID, "IdFromTx", eventLog.DelegatorId.Uint64())
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Invalid txhash, id's dont match. Id from tx hash is %v", eventLog.DelegatorId.Uint64()).Result()
	}

	// sequence id
	sequence := (receipt.BlockNumber.Uint64() * hmTypes.DefaultLogIndexUnit) + msg.LogIndex

	// check if incoming tx is older
	if k.HasStakingSequence(ctx, sequence) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	err = k.BondDelegator(ctx, msg.ID, hmTypes.ValidatorID(eventLog.ValidatorId.Uint64()), eventLog.Amount)
	if err != nil {
		return hmCommon.ErrDelegatorBond(k.Codespace()).Result()
	}

	// save staking sequence
	k.SetStakingSequence(ctx, sequence)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDelegatorBond,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(eventLog.ValidatorId.Uint64(), 10)),
			sdk.NewAttribute(types.AttributeKeyDelegatorID, strconv.FormatUint(eventLog.DelegatorId.Uint64(), 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// HandleMsgDelegatorUnBond msg delegator unbond with validator
// ** stake calculations **
// 1. On Bonding event, Validator will send MsgDelegatorUnBond transaction to heimdall.
// 2. Delegator is updated with Validator ID = 0.
// 3. VotingPower of bonded validator is reduced.
// 4. DelegatedPower of the bonded validator is reduced after reward calculation.

// ** reward calculations **
// 1. Exchange rate is calculated instantly.  ExchangeRate = (delegatedpower + delegatorRewardPool) / totaldelegatorshares
// 2. Based on exchange rate and no of shares delegator holds, totalReturns for delegator is calculated.  `totalReturns = exchangeRate * noOfShares`
// 3. Delegator RewardAmount += totalReturns - delegatorVotingPower
// 4. Add RewardAmount to DelegatorAccount .
// 5. Reduce TotalDelegatorShares of bonded validator.
// 6. Reduce DelgatorRewardPool of bonded validator.
// 7. make shares = 0 on Delegator Account.
func HandleMsgDelegatorUnBond(ctx sdk.Context, msg types.MsgDelegatorUnBond, k Keeper, contractCaller helper.IContractCaller) sdk.Result {
	k.Logger(ctx).Debug("Handling delegator unbond", "msg", msg)

	// get main tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(msg.TxHash.EthHash())
	if err != nil || receipt == nil {
		return hmCommon.ErrWaitForConfirmation(k.Codespace()).Result()
	}

	eventLog, err := contractCaller.DecodeDelegatorUnBondEvent(receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Unable to fetch delegator unbond log for txHash").Result()
	}

	if eventLog.DelegatorId.Uint64() != msg.ID.Uint64() {
		k.Logger(ctx).Error("Delegator ID in message doesnt match delegator id in logs", "MsgID", msg.ID, "IdFromTx", eventLog.DelegatorId.Uint64())
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Invalid txhash, id's dont match. Id from tx hash is %v", eventLog.DelegatorId.Uint64()).Result()
	}

	// sequence id
	sequence := (receipt.BlockNumber.Uint64() * hmTypes.DefaultLogIndexUnit) + msg.LogIndex

	// check if incoming tx is older
	if k.HasStakingSequence(ctx, sequence) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	err = k.UnBondDelegator(ctx, msg.ID, hmTypes.ValidatorID(eventLog.ValidatorId.Uint64()), eventLog.Amount)
	if err != nil {
		return hmCommon.ErrDelegatorUnBond(k.Codespace()).Result()
	}
	// save staking sequence
	k.SetStakingSequence(ctx, sequence)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDelegatorBond,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(eventLog.ValidatorId.Uint64(), 10)),
			sdk.NewAttribute(types.AttributeKeyDelegatorID, strconv.FormatUint(eventLog.DelegatorId.Uint64(), 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// HandleMsgDelegatorReBond msg delegator rebond with validator
// 1. Unbond from old validator
// 2. Bond with new validator
func HandleMsgDelegatorReBond(ctx sdk.Context, msg types.MsgDelegatorReBond, k Keeper, contractCaller helper.IContractCaller) sdk.Result {
	k.Logger(ctx).Debug("Handling delegator rebond", "msg", msg)

	// get main tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(msg.TxHash.EthHash())
	if err != nil || receipt == nil {
		return hmCommon.ErrWaitForConfirmation(k.Codespace()).Result()
	}

	eventLog, err := contractCaller.DecodeDelegatorReBondEvent(receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Unable to fetch delegator rebond log for txHash").Result()
	}

	if eventLog.DelegatorId.Uint64() != msg.ID.Uint64() {
		k.Logger(ctx).Error("Delegator ID in message doesnt match delegator id in logs", "MsgID", msg.ID, "IdFromTx", eventLog.DelegatorId.Uint64())
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Invalid txhash, id's dont match. Id from tx hash is %v", eventLog.DelegatorId.Uint64()).Result()
	}

	// sequence id
	sequence := (receipt.BlockNumber.Uint64() * hmTypes.DefaultLogIndexUnit) + msg.LogIndex

	// check if incoming tx is older
	if k.HasStakingSequence(ctx, sequence) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	err = k.ReBondDelegator(ctx, msg.ID, eventLog.Amount, hmTypes.ValidatorID(eventLog.OldValidatorId.Uint64()), hmTypes.ValidatorID(eventLog.NewValidatorId.Uint64()))
	if err != nil {
		return hmCommon.ErrDelegatorUnBond(k.Codespace()).Result()
	}

	// save staking sequence
	k.SetStakingSequence(ctx, sequence)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDelegatorBond,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(eventLog.OldValidatorId.Uint64(), 10)),
			sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(eventLog.NewValidatorId.Uint64(), 10)),
			sdk.NewAttribute(types.AttributeKeyDelegatorID, strconv.FormatUint(eventLog.DelegatorId.Uint64(), 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// HandleMsgDelStakeUpdate msg delegator stake Update
// 1. if old amount is greater than new amount, It's becoz of slashing. Burn shares. Update slashed amount. No change of rewards
// 2. if old amount is lesser than new amount, It's becoz of new stake added. Add shares.
func HandleMsgDelStakeUpdate(ctx sdk.Context, msg types.MsgDelStakeUpdate, k Keeper, contractCaller helper.IContractCaller) sdk.Result {
	k.Logger(ctx).Debug("Handling delegator stake update", "msg", msg)

	// get main tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(msg.TxHash.EthHash())
	if err != nil || receipt == nil {
		return hmCommon.ErrWaitForConfirmation(k.Codespace()).Result()
	}

	eventLog, err := contractCaller.DecodeDelStakeUpdateEvent(receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Unable to fetch delegator stake update log for txHash").Result()
	}

	if eventLog.DelegatorId.Uint64() != msg.ID.Uint64() {
		k.Logger(ctx).Error("Delegator ID in message doesnt match delegator id in logs", "MsgID", msg.ID, "IdFromTx", eventLog.DelegatorId.Uint64())
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Invalid txhash, id's dont match. Id from tx hash is %v", eventLog.DelegatorId.Uint64()).Result()
	}

	// sequence id
	sequence := (receipt.BlockNumber.Uint64() * hmTypes.DefaultLogIndexUnit) + msg.LogIndex

	// check if incoming tx is older
	if k.HasStakingSequence(ctx, sequence) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	err = k.DelStakeUpdate(ctx, msg.ID, hmTypes.ValidatorID(eventLog.ValidatorId.Uint64()), eventLog.OldAmount, eventLog.NewAmount)
	if err != nil {
		return hmCommon.ErrDelegatorStakeUpdate(k.Codespace()).Result()
	}

	// save staking sequence
	k.SetStakingSequence(ctx, sequence)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeDelegatorBond,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyDelegatorID, strconv.FormatUint(eventLog.DelegatorId.Uint64(), 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
