package staking

import (
	conf "github.com/maticnetwork/heimdall/helper"
	log "github.com/maticnetwork/heimdall/log"
)

// StakingLogger for staking module logger
var StakingLogger log.Logger

func init() {
	StakingLogger = conf.Logger.With("module", "staking")
}
