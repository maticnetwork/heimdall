package app

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/types"
)

// ModuleCommunicator retriever
type ModuleCommunicator struct {
	App *HeimdallApp
}

// GetACKCount returns ack count
func (d ModuleCommunicator) GetACKCount(ctx sdk.Context) uint64 {
	return d.App.CheckpointKeeper.GetACKCount(ctx)
}

// IsCurrentValidatorByAddress check if validator is current validator
func (d ModuleCommunicator) IsCurrentValidatorByAddress(ctx sdk.Context, address []byte) bool {
	return d.App.StakingKeeper.IsCurrentValidatorByAddress(ctx, address)
}

// GetAllDividendAccounts fetches all dividend accounts from topup module
func (d ModuleCommunicator) GetAllDividendAccounts(ctx sdk.Context) []*types.DividendAccount {
	return d.App.TopupKeeper.GetAllDividendAccounts(ctx)
}

// GetValidatorFromValID get validator from validator id
func (d ModuleCommunicator) GetValidatorFromValID(ctx sdk.Context, valID types.ValidatorID) (validator types.Validator, ok bool) {
	return d.App.StakingKeeper.GetValidatorFromValID(ctx, valID)
}

// CreateValiatorSigningInfo creates ValidatorSigningInfo used by slashing module
func (d ModuleCommunicator) CreateValiatorSigningInfo(ctx sdk.Context, valID types.ValidatorID, valSigningInfo types.ValidatorSigningInfo) {
	return
}
