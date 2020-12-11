package types

import (
	"bytes"
	"fmt"
	"math"
	"math/big"
	"sort"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto/merkle"
)

const (
	// MaxTotalVotingPower - the maximum allowed total voting power.
	// It needs to be sufficiently small to, in all cases:
	// 1. prevent clipping in incrementProposerPriority()
	// 2. let (diff+diffMax-1) not overflow in IncrementProposerPriority()
	// (Proof of 1 is tricky, left to the reader).
	// It could be higher, but this is sufficiently large for our purposes,
	// and leaves room for defensive purposes.
	MaxTotalVotingPower = int64(math.MaxInt64) / 8

	// PriorityWindowSizeFactor - is a constant that when multiplied with the
	// total voting power gives the maximum allowed distance between validator
	// priorities.
	PriorityWindowSizeFactor = 2
)

// ErrTotalVotingPowerOverflow is returned if the total voting power of the
// resulting validator set exceeds MaxTotalVotingPower.
var ErrTotalVotingPowerOverflow = fmt.Errorf("total voting power of resulting valset exceeds max %d",
	MaxTotalVotingPower)

// NewValidatorSet initializes a ValidatorSet by copying over the values from
// `valz`, a list of Validators. If valz is nil or empty, the new ValidatorSet
// will have an empty list of Validators.
//
// The addresses of validators in `valz` must be unique otherwise the function
// panics.
//
// Note the validator set size has an implied limit equal to that of the
// MaxVotesCount - commits by a validator set larger than this will fail
// validation.
func NewValidatorSet(valz []*Validator) *ValidatorSet {
	vals := &ValidatorSet{}
	err := vals.updateWithChangeSet(valz, false)
	if err != nil {
		panic(fmt.Sprintf("Cannot create validator set: %v", err))
	}
	if len(valz) > 0 {
		vals.IncrementProposerPriority(1)
	}
	return vals
}

func (vals *ValidatorSet) ValidateBasic() error {
	if vals.IsNilOrEmpty() {
		return errors.New("validator set is nil or empty")
	}

	for idx, val := range vals.Validators {
		if err := val.ValidateBasic(); err != nil {
			return fmt.Errorf("invalid validator #%d: %w", idx, err)
		}
	}

	if err := vals.Proposer.ValidateBasic(); err != nil {
		return fmt.Errorf("proposer failed validate basic, error: %w", err)
	}

	return nil
}

// IsNilOrEmpty returns true if validator set is nil or empty.
func (vals *ValidatorSet) IsNilOrEmpty() bool {
	return vals == nil || len(vals.Validators) == 0
}

// CopyIncrementProposerPriority increments ProposerPriority and updates the
// proposer on a copy, and returns it.
func (vals *ValidatorSet) CopyIncrementProposerPriority(times int32) *ValidatorSet {
	copy := vals.Copy()
	copy.IncrementProposerPriority(times)
	return copy
}

// IncrementProposerPriority increments ProposerPriority of each validator and
// updates the proposer. Panics if validator set is empty.
// `times` must be positive.
func (vals *ValidatorSet) IncrementProposerPriority(times int32) {
	if vals.IsNilOrEmpty() {
		panic("empty validator set")
	}
	if times <= 0 {
		panic("Cannot call IncrementProposerPriority with non-positive times")
	}

	// Cap the difference between priorities to be proportional to 2*totalPower by
	// re-normalizing priorities, i.e., rescale all priorities by multiplying with:
	//  2*totalVotingPower/(maxPriority - minPriority)
	diffMax := PriorityWindowSizeFactor * vals.GetTotalVotingPower()
	vals.RescalePriorities(diffMax)
	vals.shiftByAvgProposerPriority()

	var proposer *Validator
	// Call IncrementProposerPriority(1) times times.
	for i := int32(0); i < times; i++ {
		proposer = vals.incrementProposerPriority()
	}

	vals.Proposer = proposer
}

// RescalePriorities rescales the priorities such that the distance between the
// maximum and minimum is smaller than `diffMax`. Panics if validator set is
// empty.
func (vals *ValidatorSet) RescalePriorities(diffMax int64) {
	if vals.IsNilOrEmpty() {
		panic("empty validator set")
	}
	// NOTE: This check is merely a sanity check which could be
	// removed if all tests would init. voting power appropriately;
	// i.e. diffMax should always be > 0
	if diffMax <= 0 {
		return
	}

	// Calculating ceil(diff/diffMax):
	// Re-normalization is performed by dividing by an integer for simplicity.
	// NOTE: This may make debugging priority issues easier as well.
	diff := computeMaxMinPriorityDiff(vals)
	ratio := (diff + diffMax - 1) / diffMax
	if diff > diffMax {
		for _, val := range vals.Validators {
			val.ProposerPriority /= ratio
		}
	}
}

func (vals *ValidatorSet) incrementProposerPriority() *Validator {
	for _, val := range vals.Validators {
		// Check for overflow for sum.
		newPrio := safeAddClip(val.ProposerPriority, val.VotingPower)
		val.ProposerPriority = newPrio
	}
	// Decrement the validator with most ProposerPriority.
	mostest := vals.getValWithMostPriority()
	// Mind the underflow.
	mostest.ProposerPriority = safeSubClip(mostest.ProposerPriority, vals.GetTotalVotingPower())

	return mostest
}

// Should not be called on an empty validator set.
func (vals *ValidatorSet) computeAvgProposerPriority() int64 {
	n := int64(len(vals.Validators))
	sum := big.NewInt(0)
	for _, val := range vals.Validators {
		sum.Add(sum, big.NewInt(val.ProposerPriority))
	}
	avg := sum.Div(sum, big.NewInt(n))
	if avg.IsInt64() {
		return avg.Int64()
	}

	// This should never happen: each val.ProposerPriority is in bounds of int64.
	panic(fmt.Sprintf("Cannot represent avg ProposerPriority as an int64 %v", avg))
}

// Compute the difference between the max and min ProposerPriority of that set.
func computeMaxMinPriorityDiff(vals *ValidatorSet) int64 {
	if vals.IsNilOrEmpty() {
		panic("empty validator set")
	}
	max := int64(math.MinInt64)
	min := int64(math.MaxInt64)
	for _, v := range vals.Validators {
		if v.ProposerPriority < min {
			min = v.ProposerPriority
		}
		if v.ProposerPriority > max {
			max = v.ProposerPriority
		}
	}
	diff := max - min
	if diff < 0 {
		return -1 * diff
	}
	return diff
}

func (vals *ValidatorSet) getValWithMostPriority() *Validator {
	var res *Validator
	for _, val := range vals.Validators {
		res = res.CompareProposerPriority(val)
	}
	return res
}

func (vals *ValidatorSet) shiftByAvgProposerPriority() {
	if vals.IsNilOrEmpty() {
		panic("empty validator set")
	}
	avgProposerPriority := vals.computeAvgProposerPriority()
	for _, val := range vals.Validators {
		val.ProposerPriority = safeSubClip(val.ProposerPriority, avgProposerPriority)
	}
}

// Makes a copy of the validator list.
func validatorListCopy(valsList []*Validator) []*Validator {
	if valsList == nil {
		return nil
	}
	valsCopy := make([]*Validator, len(valsList))
	for i, val := range valsList {
		valsCopy[i] = val.Copy()
	}
	return valsCopy
}

// Copy each validator into a new ValidatorSet.
func (vals *ValidatorSet) Copy() *ValidatorSet {
	return &ValidatorSet{
		Validators:       validatorListCopy(vals.Validators),
		Proposer:         vals.Proposer,
		TotalVotingPower: vals.TotalVotingPower,
	}
}

// HasAddress returns true if address given is in the validator set, false -
// otherwise.
func (vals *ValidatorSet) HasAddress(address []byte) bool {
	for _, val := range vals.Validators {
		if bytes.Equal(val.GetSigner(), address) {
			return true
		}
	}
	return false
}

// GetByAddress returns an index of the validator with address and validator
// itself (copy) if found. Otherwise, -1 and nil are returned.
func (vals *ValidatorSet) GetByAddress(address []byte) (index int32, val *Validator) {
	for idx, val := range vals.Validators {
		if bytes.Equal(val.GetSigner(), address) {
			return int32(idx), val.Copy()
		}
	}
	return -1, nil
}

// GetByIndex returns the validator's address and validator itself (copy) by
// index.
// It returns nil values if index is less than 0 or greater or equal to
// len(ValidatorSet.Validators).
func (vals *ValidatorSet) GetByIndex(index int32) (address []byte, val *Validator) {
	if index < 0 || int(index) >= len(vals.Validators) {
		return nil, nil
	}
	val = vals.Validators[index]
	return val.GetSigner(), val.Copy()
}

// TotalValidators returns the length of the validator set.
func (vals *ValidatorSet) TotalValidators() int {
	return len(vals.Validators)
}

// Forces recalculation of the set's total voting power.
// Panics if total voting power is bigger than MaxTotalVotingPower.
func (vals *ValidatorSet) updateTotalVotingPower() {
	sum := int64(0)
	for _, val := range vals.Validators {
		// mind overflow
		sum = safeAddClip(sum, val.VotingPower)
		if sum > MaxTotalVotingPower {
			panic(fmt.Sprintf(
				"Total voting power should be guarded to not exceed %v; got: %v",
				MaxTotalVotingPower,
				sum))
		}
	}

	vals.TotalVotingPower = sum
}

// GetTotalVotingPower returns the sum of the voting powers of all validators.
// It recomputes the total voting power if required.
func (vals *ValidatorSet) GetTotalVotingPower() int64 {
	if vals.TotalVotingPower == 0 {
		vals.updateTotalVotingPower()
	}
	return vals.TotalVotingPower
}

// GetProposer returns the current proposer. If the validator set is empty, nil
// is returned.
func (vals *ValidatorSet) GetProposer() (proposer *Validator) {
	if len(vals.Validators) == 0 {
		return nil
	}
	if vals.Proposer == nil {
		vals.Proposer = vals.findProposer()
	}
	return vals.Proposer.Copy()
}

func (vals *ValidatorSet) findProposer() *Validator {
	var proposer *Validator
	for _, val := range vals.Validators {
		if proposer == nil || !bytes.Equal(val.GetSigner(), proposer.GetSigner()) {
			proposer = proposer.CompareProposerPriority(val)
		}
	}
	return proposer
}

// Hash returns the Merkle root hash build using validators (as leaves) in the
// set.
func (vals *ValidatorSet) Hash() []byte {
	bzs := make([][]byte, len(vals.Validators))
	for i, val := range vals.Validators {
		bzs[i] = val.Bytes()
	}
	return merkle.HashFromByteSlices(bzs)
}

// Iterate will run the given function over the set.
func (vals *ValidatorSet) Iterate(fn func(index int, val *Validator) bool) {
	for i, val := range vals.Validators {
		stop := fn(i, val.Copy())
		if stop {
			break
		}
	}
}

// Checks changes against duplicates, splits the changes in updates and
// removals, sorts them by address.
//
// Returns:
// updates, removals - the sorted lists of updates and removals
// err - non-nil if duplicate entries or entries with negative voting power are seen
//
// No changes are made to 'origChanges'.
func processChanges(origChanges []*Validator) (updates, removals []*Validator, err error) {
	// Make a deep copy of the changes and sort by address.
	changes := validatorListCopy(origChanges)
	sort.Sort(ValidatorsByAddress(changes))

	removals = make([]*Validator, 0, len(changes))
	updates = make([]*Validator, 0, len(changes))
	var prevAddr sdk.AccAddress

	// Scan changes by address and append valid validators to updates or removals lists.
	for _, valUpdate := range changes {
		if bytes.Equal(valUpdate.GetSigner(), prevAddr) {
			err = fmt.Errorf("duplicate entry %v in %v", valUpdate, changes)
			return nil, nil, err
		}

		switch {
		case valUpdate.VotingPower < 0:
			err = fmt.Errorf("voting power can't be negative: %d", valUpdate.VotingPower)
			return nil, nil, err
		case valUpdate.VotingPower > MaxTotalVotingPower:
			err = fmt.Errorf("to prevent clipping/overflow, voting power can't be higher than %d, got %d",
				MaxTotalVotingPower, valUpdate.VotingPower)
			return nil, nil, err
		case valUpdate.VotingPower == 0:
			removals = append(removals, valUpdate)
		default:
			updates = append(updates, valUpdate)
		}

		prevAddr = valUpdate.GetSigner()
	}

	return updates, removals, err
}

// verifyUpdates verifies a list of updates against a validator set, making sure the allowed
// total voting power would not be exceeded if these updates would be applied to the set.
//
// Inputs:
// updates - a list of proper validator changes, i.e. they have been verified by processChanges for duplicates
//   and invalid values.
// vals - the original validator set. Note that vals is NOT modified by this function.
// removedPower - the total voting power that will be removed after the updates are verified and applied.
//
// Returns:
// tvpAfterUpdatesBeforeRemovals -  the new total voting power if these updates would be applied without the removals.
//   Note that this will be < 2 * MaxTotalVotingPower in case high power validators are removed and
//   validators are added/ updated with high power values.
//
// err - non-nil if the maximum allowed total voting power would be exceeded
func verifyUpdates(
	updates []*Validator,
	vals *ValidatorSet,
	removedPower int64,
) (tvpAfterUpdatesBeforeRemovals int64, err error) {

	delta := func(update *Validator, vals *ValidatorSet) int64 {
		_, val := vals.GetByAddress(update.GetSigner())
		if val != nil {
			return update.VotingPower - val.VotingPower
		}
		return update.VotingPower
	}

	updatesCopy := validatorListCopy(updates)
	sort.Slice(updatesCopy, func(i, j int) bool {
		return delta(updatesCopy[i], vals) < delta(updatesCopy[j], vals)
	})

	tvpAfterRemovals := vals.GetTotalVotingPower() - removedPower
	for _, upd := range updatesCopy {
		tvpAfterRemovals += delta(upd, vals)
		if tvpAfterRemovals > MaxTotalVotingPower {
			return 0, ErrTotalVotingPowerOverflow
		}
	}
	return tvpAfterRemovals + removedPower, nil
}

func numNewValidators(updates []*Validator, vals *ValidatorSet) int {
	numNewValidators := 0
	for _, valUpdate := range updates {
		if !vals.HasAddress(valUpdate.GetSigner()) {
			numNewValidators++
		}
	}
	return numNewValidators
}

// computeNewPriorities computes the proposer priority for the validators not present in the set based on
// 'updatedTotalVotingPower'.
// Leaves unchanged the priorities of validators that are changed.
//
// 'updates' parameter must be a list of unique validators to be added or updated.
//
// 'updatedTotalVotingPower' is the total voting power of a set where all updates would be applied but
//   not the removals. It must be < 2*MaxTotalVotingPower and may be close to this limit if close to
//   MaxTotalVotingPower will be removed. This is still safe from overflow since MaxTotalVotingPower is maxInt64/8.
//
// No changes are made to the validator set 'vals'.
func computeNewPriorities(updates []*Validator, vals *ValidatorSet, updatedTotalVotingPower int64) {
	for _, valUpdate := range updates {
		address := valUpdate.GetSigner()
		_, val := vals.GetByAddress(address)
		if val == nil {
			// add val
			// Set ProposerPriority to -C*totalVotingPower (with C ~= 1.125) to make sure validators can't
			// un-bond and then re-bond to reset their (potentially previously negative) ProposerPriority to zero.
			//
			// Contract: updatedVotingPower < 2 * MaxTotalVotingPower to ensure ProposerPriority does
			// not exceed the bounds of int64.
			//
			// Compute ProposerPriority = -1.125*totalVotingPower == -(updatedVotingPower + (updatedVotingPower >> 3)).
			valUpdate.ProposerPriority = -(updatedTotalVotingPower + (updatedTotalVotingPower >> 3))
		} else {
			valUpdate.ProposerPriority = val.ProposerPriority
		}
	}

}

// Merges the vals' validator list with the updates list.
// When two elements with same address are seen, the one from updates is selected.
// Expects updates to be a list of updates sorted by address with no duplicates or errors,
// must have been validated with verifyUpdates() and priorities computed with computeNewPriorities().
func (vals *ValidatorSet) applyUpdates(updates []*Validator) {
	existing := vals.Validators
	sort.Sort(ValidatorsByAddress(existing))

	merged := make([]*Validator, len(existing)+len(updates))
	i := 0

	for len(existing) > 0 && len(updates) > 0 {
		if bytes.Compare(existing[0].GetSigner(), updates[0].GetSigner()) < 0 { // unchanged validator
			merged[i] = existing[0]
			existing = existing[1:]
		} else {
			// Apply add or update.
			merged[i] = updates[0]
			if bytes.Equal(existing[0].GetSigner(), updates[0].GetSigner()) {
				// Validator is present in both, advance existing.
				existing = existing[1:]
			}
			updates = updates[1:]
		}
		i++
	}

	// Add the elements which are left.
	for j := 0; j < len(existing); j++ {
		merged[i] = existing[j]
		i++
	}
	// OR add updates which are left.
	for j := 0; j < len(updates); j++ {
		merged[i] = updates[j]
		i++
	}

	vals.Validators = merged[:i]
}

// Checks that the validators to be removed are part of the validator set.
// No changes are made to the validator set 'vals'.
func verifyRemovals(deletes []*Validator, vals *ValidatorSet) (votingPower int64, err error) {
	removedVotingPower := int64(0)
	for _, valUpdate := range deletes {
		address := valUpdate.GetSigner()
		_, val := vals.GetByAddress(address)
		if val == nil {
			return removedVotingPower, fmt.Errorf("failed to find validator %X to remove", address)
		}
		removedVotingPower += val.VotingPower
	}
	if len(deletes) > len(vals.Validators) {
		panic("more deletes than validators")
	}
	return removedVotingPower, nil
}

// Removes the validators specified in 'deletes' from validator set 'vals'.
// Should not fail as verification has been done before.
// Expects vals to be sorted by address (done by applyUpdates).
func (vals *ValidatorSet) applyRemovals(deletes []*Validator) {
	existing := vals.Validators

	merged := make([]*Validator, len(existing)-len(deletes))
	i := 0

	// Loop over deletes until we removed all of them.
	for len(deletes) > 0 {
		if bytes.Equal(existing[0].GetSigner(), deletes[0].GetSigner()) {
			deletes = deletes[1:]
		} else { // Leave it in the resulting slice.
			merged[i] = existing[0]
			i++
		}
		existing = existing[1:]
	}

	// Add the elements which are left.
	for j := 0; j < len(existing); j++ {
		merged[i] = existing[j]
		i++
	}

	vals.Validators = merged[:i]
}

// Main function used by UpdateWithChangeSet() and NewValidatorSet().
// If 'allowDeletes' is false then delete operations (identified by validators with voting power 0)
// are not allowed and will trigger an error if present in 'changes'.
// The 'allowDeletes' flag is set to false by NewValidatorSet() and to true by UpdateWithChangeSet().
func (vals *ValidatorSet) updateWithChangeSet(changes []*Validator, allowDeletes bool) error {
	if len(changes) == 0 {
		return nil
	}

	// Check for duplicates within changes, split in 'updates' and 'deletes' lists (sorted).
	updates, deletes, err := processChanges(changes)
	if err != nil {
		return err
	}

	if !allowDeletes && len(deletes) != 0 {
		return fmt.Errorf("cannot process validators with voting power 0: %v", deletes)
	}

	// Check that the resulting set will not be empty.
	if numNewValidators(updates, vals) == 0 && len(vals.Validators) == len(deletes) {
		return errors.New("applying the validator changes would result in empty set")
	}

	// Verify that applying the 'deletes' against 'vals' will not result in error.
	// Get the voting power that is going to be removed.
	removedVotingPower, err := verifyRemovals(deletes, vals)
	if err != nil {
		return err
	}

	// Verify that applying the 'updates' against 'vals' will not result in error.
	// Get the updated total voting power before removal. Note that this is < 2 * MaxTotalVotingPower
	tvpAfterUpdatesBeforeRemovals, err := verifyUpdates(updates, vals, removedVotingPower)
	if err != nil {
		return err
	}

	// Compute the priorities for updates.
	computeNewPriorities(updates, vals, tvpAfterUpdatesBeforeRemovals)

	// Apply updates and removals.
	vals.applyUpdates(updates)
	vals.applyRemovals(deletes)

	vals.updateTotalVotingPower() // will panic if total voting power > MaxTotalVotingPower

	// Scale and center.
	vals.RescalePriorities(PriorityWindowSizeFactor * vals.GetTotalVotingPower())
	vals.shiftByAvgProposerPriority()

	sort.Sort(ValidatorsByVotingPower(vals.Validators))

	return nil
}

// UpdateWithChangeSet attempts to update the validator set with 'changes'.
// It performs the following steps:
// - validates the changes making sure there are no duplicates and splits them in updates and deletes
// - verifies that applying the changes will not result in errors
// - computes the total voting power BEFORE removals to ensure that in the next steps the priorities
//   across old and newly added validators are fair
// - computes the priorities of new validators against the final set
// - applies the updates against the validator set
// - applies the removals against the validator set
// - performs scaling and centering of priority values
// If an error is detected during verification steps, it is returned and the validator set
// is not changed.
func (vals *ValidatorSet) UpdateWithChangeSet(changes []*Validator) error {
	return vals.updateWithChangeSet(changes, true)
}

// findPreviousProposer reverses the compare proposer priority function to find the validator
// with the lowest proposer priority which would have been the previous proposer.
//
// Is used when recreating a validator set from an existing array of validators.
func (vals *ValidatorSet) findPreviousProposer() *Validator {
	var previousProposer *Validator
	for _, val := range vals.Validators {
		if previousProposer == nil {
			previousProposer = val
			continue
		}
		if previousProposer == previousProposer.CompareProposerPriority(val) {
			previousProposer = val
		}
	}
	return previousProposer
}

//-----------------

// IsErrNotEnoughVotingPowerSigned returns true if err is
// ErrNotEnoughVotingPowerSigned.
func IsErrNotEnoughVotingPowerSigned(err error) bool {
	return errors.As(err, &ErrNotEnoughVotingPowerSigned{})
}

// ErrNotEnoughVotingPowerSigned is returned when not enough validators signed
// a commit.
type ErrNotEnoughVotingPowerSigned struct {
	Got    int64
	Needed int64
}

func (e ErrNotEnoughVotingPowerSigned) Error() string {
	return fmt.Sprintf("invalid commit -- insufficient voting power: got %d, needed more than %d", e.Got, e.Needed)
}

//----------------

// String returns a string representation of ValidatorSet.
//
// See StringIndented.
func (vals *ValidatorSet) String() string {
	return vals.StringIndented("")
}

// StringIndented returns an intended String.
//
// See Validator#String.
func (vals *ValidatorSet) StringIndented(indent string) string {
	if vals == nil {
		return "nil-ValidatorSet"
	}
	var valStrings []string
	vals.Iterate(func(index int, val *Validator) bool {
		valStrings = append(valStrings, val.String())
		return false
	})
	return fmt.Sprintf(`ValidatorSet{
%s  Proposer: %v
%s  Validators:
%s    %v
%s}`,
		indent, vals.GetProposer().String(),
		indent,
		indent, strings.Join(valStrings, "\n"+indent+"    "),
		indent)

}

//-------------------------------------

// ValidatorsByVotingPower implements sort.Interface for []*Validator based on
// the VotingPower and Address fields.
type ValidatorsByVotingPower []*Validator

func (valz ValidatorsByVotingPower) Len() int { return len(valz) }

func (valz ValidatorsByVotingPower) Less(i, j int) bool {
	if valz[i].VotingPower == valz[j].VotingPower {
		return bytes.Compare(valz[i].GetSigner(), valz[j].GetSigner()) == -1
	}
	return valz[i].VotingPower > valz[j].VotingPower
}

func (valz ValidatorsByVotingPower) Swap(i, j int) {
	valz[i], valz[j] = valz[j], valz[i]
}

// ValidatorsByAddress implements sort.Interface for []*Validator based on
// the Address field.
type ValidatorsByAddress []*Validator

func (valz ValidatorsByAddress) Len() int { return len(valz) }

func (valz ValidatorsByAddress) Less(i, j int) bool {
	return bytes.Compare(valz[i].GetSigner(), valz[j].GetSigner()) == -1
}

func (valz ValidatorsByAddress) Swap(i, j int) {
	valz[i], valz[j] = valz[j], valz[i]
}

// ValidatorSetFromExistingValidators takes an existing array of validators and rebuilds
// the exact same validator set that corresponds to it without changing the proposer priority or power
// if any of the validators fail validate basic then an empty set is returned.
func ValidatorSetFromExistingValidators(valz []*Validator) (*ValidatorSet, error) {
	for _, val := range valz {
		err := val.ValidateBasic()
		if err != nil {
			return nil, fmt.Errorf("can't create validator set: %w", err)
		}
	}
	vals := &ValidatorSet{
		Validators: valz,
	}
	vals.Proposer = vals.findPreviousProposer()
	vals.updateTotalVotingPower()
	sort.Sort(ValidatorsByVotingPower(vals.Validators))
	return vals, nil
}

// safe addition/subtraction/multiplication

func safeAdd(a, b int64) (int64, bool) {
	if b > 0 && a > math.MaxInt64-b {
		return -1, true
	} else if b < 0 && a < math.MinInt64-b {
		return -1, true
	}
	return a + b, false
}

func safeSub(a, b int64) (int64, bool) {
	if b > 0 && a < math.MinInt64+b {
		return -1, true
	} else if b < 0 && a > math.MaxInt64+b {
		return -1, true
	}
	return a - b, false
}

func safeAddClip(a, b int64) int64 {
	c, overflow := safeAdd(a, b)
	if overflow {
		if b < 0 {
			return math.MinInt64
		}
		return math.MaxInt64
	}
	return c
}

func safeSubClip(a, b int64) int64 {
	c, overflow := safeSub(a, b)
	if overflow {
		if b > 0 {
			return math.MinInt64
		}
		return math.MaxInt64
	}
	return c
}
