//nolint:gosec
package simulation

import (
	"strconv"
	"time"

	"github.com/maticnetwork/heimdall/checkpoint/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/maticnetwork/heimdall/types/module"
)

// RandomizedGenState return dummy genesis
func RandomizedGenState(simState *module.SimulationState) {
	lastNoACK := 0
	ackCount := 1
	startBlock := uint64(0)
	endBlock := uint64(256)
	rootHash := hmTypes.HexToHeimdallHash("123")

	proposerAddress := hmTypes.HexToHeimdallAddress("123")
	timestamp := uint64(time.Now().Unix())
	borChainID := "1234"

	bufferedCheckpoint := hmTypes.CreateBlock(
		startBlock,
		endBlock,
		rootHash,
		proposerAddress,
		borChainID,
		timestamp,
	)

	checkpoints := make([]hmTypes.Checkpoint, ackCount)

	for i := range checkpoints {
		checkpoints[i] = bufferedCheckpoint
	}

	milestones := make([]hmTypes.Milestone, ackCount)

	for i := range milestones {
		milestones[i] = hmTypes.CreateMilestone(startBlock, endBlock, rootHash,
			proposerAddress, borChainID, strconv.Itoa(i), timestamp)
	}

	params := types.DefaultParams()
	genesisState := types.NewGenesisState(
		params,
		&bufferedCheckpoint,
		uint64(lastNoACK),
		uint64(ackCount),
		checkpoints,
		milestones,
	)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(genesisState)
}
