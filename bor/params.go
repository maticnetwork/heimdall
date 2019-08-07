package bor

import (
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	// DefaultParamspace for params keeper
	DefaultParamspace = "bor"

	// DefaultSprintDuration sprint for blocks
	DefaultSprintDuration = 64
	// DefaultSpanDuration number of blocks for which span is frozen on heimdall
	DefaultSpanDuration = 100 * DefaultSprintDuration
)

// ParamStoreKeySprintDuration is store's key for SprintDuration
var ParamStoreKeySprintDuration = []byte("sprint-duration")

// ParamStoreKeySpanDuration is store's key for SpanDuration
var ParamStoreKeySpanDuration = []byte("span-duration")

// ParamKeyTable type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		ParamStoreKeySprintDuration, DefaultSprintDuration,
		ParamStoreKeySpanDuration, DefaultSpanDuration,
	)
}
