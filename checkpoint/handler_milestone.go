package checkpoint

import (
	"bytes"
	"math"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
)

// handleMsgMilestone validates milestone transaction
func handleMsgMilestone(ctx sdk.Context, msg types.MsgMilestone, k Keeper) sdk.Result {
	logger := k.MilestoneLogger(ctx)
	milestoneLength := helper.MilestoneLength

	//
	//Get milestone validator set
	//

	//Get the milestone proposer
	validatorSet := k.sk.GetMilestoneValidatorSet(ctx)
	if validatorSet.Proposer == nil {
		logger.Error("No proposer in validator set", "msgProposer", msg.Proposer.String())
		return common.ErrInvalidMsg(k.Codespace(), "No proposer in stored validator set").Result()
	}

	//
	// Validate proposer
	//

	//check for the milestone proposer
	if !bytes.Equal(msg.Proposer.Bytes(), validatorSet.Proposer.Signer.Bytes()) {
		logger.Error(
			"Invalid proposer in msg",
			"proposer", validatorSet.Proposer.Signer.String(),
			"msgProposer", msg.Proposer.String(),
		)

		return common.ErrInvalidMsg(k.Codespace(), "Invalid proposer in msg").Result()
	}

	if ctx.BlockHeight()-k.GetMilestoneBlockNumber(ctx) < 2 {
		logger.Error(
			"Previous milestone still in voting phase",
			"previousMilestoneBlock", k.GetMilestoneBlockNumber(ctx),
			"currentMilestoneBlock", ctx.BlockHeight(),
		)

		return common.ErrPrevMilestoneInVoting(k.Codespace()).Result()
	}

	//Increment the priority in the milestone validator set
	k.sk.MilestoneIncrementAccum(ctx, 1)

	// Calculate the milestone length
	if msg.EndBlock > math.MaxInt64 || msg.StartBlock > math.MaxInt64 {
		return common.ErrMilestoneInvalid("Block number exceeds int64 max value").Result()
	}
	msgMilestoneLength := int64(msg.EndBlock) - int64(msg.StartBlock) + 1

	// Check for the minimum length of milestone
	if msgMilestoneLength < int64(milestoneLength) {
		logger.Error("Length of the milestone should be greater than configured minimum milestone length",
			"StartBlock", msg.StartBlock,
			"EndBlock", msg.EndBlock,
			"Minimum Milestone Length", milestoneLength,
		)

		return common.ErrMilestoneInvalid(k.Codespace()).Result()
	}

	// fetch last stored milestone from store
	if lastMilestone, err := k.GetLastMilestone(ctx); err == nil {
		// make sure new milestone is in continuity
		if lastMilestone.EndBlock+1 != msg.StartBlock {
			logger.Error("Milestone not in continuity ",
				"lastMilestoneEndBlock", lastMilestone.EndBlock,
				"receivedMsgStartBlock", msg.StartBlock,
			)

			return common.ErrMilestoneNotInContinuity(k.Codespace()).Result()
		}
	} else if msg.StartBlock != helper.GetMilestoneBorBlockHeight() {
		logger.Error("First milestone to start from", "block", helper.GetMilestoneBorBlockHeight(), "milestone start block", msg.StartBlock, "error", err)
		return common.ErrNoMilestoneFound(k.Codespace()).Result()
	}

	k.SetMilestoneBlockNumber(ctx, ctx.BlockHeight())

	//Set the MilestoneID in the cache
	types.SetMilestoneID(msg.MilestoneID)

	// Emit event for milestone
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMilestone,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyStartBlock, strconv.FormatUint(msg.StartBlock, 10)),
			sdk.NewAttribute(types.AttributeKeyEndBlock, strconv.FormatUint(msg.EndBlock, 10)),
			sdk.NewAttribute(types.AttributeKeyHash, msg.Hash.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// Handles milestone timeout transaction
func handleMsgMilestoneTimeout(ctx sdk.Context, _ types.MsgMilestoneTimeout, k Keeper) sdk.Result {
	logger := k.MilestoneLogger(ctx)

	// Get current block time
	currentTime := ctx.BlockTime()

	// Get buffer time from params
	bufferTime := helper.MilestoneBufferTime

	// Fetch last checkpoint from store
	// TODO figure out how to handle this error
	lastMilestone, err := k.GetLastMilestone(ctx)

	if err != nil {
		logger.Error("Didn't find the last milestone", "err", err)
		return common.ErrNoMilestoneFound(k.Codespace()).Result()
	}

	if lastMilestone.TimeStamp > math.MaxInt64 {
		return common.ErrMilestoneInvalid("Milestone timestamp exceeds int64 max value").Result()
	}
	lastMilestoneTime := time.Unix(int64(lastMilestone.TimeStamp), 0)

	// If last milestone happens before milestone buffer time -- thrown an error
	if lastMilestoneTime.After(currentTime) || (currentTime.Sub(lastMilestoneTime) < bufferTime) {
		logger.Error("Invalid Milestone Timeout msg", "lastMilestoneTime", lastMilestoneTime, "current time", currentTime,
			"buffer Time", bufferTime.String(),
		)

		return common.ErrInvalidMilestoneTimeout(k.Codespace()).Result()
	}

	// Check last no ack - prevents repetitive no-ack
	lastMilestoneTimeout := k.GetLastMilestoneTimeout(ctx)
	if lastMilestoneTimeout > math.MaxInt64 {
		return common.ErrMilestoneInvalid("Milestone timeout exceeds int64 max value").Result()
	}
	lastMilestoneTimeoutTime := time.Unix(int64(lastMilestoneTimeout), 0)

	if lastMilestoneTimeoutTime.After(currentTime) || (currentTime.Sub(lastMilestoneTimeoutTime) < bufferTime) {
		logger.Debug("Too many milestone timeout messages", "lastMilestoneTimeoutTime", lastMilestoneTimeoutTime, "current time", currentTime,
			"buffer Time", bufferTime.String())

		return common.ErrTooManyNoACK(k.Codespace()).Result()
	}

	// Set new last milestone-timeout
	//nolint:gosec
	newLastMilestoneTimeout := uint64(currentTime.Unix())
	k.SetLastMilestoneTimeout(ctx, newLastMilestoneTimeout)
	logger.Debug("Last milestone-timeout set", "lastMilestoneTimeout", newLastMilestoneTimeout)

	//
	// Update to new proposer
	//

	// Increment accum (selects new proposer)
	k.sk.MilestoneIncrementAccum(ctx, 1)

	// Get new proposer
	vs := k.sk.GetMilestoneValidatorSet(ctx)
	newProposer := vs.GetProposer()
	logger.Debug(
		"New milestone proposer selected",
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
