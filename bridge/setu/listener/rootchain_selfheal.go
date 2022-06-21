package listener

import (
	"context"
	"sync"
	"time"

	"github.com/maticnetwork/bor/core/types"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"

	"github.com/maticnetwork/heimdall/bridge/setu/util"
	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/contracts/statesender"
	"github.com/maticnetwork/heimdall/helper"
)

var (
	meter              = global.Meter(helper.NetworkName + ".heimdall.self-healing")
	stateSyncedCounter syncint64.Counter
	stakeUpdateCounter syncint64.Counter
)

func init() {
	var err error

	if stateSyncedCounter, err = meter.SyncInt64().Counter(
		"StateSynced",
		instrument.WithUnit("0"),
		instrument.WithDescription("StateSynced missing event counter"),
	); err != nil {
		panic(err)
	}

	if stakeUpdateCounter, err = meter.SyncInt64().Counter(
		"StakeUpdate",
		instrument.WithUnit("0"),
		instrument.WithDescription("StakeUpdate missing event counter"),
	); err != nil {
		panic(err)
	}
}

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

			stakeUpdateCounter.Add(ctx, 1,
				attribute.String("id", stakeUpdate.ValidatorId.String()),
				attribute.String("nonce", stakeUpdate.Nonce.String()),
				attribute.String("contract_address", stakeUpdate.Raw.Address.String()),
				attribute.Int64("block_number", int64(stakeUpdate.Raw.BlockNumber)),
				attribute.String("tx_hash", stakeUpdate.Raw.TxHash.String()),
			)

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

		stateSyncedCounter.Add(ctx, 1,
			attribute.String("id", stateSynced.Id.String()),
			attribute.String("contract_address", stateSynced.Raw.Address.String()),
			attribute.Int64("block_number", int64(stateSynced.Raw.BlockNumber)),
			attribute.String("tx_hash", stateSynced.Raw.TxHash.String()),
		)

		var ignore bool
		if ignore, err = rl.processEvent(ctx, stateSynced.Raw); err != nil {
			rl.Logger.Error("Unable to update state id on heimdall", "error", err)
			i--
			continue
		}

		if !ignore {
			time.Sleep(1 * time.Second)

			statusCheck := 0
			for {
				if _, err = util.GetClerkEventRecord(rl.cliCtx, i); err == nil {
					rl.Logger.Info("State found on heimdall", "id", i)
					break
				}

				rl.Logger.Info("State not found on heimdall", "id", i)
				time.Sleep(1 * time.Second)

				if statusCheck++; statusCheck > 15 {
					break
				}
			}

			if statusCheck > 15 {
				i--
				continue
			}
		}
	}
}

func (rl *RootChainListener) processEvent(ctx context.Context, event types.Log) (bool, error) {
	blockTime, err := rl.contractConnector.GetMainChainBlockTime(ctx, event.BlockNumber)
	if err != nil {
		rl.Logger.Error("Unable to get block time", "error", err)
		return false, err
	}

	if time.Since(blockTime) < helper.GetConfig().SHMaxDepthDuration {
		rl.Logger.Info("Block time is less than an hour, skipping state sync")
		return true, err
	}

	pubkey := helper.GetPubKey()
	pubkeyBytes := pubkey[1:]

	topic := event.Topics[0].Bytes()
	for _, abiObject := range rl.abis {
		selectedEvent := helper.EventByID(abiObject, topic)
		if selectedEvent == nil {
			continue
		}

		rl.handleLog(event, selectedEvent, pubkeyBytes)
	}

	return false, nil
}
