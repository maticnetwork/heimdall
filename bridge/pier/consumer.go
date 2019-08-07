// consumes all events from respective queues
// Deposit Event --> Mint transaction on BOR on the basis of validator set% deposit index
// Withdraw Event --> Burn transaction on BOR
// Validator Join/Exit/Power-change --> Validator set changes on BOR
// Checkpoint Propose --> MsgCheckpoint on Heimdall
// Checkpoint ACK --> MsgCheckpointACK on Heimdall
// Checkpoint NO-ACK --> Sends MsgCheckpointNoACK after x interval on Heimdall
// Validator Join/Exit/Power-change --> Validator set changes on Heimdall
package pier

import (
	"log"

	"github.com/spf13/viper"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tendermint/libs/common"
)

type ConsumerService struct {
	// Base service
	common.BaseService

	// storage client
	storageClient *leveldb.DB

	qConnector QueueConnector
}

// NewAckService returns new service object
func NewConsumerService(connector QueueConnector) *ConsumerService {
	// create logger
	logger := Logger.With("module", NoackService)
	// creating checkpointer object
	consumerService := &ConsumerService{
		storageClient: getBridgeDBInstance(viper.GetString(BridgeDBFlag)),
		qConnector:    connector,
	}
	consumerService.BaseService = *common.NewBaseService(logger, NoackService, consumerService)
	return consumerService
}

// OnStart starts new block subscription
func (consumer *ConsumerService) OnStart() error {
	consumer.BaseService.OnStart() // Always call the overridden method.
	if err := consumer.qConnector.ConsumeHeimdallQ(); err != nil {
		log.Fatalf("Cannot consume")
	}
	return nil
}

// OnStop stops all necessary go routines
func (consumer *ConsumerService) OnStop() {
	// Always call the overridden method.
	consumer.BaseService.OnStop()
	// close db
	closeBridgeDBInstance()
}
