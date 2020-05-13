//noalias
package types

// Slashing module event types
const (
	EventTypeSlash       = "slash"
	EventTypeSlashLimit  = "slash-limit"
	EventTypeTickConfirm = "tick-confirm"
	EventTypeTickAck     = "tick-ack"
	EventTypeUnjail      = "unjail"
	EventTypeLiveness    = "liveness"

	AttributeKeyAddress        = "address"
	AttributeKeyValID          = "valid"
	AttributeKeyHeight         = "height"
	AttributeKeyPower          = "power"
	AttributeKeySlashedAmount  = "slashed-amount"
	AttributeKeySlashInfoBytes = "slash-info-bytes"
	AttributeKeyProposer       = "proposer"
	AttributeKeyReason         = "reason"
	AttributeKeyJailed         = "jailed"
	AttributeKeyMissedBlocks   = "missed_blocks"

	AttributeValueDoubleSign       = "double_sign"
	AttributeValueMissingSignature = "missing_signature"
	AttributeValueCategory         = ModuleName
)
