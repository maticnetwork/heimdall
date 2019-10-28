package staking

import (
	"encoding/hex"
	"errors"
	"math/big"

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
	ValidatorAccountMapKey = []byte{0x16} // prefix for each key for Validator Account Map
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

// GetValidatorAccountMapKey returns validator account map
func GetValidatorAccountMapKey(valID []byte) []byte {
	return append(ValidatorAccountMapKey, valID...)
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

// RewardValidator will update validator account with new reward
func (k *Keeper) RewardValidator(ctx sdk.Context, valID types.ValidatorID, reward *big.Int) (err error) {
	var validatorAccount types.ValidatorAccount
	if k.CheckIfValidatorAccountExists(ctx, valID) {
		validatorAccount, err = k.GetValidatorAccountByValID(ctx, valID)
		if err != nil {
			return err
		}
	} else {
		validatorAccount = types.ValidatorAccount{
			ID:            valID,
			RewardAmount:  big.NewInt(0).String(),
			SlashedAmount: big.NewInt(0).String(),
		}
	}

	// Add reward to reward balance
	rewardBalance, _ := big.NewInt(0).SetString(validatorAccount.RewardAmount, 10)
	updatedReward := big.NewInt(0).Add(reward, rewardBalance)
	validatorAccount.RewardAmount = updatedReward.String()
	err = k.AddValidatorAccount(ctx, validatorAccount)
	return
}

// GetRewardByValidatorID Returns Total Rewards of Validator
func (k *Keeper) GetRewardByValidatorID(ctx sdk.Context, valID types.ValidatorID) (*big.Int, error) {
	validatorAccount, err := k.GetValidatorAccountByValID(ctx, valID)
	if err != nil {
		return big.NewInt(0), err
	}
	validatorReward, _ := big.NewInt(0).SetString(validatorAccount.RewardAmount, 10)
	return validatorReward, nil
}

// SlashValidator will update validatoraccount with new slashed amount
func (k *Keeper) SlashValidator(ctx sdk.Context, valID types.ValidatorID, slashAmount *big.Int) (err error) {
	k.Logger(ctx).Info("Slashing validator - ", "valID", valID, "slashAmount", slashAmount)
	var validatorAccount types.ValidatorAccount
	if k.CheckIfValidatorAccountExists(ctx, valID) {
		validatorAccount, err = k.GetValidatorAccountByValID(ctx, valID)
		if err != nil {
			return err
		}
	} else {
		validatorAccount = types.ValidatorAccount{
			ID:            valID,
			RewardAmount:  big.NewInt(0).String(),
			SlashedAmount: big.NewInt(0).String(),
		}
	}

	// Add slashamount to slash balance
	slashBalance, _ := big.NewInt(0).SetString(validatorAccount.SlashedAmount, 10)
	updatedSlash := big.NewInt(0).Add(slashAmount, slashBalance)
	validatorAccount.SlashedAmount = updatedSlash.String()
	k.Logger(ctx).Info("Validator account after slashing - ", "valaccount", validatorAccount)
	err = k.AddValidatorAccount(ctx, validatorAccount)
	return
}

// GetSlashedAmountByValidatorID returns total slashed amount of validator
func (k *Keeper) GetSlashedAmountByValidatorID(ctx sdk.Context, valID types.ValidatorID) (*big.Int, error) {
	validatorAccount, err := k.GetValidatorAccountByValID(ctx, valID)
	if err != nil {
		return big.NewInt(0), err
	}

	slashedAmount, _ := big.NewInt(0).SetString(validatorAccount.SlashedAmount, 10)
	return slashedAmount, nil
}

// GetValidatorAccountByValID will return ValidatorAccount of valID
func (k *Keeper) GetValidatorAccountByValID(ctx sdk.Context, valID types.ValidatorID) (validatorAccount types.ValidatorAccount, err error) {

	// check if validator account exists
	if !k.CheckIfValidatorAccountExists(ctx, valID) {
		return validatorAccount, errors.New("Validator Account not found")
	}

	store := ctx.KVStore(k.storeKey)
	key := GetValidatorAccountMapKey(valID.Bytes())

	// unmarshall validator account and return
	validatorAccount, err = types.UnMarshallValidatorAccount(k.cdc, store.Get(key))
	if err != nil {
		return validatorAccount, err
	}

	return validatorAccount, nil
}

// CheckIfValidatorAccountExists will return true if validator account exists
func (k *Keeper) CheckIfValidatorAccountExists(ctx sdk.Context, valID types.ValidatorID) (ok bool) {
	store := ctx.KVStore(k.storeKey)
	key := GetValidatorAccountMapKey(valID.Bytes())
	if !store.Has(key) {
		return false
	}
	return true
}

// AddValidatorAccount adds ValidatorAccount index with ValID
func (k *Keeper) AddValidatorAccount(ctx sdk.Context, validatorAccount types.ValidatorAccount) error {
	store := ctx.KVStore(k.storeKey)
	// marshall validator account
	bz, err := types.MarshallValidatorAccount(k.cdc, validatorAccount)
	if err != nil {
		return err
	}

	store.Set(GetValidatorAccountMapKey(validatorAccount.ID.Bytes()), bz)
	k.Logger(ctx).Debug("ValidatorAccount Stored", "key", hex.EncodeToString(GetValidatorAccountMapKey(validatorAccount.ID.Bytes())), "validatoraccount", validatorAccount.String())
	return nil
}

// GetAllValidatorAccounts returns all validatorAccounts
func (k *Keeper) GetAllValidatorAccounts(ctx sdk.Context) (validatorAccounts []types.ValidatorAccount) {
	// iterate through validatorAccounts and create validatorAccounts update array
	k.IterateValidatorAccountsAndApplyFn(ctx, func(validatorAccount types.ValidatorAccount) error {
		// append to list of validatorUpdates
		validatorAccounts = append(validatorAccounts, validatorAccount)
		return nil
	})

	return
}

// IterateValidatorAccountsAndApplyFn iterate validatorAccounts and apply the given function.
func (k *Keeper) IterateValidatorAccountsAndApplyFn(ctx sdk.Context, f func(validatorAccount types.ValidatorAccount) error) {
	store := ctx.KVStore(k.storeKey)

	// get validator iterator
	iterator := sdk.KVStorePrefixIterator(store, ValidatorAccountMapKey)
	defer iterator.Close()

	// loop through validatoraccounts
	for ; iterator.Valid(); iterator.Next() {
		// unmarshall validatoraccount
		validatorAccount, _ := types.UnMarshallValidatorAccount(k.cdc, iterator.Value())
		// call function and return if required
		if err := f(validatorAccount); err != nil {
			return
		}
	}
}

// CalculateSignerRewards calculates new rewards for signers
func (k *Keeper) CalculateSignerRewards(ctx sdk.Context, voteBytes []byte, sigInput []byte) (map[types.ValidatorID]*big.Int, error) {
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

	currentCheckpointReward := k.ComputeCurrentCheckpointReward(ctx, totalSignerPower)
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
		k.RewardValidator(ctx, valID, reward)
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

// ComputeCurrentCheckpointReward returns the reward to be distributed for current checkpoint
func (k *Keeper) ComputeCurrentCheckpointReward(ctx sdk.Context, totalSignerPower int64) *big.Float {
	checkpointReward := k.GetCheckpointReward(ctx)
	currValSet := k.GetValidatorSet(ctx)
	totalPower := currValSet.TotalVotingPower()
	bigval := big.NewInt(0).Mul(checkpointReward, big.NewInt(totalSignerPower))

	bigfloatval := new(big.Float)
	bigfloatval.SetInt(bigval)

	totalPow := new(big.Float)
	totalPow.SetInt(big.NewInt(totalPower))

	return new(big.Float).Quo(bigfloatval, totalPow)
}

// GetCheckpointReward returns the reward Amount
func (k *Keeper) GetCheckpointReward(ctx sdk.Context) *big.Int {
	var checkpointRewardBytes []byte
	k.paramSpace.Get(ctx, ParamStoreKeyCheckpointReward, &checkpointRewardBytes)
	checkpointReward := big.NewInt(0).SetBytes(checkpointRewardBytes)
	return checkpointReward
}

// SetCheckpointReward sets the checkpoint reward Amount
func (k *Keeper) SetCheckpointReward(ctx sdk.Context, checkpointReward *big.Int) {
	k.paramSpace.Set(ctx, ParamStoreKeyCheckpointReward, checkpointReward.Bytes())
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
