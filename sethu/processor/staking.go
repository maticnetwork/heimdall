package processor

import "github.com/maticnetwork/heimdall/sethu/util"

// RootChainListener syncs validators and checkpoints
type StakingProcessor struct {
	BaseProcessor
}

// Start starts new block subscription
func (sp *StakingProcessor) Start() error {
	sp.Logger.Info("Starting Processor", "name", sp.String())

	amqpMsgs, _ := sp.queueConnector.ConsumeMsg(util.StakingQueueName)
	// handle all amqp messages
	for amqpMsg := range amqpMsgs {
		sp.Logger.Info("Received Message from staking queue", "Msg - ", string(amqpMsg.Body), "AppID", amqpMsg.AppId)

		// send ack
		amqpMsg.Ack(false)

	}
	return nil
}
