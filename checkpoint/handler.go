package checkpoint

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		// NOTE msg already has validate basic run
		switch msg := msg.(type) {
		case MsgCheckpoint:
			return handleMsgCheckpoint(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("Invalid message in checkpoint module ").Result()
		}
	}
}

func handleMsgCheckpoint(ctx sdk.Context, msg MsgCheckpoint, k Keeper) sdk.Result {
	//fmt.Printf("entered handler with message %v and %v and %v", msg.StartBlock, msg.EndBlock, msg.RootHash.String())
	//TODO check last block in last checkpoint (startBlock of new checkpoint == last block of prev endpoint)
	// TODO insert checkpoint in state
	logger := ctx.Logger().With("module", "x/baseapp")
	valid := validateCheckpoint(int(msg.StartBlock), int(msg.EndBlock), msg.RootHash.String())

	var res int64
	if valid {
		logger.Debug("root hash matched !! ")
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
