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
	"sync"
	"time"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
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
	CurrentValidatorSetURL = "staking/validator-set"

	TransactionTimeout = 1 * time.Minute
	CommitTimeout      = 2 * time.Minute

	BridgeDBFlag = "bridge-db"
)

var logger log.Logger
var loggerOnce sync.Once

// Logger returns logger singleton instance
func Logger() log.Logger {
	loggerOnce.Do(func() {
		logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
		option, _ := log.AllowLevel(viper.GetString("log_level"))
		logger = log.NewFilter(logger, option)
	})

	return logger
}

// IsProposer  checks if we are proposer
func IsProposer(cliCtx cliContext.CLIContext) (bool, error) {
	var proposers []hmtypes.Validator
	count := uint64(1)
	result, err := FetchFromAPI(cliCtx,
		GetHeimdallServerEndpoint(fmt.Sprintf(ProposersURL, strconv.FormatUint(count, 10))),
	)

	if err != nil {
		logger.Error("Error fetching proposers", "url", ProposersURL, "error", err)
		return false, err
	}

	err = json.Unmarshal(result.Result, &proposers)
	if err != nil {
		logger.Error("error unmarshalling proposer slice", "error", err)
		return false, err
	}

	if bytes.Equal(proposers[0].Signer.Bytes(), helper.GetAddress()) {
		return true, nil
	}

	return false, nil
}

// IsInProposerList checks if we are in current proposer
func IsInProposerList(cliCtx cliContext.CLIContext, count uint64) (bool, error) {
	logger.Debug("Skipping proposers", "count", strconv.FormatUint(count, 10))
	response, err := FetchFromAPI(
		cliCtx,
		GetHeimdallServerEndpoint(fmt.Sprintf(ProposersURL, strconv.FormatUint(count, 10))),
	)
	if err != nil {
		logger.Error("Unable to send request for next proposers", "url", ProposersURL, "error", err)
		return false, err
	}

	// unmarshall data from buffer
	var proposers []hmtypes.Validator
	if err := json.Unmarshal(response.Result, &proposers); err != nil {
		logger.Error("Error unmarshalling validator data ", "error", err)
		return false, err
	}

	logger.Debug("Fetched proposers list", "numberOfProposers", count)
	for _, proposer := range proposers {
		if bytes.Equal(proposer.Signer.Bytes(), helper.GetAddress()) {
			return true, nil
		}
	}
	return false, nil
}

// CalculateTaskDelay calculates delay required for current validator to propose the tx
// It solves for multiple validators sending same transaction.
func CalculateTaskDelay(cliCtx cliContext.CLIContext) (bool, time.Duration) {
	// calculate validator position
	valPosition := 0
	isCurrentValidator := false
	response, err := FetchFromAPI(cliCtx, GetHeimdallServerEndpoint(CurrentValidatorSetURL))
	if err != nil {
		logger.Error("Unable to send request for current validatorset", "url", CurrentValidatorSetURL, "error", err)
		return isCurrentValidator, 0
	}
	// unmarshall data from buffer
	var currentValidators []hmtypes.Validator
	if err := json.Unmarshal(response.Result, &currentValidators); err != nil {
		logger.Error("Error unmarshalling current validatorset data ", "error", err)
		return isCurrentValidator, 0
	}
	logger.Debug("Fetched current validatorset list", "currentValidatorcount", len(currentValidators))
	for i, validator := range currentValidators {
		if bytes.Equal(validator.Signer.Bytes(), helper.GetAddress()) {
			valPosition = i
			isCurrentValidator = true
			break
		}
	}

	// calculate delay
	delayBetweenEachVal := 3 * time.Second
	taskDelay := time.Duration(valPosition) * delayBetweenEachVal
	return isCurrentValidator, taskDelay
}

// IsCurrentProposer checks if we are current proposer
func IsCurrentProposer(cliCtx cliContext.CLIContext) (bool, error) {
	var proposer hmtypes.Validator
	result, err := FetchFromAPI(cliCtx, GetHeimdallServerEndpoint(CurrentProposerURL))
	if err != nil {
		logger.Error("Error fetching proposers", "error", err)
		return false, err
	}

	err = json.Unmarshal(result.Result, &proposer)
	if err != nil {
		logger.Error("error unmarshalling validator", "error", err)
		return false, err
	}
	logger.Debug("Current proposer fetched", "validator", proposer.String())

	if bytes.Equal(proposer.Signer.Bytes(), helper.GetAddress()) {
		return true, nil
	}

	return false, nil
}

// IsEventSender check if we are the EventSender
func IsEventSender(cliCtx cliContext.CLIContext, validatorID uint64) bool {
	var validator hmtypes.Validator

	result, err := FetchFromAPI(cliCtx,
		GetHeimdallServerEndpoint(fmt.Sprintf(ValidatorURL, strconv.FormatUint(validatorID, 10))),
	)
	if err != nil {
		logger.Error("Error fetching proposers", "error", err)
		return false
	}

	err = json.Unmarshal(result.Result, &validator)
	if err != nil {
		logger.Error("error unmarshalling proposer slice", "error", err)
		return false
	}
	logger.Debug("Current event sender received", "validator", validator.String())

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

	logger.Debug("Error while fetching data from URL", "status", resp.StatusCode, "URL", URL)
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

// IsCatchingUp checks if the heimdall node you are connected to is fully synced or not
// returns true when synced
func IsCatchingUp(cliCtx cliContext.CLIContext) bool {
	resp, err := helper.GetNodeStatus(cliCtx)
	if err != nil {
		return true
	}
	return resp.SyncInfo.CatchingUp
}
