package processor

import (
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/maticnetwork/heimdall/sethu/broadcaster"
	"github.com/maticnetwork/heimdall/sethu/queue"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"

	httpClient "github.com/tendermint/tendermint/rpc/client"
)

// ProcessorService starts and stops all event processors
type ProcessorService struct {
	// Base service
	common.BaseService

	// queue connector
	queueConnector *queue.QueueConnector

	// tx broadcaster
	txBroadcsater *broadcaster.TxBroadcaster

	processors []Processor
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
func NewProcessorService(cdc *codec.Codec, queueConnector *queue.QueueConnector, httpClient *httpClient.HTTP, txBroadcaster *broadcaster.TxBroadcaster) *ProcessorService {
	// create logger
	logger := Logger.With("module", processorServiceStr)

	// creating processor object
	processorService := &ProcessorService{
		queueConnector: queueConnector,
	}

	processorService.BaseService = *common.NewBaseService(logger, processorServiceStr, processorService)

	checkpointProcessor := NewCheckpointProcessor()
	checkpointProcessor.BaseProcessor = *NewBaseProcessor(cdc, queueConnector, httpClient, txBroadcaster, logger, "checkpoint", checkpointProcessor)
	processorService.processors = append(processorService.processors, checkpointProcessor)

	stakingProcessor := &StakingProcessor{}
	stakingProcessor.BaseProcessor = *NewBaseProcessor(cdc, queueConnector, httpClient, txBroadcaster, logger, "staking", stakingProcessor)
	processorService.processors = append(processorService.processors, stakingProcessor)

	return processorService
}

// OnStart starts new block subscription
func (processorService *ProcessorService) OnStart() error {
	processorService.BaseService.OnStart() // Always call the overridden method.
	processorService.Logger.Info("Processor Service Started")

	// start processors
	for _, processor := range processorService.processors {
		go processor.Start()
	}

	return nil
}

// OnStop stops all necessary go routines
func (processorService *ProcessorService) OnStop() {
	processorService.BaseService.OnStop() // Always call the overridden method.
	processorService.Logger.Info("Processor Service Stopped")

}
