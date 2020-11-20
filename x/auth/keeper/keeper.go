package keeper

import (
	"fmt"

	gogotypes "github.com/gogo/protobuf/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/x/auth/types"
)

type (
	Keeper struct {
		cdc           codec.Marshaler
		storeKey      sdk.StoreKey
		memKey        sdk.StoreKey
		paramSubspace paramtypes.Subspace

		// The prototypical Account constructor.
		proto func() types.Account
	}
)

func NewKeeper(
	cdc codec.Marshaler,
	storeKey, memKey sdk.StoreKey,
	paramstore paramtypes.Subspace,
	proto func() types.Account) *Keeper {
	return &Keeper{
		cdc:           cdc,
		storeKey:      storeKey,
		memKey:        memKey,
		paramSubspace: paramstore,
	}
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// NewAccountWithAddress implements sdk.Keeper.
func (k Keeper) NewAccountWithAddress(ctx sdk.Context, addr string) types.Account {
	acc := k.proto()
	err := acc.SetAddress(addr)
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	err = acc.SetAccountNumber(k.GetNextAccountNumber(ctx))
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	return acc
}

// NewAccount creates a new account
func (k Keeper) NewAccount(ctx sdk.Context, acc types.Account) types.Account {
	if err := acc.SetAccountNumber(k.GetNextAccountNumber(ctx)); err != nil {
		panic(err)
	}
	return acc
}

// GetAccount implements sdk.Keeper.
func (k Keeper) GetAccount(ctx sdk.Context, addr string) types.Account {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.AddressStoreKey(addr))
	if bz == nil {
		return nil
	}
	acc, _ := k.decodeAccount(bz)
	return acc
}

// GetAllAccounts returns all accounts in the Keeper.
func (k Keeper) GetAllAccounts(ctx sdk.Context) []types.Account {
	accounts := []types.Account{}
	appendAccount := func(acc types.Account) (stop bool) {
		accounts = append(accounts, acc)
		return false
	}
	k.IterateAccounts(ctx, appendAccount)
	return accounts
}

// SetAccount implements sdk.Keeper
// allows addition of new accounts
func (k Keeper) SetAccount(ctx sdk.Context, acc types.Account) {
	addr := acc.GetAddress()
	store := ctx.KVStore(k.storeKey)
	bz, err := codec.MarshalAny(k.cdc, addr)
	if err != nil {
		panic(err)
	}
	store.Set(types.AddressStoreKey(addr), bz)
}

// RemoveAccount removes an account for the account mapper store.
// NOTE: this will cause supply invariant violation if called
func (k Keeper) RemoveAccount(ctx sdk.Context, acc types.Account) {
	addr := acc.GetAddress()
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.AddressStoreKey(addr))
}

// IterateAccounts implements sdk.Keeper.
func (k Keeper) IterateAccounts(ctx sdk.Context, process func(types.Account) (stop bool)) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.AddressStoreKeyPrefix)
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		acc, _ := k.decodeAccount(val)
		if process(acc) {
			return
		}
		iter.Next()
	}
}

// GetPubKey Returns the PubKey of the account at address
func (k Keeper) GetPubKey(ctx sdk.Context, addr string) (crypto.PubKey, error) {
	acc := k.GetAccount(ctx, addr)
	if acc == nil {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", addr)
	}
	return acc.GetPubKey(), nil
}

// GetSequence Returns the Sequence of the account at address
func (k Keeper) GetSequence(ctx sdk.Context, addr string) (uint64, error) {
	acc := k.GetAccount(ctx, addr)
	if acc == nil {
		return 0, sdkerrors.Wrapf(sdkerrors.ErrUnknownAddress, "account %s does not exist", addr)
	}
	return acc.GetSequence(), nil
}

// GetNextAccountNumber Returns and increments the global account number counter
func (k Keeper) GetNextAccountNumber(ctx sdk.Context) uint64 {
	var accNumber uint64
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GlobalAccountNumberKey)
	if bz == nil {
		// initialize the account numbers
		accNumber = 0
	} else {
		val := gogotypes.UInt64Value{}

		err := k.cdc.UnmarshalBinaryBare(bz, &val)
		if err != nil {
			panic(err)
		}

		accNumber = val.GetValue()
	}

	bz = k.cdc.MustMarshalBinaryBare(&gogotypes.UInt64Value{Value: accNumber + 1})
	store.Set(types.GlobalAccountNumberKey, bz)

	return accNumber
}

//
// proposer
//

// GetBlockProposer returns block proposer
func (k Keeper) GetBlockProposer(ctx sdk.Context) (hmCommonTypes.HeimdallAddress, bool) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has(types.ProposerKey()) {
		return hmCommonTypes.HeimdallAddress{}, false
	}

	bz := store.Get(types.ProposerKey())
	return hmCommonTypes.BytesToHeimdallAddress(bz), true
}

// SetBlockProposer sets block proposer
func (k Keeper) SetBlockProposer(ctx sdk.Context, addr hmCommonTypes.HeimdallAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ProposerKey(), addr.Bytes())
}

// RemoveBlockProposer removes block proposer from store
func (k Keeper) RemoveBlockProposer(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.ProposerKey())
}

// -----------------------------------------------------------------------------
// Params

// SetParams sets the auth module's parameters.
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	k.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the auth module's parameters.
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSubspace.GetParamSet(ctx, &params)
	return
}

// -----------------------------------------------------------------------------
// Misc.

func (k Keeper) decodeAccount(bz []byte) (types.Account, error) {
	var acc types.Account
	if err := codec.UnmarshalAny(k.cdc, &acc, bz); err != nil {
		return nil, err
	}

	return acc, nil
}
