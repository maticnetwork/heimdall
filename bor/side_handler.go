package bor

import (
	"bytes"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/bor/types"
	"github.com/maticnetwork/heimdall/common"
	hmCommon "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"

	hmTypes "github.com/maticnetwork/heimdall/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmTypes "github.com/tendermint/tendermint/types"
)

// NewSideTxHandler returns a side handler for "span" type messages.
func NewSideTxHandler(k Keeper, contractCaller helper.IContractCaller) hmTypes.SideTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgProposeSpan:
			return SideHandleMsgSpan(ctx, k, msg, contractCaller)
		default:
			return abci.ResponseDeliverSideTx{
				Code: uint32(sdk.CodeUnknownRequest),
			}
		}
	}
}

// NewPostTxHandler returns a side handler for "span" type messages.
func NewPostTxHandler(k Keeper, contractCaller helper.IContractCaller) hmTypes.PostTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg, sideTxResult abci.SideTxResultType) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgProposeSpan:
			return PostHandleMsgEventSpan(ctx, k, msg, sideTxResult)
		default:
			errMsg := "Unrecognized Span Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// SideHandleMsgSpan validates external calls required for processing proposed span
func SideHandleMsgSpan(ctx sdk.Context, k Keeper, msg types.MsgProposeSpan, contractCaller helper.IContractCaller) (result abci.ResponseDeliverSideTx) {
	k.Logger(ctx).Debug("✅ Validating External call for span msg",
		"msgSeed", msg.Seed.String(),
	)

	// calculate next span seed locally
	nextSpanSeed, err := k.GetNextSpanSeed(ctx)
	if err != nil {
		k.Logger(ctx).Error("Error fetching next span seed from mainchain")
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	// check if span seed matches or not.
	if !bytes.Equal(msg.Seed.Bytes(), nextSpanSeed.Bytes()) {
		k.Logger(ctx).Error(
			"Span Seed does not match",
			"msgSeed", msg.Seed.String(),
			"mainchainSeed", nextSpanSeed.String(),
		)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	// fetch current child block
	childBlock, err := contractCaller.GetMaticChainBlock(nil)
	if err != nil {
		k.Logger(ctx).Error("Error fetching current child block", "error", err)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	lastSpan, err := k.GetLastSpan(ctx)
	if err != nil {
		k.Logger(ctx).Error("Error fetching last span", "error", err)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	currentBlock := childBlock.Number.Uint64()
	// check if span proposed is in-turn or not
	if !(lastSpan.StartBlock <= currentBlock && currentBlock <= lastSpan.EndBlock) {
		k.Logger(ctx).Error(
			"Span proposed is not in-turn",
			"currentChildBlock", currentBlock,
			"msgStartblock", msg.StartBlock,
			"msgEndBlock", msg.EndBlock,
		)
		return hmCommon.ErrorSideTx(k.Codespace(), common.CodeInvalidMsg)
	}

	k.Logger(ctx).Debug("✅ Successfully validated External call for span msg")
	result.Result = abci.SideTxResultType_Yes
	return
}

// PostHandleMsgEventSpan handles state persisting span msg
func PostHandleMsgEventSpan(ctx sdk.Context, k Keeper, msg types.MsgProposeSpan, sideTxResult abci.SideTxResultType) sdk.Result {
	// Skip handler if span is not approved
	if sideTxResult != abci.SideTxResultType_Yes {
		k.Logger(ctx).Debug("Skipping new span since side-tx didn't get yes votes")
		return common.ErrSideTxValidation(k.Codespace()).Result()
	}

	// check for replay
	if k.HasSpan(ctx, msg.ID) {
		k.Logger(ctx).Debug("Skipping new span as it's already processed")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	k.Logger(ctx).Debug("Persisting span state", "sideTxResult", sideTxResult)

	// freeze for new span
	err := k.FreezeSet(ctx, msg.ID, msg.StartBlock, msg.EndBlock, msg.ChainID, msg.Seed)
	if err != nil {
		k.Logger(ctx).Error("Unable to freeze validator set for span", "Error", err)
		return common.ErrUnableToFreezeValSet(k.Codespace()).Result()
	}

	// TX bytes
	txBytes := ctx.TxBytes()
	hash := tmTypes.Tx(txBytes).Hash()

	// add events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeProposeSpan,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                  // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),                // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmTypes.BytesToHeimdallHash(hash).Hex()), // tx hash
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()),             // result
			sdk.NewAttribute(types.AttributeKeySpanID, strconv.FormatUint(msg.ID, 10)),
			sdk.NewAttribute(types.AttributeKeySpanStartBlock, strconv.FormatUint(msg.StartBlock, 10)),
			sdk.NewAttribute(types.AttributeKeySpanEndBlock, strconv.FormatUint(msg.EndBlock, 10)),
		),
	})

	// draft result with events
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
