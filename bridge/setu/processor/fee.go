package processor

import (
	"encoding/json"

	"github.com/maticnetwork/bor/core/types"
	bankTypes "github.com/maticnetwork/heimdall/bank/types"
	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// FeeProcessor - process fee related events
type FeeProcessor struct {
	BaseProcessor
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

// // ProcessMsg - identify Topup msg type and delegate to msg/event handlers
// func (fp *FeeProcessor) ProcessMsg(amqpMsg amqp.Delivery) {
// 	switch amqpMsg.AppId {
// 	case "rootchain":
// 		var vLog = types.Log{}
// 		if err := json.Unmarshal(amqpMsg.Body, &vLog); err != nil {
// 			fp.Logger.Error("Error while unmarshalling event from rootchain", "error", err)
// 			amqpMsg.Reject(false)
// 			return
// 		}
// 		if err := fp.processTopupFeeEvent(amqpMsg.Type, &vLog); err != nil {
// 			fp.Logger.Error("Error while processing topup event from rootchain", "error", err)
// 			amqpMsg.Reject(true)
// 			return
// 		}
// 	default:
// 		fp.Logger.Info("AppID mismatch", "appId", amqpMsg.AppId)
// 	}

// 	// send ack
// 	amqpMsg.Ack(false)
// }

// processTopupFeeEvent - processes topup fee event
func (fp *FeeProcessor) sendTopUpFeeToHeimdall(eventName string, logBytes string) error {
	var vLog = types.Log{}
	if err := json.Unmarshal([]byte(logBytes), &vLog); err != nil {
		fp.Logger.Error("Error while unmarshalling event from rootchain", "error", err)
		return err
	}

	event := new(stakinginfo.StakinginfoTopUpFee)
	if err := helper.UnpackLog(fp.rootchainAbi, event, eventName, &vLog); err != nil {
		fp.Logger.Error("Error while parsing event", "name", eventName, "error", err)
	} else {

		fp.Logger.Info("✅ Received task to send topup to heimdall",
			"event", eventName,
			"validatorId", event.ValidatorId,
			"Fee", event.Fee,
		)

		// create msg checkpoint ack message
		msg := bankTypes.NewMsgTopup(helper.GetFromAddress(fp.cliCtx), event.ValidatorId.Uint64(), hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()), uint64(vLog.Index))

		// return broadcast to heimdall
		if err := fp.txBroadcaster.BroadcastToHeimdall(msg); err != nil {
			fp.Logger.Error("Error while broadcasting TopupFee msg to heimdall", "error", err)
			return err
		}
	}
	return nil
}
