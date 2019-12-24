package types

// Checkpoint tags
var (
	EventTypeCheckpoint      = "checkpoint"
	EventTypeCheckpointAck   = "checkpoint-ack"
	EventTypeCheckpointNoAck = "checkpoint-noack"

	AttributeKeyProposer    = "proposer"
	AttributeKeyStartBlock  = "start-block"
	AttributeKeyEndBlock    = "end-block"
	AttributeKeyHeaderIndex = "header-index"
	AttributeKeyNewProposer = "new-proposer"

	AttributeValueCategory = ModuleName
)
