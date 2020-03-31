package pier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
	hmtypes "github.com/maticnetwork/heimdall/types"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tendermint/libs/common"
	httpClient "github.com/tendermint/tendermint/rpc/client"

	"github.com/maticnetwork/heimdall/helper"
)

// Result represents single req result
type Result struct {
	Result uint64 `json:"result"`
}

type AckService struct {
	// Base service
	common.BaseService

	// storage client
	storageClient *leveldb.DB

	// header listener subscription
	cancelACKProcess context.CancelFunc

	// cli context
	cliCtx cliContext.CLIContext

	// queue connector
	queueConnector *QueueConnector

	// http client to subscribe to
	httpClient *httpClient.HTTP

	// contract caller
	contractConnector helper.ContractCaller
}

// NewAckService returns new service object
func NewAckService(cdc *codec.Codec, queueConnector *QueueConnector, httpClient *httpClient.HTTP) *AckService {
	// create logger
	logger := Logger.With("module", NoackService)

	cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
	cliCtx.BroadcastMode = client.BroadcastSync
	cliCtx.TrustNode = true
	contractCaller, err := helper.NewContractCaller()

	if err != nil {
		logger.Error("Error while getting contract instance", "error", err)
		panic(err)
	}
	// creating checkpointer object
	ackservice := &AckService{
		storageClient: getBridgeDBInstance(viper.GetString(BridgeDBFlag)),

		cliCtx:            cliCtx,
		queueConnector:    queueConnector,
		httpClient:        httpClient,
		contractConnector: contractCaller,
	}

	ackservice.BaseService = *common.NewBaseService(logger, NoackService, ackservice)
	return ackservice
}

// OnStart starts new block subscription
func (ackService *AckService) OnStart() error {
	ackService.BaseService.OnStart() // Always call the overridden method.

	// create cancellable context
	ackCtx, cancelACKProcess := context.WithCancel(context.Background())
	ackService.cancelACKProcess = cancelACKProcess
	// start polling for checkpoint in buffer
	go ackService.startPollingCheckpoint(ackCtx, helper.GetConfig().NoACKPollInterval)

	// subscribed to new head
	ackService.Logger.Debug("Started ACK service")

	return nil
}

// OnStop stops all necessary go routines
func (ackService *AckService) OnStop() {
	ackService.BaseService.OnStop() // Always call the overridden method.

	// cancel ack process
	ackService.cancelACKProcess()
	// close bridge db instance
	closeBridgeDBInstance()
}

func (ackService *AckService) startPollingCheckpoint(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	// stop ticker when everything done
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			go ackService.checkForCheckpoint()
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (ackService *AckService) checkForCheckpoint() {
	configParams, _ := GetConfigManagerParams(ackService.cliCtx)

	rootChainInstance, err := ackService.contractConnector.GetRootChainInstance(configParams.ChainParams.RootChainAddress.EthAddress())
	if err != nil {
		return
	}

	lastHeaderNumber, err := ackService.contractConnector.CurrentHeaderBlock(rootChainInstance)
	if err != nil {
		ackService.Logger.Error("Error while fetching current header block number", "error", err)
		return
	}

	if lastHeaderNumber == 0 {
		// First checkpoint required
		return
	}
	// get big int header number

	// header block
	_, _, _, createdAt, _, err := ackService.contractConnector.GetHeaderInfo(lastHeaderNumber, rootChainInstance)
	if err != nil {
		ackService.Logger.Error("Error while fetching header block object", "error", err)
		return
	}

	// process checkpoint
	go ackService.processCheckpoint(int64(createdAt))
}

func (ackService *AckService) processCheckpoint(lastCreatedAt int64) {
	var index float64
	// if last created at ==0 , no checkpoint yet
	if lastCreatedAt == 0 {
		index = 1
	}

	checkpointCreationTime := time.Unix(lastCreatedAt, 0)
	currentTime := time.Now().UTC()
	timeDiff := currentTime.Sub(checkpointCreationTime)
	// check if last checkpoint was < NoACK wait time
	if timeDiff.Seconds() >= helper.GetConfig().NoACKWaitTime.Seconds() && index == 0 {
		index = math.Floor(timeDiff.Seconds() / helper.GetConfig().NoACKWaitTime.Seconds())
	}

	if index == 0 {
		return
	}

	params, err := ackService.getCheckpointParams()
	if err != nil {
		return
	}

	// check if difference between no-ack time and current time
	lastNoAck := ackService.getLastNoAckTime()

	lastNoAckTime := time.Unix(int64(lastNoAck), 0)
	timeDiff = currentTime.Sub(lastNoAckTime)
	// if last no ack == 0 , first no-ack to be sent
	if currentTime.Sub(lastNoAckTime).Seconds() < params.CheckpointBufferTime.Seconds() && lastNoAck != 0 {
		ackService.Logger.Debug("Cannot send multiple no-ack in short time", "timeDiff", currentTime.Sub(lastNoAckTime).Seconds(), "ExpectedDiff", params.CheckpointBufferTime.Seconds())
		return
	}

	ackService.Logger.Debug("Fetching next proposers", "count", index)

	// check if same checkpoint still exists
	if ackService.isValidProposer(uint64(index), helper.GetAddress()) {
		ackService.Logger.Debug(
			"⛑ Sending NO ACK message",
			"currentTime", currentTime.String(),
			"proposerCount", index,
		)

		// send NO ACK
		msg := checkpointTypes.NewMsgCheckpointNoAck(
			hmtypes.BytesToHeimdallAddress(helper.GetAddress()),
		)

		// send
		err := ackService.queueConnector.BroadcastToHeimdall(msg)
		if err != nil {
			ackService.Logger.Error("Error while sending no-ack tx to Heimdall queue", "error", err)
			return
		}

		ackService.Logger.Info("No-ack transaction sent successfully", "index", index)
	}
}

func (ackService *AckService) getLastNoAckTime() uint64 {
	response, err := helper.FetchFromAPI(ackService.cliCtx, helper.GetHeimdallServerEndpoint(LastNoAckURL))
	if err != nil {
		ackService.Logger.Error("Unable to send request for checkpoint buffer", "Error", err)
		return 0
	}

	var noackObject Result
	if err := json.Unmarshal(response.Result, &noackObject); err != nil {
		ackService.Logger.Error("Error unmarshalling no-ack data ", "error", err)
		return 0
	}

	return noackObject.Result
}

func (ackService *AckService) isValidProposer(count uint64, address []byte) bool {
	ackService.Logger.Debug("Skipping proposers", "count", strconv.FormatUint(count, 10))
	response, err := helper.FetchFromAPI(
		ackService.cliCtx,
		helper.GetHeimdallServerEndpoint(fmt.Sprintf(ProposersURL, strconv.FormatUint(count, 10))),
	)
	if err != nil {
		ackService.Logger.Error("Unable to send request for next proposers", "Error", err)
		return false
	}

	// unmarshall data from buffer
	var proposers []hmtypes.Validator
	if err := json.Unmarshal(response.Result, &proposers); err != nil {
		ackService.Logger.Error("Error unmarshalling validator data ", "error", err)
		return false
	}

	ackService.Logger.Debug("Fetched proposers list", "numberOfProposers", count)
	for _, proposer := range proposers {
		if bytes.Equal(proposer.Signer.Bytes(), address) {
			return true
		}
	}

	return false
}

func (ackService *AckService) getCheckpointParams() (*checkpointTypes.Params, error) {
	response, err := helper.FetchFromAPI(
		ackService.cliCtx,
		helper.GetHeimdallServerEndpoint(CheckpointParamsURL),
	)

	if err != nil {
		return nil, err
	}

	var params checkpointTypes.Params
	if err := json.Unmarshal(response.Result, &params); err != nil {
		ackService.Logger.Error("Error unmarshalling checkpoint params", "error", err)
		return nil, err
	}

	return &params, nil
}
