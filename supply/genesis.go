package supply

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	auth "github.com/maticnetwork/heimdall/auth"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	supplyTypes "github.com/maticnetwork/heimdall/supply/types"
	"github.com/maticnetwork/heimdall/types"
)

// GenesisState is the supply state that must be provided at genesis.
type GenesisState struct {
	Supply supplyTypes.Supply `json:"supply" yaml:"supply"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(supply supplyTypes.Supply) GenesisState {
	return GenesisState{supply}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(supplyTypes.DefaultSupply())
}

// InitGenesis sets supply information for genesis.
//
// CONTRACT: all types of accounts must have been already initialized/created
func InitGenesis(ctx sdk.Context, keeper Keeper, ak auth.AccountKeeper, data GenesisState) {
	// manually set the total supply based on accounts if not provided
	if data.Supply.Total.Empty() {
		var totalSupply types.Coins
		ak.IterateAccounts(ctx,
			func(acc authTypes.Account) (stop bool) {
				totalSupply = totalSupply.Add(acc.GetCoins())
				return false
			},
		)
		data.Supply.Total = totalSupply
	}
	keeper.SetSupply(ctx, data.Supply)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) GenesisState {
	return NewGenesisState(keeper.GetSupply(ctx))
}

// ValidateGenesis performs basic validation of supply genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error {
	return data.Supply.ValidateBasic()
}
