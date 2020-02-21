package processor

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/maticnetwork/bor/accounts/abi"
	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/bridge/pier"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/sethu/util"
	"github.com/streadway/amqp"

	sdk "github.com/cosmos/cosmos-sdk/types"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	hmTypes "github.com/maticnetwork/heimdall/types"

	ethCommon "github.com/maticnetwork/bor/common"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
)

// CheckpointProcessor
type CheckpointProcessor struct {
	BaseProcessor

	abi *abi.ABI
}

func NewCheckpointProcessor() *CheckpointProcessor {
	contractCaller, err := helper.NewContractCaller()
	if err != nil {
		panic(err)
	}

	checkpointProcessor := &CheckpointProcessor{
		abi: &contractCaller.RootChainABI,
	}

	return checkpointProcessor
}

// Start starts new block subscription
func (cp *CheckpointProcessor) Start() error {
	cp.Logger.Info("Starting Processor", "name", cp.String())

	amqpMsgs, err := cp.queueConnector.ConsumeMsg(util.CheckpointQueueName)
	if err != nil {
		cp.Logger.Info("error consuming msg", "error", err)
	}
	// handle all amqp messages
	for amqpMsg := range amqpMsgs {
		cp.Logger.Info("Received Message from checkpoint queue", "Msg - ", string(amqpMsg.Body), "AppID", amqpMsg.AppId)

		// Process msg
		go cp.ProcessMsg(amqpMsg)
	}
	return nil
}

// ProcessMsg - identify checkpoint msg type and call the handler
func (cp *CheckpointProcessor) ProcessMsg(amqpMsg amqp.Delivery) {
	cp.Logger.Info("Processing msg from queue", "sender", amqpMsg.AppId)

	switch amqpMsg.AppId {
	case "maticchain":
		var header = types.Header{}
		if err := header.UnmarshalJSON(amqpMsg.Body); err != nil {
			cp.Logger.Error("Error while unmarshalling the header block", "error", err)
			amqpMsg.Reject(false)
		}
		cp.Logger.Info("Processing new header", "headerNumber", header.Number)
		if err := cp.HandleHeaderBlock(&header); err != nil {
			cp.Logger.Error("Error while processing the header block", "error", err)
			amqpMsg.Reject(false)
		}
	case "heimdall":
		var event = sdk.StringEvent{}
		if err := json.Unmarshal(amqpMsg.Body, &event); err != nil {
			cp.Logger.Error("Error unmarshalling event from heimdall", "error", err)
			amqpMsg.Reject(false)
		}
		if err := cp.HandleCheckpointConfirmation(event); err != nil {
			cp.Logger.Error("Error while processing checkpoint event from heimdall", "error", err)
			amqpMsg.Reject(false)
		}
	case "rootchain":
		var log = types.Log{}
		if err := json.Unmarshal(amqpMsg.Body, &log); err != nil {
			cp.Logger.Error("Error while unmarshalling event from rootchain", "error", err)
			amqpMsg.Reject(false)
		}
		if err := cp.HandleCheckpointAck(amqpMsg.Type, log); err != nil {
			cp.Logger.Error("Error while processing checkpoint ack event from rootchain", "error", err)
			amqpMsg.Reject(false)
		}
	default:
		cp.Logger.Info("AppID mismatch", "appId", amqpMsg.AppId)
	}

	amqpMsg.Ack(false)
}

// HandleHeaderBlock - Broadcasts next checkpoint to heimdall
func (cp *CheckpointProcessor) HandleHeaderBlock(newHeader *types.Header) error {
	// expectedCheckpointState, err := cp.nextExpectedCheckpoint(newHeader.Number.Uint64())
	// if err != nil {
	// 	cp.Logger.Error("Error while calculate next expected checkpoint", "error", err)
	// 	return err
	// }

	// start := expectedCheckpointState.newStart
	// end := expectedCheckpointState.newEnd
	if err := cp.sendCheckpointToHeimdall(0, 62); err != nil {
		cp.Logger.Error("Error while sending checkpoint", "error", err)
		return err
	}
	return nil
}

// HandleCheckpointConfirmation - handles checkpoint confirmation event from heimdall.
// 1. check if this checkpoint has to be submitted to rootchain
// 2. create transaction to be submmitted to rootchain
// 3. Broadcast transaction to rootchain
func (cp *CheckpointProcessor) HandleCheckpointConfirmation(event sdk.StringEvent) error {
	cp.Logger.Error("Processed event successfullly", "eventtype", event.Type)
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
	cp.sendCheckpointToRootchain(startBlock, endBlock)
	return nil
}

func (cp *CheckpointProcessor) HandleCheckpointAck(eventName string, vLog types.Log) error {
	cp.Logger.Info("Processing checkpoint ack event", "vlog", vLog, "eventName", eventName)
	event := new(rootchain.RootchainNewHeaderBlock)
	cp.Logger.Info("Processing checkpoint ack event", "abi", cp.abi.Events)
	if err := helper.UnpackLog(cp.abi, event, eventName, &vLog); err != nil {
		cp.Logger.Error("Error while parsing event", "name", eventName, "error", err)
	} else {
		cp.Logger.Info(
			"⬜ New event found",
			"event", eventName,
			"start", event.Start,
			"end", event.End,
			"reward", event.Reward,
			"root", "0x"+hex.EncodeToString(event.Root[:]),
			"proposer", event.Proposer.Hex(),
			"headerNumber", event.HeaderBlockId,
		)

		// create msg checkpoint ack message
		msg := checkpointTypes.NewMsgCheckpointAck(helper.GetFromAddress(cp.cliCtx), event.HeaderBlockId.Uint64(), hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()), uint64(vLog.Index))
		// return broadcast to heimdall
		if isBroadcasted := cp.txBroadcaster.BroadcastToHeimdall(msg); isBroadcasted == false {
			cp.Logger.Error("Error while broadcasting checkpoint to heimdall", "error", "sf")
		}
	}
	return nil
}

// HandleCheckpointConfirmation - handles checkpoint confirmation event from heimdall.
// 1. check if this checkpoint has to be submitted to rootchain
// 2. create transaction to be submmitted to rootchain
// 3. Broadcast transaction to rootchain
func (cp *CheckpointProcessor) HandleCheckpointAckConfirmation(event sdk.StringEvent) error {
	cp.Logger.Error("Processed event successfullly", "eventtype", event.Type)
	return nil
}

// // nextExpectedCheckpoint - fetched contract checkpoint state and returns the next probable checkpoint that needs to be sent
// func (cp *CheckpointProcessor) nextExpectedCheckpoint(latestChildBlock uint64) (*ContractCheckpoint, error) {
// 	// fetch current header block from mainchain contract
// 	_currentHeaderBlock, err := cp.contractConnector.CurrentHeaderBlock()
// 	if err != nil {
// 		c.Logger.Error("Error while fetching current header block number from rootchain", "error", err)
// 		return nil, err
// 	}

// 	// current header block
// 	currentHeaderBlockNumber := big.NewInt(0).SetUint64(_currentHeaderBlock)

// 	// get header info
// 	// currentHeaderBlock = currentHeaderBlock.Sub(currentHeaderBlock, helper.GetConfig().ChildBlockInterval)
// 	_, currentStart, currentEnd, lastCheckpointTime, _, err := c.contractConnector.GetHeaderInfo(currentHeaderBlockNumber.Uint64())
// 	if err != nil {
// 		c.Logger.Error("Error while fetching current header block object from rootchain", "error", err)
// 		return nil, err
// 	}

// 	//
// 	// find next start/end
// 	//

// 	var start, end uint64
// 	start = currentEnd

// 	// add 1 if start > 0
// 	if start > 0 {
// 		start = start + 1
// 	}

// 	// get diff
// 	diff := latestChildBlock - start + 1

// 	// process if diff > 0 (positive)
// 	if diff > 0 {
// 		expectedDiff := diff - diff%helper.GetConfig().AvgCheckpointLength
// 		if expectedDiff > 0 {
// 			expectedDiff = expectedDiff - 1
// 		}

// 		// cap with max checkpoint length
// 		if expectedDiff > helper.GetConfig().MaxCheckpointLength-1 {
// 			expectedDiff = helper.GetConfig().MaxCheckpointLength - 1
// 		}

// 		// get end result
// 		end = expectedDiff + start

// 		c.Logger.Debug("Calculating checkpoint eligibility",
// 			"latest", latestChildBlock,
// 			"start", start,
// 			"end", end,
// 		)
// 	}

// 	// Handle when block producers go down
// 	if end == 0 || end == start || (0 < diff && diff < helper.GetConfig().AvgCheckpointLength) {
// 		c.Logger.Debug("Fetching last header block to calculate time")

// 		currentTime := time.Now().UTC().Unix()
// 		defaultForcePushInterval := helper.GetConfig().MaxCheckpointLength * 2 // in seconds (1024 * 2 seconds)
// 		if currentTime-int64(lastCheckpointTime) > int64(defaultForcePushInterval) {
// 			end = latestChildBlock
// 			c.Logger.Info("Force push checkpoint",
// 				"currentTime", currentTime,
// 				"lastCheckpointTime", lastCheckpointTime,
// 				"defaultForcePushInterval", defaultForcePushInterval,
// 				"start", start,
// 				"end", end,
// 			)
// 		}
// 	}

// 	// if end == 0 || start >= end {
// 	// 	c.Logger.Info("Waiting for 256 blocks or invalid start end formation", "start", start, "end", end)
// 	// 	return nil, errors.New("Invalid start end formation")
// 	// }

// 	return NewContractCheckpoint(start, end, &HeaderBlock{
// 		start:  currentStart,
// 		end:    currentEnd,
// 		number: currentHeaderBlockNumber,
// 	}), nil
// }

// sendCheckpointToHeimdall - creates checkpoint msg and broadcasts to heimdall
func (cp *CheckpointProcessor) sendCheckpointToHeimdall(start uint64, end uint64) error {
	if end == 0 || start >= end {
		cp.Logger.Info("Waiting for blocks or invalid start end formation", "start", start, "end", end)
		return errors.New("No new valid checkpoint, yet. Waiting for more blocks or time")
	}

	// Get root hash
	root, err := checkpointTypes.GetHeaders(start, end)
	if err != nil {
		return err
	}
	cp.Logger.Info("Root hash calculated", "rootHash", root)
	accountRootHash := hmTypes.ZeroHeimdallHash
	//Get DividendAccountRoot from HeimdallServer
	if accountRootHash, err = cp.fetchDividendAccountRoot(); err != nil {
		cp.Logger.Info("Error while fetching initial account root hash from HeimdallServer", "err", err)
		return err
	}

	cp.Logger.Info("✅Creating and broadcasting new checkpoint",
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
		uint64(time.Now().UTC().Unix()),
	)

	// return broadcast to heimdall
	if isBroadcasted := cp.txBroadcaster.BroadcastToHeimdall(msg); isBroadcasted == false {
		cp.Logger.Error("Error while broadcasting checkpoint to heimdall", "error", "sf")
	}

	return nil
}

// verify event on heimdall and fetch checkpoint details
func (cp *CheckpointProcessor) sendCheckpointToRootchain(startBlock uint64, endBlock uint64) bool {
	// create tag query
	var tags []string
	tags = append(tags, fmt.Sprintf("checkpoint.start-block='%v'", startBlock))
	tags = append(tags, fmt.Sprintf("checkpoint.end-block='%v'", endBlock))
	tags = append(tags, "message.action='checkpoint'")

	// search txs
	searchResult, err := helper.QueryTxsByEvents(cp.cliCtx, tags, 1, 1) // first page, 1 limit
	if err != nil {
		cp.Logger.Error("Error while searching txs", "error", err)
		return false
	}

	// loop through tx
	if searchResult.Count > 0 {
		for _, tx := range searchResult.Txs {
			txHash, err := hex.DecodeString(tx.TxHash)
			if err != nil {
				cp.Logger.Error("Error while searching txs", "error", err)
			} else {
				if err := cp.commitCheckpoint(tx.Height, txHash, startBlock, endBlock); err == nil {
					return true
				}
			}
		}
	}

	return true
}

// dispatchCheckpoint prepares the data required for mainchain checkpoint submission
// and sends a transaction to mainchain
func (cp *CheckpointProcessor) commitCheckpoint(height int64, txHash []byte, start uint64, end uint64) error {
	cp.Logger.Info("Preparing checkpoint to be pushed on chain", "height", height, "txHash", txHash, "start", start, "end", end)

	// proof
	tx, err := helper.QueryTxWithProof(cp.cliCtx, txHash)
	if err != nil {
		return err
	}

	// get votes
	votes, sigs, chainID, err := pier.FetchVotes(height, cp.httpClient)
	if err != nil {
		return err
	}

	// current child block from contract
	// TODO : Uncomment below
	// currentChildBlock, err := cp.contractConnector.GetLastChildBlock()
	// if err != nil {
	// 	return err
	// }
	currentChildBlock := uint64(0)
	cp.Logger.Info("Fetched current child block", "currentChildBlock", currentChildBlock)

	// fetch current proposer from heimdall
	validatorAddress := ethCommon.BytesToAddress(helper.GetPubKey().Address().Bytes())
	var proposer hmTypes.Validator

	// fetch latest start block from heimdall via rest query
	response, err := pier.FetchFromAPI(cp.cliCtx, pier.GetHeimdallServerEndpoint(pier.CurrentProposerURL))
	if err != nil {
		cp.Logger.Error("Failed to get current proposer through rest")
		return err
	}

	// get proposer from response
	if err := json.Unmarshal(response.Result, &proposer); err != nil {
		cp.Logger.Error("Error unmarshalling validator", "error", err)
		return err
	}

	// check if we are current proposer
	if !bytes.Equal(proposer.Signer.Bytes(), validatorAddress.Bytes()) {
		return errors.New("We are not proposer, aborting dispatch to mainchain")
	} else {
		cp.Logger.Info("We are proposer! Validating if checkpoint needs to be pushed", "commitedLastBlock", currentChildBlock, "startBlock", start)
		// check if we need to send checkpoint or not
		if ((currentChildBlock + 1) == start) || (currentChildBlock == 0 && start == 0) {
			cp.Logger.Info("Checkpoint Valid", "startBlock", start)
			cp.contractConnector.SendCheckpoint(helper.GetVoteBytes(votes, chainID), sigs, tx.Tx[authTypes.PulpHashLength:])
		} else if currentChildBlock > start {
			cp.Logger.Info("Start block does not match, checkpoint already sent", "commitedLastBlock", currentChildBlock, "startBlock", start)
		} else if currentChildBlock > end {
			cp.Logger.Info("Checkpoint already sent", "commitedLastBlock", currentChildBlock, "startBlock", start)
		} else {
			cp.Logger.Info("No need to send checkpoint")
		}
	}
	return nil
}

// fetchDividendAccountRoot - fetches dividend accountroothash
func (cp *CheckpointProcessor) fetchDividendAccountRoot() (accountroothash hmTypes.HeimdallHash, err error) {
	cp.Logger.Info("Sending Rest call to Get Dividend AccountRootHash")
	response, err := pier.FetchFromAPI(cp.cliCtx, pier.GetHeimdallServerEndpoint(pier.DividendAccountRootURL))
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
