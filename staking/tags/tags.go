package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Checkpoint tags
var (
	Action              = sdk.TagAction
	NewProposerSelected = "new-proposer"
	ValidatorJoin       = "validator-join"
	ValidatorUpdate     = "validator-update"
	ValidatorExit       = "validator-exit"
	DeactivationEpoch   = "deactivation-epoch"
	ActivationEpoch     = "activation-epoch"
	ValidatorID         = "validator-id"
)
