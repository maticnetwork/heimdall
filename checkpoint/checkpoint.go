package checkpoint

import (
	"github.com/maticnetwork/heimdall/helper"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

// CheckpointLogger for staking module logger
var CheckpointLogger tmlog.Logger

func init() {
	CheckpointLogger = helper.Logger.With("module", "checkpoint")
}
