package bor

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/common"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

// NewHandler returns a handler for "bor" type messages.
func NewHandler(k common.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		common.InitBorLogger(&ctx)
		switch msg := msg.(type) {
		case MsgProposeSpan:
			return HandleMsgProposeSpan(ctx, msg, k, common.BorLogger)
		default:
			return sdk.ErrTxDecode("Invalid message in bor module").Result()
		}
	}
}

// HandleMsgProposeSpan handles proposeSpan msg
func HandleMsgProposeSpan(ctx sdk.Context, msg MsgProposeSpan, k common.Keeper, logger tmlog.Logger) sdk.Result {
	logger.Debug("Proposing span", "TxData", msg)
	// check if last span is up or if greater diff than threshold is found between validator set
	lastSpan, err := k.GetLastSpan(ctx)
	if err != nil {
		return common.ErrSpanNotInCountinuity(k.Codespace).Result()
	}
	// check if lastStart + 1 =  newStart
	if lastSpan.StartBlock+1 != msg.StartBlock {
		return common.ErrSpanNotInCountinuity(k.Codespace).Result()
	}

	// freeze for new span
	err = k.FreezeSet(ctx, msg.StartBlock)
	if err != nil {
		return common.ErrSpanNotInCountinuity(k.Codespace).Result()
	}
	// send tags
	return sdk.Result{}
}
