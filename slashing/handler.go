package slashing

import (

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/helper"
)

// NewHandler creates an sdk.Handler for all the slashing type messages
func NewHandler(k Keeper, contractCaller helper.IContractCaller) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		return sdk.ErrTxDecode("Invalid message in slashing module").Result()
		// switch msg := msg.(type) {
/* 		case types.MsgUnjail:
			return handleMsgUnjail(ctx, msg, k, contractCaller)
		case types.MsgTick:
			return handlerMsgTick(ctx, msg, k, contractCaller)
		case types.MsgTickAck:
			return handleMsgTickAck(ctx, msg, k, contractCaller) */
		// default:
		// 	return sdk.ErrTxDecode("Invalid message in slashing module").Result()
		// }
	}
}

