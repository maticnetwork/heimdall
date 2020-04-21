//noalias
package types

// Slashing module event types
const (
	EventTypeSlash       = "slash"
	EventTypeSlashLimit  = "slash-limit"
	EventTypeTickConfirm = "tick-confirm"
	EventTypeLiveness    = "liveness"

	AttributeKeyAddress       = "address"
	AttributeKeyHeight        = "height"
	AttributeKeyPower         = "power"
	AttributeKeySlashedAmount = "slashed-amount"
	AttributeKeyReason        = "reason"
	AttributeKeyJailed        = "jailed"
	AttributeKeyMissedBlocks  = "missed_blocks"

	AttributeValueDoubleSign       = "double_sign"
	AttributeValueMissingSignature = "missing_signature"
	AttributeValueCategory         = ModuleName
)
