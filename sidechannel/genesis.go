package sidechannel

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/sidechannel/types"
)

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data types.GenesisState) {
	for _, pastCommit := range data.PastCommits {
		// set validators
		keeper.SetValidators(ctx, pastCommit.Height, pastCommit.Validators)
		// set all txs
		for _, tx := range pastCommit.Txs {
			keeper.SetTx(ctx, pastCommit.Height, tx)
		}
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper.
func ExportGenesis(ctx sdk.Context, keeper Keeper) types.GenesisState {
	mappedResult := make(map[int64]bool)

	// get all txs
	txs := make(map[int64]tmTypes.Txs)
	keeper.IterateTxsAndApplyFn(ctx, func(height int64, tx tmTypes.Tx) error {
		if _, ok := txs[height]; !ok {
			txs[height] = make(tmTypes.Txs, 0)
		}
		txs[height] = append(txs[height], tx)
		mappedResult[height] = true
		return nil
	})

	// get all validators
	validators := make(map[int64][]abci.Validator)
	keeper.IterateValidatorsAndApplyFn(ctx, func(height int64, v []abci.Validator) error {
		validators[height] = v
		mappedResult[height] = true
		return nil
	})

	result := make([]types.PastCommit, 0)
	for height := range mappedResult {
		ty := types.PastCommit{}
		if r, ok := txs[height]; ok {
			ty.Txs = r
		}

		if r, ok := validators[height]; ok {
			ty.Validators = r
		}
		result = append(result, ty)
	}

	return types.NewGenesisState(result)
}
