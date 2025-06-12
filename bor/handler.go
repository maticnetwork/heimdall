package bor

import (
	"errors"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/bor/types"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
)

// NewHandler returns a handler for "bor" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgProposeSpan,
			types.MsgProposeSpanV2:
			return HandleMsgProposeSpan(ctx, msg, k)
		case types.MsgBackfillSpans:
			return HandleMsgBackfillSpans(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("Invalid message in bor module").Result()
		}
	}
}

// HandleMsgProposeSpan handles proposeSpan msg
func HandleMsgProposeSpan(ctx sdk.Context, msg sdk.Msg, k Keeper) sdk.Result {
	var proposeMsg types.MsgProposeSpanV2
	switch msg := msg.(type) {
	case types.MsgProposeSpan:
		if ctx.BlockHeight() >= helper.GetDanelawHeight() {
			err := errors.New("msg span is not allowed after Danelaw hardfork height")
			k.Logger(ctx).Error(err.Error())
			return sdk.ErrTxDecode(err.Error()).Result()
		}
		proposeMsg = types.MsgProposeSpanV2{
			ID:         msg.ID,
			Proposer:   msg.Proposer,
			StartBlock: msg.StartBlock,
			EndBlock:   msg.EndBlock,
			ChainID:    msg.ChainID,
			Seed:       msg.Seed,
		}
	case types.MsgProposeSpanV2:
		if ctx.BlockHeight() < helper.GetDanelawHeight() {
			err := errors.New("msg span v2 is not allowed before Danelaw hardfork height")
			k.Logger(ctx).Error(err.Error())
			return sdk.ErrTxDecode(err.Error()).Result()
		}
		proposeMsg = msg
	}

	k.Logger(ctx).Debug("✅ Validating proposed span msg",
		"proposer", proposeMsg.Proposer.String(),
		"spanId", proposeMsg.ID,
		"startBlock", proposeMsg.StartBlock,
		"endBlock", proposeMsg.EndBlock,
		"seed", proposeMsg.Seed.String(),
	)

	// chainManager params
	params := k.chainKeeper.GetParams(ctx)
	chainParams := params.ChainParams

	// check chain id
	if chainParams.BorChainID != proposeMsg.ChainID {
		k.Logger(ctx).Error("Invalid Bor chain id", "msgChainID", proposeMsg.ChainID)
		return common.ErrInvalidBorChainID(k.Codespace()).Result()
	}

	// check if last span is up or if greater diff than threshold is found between validator set
	lastSpan, err := k.GetLastSpan(ctx)
	if err != nil {
		k.Logger(ctx).Error("Unable to fetch last span", "Error", err)
		return common.ErrSpanNotFound(k.Codespace()).Result()
	}

	// Validate span continuity
	if lastSpan.ID+1 != proposeMsg.ID || proposeMsg.StartBlock != lastSpan.EndBlock+1 || proposeMsg.EndBlock < proposeMsg.StartBlock {
		k.Logger(ctx).Error("Blocks not in continuity",
			"lastSpanId", lastSpan.ID,
			"spanId", proposeMsg.ID,
			"lastSpanStartBlock", lastSpan.StartBlock,
			"lastSpanEndBlock", lastSpan.EndBlock,
			"spanStartBlock", proposeMsg.StartBlock,
			"spanEndBlock", proposeMsg.EndBlock,
		)

		return common.ErrSpanNotInContinuity(k.Codespace()).Result()
	}

	// Validate Span duration
	spanDuration := k.GetParams(ctx).SpanDuration
	if spanDuration != (proposeMsg.EndBlock - proposeMsg.StartBlock + 1) {
		k.Logger(ctx).Error("Span duration of proposed span is wrong",
			"proposedSpanDuration", proposeMsg.EndBlock-proposeMsg.StartBlock+1,
			"paramsSpanDuration", spanDuration,
		)

		return common.ErrInvalidSpanDuration(k.Codespace()).Result()
	}

	// add events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeProposeSpan,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySpanID, strconv.FormatUint(proposeMsg.ID, 10)),
			sdk.NewAttribute(types.AttributeKeySpanStartBlock, strconv.FormatUint(proposeMsg.StartBlock, 10)),
			sdk.NewAttribute(types.AttributeKeySpanEndBlock, strconv.FormatUint(proposeMsg.EndBlock, 10)),
		),
	})

	// draft result with events
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func HandleMsgBackfillSpans(ctx sdk.Context, msg types.MsgBackfillSpans, k Keeper) sdk.Result {

	k.Logger(ctx).Debug("✅ validating proposed backfill spans msg",
		"proposer", msg.Proposer,
		"latestSpanId", msg.LatestSpanID,
		"latestBorSpanId", msg.LatestBorSpanID,
		"chainId", msg.ChainID,
	)

	params := k.chainKeeper.GetParams(ctx)
	chainParams := params.ChainParams

	if chainParams.BorChainID != msg.ChainID {
		k.Logger(ctx).Error("invalid bor chain id", "expected", chainParams.BorChainID, "got", msg.ChainID)
		return common.ErrInvalidBorChainID(k.Codespace()).Result()
	}

	latestSpan, err := k.GetLastSpan(ctx)
	if err != nil {
		k.Logger(ctx).Error("failed to get latest span", "error", err)
		return common.ErrUnableToGetLastSpan(k.Codespace()).Result()
	}

	if msg.LatestSpanID != latestSpan.ID && msg.LatestSpanID != latestSpan.ID-1 {
		k.Logger(ctx).Error("invalid latest span id", "expected",
			fmt.Sprintf("%d or %d", latestSpan.ID, latestSpan.ID-1), "got", msg.LatestSpanID)
		return common.ErrInvalidLastSpanID(k.Codespace(), msg.LatestSpanID).Result()
	}

	if msg.LatestBorSpanID <= msg.LatestSpanID {
		k.Logger(ctx).Error("invalid bor span id, expected greater than latest span id",
			"latestSpanId", latestSpan.ID,
			"latestBorSpanId", msg.LatestBorSpanID,
		)
		return common.ErrInvalidLastBorSpanID(k.Codespace(), msg.LatestSpanID).Result()
	}

	latestMilestone, err := k.checkpointKeeper.GetLastMilestone(ctx)
	if err != nil {
		k.Logger(ctx).Error("failed to get latest milestone", "error", err)
		return common.ErrUnableToGetLastMilestone(k.Codespace()).Result()
	}

	if latestMilestone == nil {
		k.Logger(ctx).Error("latest milestone is nil")
		return common.ErrLatestMilestoneNotFound(k.Codespace()).Result()
	}

	borSpanId, err := types.CalcCurrentBorSpanId(latestMilestone.EndBlock, latestSpan)
	if err != nil {
		k.Logger(ctx).Error("failed to calculate bor span id", "error", err)
		return common.ErrUnableToCalculateBorSpanID(k.Codespace()).Result()
	}

	if borSpanId != msg.LatestBorSpanID {
		k.Logger(ctx).Error(
			"bor span id mismatch",
			"calculatedBorSpanId", borSpanId,
			"msgLatestBorSpanId", msg.LatestBorSpanID,
			"latestMilestoneEndBlock", latestMilestone.EndBlock,
			"latestSpanStartBlock", latestSpan.StartBlock,
			"latestSpanEndBlock", latestSpan.EndBlock,
			"latestSpanId", latestSpan.ID,
		)
		return common.ErrBorSpanIDMismatch(k.Codespace(), borSpanId, msg.LatestBorSpanID).Result()
	}

	return sdk.Result{}
}
