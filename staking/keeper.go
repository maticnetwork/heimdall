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

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

var (
	DefaultValue = []byte{0x01} // Value to store in CacheCheckpoint and CacheCheckpointACK & ValidatorSetChange Flag

	ValidatorsKey          = []byte{0x21} // prefix for each key to a validator
	ValidatorMapKey        = []byte{0x22} // prefix for each key for validator map
	CurrentValidatorSetKey = []byte{0x23} // Key to store current validator set
	DividendAccountMapKey  = []byte{0x42} // prefix for each key for Dividend Account Map
	StakingSequenceKey     = []byte{0x24} // prefix for each key for staking sequence map
)

type AckRetriever interface {
	GetACKCount(ctx sdk.Context) uint64
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
		paramSpace:   paramSpace.WithKeyTable(types.ParamKeyTable()),
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
	ackCount := k.ackRetriever.GetACKCount(ctx)

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
	ackCount := k.ackRetriever.GetACKCount(ctx)

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

// AddDividendAccount adds DividendAccount index with DividendID
func (k *Keeper) AddDividendAccount(ctx sdk.Context, dividendAccount hmTypes.DividendAccount) error {
	store := ctx.KVStore(k.storeKey)
	// marshall validator account
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
	k.IterateDividendAccountsAndApplyFn(ctx, func(dividendAccount hmTypes.DividendAccount) error {
		// append to list of dividendUpdates
		dividendAccounts = append(dividendAccounts, dividendAccount)
		return nil
	})

	return
}

// IterateDividendAccountsAndApplyFn iterate dividendAccounts and apply the given function.
func (k *Keeper) IterateDividendAccountsAndApplyFn(ctx sdk.Context, f func(dividendAccount hmTypes.DividendAccount) error) {
	store := ctx.KVStore(k.storeKey)

	// get validator iterator
	iterator := sdk.KVStorePrefixIterator(store, DividendAccountMapKey)
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

// 1. Delegator is updated with Validator ID.
// 2. VotingPower of the bonded validator is updated.
// 3. shares are added to Delegator proportional to his stake and exchange rate. // delegatorshares = (delegatorstake / exchangeRate)
// 4. Exchange rate is calculated instantly.  //   ExchangeRate = (delegatedpower + delegatorRewardPool) / totaldelegatorshares
// 5. TotalDelegatorShares of bonded validator is updated.
// 6. DelegatedPower of bonded validator is updated.
func (k *Keeper) BondDelegator(ctx sdk.Context, delegatorID hmTypes.DelegatorID, valID hmTypes.ValidatorID, amount *big.Int) (err error) {

	var dividendAccount hmTypes.DividendAccount
	if k.CheckIfDividendAccountExists(ctx, hmTypes.DividendAccountID(delegatorID)) {
		dividendAccount, err = k.GetDividendAccountByID(ctx, hmTypes.DividendAccountID(delegatorID))
		if err != nil {
			return err
		}
	} else {
		dividendAccount = hmTypes.DividendAccount{
			ID:            hmTypes.DividendAccountID(valID),
			RewardAmount:  big.NewInt(0).String(),
			Shares:        big.NewInt(0).String(),
			SlashedAmount: big.NewInt(0).String(),
		}
	}

	// 3. VotingPower of the bonded validator is updated.
	// pull validator from store
	validator, ok := k.GetValidatorFromValID(ctx, valID)
	if !ok {
		k.Logger(ctx).Error("Fetching of bonded validator from store failed", "validatorId", valID)
		return errors.New("Validator not found")
	}

	// 4. shares are added to Delegator proportional to his stake and exchange rate.
	// newDelegatorshares = (delegatorstake / exchangeRate)
	delegatorStake := big.NewFloat(0).SetInt(amount)
	newDelegatorshares := new(big.Float).Quo(delegatorStake, validator.ExchangeRate())
	newShares := new(big.Int)
	newDelegatorshares.Int(newShares)

	// add shares to delegator dividend account
	oldShares, _ := big.NewInt(0).SetString(dividendAccount.Shares, 10)
	dividendAccount.Shares = big.NewInt(0).Add(oldShares, newShares).String()

	// 6. TotalDelegatorShares of bonded validator is updated.
	oldTotalDelegatorShares, _ := big.NewInt(0).SetString(validator.TotalDelegatorShares, 10)
	validator.TotalDelegatorShares = big.NewInt(0).Add(oldTotalDelegatorShares, newShares).String()

	p, err := helper.GetPowerFromAmount(amount)
	if err != nil {
		k.Logger(ctx).Error("Unable to convert amount to power", "Amount", amount)
		return errors.New("Unable to convert amount to power")
	}

	validator.VotingPower += p.Int64()

	// 7. DelegatedPower of bonded validator is updated.
	validator.DelegatedPower += p.Int64()

	// save delegator account
	err = k.AddDividendAccount(ctx, dividendAccount)
	if err != nil {
		k.Logger(ctx).Error("Unable to update delegatorAccount", "error", err, "DelegatorID", delegatorID)
		return errors.New("Delegator DividendAccount updation failed")
	}

	// save validator
	err = k.AddValidator(ctx, validator)
	if err != nil {
		k.Logger(ctx).Error("Unable to update validator", "error", err, "ValidatorID", validator.ID)
		return errors.New("Validator updation failed")
	}
	return nil
}

// HandleMsgDelegatorUnBond msg delegator unbond with validator
// ** stake calculations **
// 1. On Bonding event, Validator will send MsgDelegatorUnBond transaction to heimdall.
// 2. Delegator is updated with Validator ID = 0.
// 3. VotingPower of bonded validator is reduced.
// 4. DelegatedPower of the bonded validator is reduced after reward calculation.

// ** reward calculations **
// 1. Exchange rate is calculated instantly.  ExchangeRate = (delegatedpower + delegatorRewardPool) / totaldelegatorshares
// 2. Based on exchange rate and no of shares delegator holds, totalReturns for delegator is calculated.  `totalReturns = exchangeRate * noOfShares`
// 3. Delegator RewardAmount += totalReturns - delegatorVotingPower
// 4. Add RewardAmount to DelegatorAccount .
// 5. Reduce TotalDelegatorShares of bonded validator.
// 6. Reduce DelgatorRewardPool of bonded validator.
// 7. make shares = 0 on Delegator Account.
// UnBondDelegator
func (k *Keeper) UnBondDelegator(ctx sdk.Context, delegatorID hmTypes.DelegatorID, valID hmTypes.ValidatorID, amount *big.Int) (err error) {

	// delegatorAccount, err := k.GetDividendAccountByID(ctx, hmTypes.DividendAccountID(delegatorID))
	// if err != nil {
	// 	k.Logger(ctx).Error("Fetching of delegator dividend account from store failed", "delegatorId", delegatorID)
	// 	return errors.New("Delegator DividendAccount not found")
	// }

	// // 3. VotingPower of the bonded validator is updated.
	// // pull validator from store
	// validator, ok := k.GetValidatorFromValID(ctx, valID)
	// if !ok {
	// 	k.Logger(ctx).Error("Fetching of bonded validator from store failed", "validatorId", valID)
	// 	return errors.New("Validator not found")
	// }

	// // Get shares of delegator account
	// delegatorshares := delegatorAccount.Shares

	// // 6. TotalDelegatorShares of bonded validator is updated.
	// validator.TotalDelegatorShares -= delegatorshares

	// validator.VotingPower -= delegator.VotingPower

	// // calculate rewards.
	// totalDelegatorReturns := validator.ExchangeRate() * delegatorshares

	// rewardAmount := totalDelegatorReturns - float32(delegator.VotingPower)

	// validator.DelgatorRewardPool = int64(float32(validator.DelgatorRewardPool) - rewardAmount)

	// // 7. DelegatedPower of bonded validator is updated.
	// validator.DelegatedPower -= delegator.VotingPower

	// // save validator
	// err = k.sk.AddValidator(ctx, validator)
	// if err != nil {
	// 	k.Logger(ctx).Error("Unable to update validator", "error", err, "ValidatorID", validator.ID)
	// 	return errors.New("error adding validator to store")
	// }

	// // 2. update validator ID of delegator.
	// delegator.ValID = 0

	// // update last udpated
	// delegator.LastUpdated = lastUpdated

	// delegatorAccount.Shares = 0

	// // save delegator account
	// err = k.AddDelegatorAccount(ctx, delegatorAccount)
	// if err != nil {
	// 	k.Logger(ctx).Error("Unable to update delegatorAccount", "error", err, "DelegatorID", delegator.ID)
	// 	return errors.New("DelegatorAccount updation failed")
	// }

	// // save delegator
	// err = k.AddDelegator(ctx, delegator)
	// if err != nil {
	// 	k.Logger(ctx).Error("Unable to update delegator", "error", err, "DelegatorID", delegator.ID)
	// 	return errors.New("error adding delegator to store")
	// }

	return nil
}

// RewardValidator will update validator dividend account with new reward
func (k *Keeper) RewardValidator(ctx sdk.Context, valID hmTypes.ValidatorID, totalRewards *big.Int) (err error) {
	// Divide total reward between validator and his delegator pool.
	validator, ok := k.GetValidatorFromValID(ctx, valID)
	if !ok {
		return errors.New("Validator not found")
	}

	// calculate validator Reward
	valPower := (validator.VotingPower - validator.DelegatedPower)
	bigvalPow := new(big.Float)
	bigvalPow.SetFloat64(float64(valPower))
	bigvalTotalPow := new(big.Float)
	bigvalTotalPow.SetFloat64(float64(validator.VotingPower))
	valRew := new(big.Float).Quo(bigvalPow, bigvalTotalPow)
	valReward := big.NewInt(0)
	valReward, _ = valRew.Int(valReward)

	// calculate delegator reward
	delPower := validator.DelegatedPower
	bigdelPow := new(big.Float)
	bigdelPow.SetFloat64(float64(delPower))
	bigvalTotalPow = new(big.Float)
	bigvalTotalPow.SetFloat64(float64(validator.VotingPower))
	delRew := new(big.Float).Quo(bigdelPow, bigvalTotalPow)
	delReward := big.NewInt(0)
	delReward, _ = delRew.Int(delReward)

	var validatorAccount hmTypes.DividendAccount
	if k.CheckIfDividendAccountExists(ctx, hmTypes.DividendAccountID(valID)) {
		validatorAccount, err = k.GetDividendAccountByID(ctx, hmTypes.DividendAccountID(valID))
		if err != nil {
			return err
		}
	} else {
		validatorAccount = hmTypes.DividendAccount{
			ID:            hmTypes.DividendAccountID(valID),
			RewardAmount:  big.NewInt(0).String(),
			SlashedAmount: big.NewInt(0).String(),
		}
	}

	// Add reward to reward balance
	rewardBalance, _ := big.NewInt(0).SetString(validatorAccount.RewardAmount, 10)
	updatedReward := big.NewInt(0).Add(valReward, rewardBalance)
	validatorAccount.RewardAmount = updatedReward.String()

	// Add delegator reward to delegatorRewardPool
	delegatorPoolRewardBalance, _ := big.NewInt(0).SetString(validator.DelgatorRewardPool, 10)
	updatedDelegatorPoolReward := big.NewInt(0).Add(delReward, delegatorPoolRewardBalance)
	validator.DelgatorRewardPool = updatedDelegatorPoolReward.String()

	if err = k.AddValidator(ctx, validator); err != nil {
		return err
	}

	if err = k.AddDividendAccount(ctx, validatorAccount); err != nil {
		return err
	}
	return
}

// GetRewardByDividendAccountID Returns Total Rewards of Dividend Account
func (k *Keeper) GetRewardByDividendAccountID(ctx sdk.Context, dividendAccountID hmTypes.DividendAccountID) (*big.Int, error) {
	dividendAccount, err := k.GetDividendAccountByID(ctx, dividendAccountID)
	if err != nil {
		return big.NewInt(0), err
	}
	reward, _ := big.NewInt(0).SetString(dividendAccount.RewardAmount, 10)
	return reward, nil
}

// CalculateSignerRewards calculates new rewards for signers
func (k *Keeper) CalculateSignerRewards(ctx sdk.Context, voteBytes []byte, sigInput []byte, checkpointReward *big.Int) (map[hmTypes.ValidatorID]*big.Int, error) {
	signerRewards := make(map[hmTypes.ValidatorID]*big.Int)
	signerPower := make(map[hmTypes.ValidatorID]int64)

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

		pubKey := hmTypes.NewPubKey(pKey)
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
func (k *Keeper) UpdateValidatorRewards(ctx sdk.Context, valrewards map[hmTypes.ValidatorID]*big.Int) {
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

// GetProposerBonusPercent returns the proposer to signer reward
func (k *Keeper) GetProposerBonusPercent(ctx sdk.Context) int64 {
	var proposerBonusPercent int64
	k.paramSpace.Get(ctx, types.ParamStoreKeyProposerBonusPercent, &proposerBonusPercent)
	return proposerBonusPercent
}

// SetProposerBonusPercent sets the Proposer to signer reward
func (k *Keeper) SetProposerBonusPercent(ctx sdk.Context, proposerBonusPercent int64) {
	k.paramSpace.Set(ctx, types.ParamStoreKeyProposerBonusPercent, proposerBonusPercent)
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
