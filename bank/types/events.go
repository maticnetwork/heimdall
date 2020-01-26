package types

// bank module event types
const (
	EventTypeTopup       = "topup"
	EventTypeFeeWithdraw = "fee-withdraw"
	EventTypeTransfer    = "transfer"

	AttributeKeyRecipient         = "recipient"
	AttributeKeySender            = "sender"
	AttributeKeyValidatorID       = "validator-id"
	AttributeKeyTopupAmount       = "topup-amount"
	AttributeKeyFeeWithdrawAmount = "fee-withdraw-amount"

	AttributeValueCategory = ModuleName
)
