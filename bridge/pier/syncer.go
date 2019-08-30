package pier

import (
	"bytes"
	"context"
	"encoding/hex"
	"math/big"
	"strconv"
	"time"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/cosmos/cosmos-sdk/client"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/contracts/depositmanager"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/staking"

	"github.com/maticnetwork/heimdall/contracts/stakemanager"
	"github.com/maticnetwork/heimdall/helper"
	hmtypes "github.com/maticnetwork/heimdall/types"
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

	// connector to queues
	qConnector QueueConnector

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

	// tx encoder
	txEncoder authTypes.TxBuilder

	// header listener subscription
	cancelHeaderProcess context.CancelFunc
}

// NewSyncer returns new service object for syncing events
func NewSyncer(connector QueueConnector) *Syncer {
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

	cliCtx := cliContext.NewCLIContext()
	cliCtx.BroadcastMode = client.BroadcastAsync

	// creating syncer object
	syncer := &Syncer{
		storageClient:        getBridgeDBInstance(viper.GetString(BridgeDBFlag)),
		cliContext:           cliCtx,
		abis:                 abis,
		qConnector:           connector,
		MainClient:           helper.GetMainClient(),
		RootChainInstance:    rootchainInstance,
		StakeManagerInstance: stakeManagerInstance,
		HeaderChannel:        make(chan *types.Header),
		txEncoder:            authTypes.NewTxBuilderFromCLI().WithTxEncoder(helper.GetTxEncoder()),
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
			helper.GetDepositManagerAddress(),
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
				case "Deposit":
					syncer.processDepositEvent(selectedEvent.Name, abiObject, &vLog)
				case "Withdraw":
					// syncer.processWithdrawEvent(selectedEvent.Name, abiObject, &vLog)
				}
			}
		}
	}
}

func (syncer *Syncer) sendTx(eventName string, msg sdk.Msg) {
	_, err := helper.BroadcastMsgs(syncer.cliContext, []sdk.Msg{msg})
	if err != nil {
		logEventBroadcastTxError(syncer.Logger, eventName, err)
		return
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
			"headerNumber", event.Number,
		)

		// create msg checkpoint ack message
		msg := checkpoint.NewMsgCheckpointAck(helper.GetFromAddress(syncer.cliContext), event.Number.Uint64())
		syncer.qConnector.DispatchToHeimdall(msg)
	}
}

func (syncer *Syncer) processStakedEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {
	event := new(stakemanager.StakemanagerStaked)
	if err := UnpackLog(abiObject, event, eventName, vLog); err != nil {
		logEventParseError(syncer.Logger, eventName, err)
	} else {
		// TOOD validator staked
		syncer.Logger.Debug(
			"New event found",
			"event", eventName,
			"validator", event.User.Hex(),
			"ID", event.ValidatorId,
			"activatonEpoch", event.ActivatonEpoch,
			"amount", event.Amount,
		)
		if bytes.Compare(event.User.Bytes(), helper.GetPubKey().Address().Bytes()) == 0 {
			msg := staking.NewMsgValidatorJoin(helper.GetFromAddress(syncer.cliContext), event.ValidatorId.Uint64(), hmtypes.NewPubKey(helper.GetPubKey().Bytes()), vLog.TxHash)
			syncer.qConnector.DispatchToHeimdall(msg)
		}
	}
}

func (syncer *Syncer) processUnstakeInitEvent(eventName string, abiObject *abi.ABI, vLog *types.Log) {
	event := new(stakemanager.StakemanagerUnstakeInit)
	if err := UnpackLog(abiObject, event, eventName, vLog); err != nil {
		logEventParseError(syncer.Logger, eventName, err)
	} else {
		// TOOD validator staked
		syncer.Logger.Debug(
			"New event found",
			"event", eventName,
			"validator", event.User.Hex(),
			"validatorID", event.ValidatorId,
			"deactivatonEpoch", event.DeactivationEpoch,
			"amount", event.Amount,
		)
		// send validator exit message
		msg := staking.NewMsgValidatorExit(helper.GetFromAddress(syncer.cliContext), event.ValidatorId.Uint64(), vLog.TxHash)
		syncer.qConnector.DispatchToHeimdall(msg)
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
		// if bytes.Compare(event.NewSigner.Bytes(), helper.GetPubKey().Address().Bytes()) == 0 {
		// 	msg := staking.NewMsgValidatorUpdate(event.ValidatorId.Uint64(), hmtypes.NewPubKey(helper.GetPubKey().Bytes()), vLog.TxHash)
		// 	syncer.qConnector.DispatchToHeimdall(msg)
		// }
	}
}

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

// EventByID looks up a event by the topic id
func EventByID(abiObject *abi.ABI, sigdata []byte) *abi.Event {
	for _, event := range abiObject.Events {
		if bytes.Equal(event.Id().Bytes(), sigdata) {
			return &event
		}
	}
	return nil
}
