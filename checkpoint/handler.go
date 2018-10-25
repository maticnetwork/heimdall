package checkpoint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	conf "github.com/maticnetwork/heimdall/helper"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgCheckpoint:
			// redirect to handle msg checkpoint
			return handleMsgCheckpoint(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("Invalid message in checkpoint module ").Result()
		}
	}
}

func handleMsgCheckpoint(ctx sdk.Context, msg MsgCheckpoint, k Keeper) sdk.Result {
	logger := conf.Logger.With("module", "checkpoint")

	// check if the roothash provided is valid for start and end
	valid := validateCheckpoint(int(msg.StartBlock), int(msg.EndBlock), msg.RootHash.String())

	// check msg.proposer with tm proposer
	var key int64
	if valid {

		// add checkpoint to state if rootHash matches
		key = k.AddCheckpoint(ctx, msg.StartBlock, msg.EndBlock, msg.RootHash, msg.Proposer)
		logger.Info("root hash matched ! ", "Key", key)

	} else {

		logger.Info("Root hash doesnt match ;(")
		// return Bad Block Error
		return ErrBadBlockDetails(k.codespace).Result()
	}

	//TODO add validation
	// send tags
	return sdk.Result{}
}
