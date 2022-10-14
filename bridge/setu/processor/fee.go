package processor

import (
	"encoding/json"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/helper"
	topupTypes "github.com/maticnetwork/heimdall/topup/types"
	hmTypes "github.com/maticnetwork/heimdall/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// FeeProcessor - process fee related events
type FeeProcessor struct {
	BaseProcessor
	stakingInfoAbi *abi.ABI
}

// NewFeeProcessor - add  abi to clerk processor
func NewFeeProcessor(stakingInfoAbi *abi.ABI) *FeeProcessor {
	feeProcessor := &FeeProcessor{
		stakingInfoAbi: stakingInfoAbi,
	}
	return feeProcessor
}

// Start starts new block subscription
func (fp *FeeProcessor) Start() error {
	fp.Logger.Info("Starting")
	return nil
}

// RegisterTasks - Registers clerk related tasks with machinery
func (fp *FeeProcessor) RegisterTasks() {
	fp.Logger.Info("Registering fee related tasks")
	if err := fp.queueConnector.Server.RegisterTask("sendTopUpFeeToHeimdall", fp.sendTopUpFeeToHeimdall); err != nil {
		fp.Logger.Error("RegisterTasks | sendTopUpFeeToHeimdall", "error", err)
	}
}

// processTopupFeeEvent - processes topup fee event
func (fp *FeeProcessor) sendTopUpFeeToHeimdall(eventName string, logBytes string) error {
	var vLog = types.Log{}
	if err := json.Unmarshal([]byte(logBytes), &vLog); err != nil {
		fp.Logger.Error("Error while unmarshalling event from rootchain", "error", err)
		return err
	}

	event := new(stakinginfo.StakinginfoTopUpFee)
	if err := helper.UnpackLog(fp.stakingInfoAbi, event, eventName, &vLog); err != nil {
		fp.Logger.Error("Error while parsing event", "name", eventName, "error", err)
	} else {
		if isOld, _ := fp.isOldTx(fp.cliCtx, vLog.TxHash.String(), uint64(vLog.Index)); isOld {
			fp.Logger.Info("Ignoring task to send topup to heimdall as already processed",
				"event", eventName,
				"user", event.User,
				"Fee", event.Fee,
				"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				"logIndex", uint64(vLog.Index),
				"blockNumber", vLog.BlockNumber,
			)
			return nil
		}

		fp.Logger.Info("✅ sending topup to heimdall",
			"event", eventName,
			"user", event.User,
			"Fee", event.Fee,
			"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			"logIndex", uint64(vLog.Index),
			"blockNumber", vLog.BlockNumber,
		)

		// create msg checkpoint ack message
		msg := topupTypes.NewMsgTopup(helper.GetFromAddress(fp.cliCtx), hmTypes.BytesToHeimdallAddress(event.User.Bytes()), sdk.NewIntFromBigInt(event.Fee), hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()), uint64(vLog.Index), vLog.BlockNumber)

		// return broadcast to heimdall
		if err := fp.txBroadcaster.BroadcastToHeimdall(msg); err != nil {
			fp.Logger.Error("Error while broadcasting TopupFee msg to heimdall", "error", err)
			return err
		}
	}
	return nil
}

// isOldTx  checks if tx is already processed or not
func (fp *FeeProcessor) isOldTx(cliCtx cliContext.CLIContext, txHash string, logIndex uint64) (bool, error) {
	queryParam := map[string]interface{}{
		"txhash":   txHash,
		"logindex": logIndex,
	}

	endpoint := helper.GetHeimdallServerEndpoint(util.TopupTxStatusURL)
	url, err := util.CreateURLWithQuery(endpoint, queryParam)
	if err != nil {
		fp.Logger.Error("Error in creating url", "endpoint", endpoint, "error", err)
		return false, err
	}

	res, err := helper.FetchFromAPI(fp.cliCtx, url)
	if err != nil {
		fp.Logger.Error("Error fetching tx status", "url", url, "error", err)
		return false, err
	}

	var status bool
	if err := json.Unmarshal(res.Result, &status); err != nil {
		fp.Logger.Error("Error unmarshalling tx status received from Heimdall Server", "error", err)
		return false, err
	}

	return status, nil
}
