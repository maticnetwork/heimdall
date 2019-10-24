package types

import (
	"math/big"

	hmTyps "github.com/maticnetwork/heimdall/types"
)

// ValidatorSlashParams defines the params for slashing a validator
type ValidatorSlashParams struct {
	ValID       hmTyps.ValidatorID
	SlashAmount *big.Int
}

// NewValidatorSlashParams creates a new instance of ValidatorSlashParams.
func NewValidatorSlashParams(validatorID hmTyps.ValidatorID, amountToSlash *big.Int) ValidatorSlashParams {
	return ValidatorSlashParams{ValID: validatorID, SlashAmount: amountToSlash}
}
