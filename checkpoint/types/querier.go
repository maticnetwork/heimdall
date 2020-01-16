package types

// query endpoints supported by the auth Querier
const (
	QueryAckCount         = "ack-count"
	QueryCheckpoint       = "checkpoint"
	QueryCheckpointBuffer = "checkpoint-buffer"
	QueryLastNoAck        = "last-no-ack"
	QueryCheckpointList   = "checkpoint-list"
)

// QueryCheckpointParams defines the params for querying accounts.
type QueryCheckpointParams struct {
	HeaderIndex uint64
}

// NewQueryCheckpointParams creates a new instance of QueryCheckpointHeaderIndex.
func NewQueryCheckpointParams(headerIndex uint64) QueryCheckpointParams {
	return QueryCheckpointParams{HeaderIndex: headerIndex}
}
