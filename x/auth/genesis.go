package auth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/x/auth/keeper"
	"github.com/maticnetwork/heimdall/x/auth/types"
	authTypes "github.com/maticnetwork/heimdall/x/auth/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, ak keeper.Keeper, processors []authTypes.AccountProcessor, genState types.GenesisState) {
	ak.SetParams(ctx, genState.Params)
	genState.Accounts = authTypes.SanitizeGenesisAccounts(genState.Accounts)

	for _, gacc := range genState.Accounts {
		acc := gacc.ToAccount()

		// convert to base account
		d := acc.(*authTypes.BaseAccount)

		// execute account processors
		for _, p := range processors {
			acc = p(gacc, d)
		}

		acc = ak.NewAccount(ctx, acc)
		ak.SetAccount(ctx, acc)
	}
}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	return types.DefaultGenesis()
}
