package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Checkpoint tags
var (
	Action         = sdk.TagAction
	RecordTxHash   = "record-tx-hash"
	RecordID       = "record-id"
	RecordContract = "record-contract"
	CreatedAt      = "created-at"
)
