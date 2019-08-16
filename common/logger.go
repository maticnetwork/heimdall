package common

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

// CheckpointLogger for checkpoint module logger
var CheckpointLogger tmlog.Logger
var StakingLogger tmlog.Logger
var HelperLogger tmlog.Logger
var BorLogger tmlog.Logger

// InitCheckpointLogger initialises logger for checkpoint module
func InitCheckpointLogger(ctx *sdk.Context) {
	CheckpointLogger = ctx.Logger().With("module", "checkpoint")
}

// InitStakingLogger initialises logger for staking module
func InitStakingLogger(ctx *sdk.Context) {
	StakingLogger = ctx.Logger().With("module", "staking")
}

// InitHelperLogger initialises logger for helper module
func InitHelperLogger(ctx *sdk.Context) {
	HelperLogger = ctx.Logger().With("module", "helper")
}

// InitBorLogger initalises logger for bor module
func InitBorLogger(ctx *sdk.Context) {
	BorLogger = ctx.Logger().With("module", "bor")
}
