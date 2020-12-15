package keeper

import (
	"fmt"
	"time"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/maticnetwork/heimdall/x/gov/types"
)

type (
	Keeper struct {
		cdc           codec.LegacyAmino
		storeKey      sdk.StoreKey
		memKey        sdk.StoreKey
		paramSubspace paramtypes.Subspace
	}
)

func NewKeeper(cdc codec.LegacyAmino, storeKey, memKey sdk.StoreKey) *Keeper {
	return &Keeper{
		cdc:      cdc,
		storeKey: storeKey,
		memKey:   memKey,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

func (keeper Keeper) SetDepositParams(ctx sdk.Context, depositParams types.DepositParams) {
	keeper.paramSubspace.Set(ctx, types.ParamStoreKeyDepositParams, &depositParams)
}

func (keeper Keeper) SetVotingParams(ctx sdk.Context, votingParams types.VotingParams) {
	keeper.paramSubspace.Set(ctx, types.ParamStoreKeyVotingParams, &votingParams)
}

func (keeper Keeper) SetTallyParams(ctx sdk.Context, tallyParams types.TallyParams) {
	keeper.paramSubspace.Set(ctx, types.ParamStoreKeyTallyParams, &tallyParams)
}

// InsertActiveProposalQueue inserts a ProposalID into the active proposal queue at endTime
func (keeper Keeper) InsertActiveProposalQueue(ctx sdk.Context, proposalID uint64, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(proposalID)
	store.Set(types.ActiveProposalQueueKey(proposalID, endTime), bz)
}

// Returns the current DepositParams from the global param store
func (keeper Keeper) GetDepositParams(ctx sdk.Context) types.DepositParams {
	var depositParams types.DepositParams
	keeper.paramSubspace.Get(ctx, types.ParamStoreKeyDepositParams, &depositParams)
	return depositParams
}

// Returns the current VotingParams from the global param store
func (keeper Keeper) GetVotingParams(ctx sdk.Context) types.VotingParams {
	var votingParams types.VotingParams
	keeper.paramSubspace.Get(ctx, types.ParamStoreKeyVotingParams, &votingParams)
	return votingParams
}

// Returns the current TallyParam from the global param store
func (keeper Keeper) GetTallyParams(ctx sdk.Context) types.TallyParams {
	var tallyParams types.TallyParams
	keeper.paramSubspace.Get(ctx, types.ParamStoreKeyTallyParams, &tallyParams)
	return tallyParams
}

// InsertInactiveProposalQueue Inserts a ProposalID into the inactive proposal queue at endTime
func (keeper Keeper) InsertInactiveProposalQueue(ctx sdk.Context, proposalID uint64, endTime time.Time) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(proposalID)
	store.Set(types.InactiveProposalQueueKey(proposalID, endTime), bz)
}

// IterateVotes iterates over the all the proposals votes and performs a callback function
func (keeper Keeper) IterateVotes(ctx sdk.Context, proposalID uint64, cb func(vote types.Vote) (stop bool)) {
	iterator := keeper.GetVotesIterator(ctx, proposalID)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var vote types.Vote
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &vote)

		if cb(vote) {
			break
		}
	}
}

// IterateDeposits iterates over the all the proposals deposits and performs a callback function
func (keeper Keeper) IterateDeposits(ctx sdk.Context, proposalID uint64, cb func(deposit types.Deposit) (stop bool)) {
	iterator := keeper.GetDepositsIterator(ctx, proposalID)

	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var deposit types.Deposit
		keeper.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &deposit)

		if cb(deposit) {
			break
		}
	}
}