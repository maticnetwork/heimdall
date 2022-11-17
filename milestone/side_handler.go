package milestone

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/milestone/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// NewSideTxHandler returns a side handler for "milestone" type messages.
func NewSideTxHandler(k Keeper, contractCaller helper.IContractCaller) hmTypes.SideTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg) abci.ResponseDeliverSideTx {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgMilestone:
			return SideHandleMsgMilestone(ctx, k, msg, contractCaller)
		default:
			return abci.ResponseDeliverSideTx{
				Code: uint32(sdk.CodeUnknownRequest),
			}
		}
	}
}

// SideHandleMsgMilestone handles MsgMilestone message for external call
func SideHandleMsgMilestone(ctx sdk.Context, k Keeper, msg types.MsgMilestone, contractCaller helper.IContractCaller) (result abci.ResponseDeliverSideTx) {
	// get params
	milestoneLength := helper.MilestoneLength

	// logger
	logger := k.Logger(ctx)

	// validate milestone
	count := k.GetCount(ctx)
	lastMilestone, err := k.GetLastMilestone(ctx)

	if count != uint64(0) && err != nil {
		logger.Error("Error while receiving the last milestone in the side handler")
		return common.ErrorSideTx(k.Codespace(), common.CodeInvalidBlockInput)

	}

	if count != uint64(0) && msg.StartBlock != lastMilestone.EndBlock+1 {
		logger.Error("Milestone is not in continuity to last stored milestone",
			"startBlock", msg.StartBlock,
			"endBlock", msg.EndBlock,
			"rootHash", msg.RootHash,
			"milestoneId", msg.MilestoneID,
			"error", err,
		)
		return common.ErrorSideTx(k.Codespace(), common.CodeInvalidBlockInput)
	}

	validMilestone, err := types.ValidateMilestone(msg.StartBlock, msg.EndBlock, msg.RootHash, msg.MilestoneID, contractCaller, milestoneLength)
	if err != nil {
		logger.Error("Error validating milestone",
			"startBlock", msg.StartBlock,
			"endBlock", msg.EndBlock,
			"rootHash", msg.RootHash,
			"milestoneId", msg.MilestoneID,
			"error", err,
		)

	} else if validMilestone {
		// vote `yes` if milestone is valid
		result.Result = abci.SideTxResultType_Yes
		return
	}

	logger.Error(
		"RootHash is not valid",
		"startBlock", msg.StartBlock,
		"endBlock", msg.EndBlock,
		"rootHash", msg.RootHash,
		"milestoneId", msg.MilestoneID,
	)

	return common.ErrorSideTx(k.Codespace(), common.CodeInvalidBlockInput)
}

//
// Tx handler
//

// NewPostTxHandler returns a side handler for "bank" type messages.
func NewPostTxHandler(k Keeper, contractCaller helper.IContractCaller) hmTypes.PostTxHandler {
	return func(ctx sdk.Context, msg sdk.Msg, sideTxResult abci.SideTxResultType) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgMilestone:
			return PostHandleMsgMilestone(ctx, k, msg, sideTxResult)
		default:
			return sdk.ErrUnknownRequest("Unrecognized milestone Msg type").Result()
		}
	}
}

// PostHandleMsgMilestone handles msg milestone
func PostHandleMsgMilestone(ctx sdk.Context, k Keeper, msg types.MsgMilestone, sideTxResult abci.SideTxResultType) sdk.Result {
	logger := k.Logger(ctx)
	timeStamp := uint64(ctx.BlockTime().Unix())

	// TX bytes
	txBytes := ctx.TxBytes()
	hash := tmTypes.Tx(txBytes).Hash()

	// Emit event for milestone
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventTypeMilestone,
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),                                  // action
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),                // module name
			sdk.NewAttribute(hmTypes.AttributeKeyTxHash, hmTypes.BytesToHeimdallHash(hash).Hex()), // tx hash
			sdk.NewAttribute(hmTypes.AttributeKeySideTxResult, sideTxResult.String()),             // result
			sdk.NewAttribute(types.AttributeKeyProposer, msg.Proposer.String()),
			sdk.NewAttribute(types.AttributeKeyStartBlock, strconv.FormatUint(msg.StartBlock, 10)),
			sdk.NewAttribute(types.AttributeKeyEndBlock, strconv.FormatUint(msg.EndBlock, 10)),
			sdk.NewAttribute(types.AttributeKeyRootHash, msg.RootHash.String()),
			sdk.NewAttribute(types.AttributeKeyMilestoneID, msg.MilestoneID),
		),
	})

	// Skip handler if milestone is not approved
	if sideTxResult != abci.SideTxResultType_Yes {
		logger.Debug("Skipping new milestone since side-tx didn't get yes votes", "startBlock", msg.StartBlock, "endBlock", msg.EndBlock, "rootHash", msg.RootHash, "milestoneId", msg.MilestoneID)
		k.SetNoAckMilestone(ctx, msg.MilestoneID)

		return sdk.Result{
			Events: ctx.EventManager().Events(),
		}

	}

	if lastMilestone, err := k.GetLastMilestone(ctx); err == nil { // fetch last milestone from store
		// make sure new milestoen is after tip
		if lastMilestone.EndBlock > msg.StartBlock {
			logger.Error(" already exists",
				"currentTip", lastMilestone.EndBlock,
				"startBlock", msg.StartBlock,
			)

			k.SetNoAckMilestone(ctx, msg.MilestoneID)

			return sdk.Result{
				Events: ctx.EventManager().Events(),
			}
		}

		// check if new milestone's start block start from current tip
		if lastMilestone.EndBlock+1 != msg.StartBlock {
			logger.Error("milestone not in countinuity",
				"currentTip", lastMilestone.EndBlock,
				"startBlock", msg.StartBlock)

			k.SetNoAckMilestone(ctx, msg.MilestoneID)

			return sdk.Result{
				Events: ctx.EventManager().Events(),
			}
		}

	} else if err != nil && msg.StartBlock != 0 {
		logger.Error("First milestone to start from", "block", 0, "Error", err)

		k.SetNoAckMilestone(ctx, msg.MilestoneID)

		return sdk.Result{
			Events: ctx.EventManager().Events(),
		}

	}

	if err := k.AddMilestone(ctx, hmTypes.Milestone{ // Save milestone to buffer store
		StartBlock:  msg.StartBlock, //Add milestone to store with root hash
		EndBlock:    msg.EndBlock,
		RootHash:    msg.RootHash,
		Proposer:    msg.Proposer,
		BorChainID:  msg.BorChainID,
		MilestoneID: msg.MilestoneID,
		TimeStamp:   timeStamp,
	}); err != nil {
		k.SetNoAckMilestone(ctx, msg.MilestoneID)
		logger.Error("Failed to set milestone ", "Error", err)
	}

	return sdk.Result{
		Events: ctx.EventManager().Events(),
	}
}
