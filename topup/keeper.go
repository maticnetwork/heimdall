package topup

import (
	"math/big"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/maticnetwork/heimdall/bank"
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
)

// ModuleCommunicator manager to access validator info
type ModuleCommunicator interface {
	// AddFeeToDividendAccount add fee to dividend account
	AddFeeToDividendAccount(ctx sdk.Context, valID hmTypes.ValidatorID, fee *big.Int) sdk.Error
	// GetValidatorFromValID get validator from validator id
	GetValidatorFromValID(ctx sdk.Context, valID hmTypes.ValidatorID) (validator hmTypes.Validator, ok bool)
}

// Keeper stores all related data
type Keeper struct {
	// The (unexposed) key used to access the store from the Context.
	key sdk.StoreKey
	// The codec codec for binary encoding/decoding of accounts.
	cdc *codec.Codec
	// code space
	codespace sdk.CodespaceType
	// param subspace
	paramSpace params.Subspace
	// bank keeper
	bk bank.Keeper
	// staking keeper
	sk staking.Keeper
	// module manager
	vm ModuleCommunicator
}

// NewKeeper create new keeper
func NewKeeper(
	cdc *codec.Codec,
	storeKey sdk.StoreKey,
	paramSpace params.Subspace,
	codespace sdk.CodespaceType,
	bankKeeper bank.Keeper,
	stakingKeeper staking.Keeper,
	vm ModuleCommunicator,
) Keeper {
	return Keeper{
		cdc:        cdc,
		key:        storeKey,
		paramSpace: paramSpace,
		codespace:  codespace,
		bk:         bankKeeper,
		sk:         stakingKeeper,
		vm:         vm,
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
func GetTopupSequenceKey(sequence big.Int) []byte {
	return append(TopupSequencePrefixKey, sequence.Bytes()...)
}

// GetTopupSequences checks if topup already exists
func (keeper Keeper) GetTopupSequences(ctx sdk.Context) (sequences []*big.Int) {
	keeper.IterateTopupSequencesAndApplyFn(ctx, func(sequence big.Int) error {
		sequences = append(sequences, &sequence)
		return nil
	})
	return
}

// IterateTopupSequencesAndApplyFn interate validators and apply the given function.
func (keeper Keeper) IterateTopupSequencesAndApplyFn(ctx sdk.Context, f func(sequence big.Int) error) {
	store := ctx.KVStore(keeper.key)

	// get sequence iterator
	iterator := sdk.KVStorePrefixIterator(store, TopupSequencePrefixKey)
	defer iterator.Close()

	// loop through validators to get valid validators
	for ; iterator.Valid(); iterator.Next() {
		sequence := new(big.Int)
		sequence.SetBytes(iterator.Key()[len(TopupSequencePrefixKey):])

		// call function and return if required
		if err := f(*sequence); err != nil {
			return
		}
	}
	return
}

// SetTopupSequence sets mapping for sequence id to bool
func (keeper Keeper) SetTopupSequence(ctx sdk.Context, sequence *big.Int) {
	store := ctx.KVStore(keeper.key)
	store.Set(GetTopupSequenceKey(*sequence), DefaultValue)
}

// HasTopupSequence checks if topup already exists
func (keeper Keeper) HasTopupSequence(ctx sdk.Context, sequence *big.Int) bool {
	store := ctx.KVStore(keeper.key)
	return store.Has(GetTopupSequenceKey(*sequence))
}
