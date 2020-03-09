package processor

import (
	"encoding/json"

	"github.com/maticnetwork/bor/accounts/abi"
	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/helper"
	topupTypes "github.com/maticnetwork/heimdall/topup/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
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
	fp.queueConnector.Server.RegisterTask("sendTopUpFeeToHeimdall", fp.sendTopUpFeeToHeimdall)
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

		fp.Logger.Info("✅ Received task to send topup to heimdall",
			"event", eventName,
			"validatorId", event.ValidatorId,
			"Fee", event.Fee,
		)

		// create msg checkpoint ack message
		msg := topupTypes.NewMsgTopup(helper.GetFromAddress(fp.cliCtx), event.ValidatorId.Uint64(), hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()), uint64(vLog.Index))

		// return broadcast to heimdall
		if err := fp.txBroadcaster.BroadcastToHeimdall(msg); err != nil {
			fp.Logger.Error("Error while broadcasting TopupFee msg to heimdall", "error", err)
			return err
		}
	}
	return nil
}
