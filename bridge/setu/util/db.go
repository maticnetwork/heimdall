package util

import (
	"fmt"
	"log"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
)

var bridgeDB *leveldb.DB
var bridgeDBOnce sync.Once
var bridgeDBCloseOnce sync.Once

// GetBridgeDBInstance get sington object for bridge-db
func GetBridgeDBInstance(filePath string) *leveldb.DB {
	bridgeDBOnce.Do(func() {
		var err error
		bridgeDB, err = leveldb.OpenFile(filePath, nil)
		fmt.Println(">>>>>>>>>>>>>>>>> bridgeDB filePath", filePath)
		fmt.Println(">>>>>>>>>>>>>>>>> bridgeDB", bridgeDB)
		if err != nil {
			log.Fatalln("Error in Bor Opening Database", err.Error())
		}
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
