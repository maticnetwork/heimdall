package common

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"

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
	EmptyBufferValue = []byte{0x00} // denotes EMPTY
	DefaultValue     = []byte{0x01} // Value to store in CacheCheckpoint and CacheCheckpointACK & ValidatorSetChange Flag

	ACKCountKey             = []byte{0x11} // key to store ACK count
	BufferCheckpointKey     = []byte{0x12} // Key to store checkpoint in buffer
	HeaderBlockKey          = []byte{0x13} // prefix key for when storing header after ACk
	CheckpointCacheKey      = []byte{0x14} // key to store Cache for checkpoint
	CheckpointACKCacheKey   = []byte{0x15} // key to store Cache for checkpointACK
	CheckpointNoACKCacheKey = []byte{0x16} // key to store last no-ack

	ValidatorsKey          = []byte{0x21} // prefix for each key to a validator
	ValidatorMapKey        = []byte{0x22} // prefix for each key for validator map
	CurrentValidatorSetKey = []byte{0x23} // Key to store current validator set
)

//--------------- Checkpoint Related Keepers

func (k *Keeper) _addCheckpoint(ctx sdk.Context, key []byte, headerBlock types.CheckpointBlockHeader) error {
	store := ctx.KVStore(k.CheckpointKey)

	// create Checkpoint block and marshall
	out, err := json.Marshal(headerBlock)
	if err != nil {
		CheckpointLogger.Error("Error marshalling checkpoint to json", "error", err)
		return err
	}

	// store in key provided
	store.Set(key, out)

	return nil
}

// AddCheckpoint adds checkpoint into final blocks
func (k *Keeper) AddCheckpoint(ctx sdk.Context, headerBlockNumber uint64, headerBlock types.CheckpointBlockHeader) error {
	key := GetHeaderKey(headerBlockNumber)
	err := k._addCheckpoint(ctx, key, headerBlock)
	if err != nil {
		return err
	}
	CheckpointLogger.Info("Adding good checkpoint to state", "checkpoint", headerBlock, "headerBlockNumber", headerBlockNumber)
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
	headerNumberBytes := []byte(strconv.FormatUint(headerNumber, 10))
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

// SetCheckpointBuffer flushes Checkpoint Buffer
func (k *Keeper) SetCheckpointBuffer(ctx sdk.Context, headerBlock types.CheckpointBlockHeader) error {
	err := k._addCheckpoint(ctx, BufferCheckpointKey, headerBlock)
	if err != nil {
		return err
	}
	CheckpointLogger.Debug("Adding good checkpoint to buffer to await ACK", "checkpoint", headerBlock)
	return nil
}

// GetCheckpointFromBuffer gets checkpoint in buffer
func (k *Keeper) GetCheckpointFromBuffer(ctx sdk.Context) (types.CheckpointBlockHeader, error) {
	store := ctx.KVStore(k.CheckpointKey)

	// checkpoint block header
	var checkpoint types.CheckpointBlockHeader

	if store.Has(BufferCheckpointKey) {
		// Get checkpoint and unmarshall
		err := json.Unmarshal(store.Get(BufferCheckpointKey), &checkpoint)
		return checkpoint, err
	}

	return checkpoint, errors.New("No checkpoint found in buffer")
}

// UpdateACKCountWithValue updates ACK with value
func (k *Keeper) UpdateACKCountWithValue(ctx sdk.Context, value uint64) {
	store := ctx.KVStore(k.CheckpointKey)

	// convert
	ackCount := []byte(strconv.FormatUint(value, 10))

	// update
	store.Set(ACKCountKey, ackCount)
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

// SetLastNoAck set last no-ack object
func (k *Keeper) SetLastNoAck(ctx sdk.Context, timestamp uint64) {
	store := ctx.KVStore(k.CheckpointKey)
	// convert timestamp to bytes
	value := []byte(strconv.FormatUint(timestamp, 10))
	// set no-ack
	store.Set(CheckpointNoACKCacheKey, value)
}

// GetLastNoAck returns last no ack
func (k *Keeper) GetLastNoAck(ctx sdk.Context) uint64 {
	store := ctx.KVStore(k.CheckpointKey)
	// check if ack count is there
	if store.Has(CheckpointNoACKCacheKey) {
		// get current ACK count
		result, err := strconv.ParseUint(string(store.Get(CheckpointNoACKCacheKey)), 10, 64)
		if err == nil {
			return uint64(result)
		}
	}

	return 0
}

// ----------------- Staking Related Keepers

// GetValidatorKey drafts the validator key for addresses
func GetValidatorKey(address []byte) []byte {
	return append(ValidatorsKey, address...)
}

// GetValidatorMapKey returns validator map
func GetValidatorMapKey(address []byte) []byte {
	return append(ValidatorMapKey, address...)
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
	store.Set(GetValidatorKey(validator.Signer.Bytes()), bz)
	StakingLogger.Debug("Validator stored", "key", hex.EncodeToString(GetValidatorKey(validator.Signer.Bytes())), "validator", validator.String())

	return nil
}

// GetValidatorInfo returns validator
func (k *Keeper) GetValidatorInfo(ctx sdk.Context, address []byte, validator *types.Validator) bool {
	store := ctx.KVStore(k.StakingKey)

	// store validator with address prefixed with validator key as index
	key := GetValidatorKey(address)
	if !store.Has(key) {
		return false
	}

	// unmarshall validator and return
	err := k.cdc.UnmarshalBinary(store.Get(key), validator)
	if err != nil {
		return false
	}

	// return true if validator
	return true
}

// GetCurrentValidators returns all validators who are in validator set
func (k *Keeper) GetCurrentValidators(ctx sdk.Context) (validators []types.Validator) {
	// get ACK count
	ackCount := k.GetACKCount(ctx)

	// iterate through validator list
	k.IterateValidatorsAndApplyFn(ctx, func(validator types.Validator) error {
		// check if validator is valid for current epoch
		if validator.IsCurrentValidator(ackCount) {
			// append if validator is current valdiator
			validators = append(validators, validator)
		}
		return nil
	})

	return
}

// GetAllValidators returns all validators
func (k *Keeper) GetAllValidators(ctx sdk.Context) (validators []*types.Validator) {
	// iterate through validators and create validator update array
	k.IterateValidatorsAndApplyFn(ctx, func(validator types.Validator) error {
		// append to list of validatorUpdates
		validators = append(validators, &validator)
		return nil
	})

	return
}

// IterateValidatorsAndApplyFn interate validators and apply the given function.
// Please do not modify validator store while iterating
func (k *Keeper) IterateValidatorsAndApplyFn(ctx sdk.Context, f func(validator types.Validator) error) {
	store := ctx.KVStore(k.StakingKey)

	// get validator iterator
	iterator := sdk.KVStorePrefixIterator(store, ValidatorsKey)
	defer iterator.Close()

	// loop through validators to get valid validators
	for ; iterator.Valid(); iterator.Next() {
		// unmarshall validator
		var validator types.Validator
		k.cdc.MustUnmarshalBinary(iterator.Value(), &validator)

		// call function and return if required
		if err := f(validator); err != nil {
			return
		}
	}
}

// AddDeactivationEpoch adds deactivation epoch
func (k *Keeper) AddDeactivationEpoch(ctx sdk.Context, validator types.Validator) error {
	// get validator from mainchain
	updatedVal, err := helper.GetValidatorInfo(validator.Address)
	if err != nil {
		StakingLogger.Error("Cannot fetch validator info while unstaking", "Error", err, "ValidatorAddress", validator.Address)
	}

	// check if validator has unstaked
	if updatedVal.EndEpoch != 0 {
		validator.EndEpoch = updatedVal.EndEpoch
		// update validator in store
		k.AddValidator(ctx, validator)
		return nil
	}

	return errors.New("Deactivation period not set")
}

// UpdateSigner updates validator with signer and pubkey
func (k *Keeper) UpdateSigner(ctx sdk.Context, newSigner common.Address, newPubkey types.PubKey, prevSigner common.Address) error {
	// get old validator from state and make power 0
	var validator types.Validator
	k.GetValidatorInfo(ctx, prevSigner.Bytes(), &validator)

	// copy power to reassign below
	validatorPower := validator.Power
	validator.Power = 0
	// update validator
	k.AddValidator(ctx, validator)

	//update signer in prev Validator
	validator.Signer = newSigner
	validator.PubKey = newPubkey
	validator.Power = validatorPower

	// add updated validator to store with new key
	k.AddValidator(ctx, validator)

	return nil
}

// UpdateValidatorSetInStore adds validator set to store
func (k *Keeper) UpdateValidatorSetInStore(ctx sdk.Context, newValidatorSet types.ValidatorSet) {
	// TODO check if we may have to delay this by 1 height to sync with tendermint validator updates
	store := ctx.KVStore(k.StakingKey)

	// marshall validator set
	bz := k.cdc.MustMarshalBinary(newValidatorSet)

	// set validator set with CurrentValidatorSetKey as key in store
	store.Set(CurrentValidatorSetKey, bz)
}

// GetValidatorSet returns current Validator Set from store
func (k *Keeper) GetValidatorSet(ctx sdk.Context) (validatorSet types.ValidatorSet) {
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
func (k *Keeper) GetNextProposer(ctx sdk.Context) *types.Validator {
	// get validator set
	validatorSet := k.GetValidatorSet(ctx)

	// Increment accum in copy
	copiedValidatorSet := validatorSet.CopyIncrementAccum(1)

	// get signer address for next signer
	return copiedValidatorSet.GetProposer()
}

// GetCurrentProposer returns current proposer
func (k *Keeper) GetCurrentProposer(ctx sdk.Context) *types.Validator {
	// get validator set
	validatorSet := k.GetValidatorSet(ctx)

	// return get proposer
	return validatorSet.GetProposer()
}

// SetValidatorAddrToSignerAddr set mapping for validator address to signer address
func (k *Keeper) SetValidatorAddrToSignerAddr(ctx sdk.Context, validatorAddr common.Address, signerAddr common.Address) {
	store := ctx.KVStore(k.StakingKey)
	store.Set(GetValidatorMapKey(validatorAddr.Bytes()), signerAddr.Bytes())
}

// GetSignerFromValidator get signer address from validator
func (k *Keeper) GetSignerFromValidator(ctx sdk.Context, validatorAddr common.Address) (common.Address, bool) {
	store := ctx.KVStore(k.StakingKey)
	key := GetValidatorMapKey(validatorAddr.Bytes())
	// check if validator address has been mapped
	if !store.Has(key) {
		return common.Address{}, false
	}

	// return address from bytes
	return common.BytesToAddress(store.Get(key)), true
}

// GetValidatorFromValAddr returns signer from validator address
func (k *Keeper) GetValidatorFromValAddr(ctx sdk.Context, validatorAddr common.Address, val *types.Validator) bool {
	store := ctx.KVStore(k.StakingKey)
	key := GetValidatorMapKey(validatorAddr.Bytes())
	// check if validator address has been mapped
	if !store.Has(key) {
		return false
	}

	// query for validator using ValidatorAddress => SignerAddress map
	return k.GetValidatorInfo(ctx, store.Get(key), val)
}

// GetValidatorToSignerMap returns validator to signer map
func (k *Keeper) GetValidatorToSignerMap(ctx sdk.Context) map[string]common.Address {
	store := ctx.KVStore(k.StakingKey)

	// get validator iterator
	iterator := sdk.KVStorePrefixIterator(store, ValidatorMapKey)
	defer iterator.Close()

	// initialize result map
	result := make(map[string]common.Address)

	// loop through validators to get valid validators
	prefixLength := len(ValidatorMapKey)
	for ; iterator.Valid(); iterator.Next() {
		key := hex.EncodeToString(iterator.Key()[prefixLength:])
		result[key] = common.BytesToAddress(iterator.Value())
	}
	return result
}
