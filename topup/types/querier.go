package types

const (
	QuerySequence = "sequence"
)

// QuerySequenceParams defines the params for querying an account Sequence.
type QuerySequenceParams struct {
	TxHash   string
	LogIndex uint64
}

// NewQuerySequenceParams creates a new instance of QuerySequenceParams.
func NewQuerySequenceParams(txHash string, logIndex uint64) QuerySequenceParams {
	return QuerySequenceParams{TxHash: txHash, LogIndex: logIndex}
}
