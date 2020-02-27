package listener

import (
	"context"
	"encoding/json"
	"math/big"
	"strconv"

	ethereum "github.com/maticnetwork/bor"
	"github.com/maticnetwork/bor/accounts/abi"
	ethCommon "github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/bridge/setu/queue"
	"github.com/maticnetwork/heimdall/helper"
)

// RootChainListener - Listens to and process events from rootchain
type RootChainListener struct {
	BaseListener
	// ABIs
	abis []*abi.ABI
}

const (
	lastRootBlockKey = "rootchain-last-block" // storage key
)

// NewRootChainListener - constructor func
func NewRootChainListener() *RootChainListener {
	contractCaller, err := helper.NewContractCaller()
	if err != nil {
		panic(err)
	}
	abis := []*abi.ABI{
		&contractCaller.RootChainABI,
		&contractCaller.StateSenderABI,
		&contractCaller.StakingInfoABI,
	}
	rootchainListener := &RootChainListener{
		abis: abis,
	}
	return rootchainListener
}

// Start starts new block subscription
func (rl *RootChainListener) Start() error {
	rl.Logger.Info("Starting")
	// create cancellable context
	ctx, cancelSubscription := context.WithCancel(context.Background())
	rl.cancelSubscription = cancelSubscription

	// create cancellable context
	headerCtx, cancelHeaderProcess := context.WithCancel(context.Background())
	rl.cancelHeaderProcess = cancelHeaderProcess

	// start header process
	go rl.StartHeaderProcess(headerCtx)

	// subscribe to new head
	subscription, err := rl.contractConnector.MainChainClient.SubscribeNewHead(ctx, rl.HeaderChannel)
	if err != nil {
		// start go routine to poll for new header using client object
		rl.Logger.Info("Start polling for rootchain header blocks", "pollInterval", helper.GetConfig().SyncerPollInterval)
		go rl.StartPolling(ctx, helper.GetConfig().SyncerPollInterval)
	} else {
		// start go routine to listen new header using subscription
		go rl.StartSubscription(ctx, subscription)
	}

	// subscribed to new head
	rl.Logger.Info("Subscribed to new head")

	return nil
}

// ProcessHeader - process headerblock from rootchain
func (rl *RootChainListener) ProcessHeader(newHeader *types.Header) {
	rl.Logger.Info("Received Headerblock", "blockNumber", newHeader.Number)
	latestNumber := newHeader.Number

	// confirmation
	confirmationBlocks := big.NewInt(0).SetUint64(helper.GetConfig().ConfirmationBlocks)
	confirmationBlocks = confirmationBlocks.Add(confirmationBlocks, big.NewInt(1))
	if latestNumber.Uint64() > confirmationBlocks.Uint64() {
		latestNumber = latestNumber.Sub(latestNumber, confirmationBlocks)
	}

	// default fromBlock
	fromBlock := latestNumber
	// get last block from storage
	hasLastBlock, _ := rl.storageClient.Has([]byte(lastRootBlockKey), nil)
	if hasLastBlock {
		lastBlockBytes, err := rl.storageClient.Get([]byte(lastRootBlockKey), nil)
		if err != nil {
			rl.Logger.Info("Error while fetching last block bytes from storage", "error", err)
			return
		}
		rl.Logger.Debug("Got last block from bridge storage", "lastBlock", string(lastBlockBytes))
		if result, err := strconv.ParseUint(string(lastBlockBytes), 10, 64); err == nil {
			if result >= newHeader.Number.Uint64() {
				return
			}
			fromBlock = big.NewInt(0).SetUint64(result + 1)
		}
	}

	// to block
	toBlock := latestNumber
	// set diff
	if toBlock.Uint64() < fromBlock.Uint64() {
		fromBlock = toBlock
	}

	// set last block to storage
	rl.storageClient.Put([]byte(lastRootBlockKey), []byte(toBlock.String()), nil)

	// query log
	rl.Logger.Info("Query event logs", "fromBlock", fromBlock, "toBlock", toBlock)
	rl.queryAndBroadcastEvents(fromBlock, toBlock)
}

func (rl *RootChainListener) queryAndBroadcastEvents(fromBlock *big.Int, toBlock *big.Int) {
	// draft a query
	query := ethereum.FilterQuery{FromBlock: fromBlock, ToBlock: toBlock, Addresses: []ethCommon.Address{helper.GetRootChainAddress(), helper.GetStakingInfoAddress(), helper.GetStateSenderAddress()}}
	// get logs from rootchain by filter
	logs, err := rl.contractConnector.MainChainClient.FilterLogs(context.Background(), query)
	if err != nil {
		rl.Logger.Error("Error while filtering logs", "error", err)
		return
	} else if len(logs) > 0 {
		rl.Logger.Debug("New logs found", "numberOfLogs", len(logs))
	}

	// process filtered log
	for _, vLog := range logs {
		topic := vLog.Topics[0].Bytes()
		for _, abiObject := range rl.abis {
			selectedEvent := helper.EventByID(abiObject, topic)
			if selectedEvent != nil {
				rl.Logger.Debug("ReceivedEvent", "eventname", selectedEvent.Name)
				switch selectedEvent.Name {

				case "NewHeaderBlock":
					logBytes, _ := json.Marshal(vLog)
					if err := rl.queueConnector.PublishMsg(logBytes, queue.CheckpointQueueRoute, rl.String(), selectedEvent.Name); err != nil {
						rl.Logger.Error("Error publishing msg to checkpoint queue", "error", err)
					}

				case "StakeUpdate", "SignerChange", "UnstakeInit", "ReStaked":
					logBytes, _ := json.Marshal(vLog)
					if err := rl.queueConnector.PublishMsg(logBytes, queue.StakingQueueRoute, rl.String(), selectedEvent.Name); err != nil {
						rl.Logger.Error("Error publishing msg to staking queue", "error", err)
					}

				case "StateSynced":
					logBytes, _ := json.Marshal(vLog)
					if err := rl.queueConnector.PublishMsg(logBytes, queue.ClerkQueueRoute, rl.String(), selectedEvent.Name); err != nil {
						rl.Logger.Error("Error publishing msg to clerk queue", "error", err)
					}

				case "TopUpFee":
					logBytes, _ := json.Marshal(vLog)
					if err := rl.queueConnector.PublishMsg(logBytes, queue.FeeQueueRoute, rl.String(), selectedEvent.Name); err != nil {
						rl.Logger.Error("Error publishing msg to topup queue", "error", err)
					}

				}
			}
		}
	}
}
