package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Bank errors reserve 5400 ~ 5499.
const (
	CodeStateAlreadySynced sdk.CodeType = 5400
	CodeStateInvalid       sdk.CodeType = 5401
)

// ErrStateAlreadySynced represents state sync error
func ErrStateAlreadySynced(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeStateAlreadySynced, "State already synced")
}

// ErrStateInvalid represents state error
func ErrStateInvalid(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeStateInvalid, "State is invalid")
}
