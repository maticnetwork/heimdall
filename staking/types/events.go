package types

// Checkpoint tags
var (
	EventTypeNewProposer   = "new-proposer"
	EventTypeValidatorJoin = "validator-join"
	EventTypeSignerUpdate  = "signer-update"
	EventTypeStakeUpdate   = "stake-update"
	EventTypeValidatorExit = "validator-exit"
	EventTypeDelegatorBond = "delegator-bond"

	AttributeKeySigner            = "signer"
	AttributeKeyDeactivationEpoch = "deactivation-epoch"
	AttributeKeyActivationEpoch   = "activation-epoch"
	AttributeKeyValidatorID       = "validator-id"
	AttributeKeyDelegatorID       = "delegator-id"
	AttributeKeyUpdatedAt         = "updated-at"

	AttributeValueCategory = ModuleName
)
