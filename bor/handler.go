package bor

import (
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
)

func NewHandler(k common.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		common.InitBorLogger(&ctx)

		switch msg := msg.(type) {
		case MsgProposeSpan:
			return HandleMsgProposeSpan(ctx, msg, k, common.BorLogger)
		default:
			return sdk.ErrTxDecode("Invalid message in checkpoint module").Result()
		}
	}
}

// HandleMsgProposeSpan handles proposeSpan msg
func HandleMsgProposeSpan(ctx sdk.Context, msg MsgCheckpoint, k common.Keeper, contractCaller helper.IContractCaller, logger tmlog.Logger) sdk.Result {
	logger.Debug("Validating Checkpoint Data", "TxData", msg)

	// send tags
	return sdk.Result{}
}
