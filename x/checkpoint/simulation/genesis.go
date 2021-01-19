package simulation

import (
	"time"

	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"
	"github.com/maticnetwork/heimdall/types/module"
	"github.com/maticnetwork/heimdall/x/checkpoint/types"
)

// RandomizedGenState return dummy genesis
func RandomizedGenState(simState *module.SimulationState) {
	lastNoACK := 0
	ackCount := 1
	startBlock := uint64(0)
	endBlock := uint64(256)
	rootHash := hmCommonTypes.HexToHeimdallHash("123")

	proposerAddress := hmCommonTypes.HexToHeimdallAddress("123")
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

	Checkpoints := make([]*hmTypes.Checkpoint, ackCount)

	for i := range Checkpoints {
		Checkpoints[i] = &bufferedCheckpoint
	}

	params := types.DefaultParams()
	genesisState := types.NewGenesisState(
		params,
		&bufferedCheckpoint,
		uint64(lastNoACK),
		uint64(ackCount),
		Checkpoints,
	)
	simState.GenState[types.ModuleName] = simState.Cdc.MustMarshalJSON(genesisState)

}
