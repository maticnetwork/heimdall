package broadcaster

import (
	"fmt"
	"os"
	"testing"

	"github.com/maticnetwork/heimdall/app"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// Parallel test - to check BroadcastToHeimdall syncronization
func TestBroadcastToHeimdall(t *testing.T) {
	t.Parallel()
	cdc := app.MakeCodec()
	// cli context
	tendermintNode := "tcp://localhost:26657"
	viper.Set(helper.NodeFlag, tendermintNode)
	viper.Set("log_level", "info")
	// cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	// cliCtx.BroadcastMode = client.BroadcastSync
	// cliCtx.TrustNode = true

	helper.InitHeimdallConfig(os.ExpandEnv("$HOME/.heimdalld"))
	_txBroadcaster := NewTxBroadcaster(cdc)

	testData := []checkpointTypes.MsgCheckpoint{
		{Proposer: hmTypes.BytesToHeimdallAddress(helper.GetAddress()), StartBlock: 0, EndBlock: 63, RootHash: hmTypes.HexToHeimdallHash("0x5bd83f679c8ce7c48d6fa52ce41532fcacfbbd99d5dab415585f397bf44a0b6e"), AccountRootHash: hmTypes.HexToHeimdallHash("0xd10b5c16c25efe0b0f5b3d75038834223934ae8c2ec2b63a62bbe42aa21e2d2d")},
		{Proposer: hmTypes.BytesToHeimdallAddress(helper.GetAddress()), StartBlock: 64, EndBlock: 1024, RootHash: hmTypes.HexToHeimdallHash("0x5bd83f679c8ce7c48d6fa52ce41532fcacfbbd99d5dab415585f397bf44a0b6e"), AccountRootHash: hmTypes.HexToHeimdallHash("0xd10b5c16c25efe0b0f5b3d75038834223934ae8c2ec2b63a62bbe42aa21e2d2d")},
		{Proposer: hmTypes.BytesToHeimdallAddress(helper.GetAddress()), StartBlock: 1025, EndBlock: 2048, RootHash: hmTypes.HexToHeimdallHash("0x5bd83f679c8ce7c48d6fa52ce41532fcacfbbd99d5dab415585f397bf44a0b6e"), AccountRootHash: hmTypes.HexToHeimdallHash("0xd10b5c16c25efe0b0f5b3d75038834223934ae8c2ec2b63a62bbe42aa21e2d2d")},
		{Proposer: hmTypes.BytesToHeimdallAddress(helper.GetAddress()), StartBlock: 2049, EndBlock: 3124, RootHash: hmTypes.HexToHeimdallHash("0x5bd83f679c8ce7c48d6fa52ce41532fcacfbbd99d5dab415585f397bf44a0b6e"), AccountRootHash: hmTypes.HexToHeimdallHash("0xd10b5c16c25efe0b0f5b3d75038834223934ae8c2ec2b63a62bbe42aa21e2d2d")},
	}

	for index, test := range testData {
		t.Run(fmt.Sprint(index), func(t *testing.T) {
			// create and send checkpoint message
			msg := checkpointTypes.NewMsgCheckpointBlock(
				test.Proposer,
				test.StartBlock,
				test.EndBlock,
				test.RootHash,
				test.AccountRootHash,
				"1234",
			)

			err := _txBroadcaster.BroadcastToHeimdall(msg, nil)
			assert.Empty(t, err, "Error broadcasting tx to heimdall", err)
		})
	}
}
