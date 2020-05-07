package types

// query endpoints supported by the auth Querier
const (
	QueryParams        = "params"
	QuerySpan          = "span"
	QuerySpanList      = "span-list"
	QueryLatestSpan    = "latest-span"
	QueryNextSpan      = "next-span"
	QueryNextProducers = "next-producers"
	QueryNextSpanSeed  = "next-span-seed"

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
