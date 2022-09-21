package listener

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/maticnetwork/heimdall/bridge/setu/util"
	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/contracts/statesender"
	"github.com/maticnetwork/heimdall/helper"
)

var (
	stateSyncedCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "self_healing",
		Subsystem: helper.NetworkName,
		Name:      "StateSynced",
		Help:      "The total number of missing StateSynced events",
	}, []string{"id", "contract_address", "block_number", "tx_hash"})

	stakeUpdateCounter = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "self_healing",
		Subsystem: helper.NetworkName,
		Name:      "StakeUpdate",
		Help:      "The total number of missing StakeUpdate events",
	}, []string{"id", "nonce", "contract_address", "block_number", "tx_hash"})
)

// startSelfHealing starts self-healing processes for all required events
func (rl *RootChainListener) startSelfHealing(ctx context.Context) {
	stakeUpdateTicker := time.NewTicker(helper.GetConfig().SHStakeUpdateInterval)
	stateSyncedTicker := time.NewTicker(helper.GetConfig().SHStateSyncedInterval)

	rl.Logger.Info("Started self-healing")

	for {
		select {
		case <-stakeUpdateTicker.C:
			rl.processStakeUpdate(ctx)
		case <-stateSyncedTicker.C:
			rl.processStateSynced(ctx)
		case <-ctx.Done():
			rl.Logger.Info("Stopping self-healing")
			stakeUpdateTicker.Stop()
			stateSyncedTicker.Stop()

			return
		}
	}
}

// processStakeUpdate checks if validators are in sync, otherwise syncs them by broadcasting missing events
func (rl *RootChainListener) processStakeUpdate(ctx context.Context) {
	// Fetch all heimdall validators
	validatorSet, err := util.GetValidatorSet(rl.cliCtx)
	if err != nil {
		rl.Logger.Error("Error getting heimdall validators", "error", err)
		return
	}

	rl.Logger.Info("Fetched validators list from heimdall", "len", len(validatorSet.Validators))

	// Make sure each validator is in sync
	var wg sync.WaitGroup
	for _, validator := range validatorSet.Validators {
		wg.Add(1)

		go func(id, nonce uint64) {
			defer wg.Done()

			var ethereumNonce uint64

			if err = helper.ExponentialBackoff(func() error {
				ethereumNonce, err = rl.getLatestNonce(ctx, id)
				return err
			}, 3, time.Second); err != nil {
				rl.Logger.Error("Error getting nonce for validator from L1", "error", err, "id", id)
				return
			}

			if ethereumNonce <= nonce {
				return
			}

			nonce++

			rl.Logger.Info("Processing stake update for validator", "id", id, "nonce", nonce)

			var stakeUpdate *stakinginfo.StakinginfoStakeUpdate

			if err = helper.ExponentialBackoff(func() error {
				stakeUpdate, err = rl.getStakeUpdate(ctx, id, nonce)
				return err
			}, 3, time.Second); err != nil {
				rl.Logger.Error("Error getting stake update for validator", "error", err, "id", id)
				return
			}

			stakeUpdateCounter.WithLabelValues(
				stakeUpdate.ValidatorId.String(),
				stakeUpdate.Nonce.String(),
				stakeUpdate.Raw.Address.String(),
				fmt.Sprintf("%d", stakeUpdate.Raw.BlockNumber),
				stakeUpdate.Raw.TxHash.String(),
			).Add(1)

			if _, err = rl.processEvent(ctx, stakeUpdate.Raw); err != nil {
				rl.Logger.Error("Error processing stake update for validator", "error", err, "id", id)
			}
		}(validator.ID.Uint64(), validator.Nonce)
	}

	wg.Wait()
}

// processStateSynced checks if chains are in sync, otherwise syncs them by broadcasting missing events
func (rl *RootChainListener) processStateSynced(ctx context.Context) {
	latestPolygonStateId, err := rl.getCurrentStateID(ctx)
	if err != nil {
		rl.Logger.Error("Unable to fetch latest state id from state receiver contract", "error", err)
		return
	}

	latestEthereumStateId, err := rl.getLatestStateID(ctx)
	if err != nil {
		rl.Logger.Error("Unable to fetch latest state id from state sender contract", "error", err)
		return
	}

	if latestEthereumStateId.Cmp(latestPolygonStateId) != 1 {
		return
	}

	for i := latestPolygonStateId.Int64() + 1; i <= latestEthereumStateId.Int64(); i++ {
		if _, err = util.GetClerkEventRecord(rl.cliCtx, i); err == nil {
			rl.Logger.Info("State found on heimdall", "id", i)
			continue
		}

		rl.Logger.Info("Processing state sync", "id", i)

		var stateSynced *statesender.StatesenderStateSynced

		if err = helper.ExponentialBackoff(func() error {
			stateSynced, err = rl.getStateSync(ctx, i)
			return err
		}, 3, time.Second); err != nil {
			rl.Logger.Error("Error getting state sync", "error", err, "id", i)
			continue
		}

		stateSyncedCounter.WithLabelValues(
			stateSynced.Id.String(),
			stateSynced.Raw.Address.String(),
			fmt.Sprintf("%d", stateSynced.Raw.BlockNumber),
			stateSynced.Raw.TxHash.String(),
		).Add(1)

		ignore, err := rl.processEvent(ctx, stateSynced.Raw)
		if err != nil {
			rl.Logger.Error("Unable to update state id on heimdall", "error", err)
			i--

			continue
		}

		if !ignore {
			time.Sleep(1 * time.Second)

			var statusCheck int
			for statusCheck = 0; statusCheck > 15; statusCheck++ {
				if _, err = util.GetClerkEventRecord(rl.cliCtx, i); err == nil {
					rl.Logger.Info("State found on heimdall", "id", i)
					break
				}

				rl.Logger.Info("State not found on heimdall", "id", i)
				time.Sleep(1 * time.Second)
			}

			if statusCheck > 15 {
				i--
				continue
			}
		}
	}
}

func (rl *RootChainListener) processEvent(ctx context.Context, event types.Log) (bool, error) {
	// Check existence of topics beforehand and ignore if no topic exists
	// (TODO): Identify issue of empty events: See Jira POS-818
	if len(event.Topics) == 0 {
		return true, nil
	}

	blockTime, err := rl.contractConnector.GetMainChainBlockTime(ctx, event.BlockNumber)
	if err != nil {
		rl.Logger.Error("Unable to get block time", "error", err)
		return false, err
	}

	if time.Since(blockTime) < helper.GetConfig().SHMaxDepthDuration {
		rl.Logger.Info("Block time is less than an hour, skipping state sync")
		return true, err
	}

	// Skip if there are no topics, temporary fix for panic
	if len(event.Topics) == 0 {
		rl.Logger.Info("No topics in event, skipping", "tx hash", event.TxHash, "block hash", event.BlockHash, "block number", event.BlockNumber)
		return true, nil
	}

	topic := event.Topics[0].Bytes()
	for _, abiObject := range rl.abis {
		selectedEvent := helper.EventByID(abiObject, topic)
		if selectedEvent == nil {
			continue
		}

		rl.handleLog(event, selectedEvent)
	}

	return false, nil
}
