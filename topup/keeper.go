package topup

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/maticnetwork/heimdall/bank"
	"github.com/maticnetwork/heimdall/chainmanager"
	"github.com/maticnetwork/heimdall/params/subspace"
	"github.com/maticnetwork/heimdall/staking"
	"github.com/maticnetwork/heimdall/topup/types"
	"github.com/tendermint/tendermint/libs/log"
)

var (
	// DefaultValue default value
	DefaultValue = []byte{0x01}
	// TopupSequencePrefixKey represents topup sequence prefix key
	TopupSequencePrefixKey = []byte{0x81}
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
func (keeper Keeper) GetTopupSequences(ctx sdk.Context) (sequences []string) {
	keeper.IterateTopupSequencesAndApplyFn(ctx, func(sequence string) error {
		sequences = append(sequences, sequence)
		return nil
	})
	return
}

// IterateTopupSequencesAndApplyFn interate validators and apply the given function.
func (keeper Keeper) IterateTopupSequencesAndApplyFn(ctx sdk.Context, f func(sequence string) error) {
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
func (keeper Keeper) SetTopupSequence(ctx sdk.Context, sequence string) {
	store := ctx.KVStore(keeper.key)
	store.Set(GetTopupSequenceKey(sequence), DefaultValue)
}

// HasTopupSequence checks if topup already exists
func (keeper Keeper) HasTopupSequence(ctx sdk.Context, sequence string) bool {
	store := ctx.KVStore(keeper.key)
	return store.Has(GetTopupSequenceKey(sequence))
}
