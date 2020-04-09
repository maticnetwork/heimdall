package gov

import (
	"errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/gov/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// GetDeposit gets the deposit of a specific depositor on a specific proposal
func (keeper Keeper) GetDeposit(ctx sdk.Context, proposalID uint64, validator hmTypes.ValidatorID) (deposit types.Deposit, found bool) {
	store := ctx.KVStore(keeper.storeKey)
	bz := store.Get(types.DepositKey(proposalID, validator))
	if bz == nil {
		return deposit, false
	}

	keeper.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &deposit)
	return deposit, true
}

func (keeper Keeper) setDeposit(ctx sdk.Context, proposalID uint64, validator hmTypes.ValidatorID, deposit types.Deposit) {
	store := ctx.KVStore(keeper.storeKey)
	bz := keeper.cdc.MustMarshalBinaryLengthPrefixed(deposit)
	store.Set(types.DepositKey(proposalID, validator), bz)
}

// AddDeposit adds or updates a deposit of a specific depositor on a specific proposal
// Activates voting period when appropriate
func (keeper Keeper) AddDeposit(ctx sdk.Context, proposalID uint64, depositorAddr hmTypes.HeimdallAddress, depositAmount sdk.Coins, validator hmTypes.ValidatorID) (sdk.Error, bool) {
	// Checks to see if proposal exists
	proposal, ok := keeper.GetProposal(ctx, proposalID)
	if !ok {
		return types.ErrUnknownProposal(keeper.codespace, proposalID), false
	}

	// Check if proposal is still depositable
	if (proposal.Status != types.StatusDepositPeriod) && (proposal.Status != types.StatusVotingPeriod) {
		return types.ErrAlreadyFinishedProposal(keeper.codespace, proposalID), false
	}

	// update the governance module's account coins pool
	err := keeper.supplyKeeper.SendCoinsFromAccountToModule(ctx, depositorAddr, types.ModuleName, depositAmount)
	if err != nil {
		return err, false
	}

	// Update proposal
	proposal.TotalDeposit = proposal.TotalDeposit.Add(depositAmount)
	keeper.SetProposal(ctx, proposal)

	// Check if deposit has provided sufficient total funds to transition the proposal into the voting period
	activatedVotingPeriod := false
	if proposal.Status == types.StatusDepositPeriod && proposal.TotalDeposit.IsAllGTE(keeper.GetDepositParams(ctx).MinDeposit) {
		keeper.activateVotingPeriod(ctx, proposal)
		activatedVotingPeriod = true
	}

	// Add or update deposit object
	deposit, found := keeper.GetDeposit(ctx, proposalID, validator)
	if found {
		deposit.Amount = deposit.Amount.Add(depositAmount)
	} else {
		deposit = types.NewDeposit(proposalID, depositAmount, validator)
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeProposalDeposit,
			sdk.NewAttribute(sdk.AttributeKeyAmount, depositAmount.String()),
			sdk.NewAttribute(types.AttributeKeyProposalID, fmt.Sprintf("%d", proposalID)),
		),
	)

	keeper.setDeposit(ctx, proposalID, validator, deposit)
	return nil, activatedVotingPeriod
}

// GetAllDeposits returns all the deposits from the store
func (keeper Keeper) GetAllDeposits(ctx sdk.Context) (deposits types.Deposits) {
	keeper.IterateAllDeposits(ctx, func(deposit types.Deposit) bool {
		deposits = append(deposits, deposit)
		return false
	})
	return
}

// GetDeposits returns all the deposits from a proposal
func (keeper Keeper) GetDeposits(ctx sdk.Context, proposalID uint64) (deposits types.Deposits) {
	keeper.IterateDeposits(ctx, proposalID, func(deposit types.Deposit) bool {
		deposits = append(deposits, deposit)
		return false
	})
	return
}

// GetDepositsIterator gets all the deposits on a specific proposal as an sdk.Iterator
func (keeper Keeper) GetDepositsIterator(ctx sdk.Context, proposalID uint64) sdk.Iterator {
	store := ctx.KVStore(keeper.storeKey)
	return sdk.KVStorePrefixIterator(store, types.DepositsKey(proposalID))
}

// RefundDeposits refunds and deletes all the deposits on a specific proposal
func (keeper Keeper) RefundDeposits(ctx sdk.Context, proposalID uint64) {
	store := ctx.KVStore(keeper.storeKey)

	keeper.IterateDeposits(ctx, proposalID, func(deposit types.Deposit) bool {
		v, ok := keeper.sk.GetValidatorFromValID(ctx, deposit.Depositor)
		if !ok {
			panic(errors.New("Not able to find validator by validator id"))
		}

		err := keeper.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, v.Signer, deposit.Amount)
		if err != nil {
			panic(err)
		}

		store.Delete(types.DepositKey(proposalID, deposit.Depositor))
		return false
	})
}
