package checkpoint

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/chainmanager"
	"github.com/maticnetwork/heimdall/checkpoint/types"
	cmn "github.com/maticnetwork/heimdall/common"
	"github.com/maticnetwork/heimdall/params/subspace"
	"github.com/maticnetwork/heimdall/staking"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

var (
	DefaultValue = []byte{0x01} // Value to store in CacheCheckpoint and CacheCheckpointACK & ValidatorSetChange Flag

	ACKCountKey             = []byte{0x11} // key to store ACK count
	BufferCheckpointKey     = []byte{0x12} // Key to store checkpoint in buffer
	CheckpointKeyDeprecated = []byte{0x13} // use CheckpointKey instead
	LastNoACKKey            = []byte{0x14} // key to store last no-ack
	CheckpointKey           = []byte{0x15} // prefix key for when storing checkpoint after ACK

	checkpointKeyV2MigrationMu  = sync.Mutex{}
	checkpointKeyV2MigrationKey = []byte{0x2}
)

// ModuleCommunicator manages different module interaction
type ModuleCommunicator interface {
	GetAllDividendAccounts(ctx sdk.Context) []hmTypes.DividendAccount
}

// Keeper stores all related data
type Keeper struct {
	cdc *codec.Codec
	// staking keeper
	sk staking.Keeper
	ck chainmanager.Keeper
	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey
	// codespace
	codespace sdk.CodespaceType
	// param space
	paramSpace subspace.Subspace

	// module communicator
	moduleCommunicator ModuleCommunicator
}

// NewKeeper create new keeper
func NewKeeper(
	cdc *codec.Codec,
	storeKey sdk.StoreKey,
	paramSpace subspace.Subspace,
	codespace sdk.CodespaceType,
	stakingKeeper staking.Keeper,
	chainKeeper chainmanager.Keeper,
	moduleCommunicator ModuleCommunicator,
) Keeper {
	keeper := Keeper{
		cdc:                cdc,
		storeKey:           storeKey,
		paramSpace:         paramSpace.WithKeyTable(types.ParamKeyTable()),
		codespace:          codespace,
		sk:                 stakingKeeper,
		ck:                 chainKeeper,
		moduleCommunicator: moduleCommunicator,
	}

	return keeper
}

// Codespace returns the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

// AddCheckpoint adds checkpoint into final blocks
func (k *Keeper) AddCheckpoint(ctx sdk.Context, checkpointNumber uint64, checkpoint hmTypes.Checkpoint) error {
	key := GetCheckpointKey(checkpointNumber)
	if err := k.addCheckpoint(ctx, key, checkpoint); err != nil {
		return err
	}

	k.Logger(ctx).Info("Adding good checkpoint to state", "checkpoint", checkpoint, "checkpointNumber", checkpointNumber)

	return nil
}

// SetCheckpointBuffer flushes Checkpoint Buffer
func (k *Keeper) SetCheckpointBuffer(ctx sdk.Context, checkpoint hmTypes.Checkpoint) error {
	err := k.addCheckpoint(ctx, BufferCheckpointKey, checkpoint)
	if err != nil {
		return err
	}

	return nil
}

// addCheckpoint adds checkpoint to store
func (k *Keeper) addCheckpoint(ctx sdk.Context, key []byte, checkpoint hmTypes.Checkpoint) error {
	store := ctx.KVStore(k.storeKey)

	// create Checkpoint block and marshall
	out, err := k.cdc.MarshalBinaryBare(checkpoint)
	if err != nil {
		k.Logger(ctx).Error("Error marshalling checkpoint", "error", err)
		return err
	}

	// store in key provided
	store.Set(key, out)

	return nil
}

// GetCheckpointByNumber to get checkpoint by checkpoint number
func (k *Keeper) GetCheckpointByNumber(ctx sdk.Context, number uint64) (hmTypes.Checkpoint, error) {
	if err := k.MigrateOnceCheckpointKeyV2(ctx); err != nil {
		return hmTypes.Checkpoint{}, err
	}

	store := ctx.KVStore(k.storeKey)
	checkpointKey := GetCheckpointKey(number)

	var _checkpoint hmTypes.Checkpoint

	if store.Has(checkpointKey) {
		if err := k.cdc.UnmarshalBinaryBare(store.Get(checkpointKey), &_checkpoint); err != nil {
			return _checkpoint, err
		}

		return _checkpoint, nil
	}

	return _checkpoint, errors.New("Invalid checkpoint Index")
}

// GetCheckpointList returns all checkpoints with params like page and limit
func (k *Keeper) GetCheckpointList(ctx sdk.Context, page uint64, limit uint64) ([]hmTypes.Checkpoint, error) {
	if err := k.MigrateOnceCheckpointKeyV2(ctx); err != nil {
		return nil, err
	}

	store := ctx.KVStore(k.storeKey)

	// create headers
	var checkpoints []hmTypes.Checkpoint

	// have max limit
	if limit > 20 {
		limit = 20
	}

	// get paginated iterator
	iterator := hmTypes.KVStorePrefixIteratorPaginated(store, CheckpointKey, uint(page), uint(limit))

	// loop through validators to get valid validators
	for ; iterator.Valid(); iterator.Next() {
		var checkpoint hmTypes.Checkpoint
		if err := k.cdc.UnmarshalBinaryBare(iterator.Value(), &checkpoint); err == nil {
			checkpoints = append(checkpoints, checkpoint)
		}
	}

	return checkpoints, nil
}

// GetLastCheckpoint gets last checkpoint, checkpoint number = TotalACKs
func (k *Keeper) GetLastCheckpoint(ctx sdk.Context) (hmTypes.Checkpoint, error) {
	if err := k.MigrateOnceCheckpointKeyV2(ctx); err != nil {
		return hmTypes.Checkpoint{}, err
	}

	store := ctx.KVStore(k.storeKey)
	acksCount := k.GetACKCount(ctx)

	lastCheckpointKey := acksCount

	// fetch checkpoint and unmarshall
	var _checkpoint hmTypes.Checkpoint

	// no checkpoint received
	// header key
	headerKey := GetCheckpointKey(lastCheckpointKey)
	if store.Has(headerKey) {
		err := k.cdc.UnmarshalBinaryBare(store.Get(headerKey), &_checkpoint)
		if err != nil {
			k.Logger(ctx).Error("Unable to fetch last checkpoint from store", "key", lastCheckpointKey, "acksCount", acksCount)
			return _checkpoint, err
		} else {
			return _checkpoint, nil
		}
	}

	return _checkpoint, cmn.ErrNoCheckpointFound(k.Codespace())
}

// GetCheckpointKey appends prefix to checkpointNumber
func GetCheckpointKey(checkpointNumber uint64) []byte {
	checkpointNumberBytes := sdk.Uint64ToBigEndian(checkpointNumber)
	return append(CheckpointKey, checkpointNumberBytes...)
}

// HasStoreValue check if value exists in store or not
func (k *Keeper) HasStoreValue(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(key)
}

// FlushCheckpointBuffer flushes Checkpoint Buffer
func (k *Keeper) FlushCheckpointBuffer(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(BufferCheckpointKey)
}

// GetCheckpointFromBuffer gets checkpoint in buffer
func (k *Keeper) GetCheckpointFromBuffer(ctx sdk.Context) (*hmTypes.Checkpoint, error) {
	store := ctx.KVStore(k.storeKey)

	// checkpoint block header
	var checkpoint hmTypes.Checkpoint

	if store.Has(BufferCheckpointKey) {
		// Get checkpoint and unmarshall
		err := k.cdc.UnmarshalBinaryBare(store.Get(BufferCheckpointKey), &checkpoint)
		return &checkpoint, err
	}

	return nil, errors.New("No checkpoint found in buffer")
}

// SetLastNoAck set last no-ack object
func (k *Keeper) SetLastNoAck(ctx sdk.Context, timestamp uint64) {
	store := ctx.KVStore(k.storeKey)
	// convert timestamp to bytes
	value := []byte(strconv.FormatUint(timestamp, 10))
	// set no-ack
	store.Set(LastNoACKKey, value)
}

// GetLastNoAck returns last no ack
func (k *Keeper) GetLastNoAck(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	// check if ack count is there
	if store.Has(LastNoACKKey) {
		// get current ACK count
		result, err := strconv.ParseUint(string(store.Get(LastNoACKKey)), 10, 64)
		if err == nil {
			return result
		}
	}

	return 0
}

// GetCheckpoints get checkpoint all checkpoints
func (k *Keeper) GetCheckpoints(ctx sdk.Context) []hmTypes.Checkpoint {
	if err := k.MigrateOnceCheckpointKeyV2(ctx); err != nil {
		panic(fmt.Errorf("issue while performing checkpoint key v2 migration: %w", err))
	}

	store := ctx.KVStore(k.storeKey)
	// get checkpoint header iterator
	iterator := sdk.KVStorePrefixIterator(store, CheckpointKey)
	defer iterator.Close()

	// create headers
	var headers []hmTypes.Checkpoint

	// loop through validators to get valid validators
	for ; iterator.Valid(); iterator.Next() {
		var checkpoint hmTypes.Checkpoint
		if err := k.cdc.UnmarshalBinaryBare(iterator.Value(), &checkpoint); err == nil {
			headers = append(headers, checkpoint)
		}
	}

	return headers
}

//
// Ack count
//

// GetACKCount returns current ACK count
func (k Keeper) GetACKCount(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	// check if ack count is there
	if store.Has(ACKCountKey) {
		// get current ACK count
		ackCount, err := strconv.ParseUint(string(store.Get(ACKCountKey)), 10, 64)
		if err != nil {
			k.Logger(ctx).Error("Unable to convert key to int")
		} else {
			return ackCount
		}
	}

	return 0
}

// UpdateACKCountWithValue updates ACK with value
func (k Keeper) UpdateACKCountWithValue(ctx sdk.Context, value uint64) {
	store := ctx.KVStore(k.storeKey)

	// convert
	ackCount := []byte(strconv.FormatUint(value, 10))

	// update
	store.Set(ACKCountKey, ackCount)
}

// UpdateACKCount updates ACK count by 1
func (k Keeper) UpdateACKCount(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)

	// get current ACK Count
	ACKCount := k.GetACKCount(ctx)

	// increment by 1
	ACKs := []byte(strconv.FormatUint(ACKCount+1, 10))

	// update
	store.Set(ACKCountKey, ACKs)
}

// -----------------------------------------------------------------------------
// Params

// SetParams sets the auth module's parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetParams gets the auth module's parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return
}

// MigrateOnceCheckpointKeyV2 is used to migrate:
//   - from old checkpoint keys that were incorrectly stored as string versions of uint64 which leads
//     to incorrect lexicographical ordering to keys that are
//   - to new checkpoint keys that are marshalled uint64s to a bigendian byte slices, so that they can be sorted
//     numerically
//
// Note: this migration is guaranteed to be performed only once since all the old keys are migrated, then deleted and
// never written again.
func (k *Keeper) MigrateOnceCheckpointKeyV2(ctx sdk.Context) error {
	store := ctx.KVStore(k.storeKey)
	if store.Has(checkpointKeyV2MigrationKey) {
		// check first before acquiring mutex
		// if migration already done, no need to use a mutex
		return nil
	}

	checkpointKeyV2MigrationMu.Lock()
	defer checkpointKeyV2MigrationMu.Unlock()

	if store.Has(checkpointKeyV2MigrationKey) {
		// while waiting on the mutex maybe another goroutine
		// has already performed the migration -> double check
		return nil
	}

	iterator := sdk.KVStorePrefixIterator(store, CheckpointKeyDeprecated)
	defer iterator.Close()

	k.Logger(ctx).Info("performing migration from lexicographically to numerically sorted checkpoint keys")

	var oldKeys [][]byte
	for ; iterator.Valid(); iterator.Next() {
		var checkpoint hmTypes.Checkpoint
		if err := k.cdc.UnmarshalBinaryBare(iterator.Value(), &checkpoint); err != nil {
			return err
		}

		key := iterator.Key()
		checkpointNumber, err := strconv.ParseUint(string(key[len(CheckpointKeyDeprecated):]), 10, 64)
		if err != nil {
			return err
		}

		if err := k.AddCheckpoint(ctx, checkpointNumber, checkpoint); err != nil {
			return err
		}

		oldKeys = append(oldKeys, key)
	}

	for _, oldKey := range oldKeys {
		store.Delete(oldKey)
	}

	store.Set(checkpointKeyV2MigrationKey, []byte{1})

	return nil
}
