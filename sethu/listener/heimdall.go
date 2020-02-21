package listener

import (
	"context"
	"encoding/json"
	"time"

	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/sethu/util"

	sdk "github.com/cosmos/cosmos-sdk/types"

	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
)

// HeimdallListener syncs validators and checkpoints
type HeimdallListener struct {
	BaseListener
}

// Start starts new block subscription
func (hl *HeimdallListener) Start() error {
	hl.Logger.Info("Starting listener", "name", hl.String())

	// create cancellable context
	ctx, _ := context.WithCancel(context.Background())
	hl.StartPolling(ctx, time.Minute)
	return nil
}

func (hl *HeimdallListener) ProcessHeader(*types.Header) {

}

// startPolling starts polling
func (hl *HeimdallListener) StartPolling(ctx context.Context, pollInterval time.Duration) {
	// How often to fire the passed in function in second
	interval := pollInterval

	hl.Logger.Info("Starting polling process")
	// Setup the ticket and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)

	var eventTypes []string
	eventTypes = append(eventTypes, "message.action='checkpoint'")
	eventTypes = append(eventTypes, "message.action='checkpoint-ack'")

	// start listening
	for {
		select {
		case <-ticker.C:
			for _, eventType := range eventTypes {
				searchResult, err := helper.QueryTxsByEvents(hl.cliCtx, []string{eventType}, 1, 50)
				if err != nil {
					hl.Logger.Error("Error searching for heimdall events", "error", err)
				}
				hl.Logger.Info("search successful", "eventType", eventType, "searchResultCount", searchResult.Count)

				for _, tx := range searchResult.Txs {
					for _, log := range tx.Logs {
						event := helper.FilterEvents(log.Events, func(et sdk.StringEvent) bool {
							return et.Type == checkpointTypes.EventTypeCheckpoint || et.Type == checkpointTypes.EventTypeCheckpointAck
						})
						if event != nil {
							hl.ProcessEvent(*event)
						}
					}
				}

			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (hl *HeimdallListener) ProcessEvent(event sdk.StringEvent) {
	hl.Logger.Info("Received Event from heimdall", "Event", event)
	eventBytes, _ := json.Marshal(event)
	hl.queueConnector.PublishMsg(eventBytes, util.CheckpointQueueRoute, hl.String(), event.Type)

}
