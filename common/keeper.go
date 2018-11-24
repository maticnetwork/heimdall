package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	tmtypes "github.com/tendermint/tendermint/types"
	"strconv"
	"time"
)

//
// master keeper
//

type Keeper struct {
	MasterKey     sdk.StoreKey
	cdc           *codec.Codec
	CheckpointKey sdk.StoreKey
	StakingKey    sdk.StoreKey
	// codespace
	Codespace sdk.CodespaceType
}

// todo add staking keys here
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

// Add checkpoint to buffer or final headerBlocks
func (k Keeper) AddCheckpointToKey(ctx sdk.Context, start uint64, end uint64, root common.Hash, proposer common.Address, key []byte) sdk.Error {
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
	data := types.CreateBlock(start, end, root, proposer)
	out, err := json.Marshal(data)
	if err != nil {
		CheckpointLogger.Error("Error marshalling checkpoint to json", "error", err)
	}

	// store in key provided
	store.Set(key, []byte(out))

	return nil
}

// gets last checkpoint , headerIndex = TotalACKs * ChildBlockInterval
func (k Keeper) GetLastCheckpoint(ctx sdk.Context) types.CheckpointBlockHeader {
	store := ctx.KVStore(k.CheckpointKey)

	ACKs := k.GetACKCount(ctx)

	// fetch last checkpoint key (NumberOfACKs*ChildBlockInterval)
	lastCheckpointKey := (helper.GetConfig().ChildBlockInterval) * (ACKs)

	// fetch checkpoint and unmarshall
	var _checkpoint types.CheckpointBlockHeader
	err := json.Unmarshal(store.Get(GetHeaderKey(lastCheckpointKey)), &_checkpoint)
	if err != nil {
		CheckpointLogger.Error("Unable to fetch last checkpoint from store", "Key", lastCheckpointKey, "ACKCount", ACKs)
	}

	// return checkpoint
	return _checkpoint
}

// appends prefix to headerNumber
func GetHeaderKey(headerNumber int) []byte {
	headerNumberBytes := strconv.Itoa(headerNumber)
	return append(HeaderBlockKey, headerNumberBytes...)
}

// sets value in cache for checkpoint ACK
func (k Keeper) SetCheckpointAckCache(ctx sdk.Context, value []byte) {
	store := ctx.KVStore(k.CheckpointKey)
	store.Set(CheckpointACKCacheKey, value)
}

// sets value in cache for checkpoint
func (k Keeper) SetCheckpointCache(ctx sdk.Context, value []byte) {
	store := ctx.KVStore(k.CheckpointKey)
	store.Set(CheckpointCacheKey, value)
}

// check if value exists in cache or not
func (k Keeper) GetCheckpointCache(ctx sdk.Context, key []byte) bool {
	store := ctx.KVStore(k.CheckpointKey)
	value := store.Get(key)
	if bytes.Equal(value, EmptyBufferValue) {
		return false
	}
	return true
}

// Flush Checkpoint Buffer
func (k Keeper) FlushCheckpointBuffer(ctx sdk.Context) {
	store := ctx.KVStore(k.CheckpointKey)
	store.Set(BufferCheckpointKey, EmptyBufferValue)
}

// Get checkpoint in buffer
func (k Keeper) GetCheckpointFromBuffer(ctx sdk.Context) (types.CheckpointBlockHeader, error) {
	store := ctx.KVStore(k.CheckpointKey)

	// Get checkpoint and unmarshall
	var checkpoint types.CheckpointBlockHeader
	err := json.Unmarshal(store.Get(BufferCheckpointKey), &checkpoint)

	return checkpoint, err
}

// update ACK count by 1
func (k Keeper) UpdateACKCount(ctx sdk.Context) {
	store := ctx.KVStore(k.CheckpointKey)

	// get current ACK Count
	ACKCount := k.GetACKCount(ctx)

	// increment by 1
	ACKs := []byte(strconv.Itoa(ACKCount + 1))

	// update
	store.Set(ACKCountKey, ACKs)
}

// Get current ACK count
func (k Keeper) GetACKCount(ctx sdk.Context) int {
	store := ctx.KVStore(k.CheckpointKey)

	// get current ACK count
	ACKs, err := strconv.Atoi(string(store.Get(ACKCountKey)))
	if err != nil {
		CheckpointLogger.Error("Unable to convert key to int")
	}

	return ACKs
}

// Set ACK Count to 0
func (k Keeper) InitACKCount(ctx sdk.Context) {
	store := ctx.KVStore(k.CheckpointKey)

	// TODO maybe this needs to be set to 1
	// set to 0
	key := []byte(strconv.Itoa(int(0)))
	store.Set(ACKCountKey, key)
}

// ----------------- Staking Related Keepers

// getValidatorKey drafts the validator key for addresses
func getValidatorKey(address []byte) []byte {
	return append(ValidatorsKey, address...)
}

// Add validator indexed with address
func (k Keeper) AddValidator(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.StakingKey)

	// marshall validator
	bz := k.cdc.MustMarshalBinary(validator)

	// store validator with address prefixed with validator key as index
	store.Set(getValidatorKey(validator.Pubkey.Address().Bytes()), bz)
}

//  returns all validators who are in validator set and removes deactivated validators
func (k Keeper) GetCurrentValidators(ctx sdk.Context) (validators []types.Validator) {
	store := ctx.KVStore(k.StakingKey)

	// get ACK count
	ACKs := k.GetACKCount(ctx)

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

		// if you encounter a deactivated validator make power 0
		if validator.EndEpoch != 0 && validator.EndEpoch > int64(ACKs) && validator.Power != 0 {
			validator.Power = 0
			k.AddValidator(ctx, validator)
		}

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

func (k Keeper) GetAllValidators(ctx sdk.Context) (validators []abci.ValidatorUpdate) {
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
			Power:  validator.Power,
			PubKey: tmtypes.TM2PB.PubKey(validator.Pubkey),
		}

		// append to list of validatorUpdates
		validators = append(validators, updateVal)

		// increment iterator
		iterator.Next()
	}
	return
}

// performs deactivation of validatowrs wrt Tendermint to pass via EndBlock
func (k Keeper) RemoveDeactivatedValidators(ctx sdk.Context) {
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
		if validator.EndEpoch != 0 && validator.EndEpoch > int64(ACKs) && validator.Power != 0 {
			validator.Power = 0
			k.AddValidator(ctx, validator)
		}

		// increment iterator
		iterator.Next()
	}
	return
}

// GetValidatorInfo returns validator info for given the address
func (k Keeper) GetValidatorInfo(ctx sdk.Context, valAddr common.Address) (validator types.Validator, error error) {
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

func (k Keeper) AddDeactivationEpoch(ctx sdk.Context, valAddr common.Address, validator types.Validator) error {

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

// update validator with signer and pubkey
func (k Keeper) UpdateSigner(ctx sdk.Context, signer common.Address, pubkey crypto.PubKey, valAddr common.Address) error {
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
	validator.Pubkey = pubkey

	// add updated validator to store with same key
	k.AddValidator(ctx, validator)

	return nil
}

// Add validator set to store
func (k Keeper) UpdateValidatorSetInStore(ctx sdk.Context, newValidatorSet types.ValidatorSet) {
	store := ctx.KVStore(k.StakingKey)

	// marshall validator set
	bz := k.cdc.MustMarshalBinary(newValidatorSet)

	// set validator set with CurrentValidatorSetKey as key in store
	store.Set(CurrentValidatorSetKey, bz)
}

// Get current Validator Set from store
func (k Keeper) GetValidatorSet(ctx sdk.Context) (validatorSet types.ValidatorSet) {
	store := ctx.KVStore(k.StakingKey)

	// get current validator set from store
	bz := store.Get(CurrentValidatorSetKey)

	// unmarhsall
	k.cdc.MustUnmarshalBinary(bz, &validatorSet)

	return validatorSet
}

// increment accum for validator set by n times and replace validator set in store
func (k Keeper) IncreamentAccum(ctx sdk.Context, times int) {
	// get validator set
	validatorSet := k.GetValidatorSet(ctx)

	// increment accum
	validatorSet.IncrementAccum(times)

	// replace
	k.UpdateValidatorSetInStore(ctx, validatorSet)
}

// returns next proposer
func (k Keeper) GetNextProposer(ctx sdk.Context) string {
	// get validator set
	validatorSet := k.GetValidatorSet(ctx)

	// increment accum
	validatorSet.IncrementAccum(1)

	return validatorSet.Proposer.String()
}

// TODO expose via API
// returns current proposer
func (k Keeper) GetCurrentProposer(ctx sdk.Context) string {
	// get validator set
	validatorSet := k.GetValidatorSet(ctx)

	return validatorSet.Proposer.String()
}

// returns true if validator set has changed false otherwise
func (k Keeper) GetValidatorSetChangeFlag(ctx sdk.Context) bool {
	store := ctx.KVStore(k.StakingKey)

	// check if validator set change flag has value
	if bytes.Equal(store.Get(ValidatorSetChangeKey), DefaultValue) {
		return true
	}

	// validator set change flag is  empty
	return false
}

// inverts flag value for validator update
func (k Keeper) InvertValidatorSetChangeFlag(ctx sdk.Context) {
	store := ctx.KVStore(k.StakingKey)

	// Check if flag has value or not
	if bytes.Equal(store.Get(ValidatorSetChangeKey), DefaultValue) {
		store.Set(ValidatorSetChangeKey, EmptyBufferValue)
	} else {
		store.Set(ValidatorSetChangeKey, DefaultValue)
	}
}
