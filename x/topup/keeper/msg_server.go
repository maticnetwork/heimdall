package keeper

import (
	"context"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"

	// authTypes "github.com/maticnetwork/heimdall/auth/types"
	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"

	// "github.com/maticnetwork/heimdall/topup/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/x/topup/types"
)

type msgServer struct {
	Keeper
	contractCaller helper.IContractCaller
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper, contractCaller helper.IContractCaller) types.MsgServer {
	return &msgServer{Keeper: keeper, contractCaller: contractCaller}
}

var _ types.MsgServer = msgServer{}

// // NewHandler returns a handler for "topup" type messages.
// func NewHandler(k Keeper, contractCaller helper.IContractCaller) sdk.Handler {
// 	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
// 		ctx = ctx.WithEventManager(sdk.NewEventManager())

// 		switch msg := msg.(type) {
// 		case types.MsgTopup:
// 			return HandleMsgTopup(ctx, k, msg, contractCaller)
// 		case types.MsgWithdrawFee:
// 			return HandleMsgWithdrawFee(ctx, k, msg)
// 		default:
// 			return sdk.ErrUnknownRequest("Unrecognized topup msg type").Result()
// 		}
// 	}
// }

// HandleMsgTopup handles topup event
func (k msgServer) Topup(goCtx context.Context, msg *types.MsgTopup) (*types.MsgTopupResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	k.Logger(ctx).Debug("✅ Validating topup msg",
		"User", msg.User,
		"Fee", msg.Fee,
		"txHash", msg.TxHash,
		"logIndex", uint64(msg.LogIndex),
		"blockNumber", msg.BlockNumber,
	)

	// TODO: Is this still relevant now that we are using bank module of cosmos?
	// if !k.bk.GetSendEnabled(ctx) {
	// 	return types.ErrSendDisabled(k.Codespace()).Result()
	// }

	// sequence id
	blockNumber := new(big.Int).SetUint64(msg.BlockNumber)
	sequence := new(big.Int).Mul(blockNumber, big.NewInt(hmTypes.DefaultLogIndexUnit))
	sequence.Add(sequence, new(big.Int).SetUint64(msg.LogIndex))

	// check if incoming tx already exists
	if k.HasTopupSequence(ctx, sequence.String()) {
		k.Logger(ctx).Error("Older invalid tx found")
		return nil, hmCommon.ErrOldTx
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeTopup,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySender, msg.FromAddress.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.User.String()),
			sdk.NewAttribute(types.AttributeKeyTopupAmount, msg.Fee.String()),
		),
	})

	// return sdk.Result{
	// 	Events: ctx.EventManager().Events(),
	// }
	return &types.MsgTopupResponse{}, nil
}

// HandleMsgWithdrawFee handle withdraw fee event
func (k msgServer) WithdrawFee(goCtx context.Context, msg *types.MsgWithdrawFee) (*types.MsgWithdrawFeeResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)
	// partial withdraw
	amount := msg.Amount

	// full withdraw
	if msg.Amount.String() == big.NewInt(0).String() {
		coins := k.bk.GetAllBalances(ctx, msg.UserAddress)
		amount = coins.AmountOf(hmTypes.FeeToken)
	}

	k.Logger(ctx).Debug("Fee amount", "fromAddress", msg.UserAddress, "balance", amount.BigInt().String())
	if amount.IsZero() {
		return nil, types.ErrNoBalanceToWithdraw
	}

	// withdraw coins of validator
	maticCoins := sdk.Coins{sdk.Coin{Denom: hmTypes.FeeToken, Amount: amount}}
	if err := k.bk.SubtractCoins(ctx, msg.UserAddress, maticCoins); err != nil {
		k.Logger(ctx).Error("Error while setting Fee balance to zero ", "fromAddress", msg.UserAddress, "err", err)
		return nil, types.ErrSetFeeBalanceZero
	}

	// Add Fee to Dividend Account
	feeAmount := amount.BigInt()
	if err := k.AddFeeToDividendAccount(ctx, msg.UserAddress, feeAmount); err != nil {
		k.Logger(ctx).Error("WithdrawFee | AddFeeToDividendAccount", "fromAddress", msg.UserAddress, "err", err)
		return nil, types.ErrAddFeeToDividendAccount
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeFeeWithdraw,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyUser, msg.UserAddress.String()),
			sdk.NewAttribute(types.AttributeKeyFeeWithdrawAmount, feeAmount.String()),
		),
	})

	// return sdk.Result{
	// 	Events: ctx.EventManager().Events(),
	// }
	return &types.MsgWithdrawFeeResponse{}, nil
}
