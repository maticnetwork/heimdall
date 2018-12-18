package pier

import (
	"github.com/tendermint/tendermint/libs/common"
	"context"
	"github.com/syndtr/goleveldb/leveldb"
	hmtypes "github.com/maticnetwork/heimdall/types"
	"github.com/tendermint/tendermint/libs/log"
	"os"
	"github.com/spf13/viper"
	"time"
	"encoding/json"
	"github.com/maticnetwork/heimdall/helper"
	"math"
	"io/ioutil"
	"bytes"
	cliContext "github.com/cosmos/cosmos-sdk/client/context"
	"github.com/maticnetwork/heimdall/checkpoint"
	"net/http"
	"strconv"
)

type AckService struct {
	// Base service
	common.BaseService

	// storage client
	storageClient *leveldb.DB

	CheckpointChannel chan *hmtypes.CheckpointBlockHeader
	// cancel function for poll/subscription
	// cancel function for poll/subscription
	cancelSubscription context.CancelFunc
	// header listener subscription
	cancelHeaderProcess context.CancelFunc
}

// NewMaticCheckpointer returns new service object
func NewAckService() *AckService {
	// create logger
	logger := log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", ackService)

	// creating checkpointer object
	ackservice := &AckService{
		storageClient:     getBridgeDBInstance(viper.GetString(bridgeDBFlag)),
		CheckpointChannel:     make(chan *hmtypes.CheckpointBlockHeader),
	}

	ackservice.BaseService = *common.NewBaseService(logger,ackService , ackservice)
	return ackservice
}


// OnStart starts new block subscription
func (ackService *AckService) OnStart() error {
	ackService.BaseService.OnStart() // Always call the overridden method.

	// start polling for checkpoint in buffer
	go ackService.StartPollingCheckpoint(defaultCheckpointPollInterval)

	// subscribed to new head
	ackService.Logger.Debug("Subscribed to new head")

	return nil
}


func(ackService *AckService) StartPollingCheckpoint(interval time.Duration){
	ticker := time.NewTicker(interval)
	// stop ticker when everything done
	defer ticker.Stop()
	// write to channel when we receive checkpoint
	found := make(chan []byte)

	for {
		select {
		case data:=<-found:

			// unmarshall data from buffer
			var headerBlock hmtypes.CheckpointBlockHeader
			if  err:=json.Unmarshal(data,&headerBlock); err!=nil{
				ackService.Logger.Error("Error unmarshalling checkpoint data ","Error",err)
			}

			ackService.Logger.Info("Found Checkpoint in buffer!","Checkpoint",headerBlock.String())

			// sleep for timestamp+5 minutes
			checkpointCreationTime:= time.Unix(int64(headerBlock.TimeStamp),0)
			timeToWait := checkpointCreationTime.Add(helper.CheckpointBufferTime)
			timeDiff:=time.Now().Sub(checkpointCreationTime)

			var index float64
			if timeDiff >= helper.CheckpointBufferTime {
				index = math.Round(timeDiff.Minutes()/helper.MinutesAliveForBuffer)
			} else{
				time.Sleep(timeToWait.Sub(time.Now()))
				index = 1
			}

			// check if checkpoint still exists in buffer
			resp:=getCheckpointBuffer(ackService)
			body,err:=ioutil.ReadAll(resp.Body)
			if err!=nil{
				ackService.Logger.Error("Unable to read data from response","Error",err)
			}

			// if same checkpoint still exists
			if bytes.Compare(data,body)==0 && getValidProposers(ackService,int(index),helper.GetPubKey().Address().Bytes()){
				ackService.Logger.Debug("Sending NO ACK message","ACK-ETA",timeToWait.String(),"CurrentTime",time.Now().String(),"ProposerCount",index)
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

				resp, err := helper.SendTendermintRequest(cliContext.NewCLIContext(), txBytes)
				if err != nil {
					ackService.Logger.Error("Error while sending request to Tendermint", "error", err)
					return
				}
				ackService.Logger.Error("No ACK transaction sent","TxHash",resp.Hash,"Checkpoint",headerBlock.String())
			}
			return
		case t := <-ticker.C:
			ackService.Logger.Debug("Awaiting Checkpoint...", t)
			go readCheckpointBuffer(ackService,found)
		}
	}
	return
}

func readCheckpointBuffer(ackService *AckService,found chan<- []byte)  {
	resp:=getCheckpointBuffer(ackService)
	if resp.StatusCode!=204{
		ackService.Logger.Info("Checkpoint found in buffer")
		data,err:=ioutil.ReadAll(resp.Body)
		if err!=nil{
			ackService.Logger.Error("Unable to read data from response","Error",err)
		}
		found <-data
	}else if resp.StatusCode==204{
		ackService.Logger.Debug("Checkpoint not found in buffer","Status",resp.StatusCode)
	}
	defer resp.Body.Close()
}

func getCheckpointBuffer(ackService *AckService) (resp *http.Response) {
	resp,err:=http.Get(checkpointBufferURL)
	if err!=nil{
		ackService.Logger.Error("Unable to send request for checkpoint buffer","Error",err)
	}
	return resp
}

func getValidProposers(ackService *AckService,count int,address []byte)(bool){
	ackService.Logger.Debug("Fetching next proposers","Count",strconv.Itoa(count))
	resp,err:=http.Get(proposersURL+"/"+strconv.Itoa(count))
	if err!=nil{
		ackService.Logger.Error("Unable to send request for next proposers","Error",err)
		return false
	}

	body,err:=ioutil.ReadAll(resp.Body)
	if err!=nil{
		ackService.Logger.Error("Unable to read data from response","Error",err)
		return false
	}

	// unmarshall data from buffer
	var proposers []hmtypes.Validator
	if err:=json.Unmarshal(body,&proposers); err!=nil{
		ackService.Logger.Error("Error unmarshalling checkpoint data ","Error",err)
		return false
	}

	ackService.Logger.Debug("Fetched proposers list from heimdall","Index",count,"Proposers",proposers)

	for _,proposer:=range proposers {
		if bytes.Compare(proposer.Address.Bytes(),address)==0{
			return true
		}
	}
	return false
}

