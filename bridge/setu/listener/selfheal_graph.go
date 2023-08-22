package listener

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/maticnetwork/heimdall/contracts/stakinginfo"
	"github.com/maticnetwork/heimdall/contracts/statesender"
	"github.com/maticnetwork/heimdall/helper"
)

// StakeUpdate represents the StakeUpdate event
type stakeUpdate struct {
	Nonce           string `json:"nonce"`
	TransactionHash string `json:"transactionHash"`
	LogIndex        string `json:"logIndex"`
}

// StateSync represents the StateSync event
type stateSync struct {
	StateID         string `json:"stateId"`
	LogIndex        string `json:"logIndex"`
	TransactionHash string `json:"transactionHash"`
}

type stakeUpdateResponse struct {
	Data struct {
		StakeUpdates []stakeUpdate `json:"stakeUpdates"`
	} `json:"data"`
}

type stateSyncResponse struct {
	Data struct {
		StateSyncs []stateSync `json:"stateSyncs"`
	} `json:"data"`
}

func (rl *RootChainListener) querySubGraph(query []byte, ctx context.Context) (data []byte, err error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, rl.subGraph.graphUrl, bytes.NewBuffer(query))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := rl.subGraph.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return io.ReadAll(response.Body)
}

func (rl *RootChainListener) getLatestNonceGraph(ctx context.Context, validatorId uint64) (uint64, error) {
	query := map[string]string{
		"query": `
		{
			stakeUpdates(first:1, orderBy: nonce, orderDirection : desc, where: {validatorId: ` + strconv.Itoa(int(validatorId)) + `}){
				nonce
		   } 
		}   
		`,
	}

	byteQuery, err := json.Marshal(query)
	if err != nil {
		return 0, err
	}

	data, err := rl.querySubGraph(byteQuery, ctx)
	if err != nil {
		return 0, fmt.Errorf("unable to fetch latest nonce from graph with err: %s", err)
	}

	var response stakeUpdateResponse
	if err = json.Unmarshal(data, &response); err != nil {
		return 0, err
	}

	if len(response.Data.StakeUpdates) == 0 {
		return 0, nil
	}

	latestValidatorNonce, err := strconv.Atoi(response.Data.StakeUpdates[0].Nonce)
	if err != nil {
		return 0, err
	}

	return uint64(latestValidatorNonce), nil
}

func (rl *RootChainListener) getLatestStateIDGraph(ctx context.Context) (*big.Int, error) {
	query := map[string]string{
		"query": `
		{
			stateSyncs(first : 1, orderBy : stateId, orderDirection : desc) {
				stateId
			}
		}
		`,
	}

	byteQuery, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	data, err := rl.querySubGraph(byteQuery, ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch latest state id from graph with err: %s", err)
	}

	var response stateSyncResponse
	if err = json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("unable to unmarshal graph response: %s", err)
	}

	if len(response.Data.StateSyncs) == 0 {
		return big.NewInt(0), nil
	}

	stateID := big.NewInt(0)
	stateID.SetString(response.Data.StateSyncs[0].StateID, 10)

	return stateID, nil
}

func (rl *RootChainListener) getStakeUpdateGraph(ctx context.Context, validatorId, nonce uint64) (*stakinginfo.StakinginfoStakeUpdate, error) {
	query := map[string]string{
		"query": `
		{
			stakeUpdates(where: {validatorId: ` + strconv.Itoa(int(validatorId)) + `, nonce: ` + strconv.Itoa(int(nonce)) + `}){
				transactionHash
				logIndex
		   } 
		}   
		`,
	}

	byteQuery, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	data, err := rl.querySubGraph(byteQuery, ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch stake update from graph with err: %s", err)
	}

	var response stakeUpdateResponse
	if err = json.Unmarshal(data, &response); err != nil {
		return nil, err
	}

	if len(response.Data.StakeUpdates) == 0 {
		return nil, fmt.Errorf("no stake update found for validator %d and nonce %d", validatorId, nonce)
	}

	receipt, err := rl.contractConnector.MainChainClient.TransactionReceipt(ctx, common.HexToHash(response.Data.StakeUpdates[0].TransactionHash))
	if err != nil {
		return nil, err
	}

	for _, logs := range receipt.Logs {
		if strconv.Itoa(int(logs.Index)) == response.Data.StakeUpdates[0].LogIndex {
			var event stakinginfo.StakinginfoStakeUpdate
			if err = helper.UnpackLog(rl.stakingInfoAbi, &event, stakeUpdateEvent, logs); err != nil {
				return nil, err
			}

			return &event, nil
		}
	}

	return nil, fmt.Errorf("no logs found for given log index %s ,validator %d and nonce %d", response.Data.StakeUpdates[0].LogIndex, validatorId, nonce)
}

func (rl *RootChainListener) getStateSyncGraph(ctx context.Context, stateId int64) (*statesender.StatesenderStateSynced, error) {
	query := map[string]string{
		"query": `
		{
			stateSyncs(where: {stateId: ` + strconv.Itoa(int(stateId)) + `}) {
				logIndex
				transactionHash
			}
		}
		`,
	}

	byteQuery, err := json.Marshal(query)
	if err != nil {
		return nil, err
	}

	data, err := rl.querySubGraph(byteQuery, ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch latest state id from graph with err: %s", err)
	}

	var response stateSyncResponse
	if err = json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("unable to unmarshal graph response: %s", err)
	}

	if len(response.Data.StateSyncs) == 0 {
		return nil, fmt.Errorf("no state sync found for state id %d", stateId)
	}

	receipt, err := rl.contractConnector.MainChainClient.TransactionReceipt(ctx, common.HexToHash(response.Data.StateSyncs[0].TransactionHash))
	if err != nil {
		return nil, err
	}

	for _, logs := range receipt.Logs {
		if strconv.Itoa(int(logs.Index)) == response.Data.StateSyncs[0].LogIndex {
			var event statesender.StatesenderStateSynced
			if err = helper.UnpackLog(rl.stateSenderAbi, &event, stateSyncedEvent, logs); err != nil {
				return nil, err
			}

			return &event, nil
		}
	}

	return nil, fmt.Errorf("no logs found for given log index %s and state id %d", response.Data.StateSyncs[0].LogIndex, stateId)
}
