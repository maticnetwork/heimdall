package pier

import (
	"context"
	"encoding/hex"
	"math/big"
	"os"
	"time"

	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/ethereum/go-ethereum"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tendermint/libs/common"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/checkpoint"
	"github.com/maticnetwork/heimdall/contracts/rootchain"
	"github.com/maticnetwork/heimdall/helper"
	"net/http"
	"io/ioutil"
	"encoding/json"
	hmtypes "github.com/maticnetwork/heimdall/types"
	"bytes"
	"strconv"
	"math"
)

// MaticCheckpointer to propose
type MaticCheckpointer struct {
	// Base service
	common.BaseService

	// storage client
	storageClient *leveldb.DB

	// ETH client
	MaticClient *ethclient.Client
	// ETH RPC client
	MaticRPCClient *rpc.Client
	// Mainchain client
	MainClient *ethclient.Client
	// Rootchain instance
	RootChainInstance *rootchain.Rootchain
	// header channel
	HeaderChannel chan *types.Header
	// cancel function for poll/subscription
	cancelSubscription context.CancelFunc
	// header listener subscription
	cancelHeaderProcess context.CancelFunc
}

// NewMaticCheckpointer returns new service object
func NewMaticCheckpointer() *MaticCheckpointer {
	// create logger
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", maticCheckpointer)

	// root chain instance
	rootchainInstance, err := helper.GetRootChainInstance()
	if err != nil {
		logger.Error("Error while getting root chain instance", "error", err)
		panic(err)
	}

	// creating checkpointer object
	checkpointer := &MaticCheckpointer{
		storageClient:     getBridgeDBInstance(viper.GetString(bridgeDBFlag)),
		MaticClient:       helper.GetMaticClient(),
		MaticRPCClient:    helper.GetMaticRPCClient(),
		MainClient:        helper.GetMainClient(),
		RootChainInstance: rootchainInstance,
		HeaderChannel:     make(chan *types.Header),
	}

	checkpointer.BaseService = *common.NewBaseService(logger, maticCheckpointer, checkpointer)
	return checkpointer
}

// StartHeaderProcess starts header process when they get new header
func (checkpointer *MaticCheckpointer) StartHeaderProcess(ctx context.Context) {
	for {
		select {
		case newHeader := <-checkpointer.HeaderChannel:
			checkpointer.sendRequest(newHeader)
		case <-ctx.Done():
			return
		}
	}
}

// OnStart starts new block subscription
func (checkpointer *MaticCheckpointer) OnStart() error {
	checkpointer.BaseService.OnStart() // Always call the overridden method.

	// create cancellable context
	ctx, cancelSubscription := context.WithCancel(context.Background())
	checkpointer.cancelSubscription = cancelSubscription

	// create cancellable context
	headerCtx, cancelHeaderProcess := context.WithCancel(context.Background())
	checkpointer.cancelHeaderProcess = cancelHeaderProcess

	// start header process
	go checkpointer.StartHeaderProcess(headerCtx)

	// start polling for checkpoint in buffer
	go checkpointer.StartPollingCheckpoint(defaultCheckpointPollInterval)

	// subscribe to new head
	subscription, err := checkpointer.MaticClient.SubscribeNewHead(ctx, checkpointer.HeaderChannel)
	if err != nil {
		// start go routine to poll for new header using client object
		go checkpointer.StartPolling(ctx, defaultPollInterval)
	} else {
		// start go routine to listen new header using subscription
		go checkpointer.StartSubscription(ctx, subscription)
	}

	// subscribed to new head
	checkpointer.Logger.Debug("Subscribed to new head")

	return nil
}

// OnStop stops all necessary go routines
func (checkpointer *MaticCheckpointer) OnStop() {
	checkpointer.BaseService.OnStop() // Always call the overridden method.

	// close bridge db instance
	closeBridgeDBInstance()

	// cancel subscription if any
	checkpointer.cancelSubscription()

	// cancel header process
	checkpointer.cancelHeaderProcess()
}

func (checkpointer *MaticCheckpointer) StartPolling(ctx context.Context, pollInterval int) {
	// How often to fire the passed in function in second
	interval := time.Duration(pollInterval) * time.Millisecond

	// Setup the ticket and the channel to signal
	// the ending of the interval
	ticker := time.NewTicker(interval)

	// start listening
	for {
		select {
		case <-ticker.C:
			header, err := checkpointer.MaticClient.HeaderByNumber(ctx, nil)
			if err == nil && header != nil {
				// send data to channel
				checkpointer.HeaderChannel <- header
			}
		case <-ctx.Done():
			ticker.Stop()
			return
		}
	}
}

func (checkpointer *MaticCheckpointer) StartSubscription(ctx context.Context, subscription ethereum.Subscription) {
	for {
		select {
		case err := <-subscription.Err():
			// stop service
			checkpointer.Logger.Error("Error while subscribing new blocks", "error", err)
			checkpointer.Stop()

			// cancel subscription
			checkpointer.cancelSubscription()
			return
		case <-ctx.Done():
			return
		}
	}
}

func (checkpointer *MaticCheckpointer) sendRequest(newHeader *types.Header) {
	checkpointer.Logger.Debug("New block detected", "blockNumber", newHeader.Number)
	lastCheckpointEnd, err := checkpointer.RootChainInstance.CurrentChildBlock(nil)
	if err != nil {
		checkpointer.Logger.Error("Error while fetching current child block from rootchain", "error", err)
		return
	}

	latest := newHeader.Number.Uint64()
	start := lastCheckpointEnd.Uint64()
	var end uint64

	// add 1 if start > 0
	if start > 0 {
		start = start + 1
	}

	// get diff
	diff := latest - start + 1

	// process if diff > 0 (positive)
	if diff > 0 {
		expectedDiff := diff - diff%defaultCheckpointLength
		if expectedDiff > 0 {
			expectedDiff = expectedDiff - 1
		}

		// cap with max checkpoint length
		if expectedDiff > maxCheckpointLength-1 {
			expectedDiff = maxCheckpointLength - 1
		}

		// get end result
		end = expectedDiff + start

		checkpointer.Logger.Debug("Calculating checkpoint eligibility", "latest", latest, "start", start, "end", end)
	}

	// Handle when block producers go down
	if end == 0 || end == start || (0 < diff && diff < defaultCheckpointLength) {
		currentHeaderBlockNumber, err := checkpointer.RootChainInstance.CurrentHeaderBlock(nil)
		if err != nil {
			checkpointer.Logger.Error("Error while fetching current header block number from rootchain", "error", err)
			return
		}

		// fetch current header block
		currentHeaderBlock, err := checkpointer.RootChainInstance.HeaderBlock(nil, currentHeaderBlockNumber.Sub(currentHeaderBlockNumber, big.NewInt(1)))
		if err != nil {
			checkpointer.Logger.Error("Error while fetching current header block object from rootchain", "error", err)
			return
		}

		lastCheckpointTime := currentHeaderBlock.CreatedAt.Int64()
		currentTime := time.Now().Unix()
		if currentTime-lastCheckpointTime > defaultForcePushInterval {
			checkpointer.Logger.Info("Force push checkpoint", "currentTime", currentTime, "lastCheckpointTime", lastCheckpointTime, "defaultForcePushInterval", defaultForcePushInterval)
			end = latest
		}
	}

	if end == 0 || start >= end {
		return
	}

	// Get root hash
	root, err := checkpoint.GetHeaders(start, end)
	if err != nil {
		return
	}

	checkpointer.Logger.Info("New checkpoint header created", "latest", latest, "start", start, "end", end, "root", hex.EncodeToString(root))

	// TODO submit checkcoint
	txBytes, err := helper.CreateTxBytes(
		checkpoint.NewMsgCheckpointBlock(
			ethCommon.BytesToAddress(helper.GetPubKey().Address().Bytes()),
			start,
			end,
			ethCommon.BytesToHash(root),
			uint64(time.Now().Unix()),
		),
	)

	if err != nil {
		checkpointer.Logger.Error("Error while creating tx bytes", "error", err)
		return
	}

	resp, err := helper.SendTendermintRequest(cliContext.NewCLIContext(), txBytes)
	if err != nil {
		checkpointer.Logger.Error("Error while sending request to Tendermint", "error", err)
		return
	}

	checkpointer.Logger.Error("Checkpoint sent successfully", "hash", hex.EncodeToString(resp.Hash), "start", start, "end", end, "root", root)
}

func(checkpointer *MaticCheckpointer) StartPollingCheckpoint(interval time.Duration){
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	found := make(chan []byte)
	for {
		select {
		case data:=<-found:

			// unmarshall data from buffer
			var headerBlock hmtypes.CheckpointBlockHeader
			if  err:=json.Unmarshal(data,&headerBlock); err!=nil{
				checkpointer.Logger.Error("Error unmarshalling checkpoint data ","Error",err)
			}

			checkpointer.Logger.Info("Found Checkpoint in buffer!","Checkpoint",headerBlock.String())

			// sleep for timestamp+5 minutes
			checkpointCreationTime:= time.Unix(int64(headerBlock.TimeStamp),0)
			timeToWait := checkpointCreationTime.Add(2*time.Minute)
			timeDiff:=time.Now().Sub(checkpointCreationTime)
			var index float64
			if timeDiff >= 2*time.Minute {
				index = math.Round(timeDiff.Minutes()/2)
			} else{
				time.Sleep(timeToWait.Sub(time.Now()))
				index = 1
			}

			// check if checkpoint still exists in buffer
			resp:=getCheckpointBuffer(checkpointer)
			body,err:=ioutil.ReadAll(resp.Body)
			if err!=nil{
				checkpointer.Logger.Error("Unable to read data from response","Error",err)
			}

			// if same checkpoint still exists
			if bytes.Compare(data,body)==0 && getValidProposers(checkpointer,int(index),helper.GetPubKey().Address().Bytes()){
				checkpointer.Logger.Debug("Sending NO ACK message","ACK-ETA",timeToWait.String(),"CurrentTime",time.Now().String(),"ProposerCount",index)
				// send NO ACK
				txBytes, err := helper.CreateTxBytes(
					checkpoint.NewMsgCheckpointNoAck(
						uint64(time.Now().Unix()),
					),
				)
				if err != nil {
					checkpointer.Logger.Error("Error while creating tx bytes", "error", err)
					return
				}

				resp, err := helper.SendTendermintRequest(cliContext.NewCLIContext(), txBytes)
				if err != nil {
					checkpointer.Logger.Error("Error while sending request to Tendermint", "error", err)
					return
				}
				checkpointer.Logger.Error("No ACK transaction sent","TxHash",resp.Hash,"Checkpoint",headerBlock.String())
			}
			return
		case t := <-ticker.C:
			checkpointer.Logger.Debug("Awaiting Checkpoint...", t)
			go func() {
				resp:=getCheckpointBuffer(checkpointer)

				if resp.StatusCode!=204{
					checkpointer.Logger.Info("Checkpoint found in buffer")
					data,err:=ioutil.ReadAll(resp.Body)
					if err!=nil{
						checkpointer.Logger.Error("Unable to read data from response","Error",err)
					}
					found <-data

				}else if resp.StatusCode==204{
					checkpointer.Logger.Debug("Checkpoint not found in buffer","Status",resp.StatusCode)
				}

				defer resp.Body.Close()
			}()
		}
	}
	return
}

func getCheckpointBuffer(checkpointer *MaticCheckpointer) (resp *http.Response) {
	resp,err:=http.Get(checkpointBufferURL)
	if err!=nil{
		checkpointer.Logger.Error("Unable to send request for checkpoint buffer","Error",err)
	}
	return resp
}

func getValidProposers(checkpointer *MaticCheckpointer,count int,address []byte)(bool){
	checkpointer.Logger.Debug("Fetching next proposers","Count",strconv.Itoa(count))
	resp,err:=http.Get(proposersURL+"/"+strconv.Itoa(count))
	if err!=nil{
		checkpointer.Logger.Error("Unable to send request for next proposers","Error",err)
		return false
	}

	body,err:=ioutil.ReadAll(resp.Body)
	if err!=nil{
		checkpointer.Logger.Error("Unable to read data from response","Error",err)
		return false
	}

	// unmarshall data from buffer
	var proposers []hmtypes.Validator
	if err:=json.Unmarshal(body,&proposers); err!=nil{
		checkpointer.Logger.Error("Error unmarshalling checkpoint data ","Error",err)
		return false
	}

	checkpointer.Logger.Debug("Fetched proposers list from heimdall","Index",count,"Proposers",proposers)

	for _,proposer:=range proposers {
		if bytes.Compare(proposer.Address.Bytes(),address)==0{
			return true
		}
	}
	return false
}
