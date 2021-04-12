package processor

import (
	"encoding/json"
	"fmt"

	"github.com/RichardKnop/machinery/v1/tasks"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/maticnetwork/bor/accounts/abi"
	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/helper"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"
	stakingTypes "github.com/maticnetwork/heimdall/x/staking/types"
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
		// convert PubKey to bytes
		pubkeyBytes, err := helper.ValidateAndCompressPubKey(signerPubKey)
		if err != nil {
			sp.Logger.Error("Error while uncompressed pubkey", "name", eventName, "error", err)
			return fmt.Errorf("Invalid uncompressed pubkey %s", err)
		}
		if isOld, _ := sp.isOldTx(sp.cliCtx, vLog.TxHash.String(), uint64(vLog.Index)); isOld {
			sp.Logger.Info("Ignoring task to send validatorjoin to heimdall as already processed",
				"event", eventName,
				"validatorID", event.ValidatorId,
				"activationEpoch", event.ActivationEpoch,
				"nonce", event.Nonce,
				"amount", event.Amount,
				"totalAmount", event.Total,
				"SignerPubkey", hmCommonTypes.NewPubKey(pubkeyBytes).String(),
				"txHash", hmCommonTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				"logIndex", uint64(vLog.Index),
				"blockNumber", vLog.BlockNumber,
			)
			return nil
		}

		// if account doesn't exists Retry with delay for topup to process first.
		if _, err := util.GetAccount(sp.cliCtx, hmCommonTypes.HeimdallAddress(event.Signer)); err != nil {
			sp.Logger.Info(
				"Heimdall Account doesn't exist. Retrying validator-join after 10 seconds",
				"event", eventName,
				"signer", event.Signer,
			)
			return tasks.NewErrRetryTaskLater("account doesn't exist", util.ValidatorJoinRetryDelay)
		}

		sp.Logger.Info(
			"✅ Received task to send validatorjoin to heimdall",
			"event", eventName,
			"validatorID", event.ValidatorId,
			"activationEpoch", event.ActivationEpoch,
			"nonce", event.Nonce,
			"amount", event.Amount,
			"totalAmount", event.Total,
			"SignerPubkey", hmCommonTypes.NewPubKey(pubkeyBytes).String(),
			"txHash", hmCommonTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			"logIndex", uint64(vLog.Index),
			"blockNumber", vLog.BlockNumber,
		)

		// msg validator exit
		accAddr, _ := sdk.AccAddressFromHex(helper.GetAddressStr())
		msg := stakingTypes.NewMsgValidatorJoin(
			accAddr,
			event.ValidatorId.Uint64(),
			event.ActivationEpoch.Uint64(),
			sdk.NewIntFromBigInt(event.Amount),
			hmCommonTypes.NewPubKey(pubkeyBytes),
			hmCommonTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			uint64(vLog.Index),
			vLog.BlockNumber,
			event.Nonce.Uint64(),
		)

		// return broadcast to heimdall
		if err := sp.txBroadcaster.BroadcastToHeimdall(&msg); err != nil {
			sp.Logger.Error("Error while broadcasting validatorJoin to heimdall", "validatorId", event.ValidatorId.Uint64(), "error", err)
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
				"nonce", event.Nonce,
				"deactivatonEpoch", event.DeactivationEpoch,
				"amount", event.Amount,
				"txHash", hmCommonTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				"logIndex", uint64(vLog.Index),
				"blockNumber", vLog.BlockNumber,
			)
			return nil
		}

		sp.Logger.Info(
			"✅ Received task to send unstake-init to heimdall",
			"event", eventName,
			"validator", event.User,
			"validatorID", event.ValidatorId,
			"nonce", event.Nonce,
			"deactivatonEpoch", event.DeactivationEpoch,
			"amount", event.Amount,
			"txHash", hmCommonTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			"logIndex", uint64(vLog.Index),
			"blockNumber", vLog.BlockNumber,
		)

		// msg validator exit
		accAddr, _ := sdk.AccAddressFromHex(helper.GetAddressStr())
		msg := stakingTypes.NewMsgValidatorExit(
			accAddr,
			event.ValidatorId.Uint64(),
			event.DeactivationEpoch.Uint64(),
			hmCommonTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			uint64(vLog.Index),
			vLog.BlockNumber,
			event.Nonce.Uint64(),
		)

		// return broadcast to heimdall
		if err := sp.txBroadcaster.BroadcastToHeimdall(&msg); err != nil {
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
			sp.Logger.Info("Ignoring task to send stake-update to heimdall as already processed",
				"event", eventName,
				"validatorID", event.ValidatorId,
				"nonce", event.Nonce,
				"newAmount", event.NewAmount,
				"txHash", hmCommonTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				"logIndex", uint64(vLog.Index),
				"blockNumber", vLog.BlockNumber,
			)
			return nil
		}
		sp.Logger.Info(
			"✅ Received task to send stake-update to heimdall",
			"event", eventName,
			"validatorID", event.ValidatorId,
			"nonce", event.Nonce,
			"newAmount", event.NewAmount,
			"txHash", hmCommonTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			"logIndex", uint64(vLog.Index),
			"blockNumber", vLog.BlockNumber,
		)

		// msg validator exit
		accAddr, _ := sdk.AccAddressFromHex(helper.GetAddressStr())
		msg := stakingTypes.NewMsgStakeUpdate(
			accAddr,
			event.ValidatorId.Uint64(),
			sdk.NewIntFromBigInt(event.NewAmount),
			hmCommonTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			uint64(vLog.Index),
			vLog.BlockNumber,
			event.Nonce.Uint64(),
		)

		// return broadcast to heimdall
		if err := sp.txBroadcaster.BroadcastToHeimdall(&msg); err != nil {
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
		// convert PubKey to bytes
		newSignerPubKey, err := helper.ValidateAndCompressPubKey(event.SignerPubkey)
		if err != nil {
			return fmt.Errorf("Invalid uncompressed pubkey %s", err)
		}

		if isOld, _ := sp.isOldTx(sp.cliCtx, vLog.TxHash.String(), uint64(vLog.Index)); isOld {
			sp.Logger.Info("Ignoring task to send unstakeinit to heimdall as already processed",
				"event", eventName,
				"validatorID", event.ValidatorId,
				"nonce", event.Nonce,
				"NewSignerPubkey", hmCommonTypes.NewPubKey(newSignerPubKey).String(),
				"oldSigner", event.OldSigner.Hex(),
				"newSigner", event.NewSigner.Hex(),
				"txHash", hmCommonTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
				"logIndex", uint64(vLog.Index),
				"blockNumber", vLog.BlockNumber,
			)
			return nil
		}
		sp.Logger.Info(
			"✅ Received task to send signer-change to heimdall",
			"event", eventName,
			"validatorID", event.ValidatorId,
			"nonce", event.Nonce,
			"NewSignerPubkey", hmCommonTypes.NewPubKey(newSignerPubKey).String(),
			"oldSigner", event.OldSigner.Hex(),
			"newSigner", event.NewSigner.Hex(),
			"txHash", hmCommonTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			"logIndex", uint64(vLog.Index),
			"blockNumber", vLog.BlockNumber,
		)

		// signer change
		accAddr, _ := sdk.AccAddressFromHex(helper.GetAddressStr())
		msg := stakingTypes.NewMsgSignerUpdate(
			accAddr,
			event.ValidatorId.Uint64(),
			hmCommonTypes.NewPubKey(newSignerPubKey),
			hmCommonTypes.BytesToHeimdallHash(vLog.TxHash.Bytes()),
			uint64(vLog.Index),
			vLog.BlockNumber,
			event.Nonce.Uint64(),
		)

		// return broadcast to heimdall
		if err := sp.txBroadcaster.BroadcastToHeimdall(&msg); err != nil {
			sp.Logger.Error("Error while broadcasting signerChainge to heimdall", "validatorId", event.ValidatorId.Uint64(), "error", err)
			return err
		}
	}
	return nil
}

// isOldTx  checks if tx is already processed or not
func (sp *StakingProcessor) isOldTx(cliCtx client.Context, txHash string, logIndex uint64) (bool, error) {
	queryParam := map[string]interface{}{
		"tx_hash":   txHash,
		"log_index": logIndex,
	}

	endpoint := helper.GetHeimdallServerEndpoint(util.StakingTxStatusURL)
	url, err := util.CreateURLWithQuery(endpoint, queryParam)
	if err != nil {
		sp.Logger.Error("Error in creating url", "endpoint", endpoint, "error", err)
		return false, err
	}

	res, err := helper.FetchFromAPI(url)
	if err != nil {
		sp.Logger.Error("Error fetching tx status", "url", url, "error", err)
		return false, err
	}

	var status bool
	if err := json.Unmarshal(res, &status); err != nil {
		sp.Logger.Error("Error unmarshalling tx status received from Heimdall Server", "error", err)
		return false, err
	}

	return status, nil
}
