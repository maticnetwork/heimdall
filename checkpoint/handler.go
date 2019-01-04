package checkpoint

import (
	"bytes"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

func NewHandler(k common.Keeper,contractCallerObj helper.ContractCallerObj) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCheckpoint:
			return HandleMsgCheckpoint(ctx, msg, k,contractCallerObj)
		case MsgCheckpointAck:
			return handleMsgCheckpointAck(ctx, msg, k,contractCallerObj)
		case MsgCheckpointNoAck:
			return handleMsgCheckpointNoAck(ctx, msg, k,contractCallerObj)
		default:
			return sdk.ErrTxDecode("Invalid message in checkpoint module").Result()
		}
	}
}

func handleMsgCheckpointAck(ctx sdk.Context, msg MsgCheckpointAck, k common.Keeper,contractCaller helper.ContractCallerObj) sdk.Result {

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

	return sdk.Result{}
}

func HandleMsgCheckpoint(ctx sdk.Context, msg MsgCheckpoint, k common.Keeper,contractCaller helper.ContractCallerObj) sdk.Result {
	if msg.TimeStamp == 0 || msg.TimeStamp > uint64(time.Now().Unix()) {
		return common.ErrBadTimeStamp(k.Codespace).Result()
	}

	checkpointBuffer, err := k.GetCheckpointFromBuffer(ctx)
	if err == nil {
		if msg.TimeStamp == 0 || checkpointBuffer.TimeStamp == 0 || ((msg.TimeStamp > checkpointBuffer.TimeStamp) && msg.TimeStamp-checkpointBuffer.TimeStamp > uint64(helper.CheckpointBufferTime.Seconds())) {
			k.FlushCheckpointBuffer(ctx)
		} else {
			return common.ErrNoACK(k.Codespace).Result()
		}
	}

	// validate checkpoint
	if !ValidateCheckpoint(msg.StartBlock, msg.EndBlock, msg.RootHash) {
		common.CheckpointLogger.Error("RootHash is not valid",
			"StartBlock", msg.StartBlock,
			"EndBlock", msg.EndBlock,
			"RootHash", msg.RootHash)
		return common.ErrBadBlockDetails(k.Codespace).Result()
	}

	// fetch last checkpoint from store
	if lastCheckpoint, err := k.GetLastCheckpoint(ctx); err == nil {
		// make sure new checkpoint is after tip
		if lastCheckpoint.EndBlock > msg.StartBlock {
			common.CheckpointLogger.Error("Checkpoint already exists",
				"currentTip", lastCheckpoint.EndBlock,
				"startBlock", msg.StartBlock)
			return common.ErrBadBlockDetails(k.Codespace).Result()
		}
	}

	// check proposer in message
	if !bytes.Equal(msg.Proposer.Bytes(), k.GetValidatorSet(ctx).Proposer.Signer.Bytes()) {
		common.CheckpointLogger.Error("Invalid proposer in message",
			"currentProposer", k.GetValidatorSet(ctx).Proposer.Signer.String(),
			"checkpointProposer", msg.Proposer.String())
		return common.ErrBadProposerDetails(k.Codespace, k.GetValidatorSet(ctx).Proposer.Signer).Result()
	}

	// check if proposer has min ether
	balance, _ := contractCaller.GetBalance(msg.Proposer)
	if balance.Cmp(helper.MinBalance) == -1 {
		common.CheckpointLogger.Error("Proposer doesnt have enough ether to send checkpoint tx", "Balance", balance, "RequiredBalance", helper.MinBalance)
		return common.ErrLowBalance(k.Codespace, msg.Proposer.String()).Result()
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

func handleMsgCheckpointNoAck(ctx sdk.Context, msg MsgCheckpointNoAck, k common.Keeper,contractCaller helper.ContractCallerObj) sdk.Result {
	// current time
	currentTime := time.Unix(int64(msg.TimeStamp), 0) // buffer time
	bufferTime := helper.CheckpointBufferTime

	// fetch last checkpoint from store
	// TODO figure out how to handle this error
	lastCheckpoint, _ := k.GetLastCheckpoint(ctx)
	lastCheckpointTime := time.Unix(int64(lastCheckpoint.TimeStamp), 0)

	// if last checkpoint is not present or last checkpoint happens before checkpoint buffer time -- thrown an error
	if lastCheckpointTime.After(currentTime) || (currentTime.Sub(lastCheckpointTime) < bufferTime) {
		return common.ErrInvalidNoACK(k.Codespace).Result()
	}

	// check last no ack - prevents repetitive no-ack
	lastAck := k.GetLastNoAck(ctx)
	lastAckTime := time.Unix(int64(lastAck), 0)

	if lastAckTime.After(currentTime) || (currentTime.Sub(lastAckTime) < bufferTime) {
		return common.ErrTooManyNoACK(k.Codespace).Result()
	}

	// set last no ack
	k.SetLastNoAck(ctx, uint64(currentTime.Unix()))

	// flush buffer
	k.FlushCheckpointBuffer(ctx)
	common.CheckpointLogger.Debug("Checkpoint buffer flushed after receiving no-ack")

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
