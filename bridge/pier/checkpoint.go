package pier

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tendermint/libs/common"
	httpClient "github.com/tendermint/tendermint/rpc/client"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/helper"
	hmtypes "github.com/maticnetwork/heimdall/types"
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
	// contract caller
	contractConnector helper.ContractCaller
	// tx encoder
	txEncoder authTypes.TxBuilder

	// cli context
	cliCtx cliContext.CLIContext
	// queue connector
	queueConnector *QueueConnector
	// http client to subscribe to
	httpClient *httpClient.HTTP
}

// NewCheckpointer returns new service object
func NewCheckpointer(cdc *codec.Codec, queueConnector *QueueConnector, httpClient *httpClient.HTTP) *Checkpointer {
	// create logger
	logger := Logger.With("module", HeimdallCheckpointer)

	// cli context
	cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	cliCtx.BroadcastMode = client.BroadcastAsync
	cliCtx.TrustNode = true

	contractCaller, err := helper.NewContractCaller()
	if err != nil {
		logger.Error("Error while getting root chain instance", "error", err)
		panic(err)
	}

	// creating checkpointer object
	checkpointer := &Checkpointer{
		storageClient:     getBridgeDBInstance(viper.GetString(BridgeDBFlag)),
		HeaderChannel:     make(chan *types.Header),
		contractConnector: contractCaller,
		txEncoder:         authTypes.NewTxBuilderFromCLI().WithTxEncoder(helper.GetTxEncoder()).WithChainID(helper.GetGenesisDoc().ChainID),

		cliCtx:         cliCtx,
		queueConnector: queueConnector,
		httpClient:     httpClient,
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

	// start header process
	go c.startHeaderProcess(headerCtx)

	// subscribe to new head
	subscription, err := c.contractConnector.MaticChainClient.SubscribeNewHead(ctx, c.HeaderChannel)
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
			if isProposer(c.cliCtx) {
				header, err := c.contractConnector.MaticChainClient.HeaderByNumber(ctx, nil)
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

	// get state
	var expectedCheckpointState *ContractCheckpoint
	var bufferedCheckpoint *HeimdallCheckpoint
	var committedCheckpoint *HeimdallCheckpoint

	var wg sync.WaitGroup
	wg.Add(3)

	c.Logger.Debug("Collecting contract and heimdall checkpoint state")

	go func() {
		defer wg.Done()
		expectedCheckpointState, _ = c.nextExpectedCheckpoint(newHeader.Number.Uint64())
	}()

	go func() {
		defer wg.Done()
		bufferedCheckpoint, _ = c.fetchBufferedCheckpoint()
	}()

	go func() {
		defer wg.Done()
		committedCheckpoint, _ = c.fetchCommittedCheckpoint()
	}()

	// wait for state collection
	wg.Wait()

	c.Logger.Info("Fetched contract and heimdall states")

	// check last contract details
	if expectedCheckpointState == nil {
		c.Logger.Error("Error fetching details from contract")
		return
	}

	// contract details fetched
	c.Logger.Debug("Contract details fetched",
		"newStart", expectedCheckpointState.newStart,
		"newEnd", expectedCheckpointState.newEnd,
		"currentHeaderNumber", expectedCheckpointState.currentHeaderBlock.number,
		"currentStart", expectedCheckpointState.currentHeaderBlock.start,
		"currentEnd", expectedCheckpointState.currentHeaderBlock.end,
	)

	// TODO
	if committedCheckpoint == nil {
	}

	//
	// Check if ACK is needed
	//

	// buffer checkpoint log
	if bufferedCheckpoint == nil {
		c.Logger.Debug("Buffer not found")
	} else if bufferedCheckpoint.start != 0 {
		c.Logger.Debug("Checkpoint found in buffer",
			"start", bufferedCheckpoint.start,
			"end", bufferedCheckpoint.end,
		)
	}

	if bufferedCheckpoint != nil &&
		expectedCheckpointState != nil &&
		expectedCheckpointState.currentHeaderBlock.start == bufferedCheckpoint.start {

		// expected checkpoint state
		c.Logger.Debug("Sending ACK",
			"bufferedCheckpointEnd", bufferedCheckpoint.end,
			"contractStart", bufferedCheckpoint.start,
		)

		// // calculate header index
		// headerNumber := expectedCheckpointState.currentHeaderBlock.Sub(
		// 	expectedCheckpointState.currentHeaderBlock,
		// 	big.NewInt(int64(helper.GetConfig().ChildBlockInterval)),
		// )

		// if err := c.broadcastACK(expectedCheckpointState.currentHeaderBlock.number.Uint64()); err != nil {
		// 	c.Logger.Error("Error while sending ACK", "Error", err.Error())
		// }
	}

	//
	// Send checkpoint if valid
	//

	start := expectedCheckpointState.newStart
	end := expectedCheckpointState.newEnd
	if err := c.sendCheckpointToHeimdall(start, end); err != nil {
		c.Logger.Error("Error while sending checkpoint", "error", err)
	}
}

// fetched contract checkpoint state and returns the next probable checkpoint that needs to be sent
func (c *Checkpointer) nextExpectedCheckpoint(latestChildBlock uint64) (*ContractCheckpoint, error) {
	// fetch current header block from mainchain contract
	_currentHeaderBlock, err := c.contractConnector.CurrentHeaderBlock()
	if err != nil {
		c.Logger.Error("Error while fetching current header block number from rootchain", "error", err)
		return nil, err
	}

	// current header block
	currentHeaderBlockNumber := big.NewInt(0).SetUint64(_currentHeaderBlock)

	// get header info
	// currentHeaderBlock = currentHeaderBlock.Sub(currentHeaderBlock, helper.GetConfig().ChildBlockInterval)
	_, currentStart, currentEnd, lastCheckpointTime, _, err := c.contractConnector.GetHeaderInfo(currentHeaderBlockNumber.Uint64())
	if err != nil {
		c.Logger.Error("Error while fetching current header block object from rootchain", "error", err)
		return nil, err
	}

	//
	// find next start/end
	//

	var start, end uint64
	start = currentEnd

	// add 1 if start > 0
	if start > 0 {
		start = start + 1
	}

	// get diff
	diff := latestChildBlock - start + 1

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

		c.Logger.Debug("Calculating checkpoint eligibility",
			"latest", latestChildBlock,
			"start", start,
			"end", end,
		)
	}

	// Handle when block producers go down
	if end == 0 || end == start || (0 < diff && diff < helper.GetConfig().AvgCheckpointLength) {
		c.Logger.Debug("Fetching last header block to calculate time")

		currentTime := time.Now().Unix()
		defaultForcePushInterval := helper.GetConfig().MaxCheckpointLength * 2 // in seconds (1024 * 2 seconds)
		if currentTime-int64(lastCheckpointTime) > int64(defaultForcePushInterval) {
			end = latestChildBlock
			c.Logger.Info("Force push checkpoint",
				"currentTime", currentTime,
				"lastCheckpointTime", lastCheckpointTime,
				"defaultForcePushInterval", defaultForcePushInterval,
				"start", start,
				"end", end,
			)
		}
	}

	// if end == 0 || start >= end {
	// 	c.Logger.Info("Waiting for 256 blocks or invalid start end formation", "start", start, "end", end)
	// 	return nil, errors.New("Invalid start end formation")
	// }

	return NewContractCheckpoint(start, end, &HeaderBlock{
		start:  currentStart,
		end:    currentEnd,
		number: currentHeaderBlockNumber,
	}), nil
}

// fetch checkpoint present in buffer from heimdall
func (c *Checkpointer) fetchBufferedCheckpoint() (*HeimdallCheckpoint, error) {
	c.Logger.Info("Fetching checkpoint in buffer")

	_checkpoint, err := c.fetchCheckpoint(GetHeimdallServerEndpoint(BufferedCheckpointURL))
	if err != nil {
		return nil, err
	}

	bufferedCheckpoint := NewHeimdallCheckpoint(_checkpoint.StartBlock, _checkpoint.EndBlock)
	return bufferedCheckpoint, nil
}

// fetches latest committed checkpoint from heimdall
func (c *Checkpointer) fetchCommittedCheckpoint() (*HeimdallCheckpoint, error) {
	c.Logger.Info("Fetching last committed checkpoint")

	_checkpoint, err := c.fetchCheckpoint(GetHeimdallServerEndpoint(LatestCheckpointURL))
	if err != nil {
		return nil, err
	}

	return NewHeimdallCheckpoint(_checkpoint.StartBlock, _checkpoint.EndBlock), nil
}

// fetches checkpoint from given URL
func (c *Checkpointer) fetchCheckpoint(url string) (checkpoint hmtypes.CheckpointBlockHeader, err error) {
	response, err := FetchFromAPI(c.cliCtx, url)
	if err != nil {
		return checkpoint, err
	}

	if err := json.Unmarshal(response.Result, &checkpoint); err != nil {
		c.Logger.Error("Error unmarshalling checkpoint", "error", err)
		return checkpoint, err
	}

	return checkpoint, nil
}

// fetches initial genesis rewardroothash
func (c *Checkpointer) fetchInitialRewardRoot() (rewardRootHash hmtypes.HeimdallHash, err error) {
	c.Logger.Info("Sending Rest call to Get Initial RewardRootHash")
	response, err := FetchFromAPI(c.cliCtx, GetHeimdallServerEndpoint(InitialRewardRootURL))
	if err != nil {
		c.Logger.Error("Error Fetching rewardroothash from HeimdallServer ", "error", err)
		return rewardRootHash, err
	}

	if err := json.Unmarshal(response.Result, &rewardRootHash); err != nil {
		c.Logger.Error("Error unmarshalling rewardroothash received from Heimdall Server", "error", err)
		return rewardRootHash, err
	}
	return rewardRootHash, nil
}

// broadcast checkpoint
func (c *Checkpointer) sendCheckpointToHeimdall(start uint64, end uint64) error {
	if end == 0 || start >= end {
		c.Logger.Info("Waiting for blocks or invalid start end formation", "start", start, "end", end)
		return errors.New("No new valid checkpoint, yet. Waiting for more blocks or time")
	}

	// Get root hash
	root, err := checkpoint.GetHeaders(start, end)
	if err != nil {
		return err
	}

	rewardRootHash := hmtypes.ZeroHeimdallHash
	c.Logger.Info("Check if it is firstcheckpoint, if so Get InitialRewardRoot from HeimdallServer")
	if start == uint64(0) {
		if rewardRootHash, err = c.fetchInitialRewardRoot(); err != nil {
			c.Logger.Info("Error while fetching initial reward root hash from HeimdallServer", "err", err)
			return err
		}
	} else {
		// Get Latest Reward Root Hash through rest call
		c.Logger.Info("Sending Request to HeimdallServer to fetch latest committed Checkpoint")
		latestCheckpoint, err := c.fetchCheckpoint(GetHeimdallServerEndpoint(LatestCheckpointURL))
		if err != nil {
			c.Logger.Info("Error while fetching Latest Checkpoint from heimdallserver", "err", err)
			return err
		}
		rewardRootHash = latestCheckpoint.RewardRootHash
	}

	c.Logger.Info("Creating and broadcasting new checkpoint",
		"start", start,
		"end", end,
		"root", hmtypes.BytesToHeimdallHash(root),
		"rewardRoot", rewardRootHash,
	)

	// create and send checkpoint message
	msg := checkpoint.NewMsgCheckpointBlock(
		hmtypes.BytesToHeimdallAddress(helper.GetAddress()),
		start,
		end,
		hmtypes.BytesToHeimdallHash(root),
		rewardRootHash,
		uint64(time.Now().Unix()),
	)

	// return broadcast to heimdall
	if err := c.queueConnector.BroadcastToHeimdall(msg); err != nil {
		return err
	}

	// wait for checkpoint to confirm and commit
	go c.commitCheckpoint(start, end)

	return nil
}

// broadcastACK broadcasts ack for a checkpoint to heimdall
// func (c *Checkpointer) broadcastACK(headerID uint64) error {
// 	// create and send checkpoint ACK message
// 	msg := checkpoint.NewMsgCheckpointAck(hmtypes.BytesToHeimdallAddress(helper.GetAddress()), headerID)
// 	// broadcast ack
// 	return c.queueConnector.BroadcastToHeimdall(msg)
// }

// wait for heimdall checkpoint tx to get confirmed and dispatch checkpoint
func (c *Checkpointer) commitCheckpoint(startBlock uint64, endBlock uint64) {
	// create tag query
	var tags []string
	tags = append(tags, fmt.Sprintf("start-block='%v'", startBlock))
	tags = append(tags, fmt.Sprintf("end-block='%v'", endBlock))
	tags = append(tags, "action='checkpoint'")

	// handler
	handler := func() bool {
		// search txs
		txs, err := helper.SearchTxs(c.cliCtx, c.cliCtx.Codec, tags, 1, 1) // first page, 1 limit
		if err != nil {
			c.Logger.Error("Error while searching txs", "error", err)
			return false
		}

		// loop through tx
		for _, tx := range txs {
			txHash, err := hex.DecodeString(tx.TxHash)
			if err != nil {
				c.Logger.Error("Error while searching txs", "error", err)
			} else {
				if err := c.dispatchCheckpoint(tx.Height, txHash, startBlock, endBlock); err == nil {
					return true
				}
			}
		}

		return false
	}

	// wait for sometime / poll
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan bool)

	// wait for 2 min
	timeout := time.AfterFunc(2*time.Minute, func() {
		quit <- true
	})
	defer timeout.Stop()

	// loop
	for {
		select {
		case <-ticker.C:
			if ok := handler(); ok {
				ticker.Stop()
				return
			}
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

// dispatchCheckpoint prepares the data required for mainchain checkpoint submission
// and sends a transaction to mainchain
func (c *Checkpointer) dispatchCheckpoint(height int64, txHash []byte, start uint64, end uint64) error {
	c.Logger.Debug("Preparing checkpoint to be pushed on chain")

	// proof
	tx, err := helper.QueryTxWithProof(c.cliCtx, txHash)
	if err != nil {
		return err
	}

	// get votes
	votes, sigs, chainID, err := FetchVotes(height, c.httpClient)
	if err != nil {
		return err
	}

	// current child block from contract
	currentChildBlock, err := c.contractConnector.GetLastChildBlock()
	if err != nil {
		return err
	}
	c.Logger.Debug("Fetched current child block", "currentChildBlock", currentChildBlock)

	// fetch current proposer from heimdall
	validatorAddress := ethCommon.BytesToAddress(helper.GetPubKey().Address().Bytes())
	var proposer hmtypes.Validator

	// fetch latest start block from heimdall via rest query
	response, err := FetchFromAPI(c.cliCtx, GetHeimdallServerEndpoint(CurrentProposerURL))
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
			c.contractConnector.SendCheckpoint(helper.GetVoteBytes(votes, chainID), sigs, tx.Tx[authTypes.PulpHashLength:])
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
