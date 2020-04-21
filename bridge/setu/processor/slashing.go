package processor

import (
	"encoding/hex"
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/bor/accounts/abi"
	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/helper"
	slashingTypes "github.com/maticnetwork/heimdall/slashing/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// SlashingProcessor - process slashing related events
type SlashingProcessor struct {
	BaseProcessor
	stakingInfoAbi *abi.ABI
}

// NewSlashingProcessor - add  abi to slashing processor
func NewSlashingProcessor(stakingInfoAbi *abi.ABI) *SlashingProcessor {
	slashingProcessor := &SlashingProcessor{
		stakingInfoAbi: stakingInfoAbi,
	}
	return slashingProcessor
}

// Start starts new block subscription
func (sp *SlashingProcessor) Start() error {
	sp.Logger.Info("Starting")
	return nil
}

// RegisterTasks - Registers slashing related tasks with machinery
func (sp *SlashingProcessor) RegisterTasks() {
	sp.Logger.Info("Registering slashing related tasks")
	sp.queueConnector.Server.RegisterTask("sendTickToHeimdall", sp.sendTickToHeimdall)
	sp.queueConnector.Server.RegisterTask("sendTickToRootchain", sp.sendTickToRootchain)
	sp.queueConnector.Server.RegisterTask("sendTickAckToHeimdall", sp.sendTickAckToHeimdall)

}

// processSlashLimitEvent - processes slash limit event
func (sp *SlashingProcessor) sendTickToHeimdall(eventBytes string, txHeight int64, txHash string) error {
	sp.Logger.Info("Recevied sendTickToHeimdall request", "eventBytes", eventBytes, "txHeight", txHeight, "txHash", txHash)
	var event = sdk.StringEvent{}
	if err := json.Unmarshal([]byte(eventBytes), &event); err != nil {
		sp.Logger.Error("Error unmarshalling event from heimdall", "error", err)
		return err
	}

	// TODO - slashing - Fetch Slashing info hash
	sp.Logger.Info("✅ Creating and broadcasting Tick tx",
		"From", hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
		"slashingInfoHash", hmTypes.ZeroHeimdallHash,
		"index", uint64(2),
	)

	// create msg Tick message
	msg := slashingTypes.NewMsgtick(
		hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
		hmTypes.ZeroHeimdallHash,
		uint64(2),
	)

	// return broadcast to heimdall
	if err := sp.txBroadcaster.BroadcastToHeimdall(msg); err != nil {
		sp.Logger.Error("Error while broadcasting Tick msg to heimdall", "error", err)
		return err
	}
	return nil
}

/*
	sendTickToRootchain - create and submit tick tx to rootchain to slashing faulty validators
	1. Fetch sigs from heimdall using txHash
	2. Fetch slashing info from heimdall via Rest call
	3. Verify if this tick tx is already submitted to rootchain using nonce data
	4. create tick tx and submit to rootchain
*/
func (sp *SlashingProcessor) sendTickToRootchain(eventBytes string, txHeight int64, txHash string) error {
	sp.Logger.Info("Recevied sendTickToRootchain request", "eventBytes", eventBytes, "txHeight", txHeight, "txHash", txHash)
	var event = sdk.StringEvent{}
	if err := json.Unmarshal([]byte(eventBytes), &event); err != nil {
		sp.Logger.Error("Error unmarshalling event from heimdall", "error", err)
		return err
	}

	sp.Logger.Info("processing tick confirmation event", "eventtype", event.Type)
	// TODO - slashing...who should submit tick to rootchain??
	isCurrentProposer, err := util.IsCurrentProposer(sp.cliCtx)
	if err != nil {
		sp.Logger.Error("Error checking isCurrentProposer in CheckpointConfirmation handler", "error", err)
		return err
	}

	// TODO - replace below nonce variable with actual slash tx nonce
	nonce := uint64(10)
	shouldSend, err := sp.shouldSendTickToRootchain(nonce)
	if err != nil {
		return err
	}

	if shouldSend && isCurrentProposer {
		txHash, err := hex.DecodeString(txHash)
		if err != nil {
			sp.Logger.Error("Error decoding txHash while sending Tick to rootchain", "txHash", txHash, "error", err)
			return err
		}
		if err := sp.createAndSendTickToRootchain(txHeight, txHash); err != nil {
			sp.Logger.Error("Error sending tick to rootchain", "error", err)
			return err
		}
	} else {
		sp.Logger.Info("I am not the current proposer or tick already sent. Ignoring", "eventType", event.Type)
		return nil
	}
	return nil
}

/*
sendTickAckToHeimdall - sends tick ack msg to heimdall
*/
func (sp *SlashingProcessor) sendTickAckToHeimdall(eventName string, logBytes string) error {
	var log = types.Log{}
	if err := json.Unmarshal([]byte(logBytes), &log); err != nil {
		sp.Logger.Error("Error while unmarshalling event from rootchain", "error", err)
		return err
	}

	event := new(stakinginfo.StakinginfoSlashed)
	if err := helper.UnpackLog(sp.stakingInfoAbi, event, eventName, &log); err != nil {
		sp.Logger.Error("Error while parsing event", "name", eventName, "error", err)
	} else {
		sp.Logger.Info(
			"✅ Received task to send tick-ack to heimdall",
			"event", eventName,
			"totalSlashedAmount", event.Amount,
			"txHash", hmTypes.BytesToHeimdallHash(log.TxHash.Bytes()),
			"logIndex", uint64(log.Index),
		)

		// TODO - check if this ack is already processed on heimdall or not.
		// TODO - check if i am the proposer of this ack or not.

		// create msg checkpoint ack message
		msg := slashingTypes.NewMsgtickAck(helper.GetFromAddress(sp.cliCtx), hmTypes.BytesToHeimdallHash(log.TxHash.Bytes()), uint64(log.Index))

		// return broadcast to heimdall
		if err := sp.txBroadcaster.BroadcastToHeimdall(msg); err != nil {
			sp.Logger.Error("Error while broadcasting tick-ack to heimdall", "error", err)
			return err
		}
	}
	return nil
}

// shouldSendTickToRootchain - verifies if this tick is already submitted to rootchain
func (sp *SlashingProcessor) shouldSendTickToRootchain(tickNonce uint64) (shouldSend bool, err error) {
	/*
		1. Fetch latest tick nonce processed on rootchain.
		2.

	*/

	return
}

// createAndSendTickToRootchain prepares the data required for rootchain tick submission
// and sends a transaction to rootchain
func (sp *SlashingProcessor) createAndSendTickToRootchain(height int64, txHash []byte) error {

	return nil
}
