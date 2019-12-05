package staking

import (
	"encoding/hex"
	"errors"
	"math/big"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/ethereum/go-ethereum/common"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/helper"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
	"github.com/maticnetwork/heimdall/types"
	"github.com/tendermint/tendermint/libs/log"
)

var (
	DefaultValue = []byte{0x01} // Value to store in CacheCheckpoint and CacheCheckpointACK & ValidatorSetChange Flag

	ValidatorsKey          = []byte{0x21} // prefix for each key to a validator
	ValidatorMapKey        = []byte{0x22} // prefix for each key for validator map
	CurrentValidatorSetKey = []byte{0x23} // Key to store current validator set
	ValidatorRewardMapKey  = []byte{0x16} // prefix for each key for Validator Reward Map

)

// type AckRetriever struct {
// 	GetACKCount(ctx sdk.Context,hm app.HeimdallApp) uint64
// }
type AckRetriever interface {
	GetACKCount(ctx sdk.Context) uint64
}

// func (d AckRetriever) GetACKCount(ctx sdk.Context) uint64 {
// 	return app.checkpointKeeper.GetACKCount(ctx)
// }
// Keeper stores all related data
type Keeper struct {
	cdc *codec.Codec
	// The (unexposed) keys used to access the stores from the Context.
	storeKey sdk.StoreKey
	// codespacecodespace
	codespace sdk.CodespaceType
	// param space
	paramSpace params.Subspace
	// ack retriever
	ackRetriever AckRetriever
}

// NewKeeper create new keeper
func NewKeeper(
	cdc *codec.Codec,
	storeKey sdk.StoreKey,
	paramSpace params.Subspace,
	codespace sdk.CodespaceType,
	ackRetriever AckRetriever,
) Keeper {
	keeper := Keeper{
		cdc:          cdc,
		storeKey:     storeKey,
		paramSpace:   paramSpace.WithKeyTable(ParamKeyTable()),
		codespace:    codespace,
		ackRetriever: ackRetriever,
	}
	return keeper
}

// Codespace returns the codespace
func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", stakingTypes.ModuleName)
}

// GetValidatorKey drafts the validator key for addresses
func GetValidatorKey(address []byte) []byte {
	return append(ValidatorsKey, address...)
}

// GetValidatorMapKey returns validator map
func GetValidatorMapKey(address []byte) []byte {
	return append(ValidatorMapKey, address...)
}

// GetValidatorRewardMapKey returns validator reward map
func GetValidatorRewardMapKey(valID []byte) []byte {
	return append(ValidatorRewardMapKey, valID...)
}

// AddValidator adds validator indexed with address
func (k *Keeper) AddValidator(ctx sdk.Context, validator types.Validator) error {
	// TODO uncomment
	//if ok:=validator.ValidateBasic(); !ok{
	//	// return error
	//}

	store := ctx.KVStore(k.storeKey)

	bz, err := types.MarshallValidator(k.cdc, validator)
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
	ackCount := k.ackRetriever.GetACKCount(ctx)

	// get validator info
	validator, err := k.GetValidatorInfo(ctx, address)
	if err != nil {
		return false
	}

	// check if validator is current validator
	return validator.IsCurrentValidator(ackCount)
}

// GetValidatorInfo returns validator
func (k *Keeper) GetValidatorInfo(ctx sdk.Context, address []byte) (validator types.Validator, err error) {
	store := ctx.KVStore(k.storeKey)

	// check if validator exists
	key := GetValidatorKey(address)
	if !store.Has(key) {
		return validator, errors.New("Validator not found")
	}

	// unmarshall validator and return
	validator, err = types.UnmarshallValidator(k.cdc, store.Get(key))
	if err != nil {
		return validator, err
	}

	// return true if validator
	return validator, nil
}

// GetCurrentValidators returns all validators who are in validator set
func (k *Keeper) GetCurrentValidators(ctx sdk.Context) (validators []types.Validator) {
	// get ack count
	ackCount := k.ackRetriever.GetACKCount(ctx)

	// Get validators
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

// GetSpanEligibleValidators returns current validators who are not getting deactivated in between next span
func (k *Keeper) GetSpanEligibleValidators(ctx sdk.Context) (validators []types.Validator) {
	// get ack count
	ackCount := k.ackRetriever.GetACKCount(ctx)

	// Get validators and iterate through validator list
	k.IterateValidatorsAndApplyFn(ctx, func(validator types.Validator) error {
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
func (k *Keeper) IterateValidatorsAndApplyFn(ctx sdk.Context, f func(validator types.Validator) error) {
	store := ctx.KVStore(k.storeKey)

	// get validator iterator
	iterator := sdk.KVStorePrefixIterator(store, ValidatorsKey)
	defer iterator.Close()

	// loop through validators to get valid validators
	for ; iterator.Valid(); iterator.Next() {
		// unmarshall validator
		validator, _ := types.UnmarshallValidator(k.cdc, iterator.Value())
		// call function and return if required
		if err := f(validator); err != nil {
			return
		}
	}
}

// AddDeactivationEpoch adds deactivation epoch
func (k *Keeper) AddDeactivationEpoch(ctx sdk.Context, validator types.Validator, updatedVal types.Validator) error {
	// check if validator has unstaked
	if updatedVal.EndEpoch != 0 {
		validator.EndEpoch = updatedVal.EndEpoch
		// update validator in store
		return k.AddValidator(ctx, validator)
	}

	return errors.New("Deactivation period not set")
}

// UpdateSigner updates validator with signer and pubkey + validator => signer map
func (k *Keeper) UpdateSigner(ctx sdk.Context, newSigner types.HeimdallAddress, newPubkey types.PubKey, prevSigner types.HeimdallAddress) error {
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
func (k *Keeper) UpdateValidatorSetInStore(ctx sdk.Context, newValidatorSet types.ValidatorSet) error {
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
func (k *Keeper) GetValidatorSet(ctx sdk.Context) (validatorSet types.ValidatorSet) {
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
func (k *Keeper) GetNextProposer(ctx sdk.Context) *types.Validator {
	// get validator set
	validatorSet := k.GetValidatorSet(ctx)

	// Increment accum in copy
	copiedValidatorSet := validatorSet.CopyIncrementProposerPriority(1)

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

// SetValidatorIDToSignerAddr sets mapping for validator ID to signer address
func (k *Keeper) SetValidatorIDToSignerAddr(ctx sdk.Context, valID types.ValidatorID, signerAddr types.HeimdallAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(GetValidatorMapKey(valID.Bytes()), signerAddr.Bytes())
}

// GetSignerFromValidatorID get signer address from validator ID
func (k *Keeper) GetSignerFromValidatorID(ctx sdk.Context, valID types.ValidatorID) (common.Address, bool) {
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
func (k *Keeper) GetValidatorFromValID(ctx sdk.Context, valID types.ValidatorID) (validator types.Validator, ok bool) {
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
func (k *Keeper) GetLastUpdated(ctx sdk.Context, valID types.ValidatorID) (updatedAt uint64, found bool) {
	// get validator
	validator, ok := k.GetValidatorFromValID(ctx, valID)
	if !ok {
		return 0, false
	}
	return validator.LastUpdated, true
}

// SetValidatorIDToReward will update valid with new reward
func (k *Keeper) SetValidatorIDToReward(ctx sdk.Context, valID types.ValidatorID, reward *big.Int) {
	// Add reward to reward balance
	store := ctx.KVStore(k.storeKey)
	rewardBalance := k.GetRewardByValidatorID(ctx, valID)
	totalReward := big.NewInt(0).Add(reward, rewardBalance)
	store.Set(GetValidatorRewardMapKey(valID.Bytes()), totalReward.Bytes())
}

// GetRewardByValidatorID Returns Total Rewards of Validator
func (k *Keeper) GetRewardByValidatorID(ctx sdk.Context, valID types.ValidatorID) *big.Int {
	store := ctx.KVStore(k.storeKey)
	key := GetValidatorRewardMapKey(valID.Bytes())
	if store.Has(key) {
		// get current reward for validatorId
		rewardBalance := big.NewInt(0).SetBytes(store.Get(key))
		return rewardBalance
	}
	return big.NewInt(0)
}

// GetAllValidatorRewards returns validator reward map
func (k *Keeper) GetAllValidatorRewards(ctx sdk.Context) map[types.ValidatorID]*big.Int {
	store := ctx.KVStore(k.storeKey)
	valRewardMap := make(map[types.ValidatorID]*big.Int)
	iterator := sdk.KVStorePrefixIterator(store, ValidatorRewardMapKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		val := iterator.Value()
		// // unmarshalling val
		reward := big.NewInt(0).SetBytes(val)

		// unmarshalling key
		valID, err := strconv.ParseUint(string(iterator.Key()[len(ValidatorRewardMapKey):]), 10, 64)
		if err != nil {
			k.Logger(ctx).Debug("Error while parsing ValId",
				"valID", string(iterator.Key()[len(ValidatorRewardMapKey):]),
			)
		}
		valRewardMap[types.ValidatorID(valID)] = reward
	}
	return valRewardMap
}

// CalculateSignerRewards calculates new rewards for signers
func (k *Keeper) CalculateSignerRewards(ctx sdk.Context, voteBytes []byte, sigInput []byte, checkpointReward *big.Int) (map[types.ValidatorID]*big.Int, error) {
	signerRewards := make(map[types.ValidatorID]*big.Int)
	signerPower := make(map[types.ValidatorID]int64)

	const sigLength = 65
	totalSignerPower := int64(0)

	// Calculate total stake Power of all Signers.
	for i := 0; i < len(sigInput); i += sigLength {
		signature := sigInput[i : i+sigLength]
		pKey, err := authTypes.RecoverPubkey(voteBytes, []byte(signature))
		if err != nil {
			k.Logger(ctx).Error("Error Recovering PubKey", "Error", err)
			return nil, err
		}

		pubKey := types.NewPubKey(pKey)
		signerAddress := pubKey.Address().Bytes()
		valInfo, err := k.GetValidatorInfo(ctx, signerAddress)

		if err != nil {
			k.Logger(ctx).Error("No Validator Found for", "SignerAddress", signerAddress, "Error", err)
			return nil, err
		}
		totalSignerPower += valInfo.VotingPower
		signerPower[valInfo.ID] = valInfo.VotingPower
	}

	currentCheckpointReward := new(big.Float)
	currentCheckpointReward.SetInt(checkpointReward)
	totalSignerReward := k.ComputeTotalSignerReward(ctx, currentCheckpointReward)

	// Weighted Distribution of totalSignerReward for signers
	for valID, pow := range signerPower {
		bigpow := new(big.Float)
		bigpow.SetInt(big.NewInt(pow))
		bigval := big.NewFloat(0).Mul(totalSignerReward, bigpow)

		bigfloatval := new(big.Float)
		bigfloatval.SetFloat64(float64(totalSignerPower))

		signerReward := new(big.Float).Quo(bigval, bigfloatval)

		signerRewards[valID], _ = signerReward.Int(signerRewards[valID])
		k.Logger(ctx).Debug("Updated Reward for Validator", "ValidatorId", valID, "Reward", signerRewards[valID])
	}

	// Proposer Bonus Reward
	proposer := k.GetCurrentProposer(ctx)
	proposerReward := k.ComputeProposerReward(ctx, currentCheckpointReward)
	proposerBonus := big.NewInt(0)
	proposerReward.Int(proposerBonus)
	signerRewards[proposer.ID] = big.NewInt(0).Add(signerRewards[proposer.ID], proposerBonus)
	k.Logger(ctx).Debug("Updated Reward for Validator with proposer bonus", "ValidatorId", proposer.ID, "Reward", signerRewards[proposer.ID])

	return signerRewards, nil
}

// UpdateValidatorRewards Updates validators with Rewards
func (k *Keeper) UpdateValidatorRewards(ctx sdk.Context, valrewards map[types.ValidatorID]*big.Int) {
	for valID, reward := range valrewards {
		k.SetValidatorIDToReward(ctx, valID, reward)
	}
}

// ComputeTotalSignerReward returns total reward of all signer
func (k *Keeper) ComputeTotalSignerReward(ctx sdk.Context, currentCheckpointReward *big.Float) *big.Float {
	proposerBonusPercent := k.GetProposerBonusPercent(ctx)
	totalSignerRewardPercent := 100 - proposerBonusPercent
	signerRewardToFloat := new(big.Float)
	signerRewardToFloat.SetInt(big.NewInt(totalSignerRewardPercent))
	bigval := big.NewFloat(0).Mul(currentCheckpointReward, signerRewardToFloat)
	bigfloatval := new(big.Float)
	bigfloatval.SetFloat64(float64(100))
	return new(big.Float).Quo(bigval, bigfloatval)
}

// ComputeProposerReward returns the proposer reward
func (k *Keeper) ComputeProposerReward(ctx sdk.Context, currentCheckpointReward *big.Float) *big.Float {
	proposerBonus := k.GetProposerBonusPercent(ctx)
	proposerBonusToFloat := new(big.Float)
	proposerBonusToFloat.SetInt(big.NewInt(proposerBonus))
	bigval := big.NewFloat(0).Mul(currentCheckpointReward, proposerBonusToFloat)
	bigfloatval := new(big.Float)
	bigfloatval.SetFloat64(float64(100))
	return new(big.Float).Quo(bigval, bigfloatval)
}

// GetProposerBonusPercent returns the proposer to signer reward
func (k *Keeper) GetProposerBonusPercent(ctx sdk.Context) int64 {
	var proposerBonusPercent int64
	k.paramSpace.Get(ctx, ParamStoreKeyProposerBonusPercent, &proposerBonusPercent)
	return proposerBonusPercent
}

// SetProposerBonusPercent sets the Proposer to signer reward
func (k *Keeper) SetProposerBonusPercent(ctx sdk.Context, proposerBonusPercent int64) {
	k.paramSpace.Set(ctx, ParamStoreKeyProposerBonusPercent, proposerBonusPercent)
}
