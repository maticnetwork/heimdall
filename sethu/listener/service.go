package listener

import (
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
}

// Global logger for bridge
var Logger log.Logger

func init() {
	Logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
}

// NewListenerService returns new service object for listneing to events
func NewListenerService(queueConnector *queue.QueueConnector) *ListenerService {
	// create logger
	logger := Logger.With("module", listenerServiceStr)

	// creating listener object
	listenerService := &ListenerService{
		queueConnector: queueConnector,
	}

	listenerService.BaseService = *common.NewBaseService(logger, listenerServiceStr, listenerService)
	return listenerService
}

// OnStart starts new block subscription
func (listenerService *ListenerService) OnStart() error {
	listenerService.BaseService.OnStart() // Always call the overridden method.
	listenerService.Logger.Info("Listener Service Started")

	listenerService.queueConnector.PublishMsg([]byte("TestMessage"), "test")

	return nil
}

// OnStop stops all necessary go routines
func (listenerService *ListenerService) OnStop() {
	listenerService.BaseService.OnStop() // Always call the overridden method.
	listenerService.Logger.Info("Listener Service Stopped")

}
