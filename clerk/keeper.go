package clerk

import (
	"errors"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/maticnetwork/heimdall/clerk/types"
	clerkTypes "github.com/maticnetwork/heimdall/clerk/types"
	"github.com/tendermint/tendermint/libs/log"
)

var (
	StateRecordPrefixKey   = []byte{0x11} // prefix key for when storing state
	StateSyncEventCountKey = []byte{0x12} // state sync event count key
)

// Keeper stores all related data
type Keeper struct {
	cdc *codec.Codec
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
	storeKey sdk.StoreKey,
	paramSpace params.Subspace,
	codespace sdk.CodespaceType,
) Keeper {
	keeper := Keeper{
		cdc:        cdc,
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

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", clerkTypes.ModuleName)
}

// SetEventRecord adds record to store
func (k *Keeper) SetEventRecord(ctx sdk.Context, record clerkTypes.EventRecord) error {
	store := ctx.KVStore(k.storeKey)
	key := GetEventRecordKey(record.ID)

	// check if already set
	if store.Has(key) {
		return errors.New("State record already exists")
	}

	// TODO check state from mainchain

	// create Checkpoint block and marshall
	out, err := k.cdc.MarshalBinaryBare(record)
	if err != nil {
		k.Logger(ctx).Error("Error marshalling record", "error", err)
		return err
	}

	// store in key provided
	store.Set(key, out)

	// return
	return nil
}

// GetEventRecord returns record from store
func (k *Keeper) GetEventRecord(ctx sdk.Context, stateId uint64) (*clerkTypes.EventRecord, error) {
	store := ctx.KVStore(k.storeKey)
	key := GetEventRecordKey(stateId)

	// check store has data
	if store.Has(key) {
		var _record clerkTypes.EventRecord
		err := k.cdc.UnmarshalBinaryBare(store.Get(key), &_record)
		if err != nil {
			return nil, err
		}

		return &_record, nil
	}

	// return no error error
	return nil, errors.New("No record found")
}

// HasEventRecord check if state record
func (k *Keeper) HasEventRecord(ctx sdk.Context, stateID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	key := GetEventRecordKey(stateID)
	return store.Has(key)
}

// GetAllEventRecords get all state records
func (k *Keeper) GetAllEventRecords(ctx sdk.Context) (records []*types.EventRecord) {
	// iterate through spans and create span update array
	k.IterateRecordsAndApplyFn(ctx, func(record types.EventRecord) error {
		// append to list of validatorUpdates
		records = append(records, &record)
		return nil
	})

	return
}

//
// GetEventRecordKey returns key for state record
//

// GetEventRecordKey appends prefix to state id
func GetEventRecordKey(stateID uint64) []byte {
	stateIDBytes := []byte(strconv.FormatUint(stateID, 10))
	return append(StateRecordPrefixKey, stateIDBytes...)
}

//
// Utils
//

// IterateRecordsAndApplyFn interate records and apply the given function.
func (k *Keeper) IterateRecordsAndApplyFn(ctx sdk.Context, f func(record types.EventRecord) error) {
	store := ctx.KVStore(k.storeKey)

	// get span iterator
	iterator := sdk.KVStorePrefixIterator(store, StateRecordPrefixKey)
	defer iterator.Close()

	// loop through spans to get valid spans
	for ; iterator.Valid(); iterator.Next() {
		// unmarshall span
		var result types.EventRecord
		k.cdc.UnmarshalBinaryBare(iterator.Value(), &result)
		// call function and return if required
		if err := f(result); err != nil {
			return
		}
	}
}

// GetStateSyncEventCount returns next validatorID
func (k *Keeper) GetStateSyncEventCount(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.storeKey)
	// check if StateSyncEventCountKey is there
	if store.Has(StateSyncEventCountKey) {
		// get current StateSyncEventCountKey
		eventCount, err := strconv.ParseUint(string(store.Get(StateSyncEventCountKey)), 10, 64)
		if err != nil {
			k.Logger(ctx).Error("Unable to convert eventCount to int")
		} else {
			return eventCount
		}
	}
	return 0
}

// IncrementStateSyncEventCount increments state sync event count
func (k *Keeper) IncrementStateSyncEventCount(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)

	// get StateSyncEventCount
	eventCount := k.GetStateSyncEventCount(ctx)

	// increment by 1
	eventCountInBytes := []byte(strconv.FormatUint(eventCount+1, 10))

	// update
	store.Set(StateSyncEventCountKey, eventCountInBytes)
}

// SetStateSyncEventCount sets state sync event count
func (k *Keeper) SetStateSyncEventCount(ctx sdk.Context, eventCount uint64) {
	store := ctx.KVStore(k.storeKey)

	// convert state sync event to bytes
	eventCountInBytes := []byte(strconv.FormatUint(eventCount, 10))

	store.Set(StateSyncEventCountKey, eventCountInBytes)
}
