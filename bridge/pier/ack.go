package pier

// import (
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"math"
// 	"math/big"
// 	"strconv"
// 	"time"

// 	"github.com/cosmos/cosmos-sdk/client"
// 	cliContext "github.com/cosmos/cosmos-sdk/client/context"
// 	"github.com/cosmos/cosmos-sdk/codec"
// 	"github.com/spf13/viper"
// 	"github.com/syndtr/goleveldb/leveldb"
// 	"github.com/tendermint/tendermint/libs/common"
// 	httpClient "github.com/tendermint/tendermint/rpc/client"

// 	checkpointTypes "github.com/maticnetwork/heimdall/checkpoint/types"
// 	"github.com/maticnetwork/heimdall/contracts/rootchain"
// 	"github.com/maticnetwork/heimdall/helper"
// 	hmtypes "github.com/maticnetwork/heimdall/types"
// )

// // Result represents single req result
// type Result struct {
// 	Result uint64 `json:"result"`
// }

// type AckService struct {
// 	// Base service
// 	common.BaseService

// 	// storage client
// 	storageClient *leveldb.DB

// 	// Rootchain instance
// 	rootChainInstance *rootchain.Rootchain

// 	// header listener subscription
// 	cancelACKProcess context.CancelFunc

// 	// cli context
// 	cliCtx cliContext.CLIContext

// 	// queue connector
// 	queueConnector *QueueConnector

// 	// http client to subscribe to
// 	httpClient *httpClient.HTTP
// }

// // NewAckService returns new service object
// func NewAckService(cdc *codec.Codec, queueConnector *QueueConnector, httpClient *httpClient.HTTP) *AckService {
// 	// create logger
// 	logger := Logger.With("module", NoackService)

// 	// root chain instance
// 	rootchainInstance, err := helper.GetRootChainInstance()
// 	if err != nil {
// 		logger.Error("Error while getting root chain instance", "error", err)
// 		panic(err)
// 	}

// 	cliCtx := cliContext.NewCLIContext().WithCodec(cdc)
// 	cliCtx.BroadcastMode = client.BroadcastSync
// 	cliCtx.TrustNode = true

// 	// creating checkpointer object
// 	ackservice := &AckService{
// 		storageClient:     getBridgeDBInstance(viper.GetString(BridgeDBFlag)),
// 		rootChainInstance: rootchainInstance,

// 		cliCtx:         cliCtx,
// 		queueConnector: queueConnector,
// 		httpClient:     httpClient,
// 	}

// 	ackservice.BaseService = *common.NewBaseService(logger, NoackService, ackservice)
// 	return ackservice
// }

// // OnStart starts new block subscription
// func (ackService *AckService) OnStart() error {
// 	ackService.BaseService.OnStart() // Always call the overridden method.

// 	// create cancellable context
// 	ackCtx, cancelACKProcess := context.WithCancel(context.Background())
// 	ackService.cancelACKProcess = cancelACKProcess
// 	// start polling for checkpoint in buffer
// 	go ackService.startPollingCheckpoint(ackCtx, helper.GetConfig().NoACKPollInterval)

// 	// subscribed to new head
// 	ackService.Logger.Debug("Started ACK service")

// 	return nil
// }

// // OnStop stops all necessary go routines
// func (ackService *AckService) OnStop() {
// 	ackService.BaseService.OnStop() // Always call the overridden method.

// 	// cancel ack process
// 	ackService.cancelACKProcess()
// 	// close bridge db instance
// 	closeBridgeDBInstance()
// }

// func (ackService *AckService) startPollingCheckpoint(ctx context.Context, interval time.Duration) {
// 	ticker := time.NewTicker(interval)
// 	// stop ticker when everything done
// 	defer ticker.Stop()

// 	for {
// 		select {
// 		case <-ticker.C:
// 			go ackService.checkForCheckpoint()
// 		case <-ctx.Done():
// 			ticker.Stop()
// 			return
// 		}
// 	}
// }

// func (ackService *AckService) checkForCheckpoint() {
// 	currentHeaderNumber, err := ackService.rootChainInstance.CurrentHeaderBlock(nil)
// 	if err != nil {
// 		ackService.Logger.Error("Error while fetching current header block number", "error", err)
// 		return
// 	}

// 	// fetch last header number
// 	lastHeaderNumber := currentHeaderNumber.Uint64() // - helper.GetConfig().ChildBlockInterval
// 	if lastHeaderNumber == 0 {
// 		// First checkpoint required
// 		return
// 	}
// 	// get big int header number
// 	headerNumber := big.NewInt(0)
// 	headerNumber.SetUint64(lastHeaderNumber)

// 	// header block
// 	headerObject, err := ackService.rootChainInstance.HeaderBlocks(nil, headerNumber)
// 	if err != nil {
// 		ackService.Logger.Error("Error while fetching header block object", "error", err)
// 		return
// 	}

// 	// process checkpoint
// 	go ackService.processCheckpoint(headerObject.CreatedAt.Int64())
// }

// func (ackService *AckService) processCheckpoint(lastCreatedAt int64) {
// 	var index float64
// 	// if last created at ==0 , no checkpoint yet
// 	if lastCreatedAt == 0 {
// 		index = 1
// 	}

// 	checkpointCreationTime := time.Unix(lastCreatedAt, 0)
// 	currentTime := time.Now().UTC()
// 	timeDiff := currentTime.Sub(checkpointCreationTime)
// 	// check if last checkpoint was < NoACK wait time
// 	if timeDiff.Seconds() >= helper.GetConfig().NoACKWaitTime.Seconds() && index == 0 {
// 		index = math.Floor(timeDiff.Seconds() / helper.GetConfig().NoACKWaitTime.Seconds())
// 	}

// 	if index == 0 {
// 		return
// 	}

// 	// check if difference between no-ack time and current time
// 	lastNoAck := ackService.getLastNoAckTime()

// 	lastNoAckTime := time.Unix(int64(lastNoAck), 0)
// 	timeDiff = currentTime.Sub(lastNoAckTime)
// 	// if last no ack == 0 , first no-ack to be sent
// 	if currentTime.Sub(lastNoAckTime).Seconds() < helper.GetConfig().CheckpointBufferTime.Seconds() && lastNoAck != 0 {
// 		ackService.Logger.Debug("Cannot send multiple no-ack in short time", "timeDiff", currentTime.Sub(lastNoAckTime).Seconds(), "ExpectedDiff", helper.GetConfig().CheckpointBufferTime.Seconds())
// 		return
// 	}

// 	ackService.Logger.Debug("Fetching next proposers", "count", index)

// 	// check if same checkpoint still exists
// 	if ackService.isValidProposer(uint64(index), helper.GetAddress()) {
// 		ackService.Logger.Debug(
// 			"⛑ Sending NO ACK message",
// 			"currentTime", currentTime.String(),
// 			"proposerCount", index,
// 		)

// 		// send NO ACK
// 		msg := checkpointTypes.NewMsgCheckpointNoAck(
// 			hmtypes.BytesToHeimdallAddress(helper.GetAddress()),
// 			uint64(time.Now().UTC().Unix()),
// 		)

// 		// send
// 		err := ackService.queueConnector.BroadcastToHeimdall(msg)
// 		if err != nil {
// 			ackService.Logger.Error("Error while sending no-ack tx to Heimdall queue", "error", err)
// 			return
// 		}

// 		ackService.Logger.Info("No-ack transaction sent successfully", "index", index)
// 	}
// }

// func (ackService *AckService) getLastNoAckTime() uint64 {
// 	response, err := FetchFromAPI(ackService.cliCtx, GetHeimdallServerEndpoint(LastNoAckURL))
// 	if err != nil {
// 		ackService.Logger.Error("Unable to send request for checkpoint buffer", "Error", err)
// 		return 0
// 	}

// 	var noackObject Result
// 	if err := json.Unmarshal(response.Result, &noackObject); err != nil {
// 		ackService.Logger.Error("Error unmarshalling no-ack data ", "error", err)
// 		return 0
// 	}

// 	return noackObject.Result
// }

// func (ackService *AckService) isValidProposer(count uint64, address []byte) bool {
// 	ackService.Logger.Debug("Skipping proposers", "count", strconv.FormatUint(count, 10))
// 	response, err := FetchFromAPI(
// 		ackService.cliCtx,
// 		GetHeimdallServerEndpoint(fmt.Sprintf(ProposersURL, strconv.FormatUint(count, 10))),
// 	)
// 	if err != nil {
// 		ackService.Logger.Error("Unable to send request for next proposers", "Error", err)
// 		return false
// 	}

// 	// unmarshall data from buffer
// 	var proposers []hmtypes.Validator
// 	if err := json.Unmarshal(response.Result, &proposers); err != nil {
// 		ackService.Logger.Error("Error unmarshalling validator data ", "error", err)
// 		return false
// 	}

// 	ackService.Logger.Debug("Fetched proposers list", "numberOfProposers", count)
// 	for _, proposer := range proposers {
// 		if bytes.Equal(proposer.Signer.Bytes(), address) {
// 			return true
// 		}
// 	}

// 	return false
// }
