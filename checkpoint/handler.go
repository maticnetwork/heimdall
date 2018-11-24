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
	root, start, end, err := helper.GetHeaderInfo(msg.HeaderBlock)
	if err != nil {
		CheckpointLogger.Error("Unable to fetch header from rootchain contract", "Error", err, "HeaderBlockIndex", msg.HeaderBlock)
		return ErrBadAck(k.codespace).Result()
	}

	CheckpointLogger.Debug("HeaderBlock Fetched", "start", start, "end", end, "Roothash", root)

	// get last checkpoint from buffer
	headerBlock, err := k.GetCheckpointFromBuffer(ctx)
	if err != nil {
		CheckpointLogger.Error("Unable to get checkpoint", "error", err, "key", BufferCheckpointKey)
	}

	// match header block and checkpoint
	if start != headerBlock.StartBlock || end != headerBlock.EndBlock || strings.Compare(root.String(), headerBlock.RootHash.String()) != 0 {
		CheckpointLogger.Error("Invalid ACK", "StartExpected", headerBlock.StartBlock, "StartReceived", start,
			"EndExpected", headerBlock.EndBlock, "EndReceived", end, "RootExpected", root.String(), "RootRecieved", headerBlock.RootHash.String())
		return ErrBadAck(k.codespace).Result()
	}

	// add checkpoint to headerBlocks
	k.AddCheckpointToKey(ctx, headerBlock.StartBlock, headerBlock.EndBlock, headerBlock.RootHash, headerBlock.Proposer, GetHeaderKey(int(msg.HeaderBlock)))
	CheckpointLogger.Info("Checkpoint Added to Store", "roothash", headerBlock.RootHash, "startBlock",
		headerBlock.StartBlock, "endBlock", headerBlock.EndBlock, "proposer", headerBlock.Proposer)

	// flush checkpoint in buffer
	k.FlushCheckpointBuffer(ctx)
	CheckpointLogger.Debug("Checkpoint Buffer Flushed", "Checkpoint", headerBlock)

	// update ack count
	k.UpdateACKCount(ctx)
	CheckpointLogger.Debug("Valid ACK Received", "CurrentACKCount", k.GetACKCount(ctx)-1, "UpdatedACKCount", k.GetACKCount(ctx))

	// check for validator updates

	// if found create new validator set and replace

	// indicate ACK received by adding in cache , cache cleared in endblock
	k.SetCheckpointAckCache(ctx, CacheExistsValue)

	return sdk.Result{}
}

func handleMsgCheckpoint(ctx sdk.Context, msg MsgCheckpoint, k Keeper) sdk.Result {
	// validate checkpoint
	if !ValidateCheckpoint(msg.StartBlock, msg.EndBlock, msg.RootHash.String()) {
		CheckpointLogger.Error("RootHash Not Valid", "StartBlock", msg.StartBlock, "EndBlock", msg.EndBlock, "RootHash", msg.RootHash)
		return ErrBadBlockDetails(k.codespace).Result()
	}

	// fetch last checkpoint from store
	lastCheckpoint := k.GetLastCheckpoint(ctx)

	// make sure new checkpoint is after tip
	if lastCheckpoint.EndBlock > msg.StartBlock {
		CheckpointLogger.Error("Checkpoint already exists", "CurrentTip", lastCheckpoint.EndBlock, "MsgStartBlock", msg.StartBlock)
		return ErrBadBlockDetails(k.codespace).Result()
	}

	// add checkpoint to buffer
	k.AddCheckpointToKey(ctx, msg.StartBlock, msg.EndBlock, msg.RootHash, msg.Proposer, BufferCheckpointKey)
	CheckpointLogger.Debug("Checkpoint added in buffer!", "roothash", msg.RootHash, "startBlock",
		msg.StartBlock, "endBlock", msg.EndBlock, "proposer", msg.Proposer)

	// indicate Checkpoint received by adding in cache , cache cleared in endblock
	k.SetCheckpointCache(ctx, CacheExistsValue)

	// send tags
	return sdk.Result{}
}
