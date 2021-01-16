package processor

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/maticnetwork/heimdall/bridge/setu/broadcaster"
	"github.com/maticnetwork/heimdall/bridge/setu/queue"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/tendermint/tendermint/libs/log"
	httpClient "github.com/tendermint/tendermint/rpc/client/http"
)

// Processor defines a block header listerner for Rootchain, Maticchain, Heimdall
type Processor interface {
	Start() error

	RegisterTasks()

	String() string

	Stop()
}

type BaseProcessor struct {
	Logger log.Logger
	name   string
	quit   chan struct{}

	// queue connector
	queueConnector *queue.QueueConnector

	// tx broadcaster
	txBroadcaster *broadcaster.TxBroadcaster

	// params context
	paramsContext *util.ParamsContext

	// The "subclass" of BaseProcessor
	impl Processor

	// cli context
	cliCtx client.Context

	// contract caller
	contractConnector helper.ContractCaller

	// http client to subscribe to
	httpClient *httpClient.HTTP

	// storage client
	storageClient *leveldb.DB
}

// NewBaseProcessor creates a new BaseProcessor.
func NewBaseProcessor(cdc codec.Marshaler, queueConnector *queue.QueueConnector, httpClient *httpClient.HTTP, txBroadcaster *broadcaster.TxBroadcaster, paramsContext *util.ParamsContext, name string, impl Processor) *BaseProcessor {
	logger := util.Logger().With("service", "processor", "module", name)

	cliCtx := client.Context{}.WithJSONMarshaler(cdc)
	cliCtx.BroadcastMode = flags.BroadcastAsync

	contractCaller, err := helper.NewContractCaller()
	if err != nil {
		logger.Error("Error while getting root chain instance", "error", err)
		panic(err)
	}

	if logger == nil {
		logger = log.NewNopLogger()
	}

	// creating syncer object
	return &BaseProcessor{
		Logger: logger,
		name:   name,
		quit:   make(chan struct{}),
		impl:   impl,

		cliCtx:            cliCtx,
		queueConnector:    queueConnector,
		contractConnector: contractCaller,
		txBroadcaster:     txBroadcaster,
		paramsContext:     paramsContext,
		httpClient:        httpClient,
		storageClient:     util.GetBridgeDBInstance(viper.GetString(util.BridgeDBFlag)),
	}
}

// String implements Service by returning a string representation of the service.
func (bp *BaseProcessor) String() string {
	return bp.name
}

// OnStop stops all necessary go routines
func (bp *BaseProcessor) Stop() {
	// override to stop any go-routines in individual processors
}
