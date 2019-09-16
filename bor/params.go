package bor

import (
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	// DefaultParamspace for params keeper
	DefaultParamspace = "bor"

	// DefaultSprintDuration sprint for blocks
	DefaultSprintDuration uint64 = 64
	// DefaultSpanDuration number of blocks for which span is frozen on heimdall
	DefaultSpanDuration uint64 = 100 * DefaultSprintDuration
	// Slot cost for validator
	SlotCost uint64 = 10
	// Number of Producers to be selected per span
	NumProducers uint64 = 4
)

// ParamStoreKeySprintDuration is store's key for SprintDuration
var ParamStoreKeySprintDuration = []byte("sprintduration")

// ParamStoreKeySpanDuration is store's key for SpanDuration
var ParamStoreKeySpanDuration = []byte("spanduration")

// ParamKeyTable type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		ParamStoreKeySprintDuration, DefaultSprintDuration,
		ParamStoreKeySpanDuration, DefaultSpanDuration,
	)
}
