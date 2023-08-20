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
)

// StakeUpdate represents the StakeUpdate event
type stakeUpdate struct {
	ValidatorID     string `json:"validatorId"`
	TotalStaked     string `json:"totalStaked"`
	Block           string `json:"block"`
	Nonce           string `json:"nonce"`
	TransactionHash string `json:"transactionHash"`
	LogIndex        string `json:"logIndex"`
}

// StateSync represents the StateSync event
type stateSync struct {
	StateID         string `json:"stateId"`
	Contract        string `json:"contract"`
	RawData         string `json:"rawData"`
	LogIndex        string `json:"logIndex"`
	TransactionHash string `json:"transactionHash"`
	BlockNumber     string `json:"blockNumber"`
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

func (rl *RootChainListener) querySubGraph(query []byte, ctx context.Context) (data []byte, err error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodPost, rl.subGraph.graphUrl, bytes.NewBuffer(query))
	if err != nil {
		return nil, err
	}

	response, err := rl.subGraph.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	return io.ReadAll(response.Body)
}
