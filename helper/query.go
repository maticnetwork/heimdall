package helper

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	httpClient "github.com/tendermint/tendermint/rpc/client/http"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/cosmos/cosmos-sdk/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
)

const (
	// CommitTimeout commit timeout
	CommitTimeout = 2 * time.Minute
)

// QueryTxWithProof query tx with proof from node
func QueryTxWithProof(cliCtx client.Context, hash []byte) (*ctypes.ResultTx, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	return node.Tx(context.Background(), hash, true)
}

// GetNodeStatus returns node status
func GetNodeStatus(cliCtx client.Context) (*ctypes.ResultStatus, error) {
	node, err := cliCtx.GetNode()
	if err != nil {
		return nil, err
	}

	return node.Status(context.Background())
}

// GetBeginBlockEvents get block through per height
func GetBeginBlockEvents(client *httpClient.HTTP, height int64) ([]abci.Event, error) {
	c, cancel := context.WithTimeout(context.Background(), CommitTimeout)
	defer cancel()

	// get block using client
	blockResults, err := client.BlockResults(c, &height)
	fmt.Printf("blockResults %+v\n", blockResults)
	//fmt.Printf("BeginBlockEvents %+v\n", blockResults.BeginBlockEvents)

	if err == nil && blockResults != nil {
		return blockResults.BeginBlockEvents, nil
	}

	// subscriber
	subscriber := fmt.Sprintf("new-block-%v", height)

	// query for event
	query := tmTypes.QueryForEvent(tmTypes.EventNewBlock).String()

	// register for the next event of this type
	eventCh, err := client.Subscribe(c, subscriber, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to subscribe")
	}

	// unsubscribe query
	defer func() {
		if subErr := client.Unsubscribe(c, subscriber, query); subErr != nil {
			Logger.Error(fmt.Sprintf("Error while Unsubscribe the %s the subscriber ", query), "Subscriber", subscriber, "Err", err)
		}
	}()

	for {
		select {
		case event := <-eventCh:
			eventData := event.Data.(tmTypes.TMEventData)
			switch t := eventData.(type) {
			case tmTypes.EventDataNewBlock:
				if t.Block.Height == height {
					return t.ResultBeginBlock.GetEvents(), nil
				}
			default:
				return nil, errors.New("timed out waiting for event")
			}
		case <-c.Done():
			return nil, errors.New("timed out waiting for event")
		}
	}
}

// GetBlockWithClient get block through per height
func GetBlockWithClient(client *httpClient.HTTP, height int64) (*tmTypes.Block, error) {
	c, cancel := context.WithTimeout(context.Background(), CommitTimeout)
	defer cancel()

	// get block using client
	block, err := client.Block(c, &height)
	if err == nil && block != nil {
		return block.Block, nil
	}

	// subscriber
	subscriber := fmt.Sprintf("new-block-%v", height)

	// query for event
	query := tmTypes.QueryForEvent(tmTypes.EventNewBlock).String()

	// register for the next event of this type
	eventCh, err := client.Subscribe(c, subscriber, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to subscribe")
	}

	// unsubscribe query
	defer func() {
		if err := client.Unsubscribe(c, subscriber, query); err != nil {
			Logger.Error("GetBlockWithClient | Unsubscribe", "Error", err)
		}
	}()

	for {
		select {
		case event := <-eventCh:
			eventData := event.Data.(tmTypes.TMEventData)
			switch t := eventData.(type) {
			case tmTypes.EventDataNewBlock:
				if t.Block.Height == height {
					return t.Block, nil
				}
			default:
				return nil, errors.New("timed out waiting for event")
			}
		case <-c.Done():
			return nil, errors.New("timed out waiting for event")
		}
	}
}

// FetchSideTxSigs fetches side tx sigs from it
func FetchSideTxSigs(
	client *httpClient.HTTP,
	height int64,
	txHash []byte,
	sideTxData []byte,
) ([]byte, error) {
	// get block client
	blockDetails, err := GetBlockWithClient(client, height)

	if err != nil {
		return nil, err
	}

	// extract votes from response
	preCommits := blockDetails.LastCommit.Signatures

	// extract side-tx signs from votes
	sigs := GetSideTxSigs(txHash, sideTxData, preCommits)

	// return
	return sigs, nil
}
