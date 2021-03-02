package topup

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/x/topup/keeper"
	"github.com/maticnetwork/heimdall/x/topup/types"
)

// InitGenesis initializes the capability module's state from a provided genesis
// state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	for _, sequence := range genState.TopupSequences {
		k.SetTopupSequence(ctx, sequence)
	}

	// Add genesis dividend accounts
	for _, dividendAccount := range genState.DividendAccounts {
		dividendAccount.User = strings.ToLower(dividendAccount.User)
		if err := k.AddDividendAccount(ctx, *dividendAccount); err != nil {
			panic((err))
		}
	}

}

// ExportGenesis returns the capability module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) types.GenesisState {
	return types.NewGenesisState(
		k.GetTopupSequences(ctx),
		k.GetAllDividendAccounts(ctx),
	)
}
