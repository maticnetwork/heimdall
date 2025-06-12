package types

import (
	"fmt"

	hmTypes "github.com/maticnetwork/heimdall/types"
)

func IsBlockCloseToSpanEnd(blockNumber, spanEnd uint64) bool {
	// Check if the block number is within 100 blocks of the span end
	return blockNumber <= spanEnd && blockNumber >= (spanEnd-100)
}

// CalcCurrentBorSpanId computes the Bor span ID corresponding to latestBorBlock,
// using latestHeimdallSpan as the reference. It returns an error if inputs are invalid
// (nil span, zero or negative span length) or if arithmetic overflow is detected.
func CalcCurrentBorSpanId(latestBorBlock uint64, latestHeimdallSpan *hmTypes.Span) (uint64, error) {
	if latestHeimdallSpan == nil {
		return 0, fmt.Errorf("nil Heimdall span provided")
	}
	if latestHeimdallSpan.EndBlock < latestHeimdallSpan.StartBlock {
		return 0, fmt.Errorf(
			"invalid Heimdall span: EndBlock (%d) must be >= StartBlock (%d)",
			latestHeimdallSpan.EndBlock,
			latestHeimdallSpan.StartBlock,
		)
	}

	if latestBorBlock < latestHeimdallSpan.StartBlock {
		return 0, fmt.Errorf(
			"latestBorBlock (%d) must be >= Heimdall span StartBlock (%d)",
			latestBorBlock,
			latestHeimdallSpan.StartBlock,
		)
	}

	if latestBorBlock <= latestHeimdallSpan.EndBlock {
		return latestHeimdallSpan.ID, nil
	}

	spanLength := latestHeimdallSpan.EndBlock - latestHeimdallSpan.StartBlock + 1

	offset := latestBorBlock - latestHeimdallSpan.StartBlock
	quotient := offset / spanLength

	spanId := latestHeimdallSpan.ID + quotient

	if spanId < latestHeimdallSpan.ID {
		return 0, fmt.Errorf(
			"overflow detected computing span ID: reference ID=%d quotient=%d",
			latestHeimdallSpan.ID, quotient,
		)
	}

	return spanId, nil
}

func GenerateBorCommittedSpans(latestBorBlock uint64, latestBorUsedSpan *hmTypes.Span) []hmTypes.Span {
	spans := []hmTypes.Span{}
	spanLength := latestBorUsedSpan.EndBlock - latestBorUsedSpan.StartBlock
	prevSpan := latestBorUsedSpan
	for latestBorBlock > prevSpan.EndBlock {
		startBlock := prevSpan.EndBlock + 1
		newSpan := hmTypes.Span{
			ID:                prevSpan.ID + 1,
			StartBlock:        startBlock,
			EndBlock:          startBlock + spanLength,
			ChainID:           latestBorUsedSpan.ChainID,
			SelectedProducers: latestBorUsedSpan.SelectedProducers,
			ValidatorSet:      latestBorUsedSpan.ValidatorSet,
		}
		spans = append(spans, newSpan)
		prevSpan = &newSpan
	}
	return spans
}
