package checkpoint

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// NewSideTxHandler returns a side handler for "bank" type messages.
func NewSideTxHandler(k Keeper, contractCaller helper.IContractCaller) hmTypes.SideTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
		switch msg := msg.(type) {
		case types.MsgCheckpoint:
			return SideHandleMsgCheckpoint(ctx, k, msg)
		case types.MsgCheckpointAck:
			return SideHandleMsgCheckpointAck(ctx, k, msg)
		default:
			return abci.ResponseDeliverSideTx{
				Code: uint32(sdk.CodeUnknownRequest),
			}
		}
	}
}

// SideHandleMsgCheckpoint handles MsgCheckpointAck message for external call
func SideHandleMsgCheckpoint(ctx sdk.Context, k Keeper, msg types.MsgCheckpoint) (result abci.ResponseDeliverSideTx) {
	fmt.Println("[==] In SideHandleMsgCheckpoint doing external call")
	fmt.Println("[==] In SideHandleMsgCheckpoint  txbytes", hex.EncodeToString(ctx.TxBytes()), "isChckTx", ctx.IsCheckTx())

	// say `yes`
	result.Result = abci.SideTxResultType_Yes

	return
}

// SideHandleMsgCheckpointAck handles MsgCheckpointAck message for external call
func SideHandleMsgCheckpointAck(ctx sdk.Context, k Keeper, msg types.MsgCheckpointAck) (result abci.ResponseDeliverSideTx) {
	fmt.Println("[==] In SideHandleMsgCheckpointAck doing external call")
	fmt.Println("[==] In SideHandleMsgCheckpointAck  txbytes", hex.EncodeToString(ctx.TxBytes()), "isChckTx", ctx.IsCheckTx())

	// say `yes`
	result.Result = abci.SideTxResultType_Yes

	return
}

//
// Tx handler
//

// NewPostTxHandler returns a side handler for "bank" type messages.
func NewPostTxHandler(k Keeper, contractCaller helper.IContractCaller) hmTypes.PostTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg, sideTxResult abci.SideTxResultType) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgCheckpoint:
			return PostHandleMsgCheckpoint(ctx, k, msg, sideTxResult)
		case types.MsgCheckpointAck:
			return PostHandleMsgCheckpointAck(ctx, k, msg, sideTxResult)
		default:
			errMsg := "Unrecognized checkpoint Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// PostHandleMsgCheckpoint handles msg checkpoint
func PostHandleMsgCheckpoint(ctx sdk.Context, k Keeper, msg types.MsgCheckpoint, sideTxResult abci.SideTxResultType) sdk.Result {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// PostHandleMsgCheckpointAck handles msg checkpoint ack
func PostHandleMsgCheckpointAck(ctx sdk.Context, k Keeper, msg types.MsgCheckpointAck, sideTxResult abci.SideTxResultType) sdk.Result {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
		),
	)

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
