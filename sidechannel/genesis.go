package sidechannel

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/sidechannel/types"
)

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	for _, pastCommit := range data.PastCommits {
		// set all txs
		if len(pastCommit.Txs) > 0 {
			for _, tx := range pastCommit.Txs {
				keeper.SetTx(ctx, pastCommit.Height, tx)
			}
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	// get all txs
	txMap := make(map[int64]tmTypes.Txs)
	keeper.IterateTxsAndApplyFn(ctx, func(height int64, tx tmTypes.Tx) error {
		if _, ok := txMap[height]; !ok {
			txMap[height] = make(tmTypes.Txs, 0)
		}
		txMap[height] = append(txMap[height], tx)
		return nil
	})

	result := make([]types.PastCommit, 0)
	for height, txs := range txMap {
		p := types.PastCommit{
			Height: height,
			Txs:    txs,
		}
		result = append(result, p)
	}

	// sort result slice
	sort.Slice(result, func(i, j int) bool {
		return result[i].Height < result[j].Height
	})

	return types.NewGenesisState(result)
}
