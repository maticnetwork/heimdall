package util

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"sync"
	"time"

	borTypes "github.com/maticnetwork/heimdall/x/bor/types"

	chainmanagerTypes "github.com/maticnetwork/heimdall/x/chainmanager/types"

	types2 "github.com/maticnetwork/heimdall/x/staking/types"

	"github.com/gogo/protobuf/jsonpb"

	checkpointTypes "github.com/maticnetwork/heimdall/x/checkpoint/types"

	mLog "github.com/RichardKnop/machinery/v1/log"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/libs/log"
	httpClient "github.com/tendermint/tendermint/rpc/client/http"
	tmTypes "github.com/tendermint/tendermint/types"

	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/maticnetwork/heimdall/helper"
	"github.com/maticnetwork/heimdall/types"
	hmTypes "github.com/maticnetwork/heimdall/types"
	hmCommonTypes "github.com/maticnetwork/heimdall/types/common"
)

const (
	AccountDetailsURL       = "/cosmos/auth/v1beta1/accounts/%v"
	LastNoAckURL            = "/heimdall/checkpoint/v1beta1/last-no-ack"
	CheckpointParamsURL     = "/heimdall/checkpoint/v1beta1/params"
	ChainManagerParamsURL   = "/heimdall/chainmanager/v1beta1/params"
	ProposersURL            = "/heimdall/staking/v1beta1/proposer/%v"
	BufferedCheckpointURL   = "/heimdall/checkpoint/v1beta1/buffer"
	LatestCheckpointURL     = "/heimdall/checkpoint/v1beta1/latest"
	LatestSpanURL           = "/heimdall/bor/v1beta1/latest-span"
	NextSpanInfoURL         = "/heimdall/bor/v1beta1/prepare-next-span"
	NextSpanSeedURL         = "/heimdall/bor/v1beta1/next-span-seed"
	DividendAccountRootURL  = "/heimdall/topup/v1beta1/dividend-account-root"
	ValidatorURL            = "/heimdall/staking/v1beta1/validator/%v"
	CurrentValidatorSetURL  = "/heimdall/staking/v1beta1/validator-set"
	StakingTxStatusURL      = "/heimdall/staking/v1beta1/isoldtx"
	TopupTxStatusURL        = "/heimdall/topup/v1beta1/isoldtx"
	ClerkTxStatusURL        = "/heimdall/clerk/v1beta1/isoldtx"
	LatestSlashInfoBytesURL = "/slashing/latest_slash_info_bytes"
	TickSlashInfoListURL    = "/slashing/tick_slash_infos"
	SlashingTxStatusURL     = "/slashing/isoldtx"
	SlashingTickCountURL    = "/slashing/tick-count"

	TransactionTimeout      = 1 * time.Minute
	CommitTimeout           = 2 * time.Minute
	BlockInterval           = 6 * time.Second
	TaskDelayBetweenEachVal = 3 * BlockInterval
	ValidatorJoinRetryDelay = 3 * BlockInterval

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

		// set no-op logger if log level is not debug for machinery
		if viper.GetString("log_level") != "debug" {
			mLog.SetDebug(NoopLogger{})
		}
	})

	return logger
}

// IsProposer  checks if we are proposer
func IsProposer(cliCtx client.Context) (bool, error) {
	var validatorSet types2.QueryValidatorSetResponse
	result, err := helper.FetchFromAPI(helper.GetHeimdallServerEndpoint(CurrentValidatorSetURL))

	if err != nil {
		logger.Error("Error fetching proposers", "url", CurrentValidatorSetURL, "error", err)
		return false, err
	}

	err = jsonpb.UnmarshalString(string(result), &validatorSet)
	if err != nil {
		logger.Error("error unmarshalling proposer slice", "error", err)
		return false, err
	}

	if bytes.Equal(validatorSet.ValidatorSet.Proposer.GetSigner(), helper.GetAddress()) {
		return true, nil
	}

	return false, nil
}

// IsInProposerList checks if we are in current proposer
func IsInProposerList(cliCtx client.Context, count uint64) (bool, error) {
	logger.Debug("Skipping proposers", "count", strconv.FormatUint(count, 10))
	response, err := helper.FetchFromAPI(helper.GetHeimdallServerEndpoint(fmt.Sprintf(ProposersURL, strconv.FormatUint(count, 10))))
	if err != nil {
		logger.Error("Unable to send request for next proposers", "url", ProposersURL, "error", err)
		return false, err
	}

	// unmarshall data from buffer
	var proposers types2.QueryProposerResponse

	if err := jsonpb.UnmarshalString(string(response), &proposers); err != nil {
		logger.Error("Error unmarshalling validator data ", "error", err)
		return false, err
	}

	logger.Debug("Fetched proposers list", "numberOfProposers", count)
	for _, proposer := range proposers.Proposers {
		if bytes.Equal(proposer.GetSigner(), helper.GetAddress()) {
			return true, nil
		}
	}
	return false, nil
}

// CalculateTaskDelay calculates delay required for current validator to propose the tx
// It solves for multiple validators sending same transaction.
func CalculateTaskDelay(cliCtx client.Context) (bool, time.Duration) {
	// calculate validator position
	valPosition := 0
	isCurrentValidator := false
	response, err := helper.FetchFromAPI(helper.GetHeimdallServerEndpoint(CurrentValidatorSetURL))
	if err != nil {
		logger.Error("Unable to send request for current validatorset", "url", CurrentValidatorSetURL, "error", err)
		return false, 0
	}
	// unmarshall data from buffer
	var validatorSet types2.QueryValidatorSetResponse
	err = jsonpb.UnmarshalString(string(response), &validatorSet)
	if err != nil {
		logger.Error("Error unmarshalling current validatorset data ", "error", err)
		return false, 0
	}

	logger.Info("Fetched current validatorset list", "currentValidatorcount", len(validatorSet.ValidatorSet.Validators))
	for i, validator := range validatorSet.ValidatorSet.Validators {
		if bytes.Equal(validator.GetSigner(), helper.GetAddress()) {
			valPosition = i + 1
			isCurrentValidator = true
			break
		}
	}

	// calculate delay
	taskDelay := time.Duration(valPosition) * TaskDelayBetweenEachVal
	return isCurrentValidator, taskDelay
}

func CalculateSpanTaskDelay(cliContext client.Context, id uint64, start uint64) (bool, time.Duration) {
	// calculate validator position
	valPosition := 0
	isNextSpanProducer := false
	nextSpan, err := FetchNextSpanDetails(cliContext, id, start)

	if err != nil {
		logger.Error("Error while sending request for next span details", "error", err)
		return false, 0
	}

	// check if current user is among next span producers
	// find the index of current validator in nextSpanProducers list
	for i, validator := range nextSpan.SelectedProducers {
		if bytes.Equal(validator.GetSigner(), helper.GetAddress()) {
			valPosition = i + 1
			isNextSpanProducer = true
			break
		}
	}

	// calculate delay
	taskDelay := time.Duration(valPosition) * TaskDelayBetweenEachVal
	return isNextSpanProducer, taskDelay
}

// IsCurrentProposer checks if we are current proposer
func IsCurrentProposer(cliCtx client.Context) (bool, error) {
	var validatorSet types2.QueryValidatorSetResponse
	result, err := helper.FetchFromAPI(helper.GetHeimdallServerEndpoint(CurrentValidatorSetURL))
	if err != nil {
		logger.Error("Error fetching proposers", "error", err)
		return false, err
	}

	err = jsonpb.UnmarshalString(string(result), &validatorSet)
	if err != nil {
		logger.Error("error unmarshalling validator", "error", err)
		return false, err
	}
	logger.Debug("Current proposer fetched", "validator", validatorSet.ValidatorSet.Proposer.Signer)

	if bytes.Equal(validatorSet.ValidatorSet.Proposer.GetSigner(), helper.GetAddress()) {
		return true, nil
	}

	return false, nil
}

// IsEventSender check if we are the EventSender
func IsEventSender(cliCtx client.Context, validatorID uint64) bool {
	var validator types2.QueryValidatorResponse

	result, err := helper.FetchFromAPI(helper.GetHeimdallServerEndpoint(fmt.Sprintf(ValidatorURL, strconv.FormatUint(validatorID, 10))))
	if err != nil {
		logger.Error("Error fetching proposers", "error", err)
		return false
	}

	err = jsonpb.UnmarshalString(string(result), &validator)
	if err != nil {
		logger.Error("error unmarshalling proposer slice", "error", err)
		return false
	}
	logger.Debug("Current event sender received", "validator", validator.String())

	return bytes.Equal(validator.Validator.GetSigner(), helper.GetAddress())
}

//CreateURLWithQuery receives the uri and parameters in key value form
//it will return the new url with the given query from the parameter
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
		return event.Data.(tmTypes.TMEventData), nil
	case <-ctx.Done():
		return nil, errors.New("timed out waiting for event")
	}
}

// IsCatchingUp checks if the heimdall node you are connected to is fully synced or not
// returns true when synced
func IsCatchingUp(cliCtx client.Context) bool {
	resp, err := helper.GetNodeStatus(cliCtx)
	if err != nil {
		return true
	}
	return resp.SyncInfo.CatchingUp
}

// GetAccount returns heimdall auth account
func GetAccount(cliCtx client.Context, address hmCommonTypes.HeimdallAddress) (account authTypes.BaseAccount, err error) {
	url := helper.GetHeimdallServerEndpoint(fmt.Sprintf(AccountDetailsURL, address))
	// call account rest api
	response, err := helper.FetchFromAPI(url)
	if err != nil {
		return
	}

	if err = json.Unmarshal(response, &account); err != nil {
		logger.Error("Error unmarshalling account details", "url", url, "Error", err)
		return
	}
	return
}

// GetChainmanagerParams return chain manager params
func GetChainmanagerParams(cliCtx client.Context) (*chainmanagerTypes.Params, error) {
	response, err := helper.FetchFromAPI(helper.GetHeimdallServerEndpoint(ChainManagerParamsURL))

	if err != nil {
		logger.Error("Error fetching chainmanager params", "err", err)
		return nil, err
	}

	var chainmanagerParamsResponse chainmanagerTypes.QueryParamsResponse
	if err := jsonpb.UnmarshalString(string(response), &chainmanagerParamsResponse); err != nil {
		logger.Error("Error unmarshalling chainmanager params", "url", ChainManagerParamsURL, "err", err)
		return nil, err
	}
	return chainmanagerParamsResponse.Params, nil
}

// GetCheckpointParams return params
func GetCheckpointParams(cliCtx client.Context) (*checkpointTypes.Params, error) {
	response, err := helper.FetchFromAPI(helper.GetHeimdallServerEndpoint(CheckpointParamsURL))

	if err != nil {
		logger.Error("Error fetching Checkpoint params", "err", err)
		return nil, err
	}

	var params checkpointTypes.QueryParamsResponse
	if err := jsonpb.UnmarshalString(string(response), &params); err != nil {
		logger.Error("Error unmarshalling Checkpoint params", "url", CheckpointParamsURL, "Error", err)
		return nil, err
	}

	return &params.Params, nil
}

// GetBufferedCheckpoint return checkpoint from buffer
func GetBufferedCheckpoint(cliCtx client.Context) (*hmTypes.Checkpoint, error) {
	response, err := helper.FetchFromAPI(helper.GetHeimdallServerEndpoint(BufferedCheckpointURL))

	if err != nil {
		logger.Debug("Error fetching buffered checkpoint", "err", err)
		return nil, err
	}

	var checkpoint hmTypes.Checkpoint
	if err := jsonpb.UnmarshalString(string(response), &checkpoint); err != nil {
		logger.Error("Error unmarshalling buffered checkpoint", "url", BufferedCheckpointURL, "err", err)
		return nil, err
	}

	return &checkpoint, nil
}

// GetlastestCheckpoint return last successful checkpoint
func GetlastestCheckpoint(cliCtx client.Context) (*hmTypes.Checkpoint, error) {
	response, err := helper.FetchFromAPI(helper.GetHeimdallServerEndpoint(LatestCheckpointURL))

	if err != nil {
		logger.Debug("Error fetching latest checkpoint", "err", err)
		return nil, err
	}

	var checkpoint hmTypes.Checkpoint
	if err := json.Unmarshal(response, &checkpoint); err != nil {
		logger.Error("Error unmarshalling latest checkpoint", "url", LatestCheckpointURL, "err", err)
		return nil, err
	}

	return &checkpoint, nil
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

// fetch next span details from heimdall.
func FetchNextSpanDetails(cliCtx client.Context, id uint64, start uint64) (*types.Span, error) {
	req, err := http.NewRequest("GET", helper.GetHeimdallServerEndpoint(NextSpanInfoURL), nil)
	if err != nil {
		logger.Error("Error creating a new request", "error", err)
		return nil, err
	}
	configParams, err := GetChainmanagerParams(cliCtx)
	if err != nil {
		logger.Error("Error while fetching chainmanager params", "error", err)
		return nil, err
	}

	q := req.URL.Query()
	q.Add("span_id", strconv.FormatUint(id, 10))
	q.Add("start_block", strconv.FormatUint(start, 10))
	q.Add("chain_id", configParams.ChainParams.BorChainID)
	q.Add("proposer", helper.GetFromAddress(cliCtx).String())
	req.URL.RawQuery = q.Encode()

	// fetch next span details
	result, err := helper.FetchFromAPI(req.URL.String())
	if err != nil {
		logger.Error("Error fetching proposers", "error", err)
		return nil, err
	}

	var msg borTypes.QueryPrepareNextSpanResponse
	if err = jsonpb.UnmarshalString(string(result), &msg); err != nil {
		logger.Error("Error unmarshalling propose tx msg ", "error", err)
		return nil, err
	}

	logger.Debug("◽ Generated proposer span msg", "msg", msg.String())
	return msg.Span, nil
}

// get Last span
func GetLastSpan(cliCtx client.Context) (*types.Span, error) {
	// fetch last span
	result, err := helper.FetchFromAPI(helper.GetHeimdallServerEndpoint(LatestSpanURL))
	if err != nil {
		logger.Error("Error while fetching latest span")
		return nil, err
	}
	var lastSpan borTypes.QueryLatestSpanResponse
	err = json.Unmarshal(result, &lastSpan)
	if err != nil {
		logger.Error("Error unmarshalling span", "error", err)
		return nil, err
	}
	return lastSpan.Span, nil
}
