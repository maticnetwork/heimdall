package checkpoint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCheckpoint:
			// redirect to handle msg checkpoint
			return handleMsgCheckpoint(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("Invalid message in checkpoint module").Result()
		}
	}
}

func handleMsgCheckpoint(ctx sdk.Context, msg MsgCheckpoint, k Keeper) sdk.Result {
	if err := msg.ValidateBasic(); err != nil { return ErrBadBlockDetails(k.codespace).Result() }

	// check if the roothash provided is valid for start and end
	valid := ValidateCheckpoint(msg.StartBlock, msg.EndBlock, msg.RootHash.String())
	CheckpointLogger.Debug("Validating last block from main chain", "lastBlock", helper.GetLastBlock(), "startBlock", msg.StartBlock)
	if valid && helper.GetLastBlock() == msg.StartBlock {
		// add checkpoint to state if rootHash matches
		key := k.AddCheckpoint(ctx, msg.StartBlock, msg.EndBlock, msg.RootHash)
		CheckpointLogger.Debug("RootHash matched!", "key", key)
	} else {
		CheckpointLogger.Debug("Invalid checkpoint.")
		// return Bad Block Error
		return ErrBadBlockDetails(k.codespace).Result()
	}

	// send tags
	return sdk.Result{}
}
