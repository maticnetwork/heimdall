package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/x/gov/types"
)

// AddVote Adds a vote on a specific proposal
func (keeper Keeper) AddVote(ctx sdk.Context, proposalID uint64, voter sdk.AccAddress, option types.VoteOption, validator hmTypes.ValidatorID) error {
	proposal, ok := keeper.GetProposal(ctx, proposalID)
	if !ok {
		return types.ErrUnknownProposal
	}
	if proposal.Status != types.StatusVotingPeriod {
		return types.ErrInactiveProposal
	}

	if !types.ValidVoteOption(option) {
		return types.ErrInvalidVote
	}

	vote := types.NewVote(proposalID, validator, option)
	keeper.SetVote(ctx, proposalID, validator, vote)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProposalVote,
			sdk.NewAttribute(types.AttributeKeyOption, option.String()),
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)

	return nil
}

// GetAllVotes returns all the votes from the store
func (keeper Keeper) GetAllVotes(ctx sdk.Context) (votes types.Votes) {
	keeper.IterateAllVotes(ctx, func(vote types.Vote) bool {
		votes = append(votes, vote)
		return false
	})
	return
}

// GetVotes returns all the votes from a proposal
func (keeper Keeper) GetVotes(ctx sdk.Context, proposalID uint64) (votes types.Votes) {
	keeper.IterateVotes(ctx, proposalID, func(vote types.Vote) bool {
		votes = append(votes, vote)
		return false
	})
	return
}

// GetVote gets the vote from an address on a specific proposal
func (keeper Keeper) GetVote(ctx sdk.Context, proposalID uint64, voter hmTypes.ValidatorID) (vote types.Vote, found bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.VoteKey(proposalID, voter))
	if bz == nil {
		return vote, false
	}

	keeper.cdc.MustUnmarshalBinaryBare(bz, &vote)
	return vote, true
}

// SetVote sets a Vote to the gov store
func (keeper Keeper) SetVote(ctx sdk.Context, proposalID uint64, voter hmTypes.ValidatorID, vote types.Vote) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryBare(&vote)
	store.Set(types.VoteKey(proposalID, voter), bz)
}

// GetVotesIterator gets all the votes on a specific proposal as an sdk.Iterator
func (keeper Keeper) GetVotesIterator(ctx sdk.Context, proposalID uint64) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return sdk.KVStorePrefixIterator(store, types.VotesKey(proposalID))
}

func (keeper Keeper) deleteVote(ctx sdk.Context, proposalID uint64, voter hmTypes.ValidatorID) {
	store := ctx.KVStore(keeper.storeKey)
	store.Delete(types.VoteKey(proposalID, voter))
}

// IterateAllVotes iterates over the all the stored votes and performs a callback function
func (keeper Keeper) IterateAllVotes(ctx sdk.Context, cb func(vote types.Vote) (stop bool)) {
	store := ctx.KVStore(keeper.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.VotesKeyPrefix)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var vote types.Vote
		keeper.cdc.MustUnmarshalBinaryBare(iterator.Value(), &vote)

		if cb(vote) {
			break
		}
	}
}
