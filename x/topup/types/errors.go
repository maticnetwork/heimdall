package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/topup module sentinel errors
var (
	ErrSendDisabled            = sdkerrors.Register(ModuleName, 101, "Send transactions are currently disabled")
	ErrInvalidInputsOutputs    = sdkerrors.Register(ModuleName, 102, "No inputs/outputs or sum inputs != sum outputs to send transaction")
	ErrNoValidatorTopup        = sdkerrors.Register(ModuleName, 103, "No validator topup")
	ErrNoBalanceToWithdraw     = sdkerrors.Register(ModuleName, 104, "No balance to withdraw")
	ErrSetFeeBalanceZero       = sdkerrors.Register(ModuleName, 7001, "Error setting fee balance to zero")
	ErrAddFeeToDividendAccount = sdkerrors.Register(ModuleName, 7002, "Error adding fee to dividend account")
)
