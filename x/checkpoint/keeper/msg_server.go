package keeper

import (
	"bytes"
	"context"
	"strconv"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/x/checkpoint/types"
)

type msgServer struct {
	Keeper
	contractCaller helper.IContractCaller
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper, contractCaller helper.IContractCaller) types.MsgServer {
	return &msgServer{Keeper: keeper, contractCaller: contractCaller}
}

func (k msgServer) Checkpoint(goCtx context.Context, msg *types.MsgCheckpoint) (*types.MsgCheckpointResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
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
			logger.Error("Checkpoint already exists in buffer", "Checkpoint", checkpointBuffer.String(), "Expires", expiryTime)
			return nil, types.ErrNoACK
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
			return nil, types.ErrOldCheckpoint
		}

		// check if new checkpoint's start block start from current tip
		if lastCheckpoint.EndBlock+1 != msg.StartBlock {
			logger.Error("Checkpoint not in continuity",
				"currentTip", lastCheckpoint.EndBlock,
				"startBlock", msg.StartBlock)
			return nil, types.ErrDisCountinuousCheckpoint
		}
	} else if err.Error() == types.ErrNoCheckpointFound.Error() && msg.StartBlock != 0 {
		logger.Error("First checkpoint to start from block 0", "Error", err)
		return nil, types.ErrBadBlockDetails
	}

	//
	// Validate account hash
	//

	// Make sure latest AccountRootHash matches
	// Calculate new account root hash
	dividendAccounts := k.moduleCommunicator.GetAllDividendAccounts(ctx)
	logger.Debug("DividendAccounts of all validators", "dividendAccountsLength", len(dividendAccounts))

	// Get account root has from dividend accounts
	accountRoot, err := types.GetAccountRootHash(dividendAccounts)
	if err != nil {
		logger.Error("Error while fetching account root hash", "error", err)
		return nil, types.ErrBadBlockDetails
	}
	logger.Debug("Validator account root hash generated", "accountRootHash", hmTypes.BytesToHeimdallHash(accountRoot).String())

	// Compare stored root hash to msg root hash
	if hmTypes.BytesToHeimdallHash(accountRoot).String() != msg.AccountRootHash {
		//if !bytes.Equal(accountRoot, []byte(msg.AccountRootHash)) {
		logger.Error(
			"AccountRootHash of current state doesn't match from msg",
			"hash", hmTypes.BytesToHeimdallHash(accountRoot).String(),
			"msgHash", msg.AccountRootHash,
		)
		return nil, types.ErrBadBlockDetails
	}

	//
	// Validate proposer
	//

	// Check proposer in message
	validatorSet := k.Sk.GetValidatorSet(ctx)
	if validatorSet.Proposer == nil {
		logger.Error("No proposer in validator set", "msgProposer", msg.Proposer)
		return nil, types.ErrInvalidMsg
	}

	if !bytes.Equal([]byte(strings.ToLower(msg.Proposer)), []byte(strings.ToLower(validatorSet.Proposer.Signer))) {
		logger.Error(
			"Invalid proposer in msg",
			"proposer", validatorSet.Proposer.Signer,
			"msgProposer", msg.Proposer,
		)
		return nil, types.ErrInvalidMsg
	}

	// Emit event for checkpoint
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCheckpoint,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer),
			sdk.NewAttribute(types.AttributeKeyStartBlock, strconv.FormatUint(msg.StartBlock, 10)),
			sdk.NewAttribute(types.AttributeKeyEndBlock, strconv.FormatUint(msg.EndBlock, 10)),
			sdk.NewAttribute(types.AttributeKeyRootHash, msg.RootHash),
			sdk.NewAttribute(types.AttributeKeyAccountHash, msg.AccountRootHash),
		),
	})

	return &types.MsgCheckpointResponse{}, nil
}

func (k msgServer) CheckpointAck(goCtx context.Context, msg *types.MsgCheckpointAck) (*types.MsgCheckpointAckResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := k.Logger(ctx)

	// Get last checkpoint from buffer
	headerBlock, err := k.GetCheckpointFromBuffer(ctx)
	if err != nil {
		logger.Error("Unable to get checkpoint", "error", err)
		return nil, types.ErrBadAck
	}

	if msg.StartBlock != headerBlock.StartBlock {
		logger.Error("Invalid start block", "startExpected", headerBlock.StartBlock, "startReceived", msg.StartBlock)
		return nil, types.ErrBadAck
	}

	// Return err if start and end matches but contract root hash doesn't match
	if msg.StartBlock == headerBlock.StartBlock && msg.EndBlock == headerBlock.EndBlock && msg.RootHash != headerBlock.RootHash {
		logger.Error("Invalid ACK",
			"startExpected", headerBlock.StartBlock,
			"startReceived", msg.StartBlock,
			"endExpected", headerBlock.EndBlock,
			"endReceived", msg.EndBlock,
			"rootExpected", headerBlock.RootHash,
			"rootReceived", msg.RootHash,
		)
		return nil, types.ErrBadAck
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCheckpointAck,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyHeaderIndex, strconv.FormatUint(msg.Number, 10)),
		),
	})
	return &types.MsgCheckpointAckResponse{}, nil
}

func (k msgServer) CheckpointNoAck(goCtx context.Context, msg *types.MsgCheckpointNoAck) (*types.MsgCheckpointNoAckResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	logger := k.Logger(ctx)

	// Get current block time
	currentTime := ctx.BlockTime()

	// Get buffer time from params
	bufferTime := k.GetParams(ctx).CheckpointBufferTime

	// Fetch last checkpoint from store
	// TODO figure out how to handle this error
	lastCheckpoint, _ := k.GetLastCheckpoint(ctx)
	lastCheckpointTime := time.Unix(int64(lastCheckpoint.TimeStamp), 0)

	// If last checkpoint is not present or last checkpoint happens before checkpoint buffer time -- thrown an error
	if lastCheckpointTime.After(currentTime) || (currentTime.Sub(lastCheckpointTime) < bufferTime) {
		logger.Debug("Invalid No ACK -- Waiting for last checkpoint ACK")
		return nil, types.ErrInvalidNoACK
	}

	// Check last no ack - prevents repetitive no-ack
	lastNoAck := k.GetLastNoAck(ctx)
	lastNoAckTime := time.Unix(int64(lastNoAck), 0)

	if lastNoAckTime.After(currentTime) || (currentTime.Sub(lastNoAckTime) < bufferTime) {
		logger.Debug("Too many no-ack")
		return nil, types.ErrTooManyNoACK
	}

	// Set new last no-ack
	newLastNoAck := uint64(currentTime.Unix())
	k.SetLastNoAck(ctx, newLastNoAck)
	logger.Debug("Last No-ACK time set", "lastNoAck", newLastNoAck)

	//
	// Update to new proposer
	//

	// Increment accum (selects new proposer)
	k.Sk.IncrementAccum(ctx, 1)

	// Get new proposer
	vs := k.Sk.GetValidatorSet(ctx)
	newProposer := vs.GetProposer()
	logger.Debug(
		"New proposer selected",
		"validator", newProposer.Signer,
		"signer", newProposer.Signer,
		"power", newProposer.VotingPower,
	)

	// add events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCheckpointNoAck,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyNewProposer, newProposer.Signer),
		),
	})

	return &types.MsgCheckpointNoAckResponse{}, nil
}
