package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/codec"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.Codec

	// codespace
	codespace sdk.CodespaceType
}

var (
	ValidatorsKey = []byte{0x02} // prefix for each key to a validator
)

// NewKeeper creates new keeper for staking
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, codespace sdk.CodespaceType) Keeper {
	keeper := Keeper{
		storeKey:  key,
		cdc:       cdc,
		codespace: codespace,
	}
	return keeper
}

// SetValidatorSet validator type will contain address, pubkey and power
func (k Keeper) SetValidatorSet(ctx sdk.Context, validators []abci.ValidatorUpdate) {
	store := ctx.KVStore(k.storeKey)

	for _, validator := range validators {
		bz, err := k.cdc.MarshalBinary(validator)
		if err != nil {
			StakingLogger.Error("Error marshalling validator", "error", err)
			panic(err)
		}
		pubkey,err := types.PB2TM.PubKey(validator.GetPubKey())
		if err!=nil{
			StakingLogger.Error("Error converting to cryptoPubkey","ValidatorPubkey",validator.GetPubKey(),"ValidatorPower",validator.Power )
		}
		store.Set(getValidatorKey(pubkey.Address().Bytes()), bz)
	}
}

// getValidatorKey drafts the validator key for addresses
func getValidatorKey(address []byte) []byte {
	return append(ValidatorsKey, address...)
}

// GetAllValidators returns all validators added for a specific validafor key
func (k Keeper) GetAllValidators(ctx sdk.Context) (validators []abci.Validator) {
	// get key
	store := ctx.KVStore(k.storeKey)
	// create iterator
	iterator := sdk.KVStorePrefixIterator(store, ValidatorsKey)

	for i := 0; ; i++ {
		if !iterator.Valid() {
			break
		}

		addr := iterator.Key()[1:]
		// unmarshall validator
		var validator abci.Validator
		err := k.cdc.UnmarshalBinary(iterator.Value(), &validator)
		if err != nil {
			return
		}
		// set Validator Address
		validator.Address = addr
		// append to list
		validators = append(validators, validator)

		iterator.Next()
	}
	iterator.Close()

	return validators
}

// GetValidatorInfo returns validator info for given the address
func (k Keeper) GetValidatorInfo(ctx sdk.Context, address []byte) (validator abci.Validator) {
	// get key
	store := ctx.KVStore(k.storeKey)
	// get validator and unmarshall
	validatorBytes := store.Get(getValidatorKey(address))
	err := k.cdc.UnmarshalBinary(validatorBytes, &validator)
	if err != nil {
		return
	}

	return validator
}

// FlushValidatorSet flushes the whole validator set
func (k Keeper) FlushValidatorSet(ctx sdk.Context) {
	// get key
	store := ctx.KVStore(k.storeKey)
	// create iterator
	iterator := sdk.KVStorePrefixIterator(store, ValidatorsKey)
	for i := 0; ; i++ {
		if !iterator.Valid() {
			break
		}
		addr := iterator.Key()[1:]

		var validator abci.Validator
		if err := k.cdc.UnmarshalBinary(iterator.Value(), &validator); err != nil {
			StakingLogger.Error("Error unmarshalling validator while flushing", "error", err)
			panic(err)
		}

		validator.Address = addr
		// make power 0
		validator.Power = int64(0)
		// marshall
		bz, err := k.cdc.MarshalBinary(validator)
		if err != nil {
			StakingLogger.Error("Error marshalling validator while flushing", "error", err)
			panic(err)
		}

		store.Set(getValidatorKey(validator.Address), bz)
		iterator.Next()
	}

	iterator.Close()
}
