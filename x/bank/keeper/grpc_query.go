package keeper

import (
	"github.com/maticnetwork/heimdall/x/bank/types"
)

var _ types.QueryServer = Keeper{}
