package keeper

import (
	"github.com/maticnetwork/heimdall/x/blog/types"
)

var _ types.QueryServer = Keeper{}
