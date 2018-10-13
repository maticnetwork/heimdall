package sideBlock

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		// NOTE msg already has validate basic run
		switch msg := msg.(type) {
		case MsgSideBlock:
			return handleMsgSideBlock(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("Invalid message in side block module ").Result()
		}
	}
}

func handleMsgSideBlock(ctx sdk.Context, msg MsgSideBlock, k Keeper) sdk.Result {
	// TODO make web3 call using block hash , validate data here and return true in tags if true
	k.addBlock(ctx, msg.BlockHash, msg.TxRoot, msg.ReceiptRoot)

	//blockhash := msg.BlockHash
	//TODO check if block already exists and throw error (low priority)
	//success := getBlockDetails(blockhash,msg.TxRoot,msg.ReceiptRoot)
	success := true
	//success :=true
	if success {
		tags := sdk.NewTags("action", []byte("SideBlock"), "SideBlock", []byte(msg.BlockHash))
		return sdk.Result{
			// TODO return block data here
			Tags: tags,
		}
	} else {
		return ErrBadBlockDetails(k.codespace).Result()

	}

}
