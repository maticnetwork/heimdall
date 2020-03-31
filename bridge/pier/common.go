package pier

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/log"
	httpClient "github.com/tendermint/tendermint/rpc/client"
	tmTypes "github.com/tendermint/tendermint/types"

	chainManagerTypes "github.com/maticnetwork/heimdall/chainmanager/types"
	"github.com/maticnetwork/heimdall/helper"
	hmtypes "github.com/maticnetwork/heimdall/types"
)

const (
	ChainSyncer          = "chain-syncer"
	HeimdallCheckpointer = "heimdall-checkpointer"
	NoackService         = "checkpoint-no-ack"
	SpanServiceStr       = "span-service"
	ClerkServiceStr      = "clerk-service"
	AMQPConsumerService  = "amqp-consumer-service"

	// TxsURL represents txs url
	TxsURL = "/txs"

	AccountDetailsURL      = "/auth/accounts/%v"
	LastNoAckURL           = "/checkpoint/last-no-ack"
	CheckpointParamsURL    = "/checkpoint/params"
	ChainManagerParamsURL  = "/chainmanager/params"
	ProposersURL           = "/staking/proposer/%v"
	BufferedCheckpointURL  = "/checkpoint/buffer"
	LatestCheckpointURL    = "/checkpoint/latest-checkpoint"
	CurrentProposerURL     = "/staking/current-proposer"
	LatestSpanURL          = "/bor/latest-span"
	NextSpanInfoURL        = "/bor/prepare-next-span"
	DividendAccountRootURL = "/staking/dividend-account-root"
	ValidatorURL           = "/staking/validator/%v"

	TransactionTimeout = 1 * time.Minute
	CommitTimeout      = 2 * time.Minute

	BridgeDBFlag = "bridge-db"
)

// Global logger for bridge
var Logger log.Logger

func init() {
	Logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
}

// checks if we are proposer
func isProposer(cliCtx cliContext.CLIContext) bool {
	var proposers []hmtypes.Validator
	count := uint64(1)
	result, err := helper.FetchFromAPI(cliCtx,
		helper.GetHeimdallServerEndpoint(fmt.Sprintf(ProposersURL, strconv.FormatUint(count, 10))),
	)
	if err != nil {
		Logger.Error("Error fetching proposers", "error", err)
		return false
	}
	err = json.Unmarshal(result.Result, &proposers)
	if err != nil {
		Logger.Error("error unmarshalling proposer slice", "error", err)
		return false
	}
	Logger.Debug("Current proposer fetched", "validator", proposers[0].String())

	if bytes.Equal(proposers[0].Signer.Bytes(), helper.GetAddress()) {
		return true
	}

	return false
}

// check if we are the EventSender
func isEventSender(cliCtx cliContext.CLIContext, validatorID uint64) bool {

	var validator hmtypes.Validator

	result, err := helper.FetchFromAPI(cliCtx,
		helper.GetHeimdallServerEndpoint(fmt.Sprintf(ValidatorURL, strconv.FormatUint(validatorID, 10))),
	)
	if err != nil {
		Logger.Error("Error fetching proposers", "error", err)
		return false
	}
	err = json.Unmarshal(result.Result, &validator)
	if err != nil {
		Logger.Error("error unmarshalling proposer slice", "error", err)
		return false
	}
	Logger.Debug("Current event sender received", "validator", validator.String())

	if bytes.Equal(validator.Signer.Bytes(), helper.GetAddress()) {
		return true
	}

	return false

}

// WaitForOneEvent subscribes to a websocket event for the given
// event time and returns upon receiving it one time, or
// when the timeout duration has expired.
//
// This handles subscribing and unsubscribing under the hood
func WaitForOneEvent(tx tmTypes.Tx, client *httpClient.HTTP) (tmTypes.TMEventData, error) {
	ctx, cancel := context.WithTimeout(context.Background(), CommitTimeout)
	defer cancel()

	// subscriber
	subscriber := hex.EncodeToString(tx.Hash())

	// query
	query := tmTypes.EventQueryTxFor(tx).String()

	// register for the next event of this type
	eventCh, err := client.Subscribe(ctx, subscriber, query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to subscribe")
	}

	// make sure to unregister after the test is over
	defer client.UnsubscribeAll(ctx, subscriber)
	select {
	case event := <-eventCh:
		return event.Data.(tmTypes.TMEventData), nil
	case <-ctx.Done():
		return nil, errors.New("timed out waiting for event")
	}
}

// FetchVotes fetches votes and extracts sigs from it
func FetchVotes(
	height int64,
	client *httpClient.HTTP,
) (votes []*tmTypes.CommitSig, sigs []byte, chainID string, err error) {
	// get block client
	blockDetails, err := helper.GetBlockWithClient(client, height+1)

	if err != nil {
		return nil, nil, "", err
	}

	// extract votes from response
	preCommits := blockDetails.LastCommit.Precommits

	// extract signs from votes
	valSigs := helper.GetSigs(preCommits)

	// extract chainID
	chainID = blockDetails.ChainID

	// return
	return preCommits, valSigs, chainID, nil
}

// IsCatchingUp checks if the heimdall node you are connected to is fully synced or not
// returns true when synced
func IsCatchingUp(cliCtx cliContext.CLIContext) bool {
	resp, err := helper.GetNodeStatus(cliCtx)
	if err != nil {
		return true
	}
	return resp.SyncInfo.CatchingUp
}

// GetConfigManagerParams return configManager params
func GetConfigManagerParams(cliCtx cliContext.CLIContext) (*chainManagerTypes.Params, error) {
	response, err := helper.FetchFromAPI(
		cliCtx,
		helper.GetHeimdallServerEndpoint(ChainManagerParamsURL),
	)

	if err != nil {
		return nil, err
	}

	var params chainManagerTypes.Params
	if err := json.Unmarshal(response.Result, &params); err != nil {
		return nil, err
	}

	return &params, nil
}
