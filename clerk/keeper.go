package clerk

import (
	"errors"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/clerk/types"
	clerkTypes "github.com/maticnetwork/heimdall/clerk/types"
)

var (
	StateRecordPrefixKey = []byte{0x11} // prefix key for when storing state
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

// SetStateRecord adds record to store
func (k *Keeper) SetStateRecord(ctx sdk.Context, record clerkTypes.Record) error {
	store := ctx.KVStore(k.storeKey)
	key := GetStateRecordKey(record.ID)

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

// GetStateRecord returns record from store
func (k *Keeper) GetStateRecord(ctx sdk.Context, stateId uint64) (*clerkTypes.Record, error) {
	store := ctx.KVStore(k.storeKey)
	key := GetStateRecordKey(stateId)

	// check store has data
	if store.Has(key) {
		var _record clerkTypes.Record
		err := k.cdc.UnmarshalBinaryBare(store.Get(key), &_record)
		if err != nil {
			return nil, err
		}

		return &_record, nil
	}

	// return no error error
	return nil, errors.New("No record found")
}

// HasStateRecord check if state record
func (k *Keeper) HasStateRecord(ctx sdk.Context, stateID uint64) bool {
	store := ctx.KVStore(k.storeKey)
	key := GetStateRecordKey(stateID)
	return store.Has(key)
}

// GetAllStateRecords get all state records
func (k *Keeper) GetAllStateRecords(ctx sdk.Context) (records []*types.Record) {
	// iterate through spans and create span update array
	k.IterateRecordsAndApplyFn(ctx, func(record types.Record) error {
		// append to list of validatorUpdates
		records = append(records, &record)
		return nil
	})

	return
}

//
// GetStateRecordKey returns key for state record
//

// GetStateRecordKey appends prefix to state id
func GetStateRecordKey(stateID uint64) []byte {
	stateIDBytes := []byte(strconv.FormatUint(stateID, 10))
	return append(StateRecordPrefixKey, stateIDBytes...)
}

//
// Utils
//

// IterateRecordsAndApplyFn interate records and apply the given function.
func (k *Keeper) IterateRecordsAndApplyFn(ctx sdk.Context, f func(record types.Record) error) {
	store := ctx.KVStore(k.storeKey)

	// get span iterator
	iterator := sdk.KVStorePrefixIterator(store, StateRecordPrefixKey)
	defer iterator.Close()

	// loop through spans to get valid spans
	for ; iterator.Valid(); iterator.Next() {
		// unmarshall span
		var result types.Record
		k.cdc.UnmarshalBinaryBare(iterator.Value(), &result)
		// call function and return if required
		if err := f(result); err != nil {
			return
		}
	}
}
