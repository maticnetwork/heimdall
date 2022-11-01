package processor

import (
	"context"
	"time"

	"github.com/maticnetwork/heimdall/bridge/setu/util"
	chainmanagerTypes "github.com/maticnetwork/heimdall/chainmanager/types"
	"github.com/maticnetwork/heimdall/helper"
	milestoneTypes "github.com/maticnetwork/heimdall/milestone/types"
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
	MilestoneParams    *milestoneTypes.Params
}

// Start starts new block subscription
func (mp *MilestoneProcessor) Start() error {
	mp.Logger.Info("Starting")

	// create cancellable context
	milestoneCtx, cancelMilestoneService := context.WithCancel(context.Background())

	mp.cancelMilestoneService = cancelMilestoneService

	// start polling for span
	mp.Logger.Info("Start polling for milestone", "pollInterval", helper.GetConfig().MilestonePollInterval)

	go mp.startPolling(milestoneCtx, helper.GetConfig().MilestonePollInterval)

	return nil
}

// RegisterTasks - nil
func (mp *MilestoneProcessor) RegisterTasks() {

}

// sendMilestoneToHeimdall - handles headerblock from maticchain
// 1. check if i am the proposer for next milestone
// 2. check if milestone has to be proposed
// 3. if so, propose milestone to heimdall.
func (mp *MilestoneProcessor) checkAndPropose() (err error) {

	// fetch milestone context
	milestoneContext, err := mp.getMilestoneContext()
	if err != nil {
		return err
	}

	milestoneParams := milestoneContext.MilestoneParams

	isProposer, err := util.IsProposer(mp.cliCtx)
	if err != nil {
		mp.Logger.Error("Error checking isProposer in HeaderBlock handler", "error", err)
		return err
	}

	if isProposer {

		//fetch the block head number from bor chain
		currentBlockNumber, err := mp.getCurrentChildBlock()
		if err != nil {
			return err
		}

		result, err := util.GetMilestoneCount(mp.cliCtx)
		if err != nil || result == nil {
			return err
		}

		var start uint64

		if result.Count != 0 {
			// fetch latest milestone
			latestMilestone, err := util.GetLatestMilestone(mp.cliCtx)
			if err != nil || latestMilestone == nil {
				return err
			}

			// current Block should be greater than or equal to latest milestone end block + Sprint length
			if latestMilestone.EndBlock+milestoneParams.SprintLength > currentBlockNumber {
				mp.Logger.Debug("Current child block is less than latest milestone end block + sprint length", "currentChildBlock", currentBlockNumber, "latestMilestoneEndBlock", latestMilestone.EndBlock)
				return nil
			}

			mp.Logger.Debug("Current child block should be greater than or equal to latest milestone's endblock + sprintLength", "currentChildBlock", currentBlockNumber, "latestMilestoneEndBlock", latestMilestone.EndBlock)

			start = latestMilestone.EndBlock + 1

		}

		end := start + milestoneParams.SprintLength - 1

		if err := mp.createAndSendMilestoneToHeimdall(milestoneContext, start, end); err != nil {
			mp.Logger.Error("Error sending milestone to heimdall", "error", err)
			return err
		}
	} else {
		mp.Logger.Info("I am not the current milestone proposer")
		return
	}

	return nil
}

// sendMilestoneToHeimdall - creates milestone msg and broadcasts to heimdall
func (mp *MilestoneProcessor) createAndSendMilestoneToHeimdall(milestoneContext *MilestoneContext, start uint64, end uint64) error {
	mp.Logger.Debug("Initiating milestone to Heimdall", "start", start, "end", end)

	// get milestone params
	milestoneParams := milestoneContext.MilestoneParams

	// Get root hash
	root, err := mp.contractConnector.GetRootHash(start, end, milestoneParams.SprintLength)
	if err != nil {
		return err
	}

	milestoneId := uuid.NewRandom().String() + "-" + hmTypes.BytesToHeimdallAddress(helper.GetAddress()).String()

	mp.Logger.Info("Root hash calculated", "rootHash", hmTypes.BytesToHeimdallHash(root))

	mp.Logger.Info("✅ Creating and broadcasting new milestone",
		"start", start,
		"end", end,
		"root", hmTypes.BytesToHeimdallHash(root),
		"milestoneId",
	)

	chainParams := milestoneContext.ChainmanagerParams.ChainParams

	// create and send milestone message
	msg := milestoneTypes.NewMsgMilestoneBlock(
		hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
		start,
		end,
		hmTypes.BytesToHeimdallHash(root),
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

// startPolling - polls heimdall and checks if new span needs to be proposed
func (mp *MilestoneProcessor) startPolling(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(1 * time.Second)
	// stop ticker when everything done
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mp.checkAndPropose()
		case <-ctx.Done():
			mp.Logger.Info("Polling stopped")
			ticker.Stop()

			return
		}
	}
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

	milestoneParams, err := util.GetMilestoneParams(mp.cliCtx)
	if err != nil {
		mp.Logger.Error("Error while fetching milestone params", "error", err)
		return nil, err
	}

	return &MilestoneContext{
		ChainmanagerParams: chainmanagerParams,
		MilestoneParams:    milestoneParams,
	}, nil
}

// Stop stops all necessary go routines
func (mp *MilestoneProcessor) Stop() {
	// cancel milestone polling
	mp.cancelMilestoneService()
}
