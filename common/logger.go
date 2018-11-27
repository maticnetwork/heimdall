package common

import (
	tmlog "github.com/tendermint/tendermint/libs/log"

	"github.com/maticnetwork/heimdall/helper"
)

// CheckpointLogger for checkpoint module logger
var CheckpointLogger tmlog.Logger
var StakingLogger tmlog.Logger

func init() {
	CheckpointLogger = helper.Logger.With("module", "checkpoint")
	StakingLogger = helper.Logger.With("module", "staking")
}
