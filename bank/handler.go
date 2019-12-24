package bank

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/auth"
	"github.com/maticnetwork/heimdall/bank/types"
	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// NewHandler returns a handler for "bank" type messages.
func NewHandler(k Keeper, contractCaller helper.IContractCaller) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case types.MsgSend:
			return handleMsgSend(ctx, k, msg)
		case types.MsgMultiSend:
			return handleMsgMultiSend(ctx, k, msg)
		case types.MsgTopup:
			return handleMsgTopup(ctx, k, msg, contractCaller)
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
	tags, err := k.SendCoins(ctx, msg.FromAddress, msg.ToAddress, msg.Amount)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
	}
}

// Handle MsgMultiSend.
func handleMsgMultiSend(ctx sdk.Context, k Keeper, msg types.MsgMultiSend) sdk.Result {
	// NOTE: totalIn == totalOut should already have been checked
	if !k.GetSendEnabled(ctx) {
		return types.ErrSendDisabled(k.Codespace()).Result()
	}
	tags, err := k.InputOutputCoins(ctx, msg.Inputs, msg.Outputs)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{
		Tags: tags,
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
	eventLog, err := contractCaller.DecodeValidatorTopupFeesEvent(receipt, msg.LogIndex)
	if err != nil || eventLog == nil {
		k.Logger(ctx).Error("Error fetching log from txhash")
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Unable to fetch logs for txHash").Result()
	}

	if eventLog.ValidatorId.Uint64() != msg.ID.Uint64() {
		k.Logger(ctx).Error("ID in message doesn't match id in logs", "MsgID", msg.ID, "IdFromTx", eventLog.ValidatorId)
		return hmCommon.ErrInvalidMsg(k.Codespace(), "Invalid txhash and id don't match. Id from tx hash is %v", eventLog.ValidatorId.Uint64()).Result()
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

	// validator topup
	topupObject, err := k.GetTopup(ctx, validator.Signer)
	if err != nil {
		return types.ErrNoValidatorTopup(k.Codespace()).Result()
	}

	// create topup object
	if topupObject == nil {
		topupObject = &types.ValidatorTopup{
			ID:          validator.ID,
			TotalTopups: hmTypes.Coins{hmTypes.Coin{Denom: "vetic", Amount: hmTypes.NewInt(1)}},
			LastUpdated: 0,
		}
	}

	// create topup amount
	topupAmount := hmTypes.Coins{hmTypes.Coin{Denom: "vetic", Amount: hmTypes.NewIntFromBigInt(eventLog.Amount)}}

	// last updated
	lastUpdated := (receipt.BlockNumber.Uint64() * hmTypes.DefaultLogIndexUnit) + msg.LogIndex

	// check if incoming tx is older
	if lastUpdated <= topupObject.LastUpdated {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	// update last udpated
	topupObject.LastUpdated = lastUpdated
	// add total topups amount
	topupObject.TotalTopups = topupObject.TotalTopups.Add(topupAmount)
	// increase coins in account
	_, addTags, ec := k.AddCoins(ctx, validator.Signer, topupAmount)
	if ec != nil {
		return ec.Result()
	}
	// transfer fees to sender (proposer)
	transferTags, ec := k.SendCoins(ctx, validator.Signer, msg.FromAddress, auth.FeeWantedPerTx)
	if ec != nil {
		return ec.Result()
	}
	// save old validator
	if err := k.SetTopup(ctx, validator.Signer, *topupObject); err != nil {
		k.Logger(ctx).Error("Unable to update signer", "error", err, "validatorId", validator.ID)
		return hmCommon.ErrSignerUpdateError(k.Codespace()).Result()
	}

	// response tags
	resTags := sdk.NewTags(
		TagValidatorID, []byte(strconv.FormatUint(uint64(msg.ID), 10)),
		TagTopupAmount, []byte(strconv.FormatUint(eventLog.Amount.Uint64(), 10)),
	)

	// new tags
	return sdk.Result{
		Tags: addTags.AppendTags(resTags).AppendTags(transferTags),
	}
}
