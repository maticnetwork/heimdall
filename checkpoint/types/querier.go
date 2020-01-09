package types

// query endpoints supported by the auth Querier
const (
	QueryAckCount           = "ack-count"
	QueryInitialAccountRoot = "initial-account-root"
	QueryAccountProof       = "dividend-account-proof"
	QueryCheckpoint         = "checkpoint"
	QueryCheckpointBuffer   = "checkpoint-buffer"
	QueryLastNoAck          = "last-no-ack"
	QueryCheckpointList     = "checkpoint-list"
)

// QueryCheckpointParams defines the params for querying accounts.
type QueryCheckpointParams struct {
	HeaderIndex uint64
}

// NewQueryCheckpointParams creates a new instance of QueryCheckpointHeaderIndex.
func NewQueryCheckpointParams(headerIndex uint64) QueryCheckpointParams {
	return QueryCheckpointParams{HeaderIndex: headerIndex}
}

// QueryCheckpointListParams defines the params for querying accounts.
type QueryCheckpointListParams struct {
	Page  uint64
	Limit uint64
}

// NewQueryCheckpointListParams creates a new instance of QueryCheckpointHeaderIndex.
func NewQueryCheckpointListParams(page uint64, limit uint64) QueryCheckpointListParams {
	return QueryCheckpointListParams{Page: page, Limit: limit}
}
