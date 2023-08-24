package types

const (
	QueryLatestMilestone      = "milestone-latest"
	QueryMilestoneByNumber    = "milestone-by-number"
	QueryCount                = "count"
	QueryLatestNoAckMilestone = "latest-no-ack-milestone"
	QueryNoAckMilestoneByID   = "no-ack-milestone-by-id"
)

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
