package keeper

import (
	"github.com/maticnetwork/heimdall/x/chainmanager/types"
)

var _ types.QueryServer = Keeper{}
