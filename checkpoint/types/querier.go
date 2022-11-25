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
	QueryCurrentProposer  = "current-proposer"
	StakingQuerierRoute   = "staking"

	////####Milestone Module #####

	QueryLatestMilestone      = "milestone-latest"
	QueryMilestoneByNumber    = "milestone-by-number"
	QueryCount                = "count"
	QueryLatestNoAckMilestone = "latest-no-ack-milestone"
	QueryNoAckMilestoneByID   = "no-ack-milestone-by-id"
)

// QueryCheckpointParams defines the params for querying accounts.
type QueryCheckpointParams struct {
	Number uint64
}

// NewQueryCheckpointParams creates a new instance of QueryCheckpointHeaderIndex.
func NewQueryCheckpointParams(number uint64) QueryCheckpointParams {
	return QueryCheckpointParams{Number: number}
}

// QueryBorChainID defines the params for querying with bor chain id
type QueryBorChainID struct {
	BorChainID string
}

// NewQueryBorChainID creates a new instance of QueryBorChainID with give chain id
func NewQueryBorChainID(chainID string) QueryBorChainID {
	return QueryBorChainID{BorChainID: chainID}
}

/////////######Milestone###############

// QueryMilestoneParams defines the params for querying accounts.
type QueryMilestoneParams struct {
	Number uint64
}

// NewQueryMilestoneParams creates a new instance of QueryMilestoneHeaderIndex.
func NewQueryMilestoneParams(number uint64) QueryMilestoneParams {
	return QueryMilestoneParams{Number: number}
}

type QueryMilestoneID struct {
	MilestoneID string
}

// NewQueryMilestoneParams creates a new instance of QueryMilestoneHeaderIndex.
func NewQueryMilestoneID(id string) QueryMilestoneID {
	return QueryMilestoneID{MilestoneID: id}
}
