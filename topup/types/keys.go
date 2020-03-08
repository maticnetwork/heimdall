package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "topup"

	// StoreKey is the store key string for bor
	StoreKey = ModuleName

	// RouterKey is the message route for bor
	RouterKey = ModuleName

	// QuerierRoute is the querier route for bor
	QuerierRoute = ModuleName

	// DefaultParamspace default name for parameter store
	DefaultParamspace = ModuleName

	// DefaultCodespace default code space
	DefaultCodespace sdk.CodespaceType = ModuleName
)
