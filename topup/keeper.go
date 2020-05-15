package topup

import (
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/bank"
	"github.com/maticnetwork/heimdall/chainmanager"
	"github.com/maticnetwork/heimdall/params/subspace"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/topup/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/tendermint/tendermint/libs/log"
)

var (
	// DefaultValue default value
	DefaultValue = []byte{0x01}
	// TopupSequencePrefixKey represents topup sequence prefix key
	TopupSequencePrefixKey = []byte{0x81}

	DividendAccountMapKey = []byte{0x82} // prefix for each key for Dividend Account Map
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
	// chain keeper
	chainKeeper chainmanager.Keeper
	// bank keeper
	bk bank.Keeper
	// staking keeper
	sk staking.Keeper
}

// NewKeeper create new keeper
func NewKeeper(
	cdc *codec.Codec,
	storeKey sdk.StoreKey,
	paramSpace subspace.Subspace,
	codespace sdk.CodespaceType,
	chainKeeper chainmanager.Keeper,
	bankKeeper bank.Keeper,
	stakingKeeper staking.Keeper,
) Keeper {
	return Keeper{
		cdc:         cdc,
		key:         storeKey,
		paramSpace:  paramSpace,
		codespace:   codespace,
		chainKeeper: chainKeeper,
		bk:          bankKeeper,
		sk:          stakingKeeper,
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
// Topup methods
//

// GetTopupSequenceKey drafts topup sequence for address
func GetTopupSequenceKey(sequence string) []byte {
	return append(TopupSequencePrefixKey, []byte(sequence)...)
}

// GetTopupSequences checks if topup already exists
func (keeper *Keeper) GetTopupSequences(ctx sdk.Context) (sequences []string) {
	keeper.IterateTopupSequencesAndApplyFn(ctx, func(sequence string) error {
		sequences = append(sequences, sequence)
		return nil
	})
	return
}

// IterateTopupSequencesAndApplyFn interate validators and apply the given function.
func (keeper *Keeper) IterateTopupSequencesAndApplyFn(ctx sdk.Context, f func(sequence string) error) {
	store := ctx.KVStore(keeper.key)

	// get sequence iterator
	iterator := sdk.KVStorePrefixIterator(store, TopupSequencePrefixKey)
	defer iterator.Close()

	// loop through validators to get valid validators
	for ; iterator.Valid(); iterator.Next() {
		sequence := string(iterator.Key()[len(TopupSequencePrefixKey):])

		// call function and return if required
		if err := f(sequence); err != nil {
			return
		}
	}
}

// SetTopupSequence sets mapping for sequence id to bool
func (keeper *Keeper) SetTopupSequence(ctx sdk.Context, sequence string) {
	store := ctx.KVStore(keeper.key)
	store.Set(GetTopupSequenceKey(sequence), DefaultValue)
}

// HasTopupSequence checks if topup already exists
func (keeper *Keeper) HasTopupSequence(ctx sdk.Context, sequence string) bool {
	store := ctx.KVStore(keeper.key)
	return store.Has(GetTopupSequenceKey(sequence))
}

// GetDividendAccountMapKey returns dividend account map
func GetDividendAccountMapKey(address []byte) []byte {
	return append(DividendAccountMapKey, address...)
}

// AddDividendAccount adds DividendAccount index with DividendID
func (k *Keeper) AddDividendAccount(ctx sdk.Context, dividendAccount hmTypes.DividendAccount) error {
	store := ctx.KVStore(k.key)
	// marshall dividend account
	bz, err := hmTypes.MarshallDividendAccount(k.cdc, dividendAccount)
	if err != nil {
		return err
	}

	store.Set(GetDividendAccountMapKey(dividendAccount.User.Bytes()), bz)
	k.Logger(ctx).Debug("DividendAccount Stored", "key", hex.EncodeToString(GetDividendAccountMapKey(dividendAccount.User.Bytes())), "dividendAccount", dividendAccount.String())
	return nil
}

// GetDividendAccountByAddress will return DividendAccount of user
func (k *Keeper) GetDividendAccountByAddress(ctx sdk.Context, address hmTypes.HeimdallAddress) (dividendAccount hmTypes.DividendAccount, err error) {

	// check if dividend account exists
	if !k.CheckIfDividendAccountExists(ctx, address) {
		return dividendAccount, errors.New("Dividend Account not found")
	}

	// Get DividendAccount key
	store := ctx.KVStore(k.key)
	key := GetDividendAccountMapKey(address.Bytes())

	// unmarshall dividend account and return
	dividendAccount, err = hmTypes.UnMarshallDividendAccount(k.cdc, store.Get(key))
	if err != nil {
		return dividendAccount, err
	}

	return dividendAccount, nil
}

// CheckIfDividendAccountExists will return true if dividend account exists
func (k *Keeper) CheckIfDividendAccountExists(ctx sdk.Context, userAddr hmTypes.HeimdallAddress) (ok bool) {
	store := ctx.KVStore(k.key)
	key := GetDividendAccountMapKey(userAddr.Bytes())
	return store.Has(key)
}

// GetAllDividendAccounts returns all DividendAccountss
func (k *Keeper) GetAllDividendAccounts(ctx sdk.Context) (dividendAccounts []hmTypes.DividendAccount) {
	// iterate through dividendAccounts and create dividendAccounts update array
	k.IterateDividendAccountsByPrefixAndApplyFn(ctx, DividendAccountMapKey, func(dividendAccount hmTypes.DividendAccount) error {
		// append to list of dividendUpdates
		dividendAccounts = append(dividendAccounts, dividendAccount)
		return nil
	})

	return
}

// AddFeeToDividendAccount adds fee to dividend account for withdrawal
func (k *Keeper) AddFeeToDividendAccount(ctx sdk.Context, userAddress hmTypes.HeimdallAddress, fee *big.Int) sdk.Error {
	// Get or create dividend account
	var dividendAccount hmTypes.DividendAccount

	if k.CheckIfDividendAccountExists(ctx, hmTypes.HeimdallAddress(userAddress)) {
		dividendAccount, _ = k.GetDividendAccountByAddress(ctx, userAddress)
	} else {
		dividendAccount = hmTypes.DividendAccount{
			User:      userAddress,
			FeeAmount: big.NewInt(0).String(),
		}
	}

	// update fee
	oldFee, _ := big.NewInt(0).SetString(dividendAccount.FeeAmount, 10)
	totalFee := big.NewInt(0).Add(oldFee, fee).String()
	dividendAccount.FeeAmount = totalFee

	k.Logger(ctx).Info("Dividend Account fee of validator ", "User", dividendAccount.User, "Fee", dividendAccount.FeeAmount)
	if err := k.AddDividendAccount(ctx, dividendAccount); err != nil {
		k.Logger(ctx).Error("AddFeeToDividendAccount | AddDividendAccount", "error", err)
	}
	return nil
}

// IterateDividendAccountsByPrefixAndApplyFn iterate dividendAccounts and apply the given function.
func (k *Keeper) IterateDividendAccountsByPrefixAndApplyFn(ctx sdk.Context, prefix []byte, f func(dividendAccount hmTypes.DividendAccount) error) {
	store := ctx.KVStore(k.key)

	// get validator iterator
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	// loop through dividendAccounts
	for ; iterator.Valid(); iterator.Next() {
		// unmarshall dividendAccount
		dividendAccount, _ := hmTypes.UnMarshallDividendAccount(k.cdc, iterator.Value())
		// call function and return if required
		if err := f(dividendAccount); err != nil {
			return
		}
	}
}
