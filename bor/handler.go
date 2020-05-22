package bor

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/bor/types"
	"github.com/maticnetwork/heimdall/common"
)

// NewHandler returns a handler for "bor" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgProposeSpan:
			return HandleMsgProposeSpan(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("Invalid message in bor module").Result()
		}
	}
}

// HandleMsgProposeSpan handles proposeSpan msg
func HandleMsgProposeSpan(ctx sdk.Context, msg types.MsgProposeSpan, k Keeper) sdk.Result {
	k.Logger(ctx).Debug("✅ Validating proposed span msg",
		"spanId", msg.ID,
		"startBlock", msg.StartBlock,
		"endBlock", msg.EndBlock,
		"seed", msg.Seed.String(),
	)

	// chainManager params
	params := k.chainKeeper.GetParams(ctx)
	chainParams := params.ChainParams

	// check chain id
	if chainParams.BorChainID != msg.ChainID {
		k.Logger(ctx).Error("Invalid Bor chain id", "msgChainID", msg.ChainID)
		return common.ErrInvalidBorChainID(k.Codespace()).Result()
	}

	// check if last span is up or if greater diff than threshold is found between validator set
	lastSpan, err := k.GetLastSpan(ctx)
	if err != nil {
		k.Logger(ctx).Error("Unable to fetch last span", "Error", err)
		return common.ErrSpanNotFound(k.Codespace()).Result()
	}

	// Validate span continuity
	if lastSpan.ID+1 != msg.ID || msg.StartBlock != lastSpan.EndBlock+1 || msg.EndBlock < msg.StartBlock {
		k.Logger(ctx).Error("Blocks not in countinuity",
			"lastSpanId", lastSpan.ID,
			"spanId", msg.ID,
			"lastSpanStartBlock", lastSpan.StartBlock,
			"lastSpanEndBlock", lastSpan.EndBlock,
			"spanStartBlock", msg.StartBlock,
			"spanEndBlock", msg.EndBlock,
		)
		return common.ErrSpanNotInCountinuity(k.Codespace()).Result()
	}

	// Validate Span duration
	spanDuration := k.GetParams(ctx).SpanDuration
	if spanDuration != (msg.EndBlock - msg.StartBlock + 1) {
		k.Logger(ctx).Error("Span duration of proposed span is wrong",
			"proposedSpanDuration", (msg.EndBlock - msg.StartBlock + 1),
			"paramsSpanDuration", spanDuration,
		)
		return common.ErrInvalidSpanDuration(k.Codespace()).Result()
	}

	// add events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeProposeSpan,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySpanID, strconv.FormatUint(msg.ID, 10)),
			sdk.NewAttribute(types.AttributeKeySpanStartBlock, strconv.FormatUint(msg.StartBlock, 10)),
			sdk.NewAttribute(types.AttributeKeySpanEndBlock, strconv.FormatUint(msg.EndBlock, 10)),
		),
	})

	// draft result with events
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
