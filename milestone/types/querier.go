package types

// query endpoints supported by the auth Querier
const (
	QueryParams               = "params"
	QueryLatestMilestone      = "milestone-latest"
	QueryMilestoneByNumber    = "milestone-by-number"
	QueryProposer             = "is-proposer"
	QueryCurrentProposer      = "current-proposer"
	StakingQuerierRoute       = "staking"
	QueryCount                = "count"
	QueryLatestNoAckMilestone = "latest-no-ack-milestone"
	QueryNoAckMilestoneByID   = "no-ack-milestone-by-id"
)

// QueryBorChainID defines the params for querying with bor chain id

// QueryMilestoneParams defines the params for querying accounts.
type QueryMilestoneParams struct {
	Number uint64
}

// NewQueryMilestoneParams creates a new instance of QueryMilestoneHeaderIndex.
func NewQueryMilestoneParams(number uint64) QueryMilestoneParams {
	return QueryMilestoneParams{Number: number}
}

type QueryBorChainID struct {
	BorChainID string
}

// NewQueryBorChainID creates a new instance of QueryBorChainID with give chain id
func NewQueryBorChainID(chainID string) QueryBorChainID {
	return QueryBorChainID{BorChainID: chainID}
}

type QueryMilestoneID struct {
	MilestoneID string
}

// NewQueryMilestoneParams creates a new instance of QueryMilestoneHeaderIndex.
func NewQueryMilestoneID(id string) QueryMilestoneID {
	return QueryMilestoneID{MilestoneID: id}
}
