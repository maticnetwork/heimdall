package types

// bank module event types
const (
	EventTypeTopup    = "topup"
	EventTypeTransfer = "transfer"

	AttributeKeyRecipient   = "recipient"
	AttributeKeySender      = "sender"
	AttributeKeyValidatorID = "validator-id"
	AttributeKeyTopupAmount = "topup-amount"

	AttributeValueCategory = ModuleName
)
