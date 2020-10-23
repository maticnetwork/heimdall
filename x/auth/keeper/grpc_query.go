package keeper

import (
	"github.com/maticnetwork/heimdall/x/auth/types"
)

var _ types.QueryServer = Keeper{}
