package util

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	mLog "github.com/RichardKnop/machinery/v1/log"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	httpClient "github.com/tendermint/tendermint/rpc/client"
	tmTypes "github.com/tendermint/tendermint/types"

	authTypes "github.com/maticnetwork/heimdall/auth/types"
	chainManagerTypes "github.com/maticnetwork/heimdall/chainmanager/types"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	milestoneTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	clerktypes "github.com/maticnetwork/heimdall/clerk/types"
	"github.com/maticnetwork/heimdall/contracts/statesender"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
	hmtypes "github.com/maticnetwork/heimdall/types"
)

type BridgeEvent string

const (
	AccountDetailsURL       = "/auth/accounts/%v"
	LastNoAckURL            = "/checkpoints/last-no-ack"
	CheckpointParamsURL     = "/checkpoints/params"
	MilestoneParamsURL      = "/milestone/params"
	MilestoneCountURL       = "/milestone/count"
	ChainManagerParamsURL   = "/chainmanager/params"
	ProposersURL            = "/staking/proposer/%v"
	MilestoneProposersURL   = "/staking/milestoneProposer/%v"
	BufferedCheckpointURL   = "/checkpoints/buffer"
	LatestCheckpointURL     = "/checkpoints/latest"
	LatestMilestoneURL      = "/milestone/latest"
	CountCheckpointURL      = "/checkpoints/count"
	CurrentProposerURL      = "/staking/current-proposer"
	LatestSpanURL           = "/bor/latest-span"
	NextSpanInfoURL         = "/bor/prepare-next-span"
	NextSpanSeedURL         = "/bor/next-span-seed"
	DividendAccountRootURL  = "/topup/dividend-account-root"
	ValidatorURL            = "/staking/validator/%v"
	CurrentValidatorSetURL  = "staking/validator-set"
	StakingTxStatusURL      = "/staking/isoldtx"
	TopupTxStatusURL        = "/topup/isoldtx"
	ClerkTxStatusURL        = "/clerk/isoldtx"
	ClerkEventRecordURL     = "/clerk/event-record/%d"
	LatestSlashInfoBytesURL = "/slashing/latest_slash_info_bytes"
	TickSlashInfoListURL    = "/slashing/tick_slash_infos"
	SlashingTxStatusURL     = "/slashing/isoldtx"
	SlashingTickCountURL    = "/slashing/tick-count"

	TendermintUnconfirmedTxsURL      = "/unconfirmed_txs"
	TendermintUnconfirmedTxsCountURL = "/num_unconfirmed_txs"

	TransactionTimeout      = 1 * time.Minute
	CommitTimeout           = 2 * time.Minute
	TaskDelayBetweenEachVal = 10 * time.Second
	RetryTaskDelay          = 12 * time.Second
	RetryStateSyncTaskDelay = 24 * time.Second

	mempoolTxnCountDivisor = 1000

	// Bridge event types
	StakingEvent  BridgeEvent = "staking"
	TopupEvent    BridgeEvent = "topup"
	ClerkEvent    BridgeEvent = "clerk"
	SlashingEvent BridgeEvent = "slashing"

	BridgeDBFlag = "bridge-db"
)

var logger log.Logger
var loggerOnce sync.Once

// Logger returns logger singleton instance
func Logger() log.Logger {
	loggerOnce.Do(func() {
		defaultLevel := "info"
		logsWriter := helper.GetLogsWriter(helper.GetConfig().LogsWriterFile)
		logger = log.NewTMLogger(log.NewSyncWriter(logsWriter))
		option, err := log.AllowLevel(viper.GetString("log_level"))
		if err != nil {
			// cosmos sdk is using different style of log format
			// and levels don't map well, config.toml
			// see: https://github.com/cosmos/cosmos-sdk/pull/8072
			logger.Error("Unable to parse logging level", "Error", err)
			logger.Info("Using default log level")
			option, err = log.AllowLevel(defaultLevel)
			if err != nil {
				logger.Error("failed to allow default log level", "Level", defaultLevel, "Error", err)
			}
		}

		logger = log.NewFilter(logger, option)

		// set no-op logger if log level is not debug for machinery
		if viper.GetString("log_level") != "debug" {
			mLog.SetDebug(NoopLogger{})
		}
	})

	return logger
}

// IsProposer  checks if we are proposer
func IsProposer(cliCtx cliContext.CLIContext) (bool, error) {
	var (
		proposers []hmtypes.Validator
		count     = uint64(1)
	)

	result, err := helper.FetchFromAPI(cliCtx,
		helper.GetHeimdallServerEndpoint(fmt.Sprintf(ProposersURL, strconv.FormatUint(count, 10))),
	)
	if err != nil {
		logger.Error("Error fetching proposers", "url", ProposersURL, "error", err)
		return false, err
	}

	err = jsoniter.ConfigFastest.Unmarshal(result.Result, &proposers)
	if err != nil {
		logger.Error("error unmarshalling proposer slice", "error", err)
		return false, err
	}

	if bytes.Equal(proposers[0].Signer.Bytes(), helper.GetAddress()) {
		return true, nil
	}

	return false, nil
}

func IsMilestoneProposer(cliCtx cliContext.CLIContext) (bool, error) {
	var (
		proposers []hmtypes.Validator
		count     = uint64(1)
	)

	result, err := helper.FetchFromAPI(cliCtx,
		helper.GetHeimdallServerEndpoint(fmt.Sprintf(MilestoneProposersURL, strconv.FormatUint(count, 10))),
	)
	if err != nil {
		logger.Error("Error fetching milestone proposers", "url", MilestoneProposersURL, "error", err)
		return false, err
	}

	err = jsoniter.ConfigFastest.Unmarshal(result.Result, &proposers)
	if err != nil {
		logger.Error("error unmarshalling milestone proposer slice", "error", err)
		return false, err
	}

	if len(proposers) == 0 {
		logger.Error("length of proposer list is 0")
		return false, errors.Errorf("Length of proposer list is 0")
	}

	if bytes.Equal(proposers[0].Signer.Bytes(), helper.GetAddress()) {
		return true, nil
	}

	return false, nil
}

// IsInProposerList checks if we are in current proposer
func IsInProposerList(cliCtx cliContext.CLIContext, count uint64) (bool, error) {
	logger.Debug("Skipping proposers", "count", strconv.FormatUint(count+1, 10))

	response, err := helper.FetchFromAPI(
		cliCtx,
		helper.GetHeimdallServerEndpoint(fmt.Sprintf(ProposersURL, strconv.FormatUint(count+1, 10))),
	)
	if err != nil {
		logger.Error("Unable to send request for next proposers", "url", ProposersURL, "error", err)
		return false, err
	}

	// unmarshall data from buffer
	var proposers []hmtypes.Validator
	if err := jsoniter.ConfigFastest.Unmarshal(response.Result, &proposers); err != nil {
		logger.Error("Error unmarshalling validator data ", "error", err)
		return false, err
	}

	logger.Debug("Fetched proposers list", "numberOfProposers", count+1)

	for i := 1; i <= int(count) && i < len(proposers); i++ {
		if bytes.Equal(proposers[i].Signer.Bytes(), helper.GetAddress()) {
			return true, nil
		}
	}

	return false, nil
}

// IsInProposerList checks if we are in current proposer
func IsInMilestoneProposerList(cliCtx cliContext.CLIContext, count uint64) (bool, error) {
	logger.Debug("Skipping proposers", "count", strconv.FormatUint(count, 10))

	response, err := helper.FetchFromAPI(
		cliCtx,
		helper.GetHeimdallServerEndpoint(fmt.Sprintf(MilestoneProposersURL, strconv.FormatUint(count, 10))),
	)
	if err != nil {
		logger.Error("Unable to send request for next proposers", "url", MilestoneProposersURL, "error", err)
		return false, err
	}

	// unmarshall data from buffer
	var proposers []hmtypes.Validator
	if err := jsoniter.ConfigFastest.Unmarshal(response.Result, &proposers); err != nil {
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
func CalculateTaskDelay(cliCtx cliContext.CLIContext, event interface{}) (bool, time.Duration) {
	defer LogElapsedTimeForStateSyncedEvent(event, "CalculateTaskDelay", time.Now())
	// calculate validator position
	valPosition := 0
	isCurrentValidator := false

	validatorSet, err := GetValidatorSet(cliCtx)
	if err != nil {
		logger.Error("Error getting current validatorset data ", "error", err)
		return isCurrentValidator, 0
	}

	logger.Info("Fetched current validatorset list", "currentValidatorcount", len(validatorSet.Validators))

	for i, validator := range validatorSet.Validators {
		if bytes.Equal(validator.Signer.Bytes(), helper.GetAddress()) {
			valPosition = i + 1
			isCurrentValidator = true

			break
		}
	}

	// Change calculation later as per the discussion
	// Currently it will multiply delay for every 1000 unconfirmed txns in mempool
	// For example if the current default delay is 12 Seconds
	// Then for upto 1000 txns it will stay as 12 only
	// For 1000-2000 It will be 24 seconds
	// For 2000-3000 it will be 36 seconds
	// Basically for every 1000 txns it will increase the factor by 1.

	mempoolFactor := GetUnconfirmedTxnCount(event) / mempoolTxnCountDivisor

	// calculate delay
	taskDelay := time.Duration(valPosition) * TaskDelayBetweenEachVal * time.Duration(mempoolFactor+1)

	return isCurrentValidator, taskDelay
}

// IsCurrentProposer checks if we are current proposer
func IsCurrentProposer(cliCtx cliContext.CLIContext) (bool, error) {
	var proposer hmtypes.Validator

	result, err := helper.FetchFromAPI(cliCtx, helper.GetHeimdallServerEndpoint(CurrentProposerURL))
	if err != nil {
		logger.Error("Error fetching proposers", "error", err)
		return false, err
	}

	if err = jsoniter.ConfigFastest.Unmarshal(result.Result, &proposer); err != nil {
		logger.Error("error unmarshalling validator", "error", err)
		return false, err
	}

	logger.Debug("Current proposer fetched", "validator", proposer.String())

	if bytes.Equal(proposer.Signer.Bytes(), helper.GetAddress()) {
		return true, nil
	}

	logger.Debug("We are not the current proposer")

	return false, nil
}

// IsEventSender check if we are the EventSender
func IsEventSender(cliCtx cliContext.CLIContext, validatorID uint64) bool {
	var validator hmtypes.Validator

	result, err := helper.FetchFromAPI(cliCtx,
		helper.GetHeimdallServerEndpoint(fmt.Sprintf(ValidatorURL, strconv.FormatUint(validatorID, 10))),
	)
	if err != nil {
		logger.Error("Error fetching proposers", "error", err)
		return false
	}

	if err = jsoniter.ConfigFastest.Unmarshal(result.Result, &validator); err != nil {
		logger.Error("error unmarshalling proposer slice", "error", err)
		return false
	}

	logger.Debug("Current event sender received", "validator", validator.String())

	return bytes.Equal(validator.Signer.Bytes(), helper.GetAddress())
}

// CreateURLWithQuery receives the uri and parameters in key value form
// it will return the new url with the given query from the parameter
func CreateURLWithQuery(uri string, param map[string]interface{}) (string, error) {
	urlObj, err := url.Parse(uri)
	if err != nil {
		return uri, err
	}

	query := urlObj.Query()
	for k, v := range param {
		query.Set(k, fmt.Sprintf("%v", v))
	}

	urlObj.RawQuery = query.Encode()

	return urlObj.String(), nil
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
	defer func() {
		if err := client.UnsubscribeAll(ctx, subscriber); err != nil {
			logger.Error("WaitForOneEvent | UnsubscribeAll", "Error", err)
		}
	}()

	select {
	case event := <-eventCh:
		return event.Data, nil
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

// GetAccount returns heimdall auth account
func GetAccount(cliCtx cliContext.CLIContext, address types.HeimdallAddress) (account authTypes.Account, err error) {
	url := helper.GetHeimdallServerEndpoint(fmt.Sprintf(AccountDetailsURL, address))

	// call account rest api
	response, err := helper.FetchFromAPI(cliCtx, url)
	if err != nil {
		return
	}

	if err = cliCtx.Codec.UnmarshalJSON(response.Result, &account); err != nil {
		logger.Error("Error unmarshalling account details", "url", url)
		return
	}

	return
}

// GetChainmanagerParams return chain manager params
func GetChainmanagerParams(cliCtx cliContext.CLIContext) (*chainManagerTypes.Params, error) {
	response, err := helper.FetchFromAPI(
		cliCtx,
		helper.GetHeimdallServerEndpoint(ChainManagerParamsURL),
	)
	if err != nil {
		logger.Error("Error fetching chainmanager params", "err", err)
		return nil, err
	}

	var params chainManagerTypes.Params
	if err = jsoniter.ConfigFastest.Unmarshal(response.Result, &params); err != nil {
		logger.Error("Error unmarshalling chainmanager params", "url", ChainManagerParamsURL, "err", err)
		return nil, err
	}

	return &params, nil
}

// GetCheckpointParams return params
func GetCheckpointParams(cliCtx cliContext.CLIContext) (*checkpointTypes.Params, error) {
	response, err := helper.FetchFromAPI(
		cliCtx,
		helper.GetHeimdallServerEndpoint(CheckpointParamsURL),
	)

	if err != nil {
		logger.Error("Error fetching Checkpoint params", "err", err)
		return nil, err
	}

	var params checkpointTypes.Params
	if err := jsoniter.ConfigFastest.Unmarshal(response.Result, &params); err != nil {
		logger.Error("Error unmarshalling Checkpoint params", "url", CheckpointParamsURL)
		return nil, err
	}

	return &params, nil
}

// GetCheckpointParams return params
func GetMilestoneParams(cliCtx cliContext.CLIContext) (*milestoneTypes.Params, error) {
	response, err := helper.FetchFromAPI(
		cliCtx,
		helper.GetHeimdallServerEndpoint(MilestoneParamsURL),
	)

	if err != nil {
		logger.Error("Error fetching Milestone params", "err", err)
		return nil, err
	}

	var params milestoneTypes.Params
	if err := json.Unmarshal(response.Result, &params); err != nil {
		logger.Error("Error unmarshalling Checkpoint params", "url", MilestoneParamsURL)
		return nil, err
	}

	return &params, nil
}

// GetBufferedCheckpoint return checkpoint from bueffer
func GetBufferedCheckpoint(cliCtx cliContext.CLIContext) (*hmtypes.Checkpoint, error) {
	response, err := helper.FetchFromAPI(
		cliCtx,
		helper.GetHeimdallServerEndpoint(BufferedCheckpointURL),
	)

	if err != nil {
		logger.Debug("Error fetching buffered checkpoint", "err", err)
		return nil, err
	}

	var checkpoint hmtypes.Checkpoint
	if err := jsoniter.ConfigFastest.Unmarshal(response.Result, &checkpoint); err != nil {
		logger.Error("Error unmarshalling buffered checkpoint", "url", BufferedCheckpointURL, "err", err)
		return nil, err
	}

	return &checkpoint, nil
}

// GetLatestCheckpoint return last successful checkpoint
func GetLatestCheckpoint(cliCtx cliContext.CLIContext) (*hmtypes.Checkpoint, error) {
	response, err := helper.FetchFromAPI(
		cliCtx,
		helper.GetHeimdallServerEndpoint(LatestCheckpointURL),
	)

	if err != nil {
		logger.Debug("Error fetching latest checkpoint", "err", err)
		return nil, err
	}

	var checkpoint hmtypes.Checkpoint
	if err = jsoniter.ConfigFastest.Unmarshal(response.Result, &checkpoint); err != nil {
		logger.Error("Error unmarshalling latest checkpoint", "url", LatestCheckpointURL, "err", err)
		return nil, err
	}

	return &checkpoint, nil
}

// GetLatestMilestone return last successful milestone
func GetLatestMilestone(cliCtx cliContext.CLIContext) (*hmtypes.Milestone, error) {
	response, err := helper.FetchFromAPI(
		cliCtx,
		helper.GetHeimdallServerEndpoint(LatestMilestoneURL),
	)

	if err != nil {
		logger.Debug("Error fetching latest milestone", "err", err)
		return nil, err
	}

	var milestone hmtypes.Milestone
	if err = json.Unmarshal(response.Result, &milestone); err != nil {
		logger.Error("Error unmarshalling latest milestone", "url", LatestMilestoneURL, "err", err)
		return nil, err
	}

	return &milestone, nil
}

// GetCheckpointParams return params
func GetMilestoneCount(cliCtx cliContext.CLIContext) (*milestoneTypes.Count, error) {
	response, err := helper.FetchFromAPI(
		cliCtx,
		helper.GetHeimdallServerEndpoint(MilestoneCountURL),
	)

	if err != nil {
		logger.Error("Error fetching Milestone count", "err", err)
		return nil, err
	}

	var count milestoneTypes.Count
	if err := json.Unmarshal(response.Result, &count); err != nil {
		logger.Error("Error unmarshalling milestone Count", "url", MilestoneCountURL)
		return nil, err
	}

	return &count, nil
}

// AppendPrefix returns publickey in uncompressed format
func AppendPrefix(signerPubKey []byte) []byte {
	// append prefix - "0x04" as heimdall uses publickey in uncompressed format. Refer below link
	// https://superuser.com/questions/1465455/what-is-the-size-of-public-key-for-ecdsa-spec256r1
	prefix := make([]byte, 1)
	prefix[0] = byte(0x04)
	signerPubKey = append(prefix[:], signerPubKey[:]...)

	return signerPubKey
}

// GetValidatorNonce fetches validator nonce and height
func GetValidatorNonce(cliCtx cliContext.CLIContext, validatorID uint64) (uint64, int64, error) {
	var validator hmtypes.Validator

	result, err := helper.FetchFromAPI(cliCtx,
		helper.GetHeimdallServerEndpoint(fmt.Sprintf(ValidatorURL, strconv.FormatUint(validatorID, 10))),
	)

	if err != nil {
		logger.Error("Error fetching validator data", "error", err)
		return 0, 0, err
	}

	if err = jsoniter.ConfigFastest.Unmarshal(result.Result, &validator); err != nil {
		logger.Error("error unmarshalling validator data", "error", err)
		return 0, 0, err
	}

	logger.Debug("Validator data received ", "validator", validator.String())

	return validator.Nonce, result.Height, nil
}

// GetValidatorSet fetches the current validator set
func GetValidatorSet(cliCtx cliContext.CLIContext) (*hmtypes.ValidatorSet, error) {
	response, err := helper.FetchFromAPI(cliCtx, helper.GetHeimdallServerEndpoint(CurrentValidatorSetURL))
	if err != nil {
		logger.Error("Unable to send request for current validatorset", "url", CurrentValidatorSetURL, "error", err)
		return nil, err
	}

	var validatorSet hmtypes.ValidatorSet
	if err = jsoniter.ConfigFastest.Unmarshal(response.Result, &validatorSet); err != nil {
		logger.Error("Error unmarshalling current validatorset data ", "error", err)
		return nil, err
	}

	return &validatorSet, nil
}

// GetBlockHeight return last successful checkpoint
func GetBlockHeight(cliCtx cliContext.CLIContext) int64 {
	response, err := helper.FetchFromAPI(
		cliCtx,
		helper.GetHeimdallServerEndpoint(CountCheckpointURL),
	)
	if err != nil {
		logger.Debug("Error fetching latest block height", "err", err)
		return 0
	}

	return response.Height
}

// GetClerkEventRecord return last successful checkpoint
func GetClerkEventRecord(cliCtx cliContext.CLIContext, stateId int64) (*clerktypes.EventRecord, error) {
	response, err := helper.FetchFromAPI(
		cliCtx,
		helper.GetHeimdallServerEndpoint(fmt.Sprintf(ClerkEventRecordURL, stateId)),
	)
	if err != nil {
		logger.Error("Error fetching event record by state ID", "error", err)
		return nil, err
	}

	var eventRecord clerktypes.EventRecord
	if err = jsoniter.ConfigFastest.Unmarshal(response.Result, &eventRecord); err != nil {
		logger.Error("Error unmarshalling event record", "error", err)
		return nil, err
	}

	return &eventRecord, nil
}

func GetUnconfirmedTxnCount(event interface{}) int {
	defer LogElapsedTimeForStateSyncedEvent(event, "GetUnconfirmedTxnCount", time.Now())

	endpoint := helper.GetConfig().TendermintRPCUrl + TendermintUnconfirmedTxsCountURL

	resp, err := helper.Client.Get(endpoint)
	if err != nil || resp.StatusCode != http.StatusOK {
		logger.Error("Error fetching mempool txs count", "url", endpoint, "error", err)
		return 0
	}

	body, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		logger.Error("Error fetching mempool txs count", "error", err)
		return 0
	}

	// a minimal response of the unconfirmed txs
	var response TendermintUnconfirmedTxs

	err = jsoniter.ConfigFastest.Unmarshal(body, &response)
	if err != nil {
		logger.Error("Error unmarshalling response received from Heimdall Server", "error", err)
		return 0
	}

	count, _ := strconv.Atoi(response.Result.Total)

	return count
}

// LogElapsedTimeForStateSyncedEvent logs useful info for StateSynced events
func LogElapsedTimeForStateSyncedEvent(event interface{}, functionName string, startTime time.Time) {
	if event == nil {
		return
	}

	var (
		typedEvent  statesender.StatesenderStateSynced
		timeElapsed = time.Since(startTime).Milliseconds()
	)

	switch e := event.(type) {
	case statesender.StatesenderStateSynced:
		typedEvent = e
	case *statesender.StatesenderStateSynced:
		if e == nil {
			return
		}

		typedEvent = *e
	default:
		return
	}

	logger.Info("StateSyncedEvent: "+functionName,
		"stateSyncId", typedEvent.Id,
		"timeElapsed", timeElapsed)
}

// IsPubKeyFirstByteValid checks the validity of the first byte of the public key.
// It must be 0x04 for uncompressed public keys
func IsPubKeyFirstByteValid(pubKey []byte) bool {
	prefix := make([]byte, 1)
	prefix[0] = byte(0x04)

	return bytes.Equal(prefix, pubKey[0:1])
}
