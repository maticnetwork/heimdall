package common

import (
	"errors"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/types"
)

//
// Bor related keepers
//
// getSpanKey appends prefix to start block
func getSpanKey(startBlock uint64) []byte {
	return append(SpanPrefixKey, []byte(strconv.FormatUint(startBlock, 10))...)
}

// AddNewSpan adds new span for bor to store
func (k *Keeper) AddNewSpan(ctx sdk.Context, span types.Span) error {
	store := ctx.KVStore(k.BorKey)
	out, err := k.cdc.MarshalBinaryBare(span)
	if err != nil {
		CheckpointLogger.Error("Error marshalling span", "error", err)
		return err
	}
	// store in key provided
	store.Set(getSpanKey(span.StartBlock), out)
	return nil
}

// GetSpan fetches span indexed by start block from store
func (k *Keeper) GetSpan(ctx sdk.Context, startBlock uint64) (span types.Span, err error) {
	store := ctx.KVStore(k.BorKey)
	spanKey := getSpanKey(startBlock)
	if !store.Has(spanKey) {
		return span, errors.New("span not found for start block")
	}
	if err := k.cdc.UnmarshalBinaryBare(store.Get(spanKey), &span); err != nil {
		return span, err
	} else {
		return span, nil
	}
}

// GetLastSpan fetches last span using lastStartBlock
func (k *Keeper) GetLastSpan(ctx sdk.Context) (lastSpan types.Span, err error) {
	store := ctx.KVStore(k.BorKey)
	var lastSpanStart uint64
	if store.Has(LastSpanStartBlockKey) {
		// get last span start block
		lastSpanStart, err = uint64(strconv.Atoi(string(store.Get(LastSpanStartBlockKey))))
		if err != nil {
			BorLogger.Error("Unable to convert start block to int")
			return lastSpan, nil
		}
	}
	return k.GetSpan(ctx, lastSpanStart)
}

// updateLastSpan updates the last span start block
func (k *Keeper) updateLastSpan(ctx sdk.Context, startBlock uint64) {
	store := ctx.KVStore(k.BorKey)
	store.Set(LastSpanStartBlockKey, []byte(strconv.FormatUint(startBlock, 10)))
}

// GetSpanDuration fetches selected span duration from store
func (k *Keeper) GetSpanDuration(ctx sdk.Context) (duration uint64, err error) {
	store := ctx.KVStore(k.BorKey)
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
