package types

// query endpoints supported by the auth Querier
const (
	QueryParams          = "params"
	QueryMilestone       = "milestone"
	QueryProposer        = "is-proposer"
	QueryCurrentProposer = "current-proposer"
	StakingQuerierRoute  = "staking"
	QueryCount           = "count"
)

// QueryBorChainID defines the params for querying with bor chain id
type QueryBorChainID struct {
	BorChainID string
}

// NewQueryBorChainID creates a new instance of QueryBorChainID with give chain id
func NewQueryBorChainID(chainID string) QueryBorChainID {
	return QueryBorChainID{BorChainID: chainID}
}
