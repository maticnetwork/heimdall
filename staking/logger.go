package staking

import (
	"github.com/maticnetwork/heimdall/helper"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

// StakingLogger for staking module logger
var StakingLogger tmlog.Logger

func init() {
	StakingLogger = helper.Logger.With("module", "staking")
}
