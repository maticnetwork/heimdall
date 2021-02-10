package common

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	tmprototypes "github.com/tendermint/tendermint/proto/tendermint/types"
)

//ModuleName Definition
var ModuleName = "common_errors"

//custom error definitions
var (
	ErrInvalidMsg              = sdkerrors.Register(ModuleName, 1400, "Invalid Message")
	ErrOldTx                   = sdkerrors.Register(ModuleName, 1401, "Old txhash not allowed")
	ErrEmptyValidatorAddr      = sdkerrors.Register(ModuleName, 1402, "Invalid validator address")
	ErrDecodeEvent             = sdkerrors.Register(ModuleName, 1403, "Event decoding error")
	ErrBadProposerDetails      = sdkerrors.Register(ModuleName, 1500, "Proposer is not valid")
	ErrWaitForConfirmation     = sdkerrors.Register(ModuleName, 2510, "Please wait for confirmation time before sending transaction")
	ErrValSignerPubKeyMismatch = sdkerrors.Register(ModuleName, 2511, "Signer Pubkey mismatch between event and msg")
	ErrValSignerMismatch       = sdkerrors.Register(ModuleName, 2512, "Signer Address doesnt match pubkey address")
	ErrValidatorAlreadyJoined  = sdkerrors.Register(ModuleName, 2513, "Validator already joined")
	ErrValidatorSave           = sdkerrors.Register(ModuleName, 2514, "Cannot save validator")
	ErrNoValidator             = sdkerrors.Register(ModuleName, 2515, "Validator information not found")
	ErrNonce                   = sdkerrors.Register(ModuleName, 2516, "Incorrect validator nonce")
	ErrNoSignerChange          = sdkerrors.Register(ModuleName, 2517, "New signer same as old signer")
	ErrValUnbonded             = sdkerrors.Register(ModuleName, 2518, "Validator already unbonded, cannot exit")
	ErrInvalidPower            = sdkerrors.Register(ModuleName, 2519, "Invalid amount for stake power")

	ErrInvalidBorChainID = sdkerrors.Register(ModuleName, 3506, "Invalid Bor chain id")

	ErrEventRecordAlreadySynced = sdkerrors.Register(ModuleName, 5400, "Event record already synced")
	ErrEventRecordInvalid       = sdkerrors.Register(ModuleName, 5401, "Event record is invalid")
	ErrEventUpdate              = sdkerrors.Register(ModuleName, 5402, "Event record update error")
	ErrSideTxValidation         = sdkerrors.Register(ModuleName, 5502, "External call majority validation failed")
	ErrValidatorSigningInfoSave = sdkerrors.Register(ModuleName, 6501, "Cannot save validator signing info")
	ErrSignerUpdateError        = sdkerrors.Register(ModuleName, 2508, "Signer update error")
	ErrValidatorNotDeactivated  = sdkerrors.Register(ModuleName, 6502, "Validator Not Deactivated")
	// TODO: Check if this is ok:
	ErrEmptyAddr = sdkerrors.Register(ModuleName, 7001, "Empty address")

	// Bor Errors --------------------------------
	ErrSpanNotInCountinuity = sdkerrors.Register(ModuleName, 3501, "Span not continuous")
	ErrInvalidSpanDuration  = sdkerrors.Register(ModuleName, 3507, "Wrong span duration")
	ErrSpanNotFound         = sdkerrors.Register(ModuleName, 3503, "Span not found")
	ErrUnableToFreezeValSet = sdkerrors.Register(ModuleName, 3502, "Unable to freeze validator set for next span")
	ErrValSetMisMatch       = sdkerrors.Register(ModuleName, 3504, "Validator set mismatch")
	ErrProducerMisMatch     = sdkerrors.Register(ModuleName, 3505, "Producer set mismatch")
)

// ErrorSideTx represents side-tx error
func ErrorSideTx(err *sdkerrors.Error) (res abci.ResponseDeliverSideTx) {
	res.Code = err.ABCICode()
	res.Codespace = err.Codespace()
	res.Result = tmprototypes.SideTxResultType_SKIP // skip side-tx vote in-case of error
	return
}
