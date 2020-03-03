package bank

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/auth"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/bank/types"
	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k Keeper, contractCaller helper.IContractCaller) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgSend:
			return handleMsgSend(ctx, k, msg)
		case types.MsgMultiSend:
			return handleMsgMultiSend(ctx, k, msg)
		case types.MsgTopup:
			return handleMsgTopup(ctx, k, msg, contractCaller)
		case types.MsgWithdrawFee:
			return handleMsgWithdrawFee(ctx, k, msg)
		default:
			errMsg := "Unrecognized bank Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgSend.
func handleMsgSend(ctx sdk.Context, k Keeper, msg types.MsgSend) sdk.Result {
	if !k.GetSendEnabled(ctx) {
		return types.ErrSendDisabled(k.Codespace()).Result()
	}

	err := k.SendCoins(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// Handle MsgMultiSend.
func handleMsgMultiSend(ctx sdk.Context, k Keeper, msg types.MsgMultiSend) sdk.Result {
	// NOTE: totalIn == totalOut should already have been checked
	if !k.GetSendEnabled(ctx) {
		return types.ErrSendDisabled(k.Codespace()).Result()
	}

	err := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)
	if err != nil {
		return err.Result()
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// Handle MsgMintFeeToken
func handleMsgTopup(ctx sdk.Context, k Keeper, msg types.MsgTopup, contractCaller helper.IContractCaller) sdk.Result {
	if !k.GetSendEnabled(ctx) {
		return types.ErrSendDisabled(k.Codespace()).Result()
	}

	// get main tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(msg.TxHash.EthHash())
	if err != nil || receipt == nil {
		return hmCommon.ErrWaitForConfirmation(k.Codespace()).Result()
	}

	// get event log for topup
	eventLog, err := contractCaller.DecodeValidatorTopupFeesEvent(helper.GetStakingInfoAddress(), receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Unable to fetch logs for txHash").Result()
	}

	if eventLog.ValidatorId.Uint64() != msg.ID.Uint64() {
		k.Logger(ctx).Error("ID in message doesn't match id in logs", "MsgID", msg.ID, "IdFromTx", eventLog.ValidatorId)
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Invalid txhash and id don't match. Id from tx hash is %v", eventLog.ValidatorId.Uint64()).Result()
	}

	// use event log signer
	signer := hmTypes.BytesToHeimdallAddress(eventLog.Signer.Bytes())
	// if validator exists use siger from local state
	validator, found := k.vm.GetValidatorFromValID(ctx, msg.ID)
	if found {
		signer = validator.Signer
	}

	// validator topup
	topupObject, err := k.GetValidatorTopup(ctx, signer)
	if err != nil {
		return types.ErrNoValidatorTopup(k.Codespace()).Result()
	}

	// create topup object
	if topupObject == nil {
		topupObject = &types.ValidatorTopup{
			ID:          msg.ID,
			TotalTopups: hmTypes.Coins{hmTypes.Coin{Denom: authTypes.FeeToken, Amount: hmTypes.NewInt(0)}},
		}
	}

	// create topup amount
	topupAmount := hmTypes.Coins{hmTypes.Coin{Denom: authTypes.FeeToken, Amount: hmTypes.NewIntFromBigInt(eventLog.Fee)}}

	// sequence id
	sequence := (receipt.BlockNumber.Uint64() * hmTypes.DefaultLogIndexUnit) + msg.LogIndex

	// check if incoming tx already exists
	if k.HasTopupSequence(ctx, sequence) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	// add total topups amount
	topupObject.TotalTopups = topupObject.TotalTopups.Add(topupAmount)

	// increase coins in account
	if _, ec := k.AddCoins(ctx, signer, topupAmount); ec != nil {
		return ec.Result()
	}

	// transfer fees to sender (proposer)
	if ec := k.SendCoins(ctx, signer, msg.FromAddress, auth.DefaultFeeWantedPerTx); ec != nil {
		return ec.Result()
	}

	// save old validator
	if err := k.SetValidatorTopup(ctx, signer, *topupObject); err != nil {
		k.Logger(ctx).Error("Unable to update signer", "error", err, "validatorId", msg.ID.String())
		return hmCommon.ErrSignerUpdateError(k.Codespace()).Result()
	}
	// save topup
	k.SetTopupSequence(ctx, sequence)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTopup,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(uint64(msg.ID), 10)),
			sdk.NewAttribute(types.AttributeKeyValidatorSigner, signer.String()),
			sdk.NewAttribute(types.AttributeKeyTopupAmount, eventLog.Fee.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// Handle MsgWithdrawFee.
func handleMsgWithdrawFee(ctx sdk.Context, k Keeper, msg types.MsgWithdrawFee) sdk.Result {

	// check if fee is already withdrawn
	coins := k.GetCoins(ctx, msg.FromAddress)
	maticBalance := coins.AmountOf(authTypes.FeeToken)
	k.Logger(ctx).Info("Fee balance for ", "fromAddress", msg.FromAddress, "validatorId", msg.ID, "balance", maticBalance.BigInt().String())
	if maticBalance.IsZero() {
		return types.ErrNoBalanceToWithdraw(k.Codespace()).Result()
	}

	// withdraw coins of validator
	maticCoins := hmTypes.Coins{hmTypes.Coin{Denom: authTypes.FeeToken, Amount: maticBalance}}
	if _, err := k.SubtractCoins(ctx, msg.FromAddress, maticCoins); err != nil {
		k.Logger(ctx).Error("Error while setting Fee balance to zero ", "fromAddress", msg.FromAddress, "validatorId", msg.ID, "err", err)
		return err.Result()
	}

	// Add Fee to Dividend Account
	feeAmount := maticBalance.BigInt()
	k.vm.AddFeeToDividendAccount(ctx, msg.ID, feeAmount)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeFeeWithdraw,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(uint64(msg.ID), 10)),
			sdk.NewAttribute(types.AttributeKeyFeeWithdrawAmount, feeAmount.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
