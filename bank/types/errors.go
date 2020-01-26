package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Bank errors reserve 100 ~ 199.
const (
	CodeSendDisabled         sdk.CodeType = 101
	CodeInvalidInputsOutputs sdk.CodeType = 102
	CodeNoValidatorTopup     sdk.CodeType = 103
	CodeNoBalanceToWithdraw  sdk.CodeType = 104
)

// ErrNoInputs is an error
func ErrNoInputs(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInputsOutputs, "no inputs to send transaction")
}

// ErrNoOutputs is an error
func ErrNoOutputs(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInputsOutputs, "no outputs to send transaction")
}

// ErrInputOutputMismatch is an error
func ErrInputOutputMismatch(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidInputsOutputs, "sum inputs != sum outputs")
}

// ErrSendDisabled is an error
func ErrSendDisabled(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSendDisabled, "send transactions are currently disabled")
}

// ErrNoValidatorTopup is an error for validator topup
func ErrNoValidatorTopup(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNoValidatorTopup, "no validator topup")
}

// ErrNoBalanceToWithdraw is an error for validator topup withdraw
func ErrNoBalanceToWithdraw(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNoBalanceToWithdraw, "No balance to withdraw")
}
