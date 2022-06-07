package processor

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/RichardKnop/machinery/v1/tasks"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/bor/accounts/abi"
	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/helper"
	stakingTypes "github.com/maticnetwork/heimdall/staking/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

const (
	defaultDelayDuration time.Duration = 15 * time.Second
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
	if err := sp.queueConnector.Server.RegisterTask("sendValidatorJoinToHeimdall", sp.sendValidatorJoinToHeimdall); err != nil {
		sp.Logger.Error("RegisterTasks | sendValidatorJoinToHeimdall", "error", err)
	}
	if err := sp.queueConnector.Server.RegisterTask("sendUnstakeInitToHeimdall", sp.sendUnstakeInitToHeimdall); err != nil {
		sp.Logger.Error("RegisterTasks | sendUnstakeInitToHeimdall", "error", err)
	}
	if err := sp.queueConnector.Server.RegisterTask("sendStakeUpdateToHeimdall", sp.sendStakeUpdateToHeimdall); err != nil {
		sp.Logger.Error("RegisterTasks | sendStakeUpdateToHeimdall", "error", err)
	}
	if err := sp.queueConnector.Server.RegisterTask("sendSignerChangeToHeimdall", sp.sendSignerChangeToHeimdall); err != nil {
		sp.Logger.Error("RegisterTasks | sendSignerChangeToHeimdall", "error", err)
	}
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
		if isOld, _ := sp.isOldTx(sp.cliCtx, vLog.TxHash.String(), uint64(vLog.Index), util.StakingEvent, event); isOld {
			sp.Logger.Info("Ignoring task to send validatorjoin to heimdall as already processed",
				"event", eventName,
				"validatorID", event.ValidatorId,
				"activationEpoch", event.ActivationEpoch,
				"nonce", event.Nonce,
				"amount", event.Amount,
				"totalAmount", event.Total,
				"SignerPubkey", hmTypes.NewPubKey(signerPubKey).String(),
				"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				"logIndex", uint64(vLog.Index),
				"blockNumber", vLog.BlockNumber,
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
			"nonce", event.Nonce,
			"amount", event.Amount,
			"totalAmount", event.Total,
			"SignerPubkey", hmTypes.NewPubKey(signerPubKey).String(),
			"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			"logIndex", uint64(vLog.Index),
			"blockNumber", vLog.BlockNumber,
		)

		// msg validator exit
		msg := stakingTypes.NewMsgValidatorJoin(
			hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
			event.ValidatorId.Uint64(),
			event.ActivationEpoch.Uint64(),
			sdk.NewIntFromBigInt(event.Amount),
			hmTypes.NewPubKey(signerPubKey),
			hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			uint64(vLog.Index),
			vLog.BlockNumber,
			event.Nonce.Uint64(),
		)

		// return broadcast to heimdall
		if err := sp.txBroadcaster.BroadcastToHeimdall(msg, event); err != nil {
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
		if isOld, _ := sp.isOldTx(sp.cliCtx, vLog.TxHash.String(), uint64(vLog.Index), util.StakingEvent, event); isOld {
			sp.Logger.Info("Ignoring task to send unstakeinit to heimdall as already processed",
				"event", eventName,
				"validator", event.User,
				"validatorID", event.ValidatorId,
				"nonce", event.Nonce,
				"deactivatonEpoch", event.DeactivationEpoch,
				"amount", event.Amount,
				"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				"logIndex", uint64(vLog.Index),
				"blockNumber", vLog.BlockNumber,
			)
			return nil
		}

		validNonce, nonceDelay, err := sp.checkValidNonce(event.ValidatorId.Uint64(), event.Nonce.Uint64())
		if err != nil {
			sp.Logger.Error("Error while validating nonce for the validator", "error", err)
			return err
		}

		if !validNonce {
			sp.Logger.Info("Ignoring task to send unstake-init to heimdall as nonce is out of order")
			return tasks.NewErrRetryTaskLater("Nonce out of order", defaultDelayDuration*time.Duration(nonceDelay))
		}

		sp.Logger.Info(
			"✅ Received task to send unstake-init to heimdall",
			"event", eventName,
			"validator", event.User,
			"validatorID", event.ValidatorId,
			"nonce", event.Nonce,
			"deactivatonEpoch", event.DeactivationEpoch,
			"amount", event.Amount,
			"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			"logIndex", uint64(vLog.Index),
			"blockNumber", vLog.BlockNumber,
		)

		// msg validator exit
		msg := stakingTypes.NewMsgValidatorExit(
			hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
			event.ValidatorId.Uint64(),
			event.DeactivationEpoch.Uint64(),
			hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			uint64(vLog.Index),
			vLog.BlockNumber,
			event.Nonce.Uint64(),
		)

		// return broadcast to heimdall
		if err := sp.txBroadcaster.BroadcastToHeimdall(msg, event); err != nil {
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
		if isOld, _ := sp.isOldTx(sp.cliCtx, vLog.TxHash.String(), uint64(vLog.Index), util.StakingEvent, event); isOld {
			sp.Logger.Info("Ignoring task to send unstakeinit to heimdall as already processed",
				"event", eventName,
				"validatorID", event.ValidatorId,
				"nonce", event.Nonce,
				"newAmount", event.NewAmount,
				"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				"logIndex", uint64(vLog.Index),
				"blockNumber", vLog.BlockNumber,
			)
			return nil
		}

		validNonce, nonceDelay, err := sp.checkValidNonce(event.ValidatorId.Uint64(), event.Nonce.Uint64())
		if err != nil {
			sp.Logger.Error("Error while validating nonce for the validator", "error", err)
			return err
		}

		if !validNonce {
			sp.Logger.Info("Ignoring task to send stake-update to heimdall as nonce is out of order")
			return tasks.NewErrRetryTaskLater("Nonce out of order", defaultDelayDuration*time.Duration(nonceDelay))
		}

		sp.Logger.Info(
			"✅ Received task to send stake-update to heimdall",
			"event", eventName,
			"validatorID", event.ValidatorId,
			"nonce", event.Nonce,
			"newAmount", event.NewAmount,
			"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			"logIndex", uint64(vLog.Index),
			"blockNumber", vLog.BlockNumber,
		)

		// msg validator exit
		msg := stakingTypes.NewMsgStakeUpdate(
			hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
			event.ValidatorId.Uint64(),
			sdk.NewIntFromBigInt(event.NewAmount),
			hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			uint64(vLog.Index),
			vLog.BlockNumber,
			event.Nonce.Uint64(),
		)

		// return broadcast to heimdall
		if err := sp.txBroadcaster.BroadcastToHeimdall(msg, event); err != nil {
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

		if isOld, _ := sp.isOldTx(sp.cliCtx, vLog.TxHash.String(), uint64(vLog.Index), util.StakingEvent, event); isOld {
			sp.Logger.Info("Ignoring task to send unstakeinit to heimdall as already processed",
				"event", eventName,
				"validatorID", event.ValidatorId,
				"nonce", event.Nonce,
				"NewSignerPubkey", hmTypes.NewPubKey(newSignerPubKey).String(),
				"oldSigner", event.OldSigner.Hex(),
				"newSigner", event.NewSigner.Hex(),
				"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				"logIndex", uint64(vLog.Index),
				"blockNumber", vLog.BlockNumber,
			)
			return nil
		}

		validNonce, nonceDelay, err := sp.checkValidNonce(event.ValidatorId.Uint64(), event.Nonce.Uint64())
		if err != nil {
			sp.Logger.Error("Error while validating nonce for the validator", "error", err)
			return err
		}

		if !validNonce {
			sp.Logger.Info("Ignoring task to send signer-change to heimdall as nonce is out of order")
			return tasks.NewErrRetryTaskLater("Nonce out of order", defaultDelayDuration*time.Duration(nonceDelay))
		}

		sp.Logger.Info(
			"✅ Received task to send signer-change to heimdall",
			"event", eventName,
			"validatorID", event.ValidatorId,
			"nonce", event.Nonce,
			"NewSignerPubkey", hmTypes.NewPubKey(newSignerPubKey).String(),
			"oldSigner", event.OldSigner.Hex(),
			"newSigner", event.NewSigner.Hex(),
			"txHash", hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			"logIndex", uint64(vLog.Index),
			"blockNumber", vLog.BlockNumber,
		)

		// signer change
		msg := stakingTypes.NewMsgSignerUpdate(
			hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
			event.ValidatorId.Uint64(),
			hmTypes.NewPubKey(newSignerPubKey),
			hmTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			uint64(vLog.Index),
			vLog.BlockNumber,
			event.Nonce.Uint64(),
		)

		// return broadcast to heimdall
		if err := sp.txBroadcaster.BroadcastToHeimdall(msg, event); err != nil {
			sp.Logger.Error("Error while broadcasting signerChainge to heimdall", "validatorId", event.ValidatorId.Uint64(), "error", err)
			return err
		}
	}
	return nil
}

func (sp *StakingProcessor) checkValidNonce(validatorId uint64, txnNonce uint64) (bool, uint64, error) {
	currentNonce, currentHeight, err := util.GetValidatorNonce(sp.cliCtx, validatorId)
	if err != nil {
		sp.Logger.Error("Failed to fetch validator nonce and height data from API", "validatorId", validatorId)
		return false, 0, err
	}

	if currentNonce+1 != txnNonce {
		diff := txnNonce - currentNonce
		if diff > 10 {
			diff = 10
		}
		sp.Logger.Error("Nonce for the given event not in order", "validatorId", validatorId, "currentNonce", currentNonce, "txnNonce", txnNonce, "delay", diff*uint64(defaultDelayDuration))
		return false, diff, nil
	}

	stakingTxnCount, err := queryTxCount(sp.cliCtx, validatorId, currentHeight)
	if err != nil {
		sp.Logger.Error("Failed to query stake txns by txquery for the given validator", "validatorId", validatorId)
		return false, 0, err
	}

	if stakingTxnCount != 0 {
		sp.Logger.Info("Recent staking txn count for the given validator is not zero", "validatorId", validatorId, "currentNonce", currentNonce, "txnNonce", txnNonce, "currentHeight", currentHeight)
		return false, 1, nil
	}

	return true, 0, nil
}

func queryTxCount(cliCtx cliContext.CLIContext, validatorId uint64, currentHeight int64) (int, error) {
	const (
		defaultPage  = 1
		defaultLimit = 30 // should be consistent with tendermint/tendermint/rpc/core/pipe.go:19
	)

	stakingTxnMsgMap := map[string]string{
		"validator-stake-update": "stake-update",
		"validator-join":         "validator-join",
		"signer-update":          "signer-update",
		"validator-exit":         "validator-exit",
	}

	for msg, action := range stakingTxnMsgMap {
		events := []string{
			fmt.Sprintf("%s.%s='%s'", sdk.EventTypeMessage, sdk.AttributeKeyAction, msg),
			fmt.Sprintf("%s.%s=%d", action, "validator-id", validatorId),
			fmt.Sprintf("%s.%s>%d", "tx", "height", currentHeight-3),
		}

		searchResult, err := helper.QueryTxsByEvents(cliCtx, events, defaultPage, defaultLimit)
		if err != nil {
			return 0, err
		}

		if searchResult.TotalCount != 0 {
			return searchResult.TotalCount, nil
		}
	}
	return 0, nil
}
