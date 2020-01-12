package types

// query endpoints supported by the auth Querier
const (
	QueryRecord     = "record"
	QueryRecordList = "record-list"
)

// QueryRecordParams defines the params for querying accounts.
type QueryRecordParams struct {
	RecordID uint64
}

// NewQueryRecordParams creates a new instance of QueryRecordParams.
func NewQueryRecordParams(recordID uint64) QueryRecordParams {
	return QueryRecordParams{RecordID: recordID}
}

// QueryRecordListParams defines the params for querying accounts.
type QueryRecordListParams struct {
	Page  uint64
	Limit uint64
}

// NewQueryRecordListParams creates a new instance of QueryRecordListParams.
func NewQueryRecordListParams(page uint64, limit uint64) QueryRecordListParams {
	return QueryRecordListParams{Page: page, Limit: limit}
}
