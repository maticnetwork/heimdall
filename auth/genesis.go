package auth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	supplyTypes "github.com/maticnetwork/heimdall/supply/types"
)

// InitGenesis - Init store state from genesis data
func InitGenesis(ctx sdk.Context, ak AccountKeeper, data authTypes.GenesisState) {
	ak.SetParams(ctx, data.Params)
	data.Accounts = authTypes.SanitizeGenesisAccounts(data.Accounts)

	for _, a := range data.Accounts {
		acc := ak.NewAccount(ctx, a)
		ak.SetAccount(ctx, acc)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, ak AccountKeeper) authTypes.GenesisState {
	params := ak.GetParams(ctx)

	var genAccounts authTypes.GenesisAccounts
	ak.IterateAccounts(ctx, func(account authTypes.Account) bool {
		if genesisAccount, ok := account.(*authTypes.GenesisAccount); ok {
			genAccounts = append(genAccounts, *genesisAccount)
		} else if supplyAccount, ok := account.(*supplyTypes.ModuleAccount); ok {
			genAccounts = append(genAccounts, authTypes.GenesisAccount{BaseAccount: supplyAccount.BaseAccount})
		}
		return false
	})

	return authTypes.NewGenesisState(params, genAccounts)
}
