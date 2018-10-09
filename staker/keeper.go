package staker

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/wire"
	abci "github.com/tendermint/tendermint/abci/types"

	"fmt"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *wire.Codec
	//validatorSet sdk.ValidatorSet

	// codespace
	codespace sdk.CodespaceType
}

var (
	ValidatorsKey = []byte{0x02} // prefix for each key to a validator
)

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
			fmt.Println("error %v", err)
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
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, ValidatorsKey)

	i := 0
	for ; ; i++ {
		if !iterator.Valid() {
			fmt.Println("BREAK BREAK BREAK ")
			break
		}

		addr := iterator.Key()[1:]
		//validator := types.MustUnmarshalValidator(k.cdc, addr, iterator.Value())
		var validator abci.Validator
		err := k.cdc.UnmarshalBinary(iterator.Value(), &validator)
		if err != nil {
			return
		}
		validator.Address = addr

		validators = append(validators, validator)
		iterator.Next()
	}
	iterator.Close()
	return validators
}

// given the address returns validator info
func (k Keeper) GetValidatorInfo(ctx sdk.Context, address []byte) (validator abci.Validator) {
	store := ctx.KVStore(k.storeKey)
	validator := store.Get(GetValidatorKey(address))
	return validator

}

// flushes the whole validator set
func (k Keeper) FlushValidatorSet(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, ValidatorsKey)
	i := 0
	for ; ; i++ {
		if !iterator.Valid() {
			break
		}
		addr := iterator.Key()[1:]
		//validator := types.MustUnmarshalValidator(k.cdc, addr, iterator.Value())
		var validator abci.Validator
		err := k.cdc.UnmarshalBinary(iterator.Value(), &validator)
		if err != nil {
			return
		}
		validator.Address = addr
		validator.Power = int64(0)
		bz, err := k.cdc.MarshalBinary(validator)
		if err != nil {
			fmt.Println("error %v", err)
		}
		store.Set(GetValidatorKey(validator.Address), bz)
		iterator.Next()
	}
	iterator.Close()
}
