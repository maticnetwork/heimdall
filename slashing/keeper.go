package slashing

import (
	"fmt"
	"math/big"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/maticnetwork/heimdall/chainmanager"
	"github.com/maticnetwork/heimdall/params/subspace"
	"github.com/maticnetwork/heimdall/slashing/types"
	"github.com/maticnetwork/heimdall/staking"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// Keeper of the slashing store
type Keeper struct {
	cdc      *codec.Codec
	storeKey sdk.StoreKey
	sk       staking.Keeper
	// codespace
	codespace  sdk.CodespaceType
	paramSpace subspace.Subspace

	// chain manager keeper
	chainKeeper chainmanager.Keeper
}

// NewKeeper creates a slashing keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, sk staking.Keeper, paramSpace subspace.Subspace, codespace sdk.CodespaceType, chainKeeper chainmanager.Keeper) Keeper {
	return Keeper{
		storeKey:    key,
		cdc:         cdc,
		sk:          sk,
		paramSpace:  paramSpace.WithKeyTable(types.ParamKeyTable()),
		codespace:   codespace,
		chainKeeper: chainKeeper,
	}
}

// Codespace returns the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}

// GetValidatorSigningInfo retruns the ValidatorSigningInfo for a specific validator
// ConsAddress
func (k Keeper) GetValidatorSigningInfo(ctx sdk.Context, valID hmTypes.ValidatorID) (info hmTypes.ValidatorSigningInfo, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValidatorSigningInfoKey(valID.Bytes()))
	if bz == nil {
		found = false
		return
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &info)
	found = true
	return
}

// HasValidatorSigningInfo returns if a given validator has signing information
// persited.
func (k Keeper) HasValidatorSigningInfo(ctx sdk.Context, valID hmTypes.ValidatorID) bool {
	_, ok := k.GetValidatorSigningInfo(ctx, valID)
	return ok
}

// SetValidatorSigningInfo sets the validator signing info to a consensus address key
func (k Keeper) SetValidatorSigningInfo(ctx sdk.Context, valID hmTypes.ValidatorID, info hmTypes.ValidatorSigningInfo) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&info)
	store.Set(types.GetValidatorSigningInfoKey(valID.Bytes()), bz)
}

// IterateValidatorSigningInfos iterates over the stored ValidatorSigningInfo
func (k Keeper) IterateValidatorSigningInfos(ctx sdk.Context,
	handler func(valID hmTypes.ValidatorID, info hmTypes.ValidatorSigningInfo) (stop bool)) {

	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.ValidatorSigningInfoKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var info hmTypes.ValidatorSigningInfo
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &info)
		if handler(info.ValID, info) {
			break
		}
	}
}

// signing info bit array

// GetValidatorMissedBlockBitArray gets the bit for the missed blocks array
func (k Keeper) GetValidatorMissedBlockBitArray(ctx sdk.Context, valID hmTypes.ValidatorID, index int64) bool {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetValidatorMissedBlockBitArrayKey(valID.Bytes(), index))
	var missed gogotypes.BoolValue
	if bz == nil {
		// lazy: treat empty key as not missed
		return false
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &missed)

	return missed.Value
}

// IterateValidatorMissedBlockBitArray iterates over the signed blocks window
// and performs a callback function
func (k Keeper) IterateValidatorMissedBlockBitArray(ctx sdk.Context,
	valID hmTypes.ValidatorID, handler func(index int64, missed bool) (stop bool)) {

	store := ctx.KVStore(k.storeKey)
	index := int64(0)
	params := k.GetParams(ctx)
	// Array may be sparse
	for ; index < params.SignedBlocksWindow; index++ {
		var missed gogotypes.BoolValue
		bz := store.Get(types.GetValidatorMissedBlockBitArrayKey(valID.Bytes(), index))
		if bz == nil {
			continue
		}

		k.cdc.MustUnmarshalBinaryBare(bz, &missed)
		if handler(index, missed.Value) {
			break
		}
	}
}

// SetValidatorMissedBlockBitArray sets the bit that checks if the validator has
// missed a block in the current window
func (k Keeper) SetValidatorMissedBlockBitArray(ctx sdk.Context, valID hmTypes.ValidatorID, index int64, missed bool) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.BoolValue{Value: missed})
	store.Set(types.GetValidatorMissedBlockBitArrayKey(valID.Bytes(), index), bz)
}

// clearValidatorMissedBlockBitArray deletes every instance of ValidatorMissedBlockBitArray in the store
func (k Keeper) clearValidatorMissedBlockBitArray(ctx sdk.Context, valID hmTypes.ValidatorID) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.GetValidatorMissedBlockBitArrayPrefixKey(valID.Bytes()))
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		store.Delete(iter.Key())
	}
}

// MinSignedPerWindow - minimum blocks signed per window
func (k Keeper) MinSignedPerWindow(ctx sdk.Context) int64 {
	var minSignedPerWindow sdk.Dec
	params := k.GetParams(ctx)
	// minSignedPerWindow = percent
	minSignedPerWindow = params.MinSignedPerWindow
	signedBlocksWindow := params.SignedBlocksWindow

	// NOTE: RoundInt64 will never panic as minSignedPerWindow is
	//       less than 1.
	return minSignedPerWindow.MulInt64(signedBlocksWindow).RoundInt64()
}

/*
// AddPubkey sets a address-pubkey relation
func (k Keeper) AddPubkey(ctx sdk.Context, pubkey crypto.PubKey) {
	addr := pubkey.Address()

	pkStr, err := sdk.Bech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, pubkey)
	if err != nil {
		panic(fmt.Errorf("error while setting address-pubkey relation: %s", addr))
	}

	k.setAddrPubkeyRelation(ctx, addr, pkStr)
}

// GetPubkey returns the pubkey from the adddress-pubkey relation
func (k Keeper) GetPubkey(ctx sdk.Context, address crypto.Address) (crypto.PubKey, error) {
	store := ctx.KVStore(k.storeKey)

	var pubkey gogotypes.StringValue
	err := k.cdc.UnmarshalBinaryBare(store.Get(types.GetAddrPubkeyRelationKey(address)), &pubkey)
	if err != nil {
		return nil, fmt.Errorf("address %s not found", sdk.ConsAddress(address))
	}

	pkStr, err := sdk.GetPubKeyFromBech32(sdk.Bech32PubKeyTypeConsPub, pubkey.Value)
	if err != nil {
		return pkStr, err
	}

	return pkStr, nil
} */

// Slash attempts to slash a validator. The slash is delegated to the Slashing
// module to make the necessary validator changes.
func (k Keeper) Slash(ctx sdk.Context, valID hmTypes.ValidatorID, fraction sdk.Dec, power, distributionHeight int64) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSlash,
			sdk.NewAttribute(types.AttributeKeyValID, valID.String()),
			sdk.NewAttribute(types.AttributeKeyPower, fmt.Sprintf("%d", power)),
			sdk.NewAttribute(types.AttributeKeyReason, types.AttributeValueDoubleSign),
		),
	)

	// k.sk.Slash(ctx, addr, distributionHeight, power, fraction)

}

// Jail attempts to jail a validator. The slash is delegated to the Slashing module
// to make the necessary validator changes.
func (k Keeper) Jail(ctx sdk.Context, valID hmTypes.ValidatorID) {
	// ctx.EventManager().EmitEvent(
	// 	sdk.NewEvent(
	// 		types.EventTypeSlash,
	// 		sdk.NewAttribute(types.AttributeKeyJailed, hmTypes.BytesToHeimdallAddress(addr).String()),
	// 	),
	// )

	// k.sk.Jail(ctx, addr)
}

/*
func (k Keeper) setAddrPubkeyRelation(ctx sdk.Context, addr crypto.Address, pubkey string) {
	store := ctx.KVStore(k.storeKey)

	bz := k.cdc.MustMarshalBinaryBare(&gogotypes.StringValue{Value: pubkey})
	store.Set(types.GetAddrPubkeyRelationKey(addr), bz)
}

func (k Keeper) deleteAddrPubkeyRelation(ctx sdk.Context, addr crypto.Address) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetAddrPubkeyRelationKey(addr))
}
*/

// -----------------------------------------------------------------------------
// Params

// SetParams sets the slashing module's parameters.
func (k *Keeper) SetParams(ctx sdk.Context, params types.Params) {
	fmt.Println("Setting params")
	k.paramSpace.SetParamSet(ctx, &params)
}

// GetParams gets the slashing module's parameters.
func (k *Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	k.paramSpace.GetParamSet(ctx, &params)
	return
}

// Slashing Info api's

func (k *Keeper) SlashInterim(ctx sdk.Context, valID hmTypes.ValidatorID, amount string) {
	// add slash to buffer
	valSlashingInfo, found := k.GetBufferValSlashingInfo(ctx, valID)
	if found {
		// Add or Update Slash Amount
		prevAmount, _ := big.NewInt(0).SetString(valSlashingInfo.SlashedAmount, 10)
		amountToAdd, _ := big.NewInt(0).SetString(amount, 10)
		updatedSlashAmount := big.NewInt(0).Add(prevAmount, amountToAdd)
		valSlashingInfo.SlashedAmount = updatedSlashAmount.String()
	} else {
		// create slashing info
		valSlashingInfo = hmTypes.NewValidatorSlashingInfo(valID, amount, false)
	}

	k.SetBufferValSlashingInfo(ctx, valID, valSlashingInfo)
	k.UpdateTotalSlashedAmount(ctx, amount)
}

// GetBufferValSlashingInfo gets the validator slashing info for a validator ID key
func (k *Keeper) GetBufferValSlashingInfo(ctx sdk.Context, valId hmTypes.ValidatorID) (info hmTypes.ValidatorSlashingInfo, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetBufferValSlashingInfoKey(valId.Bytes()))
	if bz == nil {
		found = false
		return
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &info)
	found = true
	return
}

// SetBufferValSlashingInfo sets the validator slashing info to a validator ID key
func (k Keeper) SetBufferValSlashingInfo(ctx sdk.Context, valID hmTypes.ValidatorID, info hmTypes.ValidatorSlashingInfo) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&info)
	store.Set(types.GetBufferValSlashingInfoKey(valID.Bytes()), bz)
}

// RemoveBufferValSlashingInfo removes the validator slashing info for a validator ID key
func (k Keeper) RemoveBufferValSlashingInfo(ctx sdk.Context, valID hmTypes.ValidatorID) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetBufferValSlashingInfoKey(valID.Bytes()))
}

// IterateBufferValSlashingInfos iterates over the stored ValidatorSlashingInfo
func (k Keeper) IterateBufferValSlashingInfos(ctx sdk.Context,
	handler func(slashingInfo hmTypes.ValidatorSlashingInfo) (stop bool)) {

	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.BufferValSlashingInfoKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var slashingInfo hmTypes.ValidatorSlashingInfo
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &slashingInfo)
		if handler(slashingInfo) {
			break
		}
	}
}

// FlushBufferValSlashingInfos removes all validator slashing infos in buffer
func (k *Keeper) FlushBufferValSlashingInfos(ctx sdk.Context) {
	// iterate through validator slashing info and create validator slashing info update array
	k.IterateBufferValSlashingInfosAndApplyFn(ctx, func(valSlashingInfo hmTypes.ValidatorSlashingInfo) error {
		// remove from buffer data
		k.RemoveBufferValSlashingInfo(ctx, valSlashingInfo.ID)
		return nil
	})
	return
}

// FlushBufferValSlashingInfos removes all validator slashing infos in buffer
func (k *Keeper) FlushTotalSlashedAmount(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	// remove from store
	if store.Has(types.TotalSlashedAmountKey) {
		store.Delete(types.TotalSlashedAmountKey)
	}
}

// IterateBufferValSlashingInfosAndApplyFn interate ValidatorSlashingInfo and apply the given function.
func (k *Keeper) IterateBufferValSlashingInfosAndApplyFn(ctx sdk.Context, f func(slashingInfo hmTypes.ValidatorSlashingInfo) error) {
	store := ctx.KVStore(k.storeKey)

	// get validator iterator
	iterator := sdk.KVStorePrefixIterator(store, types.BufferValSlashingInfoKey)
	defer iterator.Close()

	// loop through validators to get valid validators
	for ; iterator.Valid(); iterator.Next() {
		// unmarshall validator
		slashingInfo, _ := hmTypes.UnmarshallValSlashingInfo(k.cdc, iterator.Value())
		// call function and return if required
		if err := f(slashingInfo); err != nil {
			return
		}
	}
}

// GetBufferValSlashingInfos returns all validator slashing infos in buffer
func (k *Keeper) GetBufferValSlashingInfos(ctx sdk.Context) (valSlashingInfos []*hmTypes.ValidatorSlashingInfo) {
	// iterate through validators and create validator update array
	k.IterateBufferValSlashingInfosAndApplyFn(ctx, func(valSlashingInfo hmTypes.ValidatorSlashingInfo) error {
		// append to list of valSlashingInfos
		valSlashingInfos = append(valSlashingInfos, &valSlashingInfo)
		return nil
	})

	return
}

func (k Keeper) UpdateTotalSlashedAmount(ctx sdk.Context, amount string) {
	store := ctx.KVStore(k.storeKey)
	slashedAmount, _ := big.NewInt(0).SetString(amount, 10)
	if store.Has(types.TotalSlashedAmountKey) {
		bz := store.Get(types.TotalSlashedAmountKey)
		prevAmountStr := string(bz)
		prevAmount, _ := big.NewInt(0).SetString(prevAmountStr, 10)
		slashedAmount = big.NewInt(0).Add(prevAmount, slashedAmount)
	}

	store.Set(types.TotalSlashedAmountKey, []byte(slashedAmount.String()))
	k.Logger(ctx).Debug("Updated Total Slashed Amount ", "amount", slashedAmount)

	// -slashing. emit event if total amount exceed limit
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeSlashLimit,
			sdk.NewAttribute(types.AttributeKeySlashedAmount, fmt.Sprintf("%d", slashedAmount)),
		),
	)
}

// GetTickValSlashingInfo gets the validator slashing info for a validator ID key
func (k *Keeper) GetTickValSlashingInfo(ctx sdk.Context, valId hmTypes.ValidatorID) (info hmTypes.ValidatorSlashingInfo, found bool) {
	store := ctx.KVStore(k.storeKey)

	bz := store.Get(types.GetTickValSlashingInfoKey(valId.Bytes()))
	if bz == nil {
		found = false
		return
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &info)
	found = true
	return
}

// GetTickValSlashingInfos returns all validator slashing infos in tick
func (k *Keeper) GetTickValSlashingInfos(ctx sdk.Context) (valSlashingInfos []*hmTypes.ValidatorSlashingInfo) {
	// iterate through validators and create slashing info update array
	k.IterateTickValSlashingInfosAndApplyFn(ctx, func(valSlashingInfo hmTypes.ValidatorSlashingInfo) error {
		// append to list of valSlashingInfos
		valSlashingInfos = append(valSlashingInfos, &valSlashingInfo)
		return nil
	})

	return
}

// SetTickValSlashingInfo sets the validator slashing info to a validator ID key
func (k Keeper) SetTickValSlashingInfo(ctx sdk.Context, valID hmTypes.ValidatorID, info hmTypes.ValidatorSlashingInfo) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(&info)
	store.Set(types.GetTickValSlashingInfoKey(valID.Bytes()), bz)
}

// RemoveTickValSlashingInfo removes the validator slashing info for a validator ID key
func (k Keeper) RemoveTickValSlashingInfo(ctx sdk.Context, valID hmTypes.ValidatorID) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetTickValSlashingInfoKey(valID.Bytes()))
}

// IterateTickValSlashingInfos iterates over the stored ValidatorSlashingInfo
func (k Keeper) IterateTickValSlashingInfos(ctx sdk.Context,
	handler func(slashingInfo hmTypes.ValidatorSlashingInfo) (stop bool)) {

	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.TickValSlashingInfoKey)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		var slashingInfo hmTypes.ValidatorSlashingInfo
		k.cdc.MustUnmarshalBinaryBare(iter.Value(), &slashingInfo)
		if handler(slashingInfo) {
			break
		}
	}
}

// CopyValSlashingInfosToTickData copies all validator slashing infos in buffer to tickdata
func (k *Keeper) CopyBufferValSlashingInfosToTickData(ctx sdk.Context) {
	// iterate through validators and create validator slashing info update array
	k.IterateBufferValSlashingInfosAndApplyFn(ctx, func(valSlashingInfo hmTypes.ValidatorSlashingInfo) error {
		// store to tick data
		k.SetTickValSlashingInfo(ctx, valSlashingInfo.ID, valSlashingInfo)
		return nil
	})

	return
}

// IterateTickValSlashingInfosAndApplyFn interate ValidatorSlashingInfo and apply the given function.
func (k *Keeper) IterateTickValSlashingInfosAndApplyFn(ctx sdk.Context, f func(slashingInfo hmTypes.ValidatorSlashingInfo) error) {
	store := ctx.KVStore(k.storeKey)

	// get validator iterator
	iterator := sdk.KVStorePrefixIterator(store, types.TickValSlashingInfoKey)
	defer iterator.Close()

	// loop through validators to get valid validators
	for ; iterator.Valid(); iterator.Next() {
		// unmarshall validator
		slashingInfo, _ := hmTypes.UnmarshallValSlashingInfo(k.cdc, iterator.Value())
		// call function and return if required
		if err := f(slashingInfo); err != nil {
			// Error slashing validator
			k.Logger(ctx).Error("Error slashing the validator", "error", err)
			return
		}
	}
}

// SlashAndJailTickValSlashingInfos reduces power of all validator slashing infos in tick data
func (k *Keeper) SlashAndJailTickValSlashingInfos(ctx sdk.Context) {
	// iterate through validator slashing info and create validator slashing info update array
	k.IterateTickValSlashingInfosAndApplyFn(ctx, func(valSlashingInfo hmTypes.ValidatorSlashingInfo) error {
		err := k.sk.Slash(ctx, valSlashingInfo)
		return err
	})
	return
}

// FlushTickValSlashingInfos removes all validator slashing infos in last Tick
func (k *Keeper) FlushTickValSlashingInfos(ctx sdk.Context) {
	// iterate through validator slashing info and create validator slashing info update array
	k.IterateTickValSlashingInfosAndApplyFn(ctx, func(valSlashingInfo hmTypes.ValidatorSlashingInfo) error {
		// remove from tick data
		k.RemoveTickValSlashingInfo(ctx, valSlashingInfo.ID)
		return nil
	})
	return
}

//
// Slashing sequence
//

// SetSlashingSequence sets Slashing sequence
func (k *Keeper) SetSlashingSequence(ctx sdk.Context, sequence string) {
	store := ctx.KVStore(k.storeKey)

	store.Set(types.GetSlashingSequenceKey(sequence), types.DefaultValue)
}

// HasSlashingSequence checks if Slashing sequence already exists
func (k *Keeper) HasSlashingSequence(ctx sdk.Context, sequence string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetSlashingSequenceKey(sequence))
}

// GetSlashingSequences checks if Slashing already exists
func (k *Keeper) GetSlashingSequences(ctx sdk.Context) (sequences []string) {
	k.IterateSlashingSequencesAndApplyFn(ctx, func(sequence string) error {
		sequences = append(sequences, sequence)
		return nil
	})
	return
}

// IterateSlashingSequencesAndApplyFn interate validators and apply the given function.
func (k *Keeper) IterateSlashingSequencesAndApplyFn(ctx sdk.Context, f func(sequence string) error) {
	store := ctx.KVStore(k.storeKey)

	// get sequence iterator
	iterator := sdk.KVStorePrefixIterator(store, types.SlashingSequenceKey)
	defer iterator.Close()

	// loop through validators to get valid validators
	for ; iterator.Valid(); iterator.Next() {
		sequence := string(iterator.Key()[len(types.SlashingSequenceKey):])

		// call function and return if required
		if err := f(sequence); err != nil {
			return
		}
	}
	return
}
