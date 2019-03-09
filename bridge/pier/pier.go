package pier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/helper"
	hmtypes "github.com/maticnetwork/heimdall/types"
)

const (
	chainSyncer       = "chain-syncer"
	maticCheckpointer = "matic-checkpointer"
	noackService      = "checkpoint-no-ack"

	// TODO fetch port from config
	lastNoAckURL      = "http://localhost:1317/checkpoint/last-no-ack"
	proposersURL      = "http://localhost:1317/staking/proposer/%v"
	lastCheckpointURL = "http://localhost:1317/checkpoint/buffer"

	bridgeDBFlag = "bridge-db"
	lastBlockKey = "last-block" // storage key
)

var defaultForcePushInterval = helper.GetConfig().MaxCheckpointLength * 2 // in seconds (1024 * 2 seconds)

var bridgeDB *leveldb.DB
var bridgeDBOnce sync.Once
var bridgeDBCloseOnce sync.Once

var pierLogger log.Logger

func init() {
	// create logger
	pierLogger = log.NewTMLogger(log.NewSyncWriter(os.Stdout)).With("module", "pier")
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

func isProposer() bool {
	count := uint64(1)
	resp, err := http.Get(fmt.Sprintf(proposersURL, strconv.FormatUint(count, 10)))
	if err != nil {
		pierLogger.Error("Unable to send request to get proposer", "Error", err)
		return false
	}
	defer resp.Body.Close()
	pierLogger.Debug("Request for proposer was successfull", "Count", count, "Status", resp.Status)
	if resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			pierLogger.Error("Unable to read data from response", "Error", err)
			return false
		}

		// unmarshall data from buffer
		var proposers []hmtypes.Validator
		if err := json.Unmarshal(body, &proposers); err != nil {
			pierLogger.Error("Error unmarshalling validator data ", "error", err)
			return false
		}

		// no proposer found
		if len(proposers) == 0 {
			return false
		}

		// get first proposer
		proposer := proposers[0]
		if bytes.Equal(proposer.Signer.Bytes(), helper.GetAddress()) {
			return true
		}
	} else {
		pierLogger.Error("Error while fetching proposer", "status", resp.StatusCode)
	}
	return false
}
