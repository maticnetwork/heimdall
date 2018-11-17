package checkpoint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
	"strings"
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

	// get last checkpoint
	key := k.GetLastCheckpointKey(ctx)
	headerBlock, err := k.GetCheckpoint(ctx, key)
	if err != nil {
		CheckpointLogger.Error("Unable to get checkpoint", "error", err, "key", key)
	}

	// match header block and checkpoint
	if start != headerBlock.StartBlock || end != headerBlock.EndBlock || strings.Compare(root.String(), headerBlock.RootHash.String()) != 0 {
		CheckpointLogger.Error("Invalid ACK", "StartExpected", headerBlock.StartBlock, "StartReceived", start, "End", headerBlock.EndBlock, end)
		return ErrBadAck(k.codespace).Result()
	}

	// update ack count
	CheckpointLogger.Debug("Valid ACK", "CurrentACKCount", k.GetACKCount(ctx), "UpdatedACKCount", k.GetACKCount(ctx)+1)
	k.UpdateACKCount(ctx)

	return sdk.Result{}
}

func handleMsgCheckpoint(ctx sdk.Context, msg MsgCheckpoint, k Keeper) sdk.Result {

	if err := msg.ValidateBasic(); err != nil {
		return ErrBadBlockDetails(k.codespace).Result()
	}
	key := k.AddCheckpoint(ctx, msg.StartBlock, msg.EndBlock, msg.RootHash, msg.Proposer)
	CheckpointLogger.Debug("Checkpoint added in state", "key", key)

	// send tags
	return sdk.Result{}
}
