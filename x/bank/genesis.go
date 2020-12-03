package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/x/bank/keeper"
	"github.com/maticnetwork/heimdall/x/bank/types"
	// authTypes "github.com/maticnetwork/heimdall/x/auth/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	k.SetSendEnabled(ctx, genState.SendEnabled)
	// manually set the total supply based on accounts if not provided
	if genState.Supply.Total.Empty() {
		var totalSupply sdk.Coins
		// TODO
		// k.ak.IterateAccounts(ctx,
		// 	func(acc authTypes.Account) (stop bool) {
		// 		totalSupply = totalSupply.Add(acc.GetCoins())
		// 		return false
		// 	},
		// )
		genState.Supply.Total = totalSupply
	}
	k.SetSupply(ctx, genState.Supply)
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.DefaultGenesis()
}
