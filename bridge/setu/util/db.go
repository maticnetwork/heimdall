package util

import (
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

var (
	bridgeDB          *leveldb.DB
	bridgeDBOnce      sync.Once
	bridgeDBCloseOnce sync.Once
)

// GetBridgeDBInstance get sington object for bridge-db
func GetBridgeDBInstance(filePath string) *leveldb.DB {
	bridgeDBOnce.Do(func() {
		bridgeDB, _ = leveldb.OpenFile(filePath, nil)
	})

	return bridgeDB
}

// CloseBridgeDBInstance closes bridge-db instance
func CloseBridgeDBInstance() {
	bridgeDBCloseOnce.Do(func() {
		if bridgeDB != nil {
			bridgeDB.Close()
		}
	})
}
