package common

import (
	"fmt"
	"strconv"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/maticnetwork/heimdall/types"
)

type CodeType = sdk.CodeType

const (
	DefaultCodespace sdk.CodespaceType = "1"

	CodeInvalidMsg CodeType = 1400

	CodeInvalidProposerInput     CodeType = 1500
	CodeInvalidBlockInput        CodeType = 1501
	CodeInvalidACK               CodeType = 1502
	CodeNoACK                    CodeType = 1503
	CodeBadTimeStamp             CodeType = 1504
	CodeInvalidNoACK             CodeType = 1505
	CodeTooManyNoAck             CodeType = 1506
	CodeLowBal                   CodeType = 1507
	CodeNoCheckpoint             CodeType = 1508
	CodeOldCheckpoint            CodeType = 1509
	CodeDisCountinuousCheckpoint CodeType = 1510
	CodeNoCheckpointBuffer       CodeType = 1511

	CodeOldValidator        CodeType = 2500
	CodeNoValidator         CodeType = 2501
	CodeValSignerMismatch   CodeType = 2502
	CodeValidatorExitDeny   CodeType = 2503
	CodeValAlreadyUnbonded  CodeType = 2504
	CodeSignerSynced        CodeType = 2505
	CodeValSave             CodeType = 2506
	CodeValAlreadyJoined    CodeType = 2507
	CodeSignerUpdateError   CodeType = 2508
	CodeNoConn              CodeType = 2509
	CodeWaitFrConfirmation  CodeType = 2510
	CodeValPubkeyMismatch   CodeType = 2511
	CodeErrDecodeEvent      CodeType = 2512
	CodeNoSignerChangeError CodeType = 2513

	CodeSpanNotCountinuous CodeType = 3501
	CodeUnableToFreezeSet  CodeType = 3502
	CodeSpanNotFound       CodeType = 3503
	CodeValSetMisMatch     CodeType = 3504
	CodeProducerMisMatch   CodeType = 3505
	CodeInvalidBorChainID  CodeType = 3506

	CodeFetchCheckpointSigners       CodeType = 4501
	CodeErrComputeGenesisAccountRoot CodeType = 4503
	CodeAccountRootMismatch          CodeType = 4504

	CodeErrAccountRootHash     CodeType = 4505
	CodeErrSetCheckpointBuffer CodeType = 4506
	CodeErrAddCheckpoint       CodeType = 4507

	CodeInvalidReceipt         CodeType = 5501
	CodeSideTxValidationFailed CodeType = 5502
)

// -------- Invalid msg

func ErrInvalidMsg(codespace sdk.CodespaceType, format string, args ...interface{}) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidMsg, format, args...)
}

// -------- Checkpoint Errors

func ErrBadProposerDetails(codespace sdk.CodespaceType, proposer types.HeimdallAddress) sdk.Error {
	return newError(codespace, CodeInvalidProposerInput, fmt.Sprintf("Proposer is not valid, current proposer is %v", proposer.String()))
}

func ErrBadBlockDetails(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeInvalidBlockInput, "Wrong roothash for given start and end block numbers")
}

func ErrSetCheckpointBuffer(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeErrSetCheckpointBuffer, "Account Root Hash not added to Checkpoint Buffer")
}

func ErrAddCheckpoint(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeErrAddCheckpoint, "Err in adding checkpoint to header blocks")
}

func ErrBadAccountRootHash(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeErrAccountRootHash, "Wrong roothash for given dividend accounts")
}

func ErrBadAck(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeInvalidACK, "Ack Not Valid")
}

func ErrOldCheckpoint(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeOldCheckpoint, "Checkpoint already received for given start and end block")
}

func ErrDisCountinuousCheckpoint(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeDisCountinuousCheckpoint, "Checkpoint not in countinuity")
}

func ErrNoACK(codespace sdk.CodespaceType, expiresAt uint64) sdk.Error {
	return newError(codespace, CodeNoACK, fmt.Sprintf("Checkpoint Already Exists In Buffer, ACK expected, expires at %s", strconv.FormatUint(expiresAt, 10)))
}

func ErrNoConn(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeNoConn, "Unable to connect to chain")
}

func ErrWaitForConfirmation(codespace sdk.CodespaceType, txConfirmationTime time.Duration) sdk.Error {
	return newError(codespace, CodeWaitFrConfirmation, fmt.Sprintf("Please wait for %v confirmation time before sending transaction", txConfirmationTime))
}

func ErrNoCheckpointFound(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeNoCheckpoint, "Checkpoint Not Found")
}

func ErrNoCheckpointBufferFound(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeNoCheckpointBuffer, "Checkpoint buffer not found")
}

func ErrInvalidNoACK(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeInvalidNoACK, "Invalid No ACK -- Waiting for last checkpoint ACK")
}

func ErrTooManyNoACK(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeTooManyNoAck, "Too many no-acks")
}

func ErrBadTimeStamp(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeBadTimeStamp, "Invalid time stamp. It must be in near past.")
}

// ----------- Staking Errors

func ErrOldValidator(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeOldValidator, "Start Epoch behind Current Epoch")
}

func ErrNoValidator(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeNoValidator, "Validator information not found")
}

func ErrValSignerPubKeyMismatch(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeValPubkeyMismatch, "Signer Pubkey mismatch between event and msg")
}

func ErrValSignerMismatch(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeValSignerMismatch, "Signer Address doesnt match pubkey address")
}

func ErrValIsNotCurrentVal(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeValidatorExitDeny, "Validator is not in validator set, exit not possible")
}

func ErrValUnbonded(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeValAlreadyUnbonded, "Validator already unbonded , cannot exit")
}

func ErrSignerUpdateError(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeSignerUpdateError, "Signer update error")
}

func ErrNoSignerChange(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeNoSignerChangeError, "New signer same as old signer")
}

func ErrOldTx(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeSignerUpdateError, "Old txhash not allowed")
}

func ErrValidatorAlreadySynced(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeSignerSynced, "No signer update found, invalid message")
}

func ErrValidatorSave(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeValSave, "Cannot save validator")
}

func ErrValidatorNotDeactivated(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeValSave, "Validator Not Deactivated")
}

func ErrValidatorAlreadyJoined(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeValAlreadyJoined, "Validator already joined")
}

// Bor Errors --------------------------------

func ErrInvalidBorChainID(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeInvalidBorChainID, "Invalid Bor chain id")
}

func ErrSpanNotInCountinuity(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeSpanNotCountinuous, "Span not countinuous")
}

func ErrSpanNotFound(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeSpanNotFound, "Span not found")
}

func ErrUnableToFreezeValSet(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeUnableToFreezeSet, "Unable to freeze validator set for next span")
}

func ErrValSetMisMatch(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeValSetMisMatch, "Validator set mismatch")
}

func ErrProducerMisMatch(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeProducerMisMatch, "Producer set mismatch")
}

//
// Side-tx errors
//

// ErrorSideTx represents side-tx error
func ErrorSideTx(codespace sdk.CodespaceType, code CodeType) (res abci.ResponseDeliverSideTx) {
	res.Code = uint32(code)
	res.Codespace = string(codespace)
	res.Result = abci.SideTxResultType_Skip // skip side-tx vote in-case of error
	return
}

func ErrSideTxValidation(codespace sdk.CodespaceType) sdk.Error {
	return newError(codespace, CodeSideTxValidationFailed, "External call majority validation failed. ")
}

//
// Private methods
//

func CodeToDefaultMsg(code CodeType) string {
	switch code {
	// case CodeInvalidBlockInput:
	// 	return "Invalid Block Input"

	case CodeInvalidMsg:
		return "Invalid Message"

	case CodeInvalidProposerInput:
		return "Proposer is not valid"
	case CodeInvalidBlockInput:
		return "Wrong roothash for given start and end block numbers"
	case CodeInvalidACK:
		return "Ack Not Valid"
	case CodeNoACK:
		return "Checkpoint Already Exists In Buffer, ACK expected"
	case CodeBadTimeStamp:
		return "Invalid time stamp. It must be in near past."
	case CodeInvalidNoACK:
		return "Invalid No ACK -- Waiting for last checkpoint ACK"
	case CodeTooManyNoAck:
		return "Too many no-acks"
	case CodeLowBal:
		return "Insufficient balance"
	case CodeNoCheckpoint:
		return "Checkpoint Not Found"
	case CodeOldCheckpoint:
		return "Checkpoint already received for given start and end block"
	case CodeDisCountinuousCheckpoint:
		return "Checkpoint not in countinuity"
	case CodeNoCheckpointBuffer:
		return "Checkpoint buffer Not Found"

	case CodeOldValidator:
		return "Start Epoch behind Current Epoch"
	case CodeNoValidator:
		return "Validator information not found"
	case CodeValSignerMismatch:
		return "Signer Address doesnt match pubkey address"
	case CodeValidatorExitDeny:
		return "Validator is not in validator set, exit not possible"
	case CodeValAlreadyUnbonded:
		return "Validator already unbonded , cannot exit"
	case CodeSignerSynced:
		return "No signer update found, invalid message"
	case CodeValSave:
		return "Cannot save validator"
	case CodeValAlreadyJoined:
		return "Validator already joined"
	case CodeSignerUpdateError:
		return "Signer update error"
	case CodeNoConn:
		return "Unable to connect to chain"
	case CodeWaitFrConfirmation:
		return "wait for confirmation time before sending transaction"
	case CodeValPubkeyMismatch:
		return "Signer Pubkey mismatch between event and msg"
	case CodeSpanNotCountinuous:
		return "Span not countinuous"
	case CodeUnableToFreezeSet:
		return "Unable to freeze validator set for next span"
	case CodeSpanNotFound:
		return "Span not found"
	case CodeValSetMisMatch:
		return "Validator set mismatch"
	case CodeProducerMisMatch:
		return "Producer set mismatch"
	case CodeInvalidBorChainID:
		return "Invalid Bor chain id"
	default:
		return sdk.CodeToDefaultMsg(code)
	}
}

func msgOrDefaultMsg(msg string, code CodeType) string {
	if msg != "" {
		return msg
	}
	return CodeToDefaultMsg(code)
}

func newError(codespace sdk.CodespaceType, code CodeType, msg string) sdk.Error {
	msg = msgOrDefaultMsg(msg, code)
	return sdk.NewError(codespace, code, msg)
}
