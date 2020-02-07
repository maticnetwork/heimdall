package listener

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/maticnetwork/heimdall/sethu/queue"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"
	"os"
)

const (
	listenerServiceStr = "listener-service"
)

// ListenerService starts and stops all chain event listeners
type ListenerService struct {
	// Base service
	common.BaseService
	queueConnector *queue.QueueConnector
	listeners      []Listener
}

// Global logger for bridge
var Logger log.Logger

func init() {
	Logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
}

// NewListenerService returns new service object for listneing to events
func NewListenerService(cdc *codec.Codec, queueConnector *queue.QueueConnector) *ListenerService {
	// create logger
	logger := Logger.With("module", listenerServiceStr)

	// creating listener object
	listenerService := &ListenerService{
		queueConnector: queueConnector,
	}

	listenerService.BaseService = *common.NewBaseService(logger, listenerServiceStr, listenerService)

	rootchainListener := &RootChainListener{}
	rootchainListener.BaseListener = *NewBaseListener(cdc, queueConnector, logger, "rootchain", rootchainListener)
	listenerService.listeners = append(listenerService.listeners, rootchainListener)

	return listenerService
}

// OnStart starts new block subscription
func (listenerService *ListenerService) OnStart() error {
	listenerService.BaseService.OnStart() // Always call the overridden method.
	listenerService.Logger.Info("Starting listeners", listenerService.listeners)

	// start chain listeners
	for _, listener := range listenerService.listeners {
		listener.Start()
	}

	return nil
}

// OnStop stops all necessary go routines
func (listenerService *ListenerService) OnStop() {
	listenerService.BaseService.OnStop() // Always call the overridden method.
	listenerService.Logger.Info("Listener Service Stopped")

}
