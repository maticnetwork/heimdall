package util

import (
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

var bridgeDB *leveldb.DB
var bridgeDBOnce sync.Once
var bridgeDBCloseOnce sync.Once

// GetBridgeDBInstance get sington object for bridge-db
func GetBridgeDBInstance(filePath string) (*leveldb.DB, error) {
	var err error
	bridgeDBOnce.Do(func() {
		bridgeDB, err = leveldb.OpenFile(filePath, nil)
	})
	if err != nil {
		// Return nil and the error
		return nil, err
	}
	// Return the database instance
	return bridgeDB, nil
}

// CloseBridgeDBInstance closes bridge-db instance
func CloseBridgeDBInstance() {
	bridgeDBCloseOnce.Do(func() {
		if bridgeDB != nil {
			bridgeDB.Close()
		}
	})
}
