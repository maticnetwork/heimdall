package processor

import (
	"bytes"
	"encoding/hex"
	"encoding/json"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/bor/accounts/abi"
	"github.com/maticnetwork/bor/core/types"
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
	// TODO - slashing. remove this. just for testing
	// sp.sendTickToHeimdall()
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
func (sp *SlashingProcessor) sendTickToHeimdall(eventBytes string, txHeight int64, txHash string) (err error) {
	sp.Logger.Info("Recevied sendTickToHeimdall request", "eventBytes", eventBytes, "txHeight", txHeight, "txHash", txHash)
	var event = sdk.StringEvent{}
	if err := json.Unmarshal([]byte(eventBytes), &event); err != nil {
		sp.Logger.Error("Error unmarshalling event from heimdall", "error", err)
		return err
	}

	latestSlashInfoHash := hmTypes.ZeroHeimdallHash
	//Get DividendAccountRoot from HeimdallServer
	if latestSlashInfoHash, err = sp.fetchLatestSlashInfoHash(); err != nil {
		sp.Logger.Info("Error while fetching latest slashinfo hash from HeimdallServer", "err", err)
		return err
	}

	sp.Logger.Info("✅ Creating and broadcasting Tick tx",
		"From", hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
		"slashingInfoHash", latestSlashInfoHash,
	)

	// create msg Tick message
	msg := slashingTypes.NewMsgTick(
		hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
		latestSlashInfoHash,
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

	slashInfoHash := hmTypes.ZeroHeimdallHash
	proposerAddr := hmTypes.ZeroHeimdallAddress
	for _, attr := range event.Attributes {
		if attr.Key == slashingTypes.AttributeKeyProposer {
			proposerAddr = hmTypes.HexToHeimdallAddress(attr.Value)
		}
		if attr.Key == slashingTypes.AttributeKeySlashInfoHash {
			slashInfoHash = hmTypes.HexToHeimdallHash(attr.Value)
		}
	}

	sp.Logger.Info("processing tick confirmation event", "eventtype", event.Type, "slashInfoHash", slashInfoHash, "proposer", proposerAddr)
	// TODO - slashing...who should submit tick to rootchain??
	isCurrentProposer, err := util.IsCurrentProposer(sp.cliCtx)
	if err != nil {
		sp.Logger.Error("Error checking isCurrentProposer", "error", err)
		return err
	}

	// Validates tx Height with rootchain contract
	shouldSend, err := sp.shouldSendTickToRootchain(uint64(txHeight))
	if err != nil {
		return err
	}

	// Fetch Tick val slashing info
	tickSlashInfoList, err := sp.fetchTickSlashInfoList()
	if err != nil {
		sp.Logger.Error("Error fetching tick slash info list", "error", err)
		return err
	}

	// Validate tickSlashInfoList
	isValidSlashInfo, err := sp.validateTickSlashInfo(tickSlashInfoList, slashInfoHash)
	if err != nil {
		sp.Logger.Error("Error validating tick slash info list", "error", err)
		return err
	}

	if shouldSend && isValidSlashInfo && isCurrentProposer {
		txHash, err := hex.DecodeString(txHash)
		if err != nil {
			sp.Logger.Error("Error decoding txHash while sending Tick to rootchain", "txHash", txHash, "error", err)
			return err
		}
		if err := sp.createAndSendTickToRootchain(txHeight, txHash, tickSlashInfoList, proposerAddr); err != nil {
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

		if isOld, _ := sp.isOldTx(sp.cliCtx, vLog.TxHash.String(), uint64(vLog.Index)); isOld {
			sp.Logger.Info("Ignoring task to tick ack to heimdall as already processed",
				"event", eventName,
				"totalSlashedAmount", event.Amount,
				"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				"logIndex", uint64(vLog.Index),
			)
			return nil
		}
		sp.Logger.Info(
			"✅ Received task to send tick-ack to heimdall",
			"event", eventName,
			"totalSlashedAmount", event.Amount,
			"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			"logIndex", uint64(vLog.Index),
		)

		// TODO - check if i am the proposer of this ack or not.

		// create msg checkpoint ack message
		msg := slashingTypes.NewMsgTickAck(helper.GetFromAddress(sp.cliCtx), hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()), uint64(vLog.Index))

		// return broadcast to heimdall
		if err := sp.txBroadcaster.BroadcastToHeimdall(msg); err != nil {
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

		if isOld, _ := sp.isOldTx(sp.cliCtx, vLog.TxHash.String(), uint64(vLog.Index)); isOld {
			sp.Logger.Info("Ignoring sending unjail to heimdall as already processed",
				"event", eventName,
				"ValidatorID", event.ValidatorId,
				"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				"logIndex", uint64(vLog.Index),
			)
			return nil
		}
		sp.Logger.Info(
			"✅ Received task to send unjail to heimdall",
			"event", eventName,
			"ValidatorID", event.ValidatorId,
			"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			"logIndex", uint64(vLog.Index),
		)

		// TODO - check if i am the proposer of unjail or not.

		// msg unjail
		msg := slashingTypes.NewMsgUnjail(
			hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
			event.ValidatorId.Uint64(),
			hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			uint64(vLog.Index),
		)

		// return broadcast to heimdall
		if err := sp.txBroadcaster.BroadcastToHeimdall(msg); err != nil {
			sp.Logger.Error("Error while broadcasting unjail to heimdall", "error", err)
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
func (sp *SlashingProcessor) createAndSendTickToRootchain(height int64, txHash []byte, slashInfoList []*hmTypes.ValidatorSlashingInfo, proposerAddr hmTypes.HeimdallAddress) error {
	sp.Logger.Info("Preparing tick to be pushed on chain", "height", height, "txHash", hmTypes.BytesToHeimdallHash(txHash))

	// proof
	tx, err := helper.QueryTxWithProof(sp.cliCtx, txHash)
	if err != nil {
		sp.Logger.Error("Error querying tick tx proof", "txHash", txHash)
		return err
	}

	// get votes
	votes, sigs, chainID, err := helper.FetchVotes(sp.httpClient, height)
	if err != nil {
		sp.Logger.Error("Error fetching votes for tick tx", "height", height)
		return err
	}

	shouldSend, err := sp.shouldSendTickToRootchain(uint64(tx.Height))
	if err != nil {
		return err
	}

	if shouldSend {

		slashingContrext, err := sp.getSlashingContext()
		if err != nil {
			return err
		}

		chainParams := slashingContrext.ChainmanagerParams.ChainParams
		slashManagerAddress := chainParams.SlashManagerAddress.EthAddress()

		// slashmanage instance
		slashManagerInstance, err := sp.contractConnector.GetSlashManagerInstance(slashManagerAddress)
		if err != nil {
			panic(err)
		}

		slashInfoBytes, err := slashingTypes.SortAndRLPEncodeSlashInfos(slashInfoList)
		if err != nil {
			sp.Logger.Error("Error rlp encoding slashInfos", "error", err)
			return err
		}

		if err := sp.contractConnector.SendTick(helper.GetVoteBytes(votes, chainID), sigs, slashInfoBytes, proposerAddr.EthAddress(), slashManagerAddress, slashManagerInstance); err != nil {
			sp.Logger.Info("Error submitting tick to slashManager contract", "error", err)
			return err
		}
	}

	return nil
}

// fetchLatestSlashInfoHash - fetches latest slashInfoHash
func (sp *SlashingProcessor) fetchLatestSlashInfoHash() (slashInfoHash hmTypes.HeimdallHash, err error) {
	sp.Logger.Info("Sending Rest call to Get Latest SlashInfoHash")
	response, err := helper.FetchFromAPI(sp.cliCtx, helper.GetHeimdallServerEndpoint(util.LatestSlashInfoHashURL))
	if err != nil {
		sp.Logger.Error("Error Fetching slashInfoHash from HeimdallServer ", "error", err)
		return slashInfoHash, err
	}
	sp.Logger.Info("Latest SlashInfoHash fetched")
	if err := json.Unmarshal(response.Result, &slashInfoHash); err != nil {
		sp.Logger.Error("Error unmarshalling latest slashinfo hash received from Heimdall Server", "error", err)
		return slashInfoHash, err
	}
	return slashInfoHash, nil
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

func (sp *SlashingProcessor) validateTickSlashInfo(slashInfoList []*hmTypes.ValidatorSlashingInfo, slashInfoHash hmTypes.HeimdallHash) (isValid bool, err error) {
	tickSlashInfoHash, err := slashingTypes.GenerateInfoHash(slashInfoList)
	if err != nil {
		sp.Logger.Error("Error generating tick slashinfo hash", "error", err)
		return
	}
	// compare tickSlashInfoHash with slashInfoHash
	if bytes.Compare(tickSlashInfoHash, slashInfoHash.Bytes()) == 0 {
		return true, nil
	} else {
		sp.Logger.Info("SlashingInfoHash mismatch", "tickSlashInfoHash", tickSlashInfoHash, "slashInfoHash", slashInfoHash)
	}

	return
}

// isOldTx  checks if tx is already processed or not
func (sp *SlashingProcessor) isOldTx(cliCtx cliContext.CLIContext, txHash string, logIndex uint64) (bool, error) {
	queryParam := map[string]interface{}{
		"txhash":   txHash,
		"logindex": logIndex,
	}

	endpoint := helper.GetHeimdallServerEndpoint(util.SlashingTxStatusURL)
	url, err := util.CreateURLWithQuery(endpoint, queryParam)

	res, err := helper.FetchFromAPI(sp.cliCtx, url)
	if err != nil {
		sp.Logger.Error("Error fetching tx status", "url", url, "error", err)
		return false, err
	}

	var status bool
	if err := json.Unmarshal(res.Result, &status); err != nil {
		sp.Logger.Error("Error unmarshalling tx status received from Heimdall Server", "error", err)
		return false, err
	}

	return status, nil
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
