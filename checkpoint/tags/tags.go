package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Checkpoint tags
var (
	Action      = sdk.TagAction
	Proposer    = "proposer"
	StartBlock  = "start-block"
	EndBlock    = "end-block"
	HeaderIndex = "header-index"
	NewProposer = "new-proposer"
)
