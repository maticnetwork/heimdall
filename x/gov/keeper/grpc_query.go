package keeper

import (
	"github.com/maticnetwork/heimdall/x/gov/types"
)

var _ types.QueryServer = Keeper{}
