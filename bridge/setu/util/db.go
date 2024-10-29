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

func GetBridgeDBInstance(filePath string) *leveldb.DB {
	bridgeDBOnce.Do(func() {
		var err error
		bridgeDB, err = leveldb.OpenFile(filePath, nil)
		if err != nil {
			fmt.Println("Error in Opening Database")
			log.Fatalln("Error in Opening Database", err)
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
