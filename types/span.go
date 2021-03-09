package types

import (
	"sort"
)

// Span stores details for a span on Bor chain
// span is indexed by start block
// NewSpan creates new span
func NewSpan(id uint64, startBlock uint64, endBlock uint64, validatorSet ValidatorSet, selectedProducers []Validator, BorChainID string) Span {
	return Span{
		ID:                id,
		StartBlock:        startBlock,
		EndBlock:          endBlock,
		ValidatorSet:      validatorSet,
		SelectedProducers: selectedProducers,
		BorChainId:        BorChainID,
	}
}

// SortSpanByID sorts spans by SpanID
func SortSpanByID(a []*Span) {
	sort.Slice(a, func(i, j int) bool {
		return a[i].ID < a[j].ID
	})
}
