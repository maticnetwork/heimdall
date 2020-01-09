package types

// bank module event types
const (
	EventTypeTopup         = "topup"
	EventTypeTopupWithdraw = "topup-withdraw"
	EventTypeTransfer      = "transfer"

	AttributeKeyRecipient           = "recipient"
	AttributeKeySender              = "sender"
	AttributeKeyValidatorID         = "validator-id"
	AttributeKeyTopupAmount         = "topup-amount"
	AttributeKeyTopupWithdrawAmount = "topup-withdraw-amount"

	AttributeValueCategory = ModuleName
)
