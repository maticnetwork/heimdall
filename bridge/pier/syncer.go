package pier

import (
	"bytes"
	"context"
	"encoding/hex"
	"math/big"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	ethereum "github.com/maticnetwork/bor"
	"github.com/maticnetwork/bor/accounts/abi"
	ethCommon "github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/bor/core/types"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tendermint/libs/common"
	httpClient "github.com/tendermint/tendermint/rpc/client"

	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	clerkTypes "github.com/maticnetwork/heimdall/clerk/types"
	"github.com/maticnetwork/heimdall/contracts/delegationmanager"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/contracts/stakemanager"
	"github.com/maticnetwork/heimdall/contracts/statesender"
	"github.com/maticnetwork/heimdall/helper"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

const (
	headerEvent      = "NewHeaderBlock"
	stakeInitEvent   = "Staked"
	unstakeInitEvent = "UnstakeInit"
	signerChange     = "SignerChange"

	lastBlockKey = "last-block" // storage key
)

// Syncer syncs validators and checkpoints
type Syncer struct {
	// Base service
	common.BaseService

	// storage client
	storageClient *leveldb.DB

	// contract caller
	contractConnector helper.ContractCaller

	// ABIs
	abis []*abi.ABI

	// header channel
	HeaderChannel chan *types.Header
	// cancel function for poll/subscription
	cancelSubscription context.CancelFunc

	// header listener subscription
	cancelHeaderProcess context.CancelFunc

	// cli context
	cliCtx cliContext.CLIContext

	// queue connector
	queueConnector *QueueConnector

	// http client to subscribe to
	httpClient *httpClient.HTTP
}

// NewSyncer returns new service object for syncing events
func NewSyncer(cdc *codec.Codec, queueConnector *QueueConnector, httpClient *httpClient.HTTP) *Syncer {
	// create logger
	logger := Logger.With("module", ChainSyncer)

	contractCaller, err := helper.NewContractCaller()
	if err != nil {
		logger.Error("Error while getting root chain instance", "error", err)
		panic(err)
	}

	abis := []*abi.ABI{
		&contractCaller.RootChainABI,
		&contractCaller.StakeManagerABI,
		&contractCaller.StateSenderABI,
		&contractCaller.DelegationManagerABI,
	}

	cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	cliCtx.BroadcastMode = client.BroadcastAsync
	cliCtx.TrustNode = true

	// creating syncer object
	syncer := &Syncer{
		storageClient: getBridgeDBInstance(viper.GetString(BridgeDBFlag)),

		cliCtx:            cliCtx,
		queueConnector:    queueConnector,
		httpClient:        httpClient,
		contractConnector: contractCaller,

		abis:          abis,
		HeaderChannel: make(chan *types.Header),
	}

	syncer.BaseService = *common.NewBaseService(logger, ChainSyncer, syncer)
	return syncer
}

// startHeaderProcess starts header process when they get new header
func (syncer *Syncer) startHeaderProcess(ctx context.Context) {
	for {
		select {
		case newHeader := <-syncer.HeaderChannel:
			syncer.processHeader(newHeader)
		case <-ctx.Done():
			return
		}
	}
}

// OnStart starts new block subscription
func (syncer *Syncer) OnStart() error {
	syncer.BaseService.OnStart() // Always call the overridden method.

	// create cancellable context
	ctx, cancelSubscription := context.WithCancel(context.Background())
	syncer.cancelSubscription = cancelSubscription

	// create cancellable context
	headerCtx, cancelHeaderProcess := context.WithCancel(context.Background())
	syncer.cancelHeaderProcess = cancelHeaderProcess

	// start header process
	go syncer.startHeaderProcess(headerCtx)

	// subscribe to new head
	subscription, err := syncer.contractConnector.MainChainClient.SubscribeNewHead(ctx, syncer.HeaderChannel)
	if err != nil {
		// start go routine to poll for new header using client object
		go syncer.startPolling(ctx, helper.GetConfig().SyncerPollInterval)
	} else {
		// start go routine to listen new header using subscription
		go syncer.startSubscription(ctx, subscription)
	}

	// subscribed to new head
	syncer.Logger.Debug("Subscribed to new head")

	return nil
}

// OnStop stops all necessary go routines
func (syncer *Syncer) OnStop() {
	syncer.BaseService.OnStop() // Always call the overridden method.

	// close db
	closeBridgeDBInstance()

	// cancel subscription if any
	syncer.cancelSubscription()

	// cancel header process
	syncer.cancelHeaderProcess()
}

// startPolling starts polling
func (syncer *Syncer) startPolling(ctx context.Context, pollInterval time.Duration) {
	// How often to fire the passed in function in second
	interval := pollInterval

	// Setup the ticket and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)

	// start listening
	for {
		select {
		case <-ticker.C:
			header, err := syncer.contractConnector.MainChainClient.HeaderByNumber(ctx, nil)
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

func (syncer *Syncer) startSubscription(ctx context.Context, subscription ethereum.Subscription) {
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

func (syncer *Syncer) processHeader(newHeader *types.Header) {
	syncer.Logger.Debug("New block detected", "blockNumber", newHeader.Number)
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
	hasLastBlock, _ := syncer.storageClient.Has([]byte(lastBlockKey), nil)
	if hasLastBlock {
		lastBlockBytes, err := syncer.storageClient.Get([]byte(lastBlockKey), nil)
		if err != nil {
			syncer.Logger.Info("Error while fetching last block bytes from storage", "error", err)
			return
		}

		syncer.Logger.Debug("Got last block from bridge storage", "lastBlock", string(lastBlockBytes))
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
	syncer.Logger.Debug("Processing header", "fromBlock", fromBlock, "toBlock", toBlock)

	// set diff
	if toBlock.Uint64() < fromBlock.Uint64() {
		fromBlock = toBlock
	}

	// set last block to storage
	syncer.storageClient.Put([]byte(lastBlockKey), []byte(toBlock.String()), nil)

	// log
	syncer.Logger.Info("Querying event logs", "fromBlock", fromBlock, "toBlock", toBlock)

	// draft a query
	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: []ethCommon.Address{
			helper.GetRootChainAddress(),
			helper.GetStakeManagerAddress(),
			helper.GetStateSenderAddress(),
			helper.GetDelegationManagerAddress(),
		},
	}

	// get all logs
	logs, err := syncer.contractConnector.MainChainClient.FilterLogs(context.Background(), query)
	if err != nil {
		syncer.Logger.Error("Error while filtering logs from syncer", "error", err)
		return
	} else if len(logs) > 0 {
		syncer.Logger.Debug("New logs found", "numberOfLogs", len(logs))
	}

	// log
	for _, vLog := range logs {
		topic := vLog.Topics[0].Bytes()
		for _, abiObject := range syncer.abis {
			selectedEvent := EventByID(abiObject, topic)
			if selectedEvent != nil {
				syncer.Logger.Debug("selectedEvent ", " event name -", selectedEvent.Name)
				switch selectedEvent.Name {
				case "NewHeaderBlock":
					syncer.processCheckpointEvent(selectedEvent.Name, abiObject, &vLog)
				case "Staked":
					syncer.processStakedEvent(selectedEvent.Name, abiObject, &vLog)
				case "UnstakeInit":
					syncer.processUnstakeInitEvent(selectedEvent.Name, abiObject, &vLog)
				case "StakeUpdate":
					syncer.processStakeUpdateEvent(selectedEvent.Name, abiObject, &vLog)
				case "SignerChange":
					syncer.processSignerChangeEvent(selectedEvent.Name, abiObject, &vLog)
				case "ReStaked":
					syncer.processReStakedEvent(selectedEvent.Name, abiObject, &vLog)
				case "Jailed":
					syncer.processJailedEvent(selectedEvent.Name, abiObject, &vLog)
				case "StateSynced":
					syncer.processStateSyncedEvent(selectedEvent.Name, abiObject, &vLog)
				case "Bonding":
					syncer.processDelegatorBondEvent(selectedEvent.Name, abiObject, &vLog)
				case "UnBonding":
					syncer.processDelegatorUnBondEvent(selectedEvent.Name, abiObject, &vLog)
					// case "Withdraw":
					// 	syncer.processWithdrawEvent(selectedEvent.Name, abiObject, &vLog)
				}
				break
			}
		}
	}
}

func (syncer *Syncer) processCheckpointEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {
	event := new(rootchain.RootchainNewHeaderBlock)
	if err := helper.UnpackLog(abiObject, event, eventName, vLog); err != nil {
		logEventParseError(syncer.Logger, eventName, err)
	} else {
		syncer.Logger.Info(
			"New event found",
			"event", eventName,
			"start", event.Start,
			"end", event.End,
			"reward", event.Reward,
			"root", "0x"+hex.EncodeToString(event.Root[:]),
			"proposer", event.Proposer.Hex(),
			"headerNumber", event.HeaderBlockId,
		)

		// create msg checkpoint ack message
		msg := checkpointTypes.NewMsgCheckpointAck(helper.GetFromAddress(syncer.cliCtx), event.HeaderBlockId.Uint64(), hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()), uint64(vLog.Index))
		syncer.queueConnector.BroadcastToHeimdall(msg)
	}
}

func (syncer *Syncer) processStakedEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {
	event := new(stakemanager.StakemanagerStaked)
	if err := helper.UnpackLog(abiObject, event, eventName, vLog); err != nil {
		logEventParseError(syncer.Logger, eventName, err)
	} else {
		syncer.Logger.Debug(
			"New event found",
			"event", eventName,
			"validator", event.User.Hex(),
			"ID", event.ValidatorId,
			"activatonEpoch", event.ActivationEpoch,
			"amount", event.Amount,
		)

		// compare user to get address
		if isEventSender(syncer.cliCtx, event.ValidatorId.Uint64()) {
			pubkey := helper.GetPubKey()
			msg := stakingTypes.NewMsgValidatorJoin(
				hmTypes.BytesToHeimdallAddress(event.User.Bytes()),
				event.ValidatorId.Uint64(),
				hmTypes.NewPubKey(pubkey[:]),
				hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				uint64(vLog.Index),
			)

			// process staked
			syncer.queueConnector.BroadcastToHeimdall(msg)
		}
	}
}

func (syncer *Syncer) processUnstakeInitEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {
	event := new(stakemanager.StakemanagerUnstakeInit)
	if err := helper.UnpackLog(abiObject, event, eventName, vLog); err != nil {
		logEventParseError(syncer.Logger, eventName, err)
	} else {
		syncer.Logger.Debug(
			"New event found",
			"event", eventName,
			"validator", event.User.Hex(),
			"validatorID", event.ValidatorId,
			"deactivatonEpoch", event.DeactivationEpoch,
			"amount", event.Amount,
		)

		// msg validator exit
		if isEventSender(syncer.cliCtx, event.ValidatorId.Uint64()) {
			msg := stakingTypes.NewMsgValidatorExit(
				hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
				event.ValidatorId.Uint64(),
				hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				uint64(vLog.Index),
			)

			// broadcast heimdall
			syncer.queueConnector.BroadcastToHeimdall(msg)
		}
	}
}

func (syncer *Syncer) processStakeUpdateEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {
	event := new(stakemanager.StakemanagerStakeUpdate)
	if err := helper.UnpackLog(abiObject, event, eventName, vLog); err != nil {
		logEventParseError(syncer.Logger, eventName, err)
	} else {
		syncer.Logger.Debug(
			"New event found",
			"event", eventName,
			"validatorID", event.ValidatorId,
			"oldAmount", event.OldAmount,
			"newAmount", event.NewAmount,
		)

		// msg validator exit
		if isEventSender(syncer.cliCtx, event.ValidatorId.Uint64()) {
			msg := stakingTypes.NewMsgStakeUpdate(
				hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
				event.ValidatorId.Uint64(),
				hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				uint64(vLog.Index),
			)

			// broadcast heimdall
			syncer.queueConnector.BroadcastToHeimdall(msg)
		}
	}
}

func (syncer *Syncer) processSignerChangeEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {
	event := new(stakemanager.StakemanagerSignerChange)
	if err := helper.UnpackLog(abiObject, event, eventName, vLog); err != nil {
		logEventParseError(syncer.Logger, eventName, err)
	} else {
		syncer.Logger.Debug(
			"New event found",
			"event", eventName,
			"validatorID", event.ValidatorId,
			"newSigner", event.NewSigner.Hex(),
			"oldSigner", event.OldSigner.Hex(),
		)

		// signer change
		if bytes.Compare(event.NewSigner.Bytes(), helper.GetAddress()) == 0 {
			pubkey := helper.GetPubKey()
			msg := stakingTypes.NewMsgSignerUpdate(
				hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
				event.ValidatorId.Uint64(),
				hmTypes.NewPubKey(pubkey[:]),
				hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				uint64(vLog.Index),
			)

			// process signer update
			syncer.queueConnector.BroadcastToHeimdall(msg)
		}
	}
}

func (syncer *Syncer) processReStakedEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {
	event := new(stakemanager.StakemanagerStaked)
	if err := helper.UnpackLog(abiObject, event, eventName, vLog); err != nil {
		logEventParseError(syncer.Logger, eventName, err)
	} else {
		syncer.Logger.Debug(
			"New event found",
			"event", eventName,
			"user", event.User.Hex(),
			"validatorId", event.ValidatorId,
			"activationEpoch", event.ActivationEpoch,
			"amount", event.Amount,
		)

		// // msg validator exit
		// msg := staking.NewMsgValidatorExit(
		// 	hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
		// 	event.ValidatorId.Uint64(),
		// 	vLog.TxHash,
		// )

		// // broadcast heimdall
		// syncer.queueConnector.BroadcastToHeimdall(msg)
	}
}

func (syncer *Syncer) processJailedEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {
	event := new(stakemanager.StakemanagerJailed)
	if err := helper.UnpackLog(abiObject, event, eventName, vLog); err != nil {
		logEventParseError(syncer.Logger, eventName, err)
	} else {
		syncer.Logger.Debug(
			"New event found",
			"event", eventName,
			"validatorID", event.ValidatorId,
			"exitEpoch", event.ExitEpoch,
		)

		// // msg validator exit
		// msg := staking.NewMsgValidatorExit(
		// 	hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
		// 	event.ValidatorId.Uint64(),
		// 	vLog.TxHash,
		// )

		// // broadcast heimdall
		// syncer.queueConnector.BroadcastToHeimdall(msg)
	}
}

//
// Process withdraw event
//

// func (syncer *Syncer) processWithdrawEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {
// 	event := new(depositmanager.DepositmanagerDeposit)
// 	if err := helper.UnpackLog(abiObject, event, eventName, vLog); err != nil {
// 		logEventParseError(syncer.Logger, eventName, err)
// 	} else {
// 		syncer.Logger.Debug(
// 			"New event found",
// 			"event", eventName,
// 			"user", event.User,
// 			"depositCount", event.DepositCount,
// 			"token", event.Token.String(),
// 		)

// 		// TODO dispatch to heimdall
// 	}
// }

//
// Process state synced event
//

func (syncer *Syncer) processStateSyncedEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {
	event := new(statesender.StatesenderStateSynced)
	if err := helper.UnpackLog(abiObject, event, eventName, vLog); err != nil {
		logEventParseError(syncer.Logger, eventName, err)
	} else {
		syncer.Logger.Debug(
			"New event found",
			"event", eventName,
			"id", event.Id,
			"contract", event.ContractAddress,
			"data", hex.EncodeToString(event.Data),
		)

		msg := clerkTypes.NewMsgEventRecord(
			hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
			hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			uint64(vLog.Index),
			event.Id.Uint64(),
		)

		// broadcast to heimdall
		syncer.queueConnector.BroadcastToHeimdall(msg)
	}
}

// processDelegatorBondEvent
func (syncer *Syncer) processDelegatorBondEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {
	event := new(delegationmanager.DelegationmanagerBonding)
	if err := helper.UnpackLog(abiObject, event, eventName, vLog); err != nil {
		logEventParseError(syncer.Logger, eventName, err)
	} else {
		syncer.Logger.Debug(
			"New delegator bond event found",
			"event", eventName,
			"DelegatorId", event.DelegatorId,
			"ValidatorId", event.ValidatorId,
			"Amount", event.Amount,
		)

		msg := stakingTypes.NewMsgDelegatorBond(
			hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
			hmTypes.DelegatorID(event.DelegatorId.Uint64()),
			hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			uint64(vLog.Index),
		)

		// broadcast to heimdall
		syncer.queueConnector.BroadcastToHeimdall(msg)

	}
}

// processDelegatorUnBondEvent
func (syncer *Syncer) processDelegatorUnBondEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {

	event := new(delegationmanager.DelegationmanagerUnBonding)

	if err := helper.UnpackLog(abiObject, event, eventName, vLog); err != nil {
		logEventParseError(syncer.Logger, eventName, err)
	} else {
		syncer.Logger.Debug(
			"New event found",
			"event", eventName,
			"DelegatorId", event.DelegatorId,
			"ValidatorId", event.ValidatorId,
		)
		msg := stakingTypes.NewMsgDelegatorUnBond(
			hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
			hmTypes.DelegatorID(event.DelegatorId.Uint64()),
			hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			uint64(vLog.Index),
		)

		// broadcast to heimdall
		syncer.queueConnector.BroadcastToHeimdall(msg)
	}
}

//
// Utils
//

// EventByID looks up a event by the topic id
func EventByID(abiObject *abi.ABI, sigdata []byte) *abi.Event {
	for _, event := range abiObject.Events {
		if bytes.Equal(event.Id().Bytes(), sigdata) {
			return &event
		}
	}
	return nil
}
