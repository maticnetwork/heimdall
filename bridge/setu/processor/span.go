package processor

// import (
// 	"encoding/json"

// 	"github.com/maticnetwork/bor/common"
// 	"github.com/maticnetwork/heimdall/bridge/setu/util"
// 	"github.com/maticnetwork/heimdall/helper"
// 	hmTypes "github.com/maticnetwork/heimdall/types"

// 	"github.com/maticnetwork/bor/core/types"
// 	borTypes "github.com/maticnetwork/heimdall/bor/types"
// )

// // SpanProcessor - process span related events
// type SpanProcessor struct {
// 	BaseProcessor
// }

// // Start starts new block subscription
// func (sp *SpanProcessor) Start() error {
// 	sp.Logger.Info("Starting")
// 	return nil
// }

// // RegisterTasks
// func (sp *SpanProcessor) RegisterTasks() {
// 	sp.Logger.Info("Registering span tasks")
// 	if err := sp.queueConnector.Server.RegisterTask("sendSpanToHeimdall", sp.sendSpanToHeimdall); err != nil {
// 		sp.Logger.Error("RegisterTasks | sendSpanToHeimdall", "error", err)
// 	}
// }

// // HandleSendSpanTask - handle send span task
// // 1. check if this span has to be broadcasted to heimdall
// // 2. create and broadcast  span transaction to heimdall
// func (sp *SpanProcessor) sendSpanToHeimdall(headerBlockStr string) error {
// 	var header = types.Header{}
// 	if err := header.UnmarshalJSON([]byte(headerBlockStr)); err != nil {
// 		sp.Logger.Error("Error while unmarshalling the header block", "error", err)
// 		return err
// 	}

// 	// Fetch last span
// 	lastSpan, err := util.GetLastSpan(sp.cliCtx)
// 	if err == nil && lastSpan != nil {
// 		sp.Logger.Debug("Found last span", "lastSpan", lastSpan.ID, "startBlock", lastSpan.StartBlock, "endBlock", lastSpan.EndBlock)

// 	}

// 	// check and propose span
// 	if lastSpan.StartBlock <= header.Number.Uint64() && header.Number.Uint64() <= lastSpan.EndBlock {

// 		// Fetch next span to be proposed
// 		nextSpanMsg, err := util.FetchNextSpanDetails(sp.cliCtx, lastSpan.ID+1, lastSpan.EndBlock+1)
// 		if err != nil {
// 			sp.Logger.Error("Error while fetching next span details", "error", err)
// 			return err
// 		}

// 		// Get NextSpanSeed from HeimdallServer
// 		var seed common.Hash
// 		if seed, err = sp.fetchNextSpanSeed(); err != nil {
// 			sp.Logger.Info("Error while fetching next span seed from HeimdallServer", "err", err)
// 			return err
// 		}

// 		// log new span
// 		sp.Logger.Info("✅ Proposing new span", "spanId", nextSpanMsg.ID, "startBlock", nextSpanMsg.StartBlock, "endBlock", nextSpanMsg.EndBlock, "seed", seed)

// 		// broadcast to heimdall
// 		msg := borTypes.MsgProposeSpan{
// 			ID:         nextSpanMsg.ID,
// 			Proposer:   hmTypes.BytesToHeimdallAddress(helper.GetAddress()),
// 			StartBlock: nextSpanMsg.StartBlock,
// 			EndBlock:   nextSpanMsg.EndBlock,
// 			ChainID:    nextSpanMsg.ChainID,
// 			Seed:       seed,
// 		}

// 		// return broadcast to heimdall
// 		if err := sp.txBroadcaster.BroadcastToHeimdall(msg); err != nil {
// 			sp.Logger.Error("Error while broadcasting span to heimdall", "spanId", nextSpanMsg.ID, "startBlock", nextSpanMsg.StartBlock, "endBlock", nextSpanMsg.EndBlock, "error", err)
// 			return err
// 		}
// 	}

// 	return nil
// }

// // fetchNextSpanSeed - fetches seed for next span
// func (sp *SpanProcessor) fetchNextSpanSeed() (nextSpanSeed common.Hash, err error) {
// 	sp.Logger.Debug("Sending Rest call to Get Seed for next span")
// 	response, err := helper.FetchFromAPI(sp.cliCtx, helper.GetHeimdallServerEndpoint(util.NextSpanSeedURL))
// 	if err != nil {
// 		sp.Logger.Error("Error Fetching nextspanseed from HeimdallServer ", "error", err)
// 		return nextSpanSeed, err
// 	}
// 	sp.Logger.Debug("Next span seed fetched")
// 	if err := json.Unmarshal(response.Result, &nextSpanSeed); err != nil {
// 		sp.Logger.Error("Error unmarshalling nextSpanSeed received from Heimdall Server", "error", err)
// 		return nextSpanSeed, err
// 	}
// 	return nextSpanSeed, nil
// }
