package types

import "github.com/ethereum/go-ethereum/common"

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

// QuerySpanSeedResponse defines the response to a span seed query
type QuerySpanSeedResponse struct {
	Seed       common.Hash    `json:"seed"`
	SeedAuthor common.Address `json:"seed_author"`
}

// NewQuerySpanSeedResponse creates a new instance of QuerySpanSeedResponse.
func NewQuerySpanSeedResponse(seed common.Hash, seedAuthor common.Address) QuerySpanSeedResponse {
	return QuerySpanSeedResponse{Seed: seed, SeedAuthor: seedAuthor}
}
