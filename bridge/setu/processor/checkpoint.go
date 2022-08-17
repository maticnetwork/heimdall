package processor

import (
	"context"
	"encoding/hex"
	"errors"
	"math"
	"math/big"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	jsoniter "github.com/json-iterator/go"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	chainmanagerTypes "github.com/maticnetwork/heimdall/chainmanager/types"
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

	// Rootchain abi
	rootchainAbi *abi.ABI
}

// Result represents single req result
type Result struct {
	Result uint64 `json:"result"`
}

// CheckpointContext represents checkpoint context
type CheckpointContext struct {
	ChainmanagerParams *chainmanagerTypes.Params
	CheckpointParams   *checkpointTypes.Params
}

// NewCheckpointProcessor - add rootchain abi to checkpoint processor
func NewCheckpointProcessor(rootchainAbi *abi.ABI) *CheckpointProcessor {
	return &CheckpointProcessor{
		rootchainAbi: rootchainAbi,
	}
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

	if err := cp.queueConnector.Server.RegisterTask("sendCheckpointToHeimdall", cp.sendCheckpointToHeimdall); err != nil {
		cp.Logger.Error("RegisterTasks | sendCheckpointToHeimdall", "error", err)
	}

	if err := cp.queueConnector.Server.RegisterTask("sendCheckpointToRootchain", cp.sendCheckpointToRootchain); err != nil {
		cp.Logger.Error("RegisterTasks | sendCheckpointToRootchain", "error", err)
	}

	if err := cp.queueConnector.Server.RegisterTask("sendCheckpointAckToHeimdall", cp.sendCheckpointAckToHeimdall); err != nil {
		cp.Logger.Error("RegisterTasks | sendCheckpointAckToHeimdall", "error", err)
	}
}

func (cp *CheckpointProcessor) startPollingForNoAck(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go cp.handleCheckpointNoAck()
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

	isProposer, err := util.IsProposer(cp.cliCtx)
	if err != nil {
		cp.Logger.Error("Error checking isProposer in HeaderBlock handler", "error", err)
		return err
	}

	if isProposer {
		// fetch checkpoint context
		checkpointContext, err := cp.getCheckpointContext()
		if err != nil {
			return err
		}

		// process latest confirmed child block only
		chainmanagerParams := checkpointContext.ChainmanagerParams

		cp.Logger.Debug("no of checkpoint confirmations required", "maticchainTxConfirmations", chainmanagerParams.MaticchainTxConfirmations)

		latestConfirmedChildBlock := header.Number.Uint64() - chainmanagerParams.MaticchainTxConfirmations
		if latestConfirmedChildBlock <= 0 {
			cp.Logger.Error("no of blocks on childchain is less than confirmations required", "childChainBlocks", header.Number.Uint64(), "confirmationsRequired", chainmanagerParams.MaticchainTxConfirmations)
			return errors.New("no of blocks on childchain is less than confirmations required")
		}

		expectedCheckpointState, err := cp.nextExpectedCheckpoint(checkpointContext, latestConfirmedChildBlock)
		if err != nil {
			cp.Logger.Error("Error while calculate next expected checkpoint", "error", err)
			return err
		}

		start := expectedCheckpointState.newStart
		end := expectedCheckpointState.newEnd

		//
		// Check checkpoint buffer
		//
		timeStamp := uint64(time.Now().Unix())
		checkpointBufferTime := uint64(checkpointContext.CheckpointParams.CheckpointBufferTime.Seconds())

		bufferedCheckpoint, err := util.GetBufferedCheckpoint(cp.cliCtx)
		if err != nil {
			cp.Logger.Debug("No buffered checkpoint", "bufferedCheckpoint", bufferedCheckpoint)
		}

		if bufferedCheckpoint != nil && !(bufferedCheckpoint.TimeStamp == 0 || ((timeStamp > bufferedCheckpoint.TimeStamp) && timeStamp-bufferedCheckpoint.TimeStamp >= checkpointBufferTime)) {
			cp.Logger.Info("Checkpoint already exits in buffer", "Checkpoint", bufferedCheckpoint.String())
			return nil
		}

		if err := cp.createAndSendCheckpointToHeimdall(checkpointContext, start, end); err != nil {
			cp.Logger.Error("Error sending checkpoint to heimdall", "error", err)
			return err
		}
	} else {
		cp.Logger.Info("I am not the proposer. skipping newheader", "headerNumber", header.Number)
		return
	}

	return nil
}

// sendCheckpointToRootchain - handles checkpoint confirmation event from heimdall.
// 1. check if i am the current proposer.
// 2. check if this checkpoint has to be submitted to rootchain
// 3. if so, create and broadcast checkpoint transaction to rootchain
func (cp *CheckpointProcessor) sendCheckpointToRootchain(eventBytes string, blockHeight int64) error {
	cp.Logger.Info("Received sendCheckpointToRootchain request", "eventBytes", eventBytes, "blockHeight", blockHeight)

	var event sdk.StringEvent
	if err := jsoniter.ConfigFastest.Unmarshal([]byte(eventBytes), &event); err != nil {
		cp.Logger.Error("Error unmarshalling event from heimdall", "error", err)
		return err
	}

	// var tx = sdk.TxResponse{}
	// if err := jsoniter.Unmarshal([]byte(txBytes), &tx); err != nil {
	// 	cp.Logger.Error("Error unmarshalling txResponse", "error", err)
	// 	return err
	// }

	cp.Logger.Info("processing checkpoint confirmation event", "eventtype", event.Type)

	isCurrentProposer, err := util.IsCurrentProposer(cp.cliCtx)
	if err != nil {
		cp.Logger.Error("Error checking isCurrentProposer in CheckpointConfirmation handler", "error", err)
		return err
	}

	var (
		startBlock uint64
		endBlock   uint64
		txHash     string
	)

	for _, attr := range event.Attributes {
		if attr.Key == checkpointTypes.AttributeKeyStartBlock {
			startBlock, _ = strconv.ParseUint(attr.Value, 10, 64)
		}

		if attr.Key == checkpointTypes.AttributeKeyEndBlock {
			endBlock, _ = strconv.ParseUint(attr.Value, 10, 64)
		}

		if attr.Key == hmTypes.AttributeKeyTxHash {
			txHash = attr.Value
		}
	}

	checkpointContext, err := cp.getCheckpointContext()
	if err != nil {
		return err
	}

	shouldSend, err := cp.shouldSendCheckpoint(checkpointContext, startBlock, endBlock)
	if err != nil {
		return err
	}

	if shouldSend && isCurrentProposer {
		txHash := common.FromHex(txHash)
		if err := cp.createAndSendCheckpointToRootchain(checkpointContext, startBlock, endBlock, blockHeight, txHash); err != nil {
			cp.Logger.Error("Error sending checkpoint to rootchain", "error", err)
			return err
		}
	}

	cp.Logger.Info("I am not the current proposer or checkpoint already sent. Ignoring", "eventType", event.Type)

	return nil
}

// sendCheckpointAckToHeimdall - handles checkpointAck event from rootchain
// 1. create and broadcast checkpointAck msg to heimdall.
func (cp *CheckpointProcessor) sendCheckpointAckToHeimdall(eventName string, checkpointAckStr string) error {
	// fetch checkpoint context
	checkpointContext, err := cp.getCheckpointContext()
	if err != nil {
		return err
	}

	var log = types.Log{}
	if err = jsoniter.ConfigFastest.Unmarshal([]byte(checkpointAckStr), &log); err != nil {
		cp.Logger.Error("Error while unmarshalling event from rootchain", "error", err)
		return err
	}

	event := new(rootchain.RootchainNewHeaderBlock)
	if err = helper.UnpackLog(cp.rootchainAbi, event, eventName, &log); err != nil {
		cp.Logger.Error("Error while parsing event", "name", eventName, "error", err)
	} else {
		checkpointNumber := big.NewInt(0).Div(event.HeaderBlockId, big.NewInt(0).SetUint64(checkpointContext.CheckpointParams.ChildBlockInterval))

		cp.Logger.Info(
			"✅ Received task to send checkpoint-ack to heimdall",
			"event", eventName,
			"start", event.Start,
			"end", event.End,
			"reward", event.Reward,
			"root", "0x"+hex.EncodeToString(event.Root[:]),
			"proposer", event.Proposer.Hex(),
			"checkpointNumber", checkpointNumber,
			"txHash", hmTypes.BytesToHeimdallHash(log.TxHash.Bytes()),
			"logIndex", uint64(log.Index),
		)

		// fetch latest checkpoint
		latestCheckpoint, err := util.GetLatestCheckpoint(cp.cliCtx)
		// event checkpoint is older than or equal to latest checkpoint
		if err == nil && latestCheckpoint != nil && latestCheckpoint.EndBlock >= event.End.Uint64() {
			cp.Logger.Debug("Checkpoint ack is already submitted", "start", event.Start, "end", event.End)
			return nil
		}

		// create msg checkpoint ack message
		msg := checkpointTypes.NewMsgCheckpointAck(
			helper.GetFromAddress(cp.cliCtx),
			checkpointNumber.Uint64(),
			hmTypes.BytesToHeimdallAddress(event.Proposer.Bytes()),
			event.Start.Uint64(),
			event.End.Uint64(),
			event.Root,
			hmTypes.BytesToHeimdallHash(log.TxHash.Bytes()),
			uint64(log.Index),
		)

		// return broadcast to heimdall
		if err = cp.txBroadcaster.BroadcastToHeimdall(msg, event); err != nil {
			cp.Logger.Error("Error while broadcasting checkpoint-ack to heimdall", "error", err)
			return err
		}
	}

	return nil
}

// handleCheckpointNoAck - Checkpoint No-Ack handler
// 1. Fetch latest checkpoint time from rootchain
// 2. check if elapsed time is more than NoAck Wait time.
// 3. Send NoAck to heimdall if required.
func (cp *CheckpointProcessor) handleCheckpointNoAck() {
	// fetch fresh checkpoint context
	checkpointContext, err := cp.getCheckpointContext()
	if err != nil {
		return
	}

	lastCreatedAt, err := cp.getLatestCheckpointTime(checkpointContext)
	if err != nil {
		cp.Logger.Error("Error fetching latest checkpoint time from rootchain", "error", err)
		return
	}

	isNoAckRequired, count := cp.checkIfNoAckIsRequired(checkpointContext, lastCreatedAt)
	if isNoAckRequired {
		var isProposer bool

		if isProposer, err = util.IsInProposerList(cp.cliCtx, count); err != nil {
			cp.Logger.Error("Error checking IsInProposerList while proposing Checkpoint No-Ack ", "error", err)
			return
		}

		// if i am the proposer and NoAck is required, then propose No-Ack
		if isProposer {
			// send Checkpoint No-Ack to heimdall
			if err := cp.proposeCheckpointNoAck(); err != nil {
				cp.Logger.Error("Error proposing Checkpoint No-Ack ", "error", err)
				return
			}
		}
	}
}

// nextExpectedCheckpoint - fetched contract checkpoint state and returns the next probable checkpoint that needs to be sent
func (cp *CheckpointProcessor) nextExpectedCheckpoint(checkpointContext *CheckpointContext, latestChildBlock uint64) (*ContractCheckpoint, error) {
	chainmanagerParams := checkpointContext.ChainmanagerParams
	checkpointParams := checkpointContext.CheckpointParams

	rootChainInstance, err := cp.contractConnector.GetRootChainInstance(chainmanagerParams.ChainParams.RootChainAddress.EthAddress())
	if err != nil {
		return nil, err
	}

	// fetch current header block from mainchain contract
	_currentHeaderBlock, err := cp.contractConnector.CurrentHeaderBlock(rootChainInstance, checkpointParams.ChildBlockInterval)
	if err != nil {
		cp.Logger.Error("Error while fetching current header block number from rootchain", "error", err)
		return nil, err
	}

	// current header block
	currentHeaderBlockNumber := big.NewInt(0).SetUint64(_currentHeaderBlock)

	// get header info
	_, currentStart, currentEnd, lastCheckpointTime, _, err := cp.contractConnector.GetHeaderInfo(currentHeaderBlockNumber.Uint64(), rootChainInstance, checkpointParams.ChildBlockInterval)
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
		expectedDiff := diff - diff%checkpointParams.AvgCheckpointLength
		if expectedDiff > 0 {
			expectedDiff = expectedDiff - 1
		}
		// cap with max checkpoint length
		if expectedDiff > checkpointParams.MaxCheckpointLength-1 {
			expectedDiff = checkpointParams.MaxCheckpointLength - 1
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
	if end == 0 || end == start || (0 < diff && diff < checkpointParams.AvgCheckpointLength) {
		cp.Logger.Debug("Fetching last header block to calculate time")

		currentTime := time.Now().UTC().Unix()
		defaultForcePushInterval := checkpointParams.MaxCheckpointLength * 2 // in seconds (1024 * 2 seconds)

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
func (cp *CheckpointProcessor) createAndSendCheckpointToHeimdall(checkpointContext *CheckpointContext, start uint64, end uint64) error {
	cp.Logger.Debug("Initiating checkpoint to Heimdall", "start", start, "end", end)

	if end == 0 || start >= end {
		cp.Logger.Info("Waiting for blocks or invalid start end formation", "start", start, "end", end)
		return nil
	}

	// get checkpoint params
	checkpointParams := checkpointContext.CheckpointParams

	// Get root hash
	root, err := cp.contractConnector.GetRootHash(start, end, checkpointParams.MaxCheckpointLength)
	if err != nil {
		return err
	}

	cp.Logger.Info("Root hash calculated", "rootHash", hmTypes.BytesToHeimdallHash(root))

	var accountRootHash hmTypes.HeimdallHash
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

	chainParams := checkpointContext.ChainmanagerParams.ChainParams

	// create and send checkpoint message
	msg := checkpointTypes.NewMsgCheckpointBlock(
		hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
		start,
		end,
		hmTypes.BytesToHeimdallHash(root),
		accountRootHash,
		chainParams.BorChainID,
	)

	// return broadcast to heimdall
	if err := cp.txBroadcaster.BroadcastToHeimdall(msg, nil); err != nil {
		cp.Logger.Error("Error while broadcasting checkpoint to heimdall", "error", err)
		return err
	}

	return nil
}

// createAndSendCheckpointToRootchain prepares the data required for rootchain checkpoint submission
// and sends a transaction to rootchain
func (cp *CheckpointProcessor) createAndSendCheckpointToRootchain(checkpointContext *CheckpointContext, start uint64, end uint64, height int64, txHash []byte) error {
	cp.Logger.Info("Preparing checkpoint to be pushed on chain", "height", height, "txHash", hmTypes.BytesToHeimdallHash(txHash), "start", start, "end", end)
	// proof
	tx, err := helper.QueryTxWithProof(cp.cliCtx, txHash)
	if err != nil {
		cp.Logger.Error("Error querying checkpoint tx proof", "txHash", txHash)
		return err
	}

	// fetch side txs sigs
	decoder := helper.GetTxDecoder(authTypes.ModuleCdc)

	stdTx, err := decoder(tx.Tx)
	if err != nil {
		cp.Logger.Error("Error while decoding checkpoint tx", "txHash", tx.Tx.Hash(), "error", err)
		return err
	}

	cmsg := stdTx.GetMsgs()[0]

	sideMsg, ok := cmsg.(hmTypes.SideTxMsg)
	if !ok {
		cp.Logger.Error("Invalid side-tx msg", "txHash", tx.Tx.Hash())
		return err
	}

	// side-tx data
	sideTxData := sideMsg.GetSideSignBytes()

	// get sigs
	sigs, err := helper.FetchSideTxSigs(cp.httpClient, height, tx.Tx.Hash(), sideTxData)
	if err != nil {
		cp.Logger.Error("Error fetching votes for checkpoint tx", "height", height)
		return err
	}

	shouldSend, err := cp.shouldSendCheckpoint(checkpointContext, start, end)
	if err != nil {
		return err
	}

	if shouldSend {
		// chain manager params
		chainParams := checkpointContext.ChainmanagerParams.ChainParams
		// root chain address
		rootChainAddress := chainParams.RootChainAddress.EthAddress()
		// root chain instance
		rootChainInstance, err := cp.contractConnector.GetRootChainInstance(rootChainAddress)
		if err != nil {
			cp.Logger.Info("Error while creating rootchain instance", "error", err)
			return err
		}

		if err := cp.contractConnector.SendCheckpoint(sideTxData, sigs, rootChainAddress, rootChainInstance); err != nil {
			cp.Logger.Info("Error submitting checkpoint to rootchain", "error", err)
			return err
		}
	}

	return nil
}

// fetchDividendAccountRoot - fetches dividend accountroothash
func (cp *CheckpointProcessor) fetchDividendAccountRoot() (accountroothash hmTypes.HeimdallHash, err error) {
	cp.Logger.Info("Sending Rest call to Get Dividend AccountRootHash")

	response, err := helper.FetchFromAPI(cp.cliCtx, helper.GetHeimdallServerEndpoint(util.DividendAccountRootURL))
	if err != nil {
		cp.Logger.Error("Error Fetching accountroothash from HeimdallServer ", "error", err)
		return accountroothash, err
	}

	cp.Logger.Info("Divident account root fetched")

	if err = jsoniter.ConfigFastest.Unmarshal(response.Result, &accountroothash); err != nil {
		cp.Logger.Error("Error unmarshalling accountroothash received from Heimdall Server", "error", err)
		return accountroothash, err
	}

	return accountroothash, nil
}

// fetchLatestCheckpointTime - get latest checkpoint time from rootchain
func (cp *CheckpointProcessor) getLatestCheckpointTime(checkpointContext *CheckpointContext) (int64, error) {
	// get chain params
	chainParams := checkpointContext.ChainmanagerParams.ChainParams
	checkpointParams := checkpointContext.CheckpointParams

	rootChainInstance, err := cp.contractConnector.GetRootChainInstance(chainParams.RootChainAddress.EthAddress())
	if err != nil {
		return 0, err
	}

	// fetch last header number
	lastHeaderNumber, err := cp.contractConnector.CurrentHeaderBlock(rootChainInstance, checkpointParams.ChildBlockInterval)
	if err != nil {
		cp.Logger.Error("Error while fetching current header block number", "error", err)
		return 0, err
	}

	// header block
	_, _, _, createdAt, _, err := cp.contractConnector.GetHeaderInfo(lastHeaderNumber, rootChainInstance, checkpointParams.ChildBlockInterval)
	if err != nil {
		cp.Logger.Error("Error while fetching header block object", "error", err)
		return 0, err
	}

	return int64(createdAt), nil
}

func (cp *CheckpointProcessor) getLastNoAckTime() uint64 {
	response, err := helper.FetchFromAPI(cp.cliCtx, helper.GetHeimdallServerEndpoint(util.LastNoAckURL))
	if err != nil {
		cp.Logger.Error("Error while sending request for last no-ack", "Error", err)
		return 0
	}

	var noAckObject Result
	if err := jsoniter.ConfigFastest.Unmarshal(response.Result, &noAckObject); err != nil {
		cp.Logger.Error("Error unmarshalling no-ack data ", "error", err)
		return 0
	}

	return noAckObject.Result
}

// checkIfNoAckIsRequired - check if NoAck has to be sent or not
func (cp *CheckpointProcessor) checkIfNoAckIsRequired(checkpointContext *CheckpointContext, lastCreatedAt int64) (bool, uint64) {
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

	// checkpoint params
	checkpointParams := checkpointContext.CheckpointParams

	// check if difference between no-ack time and current time
	lastNoAck := cp.getLastNoAckTime()

	lastNoAckTime := time.Unix(int64(lastNoAck), 0)
	// if last no ack == 0 , first no-ack to be sent
	if currentTime.Sub(lastNoAckTime).Seconds() < checkpointParams.CheckpointBufferTime.Seconds() && lastNoAck != 0 {
		cp.Logger.Debug("Cannot send multiple no-ack in short time", "timeDiff", currentTime.Sub(lastNoAckTime).Seconds(), "ExpectedDiff", checkpointParams.CheckpointBufferTime.Seconds())
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
	if err := cp.txBroadcaster.BroadcastToHeimdall(msg, nil); err != nil {
		cp.Logger.Error("Error while broadcasting checkpoint-no-ack to heimdall", "msg", msg, "error", err)
		return err
	}

	cp.Logger.Info("No-ack transaction sent successfully")

	return nil
}

// shouldSendCheckpoint checks if checkpoint with given start,end should be sent to rootchain or not.
func (cp *CheckpointProcessor) shouldSendCheckpoint(checkpointContext *CheckpointContext, start uint64, end uint64) (bool, error) {
	chainmanagerParams := checkpointContext.ChainmanagerParams

	rootChainInstance, err := cp.contractConnector.GetRootChainInstance(chainmanagerParams.ChainParams.RootChainAddress.EthAddress())
	if err != nil {
		cp.Logger.Error("Error while creating rootchain instance", "error", err)
		return false, err
	}

	// current child block from contract
	currentChildBlock, err := cp.contractConnector.GetLastChildBlock(rootChainInstance)
	if err != nil {
		cp.Logger.Error("Error fetching current child block", "currentChildBlock", currentChildBlock, "error", err)
		return false, err
	}

	cp.Logger.Debug("Fetched current child block", "currentChildBlock", currentChildBlock)

	shouldSend := false
	// validate if checkpoint needs to be pushed to rootchain and submit
	cp.Logger.Info("Validating if checkpoint needs to be pushed", "commitedLastBlock", currentChildBlock, "startBlock", start)
	// check if we need to send checkpoint or not
	if ((currentChildBlock + 1) == start) || (currentChildBlock == 0 && start == 0) {
		cp.Logger.Info("Checkpoint Valid", "startBlock", start)

		shouldSend = true
	} else if currentChildBlock > start {
		cp.Logger.Info("Start block does not match, checkpoint already sent", "commitedLastBlock", currentChildBlock, "startBlock", start)
	} else if currentChildBlock > end {
		cp.Logger.Info("Checkpoint already sent", "commitedLastBlock", currentChildBlock, "startBlock", start)
	} else {
		cp.Logger.Info("No need to send checkpoint")
	}

	return shouldSend, nil
}

// Stop stops all necessary go routines
func (cp *CheckpointProcessor) Stop() {
	// cancel No-Ack polling
	cp.cancelNoACKPolling()
}

//
// utils
//

func (cp *CheckpointProcessor) getCheckpointContext() (*CheckpointContext, error) {
	chainmanagerParams, err := util.GetChainmanagerParams(cp.cliCtx)
	if err != nil {
		cp.Logger.Error("Error while fetching chain manager params", "error", err)
		return nil, err
	}

	checkpointParams, err := util.GetCheckpointParams(cp.cliCtx)
	if err != nil {
		cp.Logger.Error("Error while fetching checkpoint params", "error", err)
		return nil, err
	}

	return &CheckpointContext{
		ChainmanagerParams: chainmanagerParams,
		CheckpointParams:   checkpointParams,
	}, nil
}
