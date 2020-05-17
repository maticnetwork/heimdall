package clerk

import (
	"errors"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/chainmanager"
	"github.com/maticnetwork/heimdall/clerk/types"
	"github.com/maticnetwork/heimdall/params/subspace"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

var (
	StateRecordPrefixKey = []byte{0x11} // prefix key for when storing state

	// DefaultValue default value
	DefaultValue = []byte{0x01}

	// RecordSequencePrefixKey represents record sequence prefix key
	RecordSequencePrefixKey = []byte{0x12}

	StateRecordPrefixKeyWithTime = []byte{0x13} // prefix key for when storing state with time
)

// Keeper stores all related data
type Keeper struct {
	cdc *codec.Codec
	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey
	// codespace
	codespace sdk.CodespaceType
	// param space
	paramSpace subspace.Subspace
	// chain param keeper
	chainKeeper chainmanager.Keeper
}

// NewKeeper create new keeper
func NewKeeper(
	cdc *codec.Codec,
	storeKey sdk.StoreKey,
	paramSpace subspace.Subspace,
	codespace sdk.CodespaceType,
	chainKeeper chainmanager.Keeper,
) Keeper {
	keeper := Keeper{
		cdc:         cdc,
		storeKey:    storeKey,
		paramSpace:  paramSpace,
		codespace:   codespace,
		chainKeeper: chainKeeper,
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

// SetEventRecordWithTime sets event record id with time
func (k *Keeper) SetEventRecordWithTime(ctx sdk.Context, record types.EventRecord) error {
	key := GetEventRecordKeyWithTime(record.ID, record.RecordTime)
	value, err := k.cdc.MarshalBinaryBare(record.ID)
	if err != nil {
		k.Logger(ctx).Error("Error marshalling record", "error", err)
		return err
	}

	if err := k.setEventRecordStore(ctx, key, value); err != nil {
		return err
	}
	return nil
}

// SetEventRecordWithID adds record to store with ID
func (k *Keeper) SetEventRecordWithID(ctx sdk.Context, record types.EventRecord) error {
	key := GetEventRecordKey(record.ID)
	value, err := k.cdc.MarshalBinaryBare(record)
	if err != nil {
		k.Logger(ctx).Error("Error marshalling record", "error", err)
		return err
	}

	if err := k.setEventRecordStore(ctx, key, value); err != nil {
		return err
	}
	return nil
}

// setEventRecordStore adds value to store by key
func (k *Keeper) setEventRecordStore(ctx sdk.Context, key, value []byte) error {
	store := ctx.KVStore(k.storeKey)
	// check if already set
	if store.Has(key) {
		return errors.New("Key already exists")
	}

	// store value in provided key
	store.Set(key, value)

	// return
	return nil
}

// SetEventRecord adds record to store
func (k *Keeper) SetEventRecord(ctx sdk.Context, record types.EventRecord) error {
	if err := k.SetEventRecordWithID(ctx, record); err != nil {
		return err
	}
	if err := k.SetEventRecordWithTime(ctx, record); err != nil {
		return err
	}
	return nil
}

// GetEventRecord returns record from store
func (k *Keeper) GetEventRecord(ctx sdk.Context, stateID uint64) (*types.EventRecord, error) {
	store := ctx.KVStore(k.storeKey)
	key := GetEventRecordKey(stateID)

	// check store has data
	if store.Has(key) {
		var _record types.EventRecord
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

// GetEventRecordList returns all records with params like page and limit
func (k *Keeper) GetEventRecordList(ctx sdk.Context, page uint64, limit uint64) ([]types.EventRecord, error) {
	store := ctx.KVStore(k.storeKey)

	// create records
	var records []types.EventRecord

	// have max limit
	if limit > 50 {
		limit = 50
	}

	// get paginated iterator
	iterator := hmTypes.KVStorePrefixIteratorPaginated(store, StateRecordPrefixKey, uint(page), uint(limit))

	// loop through records to get valid records
	for ; iterator.Valid(); iterator.Next() {
		var record types.EventRecord
		if err := k.cdc.UnmarshalBinaryBare(iterator.Value(), &record); err == nil {
			records = append(records, record)
		}
	}

	return records, nil
}

// GetEventRecordListWithTime returns all records with params like fromTime and toTime
func (k *Keeper) GetEventRecordListWithTime(ctx sdk.Context, fromTime, toTime time.Time, page, limit uint64) ([]types.EventRecord, error) {
	var iterator sdk.Iterator
	store := ctx.KVStore(k.storeKey)

	// create records
	var records []types.EventRecord

	// have max limit
	if limit > 50 {
		limit = 50
	}

	if page == 0 && limit == 0 {
		iterator = store.Iterator(GetEventRecordKeyWithTimePrefix(fromTime), GetEventRecordKeyWithTimePrefix(toTime))
	} else {
		iterator = hmTypes.KVStorePrefixRangeIteratorPaginated(store, uint(page), uint(limit), GetEventRecordKeyWithTimePrefix(fromTime), GetEventRecordKeyWithTimePrefix(toTime))
	}

	// get range iterator
	defer iterator.Close()
	// loop through records to get valid records
	for ; iterator.Valid(); iterator.Next() {
		var stateID uint64
		if err := k.cdc.UnmarshalBinaryBare(iterator.Value(), &stateID); err == nil {
			record, err := k.GetEventRecord(ctx, stateID)
			if err != nil {
				k.Logger(ctx).Error("GetEventRecordListWithTime | GetEventRecord", "error", err)
				continue
			}
			records = append(records, *record)
		}
	}

	return records, nil
}

//
// GetEventRecordKey returns key for state record
//

// GetEventRecordKey appends prefix to state id
func GetEventRecordKey(stateID uint64) []byte {
	stateIDBytes := []byte(strconv.FormatUint(stateID, 10))
	return append(StateRecordPrefixKey, stateIDBytes...)
}

// GetEventRecordKeyWithTime appends prefix to state id and record time
func GetEventRecordKeyWithTime(stateID uint64, recordTime time.Time) []byte {
	stateIDBytes := []byte(strconv.FormatUint(stateID, 10))
	return append(GetEventRecordKeyWithTimePrefix(recordTime), stateIDBytes...)
}

// GetEventRecordKeyWithTimePrefix gives prefix for record time key
func GetEventRecordKeyWithTimePrefix(recordTime time.Time) []byte {
	recordTimeBytes := sdk.FormatTimeBytes(recordTime)
	return append(StateRecordPrefixKeyWithTime, recordTimeBytes...)
}

// GetRecordSequenceKey returns record sequence key
func GetRecordSequenceKey(sequence string) []byte {
	return append(RecordSequencePrefixKey, []byte(sequence)...)
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
		if err := k.cdc.UnmarshalBinaryBare(iterator.Value(), &result); err != nil {
			k.Logger(ctx).Error("IterateRecordsAndApplyFn | UnmarshalBinaryBare", "error", err)
			return
		}
		// call function and return if required
		if err := f(result); err != nil {
			return
		}
	}
}

// GetRecordSequences checks if record already exists
func (k *Keeper) GetRecordSequences(ctx sdk.Context) (sequences []string) {
	k.IterateRecordSequencesAndApplyFn(ctx, func(sequence string) error {
		sequences = append(sequences, sequence)
		return nil
	})
	return
}

// IterateRecordSequencesAndApplyFn interate records and apply the given function.
func (k *Keeper) IterateRecordSequencesAndApplyFn(ctx sdk.Context, f func(sequence string) error) {
	store := ctx.KVStore(k.storeKey)

	// get sequence iterator
	iterator := sdk.KVStorePrefixIterator(store, RecordSequencePrefixKey)
	defer iterator.Close()

	// loop through sequences
	for ; iterator.Valid(); iterator.Next() {
		sequence := string(iterator.Key()[len(RecordSequencePrefixKey):])

		// call function and return if required
		if err := f(sequence); err != nil {
			return
		}
	}
}

// SetRecordSequence sets mapping for sequence id to bool
func (k *Keeper) SetRecordSequence(ctx sdk.Context, sequence string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetRecordSequenceKey(sequence), DefaultValue)
}

// HasRecordSequence checks if record already exists
func (k *Keeper) HasRecordSequence(ctx sdk.Context, sequence string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(GetRecordSequenceKey(sequence))
}
