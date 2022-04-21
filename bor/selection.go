package bor

import (
	"encoding/binary"
	"math"
	"math/rand"

	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/heimdall/bor/types"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// XXXSelectNextProducers selects producers for next span by converting power to tickets
func XXXSelectNextProducers(blkHash common.Hash, spanEligibleVals []hmTypes.Validator, producerCount uint64) (selectedIDs []uint64, err error) {
	if len(spanEligibleVals) <= int(producerCount) {
		for _, val := range spanEligibleVals {
			selectedIDs = append(selectedIDs, uint64(val.ID))
		}
		return
	}

	// extract seed from hash
	seed := helper.ToBytes32(blkHash.Bytes()[:32])
	validatorIndices := convertToSlots(spanEligibleVals)
	selectedIDs, err = ShuffleList(validatorIndices, seed)
	if err != nil {
		return
	}
	return selectedIDs[:producerCount], nil
}

// converts validator power to slots
// TODO remove 2nd loop
func convertToSlots(vals []hmTypes.Validator) (validatorIndices []uint64) {
	for _, val := range vals {
		for val.VotingPower >= types.SlotCost {
			validatorIndices = append(validatorIndices, uint64(val.ID))
			val.VotingPower = val.VotingPower - types.SlotCost
		}
	}
	return validatorIndices
}

//
// New selection algorithm
//

// SelectNextProducers selects producers for next span by converting power to tickets
func SelectNextProducers(blkHash common.Hash, spanEligibleValidators []hmTypes.Validator, producerCount uint64) ([]uint64, error) {
	selectedProducers := make([]uint64, 0)

	if len(spanEligibleValidators) <= int(producerCount) {
		for _, validator := range spanEligibleValidators {
			selectedProducers = append(selectedProducers, uint64(validator.ID))
		}

		return selectedProducers, nil
	}

	// extract seed from hash
	seedBytes := helper.ToBytes32(blkHash.Bytes()[:32])
	seed := int64(binary.BigEndian.Uint64(seedBytes[:]))
	rand.Seed(seed)

	// weighted range from validators' voting power
	votingPower := make([]uint64, len(spanEligibleValidators))
	for idx, validator := range spanEligibleValidators {
		votingPower[idx] = uint64(validator.VotingPower)
	}

	weightedRanges, totalVotingPower := createWeightedRanges(votingPower)
	// select producers, with replacement
	for i := uint64(0); i < producerCount; i++ {
		/*
			random must be in [1, totalVotingPower] to avoid situation such as
			2 validators with 1 staking power each.
			Weighted range will look like (1, 2)
			Rolling inclusive will have a range of 0 - 2, making validator with staking power 1 chance of selection = 66%
		*/
		targetWeight := randomRangeInclusive(1, totalVotingPower)
		index := binarySearch(weightedRanges, targetWeight)
		selectedProducers = append(selectedProducers, spanEligibleValidators[index].ID.Uint64())
	}

	return selectedProducers[:producerCount], nil
}

func binarySearch(array []uint64, search uint64) int {
	if len(array) == 0 {
		return -1
	}
	l := 0
	r := len(array) - 1
	for l < r {
		mid := (l + r) / 2
		if array[mid] >= search {
			r = mid
		} else {
			l = mid + 1
		}
	}
	return l
}

// randomRangeInclusive produces unbiased pseudo random in the range [min, max]. Uses rand.Uint64() and can be seeded beforehand.
func randomRangeInclusive(min uint64, max uint64) uint64 {
	if max <= min {
		return max
	}

	rangeLength := max - min + 1
	maxAllowedValue := math.MaxUint64 - math.MaxUint64%rangeLength - 1
	randomValue := rand.Uint64()

	// reject anything that is beyond the reminder to avoid bias
	for randomValue >= maxAllowedValue {
		randomValue = rand.Uint64()
	}

	return min + randomValue%rangeLength
}

// createWeightedRanges converts array [1, 2, 3] into cumulative form [1, 3, 6]
func createWeightedRanges(weights []uint64) ([]uint64, uint64) {
	weightedRanges := make([]uint64, len(weights))
	totalWeight := uint64(0)
	for i := 0; i < len(weightedRanges); i++ {
		totalWeight += weights[i]
		weightedRanges[i] = totalWeight
	}
	return weightedRanges, totalWeight
}
