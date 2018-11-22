package staking

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

const (
	DefaultCodespace sdk.CodespaceType = 2
	CodeOldValidator CodeType          = 2500
	CodeNoValidator  CodeType          = 2500
)

func ErrOldValidator(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeOldValidator, "Start Epoch behind Current Epoch")
}

func ErrNoValidator(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeNoValidator, "Valdiator information not found")
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
