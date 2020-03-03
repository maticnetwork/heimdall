package processor

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/common"
	httpClient "github.com/tendermint/tendermint/rpc/client"

	"github.com/maticnetwork/heimdall/bridge/setu/broadcaster"
	"github.com/maticnetwork/heimdall/bridge/setu/queue"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	"github.com/maticnetwork/heimdall/helper"
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

// NewProcessorService returns new service object for processing queue msg
func NewProcessorService(
	cdc *codec.Codec,
	queueConnector *queue.QueueConnector,
	httpClient *httpClient.HTTP,
	txBroadcaster *broadcaster.TxBroadcaster,
) *ProcessorService {
	// create logger
	logger := util.Logger().With("module", processorServiceStr)

	// creating processor object
	processorService := &ProcessorService{
		queueConnector: queueConnector,
	}

	contractCaller, err := helper.NewContractCaller()
	if err != nil {
		panic(err)
	}

	processorService.BaseService = *common.NewBaseService(logger, processorServiceStr, processorService)

	//
	// Intitialize processors
	//

	// initialize checkpoint processor
	checkpointProcessor := NewCheckpointProcessor()
	checkpointProcessor.BaseProcessor = *NewBaseProcessor(cdc, queueConnector, httpClient, txBroadcaster, &contractCaller.RootChainABI, "checkpoint", checkpointProcessor)

	// initialize staking processor
	stakingProcessor := &StakingProcessor{}
	stakingProcessor.BaseProcessor = *NewBaseProcessor(cdc, queueConnector, httpClient, txBroadcaster, &contractCaller.RootChainABI, "staking", stakingProcessor)

	// initialize clerk processor
	clerkProcessor := &ClerkProcessor{}
	clerkProcessor.BaseProcessor = *NewBaseProcessor(cdc, queueConnector, httpClient, txBroadcaster, &contractCaller.RootChainABI, "clerk", clerkProcessor)

	// initialize fee processor
	feeProcessor := &FeeProcessor{}
	feeProcessor.BaseProcessor = *NewBaseProcessor(cdc, queueConnector, httpClient, txBroadcaster, &contractCaller.RootChainABI, "fee", feeProcessor)

	// initialize span processor
	spanProcessor := &SpanProcessor{}
	spanProcessor.BaseProcessor = *NewBaseProcessor(cdc, queueConnector, httpClient, txBroadcaster, &contractCaller.RootChainABI, "span", spanProcessor)

	//
	// Select processors
	//

	// add into processor list
	startAll := viper.GetBool("all")
	onlyServices := viper.GetStringSlice("only")

	if startAll {
		processorService.processors = append(processorService.processors,
			checkpointProcessor,
			stakingProcessor,
			clerkProcessor,
			feeProcessor,
			spanProcessor,
		)
	} else {
		for _, service := range onlyServices {
			switch service {
			case "checkpoint":
				processorService.processors = append(processorService.processors, checkpointProcessor)
			case "staking":
				processorService.processors = append(processorService.processors, stakingProcessor)
			case "clerk":
				processorService.processors = append(processorService.processors, clerkProcessor)
			case "fee":
				processorService.processors = append(processorService.processors, feeProcessor)
			case "span":
				processorService.processors = append(processorService.processors, spanProcessor)
			}
		}
	}

	if len(processorService.processors) == 0 {
		panic("No processors selected. Use --all or --only <coma-seprated processors>")
	}

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
