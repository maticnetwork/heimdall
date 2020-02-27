package processor

import (
	"os"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/maticnetwork/heimdall/bridge/setu/broadcaster"
	"github.com/maticnetwork/heimdall/bridge/setu/queue"
	"github.com/maticnetwork/heimdall/helper"
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

	contractCaller, err := helper.NewContractCaller()
	if err != nil {
		panic(err)
	}

	processorService.BaseService = *common.NewBaseService(logger, processorServiceStr, processorService)

	// initialize checkpoint processor
	checkpointProcessor := NewCheckpointProcessor()
	checkpointProcessor.BaseProcessor = *NewBaseProcessor(cdc, queueConnector, httpClient, txBroadcaster, &contractCaller.RootChainABI, "checkpoint", checkpointProcessor)
	processorService.processors = append(processorService.processors, checkpointProcessor)

	// initialize staking processor
	stakingProcessor := &StakingProcessor{}
	stakingProcessor.BaseProcessor = *NewBaseProcessor(cdc, queueConnector, httpClient, txBroadcaster, &contractCaller.RootChainABI, "staking", stakingProcessor)
	processorService.processors = append(processorService.processors, stakingProcessor)

	// initialize clerk processor
	clerkProcessor := &ClerkProcessor{}
	clerkProcessor.BaseProcessor = *NewBaseProcessor(cdc, queueConnector, httpClient, txBroadcaster, &contractCaller.RootChainABI, "clerk", clerkProcessor)
	processorService.processors = append(processorService.processors, clerkProcessor)

	// initialize fee processor
	feeProcessor := &FeeProcessor{}
	feeProcessor.BaseProcessor = *NewBaseProcessor(cdc, queueConnector, httpClient, txBroadcaster, &contractCaller.RootChainABI, "fee", feeProcessor)
	processorService.processors = append(processorService.processors, feeProcessor)

	// initialize span processor
	spanProcessor := &SpanProcessor{}
	spanProcessor.BaseProcessor = *NewBaseProcessor(cdc, queueConnector, httpClient, txBroadcaster, &contractCaller.RootChainABI, "span", spanProcessor)
	processorService.processors = append(processorService.processors, spanProcessor)

	return processorService
}

// OnStart starts new block subscription
func (processorService *ProcessorService) OnStart() error {
	processorService.BaseService.OnStart() // Always call the overridden method.

	// start processors
	for _, processor := range processorService.processors {
		go processor.Start()
	}

	processorService.Logger.Info("all processors Started")
	return nil
}

// OnStop stops all necessary go routines
func (processorService *ProcessorService) OnStop() {
	processorService.BaseService.OnStop() // Always call the overridden method.
	// start chain listeners
	for _, processor := range processorService.processors {
		processor.Stop()
	}

	processorService.Logger.Info("all processors stopped")
	return
}
