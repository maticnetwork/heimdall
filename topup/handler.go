package topup

import (
	"math/big"

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

	// chainManager params
	params := k.chainKeeper.GetParams(ctx)
	chainParams := params.ChainParams

	// get main tx receipt
	receipt, err := contractCaller.GetConfirmedTxReceipt(ctx.BlockTime(), msg.TxHash.EthHash(), params.TxConfirmationTime)
	if err != nil || receipt == nil {
		return hmCommon.ErrWaitForConfirmation(k.Codespace(), params.TxConfirmationTime).Result()
	}

	// get event log for topup
	eventLog, err := contractCaller.DecodeValidatorTopupFeesEvent(chainParams.StakingInfoAddress.EthAddress(), receipt, msg.LogIndex)
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
	validator, found := k.sk.GetValidatorFromValID(ctx, msg.ID)
	if found {
		signer = validator.Signer
	}

	// create topup amount
	topupAmount := sdk.Coins{sdk.Coin{Denom: authTypes.FeeToken, Amount: sdk.NewIntFromBigInt(eventLog.Fee)}}

	// sequence id

	sequence := new(big.Int).Mul(receipt.BlockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// check if incoming tx already exists
	if k.HasTopupSequence(ctx, sequence.String()) {
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
	k.SetTopupSequence(ctx, sequence.String())

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

	amount := msg.Amount

	validator, err := k.sk.GetValidatorInfo(ctx, msg.ValidatorAddress.Bytes())
	if err != nil {
		return hmCommon.ErrInvalidMsg(k.Codespace(), "No validator found with signer %s", msg.ValidatorAddress.String()).Result()
	}

	if msg.Amount.String() == big.NewInt(0).String() {
		// fetch balance
		// check if fee is already withdrawn
		coins := k.bk.GetCoins(ctx, msg.ValidatorAddress)
		amount = coins.AmountOf(authTypes.FeeToken)
	}

	k.Logger(ctx).Info("Fee amount for ", "fromAddress", msg.ValidatorAddress, "validatorId", validator.ID, "balance", amount.BigInt().String())
	if amount.IsZero() {
		return types.ErrNoBalanceToWithdraw(k.Codespace()).Result()
	}

	// withdraw coins of validator
	maticCoins := sdk.Coins{sdk.Coin{Denom: authTypes.FeeToken, Amount: amount}}
	if _, err := k.bk.SubtractCoins(ctx, msg.ValidatorAddress, maticCoins); err != nil {
		k.Logger(ctx).Error("Error while setting Fee balance to zero ", "fromAddress", msg.ValidatorAddress, "validatorId", validator.ID, "err", err)
		return err.Result()
	}

	// Add Fee to Dividend Account
	feeAmount := amount.BigInt()
	if err := k.sk.AddFeeToDividendAccount(ctx, validator.ID, feeAmount); err != nil {
		k.Logger(ctx).Error("handleMsgWithdrawFee | AddFeeToDividendAccount", "fromAddress", msg.ValidatorAddress, "validatorId", validator.ID, "err", err)
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeFeeWithdraw,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyValidatorID, validator.ID.String()),
			sdk.NewAttribute(types.AttributeKeyValidatorSigner, msg.ValidatorAddress.String()),
			sdk.NewAttribute(types.AttributeKeyFeeWithdrawAmount, feeAmount.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
