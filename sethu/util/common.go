package util

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"time"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/log"
	httpClient "github.com/tendermint/tendermint/rpc/client"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/maticnetwork/heimdall/helper"
	hmtypes "github.com/maticnetwork/heimdall/types"
	rest "github.com/maticnetwork/heimdall/types/rest"
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
func IsProposer(cliCtx cliContext.CLIContext) (bool, error) {
	var proposers []hmtypes.Validator
	count := uint64(1)
	result, err := FetchFromAPI(cliCtx,
		GetHeimdallServerEndpoint(fmt.Sprintf(ProposersURL, strconv.FormatUint(count, 10))),
	)
	if err != nil {
		Logger.Error("Error fetching proposers", "url", ProposersURL, "error", err)
		return false, err
	}
	err = json.Unmarshal(result.Result, &proposers)
	if err != nil {
		Logger.Error("error unmarshalling proposer slice", "error", err)
		return false, err
	}

	if bytes.Equal(proposers[0].Signer.Bytes(), helper.GetAddress()) {
		return true, nil
	}

	return false, nil
}

func IsInProposerList(cliCtx cliContext.CLIContext, count uint64) (bool, error) {
	Logger.Debug("Skipping proposers", "count", strconv.FormatUint(count, 10))
	response, err := FetchFromAPI(
		cliCtx,
		GetHeimdallServerEndpoint(fmt.Sprintf(ProposersURL, strconv.FormatUint(count, 10))),
	)
	if err != nil {
		Logger.Error("Unable to send request for next proposers", "url", ProposersURL, "error", err)
		return false, err
	}

	// unmarshall data from buffer
	var proposers []hmtypes.Validator
	if err := json.Unmarshal(response.Result, &proposers); err != nil {
		Logger.Error("Error unmarshalling validator data ", "error", err)
		return false, err
	}

	Logger.Debug("Fetched proposers list", "numberOfProposers", count)
	for _, proposer := range proposers {
		if bytes.Equal(proposer.Signer.Bytes(), helper.GetAddress()) {
			return true, nil
		}
	}
	return false, nil
}

// checks if we are current proposer
func IsCurrentProposer(cliCtx cliContext.CLIContext) (bool, error) {
	var proposer hmtypes.Validator
	result, err := FetchFromAPI(cliCtx, GetHeimdallServerEndpoint(CurrentProposerURL))
	if err != nil {
		Logger.Error("Error fetching proposers", "error", err)
		return false, err
	}
	err = json.Unmarshal(result.Result, &proposer)
	if err != nil {
		Logger.Error("error unmarshalling validator", "error", err)
		return false, err
	}
	Logger.Debug("Current proposer fetched", "validator", proposer.String())

	if bytes.Equal(proposer.Signer.Bytes(), helper.GetAddress()) {
		return true, nil
	}

	return false, nil
}

// check if we are the EventSender
func IsEventSender(cliCtx cliContext.CLIContext, validatorID uint64) bool {

	var validator hmtypes.Validator

	result, err := FetchFromAPI(cliCtx,
		GetHeimdallServerEndpoint(fmt.Sprintf(ValidatorURL, strconv.FormatUint(validatorID, 10))),
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

// GetHeimdallServerEndpoint returns heimdall server endpoint
func GetHeimdallServerEndpoint(endpoint string) string {
	u, _ := url.Parse(helper.GetConfig().HeimdallServerURL)
	u.Path = path.Join(u.Path, endpoint)
	return u.String()
}

// FetchFromAPI fetches data from any URL
func FetchFromAPI(cliCtx cliContext.CLIContext, URL string) (result rest.ResponseWithHeight, err error) {
	resp, err := http.Get(URL)
	if err != nil {
		return result, err
	}
	defer resp.Body.Close()

	// response
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return result, err
		}
		// unmarshall data from buffer
		var response rest.ResponseWithHeight
		if err := cliCtx.Codec.UnmarshalJSON(body, &response); err != nil {
			return result, err
		}
		return response, nil
	}

	Logger.Debug("Error while fetching data from URL", "status", resp.StatusCode, "URL", URL)
	return result, fmt.Errorf("Error while fetching data from url: %v, status: %v", URL, resp.StatusCode)
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
