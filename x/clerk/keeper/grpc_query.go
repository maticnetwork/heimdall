package keeper

import (
	"github.com/maticnetwork/heimdall/x/clerk/types"
)

var _ types.QueryServer = Keeper{}
