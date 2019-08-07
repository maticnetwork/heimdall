package checkpoint

import (
	"errors"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	cmn "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/types"
)

var (
	DefaultValue = []byte{0x01} // Value to store in CacheCheckpoint and CacheCheckpointACK & ValidatorSetChange Flag

	BufferCheckpointKey     = []byte{0x12} // Key to store checkpoint in buffer
	HeaderBlockKey          = []byte{0x13} // prefix key for when storing header after ACK
	CheckpointCacheKey      = []byte{0x14} // key to store Cache for checkpoint
	CheckpointACKCacheKey   = []byte{0x15} // key to store Cache for checkpointACK
	CheckpointNoACKCacheKey = []byte{0x16} // key to store last no-ack
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
	keeper := Keeper{
		cdc:        cdc,
		sk:         stakingKeeper,
		storeKey:   storeKey,
		paramSpace: paramSpace,
		codespace:  codespace,
	}
	return keeper
}

// Codespace returns the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// AddCheckpoint adds checkpoint into final blocks
func (k *Keeper) AddCheckpoint(ctx sdk.Context, headerBlockNumber uint64, headerBlock types.CheckpointBlockHeader) error {
	key := GetHeaderKey(headerBlockNumber)
	err := k.addCheckpoint(ctx, key, headerBlock)
	if err != nil {
		return err
	}
	cmn.CheckpointLogger.Info("Adding good checkpoint to state", "checkpoint", headerBlock, "headerBlockNumber", headerBlockNumber)
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
	store := ctx.KVStore(k.storeKey)

	// create Checkpoint block and marshall
	out, err := k.cdc.MarshalBinaryBare(headerBlock)
	if err != nil {
		cmn.CheckpointLogger.Error("Error marshalling checkpoint", "error", err)
		return err
	}

	// store in key provided
	store.Set(key, out)

	return nil
}

// To get checkpoint by header block index 10,000 ,20,000 and so on
func (k *Keeper) GetCheckpointByIndex(ctx sdk.Context, headerIndex uint64) (types.CheckpointBlockHeader, error) {
	store := ctx.KVStore(k.storeKey)
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
	store := ctx.KVStore(k.storeKey)
	acksCount := k.sk.GetACKCount(ctx)

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
				cmn.CheckpointLogger.Error("Unable to fetch last checkpoint from store", "key", lastCheckpointKey, "acksCount", acksCount)
				return _checkpoint, err
			} else {
				return _checkpoint, nil
			}
		}
	}
	return _checkpoint, cmn.ErrNoCheckpointFound(k.Codespace())
}

// GetHeaderKey appends prefix to headerNumber
func GetHeaderKey(headerNumber uint64) []byte {
	headerNumberBytes := []byte(strconv.FormatUint(headerNumber, 10))
	return append(HeaderBlockKey, headerNumberBytes...)
}

// SetCheckpointAckCache sets value in cache for checkpoint ACK
func (k *Keeper) SetCheckpointAckCache(ctx sdk.Context, value []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(CheckpointACKCacheKey, value)
}

func (k *Keeper) FlushACKCache(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(CheckpointACKCacheKey)
}

func (k *Keeper) FlushCheckpointCache(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(CheckpointCacheKey)
}

// SetCheckpointCache sets value in cache for checkpoint
func (k *Keeper) SetCheckpointCache(ctx sdk.Context, value []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Set(CheckpointCacheKey, value)
}

// GetCheckpointCache check if value exists in cache or not
func (k *Keeper) GetCheckpointCache(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.storeKey)
	if store.Has(key) {
		return true
	}
	return false
}

// FlushCheckpointBuffer flushes Checkpoint Buffer
func (k *Keeper) FlushCheckpointBuffer(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(BufferCheckpointKey)
}

// GetCheckpointFromBuffer gets checkpoint in buffer
func (k *Keeper) GetCheckpointFromBuffer(ctx sdk.Context) (types.CheckpointBlockHeader, error) {
	store := ctx.KVStore(k.storeKey)

	// checkpoint block header
	var checkpoint types.CheckpointBlockHeader

	if store.Has(BufferCheckpointKey) {
		// Get checkpoint and unmarshall
		err := k.cdc.UnmarshalBinaryBare(store.Get(BufferCheckpointKey), &checkpoint)
		return checkpoint, err
	}

	return checkpoint, errors.New("No checkpoint found in buffer")
}

// SetLastNoAck set last no-ack object
func (k *Keeper) SetLastNoAck(ctx sdk.Context, timestamp uint64) {
	store := ctx.KVStore(k.storeKey)
	// convert timestamp to bytes
	value := []byte(strconv.FormatUint(timestamp, 10))
	// set no-ack
	store.Set(CheckpointNoACKCacheKey, value)
}

// GetLastNoAck returns last no ack
func (k *Keeper) GetLastNoAck(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
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

// GetCheckpointHeaders get checkpoint headers
func (k *Keeper) GetCheckpointHeaders(ctx sdk.Context) []types.CheckpointBlockHeader {
	store := ctx.KVStore(k.storeKey)
	// get checkpoint header iterator
	iterator := sdk.KVStorePrefixIterator(store, HeaderBlockKey)
	defer iterator.Close()

	// create headers
	var headers []types.CheckpointBlockHeader

	// loop through validators to get valid validators
	for ; iterator.Valid(); iterator.Next() {
		var checkpointHeader types.CheckpointBlockHeader
		if err := k.cdc.UnmarshalBinaryBare(iterator.Value(), &checkpointHeader); err == nil {
			headers = append(headers, checkpointHeader)
		}
	}
	return headers
}
