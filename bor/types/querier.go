package types

// query endpoints supported by the auth Querier
const (
	QueryParams        = "params"
	QuerySpan          = "span"
	QuerySpanList      = "span-list"
	QueryLatestSpan    = "latest-span"
	QueryNextSpan      = "next-span"
	QueryNextProducers = "next-producers"

	ParamSpan          = "span"
	ParamSprint        = "sprint"
	ParamProducerCount = "producer-count"
	ParamLastEthBlock  = "last-eth-block"
)

// QuerySpanParams defines the params for querying accounts.
type QuerySpanParams struct {
	RecordID uint64
}

// NewQuerySpanParams creates a new instance of QuerySpanParams.
func NewQuerySpanParams(recordID uint64) QuerySpanParams {
	return QuerySpanParams{RecordID: recordID}
}

// QuerySpanListParams defines the params for querying accounts.
type QuerySpanListParams struct {
	Page  uint64
	Limit uint64
}

// NewQuerySpanListParams creates a new instance of QuerySpanListParams.
func NewQuerySpanListParams(page uint64, limit uint64) QuerySpanListParams {
	return QuerySpanListParams{Page: page, Limit: limit}
}
