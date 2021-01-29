package sidechannel

import (
	"sort"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/x/sidechannel/keeper"
	"github.com/maticnetwork/heimdall/x/sidechannel/types"
)

// DefaultGenesisState returns a default genesis state
//func DefaultGenesisState() *types.GenesisState {
//	result := make([]*types.PastCommit, 0)
//	return &types.GenesisState{
//		Params: types.Params{
//			Enabled: true,
//		},
//		PastCommits: result,
//	}
//}

// ValidateGenesis performs basic validation of topup genesis data returning an
// error for any failed validation criteria.
//func ValidateGenesis(data types.GenesisState) error {
//	for _, pastCommit := range data.PastCommits {
//		if pastCommit.Height <= 2 {
//			return fmt.Errorf("Past commit height must be greater 2")
//		}
//
//		if len(pastCommit.Txs) == 0 {
//			return fmt.Errorf("Txs must be present")
//		}
//	}
//	return nil
//}

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data *types.GenesisState) []abci.ValidatorUpdate {
	for _, pastCommit := range data.PastCommits {
		// set all txs
		if len(pastCommit.Txs) > 0 {
			for _, tx := range pastCommit.Txs {
				k.SetTx(ctx, pastCommit.Height, tx)
			}
		}
	}

	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	// get all txs
	txMap := make(map[uint64][][]byte)
	k.IterateTxsAndApplyFn(ctx, func(height uint64, tx tmtypes.Tx) error {
		if _, ok := txMap[height]; !ok {
			txMap[height] = make([][]byte, 0)
		}
		txMap[height] = append(txMap[height], tx)
		return nil
	})

	result := make([]*types.PastCommit, 0)
	for height, txs := range txMap {
		p := &types.PastCommit{
			Height: uint64(height),
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
