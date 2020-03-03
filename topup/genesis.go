package topup

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/topup/types"
)

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	keeper.SetSendEnabled(ctx, data.SendEnabled)
	// TODO: change later
	// if data.TopupSequence > 0 {
	// 	keeper.SetTopupSequence(ctx, data.TopupSequence)
	// }
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	// TODO: Implement later
	return types.NewGenesisState(keeper.GetSendEnabled(ctx))
}
