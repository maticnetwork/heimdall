package tags

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Checkpoint tags
var (
	Action         = sdk.TagAction
	BorSyncID      = "bor-sync-id"
	SpanID         = "span-id"
	SpanStartBlock = "start-block"
)
