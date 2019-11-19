package bor

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/bor/tags"
	"github.com/maticnetwork/heimdall/common"
)

// NewHandler returns a handler for "bor" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgProposeSpan:
			return HandleMsgProposeSpan(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("Invalid message in bor module").Result()
		}
	}
}

// HandleMsgProposeSpan handles proposeSpan msg
func HandleMsgProposeSpan(ctx sdk.Context, msg MsgProposeSpan, k Keeper) sdk.Result {
	k.Logger(ctx).Debug("Proposing span", "TxData", msg)

	// check if last span is up or if greater diff than threshold is found between validator set
	lastSpan, err := k.GetLastSpan(ctx)
	if err != nil {
		k.Logger(ctx).Error("Unable to fetch last span", "Error", err)
		return common.ErrSpanNotFound(k.Codespace()).Result()
	}

	// check all conditions
	if lastSpan.ID+1 != msg.ID || msg.StartBlock < lastSpan.StartBlock || msg.EndBlock < msg.StartBlock {
		k.Logger(ctx).Error("Blocks not in countinuity",
			"lastSpanId", lastSpan.ID,
			"lastSpanStartBlock", lastSpan.StartBlock,
			"spanId", msg.ID,
			"spanStartBlock", msg.StartBlock,
		)
		return common.ErrSpanNotInCountinuity(k.Codespace()).Result()
	}

	// freeze for new span
	err = k.FreezeSet(ctx, msg.ID, msg.StartBlock, msg.ChainID)
	if err != nil {
		k.Logger(ctx).Error("Unable to freeze validator set for span", "Error", err)
		return common.ErrUnableToFreezeValSet(k.Codespace()).Result()
	}

	// get last span
	lastSpan, err = k.GetLastSpan(ctx)
	if err != nil {
		k.Logger(ctx).Error("Unable to fetch last span", "Error", err)
		return common.ErrSpanNotFound(k.Codespace()).Result()
	}

	// TODO add check for duration
	result := sdk.Result{
		Tags: sdk.NewTags(
			tags.Success, []byte("true"),
			tags.BorSyncID, []byte(strconv.FormatUint(uint64(msg.ID), 10)),
			tags.SpanID, []byte(strconv.FormatUint(uint64(msg.ID), 10)),
			tags.SpanStartBlock, []byte(strconv.FormatUint(uint64(msg.StartBlock), 10)),
		),
	}
	return result
}
