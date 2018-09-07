package checkpoint


import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"fmt"
)


func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		// NOTE msg already has validate basic run
		switch msg := msg.(type) {
		case MsgCheckpoint:
			return handleMsgCheckpoint(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("Invalid message in checkpoint module ").Result()
		}
	}
}

func handleMsgCheckpoint(ctx sdk.Context, msg MsgCheckpoint, k Keeper) sdk.Result {
	fmt.Printf("entered handler with message %v",msg)
	//TODO validation
	//k.addCheckpoint(ctx,msg.StartBlock,msg.EndBlock,msg.rootHash)
	//TODO add validation
	// send tags
	return sdk.Result{}
}
