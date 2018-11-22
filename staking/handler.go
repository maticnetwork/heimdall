package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper, checkpointKeeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgValidatorJoin:
			return handleMsgValidatorJoin(ctx, msg, k, checkpointKeeper)
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
func handleMsgValidatorJoin(context sdk.Context, join MsgValidatorJoin, keeper Keeper, checkpointKeeper Keeper) sdk.Result {
	// validate if start epoch is after current tip

	// fetch validator from mainchain

	// validate pubkey matches signer address

	// add validator to store

	return sdk.Result{}
}
