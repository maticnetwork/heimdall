package auth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
)

// // GenesisAccounts defines a slice of GenesisAccount objects
// type GenesisAccounts []GenesisAccount

// // Contains returns true if the given address exists in a slice of GenesisAccount
// // objects.
// func (accounts GenesisAccounts) Contains(addr sdk.Address) bool {
// 	for _, acc := range accounts {
// 		if acc.Address.Equals(addr) {
// 			return true
// 		}
// 	}

// 	return false
// }

// // GenesisState - all auth state that must be provided at genesis
// type GenesisState struct {
// 	CollectedFees types.Coins      `json:"collected_fees"`
// 	Params        authTypes.Params `json:"params"`
// 	Accounts      GenesisAccounts  `json:"accounts" yaml:"accounts"`
// }

// // NewGenesisState - Create a new genesis state
// func NewGenesisState(params authTypes.Params) GenesisState {
// 	return GenesisState{
// 		Params: params,
// 	}
// }

// // DefaultGenesisState - Return a default genesis state
// func DefaultGenesisState() GenesisState {
// 	return NewGenesisState(authTypes.DefaultParams())
// }

// // InitGenesis - Init store state from genesis data
// func InitGenesis(ctx sdk.Context, ak AccountKeeper, data GenesisState) {
// 	ak.SetParams(ctx, data.Params)
// }

// // ExportGenesis returns a GenesisState for a given context and keeper
// func ExportGenesis(ctx sdk.Context, ak AccountKeeper) GenesisState {
// 	params := ak.GetParams(ctx)

// 	return NewGenesisState(params)
// }

// // ValidateGenesis performs basic validation of auth genesis data returning an
// // error for any failed validation criteria.
// func ValidateGenesis(data GenesisState) error {
// 	if data.Params.TxSigLimit == 0 {
// 		return fmt.Errorf("invalid tx signature limit: %d", data.Params.TxSigLimit)
// 	}
// 	if data.Params.SigVerifyCostED25519 == 0 {
// 		return fmt.Errorf("invalid ED25519 signature verification cost: %d", data.Params.SigVerifyCostED25519)
// 	}
// 	if data.Params.SigVerifyCostSecp256k1 == 0 {
// 		return fmt.Errorf("invalid SECK256k1 signature verification cost: %d", data.Params.SigVerifyCostSecp256k1)
// 	}
// 	if data.Params.MaxMemoCharacters == 0 {
// 		return fmt.Errorf("invalid max memo characters: %d", data.Params.MaxMemoCharacters)
// 	}
// 	if data.Params.TxSizeCostPerByte == 0 {
// 		return fmt.Errorf("invalid tx size cost per byte: %d", data.Params.TxSizeCostPerByte)
// 	}
// 	return nil
// }

// InitGenesis - Init store state from genesis data
//
// CONTRACT: old coins from the FeeCollectionKeeper need to be transferred through
// a genesis port script to the new fee collector account
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
		genAccount := account.(authTypes.GenesisAccount)
		genAccounts = append(genAccounts, genAccount)
		return false
	})

	return authTypes.NewGenesisState(params, genAccounts)
}
