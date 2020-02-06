package listener

import "github.com/maticnetwork/bor/core/types"

// RootChainListener syncs validators and checkpoints
type RootChainListener struct {
	BaseListener
}

func (bl *RootChainListener) ProcessHeader(newHeader *types.Header) {
	bl.Logger.Info("Received Headerblock", "header", newHeader)
}
