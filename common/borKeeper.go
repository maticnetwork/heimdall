package common

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/maticnetwork/heimdall/helper"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/types"
	tmTypes "github.com/tendermint/tendermint/types"
)

//
// Bor related keepers
//

// GetSpanKey appends prefix to start block
func GetSpanKey(startBlock uint64) []byte {
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
	store.Set(GetSpanKey(span.StartBlock), out)
	// update last span
	k.UpdateLastSpan(ctx, span.StartBlock)
	// set cache for span -- which is to be cleared in end block
	k.SetSpanCache(ctx)
	return nil
}

// AddSigs adds collected signatures to the last span added
func (k *Keeper) AddSigs(ctx sdk.Context, tmVotes []byte) error {
	lastSpan, err := k.GetLastSpan(ctx)
	if err != nil {
		return err
	}

	var votes []tmTypes.Vote
	err = json.Unmarshal(tmVotes, &votes)
	if err != nil {
		return err
	}

	sigs := helper.GetSigs(votes)

	lastSpan.AddSigs(sigs)
	if err := k.AddNewSpan(ctx, lastSpan); err != nil {
		return err
	}

	// clear span cache
	k.FlushSpanCache(ctx)
	return nil
}

// GetSpan fetches span indexed by start block from store
func (k *Keeper) GetSpan(ctx sdk.Context, startBlock uint64) (span types.Span, err error) {
	store := ctx.KVStore(k.BorKey)
	spanKey := GetSpanKey(startBlock)

	// If we are starting from 0 there will be no spanKey present
	if !store.Has(spanKey) && startBlock != 0 {
		return span, errors.New("span not found for start block")
	} else if !store.Has(spanKey) && startBlock == 0 {
		return span, nil
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
		lastSpanStartInt, err := strconv.Atoi(string(store.Get(LastSpanStartBlockKey)))
		if err != nil {
			BorLogger.Error("Unable to convert start block to int")
			return lastSpan, nil
		}
		lastSpanStart = uint64(lastSpanStartInt)
	}
	return k.GetSpan(ctx, lastSpanStart)
}

// FreezeSet freezes validator set for next span
func (k *Keeper) FreezeSet(ctx sdk.Context, startBlock uint64) error {
	duration, err := k.GetSpanDuration(ctx)
	if err != nil {
		return err
	}
	newSpan := types.NewSpan(startBlock, startBlock+duration, k.GetValidatorSet(ctx), k.SelectNextProducers(ctx), helper.GetBorChainID())
	return k.AddNewSpan(ctx, newSpan)
}

// SelectNextProducers selects producers for next span
func (k *Keeper) SelectNextProducers(ctx sdk.Context) (vals []types.Validator) {
	currVals := k.GetCurrentValidators(ctx)
	// TODO add producer selection here, currently sending all validators
	return currVals
}

// UpdateLastSpan updates the last span start block
func (k *Keeper) UpdateLastSpan(ctx sdk.Context, startBlock uint64) {
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
		return duration, errors.New("duration not found")
	}
}

// SetSpanDuration sets span duration
func (k *Keeper) SetSpanDuration(ctx sdk.Context, duration uint64) {
	store := ctx.KVStore(k.BorKey)
	store.Set(SpanDurationKey, []byte(strconv.FormatUint(duration, 10)))
}

// SetSpanCache sets span cache
// to be set when we freeze span
// cache to be cleared in end block
func (k *Keeper) SetSpanCache(ctx sdk.Context) {
	store := ctx.KVStore(k.BorKey)
	// fill in default cache value
	store.Set(SpanCacheKey, DefaultValue)
}

// FlushSpanCache deletes cache stored in SpanCache
// to be called from end block to acknowledge signature aggregation
func (k *Keeper) FlushSpanCache(ctx sdk.Context) {
	store := ctx.KVStore(k.BorKey)
	store.Delete(SpanCacheKey)
}

// GetSpanCache check if value exists in span cache or not
// returns true when found and false when not present
func (k *Keeper) GetSpanCache(ctx sdk.Context) bool {
	store := ctx.KVStore(k.BorKey)
	if store.Has(SpanCacheKey) {
		return true
	}
	return false
}

// TODO add setter for span duration which could be changed by submitting a transaction
