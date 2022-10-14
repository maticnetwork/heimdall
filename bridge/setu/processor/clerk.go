package processor

import (
	"encoding/hex"
	"encoding/json"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
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
	if err := cp.queueConnector.Server.RegisterTask("sendStateSyncedToHeimdall", cp.sendStateSyncedToHeimdall); err != nil {
		cp.Logger.Error("RegisterTasks | sendStateSyncedToHeimdall", "error", err)
	}
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

		if util.GetBlockHeight(cp.cliCtx) > helper.SpanOverrideBlockHeight && len(event.Data) > helper.MaxStateSyncSize {
			cp.Logger.Info(`Data is too large to process, Resetting to ""`, "data", hex.EncodeToString(event.Data))
			event.Data = hmTypes.HexToHexBytes("")
		} else if len(event.Data) > helper.LegacyMaxStateSyncSize {
			cp.Logger.Info(`Data is too large to process, Resetting to ""`, "data", hex.EncodeToString(event.Data))
			event.Data = hmTypes.HexToHexBytes("")
		}

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

// isOldTx  checks if tx is already processed or not
func (cp *ClerkProcessor) isOldTx(cliCtx cliContext.CLIContext, txHash string, logIndex uint64) (bool, error) {
	queryParam := map[string]interface{}{
		"txhash":   txHash,
		"logindex": logIndex,
	}

	endpoint := helper.GetHeimdallServerEndpoint(util.ClerkTxStatusURL)
	url, err := util.CreateURLWithQuery(endpoint, queryParam)
	if err != nil {
		cp.Logger.Error("Error in creating url", "endpoint", endpoint, "error", err)
		return false, err
	}

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
