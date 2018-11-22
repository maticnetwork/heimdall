package pier

import (
	"context"
	"fmt"
	"os"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/contracts/stakemanager"
	"github.com/maticnetwork/heimdall/helper"
)

const (
	redisURL    = "redis-url"
	chainSyncer = "chain-syncer"
)

// ChainSyncer syncs validators and checkpoints
type ChainSyncer struct {
	// Base service
	common.BaseService

	// Redis client
	redisClient *redis.Client

	// Mainchain client
	MainClient *ethclient.Client
	// Rootchain instance
	RootChainInstance *rootchain.Rootchain
	// Stake manager instance
	StakeManagerInstance *stakemanager.Stakemanager
	// header channel
	HeaderChannel chan *types.Header
	// cancel function for poll/subscription
	cancelSubscription context.CancelFunc
	// header listener subscription
	cancelHeaderProcess context.CancelFunc
}

// NewChainSyncer returns new service object
func NewChainSyncer() *ChainSyncer {
	// create logger
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", chainSyncer)

	// root chain instance
	rootchainInstance, err := helper.GetRootChainInstance()
	if err != nil {
		logger.Error("Error while getting root chain instance", "error", err)
		panic(err)
	}

	// stake manager instance
	stakeManagerInstance, err := helper.GetStakeManagerInstance()
	if err != nil {
		logger.Error("Error while getting stake manager instance", "error", err)
		panic(err)
	}

	redisOptions, err := redis.ParseURL(viper.GetString(redisURL))
	fmt.Println("redisOptions", redisOptions, "viper.GetString(redisURL)", viper.GetString(redisURL))
	if err == nil {
		logger.Error("Error while redis instance", "error", err)
		panic(err)
	}

	// creating syncer object
	syncer := &ChainSyncer{
		redisClient:          redis.NewClient(redisOptions),
		MainClient:           helper.GetMainClient(),
		RootChainInstance:    rootchainInstance,
		StakeManagerInstance: stakeManagerInstance,
		HeaderChannel:        make(chan *types.Header),
	}

	syncer.BaseService = *common.NewBaseService(logger, ChainSyncer, syncer)
	return syncer
}

// StartHeaderProcess starts header process when they get new header
func (syncer *ChainSyncer) StartHeaderProcess(ctx context.Context) {
	for {
		select {
		case newHeader := <-syncer.HeaderChannel:
			syncer.sendRequest(newHeader)
		case <-ctx.Done():
			return
		}
	}
}

// OnStart starts new block subscription
func (syncer *ChainSyncer) OnStart() error {
	syncer.BaseService.OnStart() // Always call the overridden method.

	// create cancellable context
	ctx, cancelSubscription := context.WithCancel(context.Background())
	syncer.cancelSubscription = cancelSubscription

	// create cancellable context
	headerCtx, cancelHeaderProcess := context.WithCancel(context.Background())
	syncer.cancelHeaderProcess = cancelHeaderProcess

	// start header process
	go syncer.StartHeaderProcess(headerCtx)

	// subscribe to new head
	subscription, err := syncer.MainClient.SubscribeNewHead(ctx, syncer.HeaderChannel)
	if err != nil {
		// start go routine to poll for new header using client object
		go syncer.StartPolling(ctx, defaultPollInterval)
	} else {
		// start go routine to listen new header using subscription
		go syncer.StartSubscription(ctx, subscription)
	}

	// subscribed to new head
	syncer.Logger.Debug("Subscribed to new head")

	return nil
}

// OnStop stops all necessary go routines
func (syncer *ChainSyncer) OnStop() {
	syncer.BaseService.OnStop() // Always call the overridden method.

	// cancel subscription if any
	syncer.cancelSubscription()

	// cancel header process
	syncer.cancelHeaderProcess()
}

func (syncer *ChainSyncer) StartPolling(ctx context.Context, pollInterval int) {
	// How often to fire the passed in function in second
	interval := time.Duration(pollInterval) * time.Millisecond

	// Setup the ticket and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)

	// start listening
	for {
		select {
		case <-ticker.C:
			header, err := syncer.MainClient.HeaderByNumber(ctx, nil)
			if err == nil && header != nil {
				// send data to channel
				syncer.HeaderChannel <- header
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (syncer *ChainSyncer) StartSubscription(ctx context.Context, subscription ethereum.Subscription) {
	for {
		select {
		case err := <-subscription.Err():
			// stop service
			syncer.Logger.Error("Error while subscribing new blocks", "error", err)
			syncer.Stop()

			// cancel subscription
			syncer.cancelSubscription()
			return
		case <-ctx.Done():
			return
		}
	}
}
