package util

import (
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tendermint/libs/log"
)

var bridgeDB *leveldb.DB
var bridgeDBOnce sync.Once
var bridgeDBCloseOnce sync.Once

var dbLogger log.Logger

func init() {
	// create logger
	dbLogger = Logger.With("module", "db")
}

// GetBridgeDBInstance get sington object for bridge-db
func GetBridgeDBInstance(filePath string) *leveldb.DB {
	bridgeDBOnce.Do(func() {
		bridgeDB, _ = leveldb.OpenFile(filePath, nil)
	})

	return bridgeDB
}

// CloseBridgeDBInstance closes bridge-db instance
func closeBridgeDBInstance() {
	bridgeDBCloseOnce.Do(func() {
		if bridgeDB != nil {
			bridgeDB.Close()
		}
	})
}
