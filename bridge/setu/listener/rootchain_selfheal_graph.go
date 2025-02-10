package listener

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"math/big"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	jsoniter "github.com/json-iterator/go"

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

// querySubGraph queries the subgraph and limits the read size
func (rl *RootChainListener) querySubGraph(query []byte, ctx context.Context) (data []byte, err error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, rl.subGraphClient.graphUrl, bytes.NewBuffer(query))
	if err != nil {
		return nil, err
	}

	request.Header.Set("Content-Type", "application/json")

	response, err := rl.subGraphClient.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Limit the number of bytes read from the response body
	limitedBody := http.MaxBytesReader(nil, response.Body, helper.APIBodyLimit)

	return io.ReadAll(limitedBody)
}

// getLatestStateID returns state ID from the latest StateSynced event
func (rl *RootChainListener) getLatestStateID(ctx context.Context) (*big.Int, error) {
	query := map[string]string{
		"query": `
		{
			stateSyncs(first : 1, orderBy : stateId, orderDirection : desc) {
				stateId
			}
		}
		`,
	}

	byteQuery, err := jsoniter.ConfigFastest.Marshal(query)
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
func (rl *RootChainListener) getStateSync(ctx context.Context, stateId int64) (*types.Log, error) {
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

	for _, log := range receipt.Logs {
		if log.Index > math.MaxInt {
			return nil, fmt.Errorf("log index value out of range for int: %d", log.Index)
		}
		if strconv.Itoa(int(log.Index)) == response.Data.StateSyncs[0].LogIndex {
			return log, nil
		}
	}

	return nil, fmt.Errorf("no log found for given log index %s and state id %d", response.Data.StateSyncs[0].LogIndex, stateId)
}

// getLatestNonce returns the nonce from the latest StakeUpdate event
func (rl *RootChainListener) getLatestNonce(ctx context.Context, validatorId uint64) (uint64, error) {
	if validatorId > math.MaxInt {
		return 0, fmt.Errorf("validator ID value out of range for int: %d", validatorId)
	}

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

	if latestValidatorNonce < 0 {
		return 0, fmt.Errorf("latest validator nonce is negative: %d", latestValidatorNonce)
	}

	return uint64(latestValidatorNonce), nil
}

// getStakeUpdate returns StakeUpdate event based on the given validator ID and nonce
func (rl *RootChainListener) getStakeUpdate(ctx context.Context, validatorId, nonce uint64) (*types.Log, error) {
	if validatorId > math.MaxInt {
		return nil, fmt.Errorf("validator ID value out of range for int: %d", validatorId)
	}
	if nonce > math.MaxInt {
		return nil, fmt.Errorf("nonce value out of range for int: %d", nonce)
	}
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

	for _, log := range receipt.Logs {
		if log.Index > math.MaxInt {
			return nil, fmt.Errorf("log index value out of range for int: %d", log.Index)
		}
		if strconv.Itoa(int(log.Index)) == response.Data.StakeUpdates[0].LogIndex {
			return log, nil
		}
	}

	return nil, fmt.Errorf("no log found for given log index %s ,validator %d and nonce %d", response.Data.StakeUpdates[0].LogIndex, validatorId, nonce)
}
