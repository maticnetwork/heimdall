package types

import "fmt"

// Span stores details for a span on Bor chain
// span is indexed by start block
type Span struct {
	StartBlock        uint64
	EndBlock          uint64
	ValidatorSet      ValidatorSet
	SelectedProducers []Validator
	Signatures        []byte
	ChainID           string
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
