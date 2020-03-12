package auth

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/params/subspace"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// AccountKeeper encodes/decodes accounts using the go-amino (binary)
// encoding/decoding library.
type AccountKeeper struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey

	// The prototypical Account constructor.
	proto func() types.Account

	// The codec codec for binary encoding/decoding of accounts.
	cdc *codec.Codec

	paramSubspace subspace.Subspace
}

// NewAccountKeeper returns a new sdk.AccountKeeper that uses go-amino to
// (binary) encode and decode concrete sdk.Accounts.
// nolint
func NewAccountKeeper(
	cdc *codec.Codec, key sdk.StoreKey, paramstore subspace.Subspace, proto func() types.Account,
) AccountKeeper {

	return AccountKeeper{
		key:           key,
		proto:         proto,
		cdc:           cdc,
		paramSubspace: paramstore.WithKeyTable(types.ParamKeyTable()),
	}
}

// Logger returns a module-specific logger.
func (ak AccountKeeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

// NewAccountWithAddress implements sdk.AccountKeeper.
func (ak AccountKeeper) NewAccountWithAddress(ctx sdk.Context, addr hmTypes.HeimdallAddress) types.Account {
	acc := ak.proto()
	err := acc.SetAddress(addr)
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	err = acc.SetAccountNumber(ak.GetNextAccountNumber(ctx))
	if err != nil {
		// Handle w/ #870
		panic(err)
	}
	return acc
}

// NewAccount creates a new account
func (ak AccountKeeper) NewAccount(ctx sdk.Context, acc types.Account) types.Account {
	if err := acc.SetAccountNumber(ak.GetNextAccountNumber(ctx)); err != nil {
		panic(err)
	}
	return acc
}

// GetAccount implements sdk.AccountKeeper.
func (ak AccountKeeper) GetAccount(ctx sdk.Context, addr hmTypes.HeimdallAddress) types.Account {
	store := ctx.KVStore(ak.key)
	bz := store.Get(types.AddressStoreKey(addr))
	if bz == nil {
		return nil
	}
	acc := ak.decodeAccount(bz)
	return acc
}

// GetAllAccounts returns all accounts in the accountKeeper.
func (ak AccountKeeper) GetAllAccounts(ctx sdk.Context) []types.Account {
	accounts := []types.Account{}
	appendAccount := func(acc types.Account) (stop bool) {
		accounts = append(accounts, acc)
		return false
	}
	ak.IterateAccounts(ctx, appendAccount)
	return accounts
}

// SetAccount implements sdk.AccountKeeper
// allows addition of new accounts
func (ak AccountKeeper) SetAccount(ctx sdk.Context, acc types.Account) {
	addr := acc.GetAddress()
	store := ctx.KVStore(ak.key)
	bz, err := ak.cdc.MarshalBinaryBare(acc)
	if err != nil {
		panic(err)
	}
	store.Set(types.AddressStoreKey(addr), bz)
}

// RemoveAccount removes an account for the account mapper store.
// NOTE: this will cause supply invariant violation if called
func (ak AccountKeeper) RemoveAccount(ctx sdk.Context, acc types.Account) {
	addr := acc.GetAddress()
	store := ctx.KVStore(ak.key)
	store.Delete(types.AddressStoreKey(addr))
}

// IterateAccounts implements sdk.AccountKeeper.
func (ak AccountKeeper) IterateAccounts(ctx sdk.Context, process func(types.Account) (stop bool)) {
	store := ctx.KVStore(ak.key)
	iter := sdk.KVStorePrefixIterator(store, types.AddressStoreKeyPrefix)
	defer iter.Close()
	for {
		if !iter.Valid() {
			return
		}
		val := iter.Value()
		acc := ak.decodeAccount(val)
		if process(acc) {
			return
		}
		iter.Next()
	}
}

// GetPubKey Returns the PubKey of the account at address
func (ak AccountKeeper) GetPubKey(ctx sdk.Context, addr hmTypes.HeimdallAddress) (crypto.PubKey, sdk.Error) {
	acc := ak.GetAccount(ctx, addr)
	if acc == nil {
		return nil, sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr))
	}
	return acc.GetPubKey(), nil
}

// GetSequence Returns the Sequence of the account at address
func (ak AccountKeeper) GetSequence(ctx sdk.Context, addr hmTypes.HeimdallAddress) (uint64, sdk.Error) {
	acc := ak.GetAccount(ctx, addr)
	if acc == nil {
		return 0, sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr))
	}
	return acc.GetSequence(), nil
}

// GetNextAccountNumber Returns and increments the global account number counter
func (ak AccountKeeper) GetNextAccountNumber(ctx sdk.Context) uint64 {
	var accNumber uint64
	store := ctx.KVStore(ak.key)
	bz := store.Get(types.GlobalAccountNumberKey)
	if bz == nil {
		accNumber = 0
	} else {
		err := ak.cdc.UnmarshalBinaryLengthPrefixed(bz, &accNumber)
		if err != nil {
			panic(err)
		}
	}

	bz = ak.cdc.MustMarshalBinaryLengthPrefixed(accNumber + 1)
	store.Set(types.GlobalAccountNumberKey, bz)

	return accNumber
}

//
// proposer
//

// GetBlockProposer returns block proposer
func (ak AccountKeeper) GetBlockProposer(ctx sdk.Context) (hmTypes.HeimdallAddress, bool) {
	store := ctx.KVStore(ak.key)
	if !store.Has(types.ProposerKey()) {
		return hmTypes.HeimdallAddress{}, false
	}

	bz := store.Get(types.ProposerKey())
	return hmTypes.BytesToHeimdallAddress(bz), true
}

// SetBlockProposer sets block proposer
func (ak AccountKeeper) SetBlockProposer(ctx sdk.Context, addr hmTypes.HeimdallAddress) {
	store := ctx.KVStore(ak.key)
	store.Set(types.ProposerKey(), addr.Bytes())
}

// RemoveBlockProposer removes block proposer from store
func (ak AccountKeeper) RemoveBlockProposer(ctx sdk.Context) {
	store := ctx.KVStore(ak.key)
	store.Delete(types.ProposerKey())
}

// -----------------------------------------------------------------------------
// Params

// SetParams sets the auth module's parameters.
func (ak AccountKeeper) SetParams(ctx sdk.Context, params types.Params) {
	ak.paramSubspace.SetParamSet(ctx, &params)
}

// GetParams gets the auth module's parameters.
func (ak AccountKeeper) GetParams(ctx sdk.Context) (params types.Params) {
	ak.paramSubspace.GetParamSet(ctx, &params)
	return
}

// -----------------------------------------------------------------------------
// Misc.

func (ak AccountKeeper) decodeAccount(bz []byte) (acc types.Account) {
	err := ak.cdc.UnmarshalBinaryBare(bz, &acc)
	if err != nil {
		panic(err)
	}
	return
}
