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
		k.Logger(ctx).Debug("slashing is not enabled")
		return
	}

	// BeginBlocker iterates through and handles any newly discovered evidence of
	// misbehavior submitted by Tendermint. Currently, only equivocation is handled.
	for _, tmEvidence := range req.ByzantineValidators {
		switch tmEvidence.Type {
		case tmtypes.ABCIEvidenceTypeDuplicateVote:
			evidence := types.ConvertDuplicateVoteEvidence(tmEvidence)
			k.HandleDoubleSign(ctx, evidence.(types.Equivocation))

		default:
			k.Logger(ctx).Error(fmt.Sprintf("ignored unknown evidence type: %s", tmEvidence.Type))
		}
	}
	// TODO - slashing remove below for loop. only for testing purpose
	/* 	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		k.HandleValidatorSignature(ctx, voteInfo.Validator.Address, voteInfo.Validator.Power, voteInfo.SignedLastBlock)
		evidence := types.Equivocation{
			ConsensusAddress: voteInfo.Validator.GetAddress(),
			Height:           ctx.BlockHeight(),
			Time:             ctx.BlockTime(),
			Power:            voteInfo.Validator.GetPower(),
		}

		k.Logger(ctx).Debug("Sending fake evidence", "validatorAddr", evidence.GetConsensusAddress())
		k.HandleDoubleSign(ctx, evidence)
	} */

	// Iterate over all the validators which *should* have signed this block
	// store whether or not they have actually signed it and slash/unbond any
	// which have missed too many blocks in a row (downtime slashing)
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		k.HandleValidatorSignature(ctx, voteInfo.Validator.Address, voteInfo.Validator.Power, voteInfo.SignedLastBlock)
		// TODO - slashing remove false. only for testing purpose
		// k.HandleValidatorSignature(ctx, voteInfo.Validator.Address, voteInfo.Validator.Power, false)
	}
}
