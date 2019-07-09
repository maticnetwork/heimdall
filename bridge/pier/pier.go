package pier

import (
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/helper"
)

var defaultForcePushInterval = helper.GetConfig().MaxCheckpointLength * 2 // in seconds (1024 * 2 seconds)

var bridgeDB *leveldb.DB
var bridgeDBOnce sync.Once
var bridgeDBCloseOnce sync.Once

var pierLogger log.Logger

func init() {
	// create logger
	pierLogger = Logger.With("module", "pier")
}

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
