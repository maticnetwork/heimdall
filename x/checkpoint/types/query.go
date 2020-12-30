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
)
