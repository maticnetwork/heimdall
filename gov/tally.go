package gov

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/gov/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// validatorGovInfo used for tallying
type validatorGovInfo struct {
	Validator   hmTypes.ValidatorID // id of the validator operator
	VotingPower int64               // voting power
	Vote        types.VoteOption    // Vote of the validator
}

func newValidatorGovInfo(
	validator hmTypes.ValidatorID,
	votingPower int64,
	vote types.VoteOption,
) validatorGovInfo {
	return validatorGovInfo{
		Validator:   validator,
		VotingPower: votingPower,
		Vote:        vote,
	}
}

// TODO: Break into several smaller functions for clarity
func tally(ctx sdk.Context, keeper Keeper, proposal types.Proposal) (passes bool, burnDeposits bool, tallyResults types.TallyResult) {
	results := make(map[types.VoteOption]sdk.Dec)
	results[types.OptionYes] = sdk.ZeroDec()
	results[types.OptionAbstain] = sdk.ZeroDec()
	results[types.OptionNo] = sdk.ZeroDec()
	results[types.OptionNoWithVeto] = sdk.ZeroDec()

	totalBondedTokens := sdk.ZeroDec()
	totalVotingPower := sdk.ZeroDec()
	currValidators := make(map[hmTypes.ValidatorID]validatorGovInfo)

	// fetch all the bonded validators, insert them into currValidators
	keeper.sk.IterateCurrentValidatorsAndApplyFn(ctx, func(validator *hmTypes.Validator) bool {
		currValidators[validator.ID] = newValidatorGovInfo(
			validator.ID,
			validator.VotingPower,
			types.OptionEmpty,
		)

		return false
	})

	keeper.IterateVotes(ctx, proposal.ProposalID, func(vote types.Vote) bool {
		// if validator, just record it in the map
		if val, ok := currValidators[vote.Voter]; ok {
			val.Vote = vote.Option
			currValidators[vote.Voter] = val
		}

		keeper.deleteVote(ctx, vote.ProposalID, vote.Voter)
		return false
	})

	// iterate over the validators again to tally their voting power
	for _, val := range currValidators {
		votingPower := sdk.NewDec(val.VotingPower)
		totalBondedTokens = totalBondedTokens.Add(votingPower)

		if val.Vote == types.OptionEmpty {
			continue
		}

		results[val.Vote] = results[val.Vote].Add(votingPower)
		totalVotingPower = totalVotingPower.Add(votingPower)
	}

	tallyParams := keeper.GetTallyParams(ctx)
	tallyResults = types.NewTallyResultFromMap(results)

	// TODO: Upgrade the spec to cover all of these cases & remove pseudocode.
	// If there is no staked coins, the proposal fails
	if totalVotingPower.IsZero() {
		return false, false, tallyResults
	}

	// If there is not enough quorum of votes, the proposal fails
	percentVoting := totalVotingPower.Quo(totalBondedTokens)
	if percentVoting.LT(tallyParams.Quorum) {
		return false, true, tallyResults
	}

	// If no one votes (everyone abstains), proposal fails
	if totalVotingPower.Sub(results[types.OptionAbstain]).Equal(sdk.ZeroDec()) {
		return false, false, tallyResults
	}

	// If more than 1/3 of voters veto, proposal fails
	if results[types.OptionNoWithVeto].Quo(totalVotingPower).GT(tallyParams.Veto) {
		return false, true, tallyResults
	}

	// If more than 1/2 of non-abstaining voters vote Yes, proposal passes
	if results[types.OptionYes].Quo(totalVotingPower.Sub(results[types.OptionAbstain])).GT(tallyParams.Threshold) {
		return true, false, tallyResults
	}

	// If more than 1/2 of non-abstaining voters vote No, proposal fails
	return false, false, tallyResults
}
