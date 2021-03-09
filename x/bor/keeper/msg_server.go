package keeper

import (
	"context"
	"strconv"

	hmCommon "github.com/maticnetwork/heimdall/common"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/x/bor/types"
)

type msgServer struct {
	Keeper
}

// NewMsgServerImpl returns an implementation of the bank MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

var _ types.MsgServer = msgServer{}

func (m msgServer) PostSendProposeSpanTx(goCtx context.Context, msg *types.MsgProposeSpan) (*types.MsgProposeSpanResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	m.Keeper.Logger(ctx).Debug("✅ Validating proposed span msg",
		"spanId", msg.SpanId,
		"startBlock", msg.StartBlock,
		"endBlock", msg.EndBlock,
		"seed", msg.Seed,
	)

	// chainManager params
	params := m.Keeper.chainKeeper.GetParams(ctx)
	chainParams := params.ChainParams

	// check chain id
	if chainParams.BorChainID != msg.BorChainId {
		m.Keeper.Logger(ctx).Error("Invalid Bor chain id", "msgChainID", msg.BorChainId)
		return nil, hmCommon.ErrInvalidBorChainID
	}

	// check if last span is up or if greater diff than threshold is found between validator set
	lastSpan, err := m.Keeper.GetLastSpan(ctx)
	if err != nil {
		m.Keeper.Logger(ctx).Error("Unable to fetch last span", "Error", err)
		return nil, hmCommon.ErrSpanNotFound
	}

	// Validate span continuity
	if lastSpan.ID+1 != msg.SpanId || msg.StartBlock != lastSpan.EndBlock+1 || msg.EndBlock < msg.StartBlock {
		m.Keeper.Logger(ctx).Error("Blocks not in continuity",
			"lastSpanId", lastSpan.ID,
			"spanId", msg.SpanId,
			"lastSpanStartBlock", lastSpan.StartBlock,
			"lastSpanEndBlock", lastSpan.EndBlock,
			"spanStartBlock", msg.StartBlock,
			"spanEndBlock", msg.EndBlock,
		)
		return nil, hmCommon.ErrSpanNotInCountinuity
	}

	// Validate Span duration
	spanDuration := m.Keeper.GetParams(ctx).SpanDuration
	if spanDuration != (msg.EndBlock - msg.StartBlock + 1) {
		m.Keeper.Logger(ctx).Error("Span duration of proposed span is wrong",
			"proposedSpanDuration", msg.EndBlock-msg.StartBlock+1,
			"paramsSpanDuration", spanDuration,
		)
		return nil, hmCommon.ErrInvalidSpanDuration
	}

	// add events
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeProposeSpan,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(types.AttributeKeySpanID, strconv.FormatUint(msg.SpanId, 10)),
			sdk.NewAttribute(types.AttributeKeySpanStartBlock, strconv.FormatUint(msg.StartBlock, 10)),
			sdk.NewAttribute(types.AttributeKeySpanEndBlock, strconv.FormatUint(msg.EndBlock, 10)),
		),
	})

	return &types.MsgProposeSpanResponse{}, nil
}
