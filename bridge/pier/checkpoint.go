package pier

import (
	"context"
	"encoding/hex"
	"math/big"
	"os"
	"time"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/checkpoint"
	checkpointTx "github.com/maticnetwork/heimdall/checkpoint/rest"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/helper"
)

// MaticCheckpointer to propose
type MaticCheckpointer struct {
	// Base service
	common.BaseService

	// Redis client
	redisClient *redis.Client

	// ETH client
	MaticClient *ethclient.Client
	// ETH RPC client
	MaticRPCClient *rpc.Client
	// Mainchain client
	MainClient *ethclient.Client
	// Rootchain instance
	RootChainInstance *rootchain.Rootchain
	// header channel
	HeaderChannel chan *types.Header
	// cancel function for poll/subscription
	cancelSubscription context.CancelFunc
	// header listener subscription
	cancelHeaderProcess context.CancelFunc
}

// NewMaticCheckpointer returns new service object
func NewMaticCheckpointer() *MaticCheckpointer {
	// create logger
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", maticCheckpointer)

	// root chain instance
	rootchainInstance, err := helper.GetRootChainInstance()
	if err != nil {
		logger.Error("Error while getting root chain instance", "error", err)
		panic(err)
	}

	redisOptions, err := redis.ParseURL(viper.GetString(redisURL))
	if err != nil {
		logger.Error("Error while redis instance", "error", err)
		panic(err)
	}

	// creating checkpointer object
	checkpointer := &MaticCheckpointer{
		redisClient:       redis.NewClient(redisOptions),
		MaticClient:       helper.GetMaticClient(),
		MaticRPCClient:    helper.GetMaticRPCClient(),
		MainClient:        helper.GetMainClient(),
		RootChainInstance: rootchainInstance,
		HeaderChannel:     make(chan *types.Header),
	}

	checkpointer.BaseService = *common.NewBaseService(logger, maticCheckpointer, checkpointer)
	return checkpointer
}

// StartHeaderProcess starts header process when they get new header
func (checkpointer *MaticCheckpointer) StartHeaderProcess(ctx context.Context) {
	for {
		select {
		case newHeader := <-checkpointer.HeaderChannel:
			checkpointer.sendRequest(newHeader)
		case <-ctx.Done():
			return
		}
	}
}

// OnStart starts new block subscription
func (checkpointer *MaticCheckpointer) OnStart() error {
	checkpointer.BaseService.OnStart() // Always call the overridden method.

	// create cancellable context
	ctx, cancelSubscription := context.WithCancel(context.Background())
	checkpointer.cancelSubscription = cancelSubscription

	// create cancellable context
	headerCtx, cancelHeaderProcess := context.WithCancel(context.Background())
	checkpointer.cancelHeaderProcess = cancelHeaderProcess

	// start header process
	go checkpointer.StartHeaderProcess(headerCtx)

	// subscribe to new head
	subscription, err := checkpointer.MaticClient.SubscribeNewHead(ctx, checkpointer.HeaderChannel)
	if err != nil {
		// start go routine to poll for new header using client object
		go checkpointer.StartPolling(ctx, defaultPollInterval)
	} else {
		// start go routine to listen new header using subscription
		go checkpointer.StartSubscription(ctx, subscription)
	}

	// subscribed to new head
	checkpointer.Logger.Debug("Subscribed to new head")

	return nil
}

// OnStop stops all necessary go routines
func (checkpointer *MaticCheckpointer) OnStop() {
	checkpointer.BaseService.OnStop() // Always call the overridden method.

	// cancel subscription if any
	checkpointer.cancelSubscription()

	// cancel header process
	checkpointer.cancelHeaderProcess()
}

func (checkpointer *MaticCheckpointer) StartPolling(ctx context.Context, pollInterval int) {
	// How often to fire the passed in function in second
	interval := time.Duration(pollInterval) * time.Millisecond

	// Setup the ticket and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)

	// start listening
	for {
		select {
		case <-ticker.C:
			header, err := checkpointer.MaticClient.HeaderByNumber(ctx, nil)
			if err == nil && header != nil {
				// send data to channel
				checkpointer.HeaderChannel <- header
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (checkpointer *MaticCheckpointer) StartSubscription(ctx context.Context, subscription ethereum.Subscription) {
	for {
		select {
		case err := <-subscription.Err():
			// stop service
			checkpointer.Logger.Error("Error while subscribing new blocks", "error", err)
			checkpointer.Stop()

			// cancel subscription
			checkpointer.cancelSubscription()
			return
		case <-ctx.Done():
			return
		}
	}
}

func (checkpointer *MaticCheckpointer) sendRequest(newHeader *types.Header) {
	checkpointer.Logger.Debug("New block detected", "blockNumber", newHeader.Number)
	lastCheckpointEnd, err := checkpointer.RootChainInstance.CurrentChildBlock(nil)
	if err != nil {
		checkpointer.Logger.Error("Error while fetching current child block from rootchain", "error", err)
		return
	}

	latest := newHeader.Number.Uint64()
	start := lastCheckpointEnd.Uint64()
	var end uint64

	// add 1 if start > 0
	if start > 0 {
		start = start + 1
	}

	// get diff
	diff := latest - start + 1

	// process if diff > 0 (positive)
	if diff > 0 {
		expectedDiff := diff - diff%defaultCheckpointLength
		if expectedDiff > 0 {
			expectedDiff = expectedDiff - 1
		}

		// cap with max checkpoint length
		if expectedDiff > maxCheckpointLength-1 {
			expectedDiff = maxCheckpointLength - 1
		}

		// get end result
		end = expectedDiff + start

		checkpointer.Logger.Debug("Calculating checkpoint eligibility", "latest", latest, "start", start, "end", end)
	}

	// Handle when block producers go down
	if end == 0 || end == start || (0 < diff && diff < defaultCheckpointLength) {
		currentHeaderBlockNumber, err := checkpointer.RootChainInstance.CurrentHeaderBlock(nil)
		if err != nil {
			checkpointer.Logger.Error("Error while fetching current header block number from rootchain", "error", err)
			return
		}

		// fetch current header block
		currentHeaderBlock, err := checkpointer.RootChainInstance.HeaderBlock(nil, currentHeaderBlockNumber.Sub(currentHeaderBlockNumber, big.NewInt(1)))
		if err != nil {
			checkpointer.Logger.Error("Error while fetching current header block object from rootchain", "error", err)
			return
		}

		lastCheckpointTime := currentHeaderBlock.CreatedAt.Int64()
		currentTime := time.Now().Unix()
		if currentTime-lastCheckpointTime > defaultForcePushInterval {
			checkpointer.Logger.Info("Force push checkpoint", "currentTime", currentTime, "lastCheckpointTime", lastCheckpointTime, "defaultForcePushInterval", defaultForcePushInterval)
			end = latest
		}
	}

	if end == 0 || start >= end {
		return
	}

	// Get root hash
	root := checkpoint.GetHeaders(start, end)
	checkpointer.Logger.Info("New checkpoint header created", "latest", latest, "start", start, "end", end, "root", root)

	// TODO submit checkcoint
	txBytes, err := checkpointTx.CreateTxBytes(checkpointTx.EpochCheckpoint{
		RootHash:   root,
		StartBlock: start,
		EndBlock:   end,
	})

	if err != nil {
		checkpointer.Logger.Error("Error while creating tx bytes", "error", err)
		return
	}

	resp, err := checkpointTx.SendTendermintRequest(cliContext.NewCLIContext(), txBytes)
	if err != nil {
		checkpointer.Logger.Error("Error while sending request to Tendermint", "error", err)
		return
	}

	checkpointer.Logger.Error("Checkpoint sent successfully", "hash", hex.EncodeToString(resp.Hash), "start", start, "end", end, "root", root)
}
