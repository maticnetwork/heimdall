package bor

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

// SelectNextProducers selects producers for next span by converting power to tickets
func SelectNextProducers(logger tmlog.Logger, blkHash common.Hash, currentVals []types.Validator, producerCount uint64) (selectedIDs []uint64, err error) {
	// extract seed from hash
	seed := helper.ToBytes32(blkHash.Bytes()[:32])
	validatorIndices := convertToSlots(currentVals)
	selectedIDs, err = ShuffleList(validatorIndices, seed)
	if err != nil {
		return
	}
	if len(selectedIDs) < int(producerCount) {
		return selectedIDs, nil
	}
	return selectedIDs[:producerCount], nil
}

// converts validator power to slots
// TODO remove 2nd loop
func convertToSlots(vals []types.Validator) (validatorIndices []uint64) {
	for _, val := range vals {
		for val.Power >= SlotCost {
			validatorIndices = append(validatorIndices, uint64(val.ID))
			val.Power = val.Power - SlotCost
		}
	}
	return validatorIndices
}
