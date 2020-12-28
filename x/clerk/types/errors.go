package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/clerk module sentinel errors
var (
	ErrSample                    = sdkerrors.Register(ModuleName, 1100, "sample error")
	CodeEventRecordAlreadySynced = sdkerrors.Register(ModuleName, 5400, "Event record already synced")
	CodeEventRecordInvalid       = sdkerrors.Register(ModuleName, 5401, "Event record is invalid")
	CodeEventRecordUpdate        = sdkerrors.Register(ModuleName, 5402, "Event record update error")
)
