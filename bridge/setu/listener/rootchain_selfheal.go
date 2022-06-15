package listener

import (
	"context"
	"math/big"
	"sync"
	"time"

	"github.com/maticnetwork/heimdall/bridge/setu/util"
	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/helper"
)

// startStakeUpdateSelfHealer starts the process to continuously re-process missing StakeUpdate events.
func (rl *RootChainListener) startStakeUpdateSelfHealer(ctx context.Context) {
	ticker := time.NewTicker(15 * time.Second)

	for {
		select {
		case <-ticker.C:
			// Fetch all heimdall validators
			validatorSet, err := util.GetValidatorSet(rl.cliCtx)
			if err != nil {
				rl.Logger.Error("Error getting heimdall validators", "error", err)
				continue
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

					if err = rl.processStakeUpdate(ctx, id, nonce+1); err != nil {
						rl.Logger.Error("Error processing stake update for validator", "error", err, "id", id)
					}
				}(validator.ID.Uint64(), validator.Nonce)
			}
			wg.Wait()
		case <-ctx.Done():
			rl.Logger.Info("Stopping stake update worker")
			ticker.Stop()
			return
		}
	}
}
func (rl *RootChainListener) processStakeUpdate(ctx context.Context, validatorId, nonce uint64) error {
	rl.Logger.Info("Processing stake update for validator", "id", validatorId, "nonce", nonce)

	var stakeUpdate *stakinginfo.StakinginfoStakeUpdate
	if err := helper.ExponentialBackoff(func() error {
		var err error
		stakeUpdate, err = rl.getStakeUpdate(ctx, validatorId, nonce)
		return err
	}, 3, time.Second); err != nil {
		rl.Logger.Error("Error getting stake update for validator", "error", err, "id", validatorId)
		return err
	}

	blockTime, err := rl.getBlockTime(ctx, stakeUpdate.Raw.BlockNumber)
	if err != nil {
		rl.Logger.Error("Unable to get block time with err", "error", err)
		return err
	}

	if time.Since(blockTime) < time.Hour {
		rl.Logger.Info("Block time is less than an hour, skipping stake-update")
		return nil
	}

	rl.handleLog(stakeUpdate.Raw)
	return nil
}

func (rl *RootChainListener) getBlockTime(ctx context.Context, blockNumber uint64) (time.Time, error) {
	blockBig := big.NewInt(0).SetUint64(blockNumber)
	block, err := rl.contractConnector.MainChainClient.BlockByNumber(ctx, blockBig)
	if err != nil {
		return time.Time{}, err
	}

	return time.Unix(int64(block.Time()), 0), nil
}
