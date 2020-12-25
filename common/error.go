package common

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	abci "github.com/tendermint/tendermint/abci/types"
	tmprototypes "github.com/tendermint/tendermint/proto/tendermint/types"
)

//ModuleName Definition
var ModuleName = "common_errors"

var (
	// CodeInvalidMsg error code
	CodeInvalidMsg uint32 = 1400
)

//custom error definitions
var (
	ErrEmptyValidatorAddr       = sdkerrors.Register(ModuleName, 2, "empty validator address")
	ErrInvalidMsg               = sdkerrors.Register(ModuleName, 1400, "Invalid Message")
	ErrOldTx                    = sdkerrors.Register(ModuleName, 1401, "Old txhash not allowed")
	ErrBadProposerDetails       = sdkerrors.Register(ModuleName, 1500, "Proper is not valid")
	ErrWaitForConfirmation      = sdkerrors.Register(ModuleName, 2510, "Please wait for confirmation time before sending transaction")
	ErrValSignerPubKeyMismatch  = sdkerrors.Register(ModuleName, 2511, "Signer Pubkey mismatch between event and msg")
	ErrValSignerMismatch        = sdkerrors.Register(ModuleName, 2502, "Signer Address doesnt match pubkey address")
	ErrValidatorAlreadyJoined   = sdkerrors.Register(ModuleName, 2507, "Validator already joined")
	ErrValidatorSave            = sdkerrors.Register(ModuleName, 2506, "Cannot save validator")
	ErrNoValidator              = sdkerrors.Register(ModuleName, 2501, "Validator information not found")
	ErrNonce                    = sdkerrors.Register(ModuleName, 2514, "Incorrect validator nonce")
	ErrNoSignerChange           = sdkerrors.Register(ModuleName, 2513, "New signer same as old signer")
	ErrValUnbonded              = sdkerrors.Register(ModuleName, 2504, "Validator already unbonded , cannot exit")
	ErrSideTxValidation         = sdkerrors.Register(ModuleName, 5502, "External call majority validation failed")
	ErrValidatorSigningInfoSave = sdkerrors.Register(ModuleName, 6501, "Cannot save validator signing info")
	ErrSignerUpdateError        = sdkerrors.Register(ModuleName, 2508, "Signer update error")
	ErrValidatorNotDeactivated  = sdkerrors.Register(ModuleName, 6502, "Validator Not Deactivated")
)

// ErrorSideTx represents side-tx error
func ErrorSideTx(code uint32) (res abci.ResponseDeliverSideTx) {
	res.Code = uint32(code)
	res.Codespace = string(ModuleName)
	res.Result = tmprototypes.SideTxResultType_SKIP // skip side-tx vote in-case of error
	return
}
