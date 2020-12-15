package types

import (
	"fmt"
	"sort"
)

// Span stores details for a span on Bor chain
// span is indexed by start block
type Span struct {
	ID                uint64       `json:"span_id" yaml:"span_id"`
	StartBlock        uint64       `json:"start_block" yaml:"start_block"`
	EndBlock          uint64       `json:"end_block" yaml:"end_block"`
	ValidatorSet      ValidatorSet `json:"validator_set" yaml:"validator_set"`
	SelectedProducers []Validator  `json:"selected_producers" yaml:"selected_producers"`
	ChainID           string       `json:"bor_chain_id" yaml:"bor_chain_id"`
}

// NewSpan creates new span
func NewSpan(id uint64, startBlock uint64, endBlock uint64, validatorSet ValidatorSet, selectedProducers []Validator, chainID string) Span {
	return Span{
		ID:                id,
		StartBlock:        startBlock,
		EndBlock:          endBlock,
		ValidatorSet:      validatorSet,
		SelectedProducers: selectedProducers,
		ChainID:           chainID,
	}
}

// String returns the string representations of span
func (s *Span) String() string {
	return fmt.Sprintf(
		"Span %v {%v (%d:%d) %v}",
		s.ID,
		s.ChainID,
		s.StartBlock,
		s.EndBlock,
		s.SelectedProducers,
	)
}

// SortSpanByID sorts spans by SpanID
func SortSpanByID(a []*Span) {
	sort.Slice(a, func(i, j int) bool {
		return a[i].ID < a[j].ID
	})
}
