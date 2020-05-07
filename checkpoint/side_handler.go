package checkpoint

import (
	"bytes"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmTypes "github.com/tendermint/tendermint/types"

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
			return SideHandleMsgCheckpointAck(ctx, k, msg, contractCaller)
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

	return common.ErrorSideTx(k.Codespace(), common.CodeInvalidBlockInput)
}

// SideHandleMsgCheckpointAck handles MsgCheckpointAck message for external call
func SideHandleMsgCheckpointAck(ctx sdk.Context, k Keeper, msg types.MsgCheckpointAck, contractCaller helper.IContractCaller) (result abci.ResponseDeliverSideTx) {
	logger := k.Logger(ctx)

	// make call to headerBlock with header number
	chainParams := k.ck.GetParams(ctx).ChainParams

	//
	// Validate data from root chain
	//

	rootChainInstance, err := contractCaller.GetRootChainInstance(chainParams.RootChainAddress.EthAddress())
	if err != nil {
		logger.Error("Unable to fetch rootchain contract instance", "error", err)
		return common.ErrorSideTx(k.Codespace(), common.CodeInvalidACK)
	}

	root, start, end, _, proposer, err := contractCaller.GetHeaderInfo(msg.HeaderBlock, rootChainInstance)
	if err != nil {
		logger.Error("Unable to fetch header from rootchain contract", "error", err, "headerBlock", msg.HeaderBlock)
		return common.ErrorSideTx(k.Codespace(), common.CodeInvalidACK)
	}

	// check if message data matches with contract data
	if msg.StartBlock != start ||
		msg.EndBlock != end ||
		!msg.Proposer.Equals(proposer) ||
		!bytes.Equal(msg.RootHash.Bytes(), root.Bytes()) {

		logger.Error("Invalid message. It doesn't match with contract state", "error", err, "headerBlock", msg.HeaderBlock)
		return common.ErrorSideTx(k.Codespace(), common.CodeInvalidACK)
	}

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

	checkpointBuffer, err := k.GetCheckpointFromBuffer(ctx)
	if err == nil && checkpointBuffer != nil {
		logger.Debug("Checkpoint already exists in buffer")

		// get checkpoint buffer time from params
		params := k.GetParams(ctx)
		expiryTime := checkpointBuffer.TimeStamp + uint64(params.CheckpointBufferTime.Seconds())

		// return with error (ack is required)
		return common.ErrNoACK(k.Codespace(), expiryTime).Result()
	}

	timeStamp := uint64(ctx.BlockTime().Unix())

	// Add checkpoint to buffer with root hash and account hash
	k.SetCheckpointBuffer(ctx, hmTypes.CheckpointBlockHeader{
		StartBlock: msg.StartBlock,
		EndBlock:   msg.EndBlock,
		RootHash:   msg.RootHash,
		Proposer:   msg.Proposer,
		BorChainID: msg.BorChainID,
		TimeStamp:  timeStamp,
	})

	logger.Debug("New checkpoint into buffer stored",
		"startBlock", msg.StartBlock,
		"endBlock", msg.EndBlock,
		"rootHash", msg.RootHash,
	)

	// TX bytes
	txBytes := ctx.TxBytes()
	hash := tmTypes.Tx(txBytes).Hash()

	// Emit event for checkpoints
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCheckpoint,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                  // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),                // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmTypes.BytesToHeimdallHash(hash).Hex()), // tx hash
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()),             // result
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyStartBlock, strconv.FormatUint(msg.StartBlock, 10)),
			sdk.NewAttribute(types.AttributeKeyEndBlock, strconv.FormatUint(msg.EndBlock, 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

// PostHandleMsgCheckpointAck handles msg checkpoint ack
func PostHandleMsgCheckpointAck(ctx sdk.Context, k Keeper, msg types.MsgCheckpointAck, sideTxResult abci.SideTxResultType) sdk.Result {
	logger := k.Logger(ctx)

	// Skip handler if checkpoint-ack is not approved
	if sideTxResult != abci.SideTxResultType_Yes {
		logger.Debug("Skipping new checkpoint-ack since side-tx didn't get yes votes", "headerBlock", msg.HeaderBlock)
		return common.ErrBadBlockDetails(k.Codespace()).Result()
	}

	// get last checkpoint from buffer
	headerBlock, err := k.GetCheckpointFromBuffer(ctx)
	if err != nil {
		logger.Error("Unable to get checkpoint buffer", "error", err)
		return common.ErrBadAck(k.Codespace()).Result()
	}

	// adjust checkpoint data if latest checkpoint is already submitted
	if headerBlock.EndBlock > msg.EndBlock {
		logger.Info("Adjusting endBlock to one already submitted on chain", "endBlock", headerBlock.EndBlock, "adjustedEndBlock", msg.EndBlock)
		headerBlock.EndBlock = msg.EndBlock
		headerBlock.RootHash = msg.RootHash
		headerBlock.Proposer = msg.Proposer
	}

	//
	// Update checkpoint state
	//

	// Add checkpoint to headerBlocks
	if err := k.AddCheckpoint(ctx, msg.HeaderBlock, *headerBlock); err != nil {
		logger.Error("Error while adding checkpoint into store", "headerBlock", msg.HeaderBlock)
		return sdk.ErrInternal("Failed to add checkpoint into store").Result()
	}
	logger.Debug("Checkpoint added to store", "headerBlock", msg.HeaderBlock)

	// Flush buffer
	k.FlushCheckpointBuffer(ctx)
	logger.Debug("Checkpoint buffer flushed after receiving checkpoint ack")

	// Update ack count in staking module
	k.UpdateACKCount(ctx)
	logger.Info("Valid ack received", "CurrentACKCount", k.GetACKCount(ctx)-1, "UpdatedACKCount", k.GetACKCount(ctx))

	// Increment accum (selects new proposer)
	k.sk.IncrementAccum(ctx, 1)

	// TX bytes
	txBytes := ctx.TxBytes()
	hash := tmTypes.Tx(txBytes).Hash()

	// Emit event for checkpoints
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCheckpointAck,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                  // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),                // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmTypes.BytesToHeimdallHash(hash).Hex()), // tx hash
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()),             // result
			sdk.NewAttribute(types.AttributeKeyHeaderIndex, strconv.FormatUint(uint64(msg.HeaderBlock), 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
