package listener

import (
	"context"
	"errors"
	"math/big"

	ethereum "github.com/maticnetwork/bor"
	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/bor/core/types"
	"github.com/maticnetwork/heimdall/helper"

	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/contracts/statesender"
)

const (
	stateSyncedEventID  = "0x103fed9db65eac19c4d870f49ab7520fe03b99f1838e5996caf47e9e43308392"
	stakeUpdatedEventID = "0x35af9eea1f0e7b300b0a14fae90139a072470e44daa3f14b5069bebbc1265bda"
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
			{common.HexToHash(stateSyncedEventID)},
			{},
			{},
		},
	})
	if err != nil {
		return nil, err
	}

	return big.NewInt(0).SetBytes(latestEvent.Topics[1].Bytes()), nil
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
			{common.HexToHash(stateSyncedEventID)},
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

	return &statesender.StatesenderStateSynced{
		Id:              big.NewInt(0).SetBytes(events[0].Topics[1].Bytes()),
		ContractAddress: common.BytesToAddress(events[0].Topics[2].Bytes()),
		Data:            events[0].Data,
		Raw:             events[0],
	}, nil
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
			{common.HexToHash(stakeUpdatedEventID)},
			{common.BytesToHash(big.NewInt(0).SetUint64(validatorId).Bytes())},
			{},
			{},
		},
	})
	if err != nil {
		return 0, err
	}

	return big.NewInt(0).SetBytes(latestEvent.Topics[2].Bytes()).Uint64(), nil
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
			{common.HexToHash(stakeUpdatedEventID)},
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
	const blocksRange = 1000
	const maxIterations = 100

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

		var events []types.Log
		if events, err = rl.contractConnector.MainChainClient.FilterLogs(ctx, filters); err != nil {
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
