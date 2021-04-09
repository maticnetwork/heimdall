package processor

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/service"
	httpClient "github.com/tendermint/tendermint/rpc/client/http"

	"github.com/maticnetwork/heimdall/bridge/setu/broadcaster"
	"github.com/maticnetwork/heimdall/bridge/setu/queue"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	"github.com/maticnetwork/heimdall/helper"
)

const (
	processorServiceStr = "processor-service"
)

// ProcessorService starts and stops all event processors
type ProcessorService struct {
	// Base service
	service.BaseService

	// queue connector
	queueConnector *queue.QueueConnector

	processors []Processor
}

// NewProcessorService returns new service object for processing queue msg
func NewProcessorService(
	cliCtx client.Context,
	queueConnector *queue.QueueConnector,
	httpClient *httpClient.HTTP,
	txBroadcaster *broadcaster.TxBroadcaster,
	paramsContext *util.ParamsContext,
) *ProcessorService {
	var logger = util.Logger().With("module", processorServiceStr)
	// creating processor object
	processorService := &ProcessorService{
		queueConnector: queueConnector,
	}

	contractCaller, err := helper.NewContractCaller()
	if err != nil {
		panic(err)
	}

	processorService.BaseService = *service.NewBaseService(logger, processorServiceStr, processorService)

	//
	// Initialize processors
	//

	// initialize checkpoint processor
	checkpointProcessor := NewCheckpointProcessor(&contractCaller.RootChainABI)
	checkpointProcessor.BaseProcessor = *NewBaseProcessor(cliCtx, queueConnector, httpClient, txBroadcaster, paramsContext, "checkpoint", checkpointProcessor)

	// initialize fee processor
	feeProcessor := NewFeeProcessor(&contractCaller.StakingInfoABI)
	feeProcessor.BaseProcessor = *NewBaseProcessor(cliCtx, queueConnector, httpClient, txBroadcaster, paramsContext, "fee", feeProcessor)

	// initialize staking processor
	stakingProcessor := NewStakingProcessor(&contractCaller.StakingInfoABI)
	stakingProcessor.BaseProcessor = *NewBaseProcessor(cliCtx, queueConnector, httpClient, txBroadcaster, paramsContext, "staking", stakingProcessor)

	// initialize clerk processor
	clerkProcessor := NewClerkProcessor(&contractCaller.StateSenderABI)
	clerkProcessor.BaseProcessor = *NewBaseProcessor(cliCtx, queueConnector, httpClient, txBroadcaster, paramsContext, "clerk", clerkProcessor)

	// initialize span processor
	spanProcessor := &SpanProcessor{}
	spanProcessor.BaseProcessor = *NewBaseProcessor(cliCtx, queueConnector, httpClient, txBroadcaster, paramsContext, "span", spanProcessor)

	// initialize slashing processor
	//slashingProcessor := NewSlashingProcessor(&contractCaller.StakingInfoABI)
	//slashingProcessor.BaseProcessor = *NewBaseProcessor(cliCtx, queueConnector, httpClient, txBroadcaster, paramsContext, "slashing", slashingProcessor)

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
		for _, serviceName := range onlyServices {
			switch serviceName {
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
				//case "slashing":
				//	processorService.processors = append(processorService.processors, nil)
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
	if err := processorService.BaseService.OnStart(); err != nil {
		processorService.Logger.Error("OnStart | OnStart", "Error", err)
	} // Always call the overridden method.

	// start processors
	for _, processor := range processorService.processors {
		processor.RegisterTasks()
		processor := processor
		go func() {
			if err := processor.Start(); err != nil {
				processorService.Logger.Error("processor is failed", "Err", err)
			}
		}()
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
}
