package types

// Span stores details for a span on Bor chain
// span is indexed by start block
type Span struct {
	StartBlock        uint64
	EndBlock          uint64
	ValidatorSet      ValidatorSet
	SelectedProducers []Validator
	Signatures        []byte
}

// NewSpan creates new span
func NewSpan(startBlock uint64, endBlock uint64, validatorSet ValidatorSet, selectedProducers []Validator, sigs [][]byte) Span {
	var signatures []byte
	for _, sign := range sigs {
		signatures = append(signatures[:], sign[:]...)
	}
	return Span{
		StartBlock:        startBlock,
		EndBlock:          endBlock,
		ValidatorSet:      validatorSet,
		SelectedProducers: selectedProducers,
		Signatures:        signatures,
	}
}

// GetSignatures returns signatures for a particular
func (s *Span) GetSignatures() (sigs [][]byte) {
	return
}
