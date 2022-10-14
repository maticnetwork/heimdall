package listener

import (
	"context"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tendermint/libs/log"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/maticnetwork/heimdall/bridge/setu/queue"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	"github.com/maticnetwork/heimdall/helper"

	httpClient "github.com/tendermint/tendermint/rpc/client"
)

// Listener defines a block header listerner for Rootchain, Maticchain, Heimdall
type Listener interface {
	Start() error

	StartHeaderProcess(context.Context)

	StartPolling(context.Context, time.Duration)

	StartSubscription(context.Context, ethereum.Subscription)

	ProcessHeader(*types.Header)

	Stop()

	String() string
}

type BaseListener struct {
	Logger log.Logger
	name   string
	quit   chan struct{}

	// The "subclass" of BaseService
	impl Listener

	// contract caller
	contractConnector helper.ContractCaller

	chainClient *ethclient.Client

	// header channel
	HeaderChannel chan *types.Header

	// cancel function for poll/subscription
	cancelSubscription context.CancelFunc

	// header listener subscription
	cancelHeaderProcess context.CancelFunc

	// cli context
	cliCtx cliContext.CLIContext

	// queue connector
	queueConnector *queue.QueueConnector

	// http client to subscribe to
	httpClient *httpClient.HTTP

	// storage client
	storageClient *leveldb.DB
}

// NewBaseListener creates a new BaseListener.
func NewBaseListener(cdc *codec.Codec, queueConnector *queue.QueueConnector, httpClient *httpClient.HTTP, chainClient *ethclient.Client, name string, impl Listener) *BaseListener {

	logger := util.Logger().With("service", "listener", "module", name)
	contractCaller, err := helper.NewContractCaller()
	if err != nil {
		logger.Error("Error while getting root chain instance", "error", err)
		panic(err)
	}

	cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	cliCtx.BroadcastMode = client.BroadcastAsync
	cliCtx.TrustNode = true

	// creating syncer object
	return &BaseListener{
		Logger:        logger,
		name:          name,
		quit:          make(chan struct{}),
		impl:          impl,
		storageClient: util.GetBridgeDBInstance(viper.GetString(util.BridgeDBFlag)),

		cliCtx:            cliCtx,
		queueConnector:    queueConnector,
		httpClient:        httpClient,
		contractConnector: contractCaller,
		chainClient:       chainClient,

		HeaderChannel: make(chan *types.Header),
	}
}

// // Start starts new block subscription
// func (bl *BaseListener) Start() error {
// 	bl.Logger.Info("Starting listener", "name", bl.String())
// 	// create cancellable context
// 	ctx, cancelSubscription := context.WithCancel(context.Background())
// 	bl.cancelSubscription = cancelSubscription

// 	// create cancellable context
// 	headerCtx, cancelHeaderProcess := context.WithCancel(context.Background())
// 	bl.cancelHeaderProcess = cancelHeaderProcess

// 	// start header process
// 	go bl.StartHeaderProcess(headerCtx)

// 	// subscribe to new head
// 	subscription, err := bl.contractConnector.MainChainClient.SubscribeNewHead(ctx, bl.HeaderChannel)
// 	if err != nil {
// 		// start go routine to poll for new header using client object
// 		go bl.StartPolling(ctx, helper.GetConfig().SyncerPollInterval)
// 	} else {
// 		// start go routine to listen new header using subscription
// 		go bl.StartSubscription(ctx, subscription)
// 	}

// 	// subscribed to new head
// 	bl.Logger.Info("Subscribed to new head")

// 	return nil
// }

// String implements Service by returning a string representation of the service.
func (bl *BaseListener) String() string {
	return bl.name
}

// startHeaderProcess starts header process when they get new header
func (bl *BaseListener) StartHeaderProcess(ctx context.Context) {
	bl.Logger.Info("Starting header process")
	for {
		select {
		case newHeader := <-bl.HeaderChannel:
			bl.impl.ProcessHeader(newHeader)
		case <-ctx.Done():
			bl.Logger.Info("Header process stopped")
			return
		}
	}
}

// startPolling starts polling
func (bl *BaseListener) StartPolling(ctx context.Context, pollInterval time.Duration) {
	// How often to fire the passed in function in second
	interval := pollInterval

	// Setup the ticket and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)

	// start listening
	for {
		select {
		case <-ticker.C:
			header, err := bl.chainClient.HeaderByNumber(ctx, nil)
			if err == nil && header != nil {
				// send data to channel
				bl.HeaderChannel <- header
			}
		case <-ctx.Done():
			bl.Logger.Info("Polling stopped")
			ticker.Stop()
			return
		}
	}
}

func (bl *BaseListener) StartSubscription(ctx context.Context, subscription ethereum.Subscription) {
	for {
		select {
		case err := <-subscription.Err():
			// stop service
			bl.Logger.Error("Error while subscribing new blocks", "error", err)
			// bl.Stop()

			// cancel subscription
			if bl.cancelSubscription != nil {
				bl.cancelSubscription()
			}
			return
		case <-ctx.Done():
			bl.Logger.Info("Subscription stopped")
			return
		}
	}
}

// OnStop stops all necessary go routines
func (bl *BaseListener) Stop() {

	// cancel subscription if any
	if bl.cancelSubscription != nil {
		bl.cancelSubscription()
	}

	// cancel header process
	bl.cancelHeaderProcess()
}
