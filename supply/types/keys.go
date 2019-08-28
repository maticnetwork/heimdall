package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "supply"

	// StoreKey is the store key string for supply
	StoreKey = ModuleName

	// RouterKey is the message route for supply
	RouterKey = ModuleName

	// QuerierRoute is the querier route for supply
	QuerierRoute = ModuleName

	// DefaultParamspace default name for parameter store
	DefaultParamspace = ModuleName

	// DefaultCodespace code space
	DefaultCodespace sdk.CodespaceType = ModuleName
)
