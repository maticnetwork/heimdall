package types

import (
	"github.com/cosmos/cosmos-sdk/x/params"
)

const (
	// DefaultSprintDuration sprint for blocks
	DefaultSprintDuration uint64 = 64

	// DefaultSpanDuration number of blocks for which span is frozen on heimdall
	DefaultSpanDuration uint64 = 100 * DefaultSprintDuration

	// DefaultFirstSpanDuration first span duration
	DefaultFirstSpanDuration uint64 = 256

	// Slot cost for validator
	SlotCost int64 = 1

	// Number of Producers to be selected per span
	DefaultProducerCount uint64 = 4
)

// ParamStoreKeySprintDuration is store's key for SprintDuration
var ParamStoreKeySprintDuration = []byte("sprintduration")

// ParamStoreKeySpanDuration is store's key for SpanDuration
var ParamStoreKeySpanDuration = []byte("spanduration")

var ParamStoreKeyNumOfProducers = []byte("producercount")

// ParamKeyTable type declaration for parameters
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable(
		ParamStoreKeySprintDuration, DefaultSprintDuration,
		ParamStoreKeySpanDuration, DefaultSpanDuration,
		ParamStoreKeyNumOfProducers, DefaultProducerCount,
	)
}
