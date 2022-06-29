package slashing

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/slashing/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {

	if !k.GetParams(ctx).EnableSlashing {
		k.Logger(ctx).Debug("slashing is not enabled. To enable, send a proposal via governance")
		return
	}

	// BeginBlocker iterates through and handles any newly discovered evidence of
	// misbehavior submitted by Tendermint. Currently, only equivocation is handled.
	for _, tmEvidence := range req.ByzantineValidators {
		switch tmEvidence.Type {
		case tmtypes.ABCIEvidenceTypeDuplicateVote:
			evidence := types.ConvertDuplicateVoteEvidence(tmEvidence)
			if err := k.HandleDoubleSign(ctx, evidence.(types.Equivocation)); err != nil {
				k.Logger(ctx).Error("Failed to handle double sign", "Error", err)
			}
		default:
			k.Logger(ctx).Error(fmt.Sprintf("ignored unknown evidence type: %s", tmEvidence.Type))
		}
	}

	// Iterate over all the validators which *should* have signed this block
	// store whether or not they have actually signed it and slash/unbond any
	// which have missed too many blocks in a row (downtime slashing)
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		if err := k.HandleValidatorSignature(
			ctx,
			voteInfo.Validator.Address,
			voteInfo.Validator.Power,
			voteInfo.SignedLastBlock,
		); err != nil {
			k.Logger(ctx).Error("Failed to handle validator signature", "Error", err, "address", voteInfo.Validator)
		}
	}
}
