package slashing

import (
	"fmt"

	abci "github.com/tendermint/tendermint/abci/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	// tmtypes "github.com/tendermint/tendermint/types"
)

func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {
	// BeginBlocker iterates through and handles any newly discovered evidence of
	// misbehavior submitted by Tendermint. Currently, only equivocation is handled.
	// for _, tmEvidence := range req.ByzantineValidators {
	// 	switch tmEvidence.Type {
	// 	case tmtypes.ABCIEvidenceTypeDuplicateVote:
	// 		evidence := ConvertDuplicateVoteEvidence(tmEvidence)
	// 		k.HandleDoubleSign(ctx, evidence.(Equivocation))

	// 	default:
	// 		k.Logger(ctx).Error(fmt.Sprintf("ignored unknown evidence type: %s", tmEvidence.Type))
	// 	}
	// }

	// Iterate over all the validators which *should* have signed this block
	// store whether or not they have actually signed it and slash/unbond any
	// which have missed too many blocks in a row (downtime slashing)
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		k.HandleValidatorSignature(ctx, voteInfo.Validator.Address, voteInfo.Validator.Power, voteInfo.SignedLastBlock)
		fmt.Println(voteInfo)
	}
}
