package listener

import (
	"bytes"
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
	chainmanagerTypes "github.com/maticnetwork/heimdall/chainmanager/types"
	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/helper"
)

// RootChainListenerContext root chain listener context
type RootChainListenerContext struct {
	ChainmanagerParams *chainmanagerTypes.Params
}

// RootChainListener - Listens to and process events from rootchain
type RootChainListener struct {
	BaseListener
	// ABIs
	abis []*abi.ABI

	stakingInfoAbi *abi.ABI
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
		abis:           abis,
		stakingInfoAbi: &contractCaller.StakingInfoABI,
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
	rl.Logger.Debug("New block detected", "blockNumber", newHeader.Number)

	// fetch context
	rootchainContext, err := rl.getRootChainContext()
	if err != nil {
		return
	}
	requiredConfirmations := rootchainContext.ChainmanagerParams.MainchainTxConfirmations
	latestNumber := newHeader.Number

	// confirmation
	confirmationBlocks := big.NewInt(0).SetUint64(requiredConfirmations)

	if latestNumber.Cmp(confirmationBlocks) <= 0 {
		rl.Logger.Error("Block number less than Confirmations required", "blockNumber", latestNumber.Uint64, "confirmationsRequired", confirmationBlocks.Uint64)
		return
	}
	latestNumber = latestNumber.Sub(latestNumber, confirmationBlocks)

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

	if toBlock.Cmp(fromBlock) == -1 {
		fromBlock = toBlock
	}

	// set last block to storage
	if err := rl.storageClient.Put([]byte(lastRootBlockKey), []byte(toBlock.String()), nil); err != nil {
		rl.Logger.Error("rl.storageClient.Put", "Error", err)
	}

	// query events
	rl.queryAndBroadcastEvents(rootchainContext, fromBlock, toBlock)
}

func (rl *RootChainListener) queryAndBroadcastEvents(rootchainContext *RootChainListenerContext, fromBlock *big.Int, toBlock *big.Int) {
	rl.Logger.Info("Query rootchain event logs", "fromBlock", fromBlock, "toBlock", toBlock)

	// current public key
	pubkey := helper.GetPubKey()
	pubkeyBytes := pubkey[1:]

	// get chain params
	chainParams := rootchainContext.ChainmanagerParams.ChainParams

	// draft a query
	query := ethereum.FilterQuery{FromBlock: fromBlock, ToBlock: toBlock, Addresses: []ethCommon.Address{
		chainParams.RootChainAddress.EthAddress(),
		chainParams.StakingInfoAddress.EthAddress(),
		chainParams.StateSenderAddress.EthAddress(),
	}}
	// get logs from rootchain by filter

	ethClientTimeout := rl.contractConnector.MainChainTimeout
	ctx, cancel := context.WithTimeout(context.Background(), ethClientTimeout)
	defer cancel()
	logs, err := rl.contractConnector.MainChainClient.FilterLogs(ctx, query)
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
					if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx); isCurrentValidator {
						rl.sendTaskWithDelay("sendCheckpointAckToHeimdall", selectedEvent.Name, logBytes, delay)
					}
				case "Staked":
					event := new(stakinginfo.StakinginfoStaked)
					if err := helper.UnpackLog(rl.stakingInfoAbi, event, selectedEvent.Name, &vLog); err != nil {
						rl.Logger.Error("Error while parsing event", "name", selectedEvent.Name, "error", err)
					}
					if bytes.Equal(event.SignerPubkey, pubkeyBytes) {
						// topup has to be processed first before validator join. so adding delay.
						delay := util.TaskDelayBetweenEachVal
						rl.sendTaskWithDelay("sendValidatorJoinToHeimdall", selectedEvent.Name, logBytes, delay)
					} else if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx); isCurrentValidator {
						// topup has to be processed first before validator join. so adding delay.
						delay = delay + util.TaskDelayBetweenEachVal
						rl.sendTaskWithDelay("sendValidatorJoinToHeimdall", selectedEvent.Name, logBytes, delay)
					}

				case "StakeUpdate":
					event := new(stakinginfo.StakinginfoStakeUpdate)
					if err := helper.UnpackLog(rl.stakingInfoAbi, event, selectedEvent.Name, &vLog); err != nil {
						rl.Logger.Error("Error while parsing event", "name", selectedEvent.Name, "error", err)
					}
					if util.IsEventSender(rl.cliCtx, event.ValidatorId.Uint64()) {
						rl.sendTaskWithDelay("sendStakeUpdateToHeimdall", selectedEvent.Name, logBytes, 0)
					} else if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx); isCurrentValidator {
						rl.sendTaskWithDelay("sendStakeUpdateToHeimdall", selectedEvent.Name, logBytes, delay)
					}

				case "SignerChange":
					event := new(stakinginfo.StakinginfoSignerChange)
					if err := helper.UnpackLog(rl.stakingInfoAbi, event, selectedEvent.Name, &vLog); err != nil {
						rl.Logger.Error("Error while parsing event", "name", selectedEvent.Name, "error", err)
					}
					if bytes.Equal(event.SignerPubkey, pubkeyBytes) {
						rl.sendTaskWithDelay("sendSignerChangeToHeimdall", selectedEvent.Name, logBytes, 0)
					} else if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx); isCurrentValidator {
						rl.sendTaskWithDelay("sendSignerChangeToHeimdall", selectedEvent.Name, logBytes, delay)
					}

				case "UnstakeInit":
					event := new(stakinginfo.StakinginfoUnstakeInit)
					if err := helper.UnpackLog(rl.stakingInfoAbi, event, selectedEvent.Name, &vLog); err != nil {
						rl.Logger.Error("Error while parsing event", "name", selectedEvent.Name, "error", err)
					}
					if util.IsEventSender(rl.cliCtx, event.ValidatorId.Uint64()) {
						rl.sendTaskWithDelay("sendUnstakeInitToHeimdall", selectedEvent.Name, logBytes, 0)
					} else if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx); isCurrentValidator {
						rl.sendTaskWithDelay("sendUnstakeInitToHeimdall", selectedEvent.Name, logBytes, delay)
					}

				case "StateSynced":
					if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx); isCurrentValidator {
						rl.sendTaskWithDelay("sendStateSyncedToHeimdall", selectedEvent.Name, logBytes, delay)
					}

				case "TopUpFee":
					event := new(stakinginfo.StakinginfoTopUpFee)
					if err := helper.UnpackLog(rl.stakingInfoAbi, event, selectedEvent.Name, &vLog); err != nil {
						rl.Logger.Error("Error while parsing event", "name", selectedEvent.Name, "error", err)
					}
					if bytes.Equal(event.User.Bytes(), helper.GetAddress()) {
						rl.sendTaskWithDelay("sendTopUpFeeToHeimdall", selectedEvent.Name, logBytes, 0)
					} else if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx); isCurrentValidator {
						rl.sendTaskWithDelay("sendTopUpFeeToHeimdall", selectedEvent.Name, logBytes, delay)
					}

				case "Slashed":
					if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx); isCurrentValidator {
						rl.sendTaskWithDelay("sendTickAckToHeimdall", selectedEvent.Name, logBytes, delay)
					}

				case "UnJailed":
					event := new(stakinginfo.StakinginfoUnJailed)
					if err := helper.UnpackLog(rl.stakingInfoAbi, event, selectedEvent.Name, &vLog); err != nil {
						rl.Logger.Error("Error while parsing event", "name", selectedEvent.Name, "error", err)
					}
					if util.IsEventSender(rl.cliCtx, event.ValidatorId.Uint64()) {
						rl.sendTaskWithDelay("sendUnjailToHeimdall", selectedEvent.Name, logBytes, 0)
					} else if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx); isCurrentValidator {
						rl.sendTaskWithDelay("sendUnjailToHeimdall", selectedEvent.Name, logBytes, delay)
					}
				}
			}
		}
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
	rl.Logger.Info("Sending task", "taskName", taskName, "currentTime", time.Now(), "delayTime", eta)
	_, err := rl.queueConnector.Server.SendTask(signature)
	if err != nil {
		rl.Logger.Error("Error sending task", "taskName", taskName, "error", err)
	}
}

//
// utils
//

func (rl *RootChainListener) getRootChainContext() (*RootChainListenerContext, error) {
	chainmanagerParams, err := util.GetChainmanagerParams(rl.cliCtx)
	if err != nil {
		rl.Logger.Error("Error while fetching chain manager params", "error", err)
		return nil, err
	}

	return &RootChainListenerContext{
		ChainmanagerParams: chainmanagerParams,
	}, nil
}
