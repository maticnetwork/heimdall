package bor

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

// SelectNextProducers selects producers for next span by converting power to tickets
func SelectNextProducers(logger tmlog.Logger, blkHash common.Hash, currentVals []types.Validator) {
	// extract seed from hash
	seed := helper.ToBytes32(blkHash.Bytes()[:32])
	logger.Info("Seed generated", "Seed", hex.EncodeToString(seed[:]), "BlkHash", blkHash.String())
	validatorIndices := convertToSlots(currentVals)
	// ShuffleList(input, seed)

	// fetch eth block for height newSpanID
}

func convertToSlots(vals []types.Validator) (validatorIndices []uint64) {
	for i, val := range vals {
		for val.Power > SlotCost {
			validatorIndices := validatorIndices.append()
		}
	}
}
