package processor

import (
	"github.com/cosmos/cosmos-sdk/client"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"

	"github.com/maticnetwork/heimdall/sethu/queue"
	"github.com/tendermint/tendermint/libs/log"
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

	// The "subclass" of BaseProcessor
	impl Processor

	// storage client
	// storageClient *leveldb.DB

}

// NewBaseProcessor creates a new BaseProcessor.
func NewBaseProcessor(cdc *codec.Codec, queueConnector *queue.QueueConnector, logger log.Logger, name string, impl Processor) *BaseProcessor {

	cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	cliCtx.BroadcastMode = client.BroadcastAsync
	cliCtx.TrustNode = true

	if logger == nil {
		logger = log.NewNopLogger()
	}

	// creating syncer object
	return &BaseProcessor{
		Logger: logger,
		name:   name,
		quit:   make(chan struct{}),
		impl:   impl,

		queueConnector: queueConnector,
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
