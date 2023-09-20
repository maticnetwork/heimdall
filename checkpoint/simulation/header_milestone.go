package simulation

import (
	"time"

	"github.com/maticnetwork/heimdall/types"
)

// GenRandMilestone return headers
func GenRandMilestone(start uint64, sprintLength uint64) (milestone types.Milestone, err error) {
	end := start + sprintLength - 1
	borChainID := "1234"
	rootHash := types.HexToHeimdallHash("123")
	proposer := types.HeimdallAddress{}
	milestoneID := "00000"
	milestone = types.CreateMilestone(
		start,
		end,
		rootHash,
		proposer,
		borChainID,
		milestoneID,
		uint64(time.Now().UTC().Unix()))

	return milestone, nil
}
