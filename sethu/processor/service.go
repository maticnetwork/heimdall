package processor

import (
	"github.com/maticnetwork/heimdall/sethu/queue"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"os"
)

// ProcessorService starts and stops all event processors
type ProcessorService struct {
	// Base service
	common.BaseService

	// queue connector
	queueConnector *queue.QueueConnector
}

const (
	processorServiceStr = "processor-service"
)

// Global logger for bridge
var Logger log.Logger

func init() {
	Logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
}

// NewProcessorService returns new service object for processing queue msg
func NewProcessorService(queueConnector *queue.QueueConnector) *ProcessorService {
	// create logger
	logger := Logger.With("module", processorServiceStr)

	// creating processor object
	processorService := &ProcessorService{
		queueConnector: queueConnector,
	}

	processorService.BaseService = *common.NewBaseService(logger, processorServiceStr, processorService)
	return processorService
}

// OnStart starts new block subscription
func (processorService *ProcessorService) OnStart() error {
	processorService.BaseService.OnStart() // Always call the overridden method.
	processorService.Logger.Info("Processor Service Started")

	amqpMsgs, _ := processorService.queueConnector.ConsumeMsg("test-queue")
	// handle all amqp messages
	for amqpMsg := range amqpMsgs {
		processorService.Logger.Info("Received Message", "Msg - ", string(amqpMsg.Body))
		// send ack
		amqpMsg.Ack(false)

	}
	return nil
}

// OnStop stops all necessary go routines
func (processorService *ProcessorService) OnStop() {
	processorService.BaseService.OnStop() // Always call the overridden method.
	processorService.Logger.Info("Processor Service Stopped")

}
