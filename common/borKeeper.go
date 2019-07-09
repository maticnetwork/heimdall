package common

import (
	"errors"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

//
// Bor related keepers
//

// GetSpanDuration fetches selected span duration from store
func (k *Keeper) GetSpanDuration(ctx sdk.Context) (duration uint64, err error) {
	store := ctx.KVStore(k.CheckpointKey)
	if store.Has(SpanDurationKey) {
		duration, err := strconv.Atoi(string(store.Get(SpanDurationKey)))
		if err != nil {
			BorLogger.Error("Unable to convert key to int")
			return uint64(duration), err
		} else {
			return uint64(duration), nil
		}
	} else {
		return duration, errors.New("Duration not found")
	}
}

// TODO add setter for span duration which could be changed based on governance
