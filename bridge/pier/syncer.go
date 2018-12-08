package pier

import (
	"bytes"
	"context"
	"encoding/hex"
	"math/big"
	"os"
	"strconv"
	"time"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/contracts/stakemanager"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/staking"
	"math/rand"
)

// EventByID looks up a event by the topic id
func EventByID(abiObject *abi.ABI, sigdata []byte) *abi.Event {
	for _, event := range abiObject.Events {
		if bytes.Equal(event.Id().Bytes(), sigdata) {
			return &event
		}
	}
	return nil
}

func logEventParseError(logger log.Logger, name string, err error) {
	logger.Error("Error while parsing event", "name", name, "error", err)
}

func logEventTxBytesError(logger log.Logger, name string, err error) {
	logger.Error("Error while createing tx bytes", "name", name, "error", err)
}

func logEventBroadcastTxError(logger log.Logger, name string, err error) {
	logger.Error("Error while broadcasting tx bytes", "name", name, "error", err)
}

// ChainSyncer syncs validators and checkpoints
type ChainSyncer struct {
	// Base service
	common.BaseService

	// storage client
	storageClient *leveldb.DB

	// cli context for tendermint request
	cliContext cliContext.CLIContext

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
}

// NewChainSyncer returns new service object
func NewChainSyncer() *ChainSyncer {
	// create logger
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", chainSyncer)

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
	abis := []*abi.ABI{
		&rootchainABI,
		&stakemanagerABI,
	}

	// creating syncer object
	syncer := &ChainSyncer{
		storageClient: getBridgeDBInstance(viper.GetString(bridgeDBFlag)),
		cliContext:    cliContext.NewCLIContext(),
		abis:          abis,

		MainClient:           helper.GetMainClient(),
		RootChainInstance:    rootchainInstance,
		StakeManagerInstance: stakeManagerInstance,
		HeaderChannel:        make(chan *types.Header),
	}

	syncer.BaseService = *common.NewBaseService(logger, chainSyncer, syncer)
	return syncer
}

// StartHeaderProcess starts header process when they get new header
func (syncer *ChainSyncer) StartHeaderProcess(ctx context.Context) {
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
func (syncer *ChainSyncer) OnStart() error {
	syncer.BaseService.OnStart() // Always call the overridden method.

	// create cancellable context
	ctx, cancelSubscription := context.WithCancel(context.Background())
	syncer.cancelSubscription = cancelSubscription

	// create cancellable context
	headerCtx, cancelHeaderProcess := context.WithCancel(context.Background())
	syncer.cancelHeaderProcess = cancelHeaderProcess

	// start header process
	go syncer.StartHeaderProcess(headerCtx)

	// subscribe to new head
	subscription, err := syncer.MainClient.SubscribeNewHead(ctx, syncer.HeaderChannel)
	if err != nil {
		// start go routine to poll for new header using client object
		go syncer.StartPolling(ctx, defaultMainPollInterval)
	} else {
		// start go routine to listen new header using subscription
		go syncer.StartSubscription(ctx, subscription)
	}

	// subscribed to new head
	syncer.Logger.Debug("Subscribed to new head")

	return nil
}

// OnStop stops all necessary go routines
func (syncer *ChainSyncer) OnStop() {
	syncer.BaseService.OnStop() // Always call the overridden method.

	// close db
	closeBridgeDBInstance()

	// cancel subscription if any
	syncer.cancelSubscription()

	// cancel header process
	syncer.cancelHeaderProcess()
}

// StartPolling starts polling
func (syncer *ChainSyncer) StartPolling(ctx context.Context, pollInterval int) {
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

func (syncer *ChainSyncer) StartSubscription(ctx context.Context, subscription ethereum.Subscription) {
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

func (syncer *ChainSyncer) processHeader(newHeader *types.Header) {
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
		if result, ierr := strconv.Atoi(string(lastBlockBytes)); ierr == nil {
			if int64(result) >= newHeader.Number.Int64() {
				return
			}

			fromBlock = int64(result + 1)
		}
	}

	// set last block to storage
	syncer.storageClient.Put([]byte(lastBlockKey), []byte(newHeader.Number.String()), nil)

	// debug log
	syncer.Logger.Debug("Processing header", "fromBlock", fromBlock, "toBlock", newHeader.Number)

	if newHeader.Number.Int64()-fromBlock > 250 {
		// return if diff > 250
		return
	}

	// log
	syncer.Logger.Info("Querying event logs", "fromBlock", fromBlock, "toBlock", newHeader.Number)

	// draft a query
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(fromBlock),
		ToBlock:   newHeader.Number,
		Addresses: []ethCommon.Address{
			helper.GetRootChainAddress(),
			helper.GetStakeManagerAddress(),
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
				// case "Staked":
				// 	syncer.processStakedEvent(selectedEvent.Name, abiObject, &vLog)
				case "UnstakeInit":
					syncer.processUnstakeInitEvent(selectedEvent.Name, abiObject, &vLog)
					// case "SignerChange":
					// 	syncer.processSignerChangeEvent(selectedEvent.Name, abiObject, &vLog)
				}
			}
		}
	}
}

func (syncer *ChainSyncer) sendTx(eventName string, msg sdk.Msg) {
	txBytes, err := helper.CreateTxBytes(msg)
	if err != nil {
		logEventTxBytesError(syncer.Logger, eventName, err)
		return
	}

	// send tendermint request
	_, err = helper.SendTendermintRequest(syncer.cliContext, txBytes)
	if err != nil {
		logEventBroadcastTxError(syncer.Logger, eventName, err)
		return
	}
}

func (syncer *ChainSyncer) processCheckpointEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {
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
			"headerNumber", event.Number,
		)

		// create msg checkpoint ack message
		msg := checkpoint.NewMsgCheckpointAck(event.Number.Uint64(), uint64(rand.Int()))
		syncer.sendTx(eventName, msg)
	}
}

func (syncer *ChainSyncer) processStakedEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {
	event := new(stakemanager.StakemanagerStaked)
	if err := UnpackLog(abiObject, event, eventName, vLog); err != nil {
		logEventParseError(syncer.Logger, eventName, err)
	} else {
		// TOOD validator staked
		syncer.Logger.Debug(
			"New event found",
			"event", eventName,
			"validator", event.User.Hex(),
			"signer", event.Signer.Hex(),
			"activatonEpoch", event.ActivatonEpoch,
			"amount", event.Amount,
		)
	}
}

func (syncer *ChainSyncer) processUnstakeInitEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {
	event := new(stakemanager.StakemanagerUnstakeInit)
	if err := UnpackLog(abiObject, event, eventName, vLog); err != nil {
		logEventParseError(syncer.Logger, eventName, err)
	} else {
		// TOOD validator staked
		syncer.Logger.Debug(
			"New event found",
			"event", eventName,
			"validator", event.User.Hex(),
			"deactivatonEpoch", event.DeactivationEpoch,
			"amount", event.Amount,
		)

		// send validator exit message
		msg := staking.NewMsgValidatorExit(event.User, uint64(rand.Int()))
		syncer.sendTx(eventName, msg)
	}
}

func (syncer *ChainSyncer) processSignerChangeEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {
	event := new(stakemanager.StakemanagerSignerChange)
	if err := UnpackLog(abiObject, event, eventName, vLog); err != nil {
		logEventParseError(syncer.Logger, eventName, err)
	} else {
		// TOOD validator signer changed
		syncer.Logger.Debug(
			"New event found",
			"event", eventName,
			"validator", event.Validator.Hex(),
			"newSigner", event.NewSigner.Hex(),
			"oldSigner", event.OldSigner.Hex(),
		)
	}
}
