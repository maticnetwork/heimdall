package types

// query endpoints supported by the auth Querier
const (
	QueryParams           = "params"
	QueryAckCount         = "ack-count"
	QueryCheckpoint       = "checkpoint"
	QueryCheckpointBuffer = "checkpoint-buffer"
	QueryLastNoAck        = "last-no-ack"
	QueryCheckpointList   = "checkpoint-list"
	QueryNextCheckpoint   = "next-checkpoint"
	QueryProposer         = "is-proposer"
)

// QueryCheckpointParams defines the params for querying accounts.
type QueryCheckpointParams struct {
	HeaderIndex uint64
}

// NewQueryCheckpointParams creates a new instance of QueryCheckpointHeaderIndex.
func NewQueryCheckpointParams(headerIndex uint64) QueryCheckpointParams {
	return QueryCheckpointParams{HeaderIndex: headerIndex}
}
