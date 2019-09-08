package bor

import (
	"encoding/json"
	"errors"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	cmn "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/types"
	tmTypes "github.com/tendermint/tendermint/types"
)

var (
	DefaultValue = []byte{0x01} // Value to store in CacheCheckpoint and CacheCheckpointACK & ValidatorSetChange Flag

	SpanDurationKey       = []byte{0x24} // Key to store span duration for Bor
	SprintDurationKey     = []byte{0x25} // Key to store span duration for Bor
	LastSpanStartBlockKey = []byte{0x35} // Key to store last span start block
	SpanPrefixKey         = []byte{0x36} // prefix key to store span
	SpanCacheKey          = []byte{0x37} // key to store Cache for span
)

// Keeper stores all related data
type Keeper struct {
	cdc *codec.Codec
	sk  staking.Keeper
	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey
	// codespace
	codespace sdk.CodespaceType
	// param space
	paramSpace params.Subspace
}

// NewKeeper create new keeper
func NewKeeper(
	cdc *codec.Codec,
	stakingKeeper staking.Keeper,
	storeKey sdk.StoreKey,
	paramSpace params.Subspace,
	codespace sdk.CodespaceType,
) Keeper {
	// create keeper
	keeper := Keeper{
		cdc:        cdc,
		sk:         stakingKeeper,
		storeKey:   storeKey,
		paramSpace: paramSpace.WithKeyTable(ParamKeyTable()),
		codespace:  codespace,
	}
	return keeper
}

// Codespace returns the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// GetSpanKey appends prefix to start block
func GetSpanKey(startBlock uint64) []byte {
	return append(SpanPrefixKey, []byte(strconv.FormatUint(startBlock, 10))...)
}

// AddNewSpan adds new span for bor to store
func (k *Keeper) AddNewSpan(ctx sdk.Context, span types.Span) error {
	store := ctx.KVStore(k.storeKey)
	out, err := k.cdc.MarshalBinaryBare(span)
	if err != nil {
		cmn.BorLogger.Error("Error marshalling span", "error", err)
		return err
	}
	store.Set(GetSpanKey(span.StartBlock), out)
	// update last span
	k.UpdateLastSpan(ctx, span.StartBlock)
	// set cache for span -- which is to be cleared in end block
	k.SetSpanCache(ctx)
	return nil
}

// AddNewRawSpan adds new span for bor to store
func (k *Keeper) AddNewRawSpan(ctx sdk.Context, span types.Span) error {
	store := ctx.KVStore(k.storeKey)
	out, err := k.cdc.MarshalBinaryBare(span)
	if err != nil {
		cmn.BorLogger.Error("Error marshalling span", "error", err)
		return err
	}
	store.Set(GetSpanKey(span.StartBlock), out)
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
	var commitSigs []*tmTypes.CommitSig
	for i := range votes {
		commitSigs[i] = votes[i].CommitSig()
	}
	sigs := helper.GetSigs(commitSigs)

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
	store := ctx.KVStore(k.storeKey)
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

// GetAllSpans fetches all indexed by start block from store
func (k *Keeper) GetAllSpans(ctx sdk.Context) (spans []*types.Span) {
	// iterate through spans and create span update array
	k.IterateSpansAndApplyFn(ctx, func(span types.Span) error {
		// append to list of validatorUpdates
		spans = append(spans, &span)
		return nil
	})

	return
}

// GetLastSpan fetches last span using lastStartBlock
func (k *Keeper) GetLastSpan(ctx sdk.Context) (lastSpan types.Span, err error) {
	store := ctx.KVStore(k.storeKey)
	var lastSpanStart uint64
	if store.Has(LastSpanStartBlockKey) {
		// get last span start block
		lastSpanStartInt, err := strconv.Atoi(string(store.Get(LastSpanStartBlockKey)))
		if err != nil {
			cmn.BorLogger.Error("Unable to convert start block to int")
			return lastSpan, nil
		}
		lastSpanStart = uint64(lastSpanStartInt)
	}
	return k.GetSpan(ctx, lastSpanStart)
}

// FreezeSet freezes validator set for next span
func (k *Keeper) FreezeSet(ctx sdk.Context, startBlock uint64, borChainID string) error {
	duration := k.GetSpanDuration(ctx)

	endBlock := startBlock
	if duration > 0 {
		endBlock = endBlock + duration - 1
	}

	newSpan := types.NewSpan(startBlock, endBlock, k.sk.GetValidatorSet(ctx), k.SelectNextProducers(ctx), borChainID)
	return k.AddNewSpan(ctx, newSpan)
}

// SelectNextProducers selects producers for next span
func (k *Keeper) SelectNextProducers(ctx sdk.Context) (vals []types.Validator) {
	currVals := k.sk.GetCurrentValidators(ctx)
	// TODO add producer selection here, currently sending all validators
	return currVals
}

// UpdateLastSpan updates the last span start block
func (k *Keeper) UpdateLastSpan(ctx sdk.Context, startBlock uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(LastSpanStartBlockKey, []byte(strconv.FormatUint(startBlock, 10)))
}

// SetSpanCache sets span cache
// to be set when we freeze span
// cache to be cleared in end block
func (k *Keeper) SetSpanCache(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	// fill in default cache value
	store.Set(SpanCacheKey, DefaultValue)
}

// FlushSpanCache deletes cache stored in SpanCache
// to be called from end block to acknowledge signature aggregation
func (k *Keeper) FlushSpanCache(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(SpanCacheKey)
}

// GetSpanCache check if value exists in span cache or not
// returns true when found and false when not present
func (k *Keeper) GetSpanCache(ctx sdk.Context) bool {
	store := ctx.KVStore(k.storeKey)
	if store.Has(SpanCacheKey) {
		return true
	}
	return false
}

//
//  Params
//

// GetSpanDuration returns the span duration
func (k *Keeper) GetSpanDuration(ctx sdk.Context) uint64 {
	var duration uint64
	k.paramSpace.Get(ctx, ParamStoreKeySpanDuration, &duration)
	return duration
}

// SetSpanDuration sets the span duration
func (k *Keeper) SetSpanDuration(ctx sdk.Context, duration uint64) {
	k.paramSpace.Set(ctx, ParamStoreKeySpanDuration, duration)
}

// GetSprintDuration returns the span duration
func (k *Keeper) GetSprintDuration(ctx sdk.Context) uint64 {
	var duration uint64
	k.paramSpace.Get(ctx, ParamStoreKeySprintDuration, &duration)
	return duration
}

// SetSprintDuration sets the sprint duration
func (k *Keeper) SetSprintDuration(ctx sdk.Context, duration uint64) {
	k.paramSpace.Set(ctx, ParamStoreKeySprintDuration, duration)
}

//
// Utils
//

// IterateSpansAndApplyFn interate spans and apply the given function.
func (k *Keeper) IterateSpansAndApplyFn(ctx sdk.Context, f func(span types.Span) error) {
	store := ctx.KVStore(k.storeKey)

	// get span iterator
	iterator := sdk.KVStorePrefixIterator(store, SpanPrefixKey)
	defer iterator.Close()

	// loop through spans to get valid spans
	for ; iterator.Valid(); iterator.Next() {
		// unmarshall span
		var result types.Span
		k.cdc.UnmarshalBinaryBare(iterator.Value(), &result)
		// call function and return if required
		if err := f(result); err != nil {
			return
		}
	}
}
