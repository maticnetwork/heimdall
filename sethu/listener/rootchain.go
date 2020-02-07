package listener

import (
	"context"

	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/sethu/queue"
)

// RootChainListener syncs validators and checkpoints
type RootChainListener struct {
	BaseListener
}

// Start starts new block subscription
func (rl *RootChainListener) Start() error {
	rl.Logger.Info("Starting listener", "name", rl.String())
	// create cancellable context
	ctx, cancelSubscription := context.WithCancel(context.Background())
	rl.cancelSubscription = cancelSubscription

	// create cancellable context
	headerCtx, cancelHeaderProcess := context.WithCancel(context.Background())
	rl.cancelHeaderProcess = cancelHeaderProcess

	// start header process
	go rl.StartHeaderProcess(headerCtx)

	// subscribe to new head
	subscription, err := rl.contractConnector.MainChainClient.SubscribeNewHead(ctx, rl.HeaderChannel)
	if err != nil {
		// start go routine to poll for new header using client object
		go rl.StartPolling(ctx, helper.GetConfig().SyncerPollInterval)
	} else {
		// start go routine to listen new header using subscription
		go rl.StartSubscription(ctx, subscription)
	}

	// subscribed to new head
	rl.Logger.Info("Subscribed to new head")

	return nil
}

func (bl *RootChainListener) ProcessHeader(newHeader *types.Header) {
	bl.Logger.Info("Received Headerblock from Rootchain", "header", newHeader)
	bl.queueConnector.PublishMsg([]byte("StakingMsg"), queue.StakingQueueRoute, bl.String())
}
