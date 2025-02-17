package bor

import (
	"errors"
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
