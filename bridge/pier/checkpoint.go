package pier

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/big"
	"net/http"
	"os"
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
	"github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/helper"
	hmtypes "github.com/maticnetwork/heimdall/types"
	"github.com/cosmos/cosmos-sdk/client"
)

// MaticCheckpointer to propose
type MaticCheckpointer struct {
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

	cliCtx cliContext.CLIContext
}

type ContractCheckpoint struct {
	start              uint64
	end                uint64
	currentHeaderBlock *big.Int
	err                error
}

func NewContractCheckpoint(_start uint64, _end uint64, _currentHeaderBlock *big.Int, _err error) ContractCheckpoint {
	return ContractCheckpoint{
		start:              _start,
		end:                _end,
		currentHeaderBlock: _currentHeaderBlock,
		err:                _err,
	}
}

type HeimdallCheckpoint struct {
	start uint64
	end   uint64
	found bool
}

func NewHeimdallCheckpoint(_start uint64, _end uint64, _found bool) HeimdallCheckpoint {
	return HeimdallCheckpoint{
		start: _start,
		end:   _end,
		found: _found,
	}
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

	cliCtx := cliContext.NewCLIContext()
	cliCtx.BroadcastMode = client.BroadcastAsync

	// creating checkpointer object
	checkpointer := &MaticCheckpointer{
		storageClient:     getBridgeDBInstance(viper.GetString(bridgeDBFlag)),
		MaticClient:       helper.GetMaticClient(),
		MaticRPCClient:    helper.GetMaticRPCClient(),
		MainClient:        helper.GetMainClient(),
		RootChainInstance: rootchainInstance,
		HeaderChannel:     make(chan *types.Header),
		cliCtx:            cliCtx,
	}

	checkpointer.BaseService = *common.NewBaseService(logger, maticCheckpointer, checkpointer)
	return checkpointer
}

// startHeaderProcess starts header process when they get new header
func (checkpointer *MaticCheckpointer) startHeaderProcess(ctx context.Context) {
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
	go checkpointer.startHeaderProcess(headerCtx)

	// subscribe to new head
	subscription, err := checkpointer.MaticClient.SubscribeNewHead(ctx, checkpointer.HeaderChannel)
	if err != nil {
		// start go routine to poll for new header using client object
		go checkpointer.startPolling(ctx, helper.GetConfig().CheckpointerPollInterval)
	} else {
		// start go routine to listen new header using subscription
		go checkpointer.startSubscription(ctx, subscription)
	}

	// subscribed to new head
	checkpointer.Logger.Debug("Subscribed to new head")

	return nil
}

// OnStop stops all necessary go routines
func (checkpointer *MaticCheckpointer) OnStop() {
	checkpointer.BaseService.OnStop() // Always call the overridden method.

	// close bridge db instance
	closeBridgeDBInstance()

	// cancel subscription if any
	checkpointer.cancelSubscription()

	// cancel header process
	checkpointer.cancelHeaderProcess()
}

func (checkpointer *MaticCheckpointer) startPolling(ctx context.Context, pollInterval int) {
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
				header, err := checkpointer.MaticClient.HeaderByNumber(ctx, nil)
				if err == nil && header != nil {
					// send data to channel
					checkpointer.HeaderChannel <- header
				} else if err != nil {
					checkpointer.Logger.Error("Unable to fetch header by number from matic", "Error", err)
				}
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (checkpointer *MaticCheckpointer) startSubscription(ctx context.Context, subscription ethereum.Subscription) {
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
	contractState := make(chan ContractCheckpoint, 1)
	var lastContractCheckpoint ContractCheckpoint
	heimdallState := make(chan HeimdallCheckpoint, 1)
	var lastHeimdallCheckpoint HeimdallCheckpoint

	var wg sync.WaitGroup
	wg.Add(2)

	checkpointer.Logger.Debug("Collecting all required data")
	go checkpointer.genHeaderDetailContract(newHeader.Number.Uint64(), &wg, contractState)
	go checkpointer.getLastCheckpointStore(&wg, heimdallState)

	wg.Wait()

	checkpointer.Logger.Info("Done waiting", "contract", contractState, "heimdall", heimdallState)
	lastContractCheckpoint = <-contractState
	if lastContractCheckpoint.err != nil {
		checkpointer.Logger.Error("Error fetching details from contract ", "Error", lastContractCheckpoint.err)
		return
	}
	checkpointer.Logger.Debug("Contract Details fetched", "Start", lastContractCheckpoint.start, "End", lastContractCheckpoint.end, "Currentheader", lastContractCheckpoint.currentHeaderBlock)

	lastHeimdallCheckpoint = <-heimdallState
	if !lastHeimdallCheckpoint.found {
		checkpointer.Logger.Info("Buffer not found , sending new checkpoint", "Found", lastHeimdallCheckpoint.found)
	} else {
		checkpointer.Logger.Debug("Checkpoint found in buffer", "Start", lastHeimdallCheckpoint.start, "End", lastHeimdallCheckpoint.end)
	}

	// ACK needs to be sent
	if lastHeimdallCheckpoint.end+1 == lastContractCheckpoint.start {
		checkpointer.Logger.Debug("Detected mainchain checkpoint,sending ACK", "HeimdallEnd", lastHeimdallCheckpoint.end, "ContractStart", lastHeimdallCheckpoint.start)
		headerNumber := lastContractCheckpoint.currentHeaderBlock.Sub(lastContractCheckpoint.currentHeaderBlock, big.NewInt(int64(helper.GetConfig().ChildBlockInterval)))
		msg := checkpoint.NewMsgCheckpointAck(headerNumber.Uint64(), uint64(time.Now().Unix()))
		txBytes, err := helper.CreateTxBytes(msg)
		if err != nil {
			checkpointer.Logger.Error("Error while creating tx bytes", "error", err)
			return
		}
		// send tendermint request
		_, err = helper.SendTendermintRequest(checkpointer.cliCtx, txBytes)
		if err != nil {
			checkpointer.Logger.Error("Error while sending request to Tendermint", "error", err)
			return
		}
		return
	}
	start := lastContractCheckpoint.start
	end := lastContractCheckpoint.end

	// Get root hash
	root, err := checkpoint.GetHeaders(start, end)
	if err != nil {
		return
	}

	checkpointer.Logger.Info("New checkpoint header created", "start", start, "end", end, "root", hex.EncodeToString(root))

	// TODO submit checkcoint
	txBytes, err := helper.CreateTxBytes(
		checkpoint.NewMsgCheckpointBlock(
			ethCommon.BytesToAddress(helper.GetAddress()),
			start,
			end,
			ethCommon.BytesToHash(root),
			uint64(time.Now().Unix()),
		),
	)

	if err != nil {
		checkpointer.Logger.Error("Error while creating tx bytes", "error", err)
		return
	}

	resp, err := helper.SendTendermintRequest(checkpointer.cliCtx, txBytes)
	if err != nil {
		checkpointer.Logger.Error("Error while sending request to Tendermint", "error", err)
		return
	}

	checkpointer.Logger.Info("Checkpoint sent successfully", "hash", resp.TxHash, "start", start, "end", end, "root", hex.EncodeToString(root))
}

func (checkpointer *MaticCheckpointer) genHeaderDetailContract(lastHeader uint64, wg *sync.WaitGroup, contractState chan<- ContractCheckpoint) {
	defer wg.Done()
	lastCheckpointEnd, err := checkpointer.RootChainInstance.CurrentChildBlock(nil)
	if err != nil {
		checkpointer.Logger.Error("Error while fetching current child block from rootchain", "error", err)
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

		checkpointer.Logger.Debug("Calculating checkpoint eligibility", "latest", lastHeader, "start", start, "end", end)
	}
	currentHeaderBlockNumber, err := checkpointer.RootChainInstance.CurrentHeaderBlock(nil)
	if err != nil {
		checkpointer.Logger.Error("Error while fetching current header block number from rootchain", "error", err)
		contractState <- NewContractCheckpoint(0, 0, currentHeaderBlockNumber, err)
		return
	}
	currentHeaderNum := currentHeaderBlockNumber

	// Handle when block producers go down
	if end == 0 || end == start || (0 < diff && diff < helper.GetConfig().AvgCheckpointLength) {
		checkpointer.Logger.Debug("Fetching last header block to calculate time")
		// fetch current header block
		currentHeaderBlock, err := checkpointer.RootChainInstance.HeaderBlock(nil, currentHeaderBlockNumber.Sub(currentHeaderBlockNumber, big.NewInt(int64(helper.GetConfig().ChildBlockInterval))))
		if err != nil {
			checkpointer.Logger.Error("Error while fetching current header block object from rootchain", "error", err)
			contractState <- NewContractCheckpoint(0, 0, currentHeaderBlockNumber, err)
			return
		}
		lastCheckpointTime := currentHeaderBlock.CreatedAt.Int64()
		currentTime := time.Now().Unix()
		if currentTime-lastCheckpointTime > int64(helper.GetConfig().MaxCheckpointLength*2) {
			checkpointer.Logger.Info("Force push checkpoint", "currentTime", currentTime, "lastCheckpointTime", lastCheckpointTime, "defaultForcePushInterval", defaultForcePushInterval, "end", lastHeader)
			end = lastHeader
		}
	}

	if end == 0 || start >= end {
		checkpointer.Logger.Info("Waiting for 256 blocks or invalid start end formation", "Start", start, "End", end)
		contractState <- NewContractCheckpoint(0, 0, currentHeaderBlockNumber, errors.New("Invalid start end formation"))
		return
	}
	contractCheckpointData := NewContractCheckpoint(start, end, currentHeaderNum, nil)
	contractState <- contractCheckpointData
	return
}

func (checkpointer *MaticCheckpointer) getLastCheckpointStore(wg *sync.WaitGroup, heimdallState chan<- HeimdallCheckpoint) {
	defer wg.Done()
	checkpointer.Logger.Info("Fetching checkpoint in buffer")
	var _checkpoint hmtypes.CheckpointBlockHeader
	resp, err := http.Get(lastCheckpointURL)
	if err != nil {
		checkpointer.Logger.Error("Unable to send request to get proposer", "Error", err)
		heimdallState <- NewHeimdallCheckpoint(_checkpoint.StartBlock, _checkpoint.EndBlock, false)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			checkpointer.Logger.Error("Unable to read data from response", "Error", err)
			heimdallState <- NewHeimdallCheckpoint(_checkpoint.StartBlock, _checkpoint.EndBlock, false)
			return
		}
		if err := json.Unmarshal(body, &_checkpoint); err != nil {
			checkpointer.Logger.Error("Error unmarshalling checkpoint", "error", err)
			heimdallState <- NewHeimdallCheckpoint(_checkpoint.StartBlock, _checkpoint.EndBlock, false)
			return
		}
		heimdallState <- NewHeimdallCheckpoint(_checkpoint.StartBlock, _checkpoint.EndBlock, true)
		return
	}
	heimdallState <- NewHeimdallCheckpoint(_checkpoint.StartBlock, _checkpoint.EndBlock, false)
	return
}
