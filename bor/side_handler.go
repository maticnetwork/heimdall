package bor

import (
	"bytes"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/heimdall/bor/types"

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
		case types.MsgProposeSpan,
			types.MsgProposeSpanV2:
			return SideHandleMsgSpan(ctx, k, msg, contractCaller)
		case types.MsgBackfillSpans:
			return SideHandleMsgBackfillSpans(ctx, k, msg, contractCaller)
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
		case types.MsgProposeSpan,
			types.MsgProposeSpanV2:
			return PostHandleMsgEventSpan(ctx, k, msg, sideTxResult)
		case types.MsgBackfillSpans:
			return PostHandleMsgBackfillSpans(ctx, k, msg, sideTxResult)
		default:
			errMsg := "Unrecognized Span Msg type: %s" + msg.Type()
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// SideHandleMsgSpan validates external calls required for processing proposed span
func SideHandleMsgSpan(ctx sdk.Context, k Keeper, msg sdk.Msg, contractCaller helper.IContractCaller) (result abci.ResponseDeliverSideTx) {
	var proposeMsg types.MsgProposeSpanV2
	switch msg := msg.(type) {
	case types.MsgProposeSpan:
		if ctx.BlockHeight() >= helper.GetDanelawHeight() {
			k.Logger(ctx).Error("Msg span is not allowed after Danelaw hardfork height")
			return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
		}
		proposeMsg = types.MsgProposeSpanV2{
			ID:         msg.ID,
			Proposer:   msg.Proposer,
			StartBlock: msg.StartBlock,
			EndBlock:   msg.EndBlock,
			ChainID:    msg.ChainID,
			Seed:       msg.Seed,
		}
	case types.MsgProposeSpanV2:
		if ctx.BlockHeight() < helper.GetDanelawHeight() {
			k.Logger(ctx).Error("Msg span v2 is not allowed before Danelaw hardfork height")
			return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
		}
		proposeMsg = msg
	}

	k.Logger(ctx).Debug("✅ Validating External call for span msg",
		"msgSeed", proposeMsg.Seed.String(),
	)

	// calculate next span seed locally
	seed, seedAuthor, err := k.GetNextSpanSeed(ctx, proposeMsg.ID)
	if err != nil {
		k.Logger(ctx).Error("Error fetching next span seed from mainchain")
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
	}

	// check if span seed matches or not.
	if !bytes.Equal(proposeMsg.Seed.Bytes(), seed.Bytes()) {
		k.Logger(ctx).Error(
			"Span Seed does not match",
			"proposer", proposeMsg.Proposer.String(),
			"chainID", proposeMsg.ChainID,
			"msgSeed", proposeMsg.Seed.String(),
			"mainchainSeed", seed.String(),
		)

		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
	}

	if ctx.BlockHeight() >= helper.GetDanelawHeight() {
		// check if span seed author matches or not.
		if !bytes.Equal(proposeMsg.SeedAuthor.Bytes(), seedAuthor.Bytes()) {
			k.Logger(ctx).Error(
				"Span Seed Author does not match",
				"proposer", proposeMsg.Proposer.String(),
				"chainID", proposeMsg.ChainID,
				"msgSeed", proposeMsg.Seed.String(),
				"msgSeedAuthor", proposeMsg.SeedAuthor.String(),
				"mainchainSeedAuthor", seedAuthor.String(),
				"mainchainSeed", seed.String(),
			)

			return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
		}
	}

	// fetch current child block
	childBlock, err := contractCaller.GetMaticChainBlock(nil)
	if err != nil {
		k.Logger(ctx).Error("Error fetching current child block", "error", err)
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
	}

	lastSpan, err := k.GetLastSpan(ctx)
	if err != nil {
		k.Logger(ctx).Error("Error fetching last span", "error", err)
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
	}

	var latestMilestoneEndBlock uint64
	latestMilestone, err := k.checkpointKeeper.GetLastMilestone(ctx)
	if err == nil {
		latestMilestoneEndBlock = latestMilestone.EndBlock
	} else {
		k.Logger(ctx).Error("error fetching latest milestone", "error", err)
	}

	currentBlock := childBlock.Number.Uint64()

	maxBlockNumber := max(latestMilestoneEndBlock, currentBlock)

	if types.IsBlockCloseToSpanEnd(maxBlockNumber, lastSpan.EndBlock) {
		k.Logger(ctx).Debug("current block is close to span end", "currentBlock", currentBlock, "lastSpanEndBlock", lastSpan.EndBlock)
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeTooCloseToSpanEnd)
	}

	// If we are past end of the last span, we need to backfill before proposing a new span
	if proposeMsg.StartBlock <= maxBlockNumber {
		k.Logger(ctx).Error("span is already in the past",
			"currentBlock", currentBlock,
			"msgStartBlock", proposeMsg.StartBlock,
			"msgEndBlock", proposeMsg.EndBlock,
			"latestMilestoneEndBlock", latestMilestoneEndBlock,
			"lastSpanEndBlock", lastSpan.EndBlock,
		)
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeSpanInThePast)
	}

	// check if span proposed is in-turn or not
	if !(lastSpan.StartBlock <= currentBlock && currentBlock <= lastSpan.EndBlock) {
		k.Logger(ctx).Error(
			"Span proposed is not in-turn",
			"currentChildBlock", currentBlock,
			"msgStartblock", proposeMsg.StartBlock,
			"msgEndBlock", proposeMsg.EndBlock,
		)

		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
	}

	k.Logger(ctx).Debug("✅ Successfully validated External call for span msg")

	result.Result = abci.SideTxResultType_Yes

	return
}

func SideHandleMsgBackfillSpans(ctx sdk.Context, k Keeper, msg types.MsgBackfillSpans, contractCaller helper.IContractCaller) (result abci.ResponseDeliverSideTx) {

	k.Logger(ctx).Debug("✅ validating external call for fill missing spans msg",
		"proposer", msg.Proposer,
		"chainId", msg.ChainID,
		"latestSpanId", msg.LatestBorSpanID,
		"latestHeimdallSpan", msg.LatestBorSpanID,
	)

	latestSpan, err := k.GetLastSpan(ctx)
	if err != nil {
		k.Logger(ctx).Error("failed to get latest span", "error", err)
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
	}

	borLastUsedSpanID, err := contractCaller.GetStartBlockHeimdallSpanID(ctx.Context(), latestSpan.EndBlock+1)
	if err != nil {
		k.Logger(ctx).Error("failed to get last used bor span id", "error", err)
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
	}
	if borLastUsedSpanID == 0 {
		k.Logger(ctx).Error("last used bor span id is 0")
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
	}

	if borLastUsedSpanID != msg.LatestSpanID {
		k.Logger(ctx).Error("last used bor span id does not match",
			"expected", borLastUsedSpanID,
			"got", msg.LatestSpanID,
		)
		return hmCommon.ErrorSideTx(k.Codespace(), hmCommon.CodeInvalidMsg)
	}

	k.Logger(ctx).Debug("✅ successfully validated external call for fill missing spans msg")

	result.Result = abci.SideTxResultType_Yes

	return
}

// PostHandleMsgEventSpan handles state persisting span msg
func PostHandleMsgEventSpan(ctx sdk.Context, k Keeper, msg sdk.Msg, sideTxResult abci.SideTxResultType) sdk.Result {
	logger := k.Logger(ctx)

	// Skip handler if span is not approved
	if sideTxResult != abci.SideTxResultType_Yes {
		logger.Debug("Skipping new span since side-tx didn't get yes votes")
		return hmCommon.ErrSideTxValidation(k.Codespace()).Result()
	}

	// check if msg is of type MsgProposeSpanV2
	var proposeMsg types.MsgProposeSpanV2
	switch msg := msg.(type) {
	case types.MsgProposeSpan:
		if ctx.BlockHeight() >= helper.GetDanelawHeight() {
			k.Logger(ctx).Error("Msg span is not allowed after Danelaw hardfork height")
			return hmCommon.ErrSideTxValidation(k.Codespace()).Result()
		}
		proposeMsg = types.MsgProposeSpanV2{
			ID:         msg.ID,
			Proposer:   msg.Proposer,
			StartBlock: msg.StartBlock,
			EndBlock:   msg.EndBlock,
			ChainID:    msg.ChainID,
			Seed:       msg.Seed,
		}
	case types.MsgProposeSpanV2:
		if ctx.BlockHeight() < helper.GetDanelawHeight() {
			k.Logger(ctx).Error("Msg span v2 is not allowed before Danelaw hardfork height")
			return hmCommon.ErrSideTxValidation(k.Codespace()).Result()
		}
		proposeMsg = msg
	}

	// check for replay
	if k.HasSpan(ctx, proposeMsg.ID) {
		logger.Debug("Skipping new span as it's already processed")
		return hmCommon.ErrOldTx(k.Codespace()).Result()
	}

	logger.Debug("Persisting span state",
		"sideTxResult", sideTxResult,
		"proposer", proposeMsg.Proposer.String(),
		"spanId", proposeMsg.ID,
		"startBlock", proposeMsg.StartBlock,
		"endBlock", proposeMsg.EndBlock,
		"seed", proposeMsg.Seed.String(),
	)

	if ctx.BlockHeight() >= helper.GetJorvikHeight() {
		var seedSpanID uint64
		if proposeMsg.ID < 2 {
			seedSpanID = proposeMsg.ID - 1
		} else {
			seedSpanID = proposeMsg.ID - 2
		}

		lastSpan, err := k.GetSpan(ctx, seedSpanID)
		if err != nil {
			logger.Error("Unable to get last span", "Error", err)
			return hmCommon.ErrUnableToGetSpan(k.Codespace()).Result()
		}

		var producer *ethCommon.Address

		if ctx.BlockHeight() < helper.GetDanelawHeight() {
			// store the seed producer
			_, producer, err = k.getBorBlockForSpanSeed(ctx, lastSpan, proposeMsg.ID)
			if err != nil {
				logger.Error("Unable to get seed producer", "Error", err)
				return hmCommon.ErrUnableToGetSeed(k.Codespace()).Result()
			}
		} else {
			producer = &proposeMsg.SeedAuthor
		}

		if err = k.StoreSeedProducer(ctx, proposeMsg.ID, producer); err != nil {
			logger.Error("Unable to store seed producer", "Error", err)
			return hmCommon.ErrUnableToStoreSeedProducer(k.Codespace()).Result()
		}
	}

	// freeze for new span
	err := k.FreezeSet(ctx, proposeMsg.ID, proposeMsg.StartBlock, proposeMsg.EndBlock, proposeMsg.ChainID, proposeMsg.Seed)
	if err != nil {
		k.Logger(ctx).Error("Unable to freeze validator set for span", "Error", err)
		return hmCommon.ErrUnableToFreezeValSet(k.Codespace()).Result()
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
			sdk.NewAttribute(types.AttributeKeySpanID, strconv.FormatUint(proposeMsg.ID, 10)),
			sdk.NewAttribute(types.AttributeKeySpanStartBlock, strconv.FormatUint(proposeMsg.StartBlock, 10)),
			sdk.NewAttribute(types.AttributeKeySpanEndBlock, strconv.FormatUint(proposeMsg.EndBlock, 10)),
		),
	})

	// draft result with events
	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}

func PostHandleMsgBackfillSpans(ctx sdk.Context, k Keeper, msg types.MsgBackfillSpans, sideTxResult abci.SideTxResultType) sdk.Result {

	if sideTxResult != abci.SideTxResultType_Yes {
		k.Logger(ctx).Debug("skipping new span since side-tx didn't get yes votes")
		return hmCommon.ErrSideTxValidation(k.Codespace()).Result()
	}

	latestMilestone, err := k.checkpointKeeper.GetLastMilestone(ctx)
	if err != nil {
		k.Logger(ctx).Error("failed to get latest milestone", "error", err)
		return hmCommon.ErrUnableToGetLastMilestone(k.Codespace()).Result()
	}

	if latestMilestone == nil {
		k.Logger(ctx).Error("latest milestone is nil")
		return hmCommon.ErrLatestMilestoneNotFound(k.Codespace()).Result()
	}

	latestSpan, err := k.GetSpan(ctx, msg.LatestSpanID)
	if err != nil {
		k.Logger(ctx).Error("failed to get latest span", "error", err)
		return hmCommon.ErrUnableToGetSpan(k.Codespace()).Result()
	}

	if latestSpan == nil {
		k.Logger(ctx).Error("latest span is nil", "latestSpanId", msg.LatestSpanID)
		return hmCommon.ErrSpanNotFound(k.Codespace()).Result()
	}

	borSpans := types.GenerateBorCommittedSpans(latestMilestone.EndBlock, latestSpan)
	spansOverlap := 0
	for i := range borSpans {
		if _, err := k.GetSpan(ctx, borSpans[i].ID); err == nil {
			spansOverlap++
		}
		if spansOverlap > 1 {
			k.Logger(ctx).Error("more than one span overlap detected", "span id", borSpans[i].ID)
			return hmCommon.ErrSpanOverlap(k.Codespace(), borSpans[i].StartBlock, borSpans[i].ID).Result()
		}
		if err = k.AddNewSpan(ctx, borSpans[i]); err != nil {
			k.Logger(ctx).Error("Unable to store spans", "error", err)
			return hmCommon.ErrUnableToAddSpan(k.Codespace()).Result()
		}
	}

	txBytes := ctx.TxBytes()
	hash := tmTypes.Tx(txBytes).Hash()

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeProposeSpan,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmTypes.BytesToHeimdallHash(hash).Hex()),
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()),
			sdk.NewAttribute(types.AttributesKeyLatestSpanId, strconv.FormatUint(msg.LatestSpanID, 10)),
			sdk.NewAttribute(types.AttributesKeyLatestBorSpanId, strconv.FormatUint(borSpans[0].ID, 10)),
		),
	})

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
