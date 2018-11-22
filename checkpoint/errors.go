package checkpoint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type CodeType = sdk.CodeType

const (
	DefaultCodespace      sdk.CodespaceType = 1
	CodeInvalidBlockinput CodeType          = 1500
	CodeInvalidACK        CodeType          = 1600
	CodeNoACK             CodeType          = 1700
)

func ErrBadBlockDetails(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeInvalidBlockinput, "Checkpoint is not valid!")
}

func ErrBadAck(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeInvalidACK, "Ack Not Valid")
}

func ErrNoACK(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeNoACK, "Checkpoint Already Exists In Buffer, ACK expected")
}

func codeToDefaultMsg(code CodeType) string {
	switch code {
	case CodeInvalidBlockinput:
		return "Invalid Block Input"
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
