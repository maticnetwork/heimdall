package checkpoint

import (
	"encoding/hex"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// NewSideTxHandler returns a side handler for "bank" type messages.
func NewSideTxHandler(k Keeper, contractCaller helper.IContractCaller) hmTypes.SideTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

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

// SideHandleMsgCheckpoint handles MsgCheckpoint message for external call
func SideHandleMsgCheckpoint(ctx sdk.Context, k Keeper, msg types.MsgCheckpoint) (result abci.ResponseDeliverSideTx) {
	// get params
	params := k.GetParams(ctx)

	// logger
	logger := k.Logger(ctx)

	// validate checkpoint
	validCheckpoint, err := types.ValidateCheckpoint(msg.StartBlock, msg.EndBlock, msg.RootHash, params.MaxCheckpointLength)
	if err != nil {
		logger.Error("Error validating checkpoint",
			"error", err,
			"startBlock", msg.StartBlock,
			"endBlock", msg.EndBlock,
		)
	} else if validCheckpoint {
		// vote `yes` if checkpoint is valid
		result.Result = abci.SideTxResultType_Yes
		return
	}

	logger.Error(
		"RootHash is not valid",
		"startBlock", msg.StartBlock,
		"endBlock", msg.EndBlock,
		"rootHash", msg.RootHash,
	)

	// vote `skip`
	result.Result = abci.SideTxResultType_Skip
	// set code and codespace
	result.Code = uint32(common.CodeInvalidBlockInput)
	result.Codespace = string(k.Codespace())

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
	logger := k.Logger(ctx)

	// Skip handler if checkpoint is not approved
	if sideTxResult != abci.SideTxResultType_Yes {
		logger.Debug("Skipping new checkpoint since side-tx didn't get yes votes", "startBlock", msg.StartBlock, "endBlock", msg.EndBlock, "rootHash", msg.RootHash)
		return common.ErrBadBlockDetails(k.Codespace()).Result()
	}

	//
	// Save checkpoint to buffer store
	//

	timeStamp := uint64(ctx.BlockTime().Unix())

	// Add checkpoint to buffer with root hash and account hash
	k.SetCheckpointBuffer(ctx, hmTypes.CheckpointBlockHeader{
		StartBlock:      msg.StartBlock,
		EndBlock:        msg.EndBlock,
		RootHash:        msg.RootHash,
		AccountRootHash: msg.AccountRootHash,
		Proposer:        msg.Proposer,
		BorChainID:      msg.BorChainID,
		TimeStamp:       timeStamp,
	})

	logger.Debug("Save new checkpoint into buffer", "startBlock", msg.StartBlock, "endBlock", msg.EndBlock, "rootHash", msg.RootHash)

	// Emit events for checkpoints
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),
		),
		sdk.NewEvent(
			types.EventTypeCheckpoint,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyStartBlock, strconv.FormatUint(uint64(msg.StartBlock), 10)),
			sdk.NewAttribute(types.AttributeKeyEndBlock, strconv.FormatUint(uint64(msg.EndBlock), 10)),
		),
	})

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
