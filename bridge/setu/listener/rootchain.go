package listener

import (
	"context"
	"encoding/json"
	"math/big"
	"strconv"
	"time"

	"github.com/RichardKnop/machinery/v1/tasks"
	ethereum "github.com/maticnetwork/bor"
	"github.com/maticnetwork/bor/accounts/abi"
	ethCommon "github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
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
	rl.queryAndBroadcastEvents(fromBlock, toBlock)
}

func (rl *RootChainListener) queryAndBroadcastEvents(fromBlock *big.Int, toBlock *big.Int) {
	rl.Logger.Info("Query rootchain event logs", "fromBlock", fromBlock, "toBlock", toBlock)

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
			logBytes, _ := json.Marshal(vLog)
			if selectedEvent != nil {
				rl.Logger.Debug("ReceivedEvent", "eventname", selectedEvent.Name)
				switch selectedEvent.Name {
				case "NewHeaderBlock":
					rl.sendTask("sendCheckpointAckToHeimdall", selectedEvent.Name, logBytes)
				case "StakeUpdate":
					rl.sendTask("sendStakeUpdateToHeimdall", selectedEvent.Name, logBytes)
				case "SignerChange":
					rl.sendTask("sendSignerChangeToHeimdall", selectedEvent.Name, logBytes)
				case "UnstakeInit":
					rl.sendTask("sendUnstakeInitToHeimdall", selectedEvent.Name, logBytes)
				case "ReStaked":
					rl.sendTask("sendReStakedToHeimdall", selectedEvent.Name, logBytes)
				case "StateSynced":
					rl.sendTask("sendStateSyncedToHeimdall", selectedEvent.Name, logBytes)
				case "TopUpFee":
					if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx); isCurrentValidator {
						rl.sendTaskWithDelay("sendTopUpFeeToHeimdall", selectedEvent.Name, logBytes, delay)
					} else {
						rl.Logger.Info("i am not present in current validatorset. ignore sending topup task")
					}
				}
			}
		}
	}
}

func (rl *RootChainListener) sendTask(taskName string, eventName string, logBytes []byte) {
	signature := &tasks.Signature{
		Name: taskName,
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: eventName,
			},
			{
				Type:  "string",
				Value: string(logBytes),
			},
		},
	}
	signature.RetryCount = 3

	_, err := rl.queueConnector.Server.SendTask(signature)
	if err != nil {
		rl.Logger.Error("Error sending checkpoint task")
	}
}

func (rl *RootChainListener) sendTaskWithDelay(taskName string, eventName string, logBytes []byte, delay time.Duration) {
	signature := &tasks.Signature{
		Name: taskName,
		Args: []tasks.Arg{
			{
				Type:  "string",
				Value: eventName,
			},
			{
				Type:  "string",
				Value: string(logBytes),
			},
		},
	}
	signature.RetryCount = 3

	// add delay for task so that multiple validators won't send same transaction at same time
	eta := time.Now().Add(delay)
	signature.ETA = &eta
	rl.Logger.Info("sending topup task", "currenttime", time.Now(), "delaytime", eta)
	_, err := rl.queueConnector.Server.SendTask(signature)
	if err != nil {
		rl.Logger.Error("Error while sending task")
	}
}
