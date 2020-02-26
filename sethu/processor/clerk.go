package processor

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ethereum "github.com/maticnetwork/bor"
	"github.com/maticnetwork/bor/core/types"
	clerkTypes "github.com/maticnetwork/heimdall/clerk/types"
	"github.com/maticnetwork/heimdall/contracts/statesender"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/sethu/queue"
	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/streadway/amqp"
)

// ClerkProcessor - sync state/deposit events
type ClerkProcessor struct {
	BaseProcessor
}

// Start starts new block subscription
func (cp *ClerkProcessor) Start() error {
	cp.Logger.Info("Starting Processor")
	amqpMsgs, err := cp.queueConnector.ConsumeMsg(queue.ClerkQueueName)
	if err != nil {
		cp.Logger.Info("Error consuming statesync msg", "error", err)
		panic(err)
	}
	// handle all amqp messages
	for amqpMsg := range amqpMsgs {
		cp.Logger.Info("Received Message from clerk queue", "Msg - ", string(amqpMsg.Body), "AppID", amqpMsg.AppId)
		go cp.ProcessMsg(amqpMsg)
	}
	return nil
}

// ProcessMsg - identify clerk msg type and delegate to msg/event handlers
func (cp *ClerkProcessor) ProcessMsg(amqpMsg amqp.Delivery) {
	cp.Logger.Info("Processing msg", "sender", amqpMsg.AppId)
	switch amqpMsg.AppId {
	case "rootchain":
		var vLog = types.Log{}
		if err := json.Unmarshal(amqpMsg.Body, &vLog); err != nil {
			cp.Logger.Error("Error while unmarshalling event from rootchain", "error", err)
			amqpMsg.Reject(false)
			return
		}
		if err := cp.HandleStateSyncEvent(amqpMsg.Type, &vLog); err != nil {
			cp.Logger.Error("Error while processing Statesync event from rootchain", "error", err)
			amqpMsg.Reject(true)
			return
		}
	case "heimdall":
		var event = sdk.StringEvent{}
		if err := json.Unmarshal(amqpMsg.Body, &event); err != nil {
			cp.Logger.Error("Error unmarshalling event from heimdall", "error", err)
			amqpMsg.Reject(false)
			return
		}
		if err := cp.HandleRecordConfirmation(event); err != nil {
			cp.Logger.Error("Error while processing record event from heimdall", "error", err)
			amqpMsg.Reject(true)
			return
		}
	default:
		cp.Logger.Info("AppID mismatch", "appId", amqpMsg.AppId)
	}
	// send ack
	amqpMsg.Ack(false)
}

// HandleStateSyncEvent - handle state sync event from rootchain
// 1. check if this deposit event has to be broadcasted to heimdall
// 2. create and broadcast  record transaction to heimdall
func (cp *ClerkProcessor) HandleStateSyncEvent(eventName string, vLog *types.Log) error {
	event := new(statesender.StatesenderStateSynced)
	if err := helper.UnpackLog(cp.rootchainAbi, event, eventName, vLog); err != nil {
		cp.Logger.Error("Error while parsing event", "name", eventName, "error", err)
	} else {
		cp.Logger.Debug(
			"⬜ New event found",
			"event", eventName,
			"id", event.Id,
			"contract", event.ContractAddress,
			"data", hex.EncodeToString(event.Data),
		)

		msg := clerkTypes.NewMsgEventRecord(
			hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
			hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			uint64(vLog.Index),
			event.Id.Uint64(),
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
func (cp *ClerkProcessor) HandleRecordConfirmation(event sdk.StringEvent) (err error) {
	cp.Logger.Info("processing record confirmation event", "eventtype", event.Type)
	var recordID uint64
	for _, attr := range event.Attributes {
		if attr.Key == clerkTypes.AttributeKeyRecordID {
			if recordID, err = strconv.ParseUint(attr.Value, 10, 64); err != nil {
				cp.Logger.Error("Error parsing recordId", "eventtype", event.Type)
				return err
			}
			break
		}
	}

	// TODO - query on heimdall for recordID check status.
	if err := cp.commitRecordID(recordID); err != nil {
		cp.Logger.Error("Error commit recordId to maticchain", "recordID", recordID)
		return err
	}
	return nil
}

// broadcastToBor - propose state to bor
func (cp *ClerkProcessor) commitRecordID(stateID uint64) error {
	// encode commit span
	encodedData, err := cp.encodeProposeStateData(stateID)
	if err != nil {
		cp.Logger.Error("Error encoding state data", "recordID", stateID)
		return err
	}
	// get validator address
	stateReceiverAddress := helper.GetStateReceiverAddress()
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
		Logger.Error("Error unpacking tx for commit state", "error", err)
		return nil, err
	}
	// return data
	return data, nil
}
