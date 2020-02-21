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
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/sethu/util"
)

// RootChainListener syncs validators and checkpoints
type RootChainListener struct {
	BaseListener
	// ABIs
	abis []*abi.ABI
}

const (
	headerEvent      = "NewHeaderBlock"
	stakeInitEvent   = "Staked"
	unstakeInitEvent = "UnstakeInit"
	signerChange     = "SignerChange"

	lastBlockKey = "last-block" // storage key
)

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
	rl.Logger.Info("Starting listener", "name", rl.String())
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
		go rl.StartPolling(ctx, helper.GetConfig().SyncerPollInterval)
	} else {
		// start go routine to listen new header using subscription
		go rl.StartSubscription(ctx, subscription)
	}

	// subscribed to new head
	rl.Logger.Info("Subscribed to new head")

	return nil
}

// func (bl *RootChainListener) ProcessHeader(newHeader *types.Header) {
// 	bl.Logger.Info("Received Headerblock from Rootchain", "header", newHeader)

// 	bl.queueConnector.PublishMsg([]byte("StakingMsg"), util.StakingQueueRoute, bl.String())
// }

func (rl *RootChainListener) ProcessHeader(newHeader *types.Header) {
	rl.Logger.Info("New block detected", "blockNumber", newHeader.Number)
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
	hasLastBlock, _ := rl.storageClient.Has([]byte(lastBlockKey), nil)
	if hasLastBlock {
		lastBlockBytes, err := rl.storageClient.Get([]byte(lastBlockKey), nil)
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

	// debug log
	rl.Logger.Debug("Processing header", "fromBlock", fromBlock, "toBlock", toBlock)

	// set diff
	if toBlock.Uint64() < fromBlock.Uint64() {
		fromBlock = toBlock
	}

	// set last block to storage
	rl.storageClient.Put([]byte(lastBlockKey), []byte(toBlock.String()), nil)

	fromBlock = big.NewInt(0).SetInt64(27)
	toBlock = big.NewInt(0).SetInt64(15827)

	// log
	rl.Logger.Info("Querying event logs", "fromBlock", fromBlock, "toBlock", toBlock)

	// draft a query
	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: []ethCommon.Address{
			helper.GetRootChainAddress(),
			helper.GetStakingInfoAddress(),
			helper.GetStateSenderAddress(),
		},
	}

	// get all logs
	logs, err := rl.contractConnector.MainChainClient.FilterLogs(context.Background(), query)
	if err != nil {
		rl.Logger.Error("Error while filtering logs from syncer", "error", err)
		return
	} else if len(logs) > 0 {
		rl.Logger.Debug("New logs found", "numberOfLogs", len(logs))
	}

	// log
	for _, vLog := range logs {
		topic := vLog.Topics[0].Bytes()
		for _, abiObject := range rl.abis {
			selectedEvent := helper.EventByID(abiObject, topic)
			if selectedEvent != nil {
				rl.Logger.Debug("selectedEvent ", " event name -", selectedEvent.Name)
				switch selectedEvent.Name {
				case "NewHeaderBlock":
					logBytes, _ := json.Marshal(vLog)
					rl.queueConnector.PublishMsg(logBytes, util.CheckpointQueueRoute, rl.String(), selectedEvent.Name)
					break
				}
			}
		}
	}
}
