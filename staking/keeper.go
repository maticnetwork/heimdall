package staking

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/maticnetwork/bor/common"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

var (
	DefaultValue = []byte{0x01} // Value to store in CacheCheckpoint and CacheCheckpointACK & ValidatorSetChange Flag

	ValidatorsKey             = []byte{0x21} // prefix for each key to a validator
	ValidatorMapKey           = []byte{0x22} // prefix for each key for validator map
	CurrentValidatorSetKey    = []byte{0x23} // Key to store current validator set
	PrevDividendAccountMapKey = []byte{0x41} // store for dividend accounts before checkpoint ack.
	DividendAccountMapKey     = []byte{0x42} // prefix for each key for Dividend Account Map
	StakingSequenceKey        = []byte{0x24} // prefix for each key for staking sequence map
)

// ModuleCommunicator manages different module interaction
type ModuleCommunicator interface {
	GetACKCount(ctx sdk.Context) uint64
	SetCoins(ctx sdk.Context, addr hmTypes.HeimdallAddress, amt hmTypes.Coins) sdk.Error
	GetCoins(ctx sdk.Context, addr hmTypes.HeimdallAddress) hmTypes.Coins
	SendCoins(ctx sdk.Context, from hmTypes.HeimdallAddress, to hmTypes.HeimdallAddress, amt hmTypes.Coins) sdk.Error
}

// Keeper stores all related data
type Keeper struct {
	cdc *codec.Codec
	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey
	// codespacecodespace
	codespace sdk.CodespaceType
	// param space
	paramSpace params.Subspace
	// module communicator
	moduleCommunicator ModuleCommunicator
}

// NewKeeper create new keeper
func NewKeeper(
	cdc *codec.Codec,
	storeKey sdk.StoreKey,
	paramSpace params.Subspace,
	codespace sdk.CodespaceType,
	moduleCommunicator ModuleCommunicator,
) Keeper {
	keeper := Keeper{
		cdc:                cdc,
		storeKey:           storeKey,
		paramSpace:         paramSpace.WithKeyTable(types.ParamKeyTable()),
		codespace:          codespace,
		moduleCommunicator: moduleCommunicator,
	}
	return keeper
}

// Codespace returns the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", types.ModuleName)
}

// GetValidatorKey drafts the validator key for addresses
func GetValidatorKey(address []byte) []byte {
	return append(ValidatorsKey, address...)
}

// GetValidatorMapKey returns validator map
func GetValidatorMapKey(address []byte) []byte {
	return append(ValidatorMapKey, address...)
}

// GetStakingSequenceKey returns staking sequence key
func GetStakingSequenceKey(sequence uint64) []byte {
	return append(StakingSequenceKey, []byte(strconv.FormatUint(sequence, 10))...)
}

// AddValidator adds validator indexed with address
func (k *Keeper) AddValidator(ctx sdk.Context, validator hmTypes.Validator) error {
	// TODO uncomment
	//if ok:=validator.ValidateBasic(); !ok{
	//	// return error
	//}

	store := ctx.KVStore(k.storeKey)

	bz, err := hmTypes.MarshallValidator(k.cdc, validator)
	if err != nil {
		return err
	}

	// store validator with address prefixed with validator key as index
	store.Set(GetValidatorKey(validator.Signer.Bytes()), bz)
	k.Logger(ctx).Debug("Validator stored", "key", hex.EncodeToString(GetValidatorKey(validator.Signer.Bytes())), "validator", validator.String())

	// add validator to validator ID => SignerAddress map
	k.SetValidatorIDToSignerAddr(ctx, validator.ID, validator.Signer)

	return nil
}

// IsCurrentValidatorByAddress check if validator is in current validator set by signer address
func (k *Keeper) IsCurrentValidatorByAddress(ctx sdk.Context, address []byte) bool {
	// get ack count
	ackCount := k.moduleCommunicator.GetACKCount(ctx)

	// get validator info
	validator, err := k.GetValidatorInfo(ctx, address)
	if err != nil {
		return false
	}

	// check if validator is current validator
	return validator.IsCurrentValidator(ackCount)
}

// GetValidatorInfo returns validator
func (k *Keeper) GetValidatorInfo(ctx sdk.Context, address []byte) (validator hmTypes.Validator, err error) {
	store := ctx.KVStore(k.storeKey)

	// check if validator exists
	key := GetValidatorKey(address)
	if !store.Has(key) {
		return validator, errors.New("Validator not found")
	}

	// unmarshall validator and return
	validator, err = hmTypes.UnmarshallValidator(k.cdc, store.Get(key))
	if err != nil {
		return validator, err
	}

	// return true if validator
	return validator, nil
}

// GetCurrentValidators returns all validators who are in validator set
func (k *Keeper) GetCurrentValidators(ctx sdk.Context) (validators []hmTypes.Validator) {
	// get ack count
	ackCount := k.moduleCommunicator.GetACKCount(ctx)

	// Get validators
	// iterate through validator list
	k.IterateValidatorsAndApplyFn(ctx, func(validator hmTypes.Validator) error {
		// check if validator is valid for current epoch
		if validator.IsCurrentValidator(ackCount) {
			// append if validator is current valdiator
			validators = append(validators, validator)
		}
		return nil
	})

	return
}

// GetSpanEligibleValidators returns current validators who are not getting deactivated in between next span
func (k *Keeper) GetSpanEligibleValidators(ctx sdk.Context) (validators []hmTypes.Validator) {
	// get ack count
	ackCount := k.moduleCommunicator.GetACKCount(ctx)

	// Get validators and iterate through validator list
	k.IterateValidatorsAndApplyFn(ctx, func(validator hmTypes.Validator) error {
		// check if validator is valid for current epoch and endEpoch is not set.
		if validator.EndEpoch == 0 && validator.IsCurrentValidator(ackCount) {
			// append if validator is current valdiator
			validators = append(validators, validator)
		}
		return nil
	})

	return
}

// GetAllValidators returns all validators
func (k *Keeper) GetAllValidators(ctx sdk.Context) (validators []*hmTypes.Validator) {
	// iterate through validators and create validator update array
	k.IterateValidatorsAndApplyFn(ctx, func(validator hmTypes.Validator) error {
		// append to list of validatorUpdates
		validators = append(validators, &validator)
		return nil
	})

	return
}

// IterateValidatorsAndApplyFn interate validators and apply the given function.
func (k *Keeper) IterateValidatorsAndApplyFn(ctx sdk.Context, f func(validator hmTypes.Validator) error) {
	store := ctx.KVStore(k.storeKey)

	// get validator iterator
	iterator := sdk.KVStorePrefixIterator(store, ValidatorsKey)
	defer iterator.Close()

	// loop through validators to get valid validators
	for ; iterator.Valid(); iterator.Next() {
		// unmarshall validator
		validator, _ := hmTypes.UnmarshallValidator(k.cdc, iterator.Value())
		// call function and return if required
		if err := f(validator); err != nil {
			return
		}
	}
}

// AddDeactivationEpoch adds deactivation epoch
func (k *Keeper) AddDeactivationEpoch(ctx sdk.Context, validator hmTypes.Validator, updatedVal hmTypes.Validator) error {
	// check if validator has unstaked
	if updatedVal.EndEpoch != 0 {
		validator.EndEpoch = updatedVal.EndEpoch
		// update validator in store
		return k.AddValidator(ctx, validator)
	}

	return errors.New("Deactivation period not set")
}

// UpdateSigner updates validator with signer and pubkey + validator => signer map
func (k *Keeper) UpdateSigner(ctx sdk.Context, newSigner hmTypes.HeimdallAddress, newPubkey hmTypes.PubKey, prevSigner hmTypes.HeimdallAddress) error {
	// get old validator from state and make power 0
	validator, err := k.GetValidatorInfo(ctx, prevSigner.Bytes())
	if err != nil {
		k.Logger(ctx).Error("Unable to fetch valiator from store")
		return err
	}

	// copy power to reassign below
	validatorPower := validator.VotingPower
	validator.VotingPower = 0

	// update validator
	k.AddValidator(ctx, validator)

	//update signer in prev Validator
	validator.Signer = newSigner
	validator.PubKey = newPubkey
	validator.VotingPower = validatorPower

	// add updated validator to store with new key
	k.AddValidator(ctx, validator)
	return nil
}

// UpdateValidatorSetInStore adds validator set to store
func (k *Keeper) UpdateValidatorSetInStore(ctx sdk.Context, newValidatorSet hmTypes.ValidatorSet) error {
	// TODO check if we may have to delay this by 1 height to sync with tendermint validator updates
	store := ctx.KVStore(k.storeKey)

	// marshall validator set
	bz, err := k.cdc.MarshalBinaryBare(newValidatorSet)
	if err != nil {
		return err
	}

	// set validator set with CurrentValidatorSetKey as key in store
	store.Set(CurrentValidatorSetKey, bz)
	return nil
}

// GetValidatorSet returns current Validator Set from store
func (k *Keeper) GetValidatorSet(ctx sdk.Context) (validatorSet hmTypes.ValidatorSet) {
	store := ctx.KVStore(k.storeKey)
	// get current validator set from store
	bz := store.Get(CurrentValidatorSetKey)
	// unmarhsall
	k.cdc.UnmarshalBinaryBare(bz, &validatorSet)

	// return validator set
	return validatorSet
}

// IncrementAccum increments accum for validator set by n times and replace validator set in store
func (k *Keeper) IncrementAccum(ctx sdk.Context, times int) {
	// get validator set
	validatorSet := k.GetValidatorSet(ctx)

	// increment accum
	validatorSet.IncrementProposerPriority(times)

	// replace
	k.UpdateValidatorSetInStore(ctx, validatorSet)
}

// GetNextProposer returns next proposer
func (k *Keeper) GetNextProposer(ctx sdk.Context) *hmTypes.Validator {
	// get validator set
	validatorSet := k.GetValidatorSet(ctx)

	// Increment accum in copy
	copiedValidatorSet := validatorSet.CopyIncrementProposerPriority(1)

	// get signer address for next signer
	return copiedValidatorSet.GetProposer()
}

// GetCurrentProposer returns current proposer
func (k *Keeper) GetCurrentProposer(ctx sdk.Context) *hmTypes.Validator {
	// get validator set
	validatorSet := k.GetValidatorSet(ctx)

	// return get proposer
	return validatorSet.GetProposer()
}

// SetValidatorIDToSignerAddr sets mapping for validator ID to signer address
func (k *Keeper) SetValidatorIDToSignerAddr(ctx sdk.Context, valID hmTypes.ValidatorID, signerAddr hmTypes.HeimdallAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetValidatorMapKey(valID.Bytes()), signerAddr.Bytes())
}

// GetSignerFromValidatorID get signer address from validator ID
func (k *Keeper) GetSignerFromValidatorID(ctx sdk.Context, valID hmTypes.ValidatorID) (common.Address, bool) {
	store := ctx.KVStore(k.storeKey)
	key := GetValidatorMapKey(valID.Bytes())
	// check if validator address has been mapped
	if !store.Has(key) {
		return helper.ZeroAddress, false
	}
	// return address from bytes
	return common.BytesToAddress(store.Get(key)), true
}

// GetValidatorFromValID returns signer from validator ID
func (k *Keeper) GetValidatorFromValID(ctx sdk.Context, valID hmTypes.ValidatorID) (validator hmTypes.Validator, ok bool) {
	signerAddr, ok := k.GetSignerFromValidatorID(ctx, valID)
	if !ok {
		return validator, ok
	}
	// query for validator signer address
	validator, err := k.GetValidatorInfo(ctx, signerAddr.Bytes())
	if err != nil {
		return validator, false
	}
	return validator, true
}

// GetLastUpdated get last updated at for validator
func (k *Keeper) GetLastUpdated(ctx sdk.Context, valID hmTypes.ValidatorID) (updatedAt uint64, found bool) {
	// get validator
	validator, ok := k.GetValidatorFromValID(ctx, valID)
	if !ok {
		return 0, false
	}
	return validator.LastUpdated, true
}

// GetDividendAccountMapKey returns dividend account map
func GetDividendAccountMapKey(id []byte) []byte {
	return append(DividendAccountMapKey, id...)
}

// GetPrevDividendAccountMapKey returns prev dividend account map
func GetPrevDividendAccountMapKey(id []byte) []byte {
	return append(PrevDividendAccountMapKey, id...)
}

// AddDividendAccount adds DividendAccount index with DividendID
func (k *Keeper) AddDividendAccount(ctx sdk.Context, dividendAccount hmTypes.DividendAccount) error {
	store := ctx.KVStore(k.storeKey)
	// marshall dividend account
	bz, err := hmTypes.MarshallDividendAccount(k.cdc, dividendAccount)
	if err != nil {
		return err
	}

	store.Set(GetDividendAccountMapKey(dividendAccount.ID.Bytes()), bz)
	k.Logger(ctx).Debug("DividendAccount Stored", "key", hex.EncodeToString(GetDividendAccountMapKey(dividendAccount.ID.Bytes())), "dividendAccount", dividendAccount.String())
	return nil
}

// GetDividendAccountByID will return DividendAccount of valID
func (k *Keeper) GetDividendAccountByID(ctx sdk.Context, dividendID hmTypes.DividendAccountID) (dividendAccount hmTypes.DividendAccount, err error) {

	// check if dividend account exists
	if !k.CheckIfDividendAccountExists(ctx, dividendID) {
		return dividendAccount, errors.New("Dividend Account not found")
	}

	// Get DividendAccount key
	store := ctx.KVStore(k.storeKey)
	key := GetDividendAccountMapKey(dividendID.Bytes())

	// unmarshall dividend account and return
	dividendAccount, err = hmTypes.UnMarshallDividendAccount(k.cdc, store.Get(key))
	if err != nil {
		return dividendAccount, err
	}

	return dividendAccount, nil
}

// CheckIfDividendAccountExists will return true if dividend account exists
func (k *Keeper) CheckIfDividendAccountExists(ctx sdk.Context, dividendID hmTypes.DividendAccountID) (ok bool) {
	store := ctx.KVStore(k.storeKey)
	key := GetDividendAccountMapKey(dividendID.Bytes())
	if !store.Has(key) {
		return false
	}
	return true
}

// GetAllDividendAccounts returns all DividendAccountss
func (k *Keeper) GetAllDividendAccounts(ctx sdk.Context) (dividendAccounts []hmTypes.DividendAccount) {
	// iterate through dividendAccounts and create dividendAccounts update array
	k.IterateDividendAccountsByPrefixAndApplyFn(ctx, DividendAccountMapKey, func(dividendAccount hmTypes.DividendAccount) error {
		// append to list of dividendUpdates
		dividendAccounts = append(dividendAccounts, dividendAccount)
		return nil
	})

	return
}

// AddFeeToDividendAccount adds fee to dividend account for withdrawal
func (k *Keeper) AddFeeToDividendAccount(ctx sdk.Context, valID hmTypes.ValidatorID, fee *big.Int) sdk.Error {
	// Get or create dividend account
	var dividendAccount hmTypes.DividendAccount

	if k.CheckIfDividendAccountExists(ctx, hmTypes.DividendAccountID(valID)) {
		dividendAccount, _ = k.GetDividendAccountByID(ctx, hmTypes.DividendAccountID(valID))
	} else {
		dividendAccount = hmTypes.DividendAccount{
			ID:            hmTypes.DividendAccountID(valID),
			FeeAmount:     big.NewInt(0).String(),
			SlashedAmount: big.NewInt(0).String(),
		}
	}

	// update fee
	oldFee, _ := big.NewInt(0).SetString(dividendAccount.FeeAmount, 10)
	totalFee := big.NewInt(0).Add(oldFee, fee).String()
	dividendAccount.FeeAmount = totalFee

	k.Logger(ctx).Info("Dividend Account fee of validator ", "ID", dividendAccount.ID, "Fee", dividendAccount.FeeAmount)
	k.AddDividendAccount(ctx, dividendAccount)
	return nil
}

// IterateDividendAccountsByPrefixAndApplyFn iterate dividendAccounts and apply the given function.
func (k *Keeper) IterateDividendAccountsByPrefixAndApplyFn(ctx sdk.Context, prefix []byte, f func(dividendAccount hmTypes.DividendAccount) error) {
	store := ctx.KVStore(k.storeKey)

	// get validator iterator
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	// loop through dividendAccounts
	for ; iterator.Valid(); iterator.Next() {
		// unmarshall dividendAccount
		dividendAccount, _ := hmTypes.UnMarshallDividendAccount(k.cdc, iterator.Value())
		// call function and return if required
		if err := f(dividendAccount); err != nil {
			return
		}
	}
}

//
// Staking sequence
//

// SetStakingSequence sets staking sequence
func (k *Keeper) SetStakingSequence(ctx sdk.Context, sequence uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetStakingSequenceKey(sequence), DefaultValue)
}

// HasStakingSequence checks if staking sequence already exists
func (k *Keeper) HasStakingSequence(ctx sdk.Context, sequence uint64) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(GetStakingSequenceKey(sequence))
}
