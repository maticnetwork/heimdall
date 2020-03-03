package processor

import (
	"bytes"
	"encoding/json"

	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/bridge/setu/queue"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/helper"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/streadway/amqp"
)

// StakingProcessor - process staking related events
type StakingProcessor struct {
	BaseProcessor
}

// Start starts new block subscription
func (sp *StakingProcessor) Start() error {
	sp.Logger.Info("Starting staking processor")

	amqpMsgs, err := sp.queueConnector.ConsumeMsg(queue.StakingQueueName)
	if err != nil {
		sp.Logger.Info("error consuming staking msg", "error", err)
		panic(err)
	}
	// handle all amqp messages
	for amqpMsg := range amqpMsgs {
		sp.Logger.Debug("Received Message", "msgBody", string(amqpMsg.Body), "AppID", amqpMsg.AppId)
		go sp.ProcessMsg(amqpMsg)
	}
	return nil
}

// ProcessMsg - identify staking msg type and delegate to msg/event handlers
func (sp *StakingProcessor) ProcessMsg(amqpMsg amqp.Delivery) {
	switch amqpMsg.AppId {
	case "rootchain":
		var vLog = types.Log{}
		if err := json.Unmarshal(amqpMsg.Body, &vLog); err != nil {
			sp.Logger.Error("Error while unmarshalling event from rootchain", "error", err)
			amqpMsg.Reject(false)
			return
		}
		switch amqpMsg.Type {
		case "UnstakeInit":
			if err := sp.processUnstakeInitEvent(amqpMsg.Type, &vLog); err != nil {
				sp.Logger.Error("Error while processing unstakeInit event from rootchain", "error", err)
				amqpMsg.Reject(true)
				return
			}
		case "StakeUpdate":
			if err := sp.processStakeUpdateEvent(amqpMsg.Type, &vLog); err != nil {
				sp.Logger.Error("Error while processing stakeUpdate event from rootchain", "error", err)
				amqpMsg.Reject(true)
				return
			}
		case "SignerChange":
			if err := sp.processSignerChangeEvent(amqpMsg.Type, &vLog); err != nil {
				sp.Logger.Error("Error while processing signerChange event from rootchain", "error", err)
				amqpMsg.Reject(true)
				return
			}
		default:
			sp.Logger.Info("Event Type mismatch", "eventType", amqpMsg.Type)
		}
	default:
		sp.Logger.Info("AppID mismatch", "appId", amqpMsg.AppId)
	}

	// send ack
	amqpMsg.Ack(false)
}

func (sp *StakingProcessor) processUnstakeInitEvent(eventName string, vLog *types.Log) error {
	event := new(stakinginfo.StakinginfoUnstakeInit)
	if err := helper.UnpackLog(sp.rootchainAbi, event, eventName, vLog); err != nil {
		sp.Logger.Error("Error while parsing event", "name", eventName, "error", err)
	} else {
		sp.Logger.Debug(
			"⬜ New event found",
			"event", eventName,
			"validator", event.User,
			"validatorID", event.ValidatorId,
			"deactivatonEpoch", event.DeactivationEpoch,
			"amount", event.Amount,
		)

		// msg validator exit
		if util.IsEventSender(sp.cliCtx, event.ValidatorId.Uint64()) {
			msg := stakingTypes.NewMsgValidatorExit(
				hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
				event.ValidatorId.Uint64(),
				hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				uint64(vLog.Index),
			)

			// return broadcast to heimdall
			if err := sp.txBroadcaster.BroadcastToHeimdall(msg); err != nil {
				sp.Logger.Error("Error while broadcasting unstakeInit to heimdall", "validatorId", event.ValidatorId.Uint64(), "error", err)
				return err
			}
		}
	}
	return nil
}

func (sp *StakingProcessor) processStakeUpdateEvent(eventName string, vLog *types.Log) error {
	event := new(stakinginfo.StakinginfoStakeUpdate)
	if err := helper.UnpackLog(sp.rootchainAbi, event, eventName, vLog); err != nil {
		sp.Logger.Error("Error while parsing event", "name", eventName, "error", err)
	} else {
		sp.Logger.Debug(
			"⬜ New event found",
			"event", eventName,
			"validatorID", event.ValidatorId,
			"newAmount", event.NewAmount,
		)

		// msg validator exit
		if util.IsEventSender(sp.cliCtx, event.ValidatorId.Uint64()) {
			msg := stakingTypes.NewMsgStakeUpdate(
				hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
				event.ValidatorId.Uint64(),
				hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				uint64(vLog.Index),
			)

			// return broadcast to heimdall
			if err := sp.txBroadcaster.BroadcastToHeimdall(msg); err != nil {
				sp.Logger.Error("Error while broadcasting stakeupdate to heimdall", "validatorId", event.ValidatorId.Uint64(), "error", err)
				return err
			}
		}
	}
	return nil
}

func (sp *StakingProcessor) processSignerChangeEvent(eventName string, vLog *types.Log) error {
	event := new(stakinginfo.StakinginfoSignerChange)
	if err := helper.UnpackLog(sp.rootchainAbi, event, eventName, vLog); err != nil {
		sp.Logger.Error("Error while parsing event", "name", eventName, "error", err)
	} else {
		sp.Logger.Debug(
			"⬜ New event found",
			"event", eventName,
			"validatorID", event.ValidatorId,
			"newSigner", event.NewSigner.Hex(),
			"oldSigner", event.OldSigner.Hex(),
		)

		// signer change
		if bytes.Compare(event.NewSigner.Bytes(), helper.GetAddress()) == 0 {
			pubkey := helper.GetPubKey()
			msg := stakingTypes.NewMsgSignerUpdate(
				hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
				event.ValidatorId.Uint64(),
				hmTypes.NewPubKey(pubkey[:]),
				hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				uint64(vLog.Index),
			)

			// return broadcast to heimdall
			if err := sp.txBroadcaster.BroadcastToHeimdall(msg); err != nil {
				sp.Logger.Error("Error while broadcasting signerChainge to heimdall", "validatorId", event.ValidatorId.Uint64(), "error", err)
				return err
			}
		}
	}
	return nil
}
