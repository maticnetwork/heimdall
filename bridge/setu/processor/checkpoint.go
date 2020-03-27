package processor

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	// "github.com/streadway/amqp"

	"github.com/maticnetwork/bor/accounts/abi"
	"github.com/maticnetwork/bor/core/types"
	authTypes "github.com/maticnetwork/heimdall/auth/types"

	// "github.com/maticnetwork/heimdall/bridge/setu/queue"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// CheckpointProcessor - processor for checkpoint queue.
type CheckpointProcessor struct {
	BaseProcessor

	// header listener subscription
	cancelNoACKPolling context.CancelFunc

	// Rootchain instance
	rootChainInstance *rootchain.Rootchain

	// Rootchain abi
	rootchainAbi *abi.ABI
}

// Result represents single req result
type Result struct {
	Result uint64 `json:"result"`
}

// NewCheckpointProcessor - add rootchain abi to checkpoint processor
func NewCheckpointProcessor(rootchainAbi *abi.ABI) *CheckpointProcessor {
	// root chain instance
	rootchainInstance, err := helper.GetRootChainInstance()
	if err != nil {
		panic(err)
	}
	checkpointProcessor := &CheckpointProcessor{
		rootChainInstance: rootchainInstance,
		rootchainAbi:      rootchainAbi,
	}
	return checkpointProcessor
}

// Start - consumes messages from checkpoint queue and call processMsg
func (cp *CheckpointProcessor) Start() error {
	cp.Logger.Info("Starting")
	// no-ack
	ackCtx, cancelNoACKPolling := context.WithCancel(context.Background())
	cp.cancelNoACKPolling = cancelNoACKPolling
	cp.Logger.Info("Start polling for no-ack", "pollInterval", helper.GetConfig().NoACKPollInterval)
	go cp.startPollingForNoAck(ackCtx, helper.GetConfig().NoACKPollInterval)
	return nil
}

// RegisterTasks - Registers checkpoint related tasks with machinery
func (cp *CheckpointProcessor) RegisterTasks() {
	cp.Logger.Info("Registering checkpoint tasks")
	cp.queueConnector.Server.RegisterTask("sendCheckpointToHeimdall", cp.sendCheckpointToHeimdall)
	cp.queueConnector.Server.RegisterTask("sendCheckpointToRootchain", cp.sendCheckpointToRootchain)
	cp.queueConnector.Server.RegisterTask("sendCheckpointAckToHeimdall", cp.sendCheckpointAckToHeimdall)
}

func (cp *CheckpointProcessor) startPollingForNoAck(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	// stop ticker when everything done
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			go cp.HandleCheckpointNoAck()
		case <-ctx.Done():
			cp.Logger.Info("No-ack Polling stopped")
			ticker.Stop()
			return
		}
	}
}

// sendCheckpointToHeimdall - handles headerblock from maticchain
// 1. check if i am the proposer for next checkpoint
// 2. check if checkpoint has to be proposed for given headerblock
// 3. if so, propose checkpoint to heimdall.
func (cp *CheckpointProcessor) sendCheckpointToHeimdall(headerBlockStr string) (err error) {
	var header = types.Header{}
	if err := header.UnmarshalJSON([]byte(headerBlockStr)); err != nil {
		cp.Logger.Error("Error while unmarshalling the header block", "error", err)
		return err
	}

	cp.Logger.Info("Processing new header", "headerNumber", header.Number)
	var isProposer bool
	if isProposer, err = util.IsProposer(cp.cliCtx); err != nil {
		cp.Logger.Error("Error checking isProposer in HeaderBlock handler", "error", err)
		return err
	}

	if isProposer {
		expectedCheckpointState, err := cp.nextExpectedCheckpoint(header.Number.Uint64())
		if err != nil {
			cp.Logger.Error("Error while calculate next expected checkpoint", "error", err)
			return err
		}
		start := expectedCheckpointState.newStart
		end := expectedCheckpointState.newEnd
		// TODO - add a check to see if this checkpoint has to be proposed or not.
		// Fetch latest checkpoint from buffer. if expectedCheckpointState.newStart == start, don't send checkpoint
		if err := cp.createAndSendCheckpointToHeimdall(start, end); err != nil {
			cp.Logger.Error("Error sending checkpoint to heimdall", "error", err)
			return err
		}
	} else {
		cp.Logger.Info("i am not the proposer. skipping newheader", "headerNumber", header.Number)
		return
	}

	return nil
}

// sendCheckpointToRootchain - handles checkpoint confirmation event from heimdall.
// 1. check if i am the current proposer.
// 2. check if this checkpoint has to be submitted to rootchain
// 3. if so, create and broadcast checkpoint transaction to rootchain
func (cp *CheckpointProcessor) sendCheckpointToRootchain(checkpointStr string) error {
	var event = sdk.StringEvent{}
	if err := json.Unmarshal([]byte(checkpointStr), &event); err != nil {
		cp.Logger.Error("Error unmarshalling event from heimdall", "error", err)
		return err
	}

	cp.Logger.Info("processing checkpoint confirmation event", "eventtype", event.Type)
	isCurrentProposer, err := util.IsCurrentProposer(cp.cliCtx)
	if err != nil {
		cp.Logger.Error("Error checking isCurrentProposer in CheckpointConfirmation handler", "error", err)
		return err
	}

	if isCurrentProposer {
		var startBlock uint64
		var endBlock uint64
		for _, attr := range event.Attributes {
			if attr.Key == checkpointTypes.AttributeKeyStartBlock {
				startBlock, _ = strconv.ParseUint(attr.Value, 10, 64)
			}
			if attr.Key == checkpointTypes.AttributeKeyEndBlock {
				endBlock, _ = strconv.ParseUint(attr.Value, 10, 64)
			}
		}
		if err := cp.createAndSendCheckpointToRootchain(startBlock, endBlock); err != nil {
			cp.Logger.Error("Error sending checkpoint to rootchain", "error", err)
			return err
		}
	} else {
		cp.Logger.Info("i am not the current proposer. skipping checkpoint confirmation", "eventType", event.Type)
		return nil
	}
	return nil
}

// sendCheckpointAckToHeimdall - handles checkpointAck event from rootchain
// 1. create and broadcast checkpointAck msg to heimdall.
func (cp *CheckpointProcessor) sendCheckpointAckToHeimdall(eventName string, checkpointAckStr string) error {
	var log = types.Log{}
	if err := json.Unmarshal([]byte(checkpointAckStr), &log); err != nil {
		cp.Logger.Error("Error while unmarshalling event from rootchain", "error", err)
		return err
	}

	event := new(rootchain.RootchainNewHeaderBlock)
	if err := helper.UnpackLog(cp.rootchainAbi, event, eventName, &log); err != nil {
		cp.Logger.Error("Error while parsing event", "name", eventName, "error", err)
	} else {
		cp.Logger.Info(
			"✅ Received task to send checkpoint-ack to heimdall",
			"event", eventName,
			"start", event.Start,
			"end", event.End,
			"reward", event.Reward,
			"root", "0x"+hex.EncodeToString(event.Root[:]),
			"proposer", event.Proposer.Hex(),
			"headerNumber", event.HeaderBlockId,
			"txHash", hmTypes.BytesToHeimdallHash(log.TxHash.Bytes()),
			"logIndex", uint64(log.Index),
		)

		// TODO - check if this ack is already processed on heimdall or not.
		// TODO - check if i am the proposer of this ack or not.

		// create msg checkpoint ack message
		msg := checkpointTypes.NewMsgCheckpointAck(helper.GetFromAddress(cp.cliCtx), event.HeaderBlockId.Uint64(), hmTypes.BytesToHeimdallHash(log.TxHash.Bytes()), uint64(log.Index))

		// return broadcast to heimdall
		if err := cp.txBroadcaster.BroadcastToHeimdall(msg); err != nil {
			cp.Logger.Error("Error while broadcasting checkpoint-ack to heimdall", "error", err)
			return err
		}
	}
	return nil
}

// HandleCheckpointNoAck - Checkpoint No-Ack handler
// 1. Fetch latest checkpoint time from rootchain
// 2. check if elapsed time is more than NoAck Wait time.
// 3. Send NoAck to heimdall if required.
func (cp *CheckpointProcessor) HandleCheckpointNoAck() {
	lastCreatedAt, err := cp.getLatestCheckpointTime()
	if err != nil {
		cp.Logger.Error("Error fetching latest checkpoint time from rootchain", "error", err)
		return
	}

	isNoAckRequired, count := cp.checkIfNoAckIsRequired(lastCreatedAt)

	var isProposer bool
	if isProposer, err = util.IsInProposerList(cp.cliCtx, count); err != nil {
		cp.Logger.Error("Error checking IsInProposerList while proposing Checkpoint No-Ack ", "error", err)
		return
	}

	// if i am the proposer and NoAck is required, then propose No-Ack
	if isNoAckRequired && isProposer {
		// send Checkpoint No-Ack to heimdall
		if err := cp.proposeCheckpointNoAck(); err != nil {
			cp.Logger.Error("Error proposing Checkpoint No-Ack ", "error", err)
			return
		}
	}
}

// nextExpectedCheckpoint - fetched contract checkpoint state and returns the next probable checkpoint that needs to be sent
func (cp *CheckpointProcessor) nextExpectedCheckpoint(latestChildBlock uint64) (*ContractCheckpoint, error) {
	// fetch current header block from mainchain contract
	_currentHeaderBlock, err := cp.contractConnector.CurrentHeaderBlock()
	if err != nil {
		cp.Logger.Error("Error while fetching current header block number from rootchain", "error", err)
		return nil, err
	}
	// current header block
	currentHeaderBlockNumber := big.NewInt(0).SetUint64(_currentHeaderBlock)

	// get header info
	// currentHeaderBlock = currentHeaderBlock.Sub(currentHeaderBlock, helper.GetConfig().ChildBlockInterval)
	_, currentStart, currentEnd, lastCheckpointTime, _, err := cp.contractConnector.GetHeaderInfo(currentHeaderBlockNumber.Uint64())
	if err != nil {
		cp.Logger.Error("Error while fetching current header block object from rootchain", "error", err)
		return nil, err
	}

	// find next start/end
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
		cp.Logger.Debug("Calculating checkpoint eligibility",
			"latest", latestChildBlock,
			"start", start,
			"end", end,
		)
	}

	// Handle when block producers go down
	if end == 0 || end == start || (0 < diff && diff < helper.GetConfig().AvgCheckpointLength) {
		cp.Logger.Debug("Fetching last header block to calculate time")

		currentTime := time.Now().UTC().Unix()
		defaultForcePushInterval := helper.GetConfig().MaxCheckpointLength * 2 // in seconds (1024 * 2 seconds)
		if currentTime-int64(lastCheckpointTime) > int64(defaultForcePushInterval) {
			end = latestChildBlock
			cp.Logger.Info("Force push checkpoint",
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

// sendCheckpointToHeimdall - creates checkpoint msg and broadcasts to heimdall
func (cp *CheckpointProcessor) createAndSendCheckpointToHeimdall(start uint64, end uint64) error {
	cp.Logger.Debug("Initiating checkpoint to Heimdall", "start", start, "end", end)

	if end == 0 || start >= end {
		cp.Logger.Info("Waiting for blocks or invalid start end formation", "start", start, "end", end)
		return nil
	}

	// Get root hash
	root, err := checkpointTypes.GetHeaders(start, end)
	if err != nil {
		return err
	}
	cp.Logger.Info("Root hash calculated", "rootHash", hmTypes.BytesToHeimdallHash(root))
	accountRootHash := hmTypes.ZeroHeimdallHash
	//Get DividendAccountRoot from HeimdallServer
	if accountRootHash, err = cp.fetchDividendAccountRoot(); err != nil {
		cp.Logger.Info("Error while fetching initial account root hash from HeimdallServer", "err", err)
		return err
	}

	cp.Logger.Info("✅ Creating and broadcasting new checkpoint",
		"start", start,
		"end", end,
		"root", hmTypes.BytesToHeimdallHash(root),
		"accountRoot", accountRootHash,
	)

	// create and send checkpoint message
	msg := checkpointTypes.NewMsgCheckpointBlock(
		hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
		start,
		end,
		hmTypes.BytesToHeimdallHash(root),
		accountRootHash,
	)

	// return broadcast to heimdall
	if err := cp.txBroadcaster.BroadcastToHeimdall(msg); err != nil {
		cp.Logger.Error("Error while broadcasting checkpoint to heimdall", "error", err)
		return err
	}

	return nil
}

// verify event on heimdall and fetch checkpoint details
func (cp *CheckpointProcessor) createAndSendCheckpointToRootchain(startBlock uint64, endBlock uint64) error {
	// create tag query
	var tags []string
	tags = append(tags, fmt.Sprintf("checkpoint.start-block='%v'", startBlock))
	tags = append(tags, fmt.Sprintf("checkpoint.end-block='%v'", endBlock))
	tags = append(tags, "message.action='checkpoint'")

	// search txs
	searchResult, err := helper.QueryTxsByEvents(cp.cliCtx, tags, 1, 1) // first page, 1 limit
	if err != nil {
		cp.Logger.Error("Error searching checkpoint txs by events", "startBlock", startBlock, "endBlock", endBlock, "error", err)
		return err
	}

	// loop through tx
	if searchResult.Count > 0 {
		for _, tx := range searchResult.Txs {
			txHash, err := hex.DecodeString(tx.TxHash)
			if err != nil {
				cp.Logger.Error("Error searching checkpoint txs by events", "startBlock", startBlock, "endBlock", endBlock, "error", err)
				return err
			} else {
				if err := cp.commitCheckpoint(tx.Height, txHash, startBlock, endBlock); err != nil {
					cp.Logger.Error("Error commiting checkpoint to rootchain", "startBlock", startBlock, "endBlock", endBlock, "error", err)
					return err
				}
				break
			}
		}
	}
	return nil
}

// dispatchCheckpoint prepares the data required for rootchain checkpoint submission
// and sends a transaction to rootchain
func (cp *CheckpointProcessor) commitCheckpoint(height int64, txHash []byte, start uint64, end uint64) error {
	cp.Logger.Info("Preparing checkpoint to be pushed on chain", "height", height, "txHash", hmTypes.BytesToHeimdallHash(txHash), "start", start, "end", end)

	// proof
	tx, err := helper.QueryTxWithProof(cp.cliCtx, txHash)
	if err != nil {
		cp.Logger.Error("Error querying checkpoint tx proof", "txHash", txHash)
		return err
	}

	// get votes
	votes, sigs, chainID, err := helper.FetchVotes(cp.httpClient, height)
	if err != nil {
		cp.Logger.Error("Error fetching votes for checkpoint tx", "height", height)
		return err
	}

	// current child block from contract
	currentChildBlock, err := cp.contractConnector.GetLastChildBlock()
	if err != nil {
		cp.Logger.Info("Error fetching current child block", "currentChildBlock", currentChildBlock)
		return err
	}
	cp.Logger.Info("Fetched current child block", "currentChildBlock", currentChildBlock)

	// validate if checkpoint needs to be pushed to rootchain and submit
	cp.Logger.Info("We are proposer! Validating if checkpoint needs to be pushed", "commitedLastBlock", currentChildBlock, "startBlock", start)
	// check if we need to send checkpoint or not
	if ((currentChildBlock + 1) == start) || (currentChildBlock == 0 && start == 0) {
		cp.Logger.Info("Checkpoint Valid", "startBlock", start)
		if err := cp.contractConnector.SendCheckpoint(helper.GetVoteBytes(votes, chainID), sigs, tx.Tx[authTypes.PulpHashLength:]); err != nil {
			cp.Logger.Info("Error submitting checkpoint to rootchain", "error", err)
			return err
		}
	} else if currentChildBlock > start {
		cp.Logger.Info("Start block does not match, checkpoint already sent", "commitedLastBlock", currentChildBlock, "startBlock", start)
	} else if currentChildBlock > end {
		cp.Logger.Info("Checkpoint already sent", "commitedLastBlock", currentChildBlock, "startBlock", start)
	} else {
		cp.Logger.Info("No need to send checkpoint")
	}
	return nil
}

// fetchDividendAccountRoot - fetches dividend accountroothash
func (cp *CheckpointProcessor) fetchDividendAccountRoot() (accountroothash hmTypes.HeimdallHash, err error) {
	cp.Logger.Info("Sending Rest call to Get Dividend AccountRootHash")
	response, err := util.FetchFromAPI(cp.cliCtx, util.GetHeimdallServerEndpoint(util.DividendAccountRootURL))
	if err != nil {
		cp.Logger.Error("Error Fetching accountroothash from HeimdallServer ", "error", err)
		return accountroothash, err
	}
	cp.Logger.Info("Divident account root fetched")
	if err := json.Unmarshal(response.Result, &accountroothash); err != nil {
		cp.Logger.Error("Error unmarshalling accountroothash received from Heimdall Server", "error", err)
		return accountroothash, err
	}
	return accountroothash, nil
}

// fetchLatestCheckpointTime - get latest checkpoint time from rootchain
func (cp *CheckpointProcessor) getLatestCheckpointTime() (int64, error) {
	currentHeaderNumber, err := cp.rootChainInstance.CurrentHeaderBlock(nil)
	if err != nil {
		cp.Logger.Error("Error while fetching current header block number", "error", err)
		return 0, err
	}

	// fetch last header number
	lastHeaderNumber := currentHeaderNumber.Uint64() // - helper.GetConfig().ChildBlockInterval
	if lastHeaderNumber == 0 {
		// First checkpoint required
		return 0, err
	}
	// get big int header number
	headerNumber := big.NewInt(0)
	headerNumber.SetUint64(lastHeaderNumber)

	// header block
	headerObject, err := cp.rootChainInstance.HeaderBlocks(nil, headerNumber)
	if err != nil {
		cp.Logger.Error("Error while fetching header block object", "error", err)
		return 0, err
	}
	return headerObject.CreatedAt.Int64(), nil
}

func (cp *CheckpointProcessor) getLastNoAckTime() uint64 {
	response, err := util.FetchFromAPI(cp.cliCtx, util.GetHeimdallServerEndpoint(util.LastNoAckURL))
	if err != nil {
		cp.Logger.Error("Error while sending request for last no-ack", "Error", err)
		return 0
	}

	var noackObject Result
	if err := json.Unmarshal(response.Result, &noackObject); err != nil {
		cp.Logger.Error("Error unmarshalling no-ack data ", "error", err)
		return 0
	}

	return noackObject.Result
}

// checkIfNoAckIsRequired - check if NoAck has to be sent or not
func (cp *CheckpointProcessor) checkIfNoAckIsRequired(lastCreatedAt int64) (bool, uint64) {
	var index float64
	// if last created at ==0 , no checkpoint yet
	if lastCreatedAt == 0 {
		index = 1
	}

	checkpointCreationTime := time.Unix(lastCreatedAt, 0)
	currentTime := time.Now().UTC()
	timeDiff := currentTime.Sub(checkpointCreationTime)
	// check if last checkpoint was < NoACK wait time
	if timeDiff.Seconds() >= helper.GetConfig().NoACKWaitTime.Seconds() && index == 0 {
		index = math.Floor(timeDiff.Seconds() / helper.GetConfig().NoACKWaitTime.Seconds())
	}

	if index == 0 {
		return false, uint64(index)
	}

	params, err := cp.getCheckpointParams()
	if err != nil {
		return false, uint64(index)
	}

	// check if difference between no-ack time and current time
	lastNoAck := cp.getLastNoAckTime()

	lastNoAckTime := time.Unix(int64(lastNoAck), 0)
	timeDiff = currentTime.Sub(lastNoAckTime)
	// if last no ack == 0 , first no-ack to be sent
	if currentTime.Sub(lastNoAckTime).Seconds() < params.CheckpointBufferTime.Seconds() && lastNoAck != 0 {
		cp.Logger.Debug("Cannot send multiple no-ack in short time", "timeDiff", currentTime.Sub(lastNoAckTime).Seconds(), "ExpectedDiff", params.CheckpointBufferTime.Seconds())
		return false, uint64(index)
	}
	return true, uint64(index)
}

// proposeCheckpointNoAck - sends Checkpoint NoAck to heimdall
func (cp *CheckpointProcessor) proposeCheckpointNoAck() (err error) {
	// send NO ACK
	msg := checkpointTypes.NewMsgCheckpointNoAck(
		hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
	)

	// return broadcast to heimdall
	if err := cp.txBroadcaster.BroadcastToHeimdall(msg); err != nil {
		cp.Logger.Error("Error while broadcasting checkpoint-no-ack to heimdall", "error", err)
		return err
	}

	cp.Logger.Info("No-ack transaction sent successfully")
	return nil
}

func (cp *CheckpointProcessor) getCheckpointParams() (*checkpointTypes.Params, error) {
	response, err := util.FetchFromAPI(
		cp.cliCtx,
		util.GetHeimdallServerEndpoint(util.CheckpointParamsURL),
	)

	if err != nil {
		return nil, err
	}

	var params checkpointTypes.Params
	if err := json.Unmarshal(response.Result, &params); err != nil {
		cp.Logger.Error("Error unmarshalling checkpoint params", "error", err)
		return nil, err
	}

	return &params, nil
}

// Stop stops all necessary go routines
func (cp *CheckpointProcessor) Stop() {
	// cancel No-Ack polling
	cp.cancelNoACKPolling()
}
