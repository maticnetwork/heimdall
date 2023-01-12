package processor

import (
	"context"
	"math/big"
	"time"

	"github.com/maticnetwork/heimdall/bridge/setu/util"
	chainmanagerTypes "github.com/maticnetwork/heimdall/chainmanager/types"
	milestoneTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/pborman/uuid"

	hmTypes "github.com/maticnetwork/heimdall/types"
)

// milestoneProcessor - process milestone related events
type MilestoneProcessor struct {
	BaseProcessor

	// header listener subscription
	cancelMilestoneService context.CancelFunc
}

// MilestoneContext represents milestone context
type MilestoneContext struct {
	ChainmanagerParams *chainmanagerTypes.Params
}

// Start starts new block subscription
func (mp *MilestoneProcessor) Start() error {
	mp.Logger.Info("Starting")

	// create cancellable context
	milestoneCtx, cancelMilestoneService := context.WithCancel(context.Background())

	mp.cancelMilestoneService = cancelMilestoneService

	// start polling for milestone
	mp.Logger.Info("Start polling for milestone", "milestoneLength", helper.MilestoneLength, "pollInterval", helper.GetConfig().MilestonePollInterval)

	go mp.startPolling(milestoneCtx, helper.MilestoneLength, helper.GetConfig().MilestonePollInterval)
	go mp.startPollingMilestoneTimeout(milestoneCtx, helper.GetConfig().MilestonePollInterval)

	return nil
}

// RegisterTasks - nil
func (mp *MilestoneProcessor) RegisterTasks() {

}

// startPolling - polls heimdall and checks if new span needs to be proposed
func (mp *MilestoneProcessor) startPolling(ctx context.Context, milestoneLength uint64, interval time.Duration) {
	ticker := time.NewTicker(interval)
	// stop ticker when everything done
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := mp.checkAndPropose(milestoneLength)
			if err != nil {
				mp.Logger.Error("Error in proposing the milestone", "error", err)
			}
		case <-ctx.Done():
			mp.Logger.Info("Polling stopped")
			return
		}
	}
}

// sendMilestoneToHeimdall - handles headerblock from maticchain
// 1. check if i am the proposer for next milestone
// 2. check if milestone has to be proposed
// 3. if so, propose milestone to heimdall.
func (mp *MilestoneProcessor) checkAndPropose(milestoneLength uint64) (err error) {
	//Milestone proposing mechanism will work only after specific block height
	if util.GetBlockHeight(mp.cliCtx) < helper.GetMilestoneHardForkHeight() {
		mp.Logger.Debug("Block height Less than fork height", "current block height", util.GetBlockHeight(mp.cliCtx), "milestone hard fork height", helper.GetMilestoneHardForkHeight())
		return nil
	}

	// fetch milestone context
	milestoneContext, err := mp.getMilestoneContext()
	if err != nil {
		return err
	}

	isProposer, err := util.IsMilestoneProposer(mp.cliCtx)
	if err != nil {
		mp.Logger.Error("Error checking isProposer in HeaderBlock handler", "error", err)
		return err
	}

	if isProposer {
		result, err := util.GetMilestoneCount(mp.cliCtx)
		if err != nil || result == nil {
			return err
		}

		var start = helper.GetMilestoneBorBlockHeight()

		if result.Count != 0 {
			// fetch latest milestone
			latestMilestone, err := util.GetLatestMilestone(mp.cliCtx)
			if err != nil || latestMilestone == nil {
				return err
			}

			start = latestMilestone.EndBlock + 1
		}

		end := start + milestoneLength - 1

		if err := mp.createAndSendMilestoneToHeimdall(milestoneContext, start, end, milestoneLength); err != nil {
			mp.Logger.Error("Error sending milestone to heimdall", "error", err)
			return err
		}
	} else {
		mp.Logger.Info("I am not the current milestone proposer")
	}

	return nil
}

// sendMilestoneToHeimdall - creates milestone msg and broadcasts to heimdall
func (mp *MilestoneProcessor) createAndSendMilestoneToHeimdall(milestoneContext *MilestoneContext, start uint64, end uint64, milestoneLength uint64) error {
	mp.Logger.Debug("Initiating milestone to Heimdall", "start", start, "end", end, "milestoneLength", milestoneLength)

	// Get root hash
	endBlock, err := mp.contractConnector.GetMaticChainBlockByNumber(big.NewInt(0).SetUint64(end))
	if err != nil {
		return err
	}

	mainEndBlock, err := mp.contractConnector.GetMainChainBlock(big.NewInt(0).SetUint64(end))

	blockHash := endBlock.Header().Hash()
	mainHash := mainEndBlock.Hash()

	num := endBlock.Number().Uint64()

	milestoneId := uuid.NewRandom().String() + "-" + hmTypes.BytesToHeimdallAddress(helper.GetAddress()).String()

	mp.Logger.Info("Root hash calculated", "root", hmTypes.BytesToHeimdallHash(blockHash[:]))

	mp.Logger.Info("✅ Creating and broadcasting new milestone",
		"start", start,
		"end", end,
		"endNew", num,
		"root", blockHash,
		"main", mainHash,
		"milestoneId", milestoneId,
		"milestoneLength", milestoneLength,
	)

	chainParams := milestoneContext.ChainmanagerParams.ChainParams

	// create and send milestone message
	msg := milestoneTypes.NewMsgMilestoneBlock(
		hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
		start,
		end,
		hmTypes.BytesToHeimdallHash(blockHash[:]),
		chainParams.BorChainID,
		milestoneId,
	)

	// return broadcast to heimdall
	if err := mp.txBroadcaster.BroadcastToHeimdall(msg, nil); err != nil {
		mp.Logger.Error("Error while broadcasting milestone to heimdall", "error", err)
		return err
	}

	return nil
}

// startPolling - polls heimdall and checks if new milestoneTimeout needs to be proposed
func (mp *MilestoneProcessor) startPollingMilestoneTimeout(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	// stop ticker when everything done
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := mp.checkAndProposeMilestoneTimeout()
			if err != nil {
				mp.Logger.Error("Error in proposing the MilestoneTimeout msg", "error", err)
			}
		case <-ctx.Done():
			mp.Logger.Info("Polling stopped")
			ticker.Stop()

			return
		}
	}
}

// sendMilestoneToHeimdall - handles headerblock from maticchain
// 1. check if i am the proposer for next milestone
// 2. check if milestone has to be proposed
// 3. if so, propose milestone to heimdall.
func (mp *MilestoneProcessor) checkAndProposeMilestoneTimeout() (err error) {
	//Milestone proposing mechanism will work only after specific block height
	if util.GetBlockHeight(mp.cliCtx) < helper.GetMilestoneHardForkHeight() {
		mp.Logger.Debug("Block height Less than fork height", "current block height", util.GetBlockHeight(mp.cliCtx), "milestone hard fork height", helper.GetMilestoneHardForkHeight())
		return nil
	}

	isMilestoneTimeoutRequired, err := mp.checkIfMilestoneTimeoutIsRequired()
	if err != nil {
		mp.Logger.Debug("Error checking sMilestoneTimeoutRequired while proposing Milestone Timeout ", "error", err)
		return
	}

	if isMilestoneTimeoutRequired {
		var isProposer bool

		if isProposer, err = util.IsInMilestoneProposerList(mp.cliCtx, 10); err != nil {
			mp.Logger.Error("Error checking IsInMilestoneProposerList while proposing Milestone Timeout ", "error", err)
			return
		}

		// if i am the proposer and NoAck is required, then propose No-Ack
		if isProposer {
			// send Checkpoint No-Ack to heimdall
			if err = mp.createAndSendMilestoneTimeoutToHeimdall(); err != nil {
				mp.Logger.Error("Error proposing Milestone-Timeout ", "error", err)
				return
			}
		}
	}

	return nil
}

// sendMilestoneTimoutToHeimdall - creates milestone-timeout msg and broadcasts to heimdall
func (mp *MilestoneProcessor) createAndSendMilestoneTimeoutToHeimdall() error {
	mp.Logger.Debug("Initiating milestone timeout to Heimdall")

	mp.Logger.Info("✅ Creating and broadcasting milestone-timeout")

	// create and send milestone message
	msg := milestoneTypes.NewMsgMilestoneTimeout(
		hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
	)

	// return broadcast to heimdall
	if err := mp.txBroadcaster.BroadcastToHeimdall(msg, nil); err != nil {
		mp.Logger.Error("Error while broadcasting milestone timeout to heimdall", "error", err)
		return err
	}

	return nil
}

func (mp *MilestoneProcessor) checkIfMilestoneTimeoutIsRequired() (bool, error) {
	latestMilestone, err := util.GetLatestMilestone(mp.cliCtx)
	if err != nil || latestMilestone == nil {
		return false, err
	}

	lastMilestoneEndBlock := latestMilestone.EndBlock
	currentChildBlockNumber, _ := mp.getCurrentChildBlock()

	if err != nil {
		return false, err
	}

	if (currentChildBlockNumber - lastMilestoneEndBlock) > helper.MilestoneBufferLength {
		return true, nil
	}

	return false, nil
}

// getCurrentChildBlock gets the current child block
func (mp *MilestoneProcessor) getCurrentChildBlock() (uint64, error) {
	childBlock, err := mp.contractConnector.GetMaticChainBlock(nil)
	if err != nil {
		return 0, err
	}

	return childBlock.Number.Uint64(), nil
}

func (mp *MilestoneProcessor) getMilestoneContext() (*MilestoneContext, error) {
	chainmanagerParams, err := util.GetChainmanagerParams(mp.cliCtx)
	if err != nil {
		mp.Logger.Error("Error while fetching chain manager params", "error", err)
		return nil, err
	}

	return &MilestoneContext{
		ChainmanagerParams: chainmanagerParams,
	}, nil
}

// Stop stops all necessary go routines
func (mp *MilestoneProcessor) Stop() {
	// cancel milestone polling
	mp.cancelMilestoneService()
}
