//noalias
package types

// Slashing module event types
const (
	EventTypeSlash       = "slash"
	EventTypeSlashLimit  = "slash-limit"
	EventTypeTickConfirm = "tick-confirm"
	EventTypeLiveness    = "liveness"

	AttributeKeyAddress       = "address"
	AttributeKeyValID         = "valid"
	AttributeKeyHeight        = "height"
	AttributeKeyPower         = "power"
	AttributeKeySlashedAmount = "slashed-amount"
	AttributeKeySlashInfoHash = "slash-info-hash"
	AttributeKeyProposer      = "proposer"
	AttributeKeyReason        = "reason"
	AttributeKeyJailed        = "jailed"
	AttributeKeyMissedBlocks  = "missed_blocks"

	AttributeValueDoubleSign       = "double_sign"
	AttributeValueMissingSignature = "missing_signature"
	AttributeValueCategory         = ModuleName
)
