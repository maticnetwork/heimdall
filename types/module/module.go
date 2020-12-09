package module

import (
	"github.com/maticnetwork/heimdall/types"
)

// SideModule is the standard form for side tx elements of an application module
type SideModule interface {
	NewSideTxHandler() types.SideTxHandler
	NewPostTxHandler() types.PostTxHandler
}
