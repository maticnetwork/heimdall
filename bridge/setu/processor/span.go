package processor

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	jsoniter "github.com/json-iterator/go"

	"github.com/ethereum/go-ethereum/common"

	borTypes "github.com/maticnetwork/heimdall/bor/types"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// SpanProcessor - process span related events
type SpanProcessor struct {
	BaseProcessor

	// header listener subscription
	cancelSpanService context.CancelFunc
}

// Start starts new block subscription
func (sp *SpanProcessor) Start() error {
	sp.Logger.Info("Starting")

	// create cancellable context
	spanCtx, cancelSpanService := context.WithCancel(context.Background())

	sp.cancelSpanService = cancelSpanService

	// start polling for span
	sp.Logger.Info("Start polling for span", "pollInterval", helper.GetConfig().SpanPollInterval)
	go sp.startPolling(spanCtx, helper.GetConfig().SpanPollInterval)

	return nil
}

// RegisterTasks - nil
func (sp *SpanProcessor) RegisterTasks() {

}

// startPolling - polls heimdall and checks if new span needs to be proposed
func (sp *SpanProcessor) startPolling(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	// stop ticker when everything done
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// nolint: contextcheck
			sp.checkAndPropose()
		case <-ctx.Done():
			sp.Logger.Info("Polling stopped")
			ticker.Stop()

			return
		}
	}
}

// checkAndPropose - will check if current user is span proposer and proposes the span
func (sp *SpanProcessor) checkAndPropose() {
	lastSpan, err := sp.getLastSpan()
	if err != nil {
		sp.Logger.Error("Unable to fetch last span", "error", err)
		return
	}

	if lastSpan == nil {
		return
	}

	lastBlock, e := sp.contractConnector.GetMaticChainBlock(nil)
	if e != nil {
		sp.Logger.Error("Error fetching current child block", "error", e)
		return
	}

	latestBorBlockNumber := lastBlock.Number.Uint64()

	if latestBorBlockNumber < lastSpan.StartBlock {
		sp.Logger.Debug("Current bor block is less than last span start block, skipping proposing span",
			"lastBlock", latestBorBlockNumber,
			"lastSpanStartBlock", lastSpan.StartBlock,
		)
		return
	}

	var latestMilestoneEndBlock uint64
	latestMilestone, err := util.GetLatestMilestone(sp.cliCtx)
	if err == nil {
		latestMilestoneEndBlock = latestMilestone.EndBlock
	} else {
		sp.Logger.Error("Error fetching latest milestone", "error", err)
	}

	// Max of latest milestone end block and latest bor block number
	// Handle cases where bor is syncing and latest milestone end block is greater than latest bor block number
	maxBlockNumber := max(latestMilestoneEndBlock, latestBorBlockNumber)

	sp.Logger.Debug("Found last span",
		"lastSpan", lastSpan.ID,
		"startBlock", lastSpan.StartBlock,
		"endBlock", lastSpan.EndBlock,
	)

	nextSpanMsg, err := sp.fetchNextSpanDetails(lastSpan.ID+1, lastSpan.EndBlock+1)
	if err != nil {
		sp.Logger.Error("Unable to fetch next span details", "error", err, "lastSpanId", lastSpan.ID)
		return
	}

	// check if current user is among next span producers
	if !sp.isSpanProposer(nextSpanMsg.SelectedProducers) {
		sp.Logger.Debug("Current user is not among next span producers, skipping proposing span", "nextSpanId", nextSpanMsg.ID)
		return
	}

	if maxBlockNumber > lastSpan.EndBlock {
		if latestMilestoneEndBlock > lastSpan.EndBlock {
			sp.Logger.Debug("Bor self commmitted spans, backfill heimdall to fill missing spans",
				"lastBlock", latestBorBlockNumber,
				"lastFinalizedBlock", latestMilestoneEndBlock,
				"lastSpanEndBlock", lastSpan.EndBlock,
			)

			if err := sp.backfillSpans(latestMilestoneEndBlock, lastSpan); err != nil {
				sp.Logger.Error("Error in backfillSpans", "error", err)
			}

		} else {
			sp.Logger.Debug("Will not backfill heimdall spans, as latest milestone end block is less than last span end block",
				"lastBlock", latestBorBlockNumber,
				"lastFinalizedBlock", latestMilestoneEndBlock,
				"lastSpanEndBlock", lastSpan.EndBlock,
			)
		}
	} else {
		if borTypes.IsBlockCloseToSpanEnd(maxBlockNumber, lastSpan.EndBlock) {
			sp.Logger.Debug("Current bor block is close to last span end block, skipping proposing span",
				"lastBlock", latestBorBlockNumber,
				"lastFinalizedBlock", latestMilestoneEndBlock,
				"lastSpanEndBlock", lastSpan.EndBlock,
			)
			return
		}

		go sp.propose(lastSpan, nextSpanMsg)
	}
}

func (sp *SpanProcessor) backfillSpans(latestFinalizedBorBlockNumber uint64, lastHeimdallSpan *types.Span) error {
	// Get what span bor used when heimdall was down and it tried to fetch the next span from heimdall
	// We need to know if it managed to fetch the latest span or it used the previous one
	// We will take it from the next start block after the last heimdall span end block
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	borLastUsedSpanID, err := sp.contractConnector.GetStartBlockHeimdallSpanID(ctx, lastHeimdallSpan.EndBlock+1)
	if err != nil {
		return fmt.Errorf("error while fetching last used span id for bor: %w", err)
	}

	if borLastUsedSpanID == 0 {
		return fmt.Errorf("bor last used span id is 0, cannot backfill spans")
	}

	sp.Logger.Error("Found bor last used span id", "borLastUsedSpanId", borLastUsedSpanID)

	borLastUsedSpan, err := sp.getSpanById(sp.cliCtx, borLastUsedSpanID)
	if err != nil {
		return fmt.Errorf("error while fetching last used span for bor: %w", err)
	}

	if borLastUsedSpan == nil {
		return fmt.Errorf("bor last used span is nil, cannot backfill spans")
	}

	params, err := util.GetChainmanagerParams(sp.cliCtx)
	if err != nil {
		return fmt.Errorf("error while fetching chainmanager params: %w", err)
	}

	addrString := hmTypes.BytesToHeimdallAddress(helper.GetAddress())

	borSpanId, err := borTypes.CalcCurrentBorSpanId(latestFinalizedBorBlockNumber, borLastUsedSpan)
	if err != nil {
		return fmt.Errorf("error while calculating bor span id: %w", err)
	}

	msg := borTypes.MsgBackfillSpans{
		Proposer:        addrString,
		ChainID:         params.ChainParams.BorChainID,
		LatestSpanID:    borLastUsedSpan.ID,
		LatestBorSpanID: borSpanId,
	}

	txRes, err := sp.txBroadcaster.BroadcastToHeimdall(&msg, nil) //nolint:contextcheck
	if err != nil {
		return fmt.Errorf("error while broadcasting backfill spans to heimdall: %w", err)
	}

	if txRes.Code != uint32(sdk.CodeOK) {
		return fmt.Errorf("backfill spans tx failed on heimdall, txHash: %s, code: %d", txRes.TxHash, txRes.Code)
	}

	sp.Logger.Info("Backfill spans tx successfully broadcasted to heimdall", "txHash", txRes.TxHash, "borLastUsedSpanId", borLastUsedSpan.ID, "borSpanId", borSpanId)

	return nil
}

// propose producers for next span if needed
func (sp *SpanProcessor) propose(lastSpan *types.Span, nextSpanMsg *types.Span) {
	// call with last span on record + new span duration and see if it has been proposed
	currentBlock, err := sp.getCurrentChildBlock()
	if err != nil {
		sp.Logger.Error("Unable to fetch current block", "error", err)
		return
	}

	if lastSpan.StartBlock <= currentBlock && currentBlock <= lastSpan.EndBlock {
		// log new span
		sp.Logger.Info("✅ Proposing new span", "spanId", nextSpanMsg.ID, "startBlock", nextSpanMsg.StartBlock, "endBlock", nextSpanMsg.EndBlock)

		seed, seedAuthor, err := sp.fetchNextSpanSeed(nextSpanMsg.ID)
		if err != nil {
			sp.Logger.Info("Error while fetching next span seed from HeimdallServer", "err", err)
			return
		}

		nodeStatus, err := helper.GetNodeStatus(sp.cliCtx)
		if err != nil {
			sp.Logger.Error("Error while fetching heimdall node status", "error", err)
			return
		}

		var txRes sdk.TxResponse

		if nodeStatus.SyncInfo.LatestBlockHeight < helper.GetDanelawHeight() {
			// broadcast to heimdall
			msg := borTypes.MsgProposeSpan{
				ID:         nextSpanMsg.ID,
				Proposer:   types.BytesToHeimdallAddress(helper.GetAddress()),
				StartBlock: nextSpanMsg.StartBlock,
				EndBlock:   nextSpanMsg.EndBlock,
				ChainID:    nextSpanMsg.ChainID,
				Seed:       seed,
			}

			// return broadcast to heimdall
			txRes, err = sp.txBroadcaster.BroadcastToHeimdall(msg, nil)
			if err != nil {
				sp.Logger.Error("Error while broadcasting span to heimdall", "spanId", nextSpanMsg.ID, "startBlock", nextSpanMsg.StartBlock, "endBlock", nextSpanMsg.EndBlock, "error", err)
				return
			}
		} else {
			msg := borTypes.MsgProposeSpanV2{
				ID:         nextSpanMsg.ID,
				Proposer:   types.BytesToHeimdallAddress(helper.GetAddress()),
				StartBlock: nextSpanMsg.StartBlock,
				EndBlock:   nextSpanMsg.EndBlock,
				ChainID:    nextSpanMsg.ChainID,
				Seed:       seed,
				SeedAuthor: seedAuthor,
			}

			txRes, err = sp.txBroadcaster.BroadcastToHeimdall(msg, nil)
			if err != nil {
				sp.Logger.Error("Error while broadcasting span to heimdall", "spanId", nextSpanMsg.ID, "startBlock", nextSpanMsg.StartBlock, "endBlock", nextSpanMsg.EndBlock, "error", err)
				return
			}
		}

		if txRes.Code != uint32(sdk.CodeOK) {
			sp.Logger.Error("span tx failed on heimdall", "txHash", txRes.TxHash, "code", txRes.Code)
			return

		}
	}
}

// checks span status
func (sp *SpanProcessor) getLastSpan() (*types.Span, error) {
	// fetch latest start block from heimdall via rest query
	result, err := helper.FetchFromAPI(sp.cliCtx, helper.GetHeimdallServerEndpoint(util.LatestSpanURL))
	if err != nil {
		sp.Logger.Error("Error while fetching latest span")
		return nil, err
	}

	var lastSpan types.Span
	if err = jsoniter.ConfigFastest.Unmarshal(result.Result, &lastSpan); err != nil {
		sp.Logger.Error("Error unmarshalling span", "error", err)
		return nil, err
	}

	return &lastSpan, nil
}

func (sp *SpanProcessor) getSpanById(cliCtx cliContext.CLIContext, id uint64) (*types.Span, error) {
	sp.Logger.Error("getSpanById called", "spanId", id, "url", fmt.Sprintf(helper.GetHeimdallServerEndpoint(util.SpanByIdURL), strconv.FormatUint(id, 10)))
	// fetch latest span from heimdall using the rest query
	result, err := helper.FetchFromAPI(cliCtx, fmt.Sprintf(helper.GetHeimdallServerEndpoint(util.SpanByIdURL), strconv.FormatUint(id, 10)))
	if err != nil {
		sp.Logger.Error("Error while fetching span by id", "spanId", id, "error", err)
		return nil, err
	}

	var lastSpan types.Span
	if err = jsoniter.ConfigFastest.Unmarshal(result.Result, &lastSpan); err != nil {
		sp.Logger.Error("Error unmarshalling span", "error", err)
		return nil, err
	}

	return &lastSpan, nil
}

// getCurrentChildBlock gets the current child block
func (sp *SpanProcessor) getCurrentChildBlock() (uint64, error) {
	childBlock, err := sp.contractConnector.GetMaticChainBlock(nil)
	if err != nil {
		return 0, err
	}

	return childBlock.Number.Uint64(), nil
}

// isSpanProposer checks if current user is span proposer
func (sp *SpanProcessor) isSpanProposer(nextSpanProducers []types.Validator) bool {
	// anyone among next span producers can become next span proposer
	for _, val := range nextSpanProducers {
		if bytes.Equal(val.Signer.Bytes(), helper.GetAddress()) {
			return true
		}
	}

	return false
}

// fetch next span details from heimdall.
func (sp *SpanProcessor) fetchNextSpanDetails(id uint64, start uint64) (*types.Span, error) {
	req, err := http.NewRequest("GET", helper.GetHeimdallServerEndpoint(util.NextSpanInfoURL), nil)
	if err != nil {
		sp.Logger.Error("Error creating a new request", "error", err)
		return nil, err
	}

	configParams, err := util.GetChainmanagerParams(sp.cliCtx)
	if err != nil {
		sp.Logger.Error("Error while fetching chainmanager params", "error", err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("span_id", strconv.FormatUint(id, 10))
	q.Add("start_block", strconv.FormatUint(start, 10))
	q.Add("chain_id", configParams.ChainParams.BorChainID)
	q.Add("proposer", helper.GetFromAddress(sp.cliCtx).String())
	req.URL.RawQuery = q.Encode()

	// fetch next span details
	result, err := helper.FetchFromAPI(sp.cliCtx, req.URL.String())
	if err != nil {
		sp.Logger.Error("Error fetching proposers", "error", err)
		return nil, err
	}

	var msg types.Span
	if err = jsoniter.ConfigFastest.Unmarshal(result.Result, &msg); err != nil {
		sp.Logger.Error("Error unmarshalling propose tx msg ", "error", err)
		return nil, err
	}

	sp.Logger.Debug("◽ Generated proposer span msg", "msg", msg.String())

	return &msg, nil
}

// fetchNextSpanSeed - fetches seed for next span
func (sp *SpanProcessor) fetchNextSpanSeed(id uint64) (common.Hash, common.Address, error) {
	sp.Logger.Info("Sending Rest call to Get Seed for next span")

	response, err := helper.FetchFromAPI(sp.cliCtx, helper.GetHeimdallServerEndpoint(fmt.Sprintf(util.NextSpanSeedURL, strconv.FormatUint(id, 10))))
	if err != nil {
		sp.Logger.Error("Error Fetching nextspanseed from HeimdallServer ", "error", err)
		return common.Hash{}, common.Address{}, err
	}

	sp.Logger.Info("Next span seed fetched")

	var nextSpanSeedResponse borTypes.QuerySpanSeedResponse

	if err = jsoniter.ConfigFastest.Unmarshal(response.Result, &nextSpanSeedResponse); err != nil {
		sp.Logger.Error("Error unmarshalling nextSpanSeed received from Heimdall Server", "error", err)
		return common.Hash{}, common.Address{}, err
	}

	return nextSpanSeedResponse.Seed, nextSpanSeedResponse.SeedAuthor, nil
}

// Stop stops all necessary go routines
func (sp *SpanProcessor) Stop() {
	// cancel span polling
	sp.cancelSpanService()
}
