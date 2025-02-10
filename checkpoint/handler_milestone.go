package checkpoint

import (
	"bytes"
	"encoding/base64"
	"io"
	"net/http"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	jsoniter "github.com/json-iterator/go"

	"github.com/maticnetwork/heimdall/bridge/setu/util"
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

	// if ctx.BlockHeight()-k.GetMilestoneBlockNumber(ctx) < 2 {
	// 	logger.Error(
	// 		"Previous milestone still in voting phase",
	// 		"previousMilestoneBlock", k.GetMilestoneBlockNumber(ctx),
	// 		"currentMilestoneBlock", ctx.BlockHeight(),
	// 	)

	// 	return common.ErrPrevMilestoneInVoting(k.Codespace()).Result()
	// }

	//Increment the priority in the milestone validator set
	k.sk.MilestoneIncrementAccum(ctx, 1)

	//
	//Check for the msg milestone
	//

	//Calculate the milestone length
	msgMilestoneLength := int64(msg.EndBlock) - int64(msg.StartBlock) + 1

	//check for the minimum length of milestone
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
		// Get highest end block from pending milestones
		highestPendingEndBlock, err := k.GetHighestPendingMilestoneEndBlock(ctx)
		logger.Info("Highest pending end block", "highestPendingEndBlock", highestPendingEndBlock, "error", err)

		// Check continuity with any pending milestone or last stored milestone
		isValid := false

		// Check against last verified milestone
		if lastMilestone.EndBlock+1 == msg.StartBlock {
			isValid = true
			logger.Info("Milestone in continuity with last stored milestone",
				"lastMilestoneEndBlock", lastMilestone.EndBlock,
				"receivedMsgStartBlock", msg.StartBlock,
			)
		}

		// Check against pending milestones in DB
		if highestPendingEndBlock+1 == msg.StartBlock {
			isValid = true
			logger.Info("Milestone in continuity with pending milestone in DB",
				"pendingEndBlock", highestPendingEndBlock,
				"receivedMsgStartBlock", msg.StartBlock,
			)
		}

		// Get all pending milestone end blocks from DB
		pendingMilestones, err := k.GetAllPendingMilestones(ctx)
		// Check against mempool transactions if available
		if err == nil {
			for _, pending := range pendingMilestones {
				if pending.EndBlock+1 == msg.StartBlock {
					isValid = true
					logger.Info("Milestone in continuity with pending milestone in mempool",
						"pendingEndBlock", pending.EndBlock,
						"receivedMsgStartBlock", msg.StartBlock,
					)
					break
				}
			}
		}

		if !isValid {
			logger.Error("Milestone not in continuity with any pending or stored milestone",
				"lastMilestoneEndBlock", lastMilestone.EndBlock,
				"receivedMsgStartBlock", msg.StartBlock,
			)
			return common.ErrMilestoneNotInContinuity(k.Codespace()).Result()
		}
	} else if msg.StartBlock != helper.GetMilestoneBorBlockHeight() {
		logger.Error("First milestone to start from", "block", helper.GetMilestoneBorBlockHeight(), "milestone start block", msg.StartBlock, "error", err)
		return common.ErrNoMilestoneFound(k.Codespace()).Result()
	}

	k.SetMilestoneBlockNumber(ctx, ctx.BlockHeight(), msg.EndBlock)
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

// handleMsgMilestoneTimeout handles milestone timeout transaction
func handleMsgMilestoneTimeout(ctx sdk.Context, msg types.MsgMilestoneTimeout, k Keeper) sdk.Result {
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

	// Remove milestone from voting phase if it exists
	// if lastMilestone.MilestoneID != "" {
	// 	k.RemoveMilestoneFromVoting(ctx, lastMilestone.MilestoneID)
	// }

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

// GetAllPendingMilestones returns all pending milestones from the store
func (k Keeper) GetAllPendingMilestones(ctx sdk.Context) ([]types.MsgMilestone, error) {
	endpoint := helper.GetConfig().TendermintRPCUrl + util.TendermintUnconfirmedTxsURL

	var pendingMilestones []types.MsgMilestone

	resp, err := helper.Client.Get(endpoint)
	if err != nil || resp.StatusCode != http.StatusOK {
		return pendingMilestones, err
	}
	defer resp.Body.Close()

	// Limit the number of bytes read from the response body
	limitedBody := http.MaxBytesReader(nil, resp.Body, helper.APIBodyLimit)

	body, err := io.ReadAll(limitedBody)
	if err != nil {
		return pendingMilestones, err
	}

	// a minimal response of the unconfirmed txs
	var response util.TendermintUnconfirmedTxs

	err = jsoniter.ConfigFastest.Unmarshal(body, &response)
	if err != nil {
		return pendingMilestones, err
	}

	for _, txn := range response.Result.Txs {
		// Tendermint encodes the transactions with base64 encoding. Decode it first.
		txBytes, err := base64.StdEncoding.DecodeString(txn)
		if err != nil {
			return pendingMilestones, err
		}

		// Unmarshal the transaction from bytes
		decodedTx, err := helper.GetTxDecoder(k.cdc)(txBytes)
		if err != nil {
			continue
		}

		// We only need to check for milestone type transactions
		if decodedTx.GetMsgs()[0].Type() == "milestone" {
			pendingMilestones = append(pendingMilestones, decodedTx.GetMsgs()[0].(types.MsgMilestone))
		}
	}

	return pendingMilestones, nil
}
