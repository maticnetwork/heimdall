package keeper

import (
	"github.com/maticnetwork/heimdall/x/topup/types"
)

var _ types.QueryServer = Keeper{}
