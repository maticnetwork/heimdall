package listener

import (
	"context"
	"errors"
	"time"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	ethereum "github.com/maticnetwork/bor"
	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/sethu/queue"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tendermint/libs/log"
	httpClient "github.com/tendermint/tendermint/rpc/client"
)

var (
	// ErrAlreadyStarted is returned when somebody tries to start an already
	// running service.
	ErrAlreadyStarted = errors.New("already started")
	// ErrAlreadyStopped is returned when somebody tries to stop an already
	// stopped service (without resetting it).
	ErrAlreadyStopped = errors.New("already stopped")
	// ErrNotStarted is returned when somebody tries to stop a not running
	// service.
	ErrNotStarted = errors.New("not started")
)

// Listener defines a block header listerner for Rootchain, Maticchain, Heimdall
type Listener interface {
	Start() error

	StartHeaderProcess(context.Context)

	StartPolling(context.Context, time.Duration)

	StartSubscription(context.Context, ethereum.Subscription)

	ProcessHeader(*types.Header)

	PublishEvent()

	Stop() error

	String() string

	SetLogger(log.Logger)
}

/*
Classical-inheritance-style service declarations. Services can be started, then
stopped, then optionally restarted.

Users can override the OnStart/OnStop methods. In the absence of errors, these
methods are guaranteed to be called at most once. If OnStart returns an error,
service won't be marked as started, so the user can call Start again.

A call to Reset will panic, unless OnReset is overwritten, allowing
OnStart/OnStop to be called again.

The caller must ensure that Start and Stop are not called concurrently.

It is ok to call Stop without calling Start first.

Typical usage:

	type FooService struct {
		BaseService
		// private fields
	}

	func NewFooService() *FooService {
		fs := &FooService{
			// init
		}
		fs.BaseService = *NewBaseService(log, "FooService", fs)
		return fs
	}

	func (fs *FooService) OnStart() error {
		fs.BaseService.OnStart() // Always call the overridden method.
		// initialize private fields
		// start subroutines, etc.
	}

	func (fs *FooService) OnStop() error {
		fs.BaseService.OnStop() // Always call the overridden method.
		// close/destroy private fields
		// stop subroutines, etc.
	}
*/
type BaseListener struct {
	Logger  log.Logger
	name    string
	started uint32 // atomic
	stopped uint32 // atomic
	quit    chan struct{}

	// The "subclass" of BaseService
	impl Listener

	// storage client
	storageClient *leveldb.DB

	// contract caller
	contractConnector helper.ContractCaller

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
}

// NewBaseListener creates a new BaseListener.
func NewBaseListener(logger log.Logger, name string, impl Listener) *BaseListener {
	if logger == nil {
		logger = log.NewNopLogger()
	}

	return &BaseListener{
		Logger: logger,
		name:   name,
		quit:   make(chan struct{}),
		impl:   impl,
	}
}

// SetLogger implements Service by setting a logger.
func (bl *BaseListener) SetLogger(l log.Logger) {
	bl.Logger = l
}

// OnStart starts new block subscription
func (bl *BaseListener) Start() error {

	// create cancellable context
	ctx, cancelSubscription := context.WithCancel(context.Background())
	bl.cancelSubscription = cancelSubscription

	// create cancellable context
	headerCtx, cancelHeaderProcess := context.WithCancel(context.Background())
	bl.cancelHeaderProcess = cancelHeaderProcess

	// start header process
	go bl.startHeaderProcess(headerCtx)

	// subscribe to new head
	subscription, err := bl.contractConnector.MainChainClient.SubscribeNewHead(ctx, bl.HeaderChannel)
	if err != nil {
		// start go routine to poll for new header using client object
		go bl.startPolling(ctx, helper.GetConfig().SyncerPollInterval)
	} else {
		// start go routine to listen new header using subscription
		go bl.startSubscription(ctx, subscription)
	}

	// subscribed to new head
	bl.Logger.Debug("Subscribed to new head")

	return nil
}

// String implements Service by returning a string representation of the service.
func (bl *BaseListener) String() string {
	return bl.name
}

// startHeaderProcess starts header process when they get new header
func (bl *BaseListener) startHeaderProcess(ctx context.Context) {
	for {
		select {
		case newHeader := <-bl.HeaderChannel:
			bl.impl.ProcessHeader(newHeader)
		case <-ctx.Done():
			return
		}
	}
}

// startPolling starts polling
func (bl *BaseListener) startPolling(ctx context.Context, pollInterval time.Duration) {
	// How often to fire the passed in function in second
	interval := pollInterval

	// Setup the ticket and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)

	// start listening
	for {
		select {
		case <-ticker.C:
			header, err := bl.contractConnector.MainChainClient.HeaderByNumber(ctx, nil)
			if err == nil && header != nil {
				// send data to channel
				bl.HeaderChannel <- header
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (bl *BaseListener) startSubscription(ctx context.Context, subscription ethereum.Subscription) {
	for {
		select {
		case err := <-subscription.Err():
			// stop service
			bl.Logger.Error("Error while subscribing new blocks", "error", err)
			// bl.Stop()

			// cancel subscription
			bl.cancelSubscription()
			return
		case <-ctx.Done():
			return
		}
	}
}

func (bl *BaseListener) processHeader(newHeader *types.Header) {}
