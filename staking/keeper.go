package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	conf "github.com/maticnetwork/heimdall/helper"
	abci "github.com/tendermint/tendermint/abci/types"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *wire.Codec

	// codespace
	codespace sdk.CodespaceType
}

var (
	ValidatorsKey = []byte{0x02} // prefix for each key to a validator
)
var StakingLogger = conf.Logger.With("module", "staking")

func NewKeeper(cdc *wire.Codec, key sdk.StoreKey, codespace sdk.CodespaceType) Keeper {
	keeper := Keeper{
		storeKey:  key,
		cdc:       cdc,
		codespace: codespace,
	}
	return keeper
}

//validator type will contain address, pubkey and power
func (k Keeper) SetValidatorSet(ctx sdk.Context, validators []abci.Validator) {
	store := ctx.KVStore(k.storeKey)

	for _, validator := range validators {
		bz, err := k.cdc.MarshalBinary(validator)
		if err != nil {
			StakingLogger.Error("Error Marshalling Validator %v", err)
		}
		store.Set(GetValidatorKey(validator.Address), bz)
	}
}

// appends the validator key to address
func GetValidatorKey(address []byte) []byte {
	return append(ValidatorsKey, address...)
}

// returns all validators added for a specific validafor key
func (k Keeper) GetAllValidators(ctx sdk.Context) (validators []abci.Validator) {
	// get key
	store := ctx.KVStore(k.storeKey)
	// create iterator
	iterator := sdk.KVStorePrefixIterator(store, ValidatorsKey)

	i := 0
	for ; ; i++ {
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

// given the address returns validator info
func (k Keeper) GetValidatorInfo(ctx sdk.Context, address []byte) (validator abci.Validator) {
	// get key
	store := ctx.KVStore(k.storeKey)
	// get validator and unmarshall
	validatorBytes := store.Get(GetValidatorKey(address))
	err := k.cdc.UnmarshalBinary(validatorBytes, &validator)
	if err != nil {
		return
	}

	return validator
}

// flushes the whole validator set
func (k Keeper) FlushValidatorSet(ctx sdk.Context) {
	// get key
	store := ctx.KVStore(k.storeKey)
	// create iterator
	iterator := sdk.KVStorePrefixIterator(store, ValidatorsKey)
	i := 0
	for ; ; i++ {
		if !iterator.Valid() {
			break
		}
		addr := iterator.Key()[1:]

		var validator abci.Validator
		err := k.cdc.UnmarshalBinary(iterator.Value(), &validator)
		if err != nil {
			return
		}

		validator.Address = addr
		// make power 0
		validator.Power = int64(0)
		// marshall
		bz, err := k.cdc.MarshalBinary(validator)
		if err != nil {
			StakingLogger.Error("Error Marshalling Validator  %v", err)
		}

		store.Set(GetValidatorKey(validator.Address), bz)

		iterator.Next()
	}

	iterator.Close()
}
