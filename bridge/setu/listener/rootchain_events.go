package listener

import (
	"context"
	"errors"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/maticnetwork/heimdall/helper"

	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/contracts/statesender"
)

const (
	blocksRange   = 1000
	maxIterations = 100
)

var (
	errNoEventsFound = errors.New("no events found")
)

// getLatestStateID returns state ID from the latest StateSynced event
func (rl *RootChainListener) getLatestStateID(ctx context.Context) (*big.Int, error) {
	rootchainContext, err := rl.getRootChainContext()
	if err != nil {
		return nil, err
	}

	latestEvent, err := rl.getLatestEvent(ctx, ethereum.FilterQuery{
		Addresses: []common.Address{
			rootchainContext.ChainmanagerParams.ChainParams.StateSenderAddress.EthAddress(),
		},
		Topics: [][]common.Hash{
			{statesender.GetStateSyncedEventID()},
			{},
			{},
		},
	})
	if err != nil {
		return nil, err
	}

	var event statesender.StatesenderStateSynced
	if err = helper.UnpackLog(rl.stateSenderAbi, &event, "StateSynced", latestEvent); err != nil {
		return nil, err
	}

	return event.Id, nil
}

// getCurrentStateID returns the current state ID handled by the polygon chain
func (rl *RootChainListener) getCurrentStateID(ctx context.Context) (*big.Int, error) {
	rootchainContext, err := rl.getRootChainContext()
	if err != nil {
		return nil, err
	}

	stateReceiverInstance, err := rl.contractConnector.GetStateReceiverInstance(
		rootchainContext.ChainmanagerParams.ChainParams.StateReceiverAddress.EthAddress(),
	)
	if err != nil {
		return nil, err
	}

	stateId, err := stateReceiverInstance.LastStateId(&bind.CallOpts{Context: ctx})
	if err != nil {
		return nil, err
	}

	return stateId, nil
}

// getStateSync returns the StateSynced event based on the given state ID
func (rl *RootChainListener) getStateSync(ctx context.Context, stateId int64) (*statesender.StatesenderStateSynced, error) {
	rootchainContext, err := rl.getRootChainContext()
	if err != nil {
		return nil, err
	}

	events, err := rl.contractConnector.MainChainClient.FilterLogs(ctx, ethereum.FilterQuery{
		Addresses: []common.Address{
			rootchainContext.ChainmanagerParams.ChainParams.StateSenderAddress.EthAddress(),
		},
		Topics: [][]common.Hash{
			{statesender.GetStateSyncedEventID()},
			{common.BytesToHash(big.NewInt(stateId).Bytes())},
			{},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return nil, errNoEventsFound
	}

	var event statesender.StatesenderStateSynced
	if err = helper.UnpackLog(rl.stateSenderAbi, &event, "StateSynced", &events[0]); err != nil {
		return nil, err
	}

	return &event, nil
}

// getLatestNonce returns the nonce from the latest StakeUpdate event
func (rl *RootChainListener) getLatestNonce(ctx context.Context, validatorId uint64) (uint64, error) {
	rootchainContext, err := rl.getRootChainContext()
	if err != nil {
		return 0, err
	}

	latestEvent, err := rl.getLatestEvent(ctx, ethereum.FilterQuery{
		Addresses: []common.Address{
			rootchainContext.ChainmanagerParams.ChainParams.StakingInfoAddress.EthAddress(),
		},
		Topics: [][]common.Hash{
			{stakinginfo.GetStakeUpdateEventID()},
			{common.BytesToHash(big.NewInt(0).SetUint64(validatorId).Bytes())},
			{},
			{},
		},
	})
	if err != nil {
		return 0, err
	}

	var event stakinginfo.StakinginfoStakeUpdate
	if err = helper.UnpackLog(rl.stakingInfoAbi, &event, "StakeUpdate", latestEvent); err != nil {
		return 0, err
	}

	return event.Nonce.Uint64(), nil
}

// getStakeUpdate returns StakeUpdate event based on the given validator ID and nonce
func (rl *RootChainListener) getStakeUpdate(ctx context.Context, validatorId, nonce uint64) (*stakinginfo.StakinginfoStakeUpdate, error) {
	rootchainContext, err := rl.getRootChainContext()
	if err != nil {
		return nil, err
	}

	events, err := rl.contractConnector.MainChainClient.FilterLogs(ctx, ethereum.FilterQuery{
		Addresses: []common.Address{
			rootchainContext.ChainmanagerParams.ChainParams.StakingInfoAddress.EthAddress(),
		},
		Topics: [][]common.Hash{
			{stakinginfo.GetStakeUpdateEventID()},
			{common.BytesToHash(big.NewInt(0).SetUint64(validatorId).Bytes())},
			{common.BytesToHash(big.NewInt(0).SetUint64(nonce).Bytes())},
			{},
		},
	})
	if err != nil {
		return nil, err
	}

	if len(events) == 0 {
		return nil, errNoEventsFound
	}

	var event stakinginfo.StakinginfoStakeUpdate
	if err = helper.UnpackLog(rl.stakingInfoAbi, &event, "StakeUpdate", &events[0]); err != nil {
		return nil, err
	}

	return &event, nil
}

// getLatestEvent returns the latest event based on the given filters
func (rl *RootChainListener) getLatestEvent(ctx context.Context, filters ethereum.FilterQuery) (*types.Log, error) {
	currentBlock, err := rl.contractConnector.MainChainClient.BlockByNumber(ctx, nil)
	if err != nil {
		return nil, err
	}

	currentBlockNumber := currentBlock.Number().Uint64()
	fromBlockNumber := currentBlockNumber - blocksRange
	toBlockNumber := currentBlockNumber

	var latestEvent *types.Log

	for i := 0; i < maxIterations; i++ {
		filters.FromBlock = big.NewInt(0).SetUint64(fromBlockNumber)
		filters.ToBlock = big.NewInt(0).SetUint64(toBlockNumber)

		events, err := rl.contractConnector.MainChainClient.FilterLogs(ctx, filters)
		if err != nil {
			return nil, err
		}

		if len(events) == 0 {
			toBlockNumber = fromBlockNumber
			fromBlockNumber -= blocksRange

			continue
		}

		latestEvent = &events[len(events)-1]

		break
	}

	if latestEvent == nil {
		return nil, errNoEventsFound
	}

	return latestEvent, nil
}
