package types

// Checkpoint tags
var (
	EventTypeCheckpoint       = "checkpoint"
	EventTypeCheckpointAdjust = "checkpoint-adjust"
	EventTypeCheckpointAck    = "checkpoint-ack"
	EventTypeCheckpointNoAck  = "checkpoint-noack"

	AttributeKeyProposer    = "proposer"
	AttributeKeyStartBlock  = "start-block"
	AttributeKeyEndBlock    = "end-block"
	AttributeKeyHeaderIndex = "header-index"
	AttributeKeyNewProposer = "new-proposer"
	AttributeKeyRootHash    = "root-hash"
	AttributeKeyAccountHash = "account-hash"
	AttributeKeyHash        = "hash"

	EventTypeMilestone        = "milestone"
	EventTypeMilestoneTimeout = "milestone-timeout"

	AttributeKeyMilestoneID = "milestone-id"

	AttributeValueCategory = ModuleName
)
