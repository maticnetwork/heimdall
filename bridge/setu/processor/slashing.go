package processor

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/bor/accounts/abi"
	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/bor/core/types"
	authTypes "github.com/maticnetwork/heimdall/auth/types"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	chainmanagerTypes "github.com/maticnetwork/heimdall/chainmanager/types"
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

// SlashingContext represents slashing context
type SlashingContext struct {
	ChainmanagerParams *chainmanagerTypes.Params
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
	sp.queueConnector.Server.RegisterTask("sendUnjailToHeimdall", sp.sendUnjailToHeimdall)

}

// processSlashLimitEvent - processes slash limit event
func (sp *SlashingProcessor) sendTickToHeimdall(eventBytes string, blockHeight int64) (err error) {
	sp.Logger.Info("Recevied sendTickToHeimdall request", "eventBytes", eventBytes, "blockHeight", blockHeight)
	var event = sdk.StringEvent{}
	if err := json.Unmarshal([]byte(eventBytes), &event); err != nil {
		sp.Logger.Error("Error unmarshalling event from heimdall", "error", err)
		return err
	}

	//Get latestSlashInoBytes from HeimdallServer
	latestSlashInoBytes, err := sp.fetchLatestSlashInoBytes()
	if err != nil {
		sp.Logger.Info("Error while fetching latestSlashInoBytes from HeimdallServer", "err", err)
		return err
	}

	var tickCount uint64
	//Get tickCount from HeimdallServer
	if tickCount, err = sp.fetchTickCount(); err != nil {
		sp.Logger.Info("Error while fetching tick count from HeimdallServer", "err", err)
		return err
	}

	sp.Logger.Info("processing slash-limit event", "eventtype", event.Type)

	sp.Logger.Info("✅ Creating and broadcasting Tick tx",
		"id", tickCount+1,
		"From", hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
		"latestSlashInoBytes", latestSlashInoBytes.String(),
	)

	// create msg Tick message
	msg := slashingTypes.NewMsgTick(
		tickCount+1,
		hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
		latestSlashInoBytes,
	)

	// return broadcast to heimdall
	if err := sp.txBroadcaster.BroadcastToHeimdall(msg, event); err != nil {
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
func (sp *SlashingProcessor) sendTickToRootchain(eventBytes string, blockHeight int64) (err error) {
	sp.Logger.Info("Recevied sendTickToRootchain request", "eventBytes", eventBytes, "blockHeight", blockHeight)
	var event = sdk.StringEvent{}
	if err := json.Unmarshal([]byte(eventBytes), &event); err != nil {
		sp.Logger.Error("Error unmarshalling event from heimdall", "error", err)
		return err
	}

	slashInfoBytes := hmTypes.HexBytes{}
	proposerAddr := hmTypes.ZeroHeimdallAddress
	for _, attr := range event.Attributes {
		if attr.Key == slashingTypes.AttributeKeyProposer {
			proposerAddr = hmTypes.HexToHeimdallAddress(attr.Value)
		}
		if attr.Key == slashingTypes.AttributeKeySlashInfoBytes {
			slashInfoBytes = hmTypes.HexToHexBytes(attr.Value)
		}
	}

	sp.Logger.Info("processing tick confirmation event", "eventtype", event.Type, "slashInfoBytes", slashInfoBytes.String(), "proposer", proposerAddr)
	// TODO - slashing...who should submit tick to rootchain??
	isCurrentProposer, err := util.IsCurrentProposer(sp.cliCtx)
	if err != nil {
		sp.Logger.Error("Error checking isCurrentProposer", "error", err)
		return err
	}

	// Fetch Tick val slashing info
	tickSlashInfoList, err := sp.fetchTickSlashInfoList()
	if err != nil {
		sp.Logger.Error("Error fetching tick slash info list", "error", err)
		return err
	}

	// Validate tickSlashInfoList
	isValidSlashInfo, err := sp.validateTickSlashInfo(tickSlashInfoList, slashInfoBytes)
	if err != nil {
		sp.Logger.Error("Error validating tick slash info list", "error", err)
		return err
	}

	var txHash string
	for _, attr := range event.Attributes {
		if attr.Key == hmTypes.AttributeKeyTxHash {
			txHash = attr.Value
		}
	}

	if isValidSlashInfo && isCurrentProposer {
		txHash := common.FromHex(txHash)
		if err := sp.createAndSendTickToRootchain(blockHeight, txHash, tickSlashInfoList, proposerAddr); err != nil {
			sp.Logger.Error("Error sending tick to rootchain", "error", err)
			return err
		}
	} else {
		sp.Logger.Info("I am not the current proposer or tick already sent or invalid tick data... Ignoring", "eventType", event.Type)
		return nil
	}
	return nil
}

/*
sendTickAckToHeimdall - sends tick ack msg to heimdall
*/
func (sp *SlashingProcessor) sendTickAckToHeimdall(eventName string, logBytes string) error {
	var vLog = types.Log{}
	if err := json.Unmarshal([]byte(logBytes), &vLog); err != nil {
		sp.Logger.Error("Error while unmarshalling event from rootchain", "error", err)
		return err
	}

	event := new(stakinginfo.StakinginfoSlashed)
	if err := helper.UnpackLog(sp.stakingInfoAbi, event, eventName, &vLog); err != nil {
		sp.Logger.Error("Error while parsing event", "name", eventName, "error", err)
	} else {

		if isOld, _ := sp.isOldTx(sp.cliCtx, vLog.TxHash.String(), uint64(vLog.Index), util.SlashingEvent, event); isOld {
			sp.Logger.Info("Ignoring task to send tick ack to heimdall as already processed",
				"event", eventName,
				"tickID", event.Nonce,
				"totalSlashedAmount", event.Amount.Uint64(),
				"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				"logIndex", uint64(vLog.Index),
				"blockNumber", vLog.BlockNumber,
			)
			return nil
		}
		sp.Logger.Info(
			"✅ Received task to send tick-ack to heimdall",
			"event", eventName,
			"tickID", event.Nonce,
			"totalSlashedAmount", event.Amount.Uint64(),
			"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			"logIndex", uint64(vLog.Index),
			"blockNumber", vLog.BlockNumber,
		)

		// TODO - check if i am the proposer of this tick ack or not.

		// create msg checkpoint ack message
		msg := slashingTypes.NewMsgTickAck(helper.GetFromAddress(sp.cliCtx), event.Nonce.Uint64(), event.Amount.Uint64(), hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()), uint64(vLog.Index), vLog.BlockNumber)

		// return broadcast to heimdall
		if err := sp.txBroadcaster.BroadcastToHeimdall(msg, event); err != nil {
			sp.Logger.Error("Error while broadcasting tick-ack to heimdall", "error", err)
			return err
		}
	}
	return nil
}

/*
sendUnjailToHeimdall - sends unjail msg to heimdall
*/
func (sp *SlashingProcessor) sendUnjailToHeimdall(eventName string, logBytes string) error {
	var vLog = types.Log{}
	if err := json.Unmarshal([]byte(logBytes), &vLog); err != nil {
		sp.Logger.Error("Error while unmarshalling event from rootchain", "error", err)
		return err
	}

	event := new(stakinginfo.StakinginfoUnJailed)
	if err := helper.UnpackLog(sp.stakingInfoAbi, event, eventName, &vLog); err != nil {
		sp.Logger.Error("Error while parsing event", "name", eventName, "error", err)
	} else {

		if isOld, _ := sp.isOldTx(sp.cliCtx, vLog.TxHash.String(), uint64(vLog.Index), util.SlashingEvent, event); isOld {
			sp.Logger.Info("Ignoring sending unjail to heimdall as already processed",
				"event", eventName,
				"ValidatorID", event.ValidatorId,
				"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				"logIndex", uint64(vLog.Index),
				"blockNumber", vLog.BlockNumber,
			)
			return nil
		}
		sp.Logger.Info(
			"✅ Received task to send unjail to heimdall",
			"event", eventName,
			"ValidatorID", event.ValidatorId,
			"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			"logIndex", uint64(vLog.Index),
			"blockNumber", vLog.BlockNumber,
		)

		// TODO - check if i am the proposer of unjail or not.

		// msg unjail
		msg := slashingTypes.NewMsgUnjail(
			hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
			event.ValidatorId.Uint64(),
			hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			uint64(vLog.Index),
			vLog.BlockNumber,
		)

		// return broadcast to heimdall
		if err := sp.txBroadcaster.BroadcastToHeimdall(msg, event); err != nil {
			sp.Logger.Error("Error while broadcasting unjail to heimdall", "error", err)
			return err
		}
	}
	return nil
}

// createAndSendTickToRootchain prepares the data required for rootchain tick submission
// and sends a transaction to rootchain
func (sp *SlashingProcessor) createAndSendTickToRootchain(height int64, txHash []byte, slashInfoList []*hmTypes.ValidatorSlashingInfo, proposerAddr hmTypes.HeimdallAddress) error {
	sp.Logger.Info("Preparing tick to be pushed on chain", "height", height, "txHash", hmTypes.BytesToHeimdallHash(txHash))

	// proof
	tx, err := helper.QueryTxWithProof(sp.cliCtx, txHash)
	if err != nil {
		sp.Logger.Error("Error querying tick tx proof", "txHash", txHash)
		return err
	}

	// fetch side txs sigs
	decoder := helper.GetTxDecoder(authTypes.ModuleCdc)
	stdTx, err := decoder(tx.Tx)
	if err != nil {
		sp.Logger.Error("Error while decoding tick tx", "txHash", tx.Tx.Hash(), "error", err)
		return err
	}

	cmsg := stdTx.GetMsgs()[0]
	sideMsg, ok := cmsg.(hmTypes.SideTxMsg)
	if !ok {
		sp.Logger.Error("Invalid side-tx msg", "txHash", tx.Tx.Hash())
		return err
	}

	// side-tx data
	sideTxData := sideMsg.GetSideSignBytes()
	sp.Logger.Info("sideTx data", "sideTxData", hex.EncodeToString(sideTxData))

	// get sigs
	// TODO pass sigs in proper form in `SendTick` for slashing
	// sigs, err := helper.FetchSideTxSigs(sp.httpClient, height, tx.Tx.Hash(), sideTxData)
	// if err != nil {
	// 	sp.Logger.Error("Error fetching votes for tick tx", "height", height)
	// 	return err
	// }

	// send tick to rootchain
	slashingContrext, err := sp.getSlashingContext()
	if err != nil {
		return err
	}
	chainParams := slashingContrext.ChainmanagerParams.ChainParams
	slashManagerAddress := chainParams.SlashManagerAddress.EthAddress()

	// slashmanage instance
	slashManagerInstance, err := sp.contractConnector.GetSlashManagerInstance(slashManagerAddress)
	if err != nil {
		sp.Logger.Info("Error while creating slashmanager instance", "error", err)
		return err
	}

	// TODO pass sigs in proper form in `SendTick` for slashing
	if err := sp.contractConnector.SendTick(sideTxData, nil, slashManagerAddress, slashManagerInstance); err != nil {
		sp.Logger.Info("Error submitting tick to slashManager contract", "error", err)
		return err
	}

	return nil
}

// fetchLatestSlashInoBytes - fetches latest slashInfoBytes
func (sp *SlashingProcessor) fetchLatestSlashInoBytes() (slashInfoBytes hmTypes.HexBytes, err error) {
	sp.Logger.Info("Sending Rest call to Get Latest SlashInfoBytes")
	response, err := helper.FetchFromAPI(sp.cliCtx, helper.GetHeimdallServerEndpoint(util.LatestSlashInfoBytesURL))
	if err != nil {
		sp.Logger.Error("Error Fetching slashInfoBytes from HeimdallServer ", "error", err)
		return slashInfoBytes, err
	}
	sp.Logger.Info("Latest slashInfoBytes fetched")
	if err := json.Unmarshal(response.Result, &slashInfoBytes); err != nil {
		sp.Logger.Error("Error unmarshalling latest slashInfoBytes received from Heimdall Server", "error", err)
		return slashInfoBytes, err
	}
	return slashInfoBytes, nil
}

// fetchTickCount - fetches tick count
func (sp *SlashingProcessor) fetchTickCount() (tickCount uint64, err error) {
	sp.Logger.Info("Sending Rest call to Get Tick count")
	response, err := helper.FetchFromAPI(sp.cliCtx, helper.GetHeimdallServerEndpoint(util.SlashingTickCountURL))
	if err != nil {
		sp.Logger.Error("Error while sending request for tick count", "Error", err)
		return tickCount, err
	}

	if err := json.Unmarshal(response.Result, &tickCount); err != nil {
		sp.Logger.Error("Error unmarshalling tick count data ", "error", err)
		return tickCount, err
	}
	return tickCount, nil
}

// fetchTickSlashInfoList - fetches tick slash Info list
func (sp *SlashingProcessor) fetchTickSlashInfoList() (slashInfoList []*hmTypes.ValidatorSlashingInfo, err error) {
	sp.Logger.Info("Sending Rest call to Get Tick SlashInfo list")
	response, err := helper.FetchFromAPI(sp.cliCtx, helper.GetHeimdallServerEndpoint(util.TickSlashInfoListURL))
	if err != nil {
		sp.Logger.Error("Error Fetching Tick slashInfoList from HeimdallServer ", "error", err)
		return slashInfoList, err
	}
	sp.Logger.Info("Tick SlashInfo List fetched")
	if err := json.Unmarshal(response.Result, &slashInfoList); err != nil {
		sp.Logger.Error("Error unmarshalling tick slashinfo list received from Heimdall Server", "error", err)
		return slashInfoList, err
	}
	return slashInfoList, nil
}

func (sp *SlashingProcessor) validateTickSlashInfo(slashInfoList []*hmTypes.ValidatorSlashingInfo, slashInfoBytes hmTypes.HexBytes) (isValid bool, err error) {
	tickSlashInfoBytes, err := slashingTypes.SortAndRLPEncodeSlashInfos(slashInfoList)
	if err != nil {
		sp.Logger.Error("Error generating tick slashinfo bytes", "error", err)
		return
	}
	// compare tickSlashInfoBytes with slashInfoBytes
	if bytes.Equal(tickSlashInfoBytes, slashInfoBytes.Bytes()) {
		return true, nil
	}
	sp.Logger.Info("SlashingInfoBytes mismatch", "tickSlashInfoBytes", hex.EncodeToString(tickSlashInfoBytes), "slashInfoBytes", slashInfoBytes)
	return false, errors.New("Validation failed. tickSlashInfoBytes mismatch")
}

//
// utils
//

func (sp *SlashingProcessor) getSlashingContext() (*SlashingContext, error) {
	chainmanagerParams, err := util.GetChainmanagerParams(sp.cliCtx)
	if err != nil {
		sp.Logger.Error("Error while fetching chain manager params", "error", err)
		return nil, err
	}

	return &SlashingContext{
		ChainmanagerParams: chainmanagerParams,
	}, nil
}
