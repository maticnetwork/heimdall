package checkpoint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

const (
	DefaultCodespace      sdk.CodespaceType = 1
	CodeInvalidBlockinput CodeType          = 1500
)

func ErrBadBlockDetails(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeInvalidBlockinput, "Checkpoint Not Valid!")
}

func codeToDefaultMsg(code CodeType) string {
	switch code {
	case CodeInvalidBlockinput:
		return "Invalid Validator"
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
