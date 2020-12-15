package common

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

//ModuleName Defination
var ModuleName = "common_errors"

//custom error definations
var (
	ErrEmptyValidatorAddr      = sdkerrors.Register(ModuleName, 2, "empty validator address")
	ErrInvalidMsg              = sdkerrors.Register(ModuleName, 1400, "Invalid Message")
	ErrOldTx                   = sdkerrors.Register(ModuleName, 1401, "Old txhash not allowed")
	ErrBadProposerDetails      = sdkerrors.Register(ModuleName, 1500, "Proper is not valid")
	ErrWaitForConfirmation     = sdkerrors.Register(ModuleName, 2510, "Please wait for confirmation time before sending transaction")
	ErrValSignerPubKeyMismatch = sdkerrors.Register(ModuleName, 2511, "Signer Pubkey mismatch between event and msg")
	ErrValSignerMismatch       = sdkerrors.Register(ModuleName, 2502, "Signer Address doesnt match pubkey address")
	ErrValidatorAlreadyJoined  = sdkerrors.Register(ModuleName, 2507, "Validator already joined")
	ErrValidatorSave           = sdkerrors.Register(ModuleName, 2506, "Cannot save validator")
	ErrNoValidator             = sdkerrors.Register(ModuleName, 2501, "Validator information not found")
	ErrNonce                   = sdkerrors.Register(ModuleName, 2514, "Incorrect validator nonce")
	ErrNoSignerChange          = sdkerrors.Register(ModuleName, 2513, "New signer same as old signer")
	ErrValUnbonded             = sdkerrors.Register(ModuleName, 2504, "Validator already unbonded , cannot exit")
)
