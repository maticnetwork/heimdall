package supply

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	auth "github.com/maticnetwork/heimdall/auth"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/supply/types"
)

// InitGenesis sets supply information for genesis.
//
// CONTRACT: all types of accounts must have been already initialized/created
func InitGenesis(ctx sdk.Context, keeper Keeper, ak auth.AccountKeeper, data types.GenesisState) {
	// manually set the total supply based on accounts if not provided
	if data.Supply.Total.Empty() {
		var totalSupply sdk.Coins
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
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	return types.NewGenesisState(keeper.GetSupply(ctx))
}
