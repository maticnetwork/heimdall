package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/checkpoint module sentinel errors
var (
	ErrSample                   = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrNoACK                    = sdkerrors.Register(ModuleName, 1503, "Checkpoint Already Exists In Buffer, ACK expected")
	ErrOldCheckpoint            = sdkerrors.Register(ModuleName, 1509, "Checkpoint already received for given start and end block")
	ErrDisCountinuousCheckpoint = sdkerrors.Register(ModuleName, 1510, "Checkpoint not in countinuity")
	ErrNoCheckpointFound        = sdkerrors.Register(ModuleName, 1508, "Checkpoint Not Found")
	ErrBadBlockDetails          = sdkerrors.Register(ModuleName, 1501, "Wrong roothash for given start and end block numbers")
	ErrInvalidMsg               = sdkerrors.Register(ModuleName, 1400, "Invalid Message")
	ErrBadAck                   = sdkerrors.Register(ModuleName, 1502, "Ack Not Valid")
	ErrInvalidNoACK             = sdkerrors.Register(ModuleName, 1505, "Invalid No ACK -- Waiting for last checkpoint ACK")
	ErrTooManyNoACK             = sdkerrors.Register(ModuleName, 1506, "Too many no-acks")
)
