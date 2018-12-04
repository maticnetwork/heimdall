package checkpoint

import (
	"bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
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

	common.CheckpointLogger.Debug("HeaderBlock fetched", "headerBlock", msg.HeaderBlock, "start", start, "end", end, "Roothash", root)

	// get last checkpoint from buffer
	headerBlock, err := k.GetCheckpointFromBuffer(ctx)
	if err != nil {
		common.CheckpointLogger.Error("Unable to get checkpoint", "error", err)
		return common.ErrBadAck(k.Codespace).Result()
	}

	// match header block and checkpoint
	if start != headerBlock.StartBlock || end != headerBlock.EndBlock || !bytes.Equal(root[:], headerBlock.RootHash[:]) {
		common.CheckpointLogger.Error("Invalid ACK", "startExpected", headerBlock.StartBlock, "startReceived", start, "endExpected", headerBlock.EndBlock, "endReceived", end, "rootExpected", root.String(), "rootRecieved", headerBlock.RootHash.String())
		return common.ErrBadAck(k.Codespace).Result()
	}

	// add checkpoint to headerBlocks
	k.AddCheckpoint(ctx, msg.HeaderBlock, headerBlock)
	common.CheckpointLogger.Info("Checkpoint added to store", "roothash", headerBlock.RootHash, "startBlock", headerBlock.StartBlock, "endBlock", headerBlock.EndBlock, "proposer", headerBlock.Proposer)

	// flush checkpoint in buffer
	k.FlushCheckpointBuffer(ctx)
	common.CheckpointLogger.Debug("Checkpoint buffer flushed after receiving checkpoint ack", "checkpoint", headerBlock)

	// update ack count
	k.UpdateACKCount(ctx)
	common.CheckpointLogger.Debug("Valid ack received", "CurrentACKCount", k.GetACKCount(ctx)-1, "UpdatedACKCount", k.GetACKCount(ctx))

	// --- Update to new validator

	// remove matured Validators
	k.RemoveDeactivatedValidators(ctx)

	// get current running validator set
	currentValidatorSet := k.GetValidatorSet(ctx)

	// apply updates
	helper.UpdateValidators(
		&currentValidatorSet,           // pointer to current validator set -- UpdateValidators will modify it
		k.GetAllValidators(ctx),        // All validators
		k.GetValidatorToSignerMap(ctx), // validator to signer map
	)

	// update validator set in store
	k.UpdateValidatorSetInStore(ctx, currentValidatorSet)

	// increment accum
	k.IncreamentAccum(ctx, 1)

	//log new proposer
	vs := k.GetValidatorSet(ctx)
	newProposer := vs.GetProposer()
	common.CheckpointLogger.Debug(
		"New proposer selected",
		"validator", newProposer.Signer.String(),
		"signer", newProposer.Signer.String(),
		"power", newProposer.Power,
	)

	// --- End

	// indicate ACK received by adding in cache, cache cleared in endblock
	k.SetCheckpointAckCache(ctx, common.DefaultValue)

	return sdk.Result{}
}

func handleMsgCheckpoint(ctx sdk.Context, msg MsgCheckpoint, k common.Keeper) sdk.Result {
	// validate checkpoint
	if !ValidateCheckpoint(msg.StartBlock, msg.EndBlock, msg.RootHash) {
		common.CheckpointLogger.Error("RootHash is not valid", "StartBlock", msg.StartBlock, "EndBlock", msg.EndBlock, "RootHash", msg.RootHash)
		return common.ErrBadBlockDetails(k.Codespace).Result()
	}

	// fetch last checkpoint from store
	if lastCheckpoint, err := k.GetLastCheckpoint(ctx); err == nil {
		// make sure new checkpoint is after tip
		if lastCheckpoint.EndBlock > msg.StartBlock {
			common.CheckpointLogger.Error("Checkpoint already exists", "currentTip", lastCheckpoint.EndBlock, "startBlock", msg.StartBlock)
			return common.ErrBadBlockDetails(k.Codespace).Result()
		}
	}

	// check proposer in message
	if !bytes.Equal(msg.Proposer.Bytes(), k.GetValidatorSet(ctx).Proposer.Signer.Bytes()) {
		common.CheckpointLogger.Error("Invalid proposer in message", "currentProposer", k.GetValidatorSet(ctx).Proposer.Signer.String(), "checkpointProposer", msg.Proposer.String())
		return common.ErrBadProposerDetails(k.Codespace).Result()
	}

	// add checkpoint to buffer
	k.SetCheckpointBuffer(ctx, hmTypes.CheckpointBlockHeader{
		StartBlock: msg.StartBlock,
		EndBlock:   msg.EndBlock,
		RootHash:   msg.RootHash,
		Proposer:   msg.Proposer,
		TimeStamp:  msg.TimeStamp,
	})

	// indicate Checkpoint received by adding in cache, cache cleared in endblock
	k.SetCheckpointCache(ctx, common.DefaultValue)

	// send tags
	return sdk.Result{}
}
