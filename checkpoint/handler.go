package checkpoint

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCheckpoint:
			return handleMsgCheckpoint(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("Invalid message in checkpoint module ").Result()
		}
	}
}

func handleMsgCheckpoint(ctx sdk.Context, msg MsgCheckpoint, k Keeper) sdk.Result {
	logger := ctx.Logger().With("module", "checkpoint")
	valid := validateCheckpoint(int(msg.StartBlock), int(msg.EndBlock), msg.RootHash.String())

	// check msg.proposer with tm proposer
	var res int64
	if valid {
		logger.Debug("root hash matched!!")
		res = k.AddCheckpoint(ctx, msg.StartBlock, msg.EndBlock, msg.RootHash, msg.Proposer)
	} else {
		logger.Debug("Root hash no match ;(")
		return ErrBadBlockDetails(k.codespace).Result()

	}
	var out CheckpointBlockHeader
	json.Unmarshal(k.GetCheckpoint(ctx, res), &out)

	//TODO add validation
	// send tags
	return sdk.Result{}
}
