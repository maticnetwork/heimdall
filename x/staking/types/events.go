package types

// Checkpoint tags
var (
	EventTypeValidatorJoin = "validator-join"
	EventTypeSignerUpdate  = "signer-update"
	EventTypeStakeUpdate   = "stake-update"
	EventTypeValidatorExit = "validator-exit"

	AttributeKeySigner         = "signer"
	AttributeKeyValidatorID    = "validator-id"
	AttributeKeyValidatorNonce = "validator-nonce"

	AttributeValueCategory = ModuleName
)
