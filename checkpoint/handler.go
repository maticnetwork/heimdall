package checkpoint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
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
	root, start, end, _ := helper.GetHeaderInfo(msg.HeaderBlock)
	key := k.GetLastCheckpointKey(ctx)
	headerBlock, err := k.GetCheckpoint(ctx, key)
	if err != nil {
		CheckpointLogger.Error("Unable to get checkpoint", "error", err, "key", key)
	}
	// TODO add roothash validation
	if start != headerBlock.StartBlock || end != headerBlock.EndBlock {
		CheckpointLogger.Error("Invalid ACK", "Start", headerBlock.StartBlock, start, "End", headerBlock.EndBlock, end)
		return ErrBadAck(k.codespace).Result()
	}
	CheckpointLogger.Debug("Valid ACK , updating count")

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
