package listener

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/maticnetwork/heimdall/bridge/setu/queue"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/tendermint/tendermint/libs/service"
	httpClient "github.com/tendermint/tendermint/rpc/client/http"
)

const (
	ListenerServiceStr = "listener"

	RootChainListenerStr  = "rootchain"
	HeimdallListenerStr   = "heimdall"
	MaticChainListenerStr = "maticchain"
)

// ListenerService starts and stops all chain event listeners
type ListenerService struct {
	// Base service
	service.BaseService
	listeners []Listener
}

// NewListenerService returns new service object for listening to events
func NewListenerService(cliCtx client.Context, queueConnector *queue.QueueConnector, httpClient *httpClient.HTTP) *ListenerService {

	var logger = util.Logger().With("service", ListenerServiceStr)
	// creating listener object
	listenerService := &ListenerService{}

	listenerService.BaseService = *service.NewBaseService(logger, ListenerServiceStr, listenerService)

	rootchainListener := NewRootChainListener()
	rootchainListener.BaseListener = *NewBaseListener(cliCtx, queueConnector, httpClient, helper.GetMainClient(), RootChainListenerStr, rootchainListener)
	listenerService.listeners = append(listenerService.listeners, rootchainListener)

	maticchainListener := NewMaticChainListener()
	maticchainListener.BaseListener = *NewBaseListener(cliCtx, queueConnector, httpClient, helper.GetMaticClient(), MaticChainListenerStr, maticchainListener)
	listenerService.listeners = append(listenerService.listeners, maticchainListener)

	heimdallListener := NewHeimdallListener()
	heimdallListener.BaseListener = *NewBaseListener(cliCtx, queueConnector, httpClient, nil, HeimdallListenerStr, heimdallListener)
	listenerService.listeners = append(listenerService.listeners, heimdallListener)

	return listenerService
}

// OnStart starts new block subscription
func (listenerService *ListenerService) OnStart() error {
	if err := listenerService.BaseService.OnStart(); err != nil {
		listenerService.Logger.Error("OnStart | OnStart", "Error", err)
	} // Always call the overridden method.

	// start chain listeners
	for _, listener := range listenerService.listeners {
		if err := listener.Start(); err != nil {
			listenerService.Logger.Error("OnStart | Start", "Error", err)
		}
	}

	listenerService.Logger.Info("all listeners Started")
	return nil
}

// OnStop stops all necessary go routines
func (listenerService *ListenerService) OnStop() {
	listenerService.BaseService.OnStop() // Always call the overridden method.

	// start chain listeners
	for _, listener := range listenerService.listeners {
		listener.Stop()
	}

	listenerService.Logger.Info("all listeners stopped")
}
