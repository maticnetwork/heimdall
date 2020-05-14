package topup

import (
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

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
			return HandleMsgTopup(ctx, k, msg, contractCaller)
		case types.MsgWithdrawFee:
			return HandleMsgWithdrawFee(ctx, k, msg)
		default:
			errMsg := "Unrecognized topup Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// HandleMsgTopup handles topup event
func HandleMsgTopup(ctx sdk.Context, k Keeper, msg types.MsgTopup, contractCaller helper.IContractCaller) sdk.Result {

	k.Logger(ctx).Debug("✅ Validating topup msg",
		"User", msg.User,
		"Fee", msg.Fee,
		"txHash", hmTypes.BytesToHeimdallHash(msg.TxHash.Bytes()),
		"logIndex", uint64(msg.LogIndex),
		"blockNumber", msg.BlockNumber,
	)

	if !k.bk.GetSendEnabled(ctx) {
		return types.ErrSendDisabled(k.Codespace()).Result()
	}

	// sequence id
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// check if incoming tx already exists
	if k.HasTopupSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// HandleMsgWithdrawFee handle withdraw fee event
func HandleMsgWithdrawFee(ctx sdk.Context, k Keeper, msg types.MsgWithdrawFee) sdk.Result {

	// partial withdraw
	amount := msg.Amount

	validator, err := k.sk.GetValidatorInfo(ctx, msg.ValidatorAddress.Bytes())
	if err != nil {
		return hmCommon.ErrInvalidMsg(k.Codespace(), "No validator found with signer %s", msg.ValidatorAddress.String()).Result()
	}

	// full withdraw
	if msg.Amount.String() == big.NewInt(0).String() {
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
	if err := k.AddFeeToDividendAccount(ctx, validator.ID, feeAmount); err != nil {
		k.Logger(ctx).Error("handleMsgWithdrawFee | AddFeeToDividendAccount", "fromAddress", msg.ValidatorAddress, "validatorId", validator.ID, "err", err)
		return err.Result()
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeFeeWithdraw,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyValidatorID, validator.ID.String()),
			sdk.NewAttribute(types.AttributeKeyValidatorUser, msg.ValidatorAddress.String()),
			sdk.NewAttribute(types.AttributeKeyFeeWithdrawAmount, feeAmount.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
