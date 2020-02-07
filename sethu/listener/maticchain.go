package listener

import (
	"context"

	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/sethu/queue"
)

// MaticChainListener syncs validators and checkpoints
type MaticChainListener struct {
	BaseListener
}

// Start starts new block subscription
func (ml *MaticChainListener) Start() error {
	ml.Logger.Info("Starting listener", "name", ml.String())
	// create cancellable context
	ctx, cancelSubscription := context.WithCancel(context.Background())
	ml.cancelSubscription = cancelSubscription

	// create cancellable context
	headerCtx, cancelHeaderProcess := context.WithCancel(context.Background())
	ml.cancelHeaderProcess = cancelHeaderProcess

	// start header process
	go ml.StartHeaderProcess(headerCtx)

	// subscribe to new head
	subscription, err := ml.contractConnector.MaticChainClient.SubscribeNewHead(ctx, ml.HeaderChannel)
	if err != nil {
		// start go routine to poll for new header using client object
		go ml.StartPolling(ctx, helper.GetConfig().SyncerPollInterval)
	} else {
		// start go routine to listen new header using subscription
		go ml.StartSubscription(ctx, subscription)
	}

	// subscribed to new head
	ml.Logger.Info("Subscribed to new head")

	return nil
}

func (ml *MaticChainListener) ProcessHeader(newHeader *types.Header) {
	ml.Logger.Info("Received Headerblock from maticchain", "header", newHeader)
	ml.queueConnector.PublishMsg([]byte("Hello Matic"), queue.StakingQueueRoute, ml.String())
}
