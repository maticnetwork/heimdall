package pier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/pkg/errors"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tendermint/libs/common"

	"github.com/cosmos/cosmos-sdk/client"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/helper"
	hmtypes "github.com/maticnetwork/heimdall/types"
	httpClient "github.com/tendermint/tendermint/rpc/client"
	tmTypes "github.com/tendermint/tendermint/types"
)

// Checkpointer to propose
type Checkpointer struct {
	// Base service
	common.BaseService
	// storage client
	storageClient *leveldb.DB
	// header channel
	HeaderChannel chan *types.Header
	// cancel function for poll/subscription
	cancelSubscription context.CancelFunc
	// header listener subscription
	cancelHeaderProcess context.CancelFunc
	// queue connnector
	qConnector QueueConnector
	// contract caller
	contractConnector helper.ContractCaller
	// http client to subscribe to
	httpClient *httpClient.HTTP
	// tx encoder
	txEncoder authTypes.TxBuilder
	// context for sending transactions to heimdall
	cliCtx cliContext.CLIContext
}

// NewCheckpointer returns new service object
func NewCheckpointer(connector QueueConnector, cdc *codec.Codec) *Checkpointer {
	// create logger
	logger := Logger.With("module", HeimdallCheckpointer)
	cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	cliCtx.BroadcastMode = client.BroadcastAsync

	contractCaller, err := helper.NewContractCaller()
	if err != nil {
		logger.Error("Error while getting root chain instance", "error", err)
		panic(err)
	}

	// creating checkpointer object
	checkpointer := &Checkpointer{
		storageClient:     getBridgeDBInstance(viper.GetString(BridgeDBFlag)),
		HeaderChannel:     make(chan *types.Header),
		qConnector:        connector,
		contractConnector: contractCaller,
		cliCtx:            cliCtx,
		txEncoder:         authTypes.NewTxBuilderFromCLI().WithTxEncoder(helper.GetTxEncoder()).WithChainID(helper.GetGenesisDoc().ChainID),
		httpClient:        httpClient.NewHTTP("tcp://0.0.0.0:26657", "/websocket"),
	}

	checkpointer.BaseService = *common.NewBaseService(logger, HeimdallCheckpointer, checkpointer)
	return checkpointer
}

// startHeaderProcess starts header process when they get new header
func (c *Checkpointer) startHeaderProcess(ctx context.Context) {
	for {
		select {
		case newHeader := <-c.HeaderChannel:
			c.sendRequest(newHeader)
		case <-ctx.Done():
			return
		}
	}
}

// OnStart starts new block subscription
func (c *Checkpointer) OnStart() error {
	c.BaseService.OnStart() // Always call the overridden method.

	// create cancellable context
	ctx, cancelSubscription := context.WithCancel(context.Background())
	c.cancelSubscription = cancelSubscription

	// create cancellable context
	headerCtx, cancelHeaderProcess := context.WithCancel(context.Background())
	c.cancelHeaderProcess = cancelHeaderProcess
	// start http client
	err := c.httpClient.Start()
	if err != nil {
		log.Fatalf("Error connecting to server %v", err)
	}

	// start header process
	go c.startHeaderProcess(headerCtx)

	// subscribe to new head
	subscription, err := c.contractConnector.MaticClient.SubscribeNewHead(ctx, c.HeaderChannel)
	if err != nil {
		// start go routine to poll for new header using client object
		go c.startPolling(ctx, helper.GetConfig().CheckpointerPollInterval)
	} else {
		// start go routine to listen new header using subscription
		go c.startSubscription(ctx, subscription)
	}

	// subscribed to new head
	c.Logger.Debug("Subscribed to new head")

	return nil
}

// OnStop stops all necessary go routines
func (c *Checkpointer) OnStop() {
	c.BaseService.OnStop() // Always call the overridden method.

	// close bridge db instance
	closeBridgeDBInstance()

	// cancel subscription if any
	c.cancelSubscription()

	// cancel header process
	c.cancelHeaderProcess()

	// stop http client
	c.httpClient.Stop()
}

func (c *Checkpointer) startPolling(ctx context.Context, pollInterval int) {
	// How often to fire the passed in function in second
	interval := time.Duration(pollInterval) * time.Millisecond

	// Setup the ticket and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)

	// start listening
	for {
		select {
		case <-ticker.C:
			fmt.Printf("Checking proposer status")
			if isProposer(c.cliCtx) {
				header, err := c.contractConnector.MaticClient.HeaderByNumber(ctx, nil)
				if err == nil && header != nil {
					// send data to channel
					c.HeaderChannel <- header
				} else if err != nil {
					c.Logger.Error("Unable to fetch header by number", "Error", err)
				}
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (c *Checkpointer) startSubscription(ctx context.Context, subscription ethereum.Subscription) {
	for {
		select {
		case err := <-subscription.Err():
			// stop service
			c.Logger.Error("Error while subscribing new blocks", "error", err)
			c.Stop()

			// cancel subscription
			c.cancelSubscription()
			return
		case <-ctx.Done():
			return
		}
	}
}

func (c *Checkpointer) sendRequest(newHeader *types.Header) {
	c.Logger.Debug("New block detected", "blockNumber", newHeader.Number)

	// fetch contract state
	contractState := make(chan ContractCheckpoint, 1)
	var lastContractCheckpoint ContractCheckpoint

	// fetch heimdall state
	bufferChan := make(chan HeimdallCheckpoint, 1)
	var bufferedCheckpoint HeimdallCheckpoint

	commitedChan := make(chan HeimdallCheckpoint, 1)
	var commitedCheckpoint HeimdallCheckpoint

	var wg sync.WaitGroup
	wg.Add(3)

	c.Logger.Debug("Collecting contract and heimdall checkpoint state")
	go c.fetchContractCheckpointState(newHeader.Number.Uint64(), &wg, contractState)
	go c.fetchBufferedCheckpoint(&wg, bufferChan)
	go c.fetchCommittedCheckpoint(&wg, commitedChan)

	// wait for state collection
	wg.Wait()

	c.Logger.Info("Fetched contract and heimdall states")

	lastContractCheckpoint = <-contractState
	if lastContractCheckpoint.err != nil {
		c.Logger.Error("Error fetching details from contract ", "Error", lastContractCheckpoint.err)
		return
	}

	c.Logger.Debug("Contract Details fetched", "Start", lastContractCheckpoint.start, "End", lastContractCheckpoint.end, "Currentheader", lastContractCheckpoint.currentHeaderBlock)

	bufferedCheckpoint = <-bufferChan
	if !bufferedCheckpoint.found {
		c.Logger.Info("Buffer not found, sending new checkpoint", "Found", bufferedCheckpoint.found)
	} else if bufferedCheckpoint.start != 0 {
		c.Logger.Debug("Checkpoint found in buffer", "Start", bufferedCheckpoint.start, "End", bufferedCheckpoint.end)
	} else {
		c.Logger.Debug("Sending first checkpoint")
	}

	commitedCheckpoint = <-commitedChan
	if !commitedCheckpoint.found {
	}

	// ACK needs to be sent
	if bufferedCheckpoint.end+1 == bufferedCheckpoint.start {
		c.Logger.Debug("Sending ACK", "BufferedCheckpointEnd", bufferedCheckpoint.end, "ContractStart", bufferedCheckpoint.start)

		// channel for collecting tx hash
		txHash := make(chan string)
		ctx, cancel := context.WithTimeout(context.Background(), TransactionTimeout)
		defer cancel()

		// calculate header index
		headerNumber := lastContractCheckpoint.currentHeaderBlock.Sub(lastContractCheckpoint.currentHeaderBlock, big.NewInt(int64(helper.GetConfig().ChildBlockInterval)))

		go c.broadcastACK(txHash, headerNumber.Uint64())
		select {
		case <-ctx.Done():
			c.Logger.Error("Error sending ACK", "Error", ctx.Err())
		case <-txHash:
			c.Logger.Info("Sent transaction")
		}
		return
	}

	start := lastContractCheckpoint.start
	end := lastContractCheckpoint.end

	// Get root hash
	root, err := checkpoint.GetHeaders(start, end)
	if err != nil {
		return
	}
	c.Logger.Info("Creating new checkpoint", "start", start, "end", end, "root", ethCommon.BytesToHash(root))

	// channel for collecting tx hash
	txHash := make(chan string)
	checkpointCtx, cancel := context.WithTimeout(context.Background(), TransactionTimeout)
	defer cancel()

	go c.broadcastCheckpoint(txHash, start, end, ethCommon.BytesToAddress(helper.GetAddress()), ethCommon.BytesToHash(root))
	select {
	case <-checkpointCtx.Done():
		c.Logger.Error("Error sending ACK", "Error", checkpointCtx.Err())
	case <-txHash:
		c.Logger.Info("Sent transaction")
	}
}

// fetched contract checkpoint state and returns the next probable checkpoint that needs to be sent
func (c *Checkpointer) fetchContractCheckpointState(lastHeader uint64, wg *sync.WaitGroup, contractState chan<- ContractCheckpoint) {
	defer wg.Done()
	lastCheckpointEnd, err := c.contractConnector.CurrentChildBlock()
	if err != nil {
		c.Logger.Error("Error while fetching current child block from rootchain", "error", err)
		return
	}
	var start, end uint64
	start = lastCheckpointEnd

	// add 1 if start > 0
	if start > 0 {
		start = start + 1
	}

	// get diff
	diff := lastHeader - start + 1

	// process if diff > 0 (positive)
	if diff > 0 {
		expectedDiff := diff - diff%helper.GetConfig().AvgCheckpointLength
		if expectedDiff > 0 {
			expectedDiff = expectedDiff - 1
		}

		// cap with max checkpoint length
		if expectedDiff > helper.GetConfig().MaxCheckpointLength-1 {
			expectedDiff = helper.GetConfig().MaxCheckpointLength - 1
		}

		// get end result
		end = expectedDiff + start

		c.Logger.Debug("Calculating checkpoint eligibility", "latest", lastHeader, "start", start, "end", end)
	}
	currentHeaderBlock, err := c.contractConnector.CurrentHeaderBlock()
	currentHeaderBlockNumber := big.NewInt(int64(currentHeaderBlock))
	if err != nil {
		c.Logger.Error("Error while fetching current header block number from rootchain", "error", err)
		contractState <- NewContractCheckpoint(0, 0, currentHeaderBlockNumber, err)
		return
	}
	currentHeaderNum := currentHeaderBlockNumber

	// Handle when block producers go down
	if end == 0 || end == start || (0 < diff && diff < helper.GetConfig().AvgCheckpointLength) {
		c.Logger.Debug("Fetching last header block to calculate time")
		// fetch current header block
		_, _, _, createdAt, err := c.contractConnector.GetHeaderInfo(currentHeaderBlockNumber.Sub(currentHeaderBlockNumber, big.NewInt(int64(helper.GetConfig().ChildBlockInterval))).Uint64())
		if err != nil {
			c.Logger.Error("Error while fetching current header block object from rootchain", "error", err)
			contractState <- NewContractCheckpoint(0, 0, currentHeaderBlockNumber, err)
			return
		}
		lastCheckpointTime := createdAt
		currentTime := time.Now().Unix()
		if currentTime-int64(lastCheckpointTime) > int64(helper.GetConfig().MaxCheckpointLength*2) {
			c.Logger.Info("Force push checkpoint", "currentTime", currentTime, "lastCheckpointTime", lastCheckpointTime, "defaultForcePushInterval", defaultForcePushInterval, "end", lastHeader)
			end = lastHeader
		}
	}

	if end == 0 || start >= end {
		c.Logger.Info("Waiting for 256 blocks or invalid start end formation", "Start", start, "End", end)
		contractState <- NewContractCheckpoint(0, 0, currentHeaderBlockNumber, errors.New("Invalid start end formation"))
		return
	}
	contractCheckpointData := NewContractCheckpoint(start, end, currentHeaderNum, nil)
	contractState <- contractCheckpointData
	return
}

// fetch checkpoint present in buffer from heimdall
func (c *Checkpointer) fetchBufferedCheckpoint(wg *sync.WaitGroup, bufferedCheckpoint chan<- HeimdallCheckpoint) {
	defer wg.Done()
	c.Logger.Info("Fetching checkpoint in buffer")
	_checkpoint, err := c.fetchCheckpoint(LatestCheckpointURL)
	if err != nil {
		c.Logger.Error("Error while fetching data from server", "error", err)
		bufferedCheckpoint <- NewHeimdallCheckpoint(_checkpoint.StartBlock, _checkpoint.EndBlock, false)
		return
	}
	bufferedCheckpoint <- NewHeimdallCheckpoint(_checkpoint.StartBlock, _checkpoint.EndBlock, true)
	return
}

// fetches latest committed checkpoint from heimdall
func (c *Checkpointer) fetchCommittedCheckpoint(wg *sync.WaitGroup, lastCheckpoint chan<- HeimdallCheckpoint) {
	defer wg.Done()
	c.Logger.Info("Fetching last committed checkpoint")
	_checkpoint, err := c.fetchCheckpoint(LatestCheckpointURL)
	if err != nil {
		c.Logger.Error("Error while fetching data from server", "error", err)
		lastCheckpoint <- NewHeimdallCheckpoint(_checkpoint.StartBlock, _checkpoint.EndBlock, false)
		return
	}
	lastCheckpoint <- NewHeimdallCheckpoint(_checkpoint.StartBlock, _checkpoint.EndBlock, true)
	return
}

// fetches checkpoint from given URL
func (c *Checkpointer) fetchCheckpoint(url string) (checkpoint hmtypes.CheckpointBlockHeader, err error) {
	resp, err := http.Get(url)
	if err != nil {
		c.Logger.Error("Unable to send request to get proposer", "Error", err)
		return checkpoint, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.Logger.Error("Unable to read data from response", "Error", err)
			return checkpoint, err
		}
		if err := json.Unmarshal(body, &checkpoint); err != nil {
			c.Logger.Error("Error unmarshalling checkpoint", "error", err)
			return checkpoint, err
		}
		return checkpoint, nil
	}
	return checkpoint, fmt.Errorf("Invalid response from rest server. Status: %v URL: %v", resp.Status, url)
}

// broadcastACK broadcasts ack for a checkpoint to heimdall
func (c *Checkpointer) broadcastACK(txhash chan string, headerIdx uint64) {
	// create and send checkpoint ACK message
	msg := checkpoint.NewMsgCheckpointAck(hmtypes.BytesToHeimdallAddress(helper.GetAddress()), headerIdx)

	resp, err := helper.BuildAndBroadcastMsgs(c.cliCtx, c.txEncoder, []sdk.Msg{msg})
	if err != nil {
		c.Logger.Error("Unable to send checkpoint ack to heimdall", "Error", err, "HeaderIndex", headerIdx)
	}

	c.Logger.Debug("Checkpoint ACK tx commited", "TxHash", resp.TxHash, "HeaderIndex", headerIdx)
	// send transaction hash to channel
	txhash <- resp.TxHash
}

// broadcastCheckpoint broadcasts checkpoint tx to heimdall
func (c *Checkpointer) broadcastCheckpoint(txhash chan string, start uint64, end uint64, proposer ethCommon.Address, root ethCommon.Hash) {
	msg := checkpoint.NewMsgCheckpointBlock(
		helper.GetFromAddress(c.cliCtx),
		start,
		end,
		root,
		uint64(time.Now().Unix()),
	)

	txBytes, err := helper.GetSignedTxBytes(c.cliCtx, c.txEncoder, []sdk.Msg{msg})
	if err != nil {
		c.Logger.Error("Error creating tx bytes", "error", err)
		return
	}

	response, err := helper.BroadcastTxBytes(c.cliCtx, txBytes, client.BroadcastSync)
	if err != nil {
		c.Logger.Error("Error sending checkpoint tx", "error", err, "start", start, "end", end)
		return
	}
	c.Logger.Info("Subscribing to checkpoint tx", "hash", response.TxHash, "start", start, "end", end, "root", root.String())

	go c.SubscribeToTx(txBytes, start, end)

	txhash <- response.TxHash
}

// SubscribeToTx subscribes to a broadcasted Tx and waits for its commitment to a block
func (c *Checkpointer) SubscribeToTx(tx tmTypes.Tx, start, end uint64) error {
	data, err := WaitForOneEvent(tx, c.httpClient)
	if err != nil {
		c.Logger.Error("Unable to wait for tx", "error", err)
		return err
	}
	fmt.Printf("data %v", data)

	switch t := data.(type) {
	case tmTypes.EventDataTx:
		c.DispatchCheckpoint(t.Height, tx, start, end)
	default:
		c.Logger.Info("No cases matched")
	}
	return nil
}

// fetchVotes fetches votes and extracts sigs from it
func (c *Checkpointer) fetchVotes(height int64) (votes []*tmTypes.CommitSig, sigs []byte, chainID string, err error) {
	c.Logger.Debug("Subscribing to new height", "height", height+1)

	var blockDetails *tmTypes.Block
	data, err := c.SubscribeNewBlock()
	if err != nil {
		fmt.Printf("Error subscribing to height %v error %v", height+1, err)
		return nil, nil, "", err
	}
	switch t := data.(type) {
	case tmTypes.EventDataNewBlock:
		blockDetails = t.Block
	default:
		c.Logger.Info("No cases matched")
	}

	// TODO ensure block.height == height+1 OR Subscribe to Height+1

	// extract votes from response
	preCommits := blockDetails.LastCommit.Precommits

	// extract signs from votes
	valSigs := helper.GetSigs(preCommits)

	// extract chainID
	chainID = blockDetails.ChainID

	// return
	return preCommits, valSigs, chainID, nil
}

// DispatchCheckpoint prepares the data required for mainchain checkpoint submission
// and sends a transaction to mainchain
func (c *Checkpointer) DispatchCheckpoint(height int64, txBytes tmTypes.Tx, start uint64, end uint64) error {
	c.Logger.Debug("Preparing checkpoint to be pushed on chain")

	// get votes
	votes, sigs, chainID, err := fetchVotes(height, c.httpClient, c.Logger)
	if err != nil {
		return err
	}
	c.Logger.Debug("Fetched votes and signatures", "votes", votes, "sigs", sigs, "chainID", chainID)

	// current child block from contract
	currentChildBlock, err := c.contractConnector.CurrentChildBlock()
	if err != nil {
		return err
	}
	c.Logger.Debug("Fetched current child block", "CurrChildBlk", currentChildBlock)

	// fetch current proposer from heimdall
	validatorAddress := ethCommon.BytesToAddress(helper.GetPubKey().Address().Bytes())
	var proposer hmtypes.Validator

	// fetch latest start block from heimdall via rest query
	response, err := FetchFromAPI(c.cliCtx, CurrentProposerURL)
	if err != nil {
		c.Logger.Error("Failed to get current proposer through rest")
		return err
	}

	// get proposer from response
	if err := json.Unmarshal(response.Result, &proposer); err != nil {
		c.Logger.Error("Error unmarshalling validator", "error", err)
		return err
	}

	// check if we are current proposer
	if !bytes.Equal(proposer.Signer.Bytes(), validatorAddress.Bytes()) {
		return errors.New("We are not proposer, aborting dispatch to mainchain")
	} else {
		c.Logger.Info("We are proposer! Validating if checkpoint needs to be pushed", "commitedLastBlock", currentChildBlock, "startBlock", start)
		// check if we need to send checkpoint or not
		if ((currentChildBlock + 1) == start) || (currentChildBlock == 0 && start == 0) {
			c.Logger.Info("Checkpoint Valid", "startBlock", start)
			c.contractConnector.SendCheckpoint(helper.GetVoteBytes(votes, chainID), sigs, txBytes[authTypes.PulpHashLength:])
		} else if currentChildBlock > start {
			c.Logger.Info("Start block does not match, checkpoint already sent", "commitedLastBlock", currentChildBlock, "startBlock", start)
		} else if currentChildBlock > end {
			c.Logger.Info("Checkpoint already sent", "commitedLastBlock", currentChildBlock, "startBlock", start)
		} else {
			c.Logger.Info("No need to send checkpoint")
		}
	}
	return nil
}

// WaitForOneEvent subscribes to a websocket event for the given
// event time and returns upon receiving it one time, or
// when the timeout duration has expired.
//
// This handles subscribing and unsubscribing under the hood
func (c *Checkpointer) WaitForOneEvent(tx tmTypes.Tx, evtTyp string) (tmTypes.TMEventData, error) {
	const subscriber = "helpers"
	ctx, cancel := context.WithTimeout(context.Background(), CommitTimeout)
	defer cancel()

	query := tmTypes.EventQueryTxFor(tx).String()

	// register for the next event of this type
	eventCh, err := c.httpClient.Subscribe(ctx, subscriber, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to subscribe")
	}

	// make sure to unregister after the test is over
	defer c.httpClient.UnsubscribeAll(ctx, subscriber)

	select {
	case event := <-eventCh:
		return event.Data.(tmTypes.TMEventData), nil
	case <-ctx.Done():
		return nil, errors.New("timed out waiting for event")
	}
}

// SubscribeNewBlock subscribes to a new block
func (c *Checkpointer) SubscribeNewBlock() (tmTypes.TMEventData, error) {
	const subscriber = "helpers"
	ctx, cancel := context.WithTimeout(context.Background(), CommitTimeout)
	defer cancel()

	// register for the next event of this type
	eventCh, err := c.httpClient.Subscribe(ctx, subscriber, tmTypes.QueryForEvent(tmTypes.EventNewBlock).String())
	if err != nil {
		return nil, errors.Wrap(err, "failed to subscribe")
	}
	// unsubscribe everything
	defer c.httpClient.UnsubscribeAll(ctx, subscriber)

	select {
	case event := <-eventCh:
		return event.Data.(tmTypes.TMEventData), nil
	case <-ctx.Done():
		return nil, errors.New("timed out waiting for event")
	}
}

// TODO create a generic event subscriber
// func (c *Checkpointer) SubscribeToEvent(query tmquery.Query) (tmTypes.TMEventData, error) {
// 	const subscriber = "helpers"
// 	ctx, cancel := context.WithTimeout(context.Background(), CommitTimeout)
// 	defer cancel()

// 	// register for the next event of this type
// 	eventCh, err := c.httpClient.Subscribe(ctx, subscriber, query.String())
// 	if err != nil {
// 		return nil, errors.Wrap(err, "failed to subscribe")
// 	}

// 	// make sure to unregister after the test is over
// 	defer c.httpClient.UnsubscribeAll(ctx, subscriber)
// 	select {
// 	case event := <-eventCh:
// 		return event.Data.(tmTypes.TMEventData), nil
// 	case <-ctx.Done():
// 		return nil, errors.New("timed out waiting for event")
// 	}
// }
