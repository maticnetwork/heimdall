package processor

import "github.com/maticnetwork/heimdall/sethu/queue"

// CheckpointProcessor
type CheckpointProcessor struct {
	BaseProcessor
}

// Start starts new block subscription
func (cp *CheckpointProcessor) Start() error {
	cp.Logger.Info("Starting Processor", "name", cp.String())

	amqpMsgs, _ := cp.queueConnector.ConsumeMsg(queue.CheckpointQueueName)
	// handle all amqp messages
	for amqpMsg := range amqpMsgs {
		cp.Logger.Info("Received Message from checkpoint queue", "Msg - ", string(amqpMsg.Body), "AppID", amqpMsg.AppId)

		// send ack
		amqpMsg.Ack(false)

	}
	return nil
}
