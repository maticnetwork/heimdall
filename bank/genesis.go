package bank

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	bankTypes "github.com/maticnetwork/heimdall/bank/types"
)

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data bankTypes.GenesisState) {
	keeper.SetSendEnabled(ctx, data.SendEnabled)
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) bankTypes.GenesisState {
	return bankTypes.NewGenesisState(keeper.GetSendEnabled(ctx))
}
