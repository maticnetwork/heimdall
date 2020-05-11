package bor

import (
	"encoding/binary"
	"math/rand"

	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

func binarySearch(array []uint64, search uint64) int {
	if len(array) == 0 {
		return -1
	}
	l := 0
	r := len(array)-1
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

// SelectNextProducers selects producers for next span by converting power to tickets
func SelectNextProducers(blkHash common.Hash, spanEligibleVals []hmTypes.Validator, producerCount uint64) (selectedIDs []uint64, err error) {
	if len(spanEligibleVals) <= int(producerCount) {
		for _, val := range spanEligibleVals {
			selectedIDs = append(selectedIDs, uint64(val.ID))
		}
		return
	}

	// extract seed from hash
	seedBytes := helper.ToBytes32(blkHash.Bytes()[:32])
	seed := int64(binary.BigEndian.Uint64(seedBytes[:]))
	rand.Seed(seed)
	ranges := convertToRanges(spanEligibleVals)
	for i := uint64(0); i < producerCount; i++ {
		x := rand.Uint64()
		index := binarySearch(ranges, x)
		selectedIDs = append(selectedIDs, spanEligibleVals[index].ID.Uint64())
	}
	return selectedIDs, nil
}

// makes ranges out of validators powers
func convertToRanges(vals []hmTypes.Validator) []uint64 {
	ranges := make([]uint64, len(vals))
	for i := 0; i < len(ranges); i++ {
		prev := uint64(0)
		if i > 0 {
			prev += ranges[i-1]
		}
		ranges[i] = uint64(vals[i].VotingPower) + prev
	}
	return ranges
}
