package bor

import (
	"bytes"
	"sort"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/bor/tags"
	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/types"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

// NewHandler returns a handler for "bor" type messages.
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		common.InitBorLogger(&ctx)
		switch msg := msg.(type) {
		case MsgProposeSpan:
			return HandleMsgProposeSpan(ctx, msg, k, common.BorLogger)
		default:
			return sdk.ErrTxDecode("Invalid message in bor module").Result()
		}
	}
}

// HandleMsgProposeSpan handles proposeSpan msg
func HandleMsgProposeSpan(ctx sdk.Context, msg MsgProposeSpan, k Keeper, logger tmlog.Logger) sdk.Result {
	logger.Debug("Proposing span", "TxData", msg)

	// check if last span is up or if greater diff than threshold is found between validator set
	lastSpan, err := k.GetLastSpan(ctx)
	if err != nil {
		logger.Error("Unable to fetch last span", "Error", err)
		return common.ErrSpanNotInCountinuity(k.Codespace).Result()
	}

	// check if lastStart + 1 =  newStart
	if lastSpan.StartBlock > msg.StartBlock {
		common.BorLogger.Error("Blocks not in countinuity ", "LastStartBlock", lastSpan.StartBlock, "MsgStartBlock", msg.StartBlock)
		return common.ErrSpanNotInCountinuity(k.Codespace).Result()
	}

	// freeze for new span
	err = k.FreezeSet(ctx, msg.StartBlock, msg.ChainID)
	if err != nil {
		logger.Error("Unable to freeze validator set for span", "Error", err)
		return common.ErrSpanNotInCountinuity(k.Codespace).Result()
	}

	currentValidators := k.sk.GetCurrentValidators(ctx)

	lastSpan, err = k.GetLastSpan(ctx)
	if err != nil {
		logger.Error("Unable to fetch last span", "Error", err)
		return common.ErrSpanNotFound(k.Codespace).Result()
	}
	resTags := sdk.NewTags(
		tags.NewSpanProposed, []byte(strconv.FormatUint(uint64(msg.StartBlock), 10)),
	)

	// TODO add check for duration

	result := sortAndCompare(types.ValToMinVal(currentValidators), types.ValToMinVal(lastSpan.SelectedProducers), msg, k.Codespace)
	result.Tags = resTags
	return result
}

func sortAndCompare(allVals []types.MinimalVal, selectedVals []types.MinimalVal, msg MsgProposeSpan, codespace sdk.CodespaceType) sdk.Result {
	if len(allVals) != len(msg.Validators) || len(selectedVals) != len(msg.SelectedProducers) {
		return common.ErrProducerMisMatch(codespace).Result()
	}

	sortedAddVals := sortByAddress(allVals)
	sortedMsgValidators := sortByAddress(msg.Validators)

	for i := range sortedMsgValidators {
		if !bytes.Equal(sortedMsgValidators[i].Signer.Bytes(), sortedAddVals[i].Signer.Bytes()) || sortedMsgValidators[i].Power != sortedAddVals[i].Power {
			common.BorLogger.Error("Validator Set does not match", "inputValSet", sortedMsgValidators, "storedValSet", sortedAddVals)
			return common.ErrValSetMisMatch(codespace).Result()
		}
	}
	sortedSelectedVals := sortByAddress(selectedVals)
	sortedMsgProducers := sortByAddress(msg.SelectedProducers)
	for i := range selectedVals {
		if !bytes.Equal(sortedSelectedVals[i].Signer.Bytes(), sortedMsgProducers[i].Signer.Bytes()) {
			common.BorLogger.Error("Producer does not match", "inputProducers", sortedMsgProducers, "storedProducers", sortedSelectedVals)
			return common.ErrProducerMisMatch(codespace).Result()
		}
	}
	return sdk.Result{}
}

func sortByAddress(a []types.MinimalVal) []types.MinimalVal {
	sort.Slice(a, func(i, j int) bool {
		return bytes.Compare(a[i].Signer.Bytes(), a[j].Signer.Bytes()) < 0
	})
	return a
}
