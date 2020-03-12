package gov

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/gov/types"
)

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, keeper Keeper) {
	logger := keeper.Logger(ctx)

	// delete inactive proposal from store and its deposits
	keeper.IterateInactiveProposalsQueue(ctx, ctx.BlockHeader().Time, func(proposal types.Proposal) bool {
		keeper.RefundDeposits(ctx, proposal.ProposalID)
		keeper.DeleteProposal(ctx, proposal.ProposalID)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeInactiveProposal,
				sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.ProposalID)),
				sdk.NewAttribute(types.AttributeKeyProposalResult, types.AttributeValueProposalDropped),
			),
		)

		logger.Info(
			fmt.Sprintf("proposal %d (%s) didn't meet minimum deposit of %s (had only %s); deleted",
				proposal.ProposalID,
				proposal.GetTitle(),
				keeper.GetDepositParams(ctx).MinDeposit,
				proposal.TotalDeposit,
			),
		)
		return false
	})

	// fetch active proposals whose voting periods have ended (are passed the block time)
	keeper.IterateActiveProposalsQueue(ctx, ctx.BlockHeader().Time, func(proposal types.Proposal) bool {
		var tagValue, logMsg string

		passes, _, tallyResults := tally(ctx, keeper, proposal)

		// refund deposits
		keeper.RefundDeposits(ctx, proposal.ProposalID)

		if passes {
			handler := keeper.router.GetRoute(proposal.ProposalRoute())
			cacheCtx, writeCache := ctx.CacheContext()

			// The proposal handler may execute state mutating logic depending
			// on the proposal content. If the handler fails, no state mutation
			// is written and the error message is logged.
			err := handler(cacheCtx, proposal.Content)
			if err == nil {
				proposal.Status = types.StatusPassed
				tagValue = types.AttributeValueProposalPassed
				logMsg = "passed"

				// write state to the underlying multi-store
				writeCache()
			} else {
				proposal.Status = types.StatusFailed
				tagValue = types.AttributeValueProposalFailed
				logMsg = fmt.Sprintf("passed, but failed on execution: %s", err.ABCILog())
			}
		} else {
			proposal.Status = types.StatusRejected
			tagValue = types.AttributeValueProposalRejected
			logMsg = "rejected"
		}

		proposal.FinalTallyResult = tallyResults

		keeper.SetProposal(ctx, proposal)
		keeper.RemoveFromActiveProposalQueue(ctx, proposal.ProposalID, proposal.VotingEndTime)

		logger.Info(
			fmt.Sprintf(
				"proposal %d (%s) tallied; result: %s",
				proposal.ProposalID, proposal.GetTitle(), logMsg,
			),
		)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeActiveProposal,
				sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposal.ProposalID)),
				sdk.NewAttribute(types.AttributeKeyProposalResult, tagValue),
			),
		)
		return false
	})
}
