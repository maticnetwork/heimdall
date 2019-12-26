package utils

import "errors"

var (
	ErrInvalidSigner        = errors.New("tx intended signer does not match the given signer")
	ErrInvalidGasAdjustment = errors.New("invalid gas adjustment")
)
