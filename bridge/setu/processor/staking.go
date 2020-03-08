package processor

import (
	"bytes"
	"encoding/json"

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
	sp.queueConnector.Server.RegisterTask("sendUnstakeInitToHeimdall", sp.sendUnstakeInitToHeimdall)
	sp.queueConnector.Server.RegisterTask("sendStakeUpdateToHeimdall", sp.sendStakeUpdateToHeimdall)
	sp.queueConnector.Server.RegisterTask("sendSignerChangeToHeimdall", sp.sendSignerChangeToHeimdall)
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
		sp.Logger.Info(
			"✅ Received task to send unstake-init to heimdall",
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
		sp.Logger.Info(
			"✅ Received task to send stake-update to heimdall",
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
		sp.Logger.Info(
			"✅ Received task to send signer-change to heimdall",
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
