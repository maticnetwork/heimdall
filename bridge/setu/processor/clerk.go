package processor

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strconv"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ethereum "github.com/maticnetwork/bor"
	"github.com/maticnetwork/bor/accounts/abi"
	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	chainmanagerTypes "github.com/maticnetwork/heimdall/chainmanager/types"
	clerkTypes "github.com/maticnetwork/heimdall/clerk/types"
	"github.com/maticnetwork/heimdall/contracts/statesender"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// ClerkContext for bridge
type ClerkContext struct {
	ChainmanagerParams *chainmanagerTypes.Params
}

// ClerkProcessor - sync state/deposit events
type ClerkProcessor struct {
	BaseProcessor
	stateSenderAbi *abi.ABI
}

// NewClerkProcessor - add statesender abi to clerk processor
func NewClerkProcessor(stateSenderAbi *abi.ABI) *ClerkProcessor {
	clerkProcessor := &ClerkProcessor{
		stateSenderAbi: stateSenderAbi,
	}
	return clerkProcessor
}

// Start starts new block subscription
func (cp *ClerkProcessor) Start() error {
	cp.Logger.Info("Starting")
	return nil
}

// RegisterTasks - Registers clerk related tasks with machinery
func (cp *ClerkProcessor) RegisterTasks() {
	cp.Logger.Info("Registering clerk tasks")
	cp.queueConnector.Server.RegisterTask("sendStateSyncedToHeimdall", cp.sendStateSyncedToHeimdall)
	cp.queueConnector.Server.RegisterTask("sendDepositRecordToMatic", cp.sendDepositRecordToMatic)

}

// HandleStateSyncEvent - handle state sync event from rootchain
// 1. check if this deposit event has to be broadcasted to heimdall
// 2. create and broadcast  record transaction to heimdall
func (cp *ClerkProcessor) sendStateSyncedToHeimdall(eventName string, logBytes string) error {
	var vLog = types.Log{}
	if err := json.Unmarshal([]byte(logBytes), &vLog); err != nil {
		cp.Logger.Error("Error while unmarshalling event from rootchain", "error", err)
		return err
	}

	clerkContext, err := cp.getClerkContext()
	if err != nil {
		return err
	}

	chainParams := clerkContext.ChainmanagerParams.ChainParams

	event := new(statesender.StatesenderStateSynced)
	if err := helper.UnpackLog(cp.stateSenderAbi, event, eventName, &vLog); err != nil {
		cp.Logger.Error("Error while parsing event", "name", eventName, "error", err)
	} else {
		if isOld, _ := cp.isOldTx(cp.cliCtx, vLog.TxHash.String(), uint64(vLog.Index)); isOld {
			cp.Logger.Info("Ignoring task to send deposit to heimdall as already processed",
				"event", eventName,
				"id", event.Id,
				"contract", event.ContractAddress,
				"data", hex.EncodeToString(event.Data),
				"borChainId", chainParams.BorChainID,
				"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				"logIndex", uint64(vLog.Index),
				"blockNumber", vLog.BlockNumber,
			)
			return nil
		}

		cp.Logger.Debug(
			"⬜ New event found",
			"event", eventName,
			"id", event.Id,
			"contract", event.ContractAddress,
			"data", hex.EncodeToString(event.Data),
			"borChainId", chainParams.BorChainID,
			"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			"logIndex", uint64(vLog.Index),
			"blockNumber", vLog.BlockNumber,
		)

		msg := clerkTypes.NewMsgEventRecord(
			hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
			hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			uint64(vLog.Index),
			vLog.BlockNumber,
			event.Id.Uint64(),
			hmTypes.BytesToHeimdallAddress(event.ContractAddress.Bytes()),
			event.Data,
			chainParams.BorChainID,
		)

		// return broadcast to heimdall
		if err := cp.txBroadcaster.BroadcastToHeimdall(msg); err != nil {
			cp.Logger.Error("Error while broadcasting clerk Record to heimdall", "error", err)
			return err
		}
	}
	return nil
}

// HandleRecordConfirmation - handles clerk record confirmation event from heimdall.
// 1. check if this record has to be broadcasted to maticchain
// 2. create and broadcast  record transaction to maticchain
func (cp *ClerkProcessor) sendDepositRecordToMatic(eventBytes string, txHeight int64, txHash string) (err error) {
	var event = sdk.StringEvent{}
	if err := json.Unmarshal([]byte(eventBytes), &event); err != nil {
		cp.Logger.Error("Error unmarshalling event from heimdall", "error", err)
		return err
	}

	cp.Logger.Info("Processing record confirmation event", "eventType", event.Type)
	var recordID uint64
	for _, attr := range event.Attributes {
		if attr.Key == clerkTypes.AttributeKeyRecordID {
			if recordID, err = strconv.ParseUint(attr.Value, 10, 64); err != nil {
				cp.Logger.Error("Error parsing recordId", "eventType", event.Type)
				return err
			}
			break
		}
	}

	// get clerk context
	clerkContext, err := cp.getClerkContext()
	if err != nil {
		return err
	}

	// TODO - query on heimdall for recordID check status.
	if err := cp.commitRecordID(clerkContext, recordID); err != nil {
		cp.Logger.Error("Error commit recordId to maticchain", "recordID", recordID)
		return err
	}
	return nil
}

// broadcastToBor - propose state to bor
func (cp *ClerkProcessor) commitRecordID(clerkContext *ClerkContext, stateID uint64) error {
	// encode commit span
	encodedData, err := cp.encodeProposeStateData(stateID)
	if err != nil {
		cp.Logger.Error("Error encoding state data", "recordID", stateID)
		return err
	}

	// get chain params
	chainParams := clerkContext.ChainmanagerParams.ChainParams

	stateReceiverAddress := chainParams.StateReceiverAddress.EthAddress()
	msg := ethereum.CallMsg{
		To:   &stateReceiverAddress,
		Data: encodedData,
	}
	// return broadcast to maticchain
	if err := cp.txBroadcaster.BroadcastToMatic(msg); err != nil {
		cp.Logger.Error("Error broadcasting record to maticchain", "error", err)
		return err
	}
	return nil
}

// encodeProposeStateData - encodes state data to be proposed to maticchain
func (cp *ClerkProcessor) encodeProposeStateData(stateID uint64) ([]byte, error) {
	// state receiver ABI
	stateReceiverABI := cp.contractConnector.StateReceiverABI
	// commit state
	data, err := stateReceiverABI.Pack("proposeState", big.NewInt(0).SetUint64(stateID))
	if err != nil {
		cp.Logger.Error("Error unpacking tx for commit state", "error", err)
		return nil, err
	}
	// return data
	return data, nil
}

// isOldTx  checks if tx is already processed or not
func (cp *ClerkProcessor) isOldTx(cliCtx cliContext.CLIContext, txHash string, logIndex uint64) (bool, error) {
	queryParam := map[string]interface{}{
		"txhash":   txHash,
		"logindex": logIndex,
	}

	endpoint := helper.GetHeimdallServerEndpoint(util.ClerkTxStatusURL)
	url, err := util.CreateURLWithQuery(endpoint, queryParam)

	res, err := helper.FetchFromAPI(cp.cliCtx, url)
	if err != nil {
		cp.Logger.Error("Error fetching tx status", "url", url, "error", err)
		return false, err
	}

	var status bool
	if err := json.Unmarshal(res.Result, &status); err != nil {
		cp.Logger.Error("Error unmarshalling tx status received from Heimdall Server", "error", err)
		return false, err
	}

	return status, nil
}

//
// utils
//

func (cp *ClerkProcessor) getClerkContext() (*ClerkContext, error) {
	chainmanagerParams, err := util.GetChainmanagerParams(cp.cliCtx)
	if err != nil {
		cp.Logger.Error("Error while fetching chain manager params", "error", err)
		return nil, err
	}

	return &ClerkContext{
		ChainmanagerParams: chainmanagerParams,
	}, nil
}
