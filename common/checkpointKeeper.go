package common

import (
	"errors"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
)

//--------------- Checkpoint Related Keepers

// AddCheckpoint adds checkpoint into final blocks
func (k *Keeper) AddCheckpoint(ctx sdk.Context, headerBlockNumber uint64, headerBlock types.CheckpointBlockHeader) error {
	key := GetHeaderKey(headerBlockNumber)
	err := k.addCheckpoint(ctx, key, headerBlock)
	if err != nil {
		return err
	}
	CheckpointLogger.Info("Adding good checkpoint to state", "checkpoint", headerBlock, "headerBlockNumber", headerBlockNumber)
	return nil
}

// SetCheckpointBuffer flushes Checkpoint Buffer
func (k *Keeper) SetCheckpointBuffer(ctx sdk.Context, headerBlock types.CheckpointBlockHeader) error {
	err := k.addCheckpoint(ctx, BufferCheckpointKey, headerBlock)
	if err != nil {
		return err
	}
	return nil
}

// addCheckpoint adds checkpoint to store
func (k *Keeper) addCheckpoint(ctx sdk.Context, key []byte, headerBlock types.CheckpointBlockHeader) error {
	store := ctx.KVStore(k.CheckpointKey)

	// create Checkpoint block and marshall
	out, err := k.cdc.MarshalBinaryBare(headerBlock)
	if err != nil {
		CheckpointLogger.Error("Error marshalling checkpoint", "error", err)
		return err
	}

	// store in key provided
	store.Set(key, out)

	return nil
}

// To get checkpoint by header block index 10,000 ,20,000 and so on
func (k *Keeper) GetCheckpointByIndex(ctx sdk.Context, headerIndex uint64) (types.CheckpointBlockHeader, error) {
	store := ctx.KVStore(k.CheckpointKey)
	headerKey := GetHeaderKey(headerIndex)
	var _checkpoint types.CheckpointBlockHeader

	if store.Has(headerKey) {
		err := k.cdc.UnmarshalBinaryBare(store.Get(headerKey), &_checkpoint)
		if err != nil {
			return _checkpoint, err
		} else {
			return _checkpoint, nil
		}
	} else {
		return _checkpoint, errors.New("Invalid header Index")
	}
}

// GetLastCheckpoint gets last checkpoint, headerIndex = TotalACKs * ChildBlockInterval
func (k *Keeper) GetLastCheckpoint(ctx sdk.Context) (types.CheckpointBlockHeader, error) {
	store := ctx.KVStore(k.CheckpointKey)
	acksCount := k.GetACKCount(ctx)

	// fetch last checkpoint key (NumberOfACKs * ChildBlockInterval)
	lastCheckpointKey := helper.GetConfig().ChildBlockInterval * acksCount

	// fetch checkpoint and unmarshall
	var _checkpoint types.CheckpointBlockHeader

	// no checkpoint received
	if acksCount >= 0 {
		// header key
		headerKey := GetHeaderKey(lastCheckpointKey)
		if store.Has(headerKey) {
			err := k.cdc.UnmarshalBinaryBare(store.Get(headerKey), &_checkpoint)
			if err != nil {
				CheckpointLogger.Error("Unable to fetch last checkpoint from store", "key", lastCheckpointKey, "acksCount", acksCount)
				return _checkpoint, err
			} else {
				return _checkpoint, nil
			}
		}
	}
	return _checkpoint, ErrNoCheckpointFound(k.Codespace)
}

// GetHeaderKey appends prefix to headerNumber
func GetHeaderKey(headerNumber uint64) []byte {
	headerNumberBytes := []byte(strconv.FormatUint(headerNumber, 10))
	return append(HeaderBlockKey, headerNumberBytes...)
}

// SetCheckpointAckCache sets value in cache for checkpoint ACK
func (k *Keeper) SetCheckpointAckCache(ctx sdk.Context, value []byte) {
	store := ctx.KVStore(k.CheckpointKey)
	store.Set(CheckpointACKCacheKey, value)
}

func (k *Keeper) FlushACKCache(ctx sdk.Context) {
	store := ctx.KVStore(k.CheckpointKey)
	store.Delete(CheckpointACKCacheKey)
}

func (k *Keeper) FlushCheckpointCache(ctx sdk.Context) {
	store := ctx.KVStore(k.CheckpointKey)
	store.Delete(CheckpointCacheKey)
}

// SetCheckpointCache sets value in cache for checkpoint
func (k *Keeper) SetCheckpointCache(ctx sdk.Context, value []byte) {
	store := ctx.KVStore(k.CheckpointKey)
	store.Set(CheckpointCacheKey, value)
}

// GetCheckpointCache check if value exists in cache or not
func (k *Keeper) GetCheckpointCache(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.CheckpointKey)
	if store.Has(key) {
		return true
	}
	return false
}

// FlushCheckpointBuffer flushes Checkpoint Buffer
func (k *Keeper) FlushCheckpointBuffer(ctx sdk.Context) {
	store := ctx.KVStore(k.CheckpointKey)
	store.Delete(BufferCheckpointKey)
}

// GetCheckpointFromBuffer gets checkpoint in buffer
func (k *Keeper) GetCheckpointFromBuffer(ctx sdk.Context) (types.CheckpointBlockHeader, error) {
	store := ctx.KVStore(k.CheckpointKey)

	// checkpoint block header
	var checkpoint types.CheckpointBlockHeader

	if store.Has(BufferCheckpointKey) {
		// Get checkpoint and unmarshall
		err := k.cdc.UnmarshalBinaryBare(store.Get(BufferCheckpointKey), &checkpoint)
		return checkpoint, err
	}

	return checkpoint, errors.New("No checkpoint found in buffer")
}

// UpdateACKCountWithValue updates ACK with value
func (k *Keeper) UpdateACKCountWithValue(ctx sdk.Context, value uint64) {
	store := ctx.KVStore(k.CheckpointKey)

	// convert
	ackCount := []byte(strconv.FormatUint(value, 10))

	// update
	store.Set(ACKCountKey, ackCount)
}

// UpdateACKCount updates ACK count by 1
func (k *Keeper) UpdateACKCount(ctx sdk.Context) {
	store := ctx.KVStore(k.CheckpointKey)

	// get current ACK Count
	ACKCount := k.GetACKCount(ctx)

	// increment by 1
	ACKs := []byte(strconv.FormatUint(ACKCount+1, 10))

	// update
	store.Set(ACKCountKey, ACKs)
}

// GetACKCount returns current ACK count
func (k *Keeper) GetACKCount(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.CheckpointKey)
	// check if ack count is there
	if store.Has(ACKCountKey) {
		// get current ACK count
		ackCount, err := strconv.Atoi(string(store.Get(ACKCountKey)))
		if err != nil {
			CheckpointLogger.Error("Unable to convert key to int")
		} else {
			return uint64(ackCount)
		}
	}
	return 0
}

// SetLastNoAck set last no-ack object
func (k *Keeper) SetLastNoAck(ctx sdk.Context, timestamp uint64) {
	store := ctx.KVStore(k.CheckpointKey)
	// convert timestamp to bytes
	value := []byte(strconv.FormatUint(timestamp, 10))
	// set no-ack
	store.Set(CheckpointNoACKCacheKey, value)
}

// GetLastNoAck returns last no ack
func (k *Keeper) GetLastNoAck(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.CheckpointKey)
	// check if ack count is there
	if store.Has(CheckpointNoACKCacheKey) {
		// get current ACK count
		result, err := strconv.ParseUint(string(store.Get(CheckpointNoACKCacheKey)), 10, 64)
		if err == nil {
			return uint64(result)
		}
	}
	return 0
}
