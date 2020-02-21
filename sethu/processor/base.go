package processor

import (
	"github.com/cosmos/cosmos-sdk/client"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/sethu/broadcaster"
	"github.com/maticnetwork/heimdall/sethu/queue"
	"github.com/tendermint/tendermint/libs/log"
	httpClient "github.com/tendermint/tendermint/rpc/client"
)

// Processor defines a block header listerner for Rootchain, Maticchain, Heimdall
type Processor interface {
	Start() error

	String() string

	SetLogger(log.Logger)
}

type BaseProcessor struct {
	Logger  log.Logger
	name    string
	started uint32 // atomic
	stopped uint32 // atomic
	quit    chan struct{}

	// queue connector
	queueConnector *queue.QueueConnector

	// tx broadcaster
	txBroadcaster *broadcaster.TxBroadcaster

	// The "subclass" of BaseProcessor
	impl Processor

	// cli context
	cliCtx cliContext.CLIContext

	// contract caller
	contractConnector helper.ContractCaller

	// http client to subscribe to
	httpClient *httpClient.HTTP

	// storage client
	// storageClient *leveldb.DB

}

// NewBaseProcessor creates a new BaseProcessor.
func NewBaseProcessor(cdc *codec.Codec, queueConnector *queue.QueueConnector, httpClient *httpClient.HTTP, txBroadcaster *broadcaster.TxBroadcaster, logger log.Logger, name string, impl Processor) *BaseProcessor {

	cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	cliCtx.BroadcastMode = client.BroadcastAsync
	cliCtx.TrustNode = true

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
		httpClient:        httpClient,
		// storageClient: getBridgeDBInstance(viper.GetString(BridgeDBFlag)),
	}
}

// SetLogger implements Service by setting a logger.
func (bp *BaseProcessor) SetLogger(l log.Logger) {
	bp.Logger = l
}

// String implements Service by returning a string representation of the service.
func (bp *BaseProcessor) String() string {
	return bp.name
}
