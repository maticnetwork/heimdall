package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
)

//
// master keeper
//

// Keeper stores all related data
type Keeper struct {
	MasterKey     sdk.StoreKey
	cdc           *codec.Codec
	CheckpointKey sdk.StoreKey
	StakingKey    sdk.StoreKey
	// codespace
	Codespace sdk.CodespaceType
}

// NewKeeper create new keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, stakingKey sdk.StoreKey, checkpointKey sdk.StoreKey, codespace sdk.CodespaceType) Keeper {
	keeper := Keeper{
		MasterKey:     key,
		cdc:           cdc,
		Codespace:     codespace,
		CheckpointKey: checkpointKey,
		StakingKey:    stakingKey,
	}
	return keeper
}

// -------------- KEYS/CONSTANTS

var (
	ACKCountKey         = []byte{0x01} // key to store ACK count
	BufferCheckpointKey = []byte{0x02} // Key to store checkpoint in buffer
	HeaderBlockKey      = []byte{0x03} // prefix key for when storing header after ACk

	EmptyBufferValue = []byte{0x04} // denotes EMPTY

	CheckpointCacheKey    = []byte{0x05} // key to store Cache for checkpoint
	CheckpointACKCacheKey = []byte{0x06} // key to store Cache for checkpointACK

	DefaultValue = []byte{0x07} // Value to store in CacheCheckpoint and CacheCheckpointACK & ValidatorSetChange Flag

	ValidatorsKey          = []byte{0x08} // prefix for each key to a validator
	CurrentValidatorSetKey = []byte{0x09} // Key to store current validator set

	ValidatorSetChangeKey = []byte{0x010} // key to store flag for validator set update
)

//--------------- Checkpoint Related Keepers

// AddCheckpointToBuffer adds checkpoint to buffer or final headerBlocks
func (k *Keeper) AddCheckpointToBuffer(ctx sdk.Context, key []byte, headerBlock types.CheckpointBlockHeader) sdk.Error {
	store := ctx.KVStore(k.CheckpointKey)

	checkpointBuffer, _ := k.GetCheckpointFromBuffer(ctx)

	// Reject new checkpoint if checkpoint exists in buffer and 5 minutes have not passed
	if bytes.Equal(key, BufferCheckpointKey) && !bytes.Equal(store.Get(BufferCheckpointKey), EmptyBufferValue) && time.Now().UTC().Before(checkpointBuffer.TimeStamp.Add(helper.CheckpointBufferTime)) {
		return ErrNoACK(k.Codespace)
	}

	// Flush Checkpoint If 5 minutes have passed since it was added to buffer and NoAck received
	if bytes.Equal(key, BufferCheckpointKey) && !bytes.Equal(store.Get(BufferCheckpointKey), EmptyBufferValue) && time.Now().UTC().After(checkpointBuffer.TimeStamp.Add(helper.CheckpointBufferTime)) {
		k.FlushCheckpointBuffer(ctx)
	}

	// create Checkpoint block and marshall
	// data := types.CreateBlock(start, end, root, proposer)
	out, err := json.Marshal(headerBlock)
	if err != nil {
		CheckpointLogger.Error("Error marshalling checkpoint to json", "error", err)
	}

	if bytes.Equal(key, BufferCheckpointKey) {
		CheckpointLogger.Info("Adding good checkpoint to buffer to await ACK", "checkpoint", headerBlock)
	} else {
		CheckpointLogger.Info("Adding good checkpoint to state", "checkpoint", headerBlock)
	}

	// store in key provided
	store.Set(key, []byte(out))

	return nil
}

// GetLastCheckpoint gets last checkpoint, headerIndex = TotalACKs * ChildBlockInterval
func (k *Keeper) GetLastCheckpoint(ctx sdk.Context) (types.CheckpointBlockHeader, error) {
	store := ctx.KVStore(k.CheckpointKey)

	acksCount := k.GetACKCount(ctx)

	// fetch last checkpoint key (NumberOfACKs * ChildBlockInterval)
	lastCheckpointKey := helper.GetConfig().ChildBlockInterval * acksCount

	// fetch checkpoint and unmarshall
	var _checkpoint types.CheckpointBlockHeader

	// no checkpoint received
	if acksCount >= 0 {
		// header key
		headerKey := GetHeaderKey(lastCheckpointKey)
		if store.Has(headerKey) {
			err := json.Unmarshal(store.Get(headerKey), &_checkpoint)
			if err != nil {
				CheckpointLogger.Error("Unable to fetch last checkpoint from store", "key", lastCheckpointKey, "acksCount", acksCount)
				return _checkpoint, err
			} else {
				return _checkpoint, nil
			}
		}
	}

	return _checkpoint, errors.New("No checkpoint found")
}

// GetHeaderKey appends prefix to headerNumber
func GetHeaderKey(headerNumber uint64) []byte {
	headerNumberBytes := strconv.FormatUint(headerNumber, 10)
	return append(HeaderBlockKey, headerNumberBytes...)
}

// SetCheckpointAckCache sets value in cache for checkpoint ACK
func (k *Keeper) SetCheckpointAckCache(ctx sdk.Context, value []byte) {
	store := ctx.KVStore(k.CheckpointKey)
	store.Set(CheckpointACKCacheKey, value)
}

// SetCheckpointCache sets value in cache for checkpoint
func (k *Keeper) SetCheckpointCache(ctx sdk.Context, value []byte) {
	store := ctx.KVStore(k.CheckpointKey)
	store.Set(CheckpointCacheKey, value)
}

// GetCheckpointCache check if value exists in cache or not
func (k *Keeper) GetCheckpointCache(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.CheckpointKey)
	value := store.Get(key)
	if bytes.Equal(value, EmptyBufferValue) {
		return false
	}
	return true
}

// FlushCheckpointBuffer flushes Checkpoint Buffer
func (k *Keeper) FlushCheckpointBuffer(ctx sdk.Context) {
	store := ctx.KVStore(k.CheckpointKey)
	store.Set(BufferCheckpointKey, EmptyBufferValue)
}

// GetCheckpointFromBuffer gets checkpoint in buffer
func (k *Keeper) GetCheckpointFromBuffer(ctx sdk.Context) (types.CheckpointBlockHeader, error) {
	store := ctx.KVStore(k.CheckpointKey)

	// Get checkpoint and unmarshall
	var checkpoint types.CheckpointBlockHeader
	err := json.Unmarshal(store.Get(BufferCheckpointKey), &checkpoint)

	return checkpoint, err
}

// UpdateACKCount updates ACK count by 1
func (k *Keeper) UpdateACKCount(ctx sdk.Context) {
	store := ctx.KVStore(k.CheckpointKey)

	// get current ACK Count
	ACKCount := k.GetACKCount(ctx)

	// increment by 1
	ACKs := []byte(strconv.FormatUint(ACKCount+1, 10))

	// update
	store.Set(ACKCountKey, ACKs)
}

// GetACKCount returns current ACK count
func (k *Keeper) GetACKCount(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.CheckpointKey)
	// check if ack count is there
	if store.Has(ACKCountKey) {
		// get current ACK count
		ackCount, err := strconv.Atoi(string(store.Get(ACKCountKey)))
		if err != nil {
			CheckpointLogger.Error("Unable to convert key to int")
		} else {
			return uint64(ackCount)
		}
	}

	return 0
}

// InitACKCount sets ACK Count to 0
func (k *Keeper) InitACKCount(ctx sdk.Context) {
	store := ctx.KVStore(k.CheckpointKey)

	// TODO maybe this needs to be set to 1
	// set to 0
	key := []byte(strconv.Itoa(0))
	store.Set(ACKCountKey, key)
}

// ----------------- Staking Related Keepers

// getValidatorKey drafts the validator key for addresses
func getValidatorKey(address []byte) []byte {
	return append(ValidatorsKey, address...)
}

// AddValidator adds validator indexed with address
func (k *Keeper) AddValidator(ctx sdk.Context, validator types.Validator) error {
	store := ctx.KVStore(k.StakingKey)

	// marshall validator
	bz, err := k.cdc.MarshalBinary(validator)
	if err != nil {
		return err
	}

	// store validator with address prefixed with validator key as index
	store.Set(getValidatorKey(validator.Address.Bytes()), bz)

	return nil
}

// GetValidator returns validator
func (k *Keeper) GetValidator(ctx sdk.Context, address []byte, validator *types.Validator) error {
	store := ctx.KVStore(k.StakingKey)

	// store validator with address prefixed with validator key as index
	key := getValidatorKey(address)
	if !store.Has(key) {
		return errors.New("Validator not found")
	}

	// unmarshall validator and return
	return k.cdc.UnmarshalBinary(store.Get(key), validator)
}

// GetCurrentValidators returns all validators who are in validator set and removes deactivated validators
func (k *Keeper) GetCurrentValidators(ctx sdk.Context) (validators []types.Validator) {
	store := ctx.KVStore(k.StakingKey)

	// get ACK count
	ACKs := k.GetACKCount(ctx)

	// remove matured validators
	k.RemoveDeactivatedValidators(ctx)

	// create iterator to iterate with Validator Key prefix
	iterator := sdk.KVStorePrefixIterator(store, ValidatorsKey)
	defer iterator.Close()

	// loop through validators to get valid validators
	for i := 0; ; i++ {
		if !iterator.Valid() {
			break
		}

		// unmarshall validator
		var validator types.Validator
		k.cdc.MustUnmarshalBinary(iterator.Value(), &validator)

		// check if validator is valid for current epoch
		if validator.IsCurrentValidator(ACKs) {
			// append if validator is current valdiator
			validators = append(validators, validator)
		}

		// increment iterator
		iterator.Next()
	}
	return
}

// GetAllValidators returns all validators
func (k *Keeper) GetAllValidators(ctx sdk.Context) (validators []abci.ValidatorUpdate) {
	store := ctx.KVStore(k.StakingKey)

	// create iterator to iterate with Validator Key prefix
	iterator := sdk.KVStorePrefixIterator(store, ValidatorsKey)
	defer iterator.Close()

	// loop through validators to get valid validators
	for i := 0; ; i++ {
		if !iterator.Valid() {
			break
		}

		// unmarshall validator
		var validator types.Validator
		k.cdc.MustUnmarshalBinary(iterator.Value(), &validator)

		// convert to Validator Update
		updateVal := abci.ValidatorUpdate{
			Power:  int64(validator.Power),
			PubKey: tmTypes.TM2PB.PubKey(validator.PubKey),
		}

		// append to list of validatorUpdates
		validators = append(validators, updateVal)

		// increment iterator
		iterator.Next()
	}
	return
}

// RemoveDeactivatedValidators performs deactivation of validatowrs wrt Tendermint to pass via EndBlock
func (k *Keeper) RemoveDeactivatedValidators(ctx sdk.Context) {
	store := ctx.KVStore(k.StakingKey)

	// get ACK count
	ACKs := k.GetACKCount(ctx)

	iterator := sdk.KVStorePrefixIterator(store, ValidatorsKey)
	defer iterator.Close()

	// loop through validators to get valid validators
	for i := 0; ; i++ {
		if !iterator.Valid() {
			break
		}

		// unmarshall validator
		var validator types.Validator
		k.cdc.MustUnmarshalBinary(iterator.Value(), &validator)

		// if you encounter a deactivated validator make power 0
		if validator.EndEpoch != 0 && validator.EndEpoch > ACKs && validator.Power != 0 {
			validator.Power = 0

			// update validator in validator list
			k.AddValidator(ctx, validator)

			// indicate change in validator set
			k.SetValidatorSetChangedFlag(ctx, true)
		}

		// increment iterator
		iterator.Next()
	}
	return
}

// IterateValidatorsAndApplyFn interate validators and apply the given function
func (k *Keeper) IterateValidatorsAndApplyFn(ctx sdk.Context, f func()) {
	// TODO add remove and getter for validator list here
	f()
}

// GetValidatorInfo returns validator info for given the address
func (k *Keeper) GetValidatorInfo(ctx sdk.Context, valAddr common.Address) (validator types.Validator, error error) {
	store := ctx.KVStore(k.StakingKey)

	// get validator and unmarshall
	validatorBytes := store.Get(getValidatorKey(valAddr.Bytes()))
	if validatorBytes == nil {
		error = fmt.Errorf("Validator Not Found")
		return
	}

	// unmarshall validator
	k.cdc.MustUnmarshalBinary(validatorBytes, &validator)

	return validator, nil
}

// AddDeactivationEpoch adds deactivation epoch
func (k *Keeper) AddDeactivationEpoch(ctx sdk.Context, valAddr common.Address, validator types.Validator) error {
	// set deactivation period
	updatedVal, err := helper.GetValidatorInfo(valAddr)
	if err != nil {
		StakingLogger.Error("Cannot fetch validator info while unstaking", "Error", err, "ValidatorAddress", valAddr)
	}

	// check if validator has unstaked
	if updatedVal.EndEpoch != 0 {
		validator.EndEpoch = updatedVal.EndEpoch
		k.AddValidator(ctx, validator)
		return nil
	} else {
		StakingLogger.Debug("Deactivation period not set")
		return ErrValidatorAlreadySynced(k.Codespace)
	}

}

// UpdateSigner updates validator with signer and pubkey
func (k *Keeper) UpdateSigner(ctx sdk.Context, signer common.Address, pubkey crypto.PubKey, valAddr common.Address) error {
	store := ctx.KVStore(k.StakingKey)

	var validator types.Validator

	// get validator and unmarshall
	validatorBytes := store.Get(getValidatorKey(valAddr.Bytes()))
	if validatorBytes == nil {
		err := fmt.Errorf("Validator Not Found")
		return err
	}

	k.cdc.MustUnmarshalBinary(validatorBytes, &validator)

	//update signer
	validator.Signer = signer
	validator.PubKey = pubkey

	// add updated validator to store with same key
	k.AddValidator(ctx, validator)

	return nil
}

// UpdateValidatorSetInStore adds validator set to store
func (k *Keeper) UpdateValidatorSetInStore(ctx sdk.Context, newValidatorSet tmTypes.ValidatorSet) {
	// TODO check if we may have to delay this by 1 height to sync with tendermint validator updates
	store := ctx.KVStore(k.StakingKey)

	// marshall validator set
	bz := k.cdc.MustMarshalBinary(newValidatorSet)

	// set validator set with CurrentValidatorSetKey as key in store
	store.Set(CurrentValidatorSetKey, bz)
}

// GetValidatorSet returns current Validator Set from store
func (k *Keeper) GetValidatorSet(ctx sdk.Context) (validatorSet tmTypes.ValidatorSet) {
	store := ctx.KVStore(k.StakingKey)
	// get current validator set from store
	bz := store.Get(CurrentValidatorSetKey)
	// unmarhsall
	_ = k.cdc.UnmarshalBinary(bz, &validatorSet)
	// return validator set
	return validatorSet
}

// IncreamentAccum increments accum for validator set by n times and replace validator set in store
func (k *Keeper) IncreamentAccum(ctx sdk.Context, times int) {
	// get validator set
	validatorSet := k.GetValidatorSet(ctx)

	// increment accum
	validatorSet.IncrementAccum(times)

	// replace
	k.UpdateValidatorSetInStore(ctx, validatorSet)
}

// GetNextProposer returns next proposer
func (k *Keeper) GetNextProposer(ctx sdk.Context) string {
	// get validator set
	validatorSet := k.GetValidatorSet(ctx)

	// Increment accum in copy
	copiedValidatorSet := validatorSet.CopyIncrementAccum(1)

	return copiedValidatorSet.Proposer.String()
}

// GetCurrentProposerAddress returns current proposer
func (k *Keeper) GetCurrentProposerAddress(ctx sdk.Context) []byte {
	// TODO expose via API
	// get validator set
	validatorSet := k.GetValidatorSet(ctx)

	return validatorSet.Proposer.Address.Bytes()
}

// ValidatorSetChanged returns true if validator set has changed false otherwise
func (k *Keeper) ValidatorSetChanged(ctx sdk.Context) bool {
	store := ctx.KVStore(k.StakingKey)

	// check if validator set change flag has value
	if bytes.Equal(store.Get(ValidatorSetChangeKey), DefaultValue) {
		return true
	}

	// validator set change flag is  empty
	return false
}

//// inverts flag value for validator update
//func (k *Keeper) InvertValidatorSetChangeFlag(ctx sdk.Context) {
//	store := ctx.KVStore(k.StakingKey)
//
//	// Check if flag has value or not
//	if bytes.Equal(store.Get(ValidatorSetChangeKey), DefaultValue) {
//		store.Set(ValidatorSetChangeKey, EmptyBufferValue)
//	} else {
//		store.Set(ValidatorSetChangeKey, DefaultValue)
//	}
//}

// SetValidatorSetChangedFlag sets validator update flag depending on value
func (k *Keeper) SetValidatorSetChangedFlag(ctx sdk.Context, value bool) {
	store := ctx.KVStore(k.StakingKey)

	// check if validator set change flag has value
	if bytes.Equal(store.Get(ValidatorSetChangeKey), DefaultValue) && !value {
		store.Set(ValidatorSetChangeKey, EmptyBufferValue)
		return
	}

	store.Set(ValidatorSetChangeKey, DefaultValue)
}
