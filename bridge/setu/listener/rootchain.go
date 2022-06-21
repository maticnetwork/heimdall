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
	"github.com/maticnetwork/heimdall/contracts/statesender"
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
	stateSenderAbi *abi.ABI
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

	return &RootChainListener{
		abis:           abis,
		stakingInfoAbi: &contractCaller.StakingInfoABI,
		stateSenderAbi: &contractCaller.StateSenderABI,
	}
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
	rl.Logger.Info("Subscribed to new head")

	// Start self-healing process
	go rl.startSelfHealing(ctx)

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

	// Prepare blocks range
	toBlock := latestNumber
	if toBlock.Cmp(fromBlock) == -1 {
		fromBlock = toBlock
	}

	// Set last block to storage
	if err = rl.storageClient.Put([]byte(lastRootBlockKey), []byte(toBlock.String()), nil); err != nil {
		rl.Logger.Error("rl.storageClient.Put", "Error", err)
	}

	// Handle events
	rl.queryAndBroadcastEvents(rootchainContext, fromBlock, toBlock)
}

// queryAndBroadcastEvents fetches supported events from the rootchain and handles all of them
func (rl *RootChainListener) queryAndBroadcastEvents(rootchainContext *RootChainListenerContext, fromBlock *big.Int, toBlock *big.Int) {
	rl.Logger.Info("Query rootchain event logs", "fromBlock", fromBlock, "toBlock", toBlock)

	ctx, cancel := context.WithTimeout(context.Background(), rl.contractConnector.MainChainTimeout)
	defer cancel()

	// get chain params
	chainParams := rootchainContext.ChainmanagerParams.ChainParams

	// Fetch events from the rootchain
	logs, err := rl.contractConnector.MainChainClient.FilterLogs(ctx, ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: []ethCommon.Address{
			chainParams.RootChainAddress.EthAddress(),
			chainParams.StakingInfoAddress.EthAddress(),
			chainParams.StateSenderAddress.EthAddress(),
		},
	})
	if err != nil {
		rl.Logger.Error("Error while filtering logs", "error", err)
		return
	} else if len(logs) > 0 {
		rl.Logger.Debug("New logs found", "numberOfLogs", len(logs))
	}

	pubkey := helper.GetPubKey()
	pubkeyBytes := pubkey[1:]

	// Process filtered log
	for _, vLog := range logs {
		topic := vLog.Topics[0].Bytes()
		for _, abiObject := range rl.abis {
			selectedEvent := helper.EventByID(abiObject, topic)
			if selectedEvent == nil {
				continue
			}
			logBytes, _ := json.Marshal(vLog)
			if selectedEvent != nil {
				rl.Logger.Debug("ReceivedEvent", "eventname", selectedEvent.Name)
				switch selectedEvent.Name {
				case "NewHeaderBlock":
					if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx, selectedEvent); isCurrentValidator {
						rl.SendTaskWithDelay("sendCheckpointAckToHeimdall", selectedEvent.Name, logBytes, delay, selectedEvent)
					}
				case "Staked":
					event := new(stakinginfo.StakinginfoStaked)
					if err := helper.UnpackLog(rl.stakingInfoAbi, event, selectedEvent.Name, &vLog); err != nil {
						rl.Logger.Error("Error while parsing event", "name", selectedEvent.Name, "error", err)
					}
					if bytes.Equal(event.SignerPubkey, pubkeyBytes) {
						// topup has to be processed first before validator join. so adding delay.
						delay := util.TaskDelayBetweenEachVal
						rl.SendTaskWithDelay("sendValidatorJoinToHeimdall", selectedEvent.Name, logBytes, delay, event)
					} else if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx, event); isCurrentValidator {
						// topup has to be processed first before validator join. so adding delay.
						delay = delay + util.TaskDelayBetweenEachVal
						rl.SendTaskWithDelay("sendValidatorJoinToHeimdall", selectedEvent.Name, logBytes, delay, event)
					}

				case "StakeUpdate":
					event := new(stakinginfo.StakinginfoStakeUpdate)
					if err := helper.UnpackLog(rl.stakingInfoAbi, event, selectedEvent.Name, &vLog); err != nil {
						rl.Logger.Error("Error while parsing event", "name", selectedEvent.Name, "error", err)
					}
					if util.IsEventSender(rl.cliCtx, event.ValidatorId.Uint64()) {
						rl.SendTaskWithDelay("sendStakeUpdateToHeimdall", selectedEvent.Name, logBytes, 0, event)
					} else if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx, event); isCurrentValidator {
						rl.SendTaskWithDelay("sendStakeUpdateToHeimdall", selectedEvent.Name, logBytes, delay, event)
					}

				case "SignerChange":
					event := new(stakinginfo.StakinginfoSignerChange)
					if err := helper.UnpackLog(rl.stakingInfoAbi, event, selectedEvent.Name, &vLog); err != nil {
						rl.Logger.Error("Error while parsing event", "name", selectedEvent.Name, "error", err)
					}
					if bytes.Equal(event.SignerPubkey, pubkeyBytes) {
						rl.SendTaskWithDelay("sendSignerChangeToHeimdall", selectedEvent.Name, logBytes, 0, event)
					} else if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx, event); isCurrentValidator {
						rl.SendTaskWithDelay("sendSignerChangeToHeimdall", selectedEvent.Name, logBytes, delay, event)
					}

				case "UnstakeInit":
					event := new(stakinginfo.StakinginfoUnstakeInit)
					if err := helper.UnpackLog(rl.stakingInfoAbi, event, selectedEvent.Name, &vLog); err != nil {
						rl.Logger.Error("Error while parsing event", "name", selectedEvent.Name, "error", err)
					}
					if util.IsEventSender(rl.cliCtx, event.ValidatorId.Uint64()) {
						rl.SendTaskWithDelay("sendUnstakeInitToHeimdall", selectedEvent.Name, logBytes, 0, event)
					} else if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx, event); isCurrentValidator {
						rl.SendTaskWithDelay("sendUnstakeInitToHeimdall", selectedEvent.Name, logBytes, delay, event)
					}

				case "StateSynced":
					event := new(statesender.StatesenderStateSynced)
					if err := helper.UnpackLog(rl.stateSenderAbi, event, selectedEvent.Name, &vLog); err != nil {
						rl.Logger.Error("Error while parsing event", "name", selectedEvent.Name, "error", err)
					}
					rl.Logger.Info("StateSyncedEvent: detected", "stateSyncId", event.Id)
					if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx, event); isCurrentValidator {
						rl.SendTaskWithDelay("sendStateSyncedToHeimdall", selectedEvent.Name, logBytes, delay, event)
					}

				case "TopUpFee":
					event := new(stakinginfo.StakinginfoTopUpFee)
					if err := helper.UnpackLog(rl.stakingInfoAbi, event, selectedEvent.Name, &vLog); err != nil {
						rl.Logger.Error("Error while parsing event", "name", selectedEvent.Name, "error", err)
					}
					if bytes.Equal(event.User.Bytes(), helper.GetAddress()) {
						rl.SendTaskWithDelay("sendTopUpFeeToHeimdall", selectedEvent.Name, logBytes, 0, event)
					} else if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx, event); isCurrentValidator {
						rl.SendTaskWithDelay("sendTopUpFeeToHeimdall", selectedEvent.Name, logBytes, delay, event)
					}

				case "Slashed":
					if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx, selectedEvent); isCurrentValidator {
						rl.SendTaskWithDelay("sendTickAckToHeimdall", selectedEvent.Name, logBytes, delay, selectedEvent)
					}

				case "UnJailed":
					event := new(stakinginfo.StakinginfoUnJailed)
					if err := helper.UnpackLog(rl.stakingInfoAbi, event, selectedEvent.Name, &vLog); err != nil {
						rl.Logger.Error("Error while parsing event", "name", selectedEvent.Name, "error", err)
					}
					if util.IsEventSender(rl.cliCtx, event.ValidatorId.Uint64()) {
						rl.SendTaskWithDelay("sendUnjailToHeimdall", selectedEvent.Name, logBytes, 0, event)
					} else if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx, event); isCurrentValidator {
						rl.SendTaskWithDelay("sendUnjailToHeimdall", selectedEvent.Name, logBytes, delay, event)
					}
				}
				rl.handleLog(vLog, selectedEvent, pubkeyBytes)
			}
		}
	}
}

// handleLog handles the given log
func (rl *RootChainListener) handleLog(vLog types.Log, selectedEvent *abi.Event, pubkeyBytes []byte) {
	logBytes, _ := json.Marshal(vLog)

	rl.Logger.Debug("ReceivedEvent", "eventname", selectedEvent.Name)
	switch selectedEvent.Name {
	case "NewHeaderBlock":
		if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx, selectedEvent); isCurrentValidator {
			rl.SendTaskWithDelay("sendCheckpointAckToHeimdall", selectedEvent.Name, logBytes, delay, selectedEvent)
		}
	case "Staked":
		event := new(stakinginfo.StakinginfoStaked)
		if err := helper.UnpackLog(rl.stakingInfoAbi, event, selectedEvent.Name, &vLog); err != nil {
			rl.Logger.Error("Error while parsing event", "name", selectedEvent.Name, "error", err)
		}
		if bytes.Equal(event.SignerPubkey, pubkeyBytes) {
			// topup has to be processed first before validator join. so adding delay.
			delay := util.TaskDelayBetweenEachVal
			rl.SendTaskWithDelay("sendValidatorJoinToHeimdall", selectedEvent.Name, logBytes, delay, event)
		} else if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx, event); isCurrentValidator {
			// topup has to be processed first before validator join. so adding delay.
			delay = delay + util.TaskDelayBetweenEachVal
			rl.SendTaskWithDelay("sendValidatorJoinToHeimdall", selectedEvent.Name, logBytes, delay, event)
		}

	case "StakeUpdate":
		event := new(stakinginfo.StakinginfoStakeUpdate)
		if err := helper.UnpackLog(rl.stakingInfoAbi, event, selectedEvent.Name, &vLog); err != nil {
			rl.Logger.Error("Error while parsing event", "name", selectedEvent.Name, "error", err)
		}
		if util.IsEventSender(rl.cliCtx, event.ValidatorId.Uint64()) {
			rl.SendTaskWithDelay("sendStakeUpdateToHeimdall", selectedEvent.Name, logBytes, 0, event)
		} else if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx, event); isCurrentValidator {
			rl.SendTaskWithDelay("sendStakeUpdateToHeimdall", selectedEvent.Name, logBytes, delay, event)
		}

	case "SignerChange":
		event := new(stakinginfo.StakinginfoSignerChange)
		if err := helper.UnpackLog(rl.stakingInfoAbi, event, selectedEvent.Name, &vLog); err != nil {
			rl.Logger.Error("Error while parsing event", "name", selectedEvent.Name, "error", err)
		}
		if bytes.Equal(event.SignerPubkey, pubkeyBytes) {
			rl.SendTaskWithDelay("sendSignerChangeToHeimdall", selectedEvent.Name, logBytes, 0, event)
		} else if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx, event); isCurrentValidator {
			rl.SendTaskWithDelay("sendSignerChangeToHeimdall", selectedEvent.Name, logBytes, delay, event)
		}

	case "UnstakeInit":
		event := new(stakinginfo.StakinginfoUnstakeInit)
		if err := helper.UnpackLog(rl.stakingInfoAbi, event, selectedEvent.Name, &vLog); err != nil {
			rl.Logger.Error("Error while parsing event", "name", selectedEvent.Name, "error", err)
		}
		if util.IsEventSender(rl.cliCtx, event.ValidatorId.Uint64()) {
			rl.SendTaskWithDelay("sendUnstakeInitToHeimdall", selectedEvent.Name, logBytes, 0, event)
		} else if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx, event); isCurrentValidator {
			rl.SendTaskWithDelay("sendUnstakeInitToHeimdall", selectedEvent.Name, logBytes, delay, event)
		}

	case "StateSynced":
		event := new(statesender.StatesenderStateSynced)
		if err := helper.UnpackLog(rl.stateSenderAbi, event, selectedEvent.Name, &vLog); err != nil {
			rl.Logger.Error("Error while parsing event", "name", selectedEvent.Name, "error", err)
		}
		rl.Logger.Info("StateSyncedEvent: detected", "stateSyncId", event.Id)
		if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx, event); isCurrentValidator {
			rl.SendTaskWithDelay("sendStateSyncedToHeimdall", selectedEvent.Name, logBytes, delay, event)
		}

	case "TopUpFee":
		event := new(stakinginfo.StakinginfoTopUpFee)
		if err := helper.UnpackLog(rl.stakingInfoAbi, event, selectedEvent.Name, &vLog); err != nil {
			rl.Logger.Error("Error while parsing event", "name", selectedEvent.Name, "error", err)
		}
		if bytes.Equal(event.User.Bytes(), helper.GetAddress()) {
			rl.SendTaskWithDelay("sendTopUpFeeToHeimdall", selectedEvent.Name, logBytes, 0, event)
		} else if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx, event); isCurrentValidator {
			rl.SendTaskWithDelay("sendTopUpFeeToHeimdall", selectedEvent.Name, logBytes, delay, event)
		}

	case "Slashed":
		if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx, selectedEvent); isCurrentValidator {
			rl.SendTaskWithDelay("sendTickAckToHeimdall", selectedEvent.Name, logBytes, delay, selectedEvent)
		}

	case "UnJailed":
		event := new(stakinginfo.StakinginfoUnJailed)
		if err := helper.UnpackLog(rl.stakingInfoAbi, event, selectedEvent.Name, &vLog); err != nil {
			rl.Logger.Error("Error while parsing event", "name", selectedEvent.Name, "error", err)
		}
		if util.IsEventSender(rl.cliCtx, event.ValidatorId.Uint64()) {
			rl.SendTaskWithDelay("sendUnjailToHeimdall", selectedEvent.Name, logBytes, 0, event)
		} else if isCurrentValidator, delay := util.CalculateTaskDelay(rl.cliCtx, event); isCurrentValidator {
			rl.SendTaskWithDelay("sendUnjailToHeimdall", selectedEvent.Name, logBytes, delay, event)
		}
	}
}

func (rl *RootChainListener) SendTaskWithDelay(taskName string, eventName string, logBytes []byte, delay time.Duration, event interface{}) {
	defer util.LogElapsedTimeForStateSyncedEvent(event, "SendTaskWithDelay", time.Now())

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

// getRootChainContext returns the root chain context
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
