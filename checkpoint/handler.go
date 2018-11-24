package checkpoint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"strings"
)

func NewHandler(k common.Keeper) sdk.Handler {
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

func handleMsgCheckpointAck(ctx sdk.Context, msg MsgCheckpointAck, k common.Keeper) sdk.Result {
	// make call to headerBlock with header number
	root, start, end, err := helper.GetHeaderInfo(msg.HeaderBlock)
	if err != nil {
		common.CheckpointLogger.Error("Unable to fetch header from rootchain contract", "Error", err, "HeaderBlockIndex", msg.HeaderBlock)
		return common.ErrBadAck(k.Codespace).Result()
	}

	common.CheckpointLogger.Debug("HeaderBlock Fetched", "start", start, "end", end, "Roothash", root)

	// get last checkpoint from buffer
	headerBlock, err := k.GetCheckpointFromBuffer(ctx)
	if err != nil {
		common.CheckpointLogger.Error("Unable to get checkpoint", "error", err, "key", common.BufferCheckpointKey)
	}

	// match header block and checkpoint
	if start != headerBlock.StartBlock || end != headerBlock.EndBlock || strings.Compare(root.String(), headerBlock.RootHash.String()) != 0 {
		common.CheckpointLogger.Error("Invalid ACK", "StartExpected", headerBlock.StartBlock, "StartReceived", start,
			"EndExpected", headerBlock.EndBlock, "EndReceived", end, "RootExpected", root.String(), "RootRecieved", headerBlock.RootHash.String())
		return common.ErrBadAck(k.Codespace).Result()
	}

	// add checkpoint to headerBlocks
	k.AddCheckpointToKey(ctx, headerBlock.StartBlock, headerBlock.EndBlock, headerBlock.RootHash, headerBlock.Proposer, common.GetHeaderKey(int(msg.HeaderBlock)))
	common.CheckpointLogger.Info("Checkpoint Added to Store", "roothash", headerBlock.RootHash, "startBlock",
		headerBlock.StartBlock, "endBlock", headerBlock.EndBlock, "proposer", headerBlock.Proposer)

	// flush checkpoint in buffer
	k.FlushCheckpointBuffer(ctx)
	common.CheckpointLogger.Debug("Checkpoint Buffer Flushed", "Checkpoint", headerBlock)

	// update ack count
	k.UpdateACKCount(ctx)
	common.CheckpointLogger.Debug("Valid ACK Received", "CurrentACKCount", k.GetACKCount(ctx)-1, "UpdatedACKCount", k.GetACKCount(ctx))

	// check for validator updates
	if k.ValidatorSetChanged(ctx) {
		// GetAllValidators from store , not current , ALL !
		updatedValidators := k.GetAllValidators(ctx)

		// get current running validator set
		currentValidatorSet := k.GetValidatorSet(ctx)

		// apply updates
		helper.UpdateValidators(&currentValidatorSet, updatedValidators)

		// update validator set in store
		k.UpdateValidatorSetInStore(ctx, currentValidatorSet)

		// indicate validator set changes in state have been done
		k.SetValidatorSetChangedFlag(ctx, false)
	} else {
		// if no updates found increment accum
		k.IncreamentAccum(ctx, 1)
	}

	// indicate ACK received by adding in cache , cache cleared in endblock
	k.SetCheckpointAckCache(ctx, common.DefaultValue)

	return sdk.Result{}
}

func handleMsgCheckpoint(ctx sdk.Context, msg MsgCheckpoint, k common.Keeper) sdk.Result {
	// validate checkpoint
	if !ValidateCheckpoint(msg.StartBlock, msg.EndBlock, msg.RootHash.String()) {
		common.CheckpointLogger.Error("RootHash Not Valid", "StartBlock", msg.StartBlock, "EndBlock", msg.EndBlock, "RootHash", msg.RootHash)
		return common.ErrBadBlockDetails(k.Codespace).Result()
	}

	// fetch last checkpoint from store
	lastCheckpoint := k.GetLastCheckpoint(ctx)

	// make sure new checkpoint is after tip
	if lastCheckpoint.EndBlock > msg.StartBlock {
		common.CheckpointLogger.Error("Checkpoint already exists", "CurrentTip", lastCheckpoint.EndBlock, "MsgStartBlock", msg.StartBlock)
		return common.ErrBadBlockDetails(k.Codespace).Result()
	}

	// check proposer in message
	if msg.Proposer.String() != k.GetValidatorSet(ctx).Proposer.Address.String() {
		common.CheckpointLogger.Error("Invalid proposer in message", "CurrentProposer", k.GetValidatorSet(ctx).Proposer.Address.String(),
			"CheckpointProposer", msg.Proposer.String())
		return common.ErrBadBlockDetails(k.Codespace).Result()
	}

	// add checkpoint to buffer
	k.AddCheckpointToKey(ctx, msg.StartBlock, msg.EndBlock, msg.RootHash, msg.Proposer, common.BufferCheckpointKey)
	common.CheckpointLogger.Debug("Checkpoint added in buffer!", "roothash", msg.RootHash, "startBlock",
		msg.StartBlock, "endBlock", msg.EndBlock, "proposer", msg.Proposer)

	// indicate Checkpoint received by adding in cache , cache cleared in endblock
	k.SetCheckpointCache(ctx, common.DefaultValue)

	// send tags
	return sdk.Result{}
}
