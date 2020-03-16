package auth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
)

// InitGenesis - Init store state from genesis data
func InitGenesis(ctx sdk.Context, ak AccountKeeper, processors []authTypes.AccountProcessor, data authTypes.GenesisState) {
	ak.SetParams(ctx, data.Params)
	data.Accounts = authTypes.SanitizeGenesisAccounts(data.Accounts)

	for _, gacc := range data.Accounts {
		acc := gacc.ToAccount()

		// convert to base account
		d := acc.(*authTypes.BaseAccount)

		// execute account processors
		for _, p := range processors {
			acc = p(&gacc, d)
		}

		acc = ak.NewAccount(ctx, acc)
		ak.SetAccount(ctx, acc)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Context, ak AccountKeeper) authTypes.GenesisState {
	params := ak.GetParams(ctx)

	var genAccounts authTypes.GenesisAccounts
	ak.IterateAccounts(ctx, func(acc authTypes.Account) bool {
		account, err := authTypes.NewGenesisAccountI(acc)
		if err != nil {
			panic(err)
		}
		genAccounts = append(genAccounts, account)
		return false
	})

	return authTypes.NewGenesisState(params, genAccounts)
}
