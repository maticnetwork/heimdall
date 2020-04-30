package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Bank errors reserve 5400 ~ 5499.
const (
	CodeEventRecordAlreadySynced sdk.CodeType = 5400
	CodeEventRecordInvalid       sdk.CodeType = 5401
	CodeEventRecordUpdate        sdk.CodeType = 5402
)

// ErrEventRecordAlreadySynced represents event sync error
func ErrEventRecordAlreadySynced(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEventRecordAlreadySynced, "Event record already synced")
}

// ErrEventRecordInvalid represents event error
func ErrEventRecordInvalid(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEventRecordInvalid, "Event record is invalid")
}

// ErrEventUpdate represents event update error
func ErrEventUpdate(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEventRecordUpdate, "Event record update error")
}
