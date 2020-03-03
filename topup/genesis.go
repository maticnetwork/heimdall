package topup

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/topup/types"
)

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	// TODO: change later
	// if data.TopupSequence > 0 {
	// keeper.SetTopupSequence(ctx, data.TopupSequence)
	// }

	// for sequence, ok := range data.TopupSequences {
	// 	if ok {
	// 		keeper.SetTopupSequence(ctx, sequence)
	// 	}
	// }
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	return types.NewGenesisState()
}
