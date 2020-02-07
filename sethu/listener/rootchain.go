package listener

import "github.com/maticnetwork/bor/core/types"

import "github.com/maticnetwork/heimdall/sethu/queue"

// RootChainListener syncs validators and checkpoints
type RootChainListener struct {
	BaseListener
}

func (bl *RootChainListener) ProcessHeader(newHeader *types.Header) {
	bl.Logger.Info("Received Headerblock", "header", newHeader)
	bl.queueConnector.PublishMsg([]byte("Hello"), queue.StakingQueueRoute, bl.String())
}
