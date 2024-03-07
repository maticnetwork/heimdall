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
		case types.MsgCheckpointAdjust:
			return SideHandleMsgCheckpointAdjust(ctx, k, msg, contractCaller)
		case types.MsgCheckpoint:
			return SideHandleMsgCheckpoint(ctx, k, msg, contractCaller)
		case types.MsgCheckpointAck:
			return SideHandleMsgCheckpointAck(ctx, k, msg, contractCaller)
		case types.MsgMilestone:
			return SideHandleMsgMilestone(ctx, k, msg, contractCaller)
		default:
			return abci.ResponseDeliverSideTx{
				Code: uint32(sdk.CodeUnknownRequest),
			}
		}
	}
}

// SideHandleMsgCheckpointAdjust handles MsgCheckpointAdjust message for external call
func SideHandleMsgCheckpointAdjust(ctx sdk.Context, k Keeper, msg types.MsgCheckpointAdjust, contractCaller helper.IContractCaller) (result abci.ResponseDeliverSideTx) {
	logger := k.Logger(ctx)
	chainParams := k.ck.GetParams(ctx).ChainParams
	params := k.GetParams(ctx)

	checkpointBuffer, err := k.GetCheckpointFromBuffer(ctx)
	if checkpointBuffer != nil {
		logger.Error("checkpoint buffer", "error", err)
		return common.ErrorSideTx(k.Codespace(), common.CodeCheckpointBuffer)
	}

	checkpointObj, err := k.GetCheckpointByNumber(ctx, msg.HeaderIndex)
	if err != nil {
		logger.Error("Unable to get checkpoint from db", "header index", msg.HeaderIndex, "error", err)
		return common.ErrorSideTx(k.Codespace(), common.CodeNoCheckpoint)
	}

	rootChainInstance, err := contractCaller.GetRootChainInstance(chainParams.RootChainAddress.EthAddress())
	if err != nil {
		logger.Error("Unable to fetch rootchain contract instance", "eth address", chainParams.RootChainAddress.EthAddress(), "error", err)
		return common.ErrorSideTx(k.Codespace(), common.CodeOldCheckpoint)
	}

	root, start, end, _, proposer, err := contractCaller.GetHeaderInfo(msg.HeaderIndex, rootChainInstance, params.ChildBlockInterval)
	if err != nil {
		logger.Error("Unable to fetch checkpoint from rootchain", "checkpointNumber", msg.HeaderIndex, "error", err)
		return common.ErrorSideTx(k.Codespace(), common.CodeNoCheckpoint)
	}

	if checkpointObj.EndBlock == end && checkpointObj.StartBlock == start && bytes.Equal(checkpointObj.RootHash.Bytes(), root.Bytes()) && bytes.Equal(checkpointObj.Proposer.Bytes(), proposer.Bytes()) {
		logger.Error("Same Checkpoint in DB")
		return common.ErrorSideTx(k.Codespace(), common.CodeCheckpointAlreadyExists)
	}

	if msg.EndBlock != end || msg.StartBlock != start || !bytes.Equal(msg.RootHash.Bytes(), root.Bytes()) || !bytes.Equal(msg.Proposer.Bytes(), proposer.Bytes()) {
		logger.Error("Checkpoint on Rootchain is not same as msg",
			"message start block", msg.StartBlock,
			"Rootchain Checkpoint start block", start,
			"message end block", msg.EndBlock,
			"Rootchain Checkpointt end block", end,
			"message proposer", msg.Proposer,
			"Rootchain Checkpoint proposer", proposer,
			"message root hash", msg.RootHash,
			"Rootchain Checkpoint root hash", root,
		)

		return common.ErrorSideTx(k.Codespace(), common.CodeCheckpointAlreadyExists)
	}

	result.Result = abci.SideTxResultType_Yes

	return
}

// SideHandleMsgCheckpoint handles MsgCheckpoint message for external call
func SideHandleMsgCheckpoint(ctx sdk.Context, k Keeper, msg types.MsgCheckpoint, contractCaller helper.IContractCaller) (result abci.ResponseDeliverSideTx) {
	// get params
	params := k.GetParams(ctx)
	maticTxConfirmations := k.ck.GetParams(ctx).MaticchainTxConfirmations

	// logger
	logger := k.Logger(ctx)

	// validate checkpoint
	validCheckpoint, err := types.ValidateCheckpoint(msg.StartBlock, msg.EndBlock, msg.RootHash, params.MaxCheckpointLength, contractCaller, maticTxConfirmations)
	if err != nil {
		logger.Error("Error validating checkpoint",
			"startBlock", msg.StartBlock,
			"endBlock", msg.EndBlock,
			"rootHash", msg.RootHash,
			"error", err,
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

	params := k.GetParams(ctx)
	chainParams := k.ck.GetParams(ctx).ChainParams

	//
	// Validate data from root chain
	//

	rootChainInstance, err := contractCaller.GetRootChainInstance(chainParams.RootChainAddress.EthAddress())
	if err != nil {
		logger.Error("Unable to fetch rootchain contract instance",
			"eth address", chainParams.RootChainAddress.EthAddress(),
			"error", err,
		)

		return common.ErrorSideTx(k.Codespace(), common.CodeInvalidACK)
	}

	root, start, end, _, proposer, err := contractCaller.GetHeaderInfo(msg.Number, rootChainInstance, params.ChildBlockInterval)
	if err != nil {
		logger.Error("Unable to fetch checkpoint from rootchain", "checkpointNumber", msg.Number, "error", err)
		return common.ErrorSideTx(k.Codespace(), common.CodeInvalidACK)
	}

	// check if message data matches with contract data
	if msg.StartBlock != start ||
		msg.EndBlock != end ||
		!msg.Proposer.Equals(proposer) ||
		!bytes.Equal(msg.RootHash.Bytes(), root.Bytes()) {
		logger.Error("Invalid message. It doesn't match with contract state",
			"checkpointNumber", msg.Number,
			"message start block", msg.StartBlock,
			"Rootchain Checkpoint start block", start,
			"message end block", msg.EndBlock,
			"Rootchain Checkpointt end block", end,
			"message proposer", msg.Proposer,
			"Rootchain Checkpoint proposer", proposer,
			"message root hash", msg.RootHash,
			"Rootchain Checkpoint root hash", root,
			"error", err,
		)

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
		case types.MsgCheckpointAdjust:
			return PostHandleMsgCheckpointAdjust(ctx, k, msg, sideTxResult, contractCaller)
		case types.MsgCheckpoint:
			return PostHandleMsgCheckpoint(ctx, k, msg, sideTxResult)
		case types.MsgCheckpointAck:
			return PostHandleMsgCheckpointAck(ctx, k, msg, sideTxResult)
		case types.MsgMilestone:
			return PostHandleMsgMilestone(ctx, k, msg, sideTxResult)
		default:
			return sdk.ErrUnknownRequest("Unrecognized checkpoint Msg type").Result()
		}
	}
}

// PostHandleMsgCheckpointAdjust handles msg checkpoint adjust
func PostHandleMsgCheckpointAdjust(ctx sdk.Context, k Keeper, msg types.MsgCheckpointAdjust, sideTxResult abci.SideTxResultType, contractCaller helper.IContractCaller) sdk.Result {
	logger := k.Logger(ctx)

	// Skip handler if checkpoint-adjust is not approved
	if sideTxResult != abci.SideTxResultType_Yes {
		logger.Debug("Skipping new checkpoint-adjust since side-tx didn't get yes votes", "checkpointNumber", msg.HeaderIndex)
		return common.ErrBadBlockDetails(k.Codespace()).Result()
	}

	checkpointBuffer, err := k.GetCheckpointFromBuffer(ctx)
	if checkpointBuffer != nil {
		logger.Error("checkpoint buffer exists", "error", err)
		return common.ErrCheckpointBufferFound(k.Codespace()).Result()
	}

	checkpointObj, err := k.GetCheckpointByNumber(ctx, msg.HeaderIndex)
	if err != nil {
		logger.Error("Unable to get checkpoint from db",
			"checkpoint number", msg.HeaderIndex,
			"error", err)

		return common.ErrNoCheckpointFound(k.Codespace()).Result()
	}

	if checkpointObj.EndBlock == msg.EndBlock && checkpointObj.StartBlock == msg.StartBlock && bytes.Equal(checkpointObj.RootHash.Bytes(), msg.RootHash.Bytes()) && bytes.Equal(checkpointObj.Proposer.Bytes(), msg.Proposer.Bytes()) {
		logger.Error("Same Checkpoint in DB")
		return common.ErrCheckpointAlreadyExists(k.Codespace()).Result()
	}

	logger.Info("Previous checkpoint details: EndBlock -", checkpointObj.EndBlock, ", RootHash -", msg.RootHash, " Proposer -", checkpointObj.Proposer)

	checkpointObj.EndBlock = msg.EndBlock
	checkpointObj.RootHash = hmTypes.BytesToHeimdallHash(msg.RootHash.Bytes())
	checkpointObj.Proposer = msg.Proposer

	logger.Info("New checkpoint details: EndBlock -", checkpointObj.EndBlock, ", RootHash -", msg.RootHash, " Proposer -", checkpointObj.Proposer)

	//
	// Update checkpoint state
	//

	// Add checkpoint to store
	if err = k.AddCheckpoint(ctx, msg.HeaderIndex, checkpointObj); err != nil {
		logger.Error("Error while adding checkpoint into store", "checkpointNumber", msg.HeaderIndex)
		return sdk.ErrInternal("Failed to add checkpoint into store").Result()
	}

	logger.Debug("Checkpoint updated to store", "checkpointNumber", msg.HeaderIndex)

	// Emit event for checkpoints
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCheckpointAck,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                      // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),    // module name
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()), // result
			sdk.NewAttribute(types.AttributeKeyHeaderIndex, strconv.FormatUint(msg.HeaderIndex, 10)),
			sdk.NewAttribute(types.AttributeKeyStartBlock, strconv.FormatUint(msg.StartBlock, 10)),
			sdk.NewAttribute(types.AttributeKeyEndBlock, strconv.FormatUint(msg.EndBlock, 10)),
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyRootHash, msg.RootHash.String()),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
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
	// Validate last checkpoint
	//

	// fetch last checkpoint from store
	if lastCheckpoint, err := k.GetLastCheckpoint(ctx); err == nil {
		// make sure new checkpoint is after tip
		if lastCheckpoint.EndBlock > msg.StartBlock {
			logger.Error("Checkpoint already exists",
				"currentTip", lastCheckpoint.EndBlock,
				"startBlock", msg.StartBlock,
			)

			return common.ErrOldCheckpoint(k.Codespace()).Result()
		}

		// check if new checkpoint's start block start from current tip
		if lastCheckpoint.EndBlock+1 != msg.StartBlock {
			logger.Error("Checkpoint not in continuity",
				"currentTip", lastCheckpoint.EndBlock,
				"startBlock", msg.StartBlock)

			return common.ErrDisContinuousCheckpoint(k.Codespace()).Result()
		}
	} else if err.Error() == common.ErrNoCheckpointFound(k.Codespace()).Error() && msg.StartBlock != 0 {
		logger.Error("First checkpoint to start from block 0", "Error", err)
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
	if err = k.SetCheckpointBuffer(ctx, hmTypes.Checkpoint{
		StartBlock: msg.StartBlock,
		EndBlock:   msg.EndBlock,
		RootHash:   msg.RootHash,
		Proposer:   msg.Proposer,
		BorChainID: msg.BorChainID,
		TimeStamp:  timeStamp,
	}); err != nil {
		logger.Error("Failed to set checkpoint buffer", "Error", err)
	}

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
			sdk.NewAttribute(types.AttributeKeyRootHash, msg.RootHash.String()),
			sdk.NewAttribute(types.AttributeKeyAccountHash, msg.AccountRootHash.String()),
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
		logger.Debug("Skipping new checkpoint-ack since side-tx didn't get yes votes", "checkpointNumber", msg.Number)
		return common.ErrBadBlockDetails(k.Codespace()).Result()
	}

	// get last checkpoint from buffer
	checkpointObj, err := k.GetCheckpointFromBuffer(ctx)
	if err != nil {
		logger.Error("Unable to get checkpoint buffer", "error", err)
		return common.ErrBadAck(k.Codespace()).Result()
	}

	// invalid start block
	if msg.StartBlock != checkpointObj.StartBlock {
		logger.Error("Invalid start block", "startExpected", checkpointObj.StartBlock, "startReceived", msg.StartBlock)
		return common.ErrBadAck(k.Codespace()).Result()
	}

	// Return err if start and end matches but contract root hash doesn't match
	if msg.StartBlock == checkpointObj.StartBlock && msg.EndBlock == checkpointObj.EndBlock && !msg.RootHash.Equals(checkpointObj.RootHash) {
		logger.Error("Invalid ACK",
			"startExpected", checkpointObj.StartBlock,
			"startReceived", msg.StartBlock,
			"endExpected", checkpointObj.EndBlock,
			"endReceived", msg.StartBlock,
			"rootExpected", checkpointObj.RootHash.String(),
			"rootReceived", msg.RootHash.String(),
		)

		return common.ErrBadAck(k.Codespace()).Result()
	}

	// adjust checkpoint data if latest checkpoint is already submitted
	if ctx.BlockHeight() < helper.GetAalborgHardForkHeight() {
		if checkpointObj.EndBlock > msg.EndBlock {
			logger.Info("Adjusting endBlock to one already submitted on chain", "endBlock", checkpointObj.EndBlock, "adjustedEndBlock", msg.EndBlock)
			checkpointObj.EndBlock = msg.EndBlock
			checkpointObj.RootHash = msg.RootHash
			checkpointObj.Proposer = msg.Proposer
		}
	} else {
		if checkpointObj.EndBlock != msg.EndBlock {
			logger.Info("Adjusting endBlock to one already submitted on chain", "endBlock", checkpointObj.EndBlock, "adjustedEndBlock", msg.EndBlock)
			checkpointObj.EndBlock = msg.EndBlock
			checkpointObj.RootHash = msg.RootHash
			checkpointObj.Proposer = msg.Proposer
		}
	}

	//
	// Update checkpoint state
	//

	// Add checkpoint to store
	if err = k.AddCheckpoint(ctx, msg.Number, *checkpointObj); err != nil {
		logger.Error("Error while adding checkpoint into store", "checkpointNumber", msg.Number)
		return sdk.ErrInternal("Failed to add checkpoint into store").Result()
	}

	logger.Debug("Checkpoint added to store", "checkpointNumber", msg.Number)

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
			sdk.NewAttribute(types.AttributeKeyHeaderIndex, strconv.FormatUint(msg.Number, 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
