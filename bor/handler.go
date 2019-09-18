package bor

import (
	"bytes"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/bor/tags"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/types"
)

// NewHandler returns a handler for "bor" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgProposeSpan:
			return HandleMsgProposeSpan(ctx, msg, k)
		default:
			return sdk.ErrTxDecode("Invalid message in bor module").Result()
		}
	}
}

// HandleMsgProposeSpan handles proposeSpan msg
func HandleMsgProposeSpan(ctx sdk.Context, msg MsgProposeSpan, k Keeper) sdk.Result {
	k.Logger(ctx).Debug("Proposing span", "TxData", msg)

	// check if last span is up or if greater diff than threshold is found between validator set
	lastSpan, err := k.GetLastSpan(ctx)
	if err != nil {
		k.Logger(ctx).Error("Unable to fetch last span", "Error", err)
		return common.ErrSpanNotFound(k.Codespace()).Result()
	}

	// check all conditions
	if lastSpan.ID+1 != msg.ID || msg.StartBlock < lastSpan.StartBlock || msg.EndBlock < msg.StartBlock {
		k.Logger(ctx).Error("Blocks not in countinuity",
			"lastSpanId", lastSpan.ID,
			"lastSpanStartBlock", lastSpan.StartBlock,
			"spanId", msg.ID,
			"spanStartBlock", msg.StartBlock,
		)
		return common.ErrSpanNotInCountinuity(k.Codespace()).Result()
	}

	// freeze for new span
	err = k.FreezeSet(ctx, msg.ID, msg.StartBlock, msg.ChainID)
	if err != nil {
		k.Logger(ctx).Error("Unable to freeze validator set for span", "Error", err)
		return common.ErrUnableToFreezeValSet(k.Codespace()).Result()
	}

	// get current validators
	currentValidators := k.sk.GetCurrentValidators(ctx)

	// get last span
	lastSpan, err = k.GetLastSpan(ctx)
	if err != nil {
		k.Logger(ctx).Error("Unable to fetch last span", "Error", err)
		return common.ErrSpanNotFound(k.Codespace()).Result()
	}
	k.Logger(ctx).Info("==> Fetched last span", "Span", lastSpan.SelectedProducers)
	// TODO add check for duration
	result, ok := sortAndCompare(types.ValToMinVal(currentValidators), types.ValToMinVal(lastSpan.SelectedProducers), msg, k.Codespace())
	if ok {
		result.Tags = sdk.NewTags(
			tags.Success, []byte("true"),
			tags.BorSyncID, []byte(strconv.FormatUint(uint64(msg.ID), 10)),
			tags.SpanID, []byte(strconv.FormatUint(uint64(msg.ID), 10)),
			tags.SpanStartBlock, []byte(strconv.FormatUint(uint64(msg.StartBlock), 10)),
		)
	}
	return result
}

func sortAndCompare(allVals []types.MinimalVal, selectedVals []types.MinimalVal, msg MsgProposeSpan, codespace sdk.CodespaceType) (sdk.Result, bool) {
	if len(allVals) != len(msg.Validators) || len(selectedVals) != len(msg.SelectedProducers) {
		return common.ErrProducerMisMatch(codespace).Result(), false
	}

	sortedAddVals := types.SortMinimalValByAddress(allVals)
	sortedMsgValidators := types.SortMinimalValByAddress(msg.Validators)
	for i := range sortedMsgValidators {
		if !bytes.Equal(sortedMsgValidators[i].Signer.Bytes(), sortedAddVals[i].Signer.Bytes()) || sortedMsgValidators[i].Power != sortedAddVals[i].Power {
			return common.ErrValSetMisMatch(codespace).Result(), false
		}
	}

	sortedSelectedVals := types.SortMinimalValByAddress(selectedVals)
	sortedMsgProducers := types.SortMinimalValByAddress(msg.SelectedProducers)
	for i := range selectedVals {
		if !bytes.Equal(sortedSelectedVals[i].Signer.Bytes(), sortedMsgProducers[i].Signer.Bytes()) {
			return common.ErrProducerMisMatch(codespace).Result(), false
		}
	}

	return sdk.Result{}, true
}
