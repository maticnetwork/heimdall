package processor

import (
	"encoding/json"

	"github.com/RichardKnop/machinery/v1/tasks"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/maticnetwork/bor/accounts/abi"
	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/helper"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// StakingProcessor - process staking related events
type StakingProcessor struct {
	BaseProcessor
	stakingInfoAbi *abi.ABI
}

// NewStakingProcessor - add  abi to staking processor
func NewStakingProcessor(stakingInfoAbi *abi.ABI) *StakingProcessor {
	stakingProcessor := &StakingProcessor{
		stakingInfoAbi: stakingInfoAbi,
	}
	return stakingProcessor
}

// Start starts new block subscription
func (sp *StakingProcessor) Start() error {
	sp.Logger.Info("Starting")
	return nil
}

// RegisterTasks - Registers staking tasks with machinery
func (sp *StakingProcessor) RegisterTasks() {
	sp.Logger.Info("Registering staking related tasks")
	sp.queueConnector.Server.RegisterTask("sendValidatorJoinToHeimdall", sp.sendValidatorJoinToHeimdall)
	sp.queueConnector.Server.RegisterTask("sendUnstakeInitToHeimdall", sp.sendUnstakeInitToHeimdall)
	sp.queueConnector.Server.RegisterTask("sendStakeUpdateToHeimdall", sp.sendStakeUpdateToHeimdall)
	sp.queueConnector.Server.RegisterTask("sendSignerChangeToHeimdall", sp.sendSignerChangeToHeimdall)
}

func (sp *StakingProcessor) sendValidatorJoinToHeimdall(eventName string, logBytes string) error {
	var vLog = types.Log{}
	if err := json.Unmarshal([]byte(logBytes), &vLog); err != nil {
		sp.Logger.Error("Error while unmarshalling event from rootchain", "error", err)
		return err
	}

	event := new(stakinginfo.StakinginfoStaked)
	if err := helper.UnpackLog(sp.stakingInfoAbi, event, eventName, &vLog); err != nil {
		sp.Logger.Error("Error while parsing event", "name", eventName, "error", err)
	} else {
		signerPubKey := event.SignerPubkey
		if len(signerPubKey) == 64 {
			signerPubKey = util.AppendPrefix(signerPubKey)
		}
		if isOld, _ := sp.isOldTx(sp.cliCtx, vLog.TxHash.String(), uint64(vLog.Index)); isOld {
			sp.Logger.Info("Ignoring task to send validatorjoin to heimdall as already processed",
				"event", eventName,
				"validatorID", event.ValidatorId,
				"activationEpoch", event.ActivationEpoch,
				"amount", event.Amount,
				"totalAmount", event.Total,
				"SignerPubkey", hmTypes.NewPubKey(signerPubKey).String(),
				"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				"logIndex", uint64(vLog.Index),
			)
			return nil
		}

		// if account doesn't exists Retry with delay for topup to process first.
		if _, err := util.GetAccount(sp.cliCtx, hmTypes.HeimdallAddress(event.Signer)); err != nil {
			sp.Logger.Info(
				"Heimdall Account doesn't exist. Retrying validator-join after 10 seconds",
				"event", eventName,
				"signer", event.Signer,
			)
			return tasks.NewErrRetryTaskLater("account doesn't exist", util.RetryTaskDelay)
		}

		sp.Logger.Info(
			"✅ Received task to send validatorjoin to heimdall",
			"event", eventName,
			"validatorID", event.ValidatorId,
			"activationEpoch", event.ActivationEpoch,
			"amount", event.Amount,
			"totalAmount", event.Total,
			"SignerPubkey", hmTypes.NewPubKey(signerPubKey).String(),
			"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			"logIndex", uint64(vLog.Index),
		)

		// msg validator exit
		msg := stakingTypes.NewMsgValidatorJoin(
			hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
			event.ValidatorId.Uint64(),
			hmTypes.NewPubKey(signerPubKey),
			hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			uint64(vLog.Index),
		)

		// return broadcast to heimdall
		if err := sp.txBroadcaster.BroadcastToHeimdall(msg); err != nil {
			sp.Logger.Error("Error while broadcasting unstakeInit to heimdall", "validatorId", event.ValidatorId.Uint64(), "error", err)
			return err
		}
	}
	return nil
}

func (sp *StakingProcessor) sendUnstakeInitToHeimdall(eventName string, logBytes string) error {
	var vLog = types.Log{}
	if err := json.Unmarshal([]byte(logBytes), &vLog); err != nil {
		sp.Logger.Error("Error while unmarshalling event from rootchain", "error", err)
		return err
	}

	event := new(stakinginfo.StakinginfoUnstakeInit)
	if err := helper.UnpackLog(sp.stakingInfoAbi, event, eventName, &vLog); err != nil {
		sp.Logger.Error("Error while parsing event", "name", eventName, "error", err)
	} else {
		if isOld, _ := sp.isOldTx(sp.cliCtx, vLog.TxHash.String(), uint64(vLog.Index)); isOld {
			sp.Logger.Info("Ignoring task to send unstakeinit to heimdall as already processed",
				"event", eventName,
				"validator", event.User,
				"validatorID", event.ValidatorId,
				"deactivatonEpoch", event.DeactivationEpoch,
				"amount", event.Amount,
				"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				"logIndex", uint64(vLog.Index),
			)
			return nil
		}

		sp.Logger.Info(
			"✅ Received task to send unstake-init to heimdall",
			"event", eventName,
			"validator", event.User,
			"validatorID", event.ValidatorId,
			"deactivatonEpoch", event.DeactivationEpoch,
			"amount", event.Amount,
			"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			"logIndex", uint64(vLog.Index),
		)

		// msg validator exit
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
	return nil
}

func (sp *StakingProcessor) sendStakeUpdateToHeimdall(eventName string, logBytes string) error {
	var vLog = types.Log{}
	if err := json.Unmarshal([]byte(logBytes), &vLog); err != nil {
		sp.Logger.Error("Error while unmarshalling event from rootchain", "error", err)
		return err
	}

	event := new(stakinginfo.StakinginfoStakeUpdate)
	if err := helper.UnpackLog(sp.stakingInfoAbi, event, eventName, &vLog); err != nil {
		sp.Logger.Error("Error while parsing event", "name", eventName, "error", err)
	} else {
		if isOld, _ := sp.isOldTx(sp.cliCtx, vLog.TxHash.String(), uint64(vLog.Index)); isOld {
			sp.Logger.Info("Ignoring task to send unstakeinit to heimdall as already processed",
				"event", eventName,
				"validatorID", event.ValidatorId,
				"newAmount", event.NewAmount,
				"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				"logIndex", uint64(vLog.Index),
			)
			return nil
		}
		sp.Logger.Info(
			"✅ Received task to send stake-update to heimdall",
			"event", eventName,
			"validatorID", event.ValidatorId,
			"newAmount", event.NewAmount,
			"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			"logIndex", uint64(vLog.Index),
		)

		// msg validator exit
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
	return nil
}

func (sp *StakingProcessor) sendSignerChangeToHeimdall(eventName string, logBytes string) error {
	var vLog = types.Log{}
	if err := json.Unmarshal([]byte(logBytes), &vLog); err != nil {
		sp.Logger.Error("Error while unmarshalling event from rootchain", "error", err)
		return err
	}

	event := new(stakinginfo.StakinginfoSignerChange)
	if err := helper.UnpackLog(sp.stakingInfoAbi, event, eventName, &vLog); err != nil {
		sp.Logger.Error("Error while parsing event", "name", eventName, "error", err)
	} else {
		newSignerPubKey := event.SignerPubkey
		if len(newSignerPubKey) == 64 {
			newSignerPubKey = util.AppendPrefix(newSignerPubKey)
		}

		if isOld, _ := sp.isOldTx(sp.cliCtx, vLog.TxHash.String(), uint64(vLog.Index)); isOld {
			sp.Logger.Info("Ignoring task to send unstakeinit to heimdall as already processed",
				"event", eventName,
				"validatorID", event.ValidatorId,
				"NewSignerPubkey", hmTypes.NewPubKey(newSignerPubKey).String(),
				"oldSigner", event.OldSigner.Hex(),
				"newSigner", event.NewSigner.Hex(),
				"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				"logIndex", uint64(vLog.Index),
			)
			return nil
		}
		sp.Logger.Info(
			"✅ Received task to send signer-change to heimdall",
			"event", eventName,
			"validatorID", event.ValidatorId,
			"NewSignerPubkey", hmTypes.NewPubKey(newSignerPubKey).String(),
			"oldSigner", event.OldSigner.Hex(),
			"newSigner", event.NewSigner.Hex(),
			"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			"logIndex", uint64(vLog.Index),
		)

		// signer change
		msg := stakingTypes.NewMsgSignerUpdate(
			hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
			event.ValidatorId.Uint64(),
			hmTypes.NewPubKey(newSignerPubKey),
			hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			uint64(vLog.Index),
		)

		// return broadcast to heimdall
		if err := sp.txBroadcaster.BroadcastToHeimdall(msg); err != nil {
			sp.Logger.Error("Error while broadcasting signerChainge to heimdall", "validatorId", event.ValidatorId.Uint64(), "error", err)
			return err
		}
	}
	return nil
}

// isOldTx  checks if tx is already processed or not
func (sp *StakingProcessor) isOldTx(cliCtx cliContext.CLIContext, txHash string, logIndex uint64) (bool, error) {
	queryParam := map[string]interface{}{
		"txhash":   txHash,
		"logindex": logIndex,
	}

	endpoint := helper.GetHeimdallServerEndpoint(util.StakingTxStatusURL)
	url, err := util.CreateURLWithQuery(endpoint, queryParam)
	if err != nil {
		sp.Logger.Error("Error in creating url", "endpoint", endpoint, "error", err)
		return false, err
	}

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
