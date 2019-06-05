package common

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

// CheckpointLogger for checkpoint module logger
var CheckpointLogger tmlog.Logger
var StakingLogger tmlog.Logger
var HelperLogger tmlog.Logger

func InitLoggers(ctx sdk.Context) {
	CheckpointLogger = ctx.Logger().With("module", "checkpoint")
	StakingLogger = ctx.Logger().With("module", "staking")
}
