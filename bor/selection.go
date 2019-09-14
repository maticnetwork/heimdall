package bor

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

// SelectNextProducers selects producers for next span by converting power to tickets
func SelectNextProducers(logger tmlog.Logger, blkHash common.Hash, currentVals []types.Validator) (selectedIDs []uint64, err error) {
	// extract seed from hash
	seed := helper.ToBytes32(blkHash.Bytes()[:32])
	logger.Info("Seed generated", "Seed", hex.EncodeToString(seed[:]), "BlkHash", blkHash.String())
	validatorIndices := convertToSlots(currentVals)
	logger.Info("Created validator indices", "Length", len(validatorIndices), "ValIndices", validatorIndices)
	return ShuffleList(validatorIndices, seed)
}

func convertToSlots(vals []types.Validator) (validatorIndices []uint64) {
	for _, val := range vals {
		for val.Power > SlotCost {
			fmt.Printf("Creating slots for validator %v with power %v", val.ID, val.Power)
			validatorIndices = append(validatorIndices, uint64(val.ID))
			val.Power = val.Power - SlotCost
		}
	}
	return validatorIndices
}
