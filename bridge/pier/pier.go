package pier

import (
	"sync"
	"time"

	"github.com/syndtr/goleveldb/leveldb"
)

const (
	chainSyncer       = "chain-syncer"
	maticCheckpointer = "matic-checkpointer"
	ACKService        = "ack-service"

	// TODO fetch port from config
	lastNoAckURL = "http://localhost:1317/checkpoint/last-no-ack"
	proposersURL = "http://localhost:1317/staking/proposer"

	bridgeDBFlag = "bridge-db"
	lastBlockKey = "last-block" // storage key

	defaultPollInterval           = 5 * 1000 // in milliseconds
	defaultMainPollInterval       = 5 * 1000 // in milliseconds
	defaultCheckpointPollInterval = 15 * time.Second
	defaultCheckpointLength       = 256                     // checkpoint number starts with 0, so length = defaultCheckpointLength -1
	maxCheckpointLength           = 4096                    // max blocks in one checkpoint
	defaultForcePushInterval      = maxCheckpointLength * 2 // in seconds (4096 * 2 seconds)
)

var bridgeDB *leveldb.DB
var bridgeDBOnce sync.Once
var bridgeDBCloseOnce sync.Once

// GetBridgeDBInstance get sington object for bridge-db
func getBridgeDBInstance(filePath string) *leveldb.DB {
	bridgeDBOnce.Do(func() {
		bridgeDB, _ = leveldb.OpenFile(filePath, nil)
	})

	return bridgeDB
}

// CloseBridgeDBInstance closes bridge-db instance
func closeBridgeDBInstance() {
	bridgeDBCloseOnce.Do(func() {
		bridgeDB.Close()
	})
}
