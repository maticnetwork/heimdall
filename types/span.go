package types

import (
	"fmt"
)

// Span stores details for a span on Bor chain
// span is indexed by start block
type Span struct {
	StartBlock        uint64       `json:"start_block" yaml:"start_block"`
	EndBlock          uint64       `json:"end_block" yaml:"end_block"`
	ValidatorSet      ValidatorSet `json:"validator_set" yaml:"validator_set"`
	SelectedProducers []Validator  `json:"selected_producers" yaml:"selected_producers"`
	Signatures        []byte       `json:"signatures" yaml:"signatures"`
	ChainID           string       `json:"bor_chain_id" yaml:"bor_chain_id"`
}

// NewSpan creates new span
func NewSpan(startBlock uint64, endBlock uint64, validatorSet ValidatorSet, selectedProducers []Validator, chainID string) Span {
	return Span{
		StartBlock:        startBlock,
		EndBlock:          endBlock,
		ValidatorSet:      validatorSet,
		SelectedProducers: selectedProducers,
		ChainID:           chainID,
	}
}

// AddSigs adds signatures to span
func (s *Span) AddSigs(sigs []byte) {
	s.Signatures = sigs
}

// GetSignatures returns signatures for a particular
func (s *Span) GetSignatures() (sigs []byte) {
	return s.Signatures
}

// String returns the string representatin of span
func (s *Span) String() string {
	return fmt.Sprintf(
		"Span {%v (%d:%d) %v %v}",
		s.ChainID,
		s.StartBlock,
		s.EndBlock,
		s.ValidatorSet,
		s.SelectedProducers,
	)
}
