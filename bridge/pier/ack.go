package pier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"time"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/helper"
	hmtypes "github.com/maticnetwork/heimdall/types"
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

	// Rootchain instance
	rootChainInstance *rootchain.Rootchain

	// header listener subscription
	cancelACKProcess context.CancelFunc

	cliCtx cliContext.CLIContext
}

// NewAckService returns new service object
func NewAckService() *AckService {
	// create logger
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", noackService)

	// root chain instance
	rootchainInstance, err := helper.GetRootChainInstance()
	if err != nil {
		logger.Error("Error while getting root chain instance", "error", err)
		panic(err)
	}

	cliCtx := cliContext.NewCLIContext()
	cliCtx.Async = true

	// creating checkpointer object
	ackservice := &AckService{
		storageClient:     getBridgeDBInstance(viper.GetString(bridgeDBFlag)),
		rootChainInstance: rootchainInstance,
		cliCtx:            cliCtx,
	}

	ackservice.BaseService = *common.NewBaseService(logger, noackService, ackservice)
	return ackservice
}

// OnStart starts new block subscription
func (ackService *AckService) OnStart() error {
	ackService.BaseService.OnStart() // Always call the overridden method.

	// create cancellable context
	ackCtx, cancelACKProcess := context.WithCancel(context.Background())
	ackService.cancelACKProcess = cancelACKProcess
	// start polling for checkpoint in buffer
	go ackService.startPollingCheckpoint(ackCtx, defaultCheckpointPollInterval)

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
	currentHeaderNumber, err := ackService.rootChainInstance.CurrentHeaderBlock(nil)
	if err != nil {
		ackService.Logger.Error("Error while fetching current header block number", "error", err)
		return
	}

	// fetch last header number
	lastHeaderNumber := currentHeaderNumber.Uint64() - helper.GetConfig().ChildBlockInterval
	if lastHeaderNumber == 0 {
		// First checkpoint required
		return
	}
	// get big int header number
	headerNumber := big.NewInt(0)
	headerNumber.SetUint64(lastHeaderNumber)

	// header block
	headerObject, err := ackService.rootChainInstance.HeaderBlock(nil, headerNumber)
	if err != nil {
		ackService.Logger.Error("Error while fetching header block object", "error", err)
		return
	}

	// process checkpoint
	go ackService.processCheckpoint(headerObject.CreatedAt.Int64())
}

func (ackService *AckService) processCheckpoint(lastCreatedAt int64) {
	var index float64
	// if last created at ==0 , no checkpoint yet
	if lastCreatedAt == 0 {
		index = 1
	}

	checkpointCreationTime := time.Unix(lastCreatedAt, 0)
	currentTime := time.Now()
	timeDiff := currentTime.Sub(checkpointCreationTime)
	// check if last checkpoint was < checkpointBufferTime
	if timeDiff.Seconds() >= helper.CheckpointBufferTime.Seconds() && index == 0 {
		index = math.Floor(timeDiff.Seconds() / helper.CheckpointBufferTime.Seconds())
		ackService.Logger.Info("index set", "Index", index)
	}

	if index == 0 {
		return
	}

	// check if difference between no-ack time and current time
	lastNoAck := ackService.getLastNoAckTime()
	// if last no ack == 0 , first no-ack to be sent
	if lastNoAck != 0 {
		lastNoAckTime := time.Unix(int64(ackService.getLastNoAckTime()), 0)
		timeDiff = currentTime.Sub(lastNoAckTime)
		ackService.Logger.Debug("created time diff", "TimeDiff", timeDiff, "lasttime", lastNoAckTime)
		if currentTime.Sub(lastNoAckTime).Seconds() < helper.CheckpointBufferTime.Seconds() {
			ackService.Logger.Debug("Cannot send multiple no-ack in short time", "timeDiff", currentTime.Sub(lastNoAckTime).Seconds(), "ExpectedDiff", helper.CheckpointBufferTime.Seconds())
			return
		}
	}

	ackService.Logger.Debug("Fetching next proposers", "Count", index)

	// check if same checkpoint still exists
	if ackService.isValidProposer(uint64(index), helper.GetPubKey().Address().Bytes()) {
		ackService.Logger.Debug(
			"Sending NO ACK message",
			"currentTime", currentTime.String(),
			"proposerCount", index,
		)

		// send NO ACK
		txBytes, err := helper.CreateTxBytes(
			checkpoint.NewMsgCheckpointNoAck(
				uint64(time.Now().Unix()),
			),
		)

		if err != nil {
			ackService.Logger.Error("Error while creating tx bytes", "error", err)
			return
		}

		resp, err := helper.SendTendermintRequest(ackService.cliCtx, txBytes)
		if err != nil {
			ackService.Logger.Error("Error while sending request to Tendermint", "error", err)
			return
		}

		ackService.Logger.Info("no-ack transaction sent successfully", "txHash", resp.Hash)
	}
}

func (ackService *AckService) getLastNoAckTime() uint64 {
	resp, err := http.Get(lastNoAckURL)
	if err != nil {
		ackService.Logger.Error("Unable to send request for checkpoint buffer", "Error", err)
		return 0
	}

	if resp.StatusCode == 200 {
		ackService.Logger.Info("Found last no-ack")
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			ackService.Logger.Error("Unable to parse no-ack body", "error", err)
			return 0
		}

		var noackObject Result
		if err := json.Unmarshal(data, &noackObject); err != nil {
			ackService.Logger.Error("Error unmarshalling no-ack data ", "error", err)
		} else {
			return noackObject.Result
		}
	}
	return 0
}

func (ackService *AckService) isValidProposer(count uint64, address []byte) bool {
	ackService.Logger.Debug("Skipping proposers", "count", strconv.FormatUint(count, 10))
	resp, err := http.Get(fmt.Sprintf(proposersURL, strconv.FormatUint(count, 10)))
	if err != nil {
		ackService.Logger.Error("Unable to send request for next proposers", "Error", err)
		return false
	}
	ackService.Logger.Debug("Request for proposer was successfull", "Count", count, "Status", resp.Status)
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ackService.Logger.Error("Unable to read data from response", "Error", err)
		return false
	}

	// unmarshall data from buffer
	var proposers []hmtypes.Validator
	if err := json.Unmarshal(body, &proposers); err != nil {
		ackService.Logger.Error("Error unmarshalling validator data ", "error", err)
		return false
	}

	ackService.Logger.Debug("Fetched proposers list from heimdall", "numberOfProposers", count)
	for _, proposer := range proposers {
		if bytes.Equal(proposer.Address.Bytes(), address) {
			return true
		}
	}

	return false
}
