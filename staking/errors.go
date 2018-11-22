package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

// TODO come up with better status numbers
const (
	DefaultCodespace       sdk.CodespaceType = 2
	CodeOldValidator       CodeType          = 2500
	CodeNoValidator        CodeType          = 3500
	CodeValSignerMismatch  CodeType          = 4500
	CodeValidatorExitDeny  CodeType          = 5500
	CodeValAlreadyUnbonded CodeType          = 6500
)

func ErrOldValidator(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeOldValidator, "Start Epoch behind Current Epoch")
}

func ErrNoValidator(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeNoValidator, "Validator information not found")
}

func ErrValSignerMismatch(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeValSignerMismatch, "Signer Address doesnt match pubkey address")
}

func ErrValIsCurrentVal(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeValidatorExitDeny, "Validator is locked in till deactivation epoch, exit denied")
}

func ErrValUnbonded(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeValAlreadyUnbonded, "Validator already unbonded , cannot exit")
}

func codeToDefaultMsg(code CodeType) string {
	switch code {
	case CodeOldValidator:
		return "Old validator"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

func msgOrDefaultMsg(msg string, code CodeType) string {
	if msg != "" {
		return msg
	}
	return codeToDefaultMsg(code)
}

func newError(codespace sdk.CodespaceType, code CodeType, msg string) sdk.Error {
	msg = msgOrDefaultMsg(msg, code)
	return sdk.NewError(codespace, code, msg)
}
