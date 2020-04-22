package checkpoint

import (
	"bytes"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// NewHandler creates new handler for handling messages for checkpoint module
func NewHandler(k Keeper, contractCaller helper.IContractCaller) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgCheckpoint:
			return handleMsgCheckpoint(ctx, msg, k, contractCaller)
		case types.MsgCheckpointAck:
			return handleMsgCheckpointAck(ctx, msg, k, contractCaller)
		case types.MsgCheckpointNoAck:
			return handleMsgCheckpointNoAck(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("Invalid message in checkpoint module").Result()
		}
	}
}

// handleMsgCheckpoint Validates checkpoint transaction
func handleMsgCheckpoint(ctx sdk.Context, msg types.MsgCheckpoint, k Keeper, contractCaller helper.IContractCaller) sdk.Result {
	k.Logger(ctx).Debug("Validating checkpoint data", "TxData", msg)
	timeStamp := uint64(ctx.BlockTime().Unix())
	params := k.GetParams(ctx)

	checkpointBuffer, err := k.GetCheckpointFromBuffer(ctx)
	if err == nil {
		checkpointBufferTime := uint64(params.CheckpointBufferTime.Seconds())

		if checkpointBuffer.TimeStamp == 0 || ((timeStamp > checkpointBuffer.TimeStamp) && timeStamp-checkpointBuffer.TimeStamp >= checkpointBufferTime) {
			k.Logger(ctx).Debug("Checkpoint has been timed out, flushing buffer", "CheckpointTimestamp", timeStamp, "PrevCheckpointTimestamp", checkpointBuffer.TimeStamp)
			k.FlushCheckpointBuffer(ctx)
		} else {
			expiryTime := checkpointBuffer.TimeStamp + checkpointBufferTime
			k.Logger(ctx).Error("Checkpoint already exits in buffer", "Checkpoint", checkpointBuffer.String(), "Expires", expiryTime)
			return common.ErrNoACK(k.Codespace(), expiryTime).Result()
		}
	}

	// validate checkpoint
	validCheckpoint, err := types.ValidateCheckpoint(msg.StartBlock, msg.EndBlock, msg.RootHash, params.AvgCheckpointLength)
	if err != nil {
		k.Logger(ctx).Error("Error validating checkpoint",
			"Error", err,
			"StartBlock", msg.StartBlock,
			"EndBlock", msg.EndBlock)
		return common.ErrBadBlockDetails(k.Codespace()).Result()
	}

	if !validCheckpoint {
		k.Logger(ctx).Error("RootHash is not valid",
			"StartBlock", msg.StartBlock,
			"EndBlock", msg.EndBlock,
			"RootHash", msg.RootHash)
		return common.ErrBadBlockDetails(k.Codespace()).Result()
	}

	k.Logger(ctx).Debug("Valid Roothash in checkpoint", "StartBlock", msg.StartBlock, "EndBlock", msg.EndBlock)

	// fetch last checkpoint from store
	if lastCheckpoint, err := k.GetLastCheckpoint(ctx); err == nil {
		// make sure new checkpoint is after tip
		if lastCheckpoint.EndBlock > msg.StartBlock {
			k.Logger(ctx).Error("Checkpoint already exists",
				"currentTip", lastCheckpoint.EndBlock,
				"startBlock", msg.StartBlock)
			return common.ErrOldCheckpoint(k.Codespace()).Result()
		}
		if lastCheckpoint.EndBlock+1 != msg.StartBlock {
			k.Logger(ctx).Error("Checkpoint not in countinuity",
				"currentTip", lastCheckpoint.EndBlock,
				"startBlock", msg.StartBlock)
			return common.ErrDisCountinuousCheckpoint(k.Codespace()).Result()
		}
	} else if err.Error() == common.ErrNoCheckpointFound(k.Codespace()).Error() && msg.StartBlock != 0 {
		k.Logger(ctx).Error("First checkpoint to start from block 1", "Error", err)
		return common.ErrBadBlockDetails(k.Codespace()).Result()
	}
	k.Logger(ctx).Debug("Valid checkpoint tip")

	// make sure latest AccountRootHash matches
	// Calculate new account root hash
	dividendAccounts := k.sk.GetAllDividendAccounts(ctx)
	k.Logger(ctx).Debug("DividendAccounts of all validators", "dividendAccounts", dividendAccounts)
	accountRoot, err := types.GetAccountRootHash(dividendAccounts)
	k.Logger(ctx).Info("Validator Account root hash generated", "AccountRootHash", hmTypes.BytesToHeimdallHash(accountRoot).String())

	if !bytes.Equal(accountRoot, msg.AccountRootHash.Bytes()) {
		k.Logger(ctx).Error("AccountRootHash of current state", hmTypes.BytesToHeimdallHash(accountRoot).String(),
			"doesn't match with AccountRootHash of msg", msg.AccountRootHash)
		return common.ErrBadBlockDetails(k.Codespace()).Result()
	}

	k.Logger(ctx).Debug("AccountRootHash matches")

	// check proposer in message
	if !bytes.Equal(msg.Proposer.Bytes(), k.sk.GetValidatorSet(ctx).Proposer.Signer.Bytes()) {
		k.Logger(ctx).Error("Invalid proposer in message",
			"currentProposer", k.sk.GetValidatorSet(ctx).Proposer.Signer.String(),
			"checkpointProposer", msg.Proposer.String())
		return common.ErrBadProposerDetails(k.Codespace(), k.sk.GetValidatorSet(ctx).Proposer.Signer).Result()
	}
	k.Logger(ctx).Debug("Valid proposer in checkpoint")

	// add checkpoint to buffer
	// Add AccountRootHash to CheckpointBuffer
	k.SetCheckpointBuffer(ctx, hmTypes.CheckpointBlockHeader{
		StartBlock:      msg.StartBlock,
		EndBlock:        msg.EndBlock,
		RootHash:        msg.RootHash,
		AccountRootHash: msg.AccountRootHash,
		Proposer:        msg.Proposer,
		TimeStamp:       timeStamp,
	})

	checkpoint, _ := k.GetCheckpointFromBuffer(ctx)
	k.Logger(ctx).Debug("Adding good checkpoint to buffer to await ACK", "checkpointStored", checkpoint.String())

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCheckpoint,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyStartBlock, strconv.FormatUint(uint64(msg.StartBlock), 10)),
			sdk.NewAttribute(types.AttributeKeyEndBlock, strconv.FormatUint(uint64(msg.EndBlock), 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// handleMsgCheckpointAck Validates if checkpoint submitted on chain is valid
func handleMsgCheckpointAck(ctx sdk.Context, msg types.MsgCheckpointAck, k Keeper, contractCaller helper.IContractCaller) sdk.Result {
	k.Logger(ctx).Debug("Validating Checkpoint ACK", "Tx", msg)

	// make call to headerBlock with header number
	chainParams := k.ck.GetParams(ctx).ChainParams

	rootChainInstance, err := contractCaller.GetRootChainInstance(chainParams.RootChainAddress.EthAddress())
	if err != nil {
		k.Logger(ctx).Error("Unable to fetch rootchain contract instance", "Error", err)
		return common.ErrBadAck(k.Codespace()).Result()
	}
	root, start, end, createdAt, proposer, err := contractCaller.GetHeaderInfo(msg.HeaderBlock, rootChainInstance)
	if err != nil {
		k.Logger(ctx).Error("Unable to fetch header from rootchain contract", "Error", err, "headerBlockIndex", msg.HeaderBlock)
		return common.ErrBadAck(k.Codespace()).Result()
	}

	k.Logger(ctx).Debug("HeaderBlock fetched",
		"headerBlock", msg.HeaderBlock,
		"start", start,
		"end", end,
		"roothash", root,
		"proposer", proposer,
		"createdAt", createdAt,
	)

	// get last checkpoint from buffer
	headerBlock, err := k.GetCheckpointFromBuffer(ctx)
	if err != nil {
		k.Logger(ctx).Error("Unable to get checkpoint", "error", err)
		return common.ErrBadAck(k.Codespace()).Result()
	}
	if start != headerBlock.StartBlock {
		k.Logger(ctx).Error("Invalid start block", "startExpected", headerBlock.StartBlock, "startReceived", start)
		return common.ErrBadAck(k.Codespace()).Result()
	} else if start == headerBlock.StartBlock && end == headerBlock.EndBlock && !bytes.Equal(root.Bytes(), headerBlock.RootHash.Bytes()) {
		k.Logger(ctx).Error("Invalid ACK",
			"startExpected", headerBlock.StartBlock,
			"startReceived", start,
			"endExpected", headerBlock.EndBlock,
			"endReceived", end,
			"rootExpected", headerBlock.RootHash.String(),
			"rootRecieved", root.String())
		return common.ErrBadAck(k.Codespace()).Result()
	}
	if headerBlock.EndBlock > end {
		k.Logger(ctx).Info("Adjusting endBlock to one already submitted on chain", "OldEndBlock", headerBlock.EndBlock, "AdjustedEndBlock", end)
		headerBlock.EndBlock = end
		headerBlock.RootHash = hmTypes.HeimdallHash(root)
		// TODO proposer also needs to be changed
	}

	// Add checkpoint to headerBlocks
	k.AddCheckpoint(ctx, msg.HeaderBlock, *headerBlock)
	k.Logger(ctx).Info("Checkpoint added to store", "headerBlock", headerBlock.String())

	// flush buffer
	k.FlushCheckpointBuffer(ctx)
	k.Logger(ctx).Debug("Checkpoint buffer flushed after receiving checkpoint ack", "checkpoint", headerBlock)

	// update ack count
	k.UpdateACKCount(ctx)
	k.Logger(ctx).Debug("Valid ack received", "CurrentACKCount", k.GetACKCount(ctx)-1, "UpdatedACKCount", k.GetACKCount(ctx))

	// --- Update to new proposer

	// increment accum
	k.sk.IncrementAccum(ctx, 1)

	//log new proposer
	vs := k.sk.GetValidatorSet(ctx)
	newProposer := vs.GetProposer()
	k.Logger(ctx).Debug(
		"New proposer selected",
		"validator", newProposer.Signer.String(),
		"signer", newProposer.Signer.String(),
		"power", newProposer.VotingPower,
	)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCheckpointAck,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyHeaderIndex, strconv.FormatUint(uint64(msg.HeaderBlock), 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// Validate checkpoint no-ack transaction
func handleMsgCheckpointNoAck(ctx sdk.Context, msg types.MsgCheckpointNoAck, k Keeper) sdk.Result {
	k.Logger(ctx).Debug("Validating checkpoint no-ack", "TxData", msg)
	// current time
	currentTime := ctx.BlockTime()

	bufferTime := k.GetParams(ctx).CheckpointBufferTime

	// fetch last checkpoint from store
	// TODO figure out how to handle this error
	lastCheckpoint, _ := k.GetLastCheckpoint(ctx)
	lastCheckpointTime := time.Unix(int64(lastCheckpoint.TimeStamp), 0)

	// if last checkpoint is not present or last checkpoint happens before checkpoint buffer time -- thrown an error
	if lastCheckpointTime.After(currentTime) || (currentTime.Sub(lastCheckpointTime) < bufferTime) {
		k.Logger(ctx).Debug("Invalid No ACK -- Waiting for last checkpoint ACK")
		return common.ErrInvalidNoACK(k.Codespace()).Result()
	}

	// check last no ack - prevents repetitive no-ack
	lastAck := k.GetLastNoAck(ctx)
	lastAckTime := time.Unix(int64(lastAck), 0)

	if lastAckTime.After(currentTime) || (currentTime.Sub(lastAckTime) < bufferTime) {
		k.Logger(ctx).Debug("Too many no-ack")
		return common.ErrTooManyNoACK(k.Codespace()).Result()
	}

	// set last no ack
	k.SetLastNoAck(ctx, uint64(currentTime.Unix()))
	k.Logger(ctx).Debug("Last No-ACK time set", "LastNoAck", k.GetLastNoAck(ctx))

	// --- Update to new proposer

	// increment accum
	k.sk.IncrementAccum(ctx, 1)

	//log new proposer
	vs := k.sk.GetValidatorSet(ctx)
	newProposer := vs.GetProposer()
	k.Logger(ctx).Debug(
		"New proposer selected",
		"validator", newProposer.Signer.String(),
		"signer", newProposer.Signer.String(),
		"power", newProposer.VotingPower,
	)

	// add events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCheckpointNoAck,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyNewProposer, newProposer.Signer.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
