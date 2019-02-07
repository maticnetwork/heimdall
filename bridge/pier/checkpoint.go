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
	cliCtx.Async = true

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
		go checkpointer.startPolling(ctx, defaultPollInterval)
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
	contractStart, contractEnd, currentHeaderBlock, err := checkpointer.GenHeaderDetailContract(newHeader.Number.Uint64())
	if err != nil {
		checkpointer.Logger.Error("Error fetching details from contract ", "Error", err)
		return
	}
	checkpointer.Logger.Debug("Contract Details fetched", "Start", contractStart, "End", contractEnd, "Currentheader", currentHeaderBlock.String())
	hmStart, hmEnd, found := getLastCheckpointStore()
	if !found {
		checkpointer.Logger.Error("Buffer not found , sending new checkpoint", "Bool", found)
	} else {
		checkpointer.Logger.Debug("Checkpoint found in buffer", "Start", hmStart, "End", hmEnd)
	}
	// ACK needs to be sent
	if hmEnd+1 == contractStart {
		checkpointer.Logger.Debug("Detected mainchain checkpoint,sending ACK", "HeimdallStart", hmStart, "HeimdallEnd", hmEnd, "ContractEnd", contractEnd)
		headerNumber := currentHeaderBlock.Sub(currentHeaderBlock, big.NewInt(10000))
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
	start := contractStart
	end := contractEnd

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

	checkpointer.Logger.Info("Checkpoint sent successfully", "hash", resp.Hash.String(), "start", start, "end", end, "root", hex.EncodeToString(root))
}

func (checkpointer *MaticCheckpointer) GenHeaderDetailContract(latest uint64) (start, end uint64, currentHeaderNum *big.Int, err error) {
	lastCheckpointEnd, err := checkpointer.RootChainInstance.CurrentChildBlock(nil)
	if err != nil {
		checkpointer.Logger.Error("Error while fetching current child block from rootchain", "error", err)
		return
	}

	start = lastCheckpointEnd.Uint64()

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
	currentHeaderBlockNumber, err := checkpointer.RootChainInstance.CurrentHeaderBlock(nil)
	if err != nil {
		checkpointer.Logger.Error("Error while fetching current header block number from rootchain", "error", err)
		return 0, 0, currentHeaderBlockNumber, err
	}
	currentHeaderNum = currentHeaderBlockNumber

	// Handle when block producers go down
	if end == 0 || end == start || (0 < diff && diff < defaultCheckpointLength) {
		// fetch current header block
		currentHeaderBlock, err := checkpointer.RootChainInstance.HeaderBlock(nil, currentHeaderBlockNumber.Sub(currentHeaderBlockNumber, big.NewInt(1)))
		if err != nil {
			checkpointer.Logger.Error("Error while fetching current header block object from rootchain", "error", err)
			return 0, 0, currentHeaderBlockNumber, err
		}
		lastCheckpointTime := currentHeaderBlock.CreatedAt.Int64()
		currentTime := time.Now().Unix()
		if currentTime-lastCheckpointTime > defaultForcePushInterval {
			checkpointer.Logger.Info("Force push checkpoint", "currentTime", currentTime, "lastCheckpointTime", lastCheckpointTime, "defaultForcePushInterval", defaultForcePushInterval)
			end = latest
		}
	}

	if end == 0 || start >= end {
		checkpointer.Logger.Error("Invalid start end formation", "Start", start, "End", end)
		return 0, 0, currentHeaderBlockNumber, errors.New("Invalid start end formation")
	}
	return
}

func getLastCheckpointStore() (uint64, uint64, bool) {
	var _checkpoint hmtypes.CheckpointBlockHeader
	resp, err := http.Get(lastCheckpointURL)
	if err != nil {
		pierLogger.Error("Unable to send request to get proposer", "Error", err)
		return _checkpoint.StartBlock, _checkpoint.EndBlock, false
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			pierLogger.Error("Unable to read data from response", "Error", err)
			return _checkpoint.StartBlock, _checkpoint.EndBlock, false
		}
		if err := json.Unmarshal(body, &_checkpoint); err != nil {
			pierLogger.Error("Error unmarshalling checkpoint", "error", err)
			return _checkpoint.StartBlock, _checkpoint.EndBlock, false
		}
		return _checkpoint.StartBlock, _checkpoint.EndBlock, true

	}
	return _checkpoint.StartBlock, _checkpoint.EndBlock, false
}
