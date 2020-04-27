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
	logger := k.Logger(ctx)

	timeStamp := uint64(ctx.BlockTime().Unix())
	params := k.GetParams(ctx)

	//
	// Check checkpoint buffer
	//

	checkpointBuffer, err := k.GetCheckpointFromBuffer(ctx)
	if err == nil {
		checkpointBufferTime := uint64(params.CheckpointBufferTime.Seconds())

		if checkpointBuffer.TimeStamp == 0 || ((timeStamp > checkpointBuffer.TimeStamp) && timeStamp-checkpointBuffer.TimeStamp >= checkpointBufferTime) {
			logger.Debug("Checkpoint has been timed out. Flushing buffer.", "checkpointTimestamp", timeStamp, "prevCheckpointTimestamp", checkpointBuffer.TimeStamp)
			k.FlushCheckpointBuffer(ctx)
		} else {
			expiryTime := checkpointBuffer.TimeStamp + checkpointBufferTime
			logger.Error("Checkpoint already exits in buffer", "Checkpoint", checkpointBuffer.String(), "Expires", expiryTime)
			return common.ErrNoACK(k.Codespace(), expiryTime).Result()
		}
	}

	//
	// Validate last checkpoint
	//

	// fetch last checkpoint from store
	if lastCheckpoint, err := k.GetLastCheckpoint(ctx); err == nil {
		// make sure new checkpoint is after tip
		if lastCheckpoint.EndBlock > msg.StartBlock {
			logger.Error("Checkpoint already exists",
				"currentTip", lastCheckpoint.EndBlock,
				"startBlock", msg.StartBlock,
			)
			return common.ErrOldCheckpoint(k.Codespace()).Result()
		}

		// check if new checkpoint's start block start from current tip
		if lastCheckpoint.EndBlock+1 != msg.StartBlock {
			logger.Error("Checkpoint not in countinuity",
				"currentTip", lastCheckpoint.EndBlock,
				"startBlock", msg.StartBlock)
			return common.ErrDisCountinuousCheckpoint(k.Codespace()).Result()
		}
	} else if err.Error() == common.ErrNoCheckpointFound(k.Codespace()).Error() && msg.StartBlock != 0 {
		logger.Error("First checkpoint to start from block 0", "Error", err)
		return common.ErrBadBlockDetails(k.Codespace()).Result()
	}

	//
	// Validate account hash
	//

	// Make sure latest AccountRootHash matches
	// Calculate new account root hash
	dividendAccounts := k.sk.GetAllDividendAccounts(ctx)
	logger.Debug("DividendAccounts of all validators", "dividendAccountsLength", len(dividendAccounts))

	// Get account root has from dividend accounts
	accountRoot, err := types.GetAccountRootHash(dividendAccounts)
	if err != nil {
		logger.Error("Error while fetching account root hash", "error", err)
		return common.ErrBadBlockDetails(k.Codespace()).Result()
	}
	logger.Debug("Validator account root hash generated", "accountRootHash", hmTypes.BytesToHeimdallHash(accountRoot).String())

	// Compare stored root hash to msg root hash
	if !bytes.Equal(accountRoot, msg.AccountRootHash.Bytes()) {
		logger.Error(
			"AccountRootHash of current state doesn't match from msg",
			"hash", hmTypes.BytesToHeimdallHash(accountRoot).String(),
			"msgHash", msg.AccountRootHash,
		)
		return common.ErrBadBlockDetails(k.Codespace()).Result()
	}

	//
	// Validate proposer
	//

	// Check proposer in message
	validatorSet := k.sk.GetValidatorSet(ctx)
	if validatorSet.Proposer == nil {
		logger.Error("No proposer in validator set", "msgProposer", msg.Proposer.String())
		return common.ErrInvalidMsg(k.Codespace(), "No proposer in stored validator set").Result()
	}

	if !bytes.Equal(msg.Proposer.Bytes(), validatorSet.Proposer.Signer.Bytes()) {
		logger.Error(
			"Invalid proposer in msg",
			"proposer", validatorSet.Proposer.Signer.String(),
			"msgProposer", msg.Proposer.String(),
		)
		return common.ErrInvalidMsg(k.Codespace(), "Invalid proposer in msg").Result()
	}

	// Emit event for checkpoint
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCheckpoint,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyStartBlock, strconv.FormatUint(msg.StartBlock, 10)),
			sdk.NewAttribute(types.AttributeKeyEndBlock, strconv.FormatUint(msg.EndBlock, 10)),
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
