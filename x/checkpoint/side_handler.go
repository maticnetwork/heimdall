package checkpoint

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	tmprototypes "github.com/tendermint/tendermint/proto/tendermint/types"
	tmTypes "github.com/tendermint/tendermint/types"

	borCommon "github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/x/checkpoint/keeper"
	"github.com/maticnetwork/heimdall/x/checkpoint/types"
)

// NewSideTxHandler returns a side handler for "bank" type messages.
func NewSideTxHandler(k keeper.Keeper, contractCaller helper.IContractCaller) hmTypes.SideTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCheckpoint:
			return SideHandleMsgCheckpoint(ctx, k, *msg, contractCaller)
		case *types.MsgCheckpointAck:
			return SideHandleMsgCheckpointAck(ctx, k, *msg, contractCaller)
		default:
			return abci.ResponseDeliverSideTx{
				Code: uint32(6), // TODO should be changed like `sdk.CodeUnknownRequest`
			}
		}
	}
}

// SideHandleMsgCheckpoint handles MsgCheckpoint message for external call
func SideHandleMsgCheckpoint(ctx sdk.Context, k keeper.Keeper, msg types.MsgCheckpoint, contractCaller helper.IContractCaller) (result abci.ResponseDeliverSideTx) {
	// get params
	params := k.GetParams(ctx)

	// logger
	logger := k.Logger(ctx)

	// validate checkpoint
	validCheckpoint, err := types.ValidateCheckpoint(msg.StartBlock, msg.EndBlock, hmCommonTypes.BytesToHeimdallHash(msg.RootHash), params.MaxCheckpointLength, contractCaller)
	if err != nil {
		logger.Error("Error validating checkpoint",
			"error", err,
			"startBlock", msg.StartBlock,
			"endBlock", msg.EndBlock,
		)
	} else if validCheckpoint {
		// vote `yes` if checkpoint is valid
		result.Result = tmprototypes.SideTxResultType_YES
		return
	}

	logger.Error(
		"RootHash is not valid",
		"startBlock", msg.StartBlock,
		"endBlock", msg.EndBlock,
		"rootHash", msg.RootHash,
	)

	return
}

// SideHandleMsgCheckpointAck handles MsgCheckpointAck message for external call
func SideHandleMsgCheckpointAck(ctx sdk.Context, k keeper.Keeper, msg types.MsgCheckpointAck, contractCaller helper.IContractCaller) (result abci.ResponseDeliverSideTx) {
	logger := k.Logger(ctx)

	params := k.GetParams(ctx)
	chainParams := k.Ck.GetParams(ctx).ChainParams

	//
	// Validate data from root chain
	//

	rootChainInstance, err := contractCaller.GetRootChainInstance(borCommon.BytesToAddress(chainParams.RootChainAddress.Bytes()))
	if err != nil {
		logger.Error("Unable to fetch rootchain contract instance", "error", err)
		// TODO fix this
		// return common.ErrorSideTx(k.Codespace(), common.CodeInvalidACK)
		return
	}

	root, start, end, _, proposer, err := contractCaller.GetHeaderInfo(msg.Number, rootChainInstance, params.ChildBlockInterval)
	if err != nil {
		logger.Error("Unable to fetch checkpoint from rootchain", "error", err, "checkpointNumber", msg.Number)
		// TODO fix this
		// return common.ErrorSideTx(k.Codespace(), common.CodeInvalidACK)
		return
	}

	// check if message data matches with contract data
	if msg.StartBlock != start ||
		msg.EndBlock != end ||
		msg.Proposer != proposer.String() ||
		!bytes.Equal([]byte(msg.RootHash), root.Bytes()) {

		logger.Error("Invalid message. It doesn't match with contract state", "error", err, "checkpointNumber", msg.Number)
		// TODO fix this
		// return common.ErrorSideTx(k.Codespace(), common.CodeInvalidACK)
		return
	}

	// say `yes`
	result.Result = tmprototypes.SideTxResultType_YES

	return
}

//
// Tx handler
//

// NewPostTxHandler returns a side handler for "bank" type messages.
func NewPostTxHandler(k keeper.Keeper, contractCaller helper.IContractCaller) hmTypes.PostTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg, sideTxResult tmprototypes.SideTxResultType) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case *types.MsgCheckpoint:
			return PostHandleMsgCheckpoint(ctx, k, *msg, sideTxResult)
		case *types.MsgCheckpointAck:
			return PostHandleMsgCheckpointAck(ctx, k, *msg, sideTxResult)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// PostHandleMsgCheckpoint handles msg checkpoint
func PostHandleMsgCheckpoint(ctx sdk.Context, k keeper.Keeper, msg types.MsgCheckpoint, sideTxResult tmprototypes.SideTxResultType) (*sdk.Result, error) {
	logger := k.Logger(ctx)

	// Skip handler if checkpoint is not approved
	if sideTxResult != tmprototypes.SideTxResultType_YES {
		logger.Debug("Skipping new checkpoint since side-tx didn't get yes votes", "startBlock", msg.StartBlock, "endBlock", msg.EndBlock, "rootHash", msg.RootHash)
		return nil, types.ErrBadBlockDetails
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
			return nil, types.ErrOldCheckpoint
		}

		// check if new checkpoint's start block start from current tip
		if lastCheckpoint.EndBlock+1 != msg.StartBlock {
			logger.Error("Checkpoint not in countinuity",
				"currentTip", lastCheckpoint.EndBlock,
				"startBlock", msg.StartBlock)
			return nil, types.ErrDisCountinuousCheckpoint
		}
	} else if err.Error() == types.ErrNoCheckpointFound.Error() && msg.StartBlock != 0 {
		logger.Error("First checkpoint to start from block 0", "Error", err)
		return nil, types.ErrBadBlockDetails
	}

	//
	// Save checkpoint to buffer store
	//

	checkpointBuffer, err := k.GetCheckpointFromBuffer(ctx)
	if err == nil && checkpointBuffer != nil {
		logger.Debug("Checkpoint already exists in buffer")

		// get checkpoint buffer time from params
		// params := k.GetParams(ctx)
		// expiryTime := checkpointBuffer.TimeStamp + uint64(params.CheckpointBufferTime.Seconds())

		// return with error (ack is required)
		return nil, types.ErrNoACK
	}

	timeStamp := uint64(ctx.BlockTime().Unix())

	// Add checkpoint to buffer with root hash and account hash
	k.SetCheckpointBuffer(ctx, &hmTypes.Checkpoint{
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
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                        // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),                      // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmCommonTypes.BytesToHeimdallHash(hash).Hex()), // tx hash
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()),                   // result
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer),
			sdk.NewAttribute(types.AttributeKeyStartBlock, strconv.FormatUint(msg.StartBlock, 10)),
			sdk.NewAttribute(types.AttributeKeyEndBlock, strconv.FormatUint(msg.EndBlock, 10)),
			sdk.NewAttribute(types.AttributeKeyRootHash, hex.EncodeToString(msg.RootHash)),
			sdk.NewAttribute(types.AttributeKeyAccountHash, hex.EncodeToString(msg.AccountRootHash)),
		),
	})

	return &sdk.Result{}, nil

}

// PostHandleMsgCheckpointAck handles msg checkpoint ack
func PostHandleMsgCheckpointAck(ctx sdk.Context, k keeper.Keeper, msg types.MsgCheckpointAck, sideTxResult tmprototypes.SideTxResultType) (*sdk.Result, error) {
	logger := k.Logger(ctx)

	// Skip handler if checkpoint-ack is not approved
	if sideTxResult != tmprototypes.SideTxResultType_YES {
		logger.Debug("Skipping new checkpoint-ack since side-tx didn't get yes votes", "checkpointNumber", msg.Number)
		return nil, types.ErrBadBlockDetails
	}

	// get last checkpoint from buffer
	checkpointObj, err := k.GetCheckpointFromBuffer(ctx)
	if err != nil {
		logger.Error("Unable to get checkpoint buffer", "error", err)
		return nil, types.ErrBadAck
	}

	// invalid start block
	if msg.StartBlock != checkpointObj.StartBlock {
		logger.Error("Invalid start block", "startExpected", checkpointObj.StartBlock, "startReceived", msg.StartBlock)
		return nil, types.ErrBadAck
	}

	// Return err if start and end matches but contract root hash doesn't match
	if msg.StartBlock == checkpointObj.StartBlock && msg.EndBlock == checkpointObj.EndBlock && !bytes.Equal([]byte(msg.RootHash), checkpointObj.RootHash) {
		logger.Error("Invalid ACK",
			"startExpected", checkpointObj.StartBlock,
			"startReceived", msg.StartBlock,
			"endExpected", checkpointObj.EndBlock,
			"endReceived", msg.StartBlock,
			"rootExpected", hex.EncodeToString(checkpointObj.RootHash),
			"rootRecieved", msg.RootHash,
		)
		return nil, types.ErrBadAck
	}

	// adjust checkpoint data if latest checkpoint is already submitted
	if checkpointObj.EndBlock > msg.EndBlock {
		logger.Info("Adjusting endBlock to one already submitted on chain", "endBlock", checkpointObj.EndBlock, "adjustedEndBlock", msg.EndBlock)
		checkpointObj.EndBlock = msg.EndBlock
		checkpointObj.RootHash = []byte(msg.RootHash)
		checkpointObj.Proposer = msg.Proposer
	}

	//
	// Update checkpoint state
	//

	// Add checkpoint to store
	if err := k.AddCheckpoint(ctx, msg.Number, checkpointObj); err != nil {
		logger.Error("Error while adding checkpoint into store", "checkpointNumber", msg.Number)
		// TODO
		// return sdk.ErrInternal("Failed to add checkpoint into store").Result()
		return nil, types.ErrNoCheckpointFound
	}
	logger.Debug("Checkpoint added to store", "checkpointNumber", msg.Number)

	// Flush buffer
	k.FlushCheckpointBuffer(ctx)
	logger.Debug("Checkpoint buffer flushed after receiving checkpoint ack")

	// Update ack count in staking module
	k.UpdateACKCount(ctx)
	logger.Info("Valid ack received", "CurrentACKCount", k.GetACKCount(ctx)-1, "UpdatedACKCount", k.GetACKCount(ctx))

	// Increment accum (selects new proposer)
	k.Sk.IncrementAccum(ctx, 1)

	// TX bytes
	txBytes := ctx.TxBytes()
	hash := tmTypes.Tx(txBytes).Hash()

	// Emit event for checkpoints
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeCheckpointAck,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                        // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),                      // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmCommonTypes.BytesToHeimdallHash(hash).Hex()), // tx hash
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()),                   // result
			sdk.NewAttribute(types.AttributeKeyHeaderIndex, strconv.FormatUint(msg.Number, 10)),
		),
	})

	return &sdk.Result{}, nil
}
