package keeper

import (
	"github.com/maticnetwork/heimdall/x/staking/types"
)

var _ types.QueryServer = Keeper{}
