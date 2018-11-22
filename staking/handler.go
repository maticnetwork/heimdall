package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgValidatorJoin:
			return handleMsgValidatorJoin(ctx, msg, k)
		case MsgValidatorExit:
			return handleMsgValidatorExit(ctx, msg, k)
		case MsgValidatorUpdate:
			return handleMsgValidatorUpdate(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("Invalid message in checkpoint module").Result()
		}
	}
}
func handleMsgValidatorUpdate(context sdk.Context, update MsgValidatorUpdate, keeper Keeper) sdk.Result {
	// verify from mainchain
	return sdk.Result{}
}
func handleMsgValidatorExit(context sdk.Context, exit MsgValidatorExit, keeper Keeper) sdk.Result {
	// verify deactivation from ACK count
	return sdk.Result{}
}
func handleMsgValidatorJoin(ctx sdk.Context, msg MsgValidatorJoin, k Keeper) sdk.Result {
	// fetch validator from mainchain
	validator, err := helper.GetValidatorInfo(msg.ValidatorAddr)
	if err != nil {
		return ErrNoValidator(k.codespace).Result()
	}

	// validate if start epoch is after current tip
	ACKs := k.checkpointKeeper.GetACKCount(ctx)
	if int(validator.StartEpoch) < ACKs {
		return ErrOldValidator(k.codespace).Result()
	}

	// validate pubkey matches signer address

	// add validator to store

	return sdk.Result{}
}
