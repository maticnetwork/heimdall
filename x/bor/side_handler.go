package bor

import (
	"bytes"
	"fmt"
	"strconv"

	tmTypes "github.com/tendermint/tendermint/types"

	hmCommon "github.com/maticnetwork/heimdall/common"

	"github.com/maticnetwork/heimdall/types/common"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/x/bor/keeper"
	"github.com/maticnetwork/heimdall/x/bor/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmprototypes "github.com/tendermint/tendermint/proto/tendermint/types"
)

// NewSideTxHandler returns a side handler for "span" type messages.
func NewSideTxHandler(k keeper.Keeper, contractCaller helper.IContractCaller) hmTypes.SideTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgProposeSpan:
			return SideHandleMsgSpan(ctx, k, *msg, contractCaller)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			fmt.Println(errMsg)
			return abci.ResponseDeliverSideTx{
				Code: uint32(6), // TODO should be changed like `sdk.CodeUnknownRequest `
			}
		}
	}
}

// NewPostTxHandler returns a side handler for "span" type messages.
func NewPostTxHandler(k keeper.Keeper) hmTypes.PostTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg, sideTxResult tmprototypes.SideTxResultType) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case *types.MsgProposeSpan:
			return PostHandleMsgEventSpan(ctx, k, *msg, sideTxResult)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// SideHandleMsgSpan validates external calls required for processing proposed span
func SideHandleMsgSpan(ctx sdk.Context, k keeper.Keeper, msg types.MsgProposeSpan, contractCaller helper.IContractCaller) (result abci.ResponseDeliverSideTx) {
	k.Logger(ctx).Debug("✅ Validating External call for span msg",
		"msgSeed", msg.Seed,
	)
	// calculate next span seed locally
	nextSpanSeed, err := k.GetNextSpanSeed(ctx, contractCaller)
	if err != nil {
		k.Logger(ctx).Error("Error fetching next span seed from mainchain")
		return hmCommon.ErrorSideTx(hmCommon.ErrInvalidMsg)
	}

	// check if span seed matches or not.

	if !bytes.Equal(common.HexToHeimdallHash(msg.Seed).Bytes(), nextSpanSeed.Bytes()) {
		k.Logger(ctx).Error(
			"Span Seed does not match",
			"msgSeed", msg.Seed,
			"mainchainSeed", nextSpanSeed.String(),
		)
		return hmCommon.ErrorSideTx(hmCommon.ErrInvalidMsg)
	}

	// fetch current child block
	childBlock, err := contractCaller.GetMaticChainBlock(nil)
	if err != nil {
		k.Logger(ctx).Error("Error fetching current child block", "error", err)
		return hmCommon.ErrorSideTx(hmCommon.ErrInvalidMsg)
	}

	lastSpan, err := k.GetLastSpan(ctx)
	if err != nil {
		k.Logger(ctx).Error("Error fetching last span", "error", err)
		return hmCommon.ErrorSideTx(hmCommon.ErrInvalidMsg)
	}

	currentBlock := childBlock.Number.Uint64()
	// check if span proposed is in-turn or not
	if !(lastSpan.StartBlock <= currentBlock && currentBlock <= lastSpan.EndBlock) {
		k.Logger(ctx).Error(
			"Span proposed is not in-turn",
			"currentChildBlock", currentBlock,
			"msgStartBlock", msg.StartBlock,
			"msgEndBlock", msg.EndBlock,
		)
		return hmCommon.ErrorSideTx(hmCommon.ErrInvalidMsg)
	}

	k.Logger(ctx).Debug("✅ Successfully validated External call for span msg")
	result.Result = tmprototypes.SideTxResultType_YES
	return
}

// PostHandleMsgEventSpan handles state persisting span msg
func PostHandleMsgEventSpan(ctx sdk.Context, k keeper.Keeper, msg types.MsgProposeSpan, sideTxResult tmprototypes.SideTxResultType) (*sdk.Result, error) {
	// Skip handler if span is not approved
	if sideTxResult != tmprototypes.SideTxResultType_YES {
		k.Logger(ctx).Debug("Skipping new validator-join since side-tx didn't get yes votes")
		return nil, hmCommon.ErrSideTxValidation
	}

	// check for replay
	if found := k.HasSpan(ctx, msg.SpanId); found {
		k.Logger(ctx).Debug("Skipping new span as it's already processed")
		return nil, hmCommon.ErrOldTx
	}

	k.Logger(ctx).Debug("Persisting span state", "sideTxResult", sideTxResult)

	// freeze for new span
	err := k.FreezeSet(ctx, msg.SpanId, msg.StartBlock, msg.EndBlock, msg.ChainId, common.HexToHeimdallHash(msg.Seed).EthHash())
	if err != nil {
		k.Logger(ctx).Error("Unable to freeze validator set for span", "Error", err)
		return nil, hmCommon.ErrUnableToFreezeValSet
	}

	// TX bytes
	txBytes := ctx.TxBytes()
	hash := tmTypes.Tx(txBytes).Hash()

	// add events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeProposeSpan,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                 // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),               // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, common.BytesToHeimdallHash(hash).Hex()), // tx hash
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()),            // result
			sdk.NewAttribute(types.AttributeKeySpanID, strconv.FormatUint(msg.SpanId, 10)),
			sdk.NewAttribute(types.AttributeKeySpanStartBlock, strconv.FormatUint(msg.StartBlock, 10)),
			sdk.NewAttribute(types.AttributeKeySpanEndBlock, strconv.FormatUint(msg.EndBlock, 10)),
		),
	})

	// draft result with events
	return &sdk.Result{
		Events: ctx.EventManager().ABCIEvents(),
	}, nil
}
