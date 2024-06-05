package processor

import (
	"encoding/base64"
	"io"
	"net/http"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tendermint/libs/log"
	httpClient "github.com/tendermint/tendermint/rpc/client"

	"github.com/maticnetwork/heimdall/bridge/setu/broadcaster"
	"github.com/maticnetwork/heimdall/bridge/setu/queue"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	clerkTypes "github.com/maticnetwork/heimdall/clerk/types"
	"github.com/maticnetwork/heimdall/helper"
)

// Processor defines a block header listener for Rootchain, Maticchain, Heimdall
type Processor interface {
	Start() error

	RegisterTasks()

	String() string

	Stop()
}

type BaseProcessor struct {
	Logger log.Logger
	name   string
	quit   chan struct{}

	// queue connector
	queueConnector *queue.QueueConnector

	// tx broadcaster
	txBroadcaster *broadcaster.TxBroadcaster

	// The "subclass" of BaseProcessor
	impl Processor

	// cli context
	cliCtx cliContext.CLIContext

	// contract caller
	contractConnector helper.ContractCaller

	// http client to subscribe to
	httpClient *httpClient.HTTP

	// storage client
	storageClient *leveldb.DB
}

// NewBaseProcessor creates a new BaseProcessor.
func NewBaseProcessor(cdc *codec.Codec, queueConnector *queue.QueueConnector, httpClient *httpClient.HTTP, txBroadcaster *broadcaster.TxBroadcaster, name string, impl Processor) *BaseProcessor {
	logger := util.Logger().With("service", "processor", "module", name)

	cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	cliCtx.BroadcastMode = client.BroadcastAsync
	cliCtx.TrustNode = true

	contractCaller, err := helper.NewContractCaller()
	if err != nil {
		logger.Error("Error while getting root chain instance", "error", err)
		panic(err)
	}

	if logger == nil {
		logger = log.NewNopLogger()
	}

	// creating syncer object
	return &BaseProcessor{
		Logger: logger,
		name:   name,
		quit:   make(chan struct{}),
		impl:   impl,

		cliCtx:            cliCtx,
		queueConnector:    queueConnector,
		contractConnector: contractCaller,
		txBroadcaster:     txBroadcaster,
		httpClient:        httpClient,
		storageClient:     util.GetBridgeDBInstance(viper.GetString(util.BridgeDBFlag)),
	}
}

// String implements Service by returning a string representation of the service.
func (bp *BaseProcessor) String() string {
	return bp.name
}

// OnStop stops all necessary go routines
func (bp *BaseProcessor) Stop() {
	// override to stop any go-routines in individual processors
}

// isOldTx checks if the transaction already exists in the chain or not
// It is a generic function, which is consumed in all processors
func (bp *BaseProcessor) isOldTx(_ cliContext.CLIContext, txHash string, logIndex uint64, eventType util.BridgeEvent, event interface{}) (bool, error) {
	defer util.LogElapsedTimeForStateSyncedEvent(event, "isOldTx", time.Now())

	queryParam := map[string]interface{}{
		"txhash":   txHash,
		"logindex": logIndex,
	}

	// define the endpoint based on the type of event
	var endpoint string

	switch eventType {
	case util.StakingEvent:
		endpoint = helper.GetHeimdallServerEndpoint(util.StakingTxStatusURL)
	case util.TopupEvent:
		endpoint = helper.GetHeimdallServerEndpoint(util.TopupTxStatusURL)
	case util.ClerkEvent:
		endpoint = helper.GetHeimdallServerEndpoint(util.ClerkTxStatusURL)
	case util.SlashingEvent:
		endpoint = helper.GetHeimdallServerEndpoint(util.SlashingTxStatusURL)
	}

	url, err := util.CreateURLWithQuery(endpoint, queryParam)
	if err != nil {
		bp.Logger.Error("Error in creating url", "endpoint", endpoint, "error", err)
		return false, err
	}

	res, err := helper.FetchFromAPI(bp.cliCtx, url)
	if err != nil {
		bp.Logger.Error("Error fetching tx status", "url", url, "error", err)
		return false, err
	}

	var status bool
	if err := jsoniter.ConfigFastest.Unmarshal(res.Result, &status); err != nil {
		bp.Logger.Error("Error unmarshalling tx status received from Heimdall Server", "error", err)
		return false, err
	}

	return status, nil
}

// checkTxAgainstMempool checks if the transaction is already in the mempool or not
// It is consumed only for `clerk` processor
func (bp *BaseProcessor) checkTxAgainstMempool(msg types.Msg, event interface{}) (bool, error) {
	defer util.LogElapsedTimeForStateSyncedEvent(event, "checkTxAgainstMempool", time.Now())

	endpoint := helper.GetConfig().TendermintRPCUrl + util.TendermintUnconfirmedTxsURL

	resp, err := helper.Client.Get(endpoint)
	if err != nil || resp.StatusCode != http.StatusOK {
		bp.Logger.Error("Error fetching mempool tx", "url", endpoint, "error", err)
		return false, err
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		bp.Logger.Error("Error fetching mempool tx", "error", err)
		return false, err
	}

	// a minimal response of the unconfirmed txs
	var response util.TendermintUnconfirmedTxs

	err = jsoniter.ConfigFastest.Unmarshal(body, &response)
	if err != nil {
		bp.Logger.Error("Error unmarshalling response received from Heimdall Server", "error", err)
		return false, err
	}

	// Iterate over txs present in the mempool
	// We can verify if the message we're about to send is present by
	// checking the type of transaction, the transaction hash and log index
	// present in the data of transaction

	status := false
Loop:
	for _, txn := range response.Result.Txs {
		// Tendermint encodes the transactions with base64 encoding. Decode it first.
		txBytes, err := base64.StdEncoding.DecodeString(txn)
		if err != nil {
			bp.Logger.Error("Error decoding tx (base64 decoder) while checking against mempool", "error", err)
			continue
		}

		// Unmarshal the transaction from bytes
		decodedTx, err := helper.GetTxDecoder(bp.cliCtx.Codec)(txBytes)
		if err != nil {
			bp.Logger.Error("Error decoding tx (tx decoder) while checking against mempool", "error", err)
			continue
		}
		txMsg := decodedTx.GetMsgs()[0]

		// We only need to check for `event-record` type transactions.
		// If required, add case for others here.
		switch txMsg.Type() {
		case "event-record":

			// typecast the txs for clerk type message
			mempoolTxMsg, ok := txMsg.(clerkTypes.MsgEventRecord)
			if !ok {
				bp.Logger.Error("Unable to typecast message to clerk event record while checking against mempool")
				continue Loop
			}

			// typecast the msg for clerk type message
			clerkMsg, ok := msg.(clerkTypes.MsgEventRecord)
			if !ok {
				bp.Logger.Error("Unable to typecast message to clerk event record while checking against mempool")
				continue Loop
			}

			// check the transaction hash in message
			if clerkMsg.GetTxHash() != mempoolTxMsg.GetTxHash() {
				continue Loop
			}

			// check the log index in the message
			if clerkMsg.GetLogIndex() != mempoolTxMsg.GetLogIndex() {
				continue Loop
			}

			// If we reach here, there's already a same transaction in the mempool
			status = true
			break Loop
		default:
			// ignore
		}
	}

	return status, nil
}
