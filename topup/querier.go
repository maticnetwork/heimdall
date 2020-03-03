package topup

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// NewQuerier returns a new sdk.Keeper instance.
func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, sdk.Error) {
		switch path[0] {
		// TODO: Implement later
		// case types.QueryBalance:
		// 	return queryBalance(ctx, req, k)

		default:
			return nil, sdk.ErrUnknownRequest("unknown bank query endpoint")
		}
	}
}
