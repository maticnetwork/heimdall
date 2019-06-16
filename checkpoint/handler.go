package checkpoint

import (
	"bytes"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

func NewHandler(k common.Keeper, contractCaller helper.IContractCaller) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		common.InitCheckpointLogger(&ctx)

		switch msg := msg.(type) {
		case MsgCheckpoint:
			return HandleMsgCheckpoint(ctx, msg, k, contractCaller, common.CheckpointLogger)
		case MsgCheckpointAck:
			return HandleMsgCheckpointAck(ctx, msg, k, contractCaller, common.CheckpointLogger)
		case MsgCheckpointNoAck:
			return HandleMsgCheckpointNoAck(ctx, msg, k, common.CheckpointLogger)
		default:
			return sdk.ErrTxDecode("Invalid message in checkpoint module").Result()
		}
	}
}

// Validates checkpoint transaction
func HandleMsgCheckpoint(ctx sdk.Context, msg MsgCheckpoint, k common.Keeper, contractCaller helper.IContractCaller, logger tmlog.Logger) sdk.Result {
	logger.Debug("Validating Checkpoint Data", "TxData", msg)
	if msg.TimeStamp == 0 || msg.TimeStamp > uint64(time.Now().Unix()) {
		logger.Error("Checkpoint timestamp must be in near past", "CurrentTime", time.Now().Unix(), "CheckpointTime", msg.TimeStamp, "Condition", msg.TimeStamp >= uint64(time.Now().Unix()))
		return common.ErrBadTimeStamp(k.Codespace).Result()
	}

	checkpointBuffer, err := k.GetCheckpointFromBuffer(ctx)
	if err == nil {
		if msg.TimeStamp == 0 || checkpointBuffer.TimeStamp == 0 || ((msg.TimeStamp > checkpointBuffer.TimeStamp) && msg.TimeStamp-checkpointBuffer.TimeStamp >= uint64(helper.GetConfig().CheckpointBufferTime.Seconds())) {
			logger.Debug("Checkpoint has been timed out, flushing buffer", "CheckpointTimestamp", msg.TimeStamp, "PrevCheckpointTimestamp", checkpointBuffer.TimeStamp)
			k.FlushCheckpointBuffer(ctx)
		} else {
			// calulates remaining time for buffer to be flushed
			checkpointTime := time.Unix(int64(checkpointBuffer.TimeStamp), 0)
			expiryTime := checkpointTime.Add(helper.GetConfig().CheckpointBufferTime)
			diff := expiryTime.Sub(time.Now()).Seconds()

			logger.Error("Checkpoint already exits in buffer", "Checkpoint", checkpointBuffer.String(), "Expires", expiryTime)

			return common.ErrNoACK(k.Codespace, diff).Result()
		}
	}
	logger.Debug("Received checkpoint from buffer", "Checkpoint", checkpointBuffer.String())

	// validate checkpoint
	if !ValidateCheckpoint(msg.StartBlock, msg.EndBlock, msg.RootHash, logger) {
		logger.Error("RootHash is not valid",
			"StartBlock", msg.StartBlock,
			"EndBlock", msg.EndBlock,
			"RootHash", msg.RootHash)
		return common.ErrBadBlockDetails(k.Codespace).Result()
	}
	logger.Debug("Valid Roothash in checkpoint", "StartBlock", msg.StartBlock, "EndBlock", msg.EndBlock)

	// fetch last checkpoint from store
	if lastCheckpoint, err := k.GetLastCheckpoint(ctx); err == nil {
		// make sure new checkpoint is after tip
		if lastCheckpoint.EndBlock > msg.StartBlock {
			logger.Error("Checkpoint already exists",
				"currentTip", lastCheckpoint.EndBlock,
				"startBlock", msg.StartBlock)
			return common.ErrOldCheckpoint(k.Codespace).Result()
		}
		if lastCheckpoint.EndBlock+1 != msg.StartBlock {
			logger.Error("Checkpoint not in countinuity",
				"currentTip", lastCheckpoint.EndBlock,
				"startBlock", msg.StartBlock)
			return common.ErrDisCountinuousCheckpoint(k.Codespace).Result()

		}
	} else if err.Error() == common.ErrNoCheckpointFound(k.Codespace).Error() && msg.StartBlock != 0 {
		logger.Error("First checkpoint to start from block 1", "Error", err)
		return common.ErrBadBlockDetails(k.Codespace).Result()
	}
	logger.Debug("Valid checkpoint tip")

	// check proposer in message
	if !bytes.Equal(msg.Proposer.Bytes(), k.GetValidatorSet(ctx).Proposer.Signer.Bytes()) {
		logger.Error("Invalid proposer in message",
			"currentProposer", k.GetValidatorSet(ctx).Proposer.Signer.String(),
			"checkpointProposer", msg.Proposer.String())
		return common.ErrBadProposerDetails(k.Codespace, k.GetValidatorSet(ctx).Proposer.Signer).Result()
	}
	logger.Debug("Valid proposer in checkpoint")

	// check if proposer has min ether
	balance, _ := contractCaller.GetBalance(msg.Proposer)
	if balance.Cmp(helper.MinBalance) == -1 {
		logger.Error("Proposer doesnt have enough ether to send checkpoint tx", "Balance", balance, "RequiredBalance", helper.MinBalance)
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

	checkpoint, _ := k.GetCheckpointFromBuffer(ctx)
	logger.Debug("Adding good checkpoint to buffer to await ACK", "checkpointStored", checkpoint.String())

	// indicate Checkpoint received by adding in cache, cache cleared in endblock
	k.SetCheckpointCache(ctx, common.DefaultValue)
	logger.Debug("Set Checkpoint Cache", "CheckpointReceived", k.GetCheckpointCache(ctx, common.CheckpointCacheKey))

	// send tags
	return sdk.Result{}
}

// Validates if checkpoint submitted on chain is valid
func HandleMsgCheckpointAck(ctx sdk.Context, msg MsgCheckpointAck, k common.Keeper, contractCaller helper.IContractCaller, logger tmlog.Logger) sdk.Result {
	logger.Debug("Validating Checkpoint ACK", "Tx", msg)

	// make call to headerBlock with header number
	root, start, end, createdAt, err := contractCaller.GetHeaderInfo(msg.HeaderBlock)
	if err != nil {
		logger.Error("Unable to fetch header from rootchain contract", "Error", err, "HeaderBlockIndex", msg.HeaderBlock)
		return common.ErrBadAck(k.Codespace).Result()
	}

	// check confirmation
	latestBlock, err := contractCaller.GetMainChainBlock(nil)
	if err != nil {
		logger.Error("Unable to connect to mainchain", "Error", err)
		return common.ErrNoConn(k.Codespace).Result()
	}
	if latestBlock.Number.Uint64()-createdAt < helper.GetConfig().ConfirmationBlocks {
		logger.Error("Not enough confirmations", "LatestBlock", latestBlock.Number.Uint64(), "TxBlock", createdAt)
		return common.ErrWaitFrConfirmation(k.Codespace).Result()
	}

	logger.Debug("HeaderBlock fetched", "headerBlock", msg.HeaderBlock, "start", start,
		"end", end, "Roothash", root, "CreatedAt", createdAt, "Latest", latestBlock.Number.Uint64())

	// get last checkpoint from buffer
	headerBlock, err := k.GetCheckpointFromBuffer(ctx)
	if err != nil {
		logger.Error("Unable to get checkpoint", "error", err)
		return common.ErrBadAck(k.Codespace).Result()
	}

	// match header block and checkpoint
	if start != headerBlock.StartBlock || end != headerBlock.EndBlock || !bytes.Equal(root.Bytes(), headerBlock.RootHash.Bytes()) {
		logger.Error("Invalid ACK",
			"startExpected", headerBlock.StartBlock,
			"startReceived", start,
			"endExpected", headerBlock.EndBlock,
			"endReceived", end,
			"rootExpected", headerBlock.RootHash.String(),
			"rootRecieved", root.String())

		return common.ErrBadAck(k.Codespace).Result()
	}

	// add checkpoint to headerBlocks
	k.AddCheckpoint(ctx, msg.HeaderBlock, headerBlock)
	logger.Info("Checkpoint added to store", "headerBlock", headerBlock.String())

	// flush buffer
	k.FlushCheckpointBuffer(ctx)
	logger.Debug("Checkpoint buffer flushed after receiving checkpoint ack", "checkpoint", headerBlock)

	// update ack count
	k.UpdateACKCount(ctx)
	logger.Debug("Valid ack received", "CurrentACKCount", k.GetACKCount(ctx)-1, "UpdatedACKCount", k.GetACKCount(ctx))

	// indicate ACK received by adding in cache, cache cleared in endblock
	k.SetCheckpointAckCache(ctx, common.DefaultValue)
	logger.Debug("Checkpoint ACK cache set", "CacheValue", k.GetCheckpointCache(ctx, common.CheckpointACKCacheKey))

	return sdk.Result{}
}

// Validate checkpoint no-ack transaction
func HandleMsgCheckpointNoAck(ctx sdk.Context, msg MsgCheckpointNoAck, k common.Keeper, logger tmlog.Logger) sdk.Result {
	logger.Debug("Validating checkpoint no-ack", "TxData", msg)
	// current time
	currentTime := time.Unix(int64(msg.TimeStamp), 0) // buffer time
	bufferTime := helper.GetConfig().CheckpointBufferTime

	// fetch last checkpoint from store
	// TODO figure out how to handle this error
	lastCheckpoint, _ := k.GetLastCheckpoint(ctx)
	lastCheckpointTime := time.Unix(int64(lastCheckpoint.TimeStamp), 0)

	// if last checkpoint is not present or last checkpoint happens before checkpoint buffer time -- thrown an error
	if lastCheckpointTime.After(currentTime) || (currentTime.Sub(lastCheckpointTime) < bufferTime) {
		logger.Debug("Invalid No ACK -- ongoing buffer period")
		return common.ErrInvalidNoACK(k.Codespace).Result()
	}

	// check last no ack - prevents repetitive no-ack
	lastAck := k.GetLastNoAck(ctx)
	lastAckTime := time.Unix(int64(lastAck), 0)

	if lastAckTime.After(currentTime) || (currentTime.Sub(lastAckTime) < bufferTime) {
		logger.Debug("Too many no-ack")
		return common.ErrTooManyNoACK(k.Codespace).Result()
	}

	// set last no ack
	k.SetLastNoAck(ctx, uint64(currentTime.Unix()))
	logger.Debug("Last No-ACK time set", "LastNoAck", k.GetLastNoAck(ctx))

	// --- Update to new proposer

	// increment accum
	k.IncreamentAccum(ctx, 1)

	//log new proposer
	vs := k.GetValidatorSet(ctx)
	newProposer := vs.GetProposer()
	logger.Debug(
		"New proposer selected",
		"validator", newProposer.Signer.String(),
		"signer", newProposer.Signer.String(),
		"power", newProposer.Power,
	)

	// --- End
	return sdk.Result{}
}
