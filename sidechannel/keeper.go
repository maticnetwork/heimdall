package sidechannel

import (
	"bytes"
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/params/subspace"
	"github.com/maticnetwork/heimdall/sidechannel/types"
	"github.com/tendermint/tendermint/libs/log"
)

// Keeper stores all related data
type Keeper struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey
	// The codec codec for binary encoding/decoding of accounts.
	cdc *codec.Codec
	// code space
	codespace sdk.CodespaceType
	// param subspace
	paramSpace subspace.Subspace
}

// NewKeeper create new keeper
func NewKeeper(
	cdc *codec.Codec,
	storeKey sdk.StoreKey,
	paramSpace subspace.Subspace,
	codespace sdk.CodespaceType,
) Keeper {
	return Keeper{
		cdc:        cdc,
		key:        storeKey,
		paramSpace: paramSpace,
		codespace:  codespace,
	}
}

// Codespace returns the keeper's codespace.
func (keeper Keeper) Codespace() sdk.CodespaceType {
	return keeper.codespace
}

// Logger returns a module-specific logger
func (keeper Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

//
// Txs methods
//

// GetTx returns tx per height and hash
func (keeper Keeper) GetTx(ctx sdk.Context, height int64, hash []byte) tmTypes.Tx {
	store := ctx.KVStore(keeper.key)
	return store.Get(types.TxStoreKey(height, hash))
}

// HasTx checks if tx exists
func (keeper Keeper) HasTx(ctx sdk.Context, height int64, hash []byte) bool {
	store := ctx.KVStore(keeper.key)
	return store.Has(types.TxStoreKey(height, hash))
}

// SetTx sets tx
func (keeper Keeper) SetTx(ctx sdk.Context, height int64, tx tmTypes.Tx) {
	store := ctx.KVStore(keeper.key)
	store.Set(types.TxStoreKey(height, tx.Hash()), tx)
}

// GetTxs returns txs per height
func (keeper Keeper) GetTxs(ctx sdk.Context, height int64) (txs tmTypes.Txs) {
	// iterate through tx and append to txs
	keeper.IterateTxAndApplyFn(ctx, height, func(tx tmTypes.Tx) error {
		txs = append(txs, tx)
		return nil
	})

	return
}

// RemoveTx removes tx per height and hash
func (keeper Keeper) RemoveTx(ctx sdk.Context, height int64, hash []byte) {
	store := ctx.KVStore(keeper.key)
	store.Delete(types.TxStoreKey(height, hash))
}

//
// Validators methods
//

// SetValidators sets validators
func (keeper Keeper) SetValidators(ctx sdk.Context, height int64, validators []abci.Validator) error {
	store := ctx.KVStore(keeper.key)

	// marshal validators
	bz, err := keeper.cdc.MarshalBinaryBare(validators)
	if err != nil {
		return err
	}

	store.Set(types.ValidatorsKey(height), bz)
	return nil
}

// GetValidators returnss all validators
func (keeper Keeper) GetValidators(ctx sdk.Context, height int64) (validators []abci.Validator) {
	store := ctx.KVStore(keeper.key)

	// marshal validators if exists
	if keeper.HasValidators(ctx, height) {
		keeper.cdc.UnmarshalBinaryBare(store.Get(types.ValidatorsKey(height)), &validators)
	}

	return
}

// HasValidators checks if store has validators at height
func (keeper Keeper) HasValidators(ctx sdk.Context, height int64) bool {
	store := ctx.KVStore(keeper.key)
	return store.Has(types.ValidatorsKey(height))
}

// RemoveValidators removes validators per height
func (keeper Keeper) RemoveValidators(ctx sdk.Context, height int64) {
	store := ctx.KVStore(keeper.key)
	store.Delete(types.ValidatorsKey(height))
}

//
// Iterators
//

// IterateTxAndApplyFn interate tx and apply the given function.
func (keeper Keeper) IterateTxAndApplyFn(ctx sdk.Context, height int64, f func(tmTypes.Tx) error) {
	store := ctx.KVStore(keeper.key)

	// get sequence iterator
	iterator := sdk.KVStorePrefixIterator(store, types.TxsStoreKey(height))
	defer iterator.Close()

	// loop through validators to get valid validators
	for ; iterator.Valid(); iterator.Next() {
		val := iterator.Value()

		// call function and return if required
		if err := f(val); err != nil {
			return
		}
	}
}

// IterateTxsAndApplyFn interate all txs and apply the given function.
func (keeper Keeper) IterateTxsAndApplyFn(ctx sdk.Context, f func(int64, tmTypes.Tx) error) {
	store := ctx.KVStore(keeper.key)

	// get sequence iterator
	iterator := sdk.KVStorePrefixIterator(store, types.TxsKeyPrefix)
	defer iterator.Close()

	prefixLength := len(types.TxsKeyPrefix)

	// loop through validators to get valid validators
	for ; iterator.Valid(); iterator.Next() {
		heightBytes := iterator.Key()[prefixLength : 8+prefixLength]

		var height uint64
		buf := bytes.NewBuffer(heightBytes)
		binary.Read(buf, binary.BigEndian, &height)

		// call function and return if required
		if err := f(int64(height), iterator.Value()); err != nil {
			return
		}
	}
}

// IterateValidatorsAndApplyFn interate all validators and apply the given function.
func (keeper Keeper) IterateValidatorsAndApplyFn(ctx sdk.Context, f func(int64, []abci.Validator) error) {
	store := ctx.KVStore(keeper.key)

	// get sequence iterator
	iterator := sdk.KVStorePrefixIterator(store, types.ValidatorsKeyPrefix)
	defer iterator.Close()

	prefixLength := len(types.ValidatorsKeyPrefix)

	// loop through validators to get valid validators
	for ; iterator.Valid(); iterator.Next() {
		var validators []abci.Validator
		if err := keeper.cdc.UnmarshalBinaryBare(iterator.Value(), &validators); err != nil {
			return
		}

		heightBytes := iterator.Key()[prefixLength : 8+prefixLength]

		var height uint64
		buf := bytes.NewBuffer(heightBytes)
		binary.Read(buf, binary.BigEndian, &height)

		// call function and return if required
		if err := f(int64(height), validators); err != nil {
			return
		}
	}
}
