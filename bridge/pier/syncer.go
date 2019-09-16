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
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tendermint/libs/common"
	httpClient "github.com/tendermint/tendermint/rpc/client"

	"github.com/maticnetwork/heimdall/checkpoint"
	clerkTypes "github.com/maticnetwork/heimdall/clerk/types"
	"github.com/maticnetwork/heimdall/contracts/depositmanager"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/contracts/stakemanager"
	"github.com/maticnetwork/heimdall/contracts/statesyncer"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

const (
	headerEvent      = "NewHeaderBlock"
	stakeInitEvent   = "Staked"
	unstakeInitEvent = "UnstakeInit"
	signerChange     = "SignerChange"
	depositEvent     = "Deposit"

	lastBlockKey = "last-block" // storage key
)

// Syncer syncs validators and checkpoints
type Syncer struct {
	// Base service
	common.BaseService

	// storage client
	storageClient *leveldb.DB

	// ABIs
	abis []*abi.ABI

	// Mainchain client
	MainClient *ethclient.Client
	// Rootchain instance
	RootChainInstance *rootchain.Rootchain
	// Stake manager instance
	StakeManagerInstance *stakemanager.Stakemanager
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
	// root chain instance
	rootchainInstance, err := helper.GetRootChainInstance()
	if err != nil {
		logger.Error("Error while getting root chain instance", "error", err)
		panic(err)
	}

	// stake manager instance
	stakeManagerInstance, err := helper.GetStakeManagerInstance()
	if err != nil {
		logger.Error("Error while getting stake manager instance", "error", err)
		panic(err)
	}

	rootchainABI, _ := helper.GetRootChainABI()
	stakemanagerABI, _ := helper.GetStakeManagerABI()
	depositManagerABI, _ := helper.GetDepositManagerABI()
	abis := []*abi.ABI{
		&rootchainABI,
		&stakemanagerABI,
		&depositManagerABI,
	}

	cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	cliCtx.BroadcastMode = client.BroadcastAsync
	cliCtx.TrustNode = true

	// creating syncer object
	syncer := &Syncer{
		storageClient: getBridgeDBInstance(viper.GetString(BridgeDBFlag)),

		cliCtx:         cliCtx,
		queueConnector: queueConnector,
		httpClient:     httpClient,

		abis:                 abis,
		MainClient:           helper.GetMainClient(),
		RootChainInstance:    rootchainInstance,
		StakeManagerInstance: stakeManagerInstance,
		HeaderChannel:        make(chan *types.Header),
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
	subscription, err := syncer.MainClient.SubscribeNewHead(ctx, syncer.HeaderChannel)
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
func (syncer *Syncer) startPolling(ctx context.Context, pollInterval int) {
	// How often to fire the passed in function in second
	interval := time.Duration(pollInterval) * time.Millisecond

	// Setup the ticket and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)

	// start listening
	for {
		select {
		case <-ticker.C:
			header, err := syncer.MainClient.HeaderByNumber(ctx, nil)
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
	// default fromBlock
	fromBlock := newHeader.Number.Int64()
	// get last block from storage
	hasLastBlock, _ := syncer.storageClient.Has([]byte(lastBlockKey), nil)
	if hasLastBlock {
		lastBlockBytes, err := syncer.storageClient.Get([]byte(lastBlockKey), nil)
		if err != nil {
			syncer.Logger.Info("Error while fetching last block bytes from storage", "error", err)
			return
		}

		syncer.Logger.Debug("Got last block from bridge storage", "lastBlock", string(lastBlockBytes))
		if result, err := strconv.Atoi(string(lastBlockBytes)); err == nil {
			if int64(result) >= newHeader.Number.Int64() {
				return
			}

			fromBlock = int64(result + 1)
		}
	}

	// confirmation
	toBlock := newHeader.Number
	confirmationBlocks := big.NewInt(0).SetUint64(helper.GetConfig().ConfirmationBlocks)
	confirmationBlocks = confirmationBlocks.Add(confirmationBlocks, big.NewInt(1))
	if toBlock.Uint64() > confirmationBlocks.Uint64() {
		toBlock = toBlock.Sub(toBlock, confirmationBlocks)
	}

	// set last block to storage
	syncer.storageClient.Put([]byte(lastBlockKey), []byte(toBlock.String()), nil)

	// debug log
	syncer.Logger.Debug("Processing header", "fromBlock", fromBlock, "toBlock", toBlock)

	if newHeader.Number.Int64()-fromBlock > 250 {
		// return if diff > 250
		return
	}

	// log
	syncer.Logger.Info("Querying event logs", "fromBlock", fromBlock, "toBlock", toBlock)

	// draft a query
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(fromBlock),
		ToBlock:   toBlock,
		Addresses: []ethCommon.Address{
			helper.GetRootChainAddress(),
			helper.GetStakeManagerAddress(),
			helper.GetDepositManagerAddress(),
			helper.GetStateSenderAddress(),
		},
	}

	// get all logs
	logs, err := syncer.MainClient.FilterLogs(context.Background(), query)
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
				switch selectedEvent.Name {
				case "NewHeaderBlock":
					syncer.processCheckpointEvent(selectedEvent.Name, abiObject, &vLog)
				case "Staked":
					syncer.processStakedEvent(selectedEvent.Name, abiObject, &vLog)
				case "UnstakeInit":
					syncer.processUnstakeInitEvent(selectedEvent.Name, abiObject, &vLog)
				case "SignerChange":
					syncer.processSignerChangeEvent(selectedEvent.Name, abiObject, &vLog)
				case "ReStaked":
					syncer.processReStakedEvent(selectedEvent.Name, abiObject, &vLog)
				case "Jailed":
					syncer.processJailedEvent(selectedEvent.Name, abiObject, &vLog)
				case "Deposit":
					syncer.processDepositEvent(selectedEvent.Name, abiObject, &vLog)
				case "StateSynced":
					syncer.processDepositEvent(selectedEvent.Name, abiObject, &vLog)
				case "Withdraw":
					syncer.processWithdrawEvent(selectedEvent.Name, abiObject, &vLog)
				}
			}
		}
	}
}

func (syncer *Syncer) processCheckpointEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {
	event := new(rootchain.RootchainNewHeaderBlock)
	if err := UnpackLog(abiObject, event, eventName, vLog); err != nil {
		logEventParseError(syncer.Logger, eventName, err)
	} else {
		syncer.Logger.Info(
			"New event found",
			"event", eventName,
			"start", event.Start,
			"end", event.End,
			"root", "0x"+hex.EncodeToString(event.Root[:]),
			"proposer", event.Proposer.Hex(),
			"headerNumber", event.HeaderBlockId,
		)

		// create msg checkpoint ack message
		msg := checkpoint.NewMsgCheckpointAck(helper.GetFromAddress(syncer.cliCtx), event.HeaderBlockId.Uint64())
		syncer.queueConnector.BroadcastToHeimdall(msg)
	}
}

func (syncer *Syncer) processStakedEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {
	event := new(stakemanager.StakemanagerStaked)
	if err := UnpackLog(abiObject, event, eventName, vLog); err != nil {
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
		if bytes.Compare(event.User.Bytes(), helper.GetAddress()) == 0 {
			pubkey := helper.GetPubKey()
			msg := staking.NewMsgValidatorJoin(
				hmTypes.BytesToHeimdallAddress(event.User.Bytes()),
				event.ValidatorId.Uint64(),
				hmTypes.NewPubKey(pubkey[:]),
				hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			)

			// process staked
			syncer.queueConnector.BroadcastToHeimdall(msg)
		}
	}
}

func (syncer *Syncer) processUnstakeInitEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {
	event := new(stakemanager.StakemanagerUnstakeInit)
	if err := UnpackLog(abiObject, event, eventName, vLog); err != nil {
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
		msg := staking.NewMsgValidatorExit(
			hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
			event.ValidatorId.Uint64(),
			hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
		)

		// broadcast heimdall
		syncer.queueConnector.BroadcastToHeimdall(msg)
	}
}

func (syncer *Syncer) processSignerChangeEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {
	event := new(stakemanager.StakemanagerSignerChange)
	if err := UnpackLog(abiObject, event, eventName, vLog); err != nil {
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
			msg := staking.NewMsgValidatorUpdate(
				hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
				event.ValidatorId.Uint64(),
				hmTypes.NewPubKey(pubkey[:]),
				hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			)

			// process signer update
			syncer.queueConnector.BroadcastToHeimdall(msg)
		}
	}
}

func (syncer *Syncer) processReStakedEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {
	event := new(stakemanager.StakemanagerStaked)
	if err := UnpackLog(abiObject, event, eventName, vLog); err != nil {
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
	if err := UnpackLog(abiObject, event, eventName, vLog); err != nil {
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
// Process deposit event
//

func (syncer *Syncer) processDepositEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {
	event := new(depositmanager.DepositmanagerDeposit)
	if err := UnpackLog(abiObject, event, eventName, vLog); err != nil {
		logEventParseError(syncer.Logger, eventName, err)
	} else {
		syncer.Logger.Debug(
			"New event found",
			"event", eventName,
			"user", event.User,
			"depositCount", event.DepositCount,
			"token", event.Token.String(),
		)

		// TODO dispatch to heimdall
	}
}

//
// Process withdraw event
//

func (syncer *Syncer) processWithdrawEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {
	event := new(depositmanager.DepositmanagerDeposit)
	if err := UnpackLog(abiObject, event, eventName, vLog); err != nil {
		logEventParseError(syncer.Logger, eventName, err)
	} else {
		syncer.Logger.Debug(
			"New event found",
			"event", eventName,
			"user", event.User,
			"depositCount", event.DepositCount,
			"token", event.Token.String(),
		)

		// TODO dispatch to heimdall
	}
}

//
// Process state synced event
//

func (syncer *Syncer) processStateSyncedEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {
	event := new(statesyncer.StatesyncerStateSynced)
	if err := UnpackLog(abiObject, event, eventName, vLog); err != nil {
		logEventParseError(syncer.Logger, eventName, err)
	} else {
		syncer.Logger.Debug(
			"New event found",
			"event", eventName,
			"id", event.Id,
			"contract", event.ContractAddress,
			"data", hex.EncodeToString(event.Data),
		)

		// TODO dispatch to heimdall
		msg := clerkTypes.NewMsgStateRecord(
			hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
			hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			event.Id.Uint64(),
			hmTypes.BytesToHeimdallAddress(event.ContractAddress.Bytes()),
			event.Data,
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
