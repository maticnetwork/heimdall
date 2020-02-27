package listener

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/bridge/setu/queue"
	"github.com/maticnetwork/heimdall/helper"

	sdk "github.com/cosmos/cosmos-sdk/types"

	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	clerkTypes "github.com/maticnetwork/heimdall/clerk/types"
)

const (
	heimdallLastBlockKey = "heimdall-last-block" // storage key
)

// HeimdallListener - Listens to and process events from heimdall
type HeimdallListener struct {
	BaseListener
}

// NewHeimdallListener - constructor func
func NewHeimdallListener() *HeimdallListener {
	return &HeimdallListener{}
}

// Start starts new block subscription
func (hl *HeimdallListener) Start() error {
	hl.Logger.Info("Starting")

	// create cancellable context
	_, cancelSubscription := context.WithCancel(context.Background())
	hl.cancelSubscription = cancelSubscription

	// create cancellable context
	headerCtx, cancelHeaderProcess := context.WithCancel(context.Background())
	hl.cancelHeaderProcess = cancelHeaderProcess

	// Heimdall pollIntervall = (minimal pollInterval of rootchain and matichain)
	pollInterval := helper.GetConfig().SyncerPollInterval
	if helper.GetConfig().CheckpointerPollInterval < helper.GetConfig().SyncerPollInterval {
		pollInterval = helper.GetConfig().CheckpointerPollInterval
	}

	hl.Logger.Info("Start polling for events", "pollInterval", pollInterval)
	hl.StartPolling(headerCtx, pollInterval)
	return nil
}

// ProcessHeader -
func (hl *HeimdallListener) ProcessHeader(*types.Header) {

}

// StartPolling - starts polling for heimdall events
func (hl *HeimdallListener) StartPolling(ctx context.Context, pollInterval time.Duration) {
	// How often to fire the passed in function in second
	interval := pollInterval

	// Setup the ticket and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)

	var eventTypes []string
	eventTypes = append(eventTypes, "message.action='checkpoint'")
	eventTypes = append(eventTypes, "message.action='checkpoint-ack'")
	eventTypes = append(eventTypes, "message.action='event-record'")

	// start listening
	for {
		select {
		case <-ticker.C:
			fromBlock, toBlock := hl.fetchFromAndToBlock()
			if fromBlock < toBlock {
				for _, eventType := range eventTypes {
					var query []string
					query = append(query, eventType)
					query = append(query, fmt.Sprintf("tx.height>=%v", fromBlock))
					query = append(query, fmt.Sprintf("tx.height<=%v", toBlock))

					hl.Logger.Debug(" heimdall event search query", "query", query)
					searchResult, err := helper.QueryTxsByEvents(hl.cliCtx, query, 1, 50)
					if err != nil {
						hl.Logger.Error("Error while searching events", "eventType", eventType, "error", err)
						break
					}

					for _, tx := range searchResult.Txs {
						for _, log := range tx.Logs {
							event := helper.FilterEvents(log.Events, func(et sdk.StringEvent) bool {
								return et.Type == checkpointTypes.EventTypeCheckpoint || et.Type == checkpointTypes.EventTypeCheckpointAck || et.Type == clerkTypes.EventTypeRecord
							})
							if event != nil {
								hl.ProcessEvent(*event)
							}
						}
					}
				}
				// set last block to storage
				hl.storageClient.Put([]byte(heimdallLastBlockKey), []byte(strconv.FormatUint(toBlock, 10)), nil)
			}

		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (hl *HeimdallListener) fetchFromAndToBlock() (fromBlock uint64, toBlock uint64) {
	// toBlock - get latest blockheight from heimdall node
	nodeStatus, _ := helper.GetNodeStatus(hl.cliCtx)
	toBlock = uint64(nodeStatus.SyncInfo.LatestBlockHeight)

	// fromBlock - get last block from storage
	hasLastBlock, _ := hl.storageClient.Has([]byte(heimdallLastBlockKey), nil)
	if hasLastBlock {
		lastBlockBytes, err := hl.storageClient.Get([]byte(heimdallLastBlockKey), nil)
		if err != nil {
			hl.Logger.Info("Error while fetching last block bytes from storage", "error", err)
			return
		}

		if result, err := strconv.ParseUint(string(lastBlockBytes), 10, 64); err == nil {
			hl.Logger.Debug("Got last block from bridge storage", "lastBlock", uint64(result))
			fromBlock = uint64(result) + 1
		} else {
			hl.Logger.Info("Error parsing last block bytes from storage", "error", err)
			toBlock = 0
			return
		}
	}
	return
}

// ProcessEvent - process event from heimdall.
func (hl *HeimdallListener) ProcessEvent(event sdk.StringEvent) {
	hl.Logger.Info("Received Event", "EventType", event.Type)
	eventBytes, _ := json.Marshal(event)

	switch event.Type {

	case clerkTypes.EventTypeRecord:
		if err := hl.queueConnector.PublishMsg(eventBytes, queue.ClerkQueueRoute, hl.String(), event.Type); err != nil {
			hl.Logger.Error("Error publishing msg to clerk queue", "EventType", event.Type)
		}

	case checkpointTypes.EventTypeCheckpoint, checkpointTypes.EventTypeCheckpointAck:
		if err := hl.queueConnector.PublishMsg(eventBytes, queue.CheckpointQueueRoute, hl.String(), event.Type); err != nil {
			hl.Logger.Error("Error publishing msg to checkpoint queue", "EventType", event.Type)
		}

	default:
		hl.Logger.Info("EventType mismatch", "eventType", event.Type)

	}
}
