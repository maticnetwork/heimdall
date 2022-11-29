package checkpoint

import (
	"bytes"
	"fmt"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
)

// handleMsgMilestone Validates milestone transaction
func handleMsgMilestone(ctx sdk.Context, msg types.MsgMilestone, k Keeper) sdk.Result {

	logger := k.Logger(ctx)
	milestoneLength := helper.MilestoneLength

	//Check for the hard fork value
	if ctx.BlockHeight() < helper.GetMilestoneHardForkHeight() {
		logger.Error("Network hasn't reached the", "Hard forked height", helper.GetMilestoneHardForkHeight())
		return common.ErrInvalidMsg(k.Codespace(), "Network hasn't reached the milestone hard forked height").Result()
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

	//Increment the priority in the milestone validator set
	k.sk.MilestoneIncrementAccum(ctx, 1)

	if !bytes.Equal(msg.Proposer.Bytes(), validatorSet.Proposer.Signer.Bytes()) {
		logger.Error(
			"Invalid proposer in msg",
			"proposer", validatorSet.Proposer.Signer.String(),
			"msgProposer", msg.Proposer.String(),
		)

		return common.ErrInvalidMsg(k.Codespace(), "Invalid proposer in msg").Result()
	}

	//
	//Check for the msg milestone
	//

	if msg.StartBlock+milestoneLength-1 != msg.EndBlock {
		logger.Error("Milestone's length doesn't match the  milestone length set in configuration",
			"StartBlock", msg.StartBlock,
			"EndBlock", msg.EndBlock,
			"Milestone Length", milestoneLength,
		)

		return common.ErrMilestoneInvalid(k.Codespace()).Result()
	}

	// fetch last milestone from store
	if lastMilestone, err := k.GetLastMilestone(ctx); err == nil {
		// make sure new milestone is in continuity
		if lastMilestone.EndBlock+1 != msg.StartBlock {
			logger.Error("Milestone not in continuity ",
				"lastMilestoneEndBlock", lastMilestone.EndBlock,
				"receivedMsgStartBlock", msg.StartBlock,
			)

			return common.ErrMilestoneNotInContinuity(k.Codespace()).Result()
		}

	} else if err != nil && msg.StartBlock != 0 {
		logger.Error("First milestone to start from block 0", "milestone start block", msg.StartBlock, "error", err)
		return common.ErrNoMilestoneFound(k.Codespace()).Result()

	}

	// Emit event for milestone
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMilestone,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyStartBlock, strconv.FormatUint(msg.StartBlock, 10)),
			sdk.NewAttribute(types.AttributeKeyEndBlock, strconv.FormatUint(msg.EndBlock, 10)),
			sdk.NewAttribute(types.AttributeKeyRootHash, msg.RootHash.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// Handles milestone timeout transaction
func handleMsgMilestoneTimeout(ctx sdk.Context, msg types.MsgMilestoneTimeout, k Keeper) sdk.Result {
	logger := k.Logger(ctx)

	// Get current block time
	currentTime := ctx.BlockTime()
	fmt.Print("sudesh", currentTime.Second())

	// Get buffer time from params
	bufferTime := helper.MilestoneBufferTime

	// Fetch last checkpoint from store
	// TODO figure out how to handle this error
	lastMilestone, err := k.GetLastMilestone(ctx)

	if err != nil {
		logger.Error("Didn't find the last milestone")
		return common.ErrNoMilestoneFound(k.Codespace()).Result()
	}

	lastMilestoneTime := time.Unix(int64(lastMilestone.TimeStamp), 0)

	// If last milestone happens before milestone buffer time -- thrown an error
	if lastMilestoneTime.After(currentTime) || (currentTime.Sub(lastMilestoneTime) < bufferTime) {

		fmt.Print("last Milestone Time", lastMilestoneTime.Second())
		fmt.Print("Current Time", currentTime.Second())

		logger.Debug("Invalid Milestone Timeout msg", "lastMilestoneTime", lastMilestoneTime, "current time", currentTime,
			"buffer Time", bufferTime.String(),
		)

		logger.Error("Invalid Milestone Timeout msg", "lastMilestoneTime", lastMilestoneTime, "current time", currentTime,
			"buffer Time", bufferTime.String(),
		)

		return common.ErrInvalidMilestoneTimeout(k.Codespace()).Result()
	}

	// Check last no ack - prevents repetitive no-ack
	lastMilestoneTimeout := k.GetLastMilestoneTimeout(ctx)
	lastMilestoneTimeoutTime := time.Unix(int64(lastMilestoneTimeout), 0)

	if lastMilestoneTimeoutTime.After(currentTime) || (currentTime.Sub(lastMilestoneTimeoutTime) < bufferTime) {
		logger.Debug("Too many milestone timeout messages", "lastMilestoneTimeoutTime", lastMilestoneTimeoutTime, "current time", currentTime,
			"buffer Time", bufferTime.String())

		return common.ErrTooManyNoACK(k.Codespace()).Result()
	}

	// Set new last milestone-timeout
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
