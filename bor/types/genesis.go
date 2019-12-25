package types

import (
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// GenesisState is the bor state that must be provided at genesis.
type GenesisState struct {
	SprintDuration uint64          `json:"sprint_duration" yaml:"sprint_duration"` // sprint duration
	SpanDuration   uint64          `json:"span_duration" yaml:"span_duration"`     // span duration ie number of blocks for which val set is frozen on heimdall
	ProducerCount  uint64          `json:"producer_count" yaml:"producer_count"`   // producer count per span
	Spans          []*hmTypes.Span `json:"spans" yaml:"spans"`                     // list of spans
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(sprintDuration uint64, spanDuration uint64, producerCount uint64, spans []*hmTypes.Span) GenesisState {
	return GenesisState{
		SprintDuration: sprintDuration,
		SpanDuration:   spanDuration,
		ProducerCount:  producerCount,
		Spans:          spans,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(DefaultSprintDuration, DefaultSpanDuration, DefaultProducerCount, nil)
}

// ValidateGenesis performs basic validation of bor genesis data returning an
// error for any failed validation criteria.
func ValidateGenesis(data GenesisState) error { return nil }

// genFirstSpan generates default first valdiator producer set
func genFirstSpan(valset hmTypes.ValidatorSet) []*hmTypes.Span {
	var firstSpan []*hmTypes.Span
	var selectedProducers []hmTypes.Validator
	for _, val := range valset.Validators {
		selectedProducers = append(selectedProducers, *val)
	}
	newSpan := hmTypes.NewSpan(0, 0, 0+DefaultSpanDuration-1, valset, selectedProducers, helper.GetConfig().BorChainID)
	firstSpan = append(firstSpan, &newSpan)
	return firstSpan
}
