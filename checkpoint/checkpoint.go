package checkpoint

import (
	conf "github.com/maticnetwork/heimdall/helper"
	log "github.com/maticnetwork/heimdall/log"
)

// CheckpointLogger for staking module logger
var CheckpointLogger log.Logger

func init() {
	CheckpointLogger = conf.Logger.With("module", "checkpoint")
}
