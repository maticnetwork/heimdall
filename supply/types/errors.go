package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// supply errors reserve 2000 ~ 2100.
const (
	CodeNoAccountCreated sdk.CodeType = 2000
	CodeNoPermission     sdk.CodeType = 2001
)

// ErrNoAccountCreated is an error
func ErrNoAccountCreated(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNoAccountCreated, "account isn't able to be created")
}

// ErrNoPermission is an error
func ErrNoPermission(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNoPermission, "module account does not have permissions")
}
