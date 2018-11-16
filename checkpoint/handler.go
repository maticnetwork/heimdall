package checkpoint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCheckpoint:
			return handleMsgCheckpoint(ctx, msg, k)
		case MsgCheckpointAck:
			return handleMsgCheckpointAck(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("Invalid message in checkpoint module").Result()
		}
	}
}
func handleMsgCheckpointAck(ctx sdk.Context, msg MsgCheckpointAck, k Keeper) sdk.Result {
	// make call to headerBlock with header number
	// validate via checking details with lastCheckpoint on TM
	return sdk.Result{}
}

func handleMsgCheckpoint(ctx sdk.Context, msg MsgCheckpoint, k Keeper) sdk.Result {
	if err := msg.ValidateBasic(); err != nil {
		return ErrBadBlockDetails(k.codespace).Result()
	}
	key := k.AddCheckpoint(ctx, msg.StartBlock, msg.EndBlock, msg.RootHash)
	CheckpointLogger.Debug("Checkpoint added in state", "key", key)

	// send tags
	return sdk.Result{}
}
