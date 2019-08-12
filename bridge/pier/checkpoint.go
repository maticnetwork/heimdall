package pier

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"sync"
	"time"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/ethereum/go-ethereum"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/helper"
	hmtypes "github.com/maticnetwork/heimdall/types"
)

// Checkpointer to propose
type Checkpointer struct {
	// Base service
	common.BaseService

	// storage client
	storageClient *leveldb.DB

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
	// queue connnector
	qConnector QueueConnector

	cliCtx cliContext.CLIContext
}

// NewCheckpointer returns new service object
func NewCheckpointer(connector QueueConnector) *Checkpointer {
	// create logger
	logger := Logger.With("module", HeimdallCheckpointer)

	// root chain instance
	rootchainInstance, err := helper.GetRootChainInstance()
	if err != nil {
		logger.Error("Error while getting root chain instance", "error", err)
		panic(err)
	}

	cliCtx := cliContext.NewCLIContext()
	cliCtx.BroadcastMode = client.BroadcastAsync

	// creating checkpointer object
	checkpointer := &Checkpointer{
		storageClient:     getBridgeDBInstance(viper.GetString(BridgeDBFlag)),
		MaticClient:       helper.GetMaticClient(),
		MaticRPCClient:    helper.GetMaticRPCClient(),
		MainClient:        helper.GetMainClient(),
		RootChainInstance: rootchainInstance,
		HeaderChannel:     make(chan *types.Header),
		qConnector:        connector,
		cliCtx:            cliCtx,
	}

	checkpointer.BaseService = *common.NewBaseService(logger, HeimdallCheckpointer, checkpointer)
	return checkpointer
}

// startHeaderProcess starts header process when they get new header
func (c *Checkpointer) startHeaderProcess(ctx context.Context) {
	for {
		select {
		case newHeader := <-c.HeaderChannel:
			c.sendRequest(newHeader)
		case <-ctx.Done():
			return
		}
	}
}

// OnStart starts new block subscription
func (c *Checkpointer) OnStart() error {
	c.BaseService.OnStart() // Always call the overridden method.

	// create cancellable context
	ctx, cancelSubscription := context.WithCancel(context.Background())
	c.cancelSubscription = cancelSubscription

	// create cancellable context
	headerCtx, cancelHeaderProcess := context.WithCancel(context.Background())
	c.cancelHeaderProcess = cancelHeaderProcess

	// start header process
	go c.startHeaderProcess(headerCtx)

	// subscribe to new head
	subscription, err := c.MaticClient.SubscribeNewHead(ctx, c.HeaderChannel)
	if err != nil {
		// start go routine to poll for new header using client object
		go c.startPolling(ctx, helper.GetConfig().CheckpointerPollInterval)
	} else {
		// start go routine to listen new header using subscription
		go c.startSubscription(ctx, subscription)
	}

	// subscribed to new head
	c.Logger.Debug("Subscribed to new head")

	return nil
}

// OnStop stops all necessary go routines
func (c *Checkpointer) OnStop() {
	c.BaseService.OnStop() // Always call the overridden method.

	// close bridge db instance
	closeBridgeDBInstance()

	// cancel subscription if any
	c.cancelSubscription()

	// cancel header process
	c.cancelHeaderProcess()
}

func (c *Checkpointer) startPolling(ctx context.Context, pollInterval int) {
	// How often to fire the passed in function in second
	interval := time.Duration(pollInterval) * time.Millisecond

	// Setup the ticket and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)

	// start listening
	for {
		select {
		case <-ticker.C:
			if isProposer() {
				header, err := c.MaticClient.HeaderByNumber(ctx, nil)
				if err == nil && header != nil {
					// send data to channel
					c.HeaderChannel <- header
				} else if err != nil {
					c.Logger.Error("Unable to fetch header by number", "Error", err)
				}
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (c *Checkpointer) startSubscription(ctx context.Context, subscription ethereum.Subscription) {
	for {
		select {
		case err := <-subscription.Err():
			// stop service
			c.Logger.Error("Error while subscribing new blocks", "error", err)
			c.Stop()

			// cancel subscription
			c.cancelSubscription()
			return
		case <-ctx.Done():
			return
		}
	}
}

func (c *Checkpointer) sendRequest(newHeader *types.Header) {
	c.Logger.Debug("New block detected", "blockNumber", newHeader.Number)

	// fetch data
	// take decision
	// spawn go routines with timeouts to take care of transactions

	// fetch contract state
	contractState := make(chan ContractCheckpoint, 1)
	var lastContractCheckpoint ContractCheckpoint

	// fetch heimdall state
	bufferChan := make(chan HeimdallCheckpoint, 1)
	var bufferedCheckpoint HeimdallCheckpoint

	commitedChan := make(chan HeimdallCheckpoint, 1)
	var commitedCheckpoint HeimdallCheckpoint

	var wg sync.WaitGroup
	wg.Add(3)

	c.Logger.Debug("Collecting contract and heimdall checkpoint state")
	go c.genHeaderDetailContract(newHeader.Number.Uint64(), &wg, contractState)
	go c.fetchBufferedCheckpoint(&wg, bufferChan)
	go c.fetchCommittedCheckpoint(&wg, commitedChan)

	// wait for state collection
	wg.Wait()

	c.Logger.Info("Fetched contract and heimdall states")

	lastContractCheckpoint = <-contractState
	if lastContractCheckpoint.err != nil {
		c.Logger.Error("Error fetching details from contract ", "Error", lastContractCheckpoint.err)
		return
	}

	c.Logger.Debug("Contract Details fetched", "Start", lastContractCheckpoint.start, "End", lastContractCheckpoint.end, "Currentheader", lastContractCheckpoint.currentHeaderBlock)

	bufferedCheckpoint = <-bufferChan
	if !bufferedCheckpoint.found {
		c.Logger.Info("Buffer not found, sending new checkpoint", "Found", bufferedCheckpoint.found)
	} else if bufferedCheckpoint.start != 0 {
		c.Logger.Debug("Checkpoint found in buffer", "Start", bufferedCheckpoint.start, "End", bufferedCheckpoint.end)
	} else {
		c.Logger.Debug("Sending first checkpoint")
	}

	// ACK needs to be sent
	if bufferedCheckpoint.end+1 == bufferedCheckpoint.start {
		c.Logger.Debug("Detected mainchain checkpoint,sending ACK", "HeimdallEnd", bufferedCheckpoint.end, "ContractStart", bufferedCheckpoint.start)
		headerNumber := lastContractCheckpoint.currentHeaderBlock.Sub(lastContractCheckpoint.currentHeaderBlock, big.NewInt(int64(helper.GetConfig().ChildBlockInterval)))
		// create and send checkpoint message
		msg := checkpoint.NewMsgCheckpointAck(headerNumber.Uint64(), uint64(time.Now().Unix()))
		resp, err := helper.CreateAndSendTx(msg, c.cliCtx)
		if err != nil {
			c.Logger.Error("Unable to send checkpoint ack to heimdall", "Error", err, "HeaderIndex", headerNumber.String())
		}
		c.Logger.Debug("Checkpoint ACK tx commited", "TxHash", resp.TxHash, "HeaderIndex", headerNumber.String())
		return
	}

	start := lastContractCheckpoint.start
	end := lastContractCheckpoint.end

	// Get root hash
	root, err := checkpoint.GetHeaders(start, end)
	if err != nil {
		return
	}
	c.Logger.Info("New checkpoint header created", "start", start, "end", end, "root", ethCommon.BytesToHash(root))
	checkpointTxRes, err := helper.CreateAndSendTx(
		checkpoint.NewMsgCheckpointBlock(
			ethCommon.BytesToAddress(helper.GetAddress()),
			start,
			end,
			ethCommon.BytesToHash(root),
			uint64(time.Now().Unix()),
		),
		c.cliCtx,
	)

	if err != nil {
		c.Logger.Error("Error sending checkpoint tx", "error", err, "start", start, "end", end)
		return
	}
	c.Logger.Info("Checkpoint sent successfully", "hash", checkpointTxRes.TxHash, "start", start, "end", end, "root", hex.EncodeToString(root))
}

func (c *Checkpointer) genHeaderDetailContract(lastHeader uint64, wg *sync.WaitGroup, contractState chan<- ContractCheckpoint) {
	defer wg.Done()
	lastCheckpointEnd, err := c.RootChainInstance.CurrentChildBlock(nil)
	if err != nil {
		c.Logger.Error("Error while fetching current child block from rootchain", "error", err)
		return
	}
	var start, end uint64
	start = lastCheckpointEnd.Uint64()

	// add 1 if start > 0
	if start > 0 {
		start = start + 1
	}

	// get diff
	diff := lastHeader - start + 1

	// process if diff > 0 (positive)
	if diff > 0 {
		expectedDiff := diff - diff%helper.GetConfig().AvgCheckpointLength
		if expectedDiff > 0 {
			expectedDiff = expectedDiff - 1
		}

		// cap with max checkpoint length
		if expectedDiff > helper.GetConfig().MaxCheckpointLength-1 {
			expectedDiff = helper.GetConfig().MaxCheckpointLength - 1
		}

		// get end result
		end = expectedDiff + start

		c.Logger.Debug("Calculating checkpoint eligibility", "latest", lastHeader, "start", start, "end", end)
	}
	currentHeaderBlockNumber, err := c.RootChainInstance.CurrentHeaderBlock(nil)
	if err != nil {
		c.Logger.Error("Error while fetching current header block number from rootchain", "error", err)
		contractState <- NewContractCheckpoint(0, 0, currentHeaderBlockNumber, err)
		return
	}
	currentHeaderNum := currentHeaderBlockNumber

	// Handle when block producers go down
	if end == 0 || end == start || (0 < diff && diff < helper.GetConfig().AvgCheckpointLength) {
		c.Logger.Debug("Fetching last header block to calculate time")
		// fetch current header block
		currentHeaderBlock, err := c.RootChainInstance.HeaderBlock(nil, currentHeaderBlockNumber.Sub(currentHeaderBlockNumber, big.NewInt(int64(helper.GetConfig().ChildBlockInterval))))
		if err != nil {
			c.Logger.Error("Error while fetching current header block object from rootchain", "error", err)
			contractState <- NewContractCheckpoint(0, 0, currentHeaderBlockNumber, err)
			return
		}
		lastCheckpointTime := currentHeaderBlock.CreatedAt.Int64()
		currentTime := time.Now().Unix()
		if currentTime-lastCheckpointTime > int64(helper.GetConfig().MaxCheckpointLength*2) {
			c.Logger.Info("Force push checkpoint", "currentTime", currentTime, "lastCheckpointTime", lastCheckpointTime, "defaultForcePushInterval", defaultForcePushInterval, "end", lastHeader)
			end = lastHeader
		}
	}

	if end == 0 || start >= end {
		c.Logger.Info("Waiting for 256 blocks or invalid start end formation", "Start", start, "End", end)
		contractState <- NewContractCheckpoint(0, 0, currentHeaderBlockNumber, errors.New("Invalid start end formation"))
		return
	}
	contractCheckpointData := NewContractCheckpoint(start, end, currentHeaderNum, nil)
	contractState <- contractCheckpointData
	return
}

// fetchBufferedCheckpoint fetch buffered checkpoint from heimdall
func (c *Checkpointer) fetchBufferedCheckpoint(wg *sync.WaitGroup, bufferedCheckpoint chan<- HeimdallCheckpoint) {
	defer wg.Done()
	c.Logger.Info("Fetching checkpoint in buffer")
	_checkpoint, err := c.fetchCheckpoint(LatestCheckpoint)
	if err != nil {
		c.Logger.Error("Error while fetching data from server", "error", err)
		bufferedCheckpoint <- NewHeimdallCheckpoint(_checkpoint.StartBlock, _checkpoint.EndBlock, false)
		return
	}
	bufferedCheckpoint <- NewHeimdallCheckpoint(_checkpoint.StartBlock, _checkpoint.EndBlock, true)
	return
}

// fetchCommittedCheckpoint fetches latest committed checkpoint from heimdall
func (c *Checkpointer) fetchCommittedCheckpoint(wg *sync.WaitGroup, lastCheckpoint chan<- HeimdallCheckpoint) {
	defer wg.Done()
	c.Logger.Info("Fetching last committed checkpoint")
	_checkpoint, err := c.fetchCheckpoint(LatestCheckpoint)
	if err != nil {
		c.Logger.Error("Error while fetching data from server", "error", err)
		lastCheckpoint <- NewHeimdallCheckpoint(_checkpoint.StartBlock, _checkpoint.EndBlock, false)
		return
	}
	lastCheckpoint <- NewHeimdallCheckpoint(_checkpoint.StartBlock, _checkpoint.EndBlock, true)
	return
}

// fetchCheckpoint fetches checkpoint from given URL
func (c *Checkpointer) fetchCheckpoint(url string) (checkpoint hmtypes.CheckpointBlockHeader, err error) {
	resp, err := http.Get(url)
	if err != nil {
		c.Logger.Error("Unable to send request to get proposer", "Error", err)
		return checkpoint, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.Logger.Error("Unable to read data from response", "Error", err)
			return checkpoint, err
		}
		if err := json.Unmarshal(body, &checkpoint); err != nil {
			c.Logger.Error("Error unmarshalling checkpoint", "error", err)
			return checkpoint, err
		}
		return checkpoint, nil
	}
	return checkpoint, fmt.Errorf("Invalid response from rest server. Status: %v URL: %v", resp.Status, url)
}
