package keeper

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/x/sidechannel/types"
)

type (
	Keeper struct {
		cdc      codec.Marshaler
		storeKey sdk.StoreKey
	}
)

func NewKeeper(cdc codec.Marshaler, storeKey sdk.StoreKey) Keeper {
	return Keeper{
		cdc:      cdc,
		storeKey: storeKey,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

//
// Txs methods
//

// GetTx returns tx per height and hash
func (k Keeper) GetTx(ctx sdk.Context, height uint64, hash []byte) tmtypes.Tx {
	store := ctx.KVStore(k.storeKey)
	return store.Get(TxStoreKey(height, hash))
}

// HasTx checks if tx exists
func (k Keeper) HasTx(ctx sdk.Context, height uint64, hash []byte) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(TxStoreKey(height, hash))
}

// SetTx sets tx
func (k Keeper) SetTx(ctx sdk.Context, height uint64, tx tmtypes.Tx) {
	store := ctx.KVStore(k.storeKey)
	store.Set(TxStoreKey(height, tx.Hash()), tx)
}

// GetTxs returns txs per height
func (k Keeper) GetTxs(ctx sdk.Context, height uint64) (txs tmtypes.Txs) {
	// iterate through tx and append to txs
	k.IterateTxAndApplyFn(ctx, height, func(tx tmtypes.Tx) error {
		txs = append(txs, tx)
		return nil
	})

	return
}

// RemoveTx removes tx per height and hash
func (k Keeper) RemoveTx(ctx sdk.Context, height uint64, hash []byte) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(TxStoreKey(height, hash))
}

//
// Validators methods
//

// SetValidators sets validators
func (k Keeper) SetValidators(ctx sdk.Context, height uint64, validators []*abci.Validator) error {
	store := ctx.KVStore(k.storeKey)

	// marshal validators
	bz, err := k.cdc.MarshalBinaryBare(&types.PreviousValidators{
		Height:     height,
		Validators: validators,
	})
	if err != nil {
		return err
	}

	store.Set(ValidatorsKey(height), bz)
	return nil
}

// GetValidators returnss all validators
func (k Keeper) GetValidators(ctx sdk.Context, height uint64) []*abci.Validator {
	store := ctx.KVStore(k.storeKey)

	// marshal validators if exists
	if k.HasValidators(ctx, height) {
		var result types.PreviousValidators
		k.cdc.UnmarshalBinaryBare(store.Get(ValidatorsKey(height)), &result)
		return result.Validators
	}

	return nil
}

// HasValidators checks if store has validators at height
func (k Keeper) HasValidators(ctx sdk.Context, height uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(ValidatorsKey(height))
}

// RemoveValidators removes validators per height
func (k Keeper) RemoveValidators(ctx sdk.Context, height uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(ValidatorsKey(height))
}

//
// Iterators
//

// IterateTxAndApplyFn interate tx and apply the given function.
func (k Keeper) IterateTxAndApplyFn(ctx sdk.Context, height uint64, f func(tmtypes.Tx) error) {
	store := ctx.KVStore(k.storeKey)

	// get sequence iterator
	iterator := sdk.KVStorePrefixIterator(store, TxsStoreKey(height))
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
func (k Keeper) IterateTxsAndApplyFn(ctx sdk.Context, f func(uint64, tmtypes.Tx) error) {
	store := ctx.KVStore(k.storeKey)

	// get sequence iterator
	iterator := sdk.KVStorePrefixIterator(store, TxsKeyPrefix)
	defer iterator.Close()

	prefixLength := len(TxsKeyPrefix)

	// loop through validators to get valid validators
	for ; iterator.Valid(); iterator.Next() {
		heightBytes := iterator.Key()[prefixLength : 8+prefixLength]

		var height uint64
		buf := bytes.NewBuffer(heightBytes)
		binary.Read(buf, binary.BigEndian, &height)

		// call function and return if required
		if err := f(height, iterator.Value()); err != nil {
			return
		}
	}
}

// IterateValidatorsAndApplyFn interate all validators and apply the given function.
func (k Keeper) IterateValidatorsAndApplyFn(ctx sdk.Context, f func(uint64, []*abci.Validator) error) {
	store := ctx.KVStore(k.storeKey)

	// get sequence iterator
	iterator := sdk.KVStorePrefixIterator(store, ValidatorsKeyPrefix)
	defer iterator.Close()

	prefixLength := len(ValidatorsKeyPrefix)

	// loop through validators to get valid validators
	for ; iterator.Valid(); iterator.Next() {
		var previousValidators types.PreviousValidators
		if err := k.cdc.UnmarshalBinaryBare(iterator.Value(), &previousValidators); err != nil {
			return
		}

		// validators
		validators := previousValidators.Validators

		// get height bytes
		heightBytes := iterator.Key()[prefixLength : 8+prefixLength]

		var height uint64
		buf := bytes.NewBuffer(heightBytes)
		binary.Read(buf, binary.BigEndian, &height)

		// call function and return if required
		if err := f(height, validators); err != nil {
			return
		}
	}
}
