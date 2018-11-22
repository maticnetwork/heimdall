package staking

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/tendermint/tendermint/crypto"
)

type Keeper struct {
	storeKey         sdk.StoreKey
	cdc              *codec.Codec
	checkpointKeeper checkpoint.Keeper

	// codespace
	codespace sdk.CodespaceType
}

var (
	ValidatorsKey = []byte{0x01} // prefix for each key to a validator
)

// getValidatorKey drafts the validator key for addresses
func getValidatorKey(address []byte) []byte {
	return append(ValidatorsKey, address...)
}

type Validator struct {
	Address    common.Address
	StartEpoch int64
	EndEpoch   int64
	Pubkey     crypto.PubKey
	Power      int64 // aka Amount
	Signer     common.Address
}

// create empty validator without pubkey
func CreateEmptyValidator() Validator {
	validator := Validator{
		Address:    common.HexToAddress(""),
		StartEpoch: int64(0),
		EndEpoch:   int64(0),
		Power:      int64(0),
		Signer:     common.HexToAddress(""),
	}
	return validator
}

// NewKeeper creates new keeper for staking
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, codespace sdk.CodespaceType) Keeper {
	keeper := Keeper{
		storeKey:  key,
		cdc:       cdc,
		codespace: codespace,
	}
	return keeper
}

// Add validator indexed with address
func (k Keeper) AddValidator(ctx sdk.Context, validator Validator) {
	store := ctx.KVStore(k.storeKey)

	// marshall validator
	bz, err := k.cdc.MarshalBinary(validator)
	if err != nil {
		StakingLogger.Error("Error marshalling validator", "error", err, "ValidatorAddress", validator.Address)
		panic(err)
	}

	// store validator with address prefixed with validator key as index
	store.Set(getValidatorKey(validator.Pubkey.Address().Bytes()), bz)
}

// GetAllValidators returns all validators added for a specific validator key
func (k Keeper) GetAllCurrentValidators(ctx sdk.Context) (validators []Validator) {
	store := ctx.KVStore(k.storeKey)

	// create iterator to iterate with Validator Key prefix
	iterator := sdk.KVStorePrefixIterator(store, ValidatorsKey)
	defer iterator.Close()

	// get current ACK Count aka epoch
	ACKs := k.checkpointKeeper.GetACKCount(ctx)

	// loop through validators to get valid validators
	for i := 0; ; i++ {
		if !iterator.Valid() {
			break
		}

		// unmarshall validator
		var validator Validator
		err := k.cdc.UnmarshalBinary(iterator.Value(), &validator)
		if err != nil {
			return
		}

		// check if validator is valid for current epoch
		if validator.StartEpoch >= int64(ACKs) && validator.EndEpoch <= int64(ACKs) {
			// add validator to validators
			validators = append(validators, validator)
		}

		// increment iterator
		iterator.Next()
	}
	return
}
