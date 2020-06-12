package processor

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	"github.com/maticnetwork/heimdall/helper"

	borTypes "github.com/maticnetwork/heimdall/bor/types"

	"github.com/maticnetwork/heimdall/types"
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
	if err == nil && lastSpan != nil {
		sp.Logger.Debug("Found last span", "lastSpan", lastSpan.ID, "startBlock", lastSpan.StartBlock, "endBlock", lastSpan.EndBlock)
		nextSpanMsg, err := sp.fetchNextSpanDetails(lastSpan.ID+1, lastSpan.EndBlock+1)

		// check if current user is among next span producers
		if err == nil && sp.isSpanProposer(nextSpanMsg.SelectedProducers) {
			go sp.propose(lastSpan, nextSpanMsg)
		} else {
			sp.Logger.Error("Unable to fetch next span details", "lastSpanId", lastSpan.ID)
			return
		}
	}
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

		//Get NextSpanSeed from HeimdallServer
		var seed common.Hash
		if seed, err = sp.fetchNextSpanSeed(); err != nil {
			sp.Logger.Info("Error while fetching next span seed from HeimdallServer", "err", err)
			return
		}

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
		if err := sp.txBroadcaster.BroadcastToHeimdall(msg); err != nil {
			sp.Logger.Error("Error while broadcasting span to heimdall", "spanId", nextSpanMsg.ID, "startBlock", nextSpanMsg.StartBlock, "endBlock", nextSpanMsg.EndBlock, "error", err)
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
	err = json.Unmarshal(result.Result, &lastSpan)
	if err != nil {
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
	if err = json.Unmarshal(result.Result, &msg); err != nil {
		sp.Logger.Error("Error unmarshalling propose tx msg ", "error", err)
		return nil, err
	}

	sp.Logger.Debug("◽ Generated proposer span msg", "msg", msg.String())
	return &msg, nil
}

// fetchNextSpanSeed - fetches seed for next span
func (sp *SpanProcessor) fetchNextSpanSeed() (nextSpanSeed common.Hash, err error) {
	sp.Logger.Info("Sending Rest call to Get Seed for next span")
	response, err := helper.FetchFromAPI(sp.cliCtx, helper.GetHeimdallServerEndpoint(util.NextSpanSeedURL))
	if err != nil {
		sp.Logger.Error("Error Fetching nextspanseed from HeimdallServer ", "error", err)
		return nextSpanSeed, err
	}
	sp.Logger.Info("Next span seed fetched")
	if err := json.Unmarshal(response.Result, &nextSpanSeed); err != nil {
		sp.Logger.Error("Error unmarshalling nextSpanSeed received from Heimdall Server", "error", err)
		return nextSpanSeed, err
	}
	return nextSpanSeed, nil
}

// OnStop stops all necessary go routines
func (sp *SpanProcessor) Stop() {

	// cancel span polling
	sp.cancelSpanService()

}
