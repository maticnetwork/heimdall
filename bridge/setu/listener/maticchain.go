package listener

import (
	"context"
	"time"

	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/helper"
)

// MaticChainListener - Listens to and process headerblocks from maticchain
type MaticChainListener struct {
	BaseListener
}

// NewMaticChainListener - constructor func
func NewMaticChainListener() *MaticChainListener {
	return &MaticChainListener{}
}

// Start starts new block subscription
func (ml *MaticChainListener) Start() error {
	ml.Logger.Info("Starting")

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
		ml.Logger.Info("Start polling for header blocks", "pollInterval", helper.GetConfig().CheckpointerPollInterval)
		go ml.StartPolling(ctx, helper.GetConfig().CheckpointerPollInterval)
	} else {
		// start go routine to listen new header using subscription
		go ml.StartSubscription(ctx, subscription)
	}

	// subscribed to new head
	ml.Logger.Info("Subscribed to new head")

	return nil
}

// ProcessHeader - process headerblock from maticchain
func (ml *MaticChainListener) ProcessHeader(newHeader *types.Header) {
	ml.Logger.Info("Publishing received headerblock to checkpoint queue", "blockNumber", newHeader.Number)
	// Marshall header block and publish to queue
	headerBytes, err := newHeader.MarshalJSON()
	if err != nil {
		ml.Logger.Error("Error marshalling header block", "error", err)
	}
	ml.sendTask("sendCheckpointToHeimdall", headerBytes)
}

func (ml *MaticChainListener) sendTask(taskName string, headerBytes []byte) {
	// create machinery task
	signature := &tasks.Signature{
		Name: taskName,
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: string(headerBytes),
			},
		},
	}
	signature.RetryCount = 3
	// Delay the task by 5 seconds
	eta := time.Now().UTC().Add(time.Second * 5)
	signature.ETA = &eta
	// send task
	_, err := ml.queueConnector.Server.SendTask(signature)
	if err != nil {
		ml.Logger.Error("Error sending task", "taskName", taskName)
	}
}
