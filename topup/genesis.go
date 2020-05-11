package topup

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/topup/types"
)

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	for _, sequence := range data.TopupSequences {
		keeper.SetTopupSequence(ctx, sequence)
	}

	// Add genesis dividend accounts
	for _, dividendAccount := range data.DividentAccounts {
		if err := keeper.AddDividendAccount(ctx, dividendAccount); err != nil {
			panic((err))
		}
	}

}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	return types.NewGenesisState(
		keeper.GetTopupSequences(ctx),
		keeper.GetAllDividendAccounts(ctx),
	)
}
