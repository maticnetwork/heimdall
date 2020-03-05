package topup

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/auth"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/topup/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// NewHandler returns a handler for "topup" type messages.
func NewHandler(k Keeper, contractCaller helper.IContractCaller) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgTopup:
			return handleMsgTopup(ctx, k, msg, contractCaller)
		case types.MsgWithdrawFee:
			return handleMsgWithdrawFee(ctx, k, msg)
		default:
			errMsg := "Unrecognized topup Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// Handle MsgMintFeeToken
func handleMsgTopup(ctx sdk.Context, k Keeper, msg types.MsgTopup, contractCaller helper.IContractCaller) sdk.Result {
	if !k.bk.GetSendEnabled(ctx) {
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

	// create topup amount
	topupAmount := hmTypes.Coins{hmTypes.Coin{Denom: authTypes.FeeToken, Amount: hmTypes.NewIntFromBigInt(eventLog.Fee)}}

	// sequence id
	sequence := (receipt.BlockNumber.Uint64() * hmTypes.DefaultLogIndexUnit) + msg.LogIndex

	// check if incoming tx already exists
	if k.HasTopupSequence(ctx, sequence) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	// increase coins in account
	if _, err := k.bk.AddCoins(ctx, signer, topupAmount); err != nil {
		return err.Result()
	}

	// transfer fees to sender (proposer)
	if err := k.bk.SendCoins(ctx, signer, msg.FromAddress, auth.DefaultFeeWantedPerTx); err != nil {
		return err.Result()
	}

	// save topup
	k.SetTopupSequence(ctx, sequence)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTopup,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyValidatorID, msg.ID.String()),
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
	coins := k.bk.GetCoins(ctx, msg.ValidatorAddress)
	maticBalance := coins.AmountOf(authTypes.FeeToken)

	validator, err := k.sk.GetValidatorInfo(ctx, msg.ValidatorAddress.Bytes())
	if err != nil {
		k.Logger(ctx).Info("Fee balance for ", "fromAddress", msg.ValidatorAddress, "validatorId", validator.ID, "balance", maticBalance.BigInt().String())
		if maticBalance.IsZero() {
			return types.ErrNoBalanceToWithdraw(k.Codespace()).Result()
		}

		// withdraw coins of validator
		maticCoins := hmTypes.Coins{hmTypes.Coin{Denom: authTypes.FeeToken, Amount: maticBalance}}
		if _, err := k.bk.SubtractCoins(ctx, msg.ValidatorAddress, maticCoins); err != nil {
			k.Logger(ctx).Error("Error while setting Fee balance to zero ", "fromAddress", msg.ValidatorAddress, "validatorId", validator.ID, "err", err)
			return err.Result()
		}

		// Add Fee to Dividend Account
		feeAmount := maticBalance.BigInt()
		k.vm.AddFeeToDividendAccount(ctx, validator.ID, feeAmount)

		ctx.EventManager().EmitEvents(sdk.Events{
			sdk.NewEvent(
				types.EventTypeFeeWithdraw,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(types.AttributeKeyValidatorID, strconv.FormatUint(uint64(validator.ID), 10)),
				// sdk.NewAttribute(types.AttributeKeyValidatorSigner, msg.ValidatorAddress.String()),
				sdk.NewAttribute(types.AttributeKeyFeeWithdrawAmount, feeAmount.String()),
			),
		})

		return sdk.Result{
			Events: ctx.EventManager().Events(),
		}
	}
	return sdk.Result{}
}
