package checkpoint

import (
	"bytes"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

func NewHandler(k common.Keeper, contractCaller helper.IContractCaller) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCheckpoint:
			return HandleMsgCheckpoint(ctx, msg, k, contractCaller)
		case MsgCheckpointAck:
			return HandleMsgCheckpointAck(ctx, msg, k, contractCaller)
		case MsgCheckpointNoAck:
			return HandleMsgCheckpointNoAck(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("Invalid message in checkpoint module").Result()
		}
	}
}

func HandleMsgCheckpointAck(ctx sdk.Context, msg MsgCheckpointAck, k common.Keeper, contractCaller helper.IContractCaller) sdk.Result {
	common.CheckpointLogger.Debug("handling ack message", "Msg", msg)
	// make call to headerBlock with header number
	root, start, end, err := contractCaller.GetHeaderInfo(msg.HeaderBlock)
	if err != nil {
		common.CheckpointLogger.Error("Unable to fetch header from rootchain contract", "Error", err, "HeaderBlockIndex", msg.HeaderBlock)
		return common.ErrBadAck(k.Codespace).Result()
	}

	common.CheckpointLogger.Debug("HeaderBlock fetched",
		"headerBlock", msg.HeaderBlock,
		"start", start,
		"end", end,
		"Roothash", root)

	// get last checkpoint from buffer
	headerBlock, err := k.GetCheckpointFromBuffer(ctx)
	if err != nil {
		common.CheckpointLogger.Error("Unable to get checkpoint", "error", err)
		return common.ErrBadAck(k.Codespace).Result()
	}

	// match header block and checkpoint
	if start != headerBlock.StartBlock || end != headerBlock.EndBlock || !bytes.Equal(root.Bytes(), headerBlock.RootHash.Bytes()) {
		common.CheckpointLogger.Error("Invalid ACK",
			"startExpected", headerBlock.StartBlock,
			"startReceived", start,
			"endExpected", headerBlock.EndBlock,
			"endReceived", end,
			"rootExpected", root.String(),
			"rootRecieved", headerBlock.RootHash.String())

		return common.ErrBadAck(k.Codespace).Result()
	}

	// add checkpoint to headerBlocks
	k.AddCheckpoint(ctx, msg.HeaderBlock, headerBlock)
	common.CheckpointLogger.Info("Checkpoint added to store", "headerBlock", headerBlock.String())

	// flush buffer
	k.FlushCheckpointBuffer(ctx)
	common.CheckpointLogger.Debug("Checkpoint buffer flushed after receiving checkpoint ack", "checkpoint", headerBlock)

	// update ack count
	k.UpdateACKCount(ctx)
	common.CheckpointLogger.Debug("Valid ack received", "CurrentACKCount", k.GetACKCount(ctx)-1, "UpdatedACKCount", k.GetACKCount(ctx))

	// indicate ACK received by adding in cache, cache cleared in endblock
	k.SetCheckpointAckCache(ctx, common.DefaultValue)
	common.CheckpointLogger.Debug("Checkpoint ACK cache set", "CacheValue", k.GetCheckpointCache(ctx, common.CheckpointACKCacheKey))

	return sdk.Result{}
}

func HandleMsgCheckpoint(ctx sdk.Context, msg MsgCheckpoint, k common.Keeper, contractCaller helper.IContractCaller) sdk.Result {
	common.CheckpointLogger.Debug("Handling checkpoint msg", "Msg", msg)
	if msg.TimeStamp == 0 || msg.TimeStamp > uint64(time.Now().Unix()) {
		common.CheckpointLogger.Error("Checkpoint timestamp must be in near past", "CurrentTime", time.Now().Unix(), "CheckpointTime", msg.TimeStamp, "Condition", msg.TimeStamp >= uint64(time.Now().Unix()))
		return common.ErrBadTimeStamp(k.Codespace).Result()
	}
	common.CheckpointLogger.Debug("Valid Timestamp")

	checkpointBuffer, err := k.GetCheckpointFromBuffer(ctx)
	if err == nil {
		if msg.TimeStamp == 0 || checkpointBuffer.TimeStamp == 0 || ((msg.TimeStamp > checkpointBuffer.TimeStamp) && msg.TimeStamp-checkpointBuffer.TimeStamp >= uint64(helper.CheckpointBufferTime.Seconds())) {
			common.CheckpointLogger.Debug("Checkpoint has been timed out, flushing buffer", "CheckpointTimestamp", msg.TimeStamp, "PrevCheckpointTimestamp", checkpointBuffer.TimeStamp)
			k.FlushCheckpointBuffer(ctx)
		} else {
			checkpointTime := time.Unix(int64(checkpointBuffer.TimeStamp), 0)
			expiryTime := checkpointTime.Add(helper.CheckpointBufferTime)
			diff := expiryTime.Sub(time.Now()).Minutes
			common.CheckpointLogger.Error("Checkpoint already exits in buffer", "Checkpoint", checkpointBuffer.String())
			return common.ErrNoACK(k.Codespace, diff).Result()
		}
	}
	common.CheckpointLogger.Debug("Received checkpoint from buffer", "Checkpoint", checkpointBuffer.String())

	// validate checkpoint
	if !ValidateCheckpoint(msg.StartBlock, msg.EndBlock, msg.RootHash) {
		common.CheckpointLogger.Error("RootHash is not valid",
			"StartBlock", msg.StartBlock,
			"EndBlock", msg.EndBlock,
			"RootHash", msg.RootHash)
		return common.ErrBadBlockDetails(k.Codespace).Result()
	}
	common.CheckpointLogger.Debug("Valid Roothash in checkpoint", "StartBlock", msg.StartBlock, "EndBlock", msg.EndBlock)

	// fetch last checkpoint from store
	if lastCheckpoint, err := k.GetLastCheckpoint(ctx); err == nil {
		// make sure new checkpoint is after tip
		if lastCheckpoint.EndBlock > msg.StartBlock {
			common.CheckpointLogger.Error("Checkpoint already exists",
				"currentTip", lastCheckpoint.EndBlock,
				"startBlock", msg.StartBlock)
			return common.ErrBadBlockDetails(k.Codespace).Result()
		}
		if lastCheckpoint.EndBlock+1 != msg.StartBlock {
			common.CheckpointLogger.Error("Checkpoint not in countinuity",
				"currentTip", lastCheckpoint.EndBlock,
				"startBlock", msg.StartBlock)
			return common.ErrBadBlockDetails(k.Codespace).Result()

		}
	} else if err.Error() == common.ErrNoCheckpointFound(k.Codespace).Error() && msg.StartBlock != 0 {
		common.CheckpointLogger.Error("First checkpoint to start from block 1", "Error", err)
		return common.ErrBadBlockDetails(k.Codespace).Result()
	}
	common.CheckpointLogger.Debug("Valid checkpoint tip")

	// check proposer in message
	if !bytes.Equal(msg.Proposer.Bytes(), k.GetValidatorSet(ctx).Proposer.Signer.Bytes()) {
		common.CheckpointLogger.Error("Invalid proposer in message",
			"currentProposer", k.GetValidatorSet(ctx).Proposer.Signer.String(),
			"checkpointProposer", msg.Proposer.String())
		return common.ErrBadProposerDetails(k.Codespace, k.GetValidatorSet(ctx).Proposer.Signer).Result()
	}
	common.CheckpointLogger.Debug("Valid proposer in checkpoint")

	// check if proposer has min ether
	balance, _ := contractCaller.GetBalance(msg.Proposer)
	if balance.Cmp(helper.MinBalance) == -1 {
		common.CheckpointLogger.Error("Proposer doesnt have enough ether to send checkpoint tx", "Balance", balance, "RequiredBalance", helper.MinBalance)
		return common.ErrLowBalance(k.Codespace, msg.Proposer.String()).Result()
	}
	common.CheckpointLogger.Debug("Proposer has enough ether to send transaction")

	// add checkpoint to buffer
	k.SetCheckpointBuffer(ctx, hmTypes.CheckpointBlockHeader{
		StartBlock: msg.StartBlock,
		EndBlock:   msg.EndBlock,
		RootHash:   msg.RootHash,
		Proposer:   msg.Proposer,
		TimeStamp:  msg.TimeStamp,
	})

	checkpoint, _ := k.GetCheckpointFromBuffer(ctx)
	common.CheckpointLogger.Debug("Adding good checkpoint to buffer to await ACK", "checkpointStored", checkpoint.String())

	// indicate Checkpoint received by adding in cache, cache cleared in endblock
	k.SetCheckpointCache(ctx, common.DefaultValue)
	common.CheckpointLogger.Debug("Set Checkpoint Cache", k.GetCheckpointCache(ctx, common.CheckpointCacheKey))

	// send tags
	return sdk.Result{}
}
func HandleMsgCheckpointNoAck(ctx sdk.Context, msg MsgCheckpointNoAck, k common.Keeper) sdk.Result {
	common.CheckpointLogger.Debug("handling no-ack", "Msg", msg)
	// current time
	currentTime := time.Unix(int64(msg.TimeStamp), 0) // buffer time
	bufferTime := helper.CheckpointBufferTime

	// fetch last checkpoint from store
	// TODO figure out how to handle this error
	lastCheckpoint, _ := k.GetLastCheckpoint(ctx)
	lastCheckpointTime := time.Unix(int64(lastCheckpoint.TimeStamp), 0)

	// if last checkpoint is not present or last checkpoint happens before checkpoint buffer time -- thrown an error
	if lastCheckpointTime.After(currentTime) || (currentTime.Sub(lastCheckpointTime) < bufferTime) {
		common.CheckpointLogger.Debug("Invalid No ACK -- ongoing buffer period")
		return common.ErrInvalidNoACK(k.Codespace).Result()
	}

	// check last no ack - prevents repetitive no-ack
	lastAck := k.GetLastNoAck(ctx)
	lastAckTime := time.Unix(int64(lastAck), 0)

	if lastAckTime.After(currentTime) || (currentTime.Sub(lastAckTime) < bufferTime) {
		common.CheckpointLogger.Debug("Too many no-ack")
		return common.ErrTooManyNoACK(k.Codespace).Result()
	}

	// set last no ack
	k.SetLastNoAck(ctx, uint64(currentTime.Unix()))
	common.CheckpointLogger.Debug("Last No-ACK time set", "LastNoAck", k.GetLastNoAck(ctx))

	// --- Update to new proposer

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
	return sdk.Result{}
}
