package types

import (
	"github.com/cosmos/cosmos-sdk/codec"

	hmTypes "github.com/maticnetwork/heimdall/types"
)

//
// Validator topup
//

// ValidatorTopup heimdall validator
type ValidatorTopup struct {
	ID          hmTypes.ValidatorID `json:"id"`
	TotalTopups hmTypes.Coins       `json:"total_topups"`
}

// MarshallValidatorTopup marshall validator topup
func MarshallValidatorTopup(cdc *codec.Codec, validatorTopup ValidatorTopup) (bz []byte, err error) {
	bz, err = cdc.MarshalBinaryBare(validatorTopup)
	if err != nil {
		return bz, err
	}
	return bz, nil
}

// UnmarshallValidatorTopup unmarshall validator topup
func UnmarshallValidatorTopup(cdc *codec.Codec, value []byte) (ValidatorTopup, error) {
	var validatorTopup ValidatorTopup
	// unmarshall validator and return
	err := cdc.UnmarshalBinaryBare(value, &validatorTopup)
	if err != nil {
		return validatorTopup, err
	}
	return validatorTopup, nil
}

// Copy creates a new copy of the ValidatorTopup so we can mutate last updated at
// Panics if the validator is nil.
func (v *ValidatorTopup) Copy() *ValidatorTopup {
	vCopy := *v
	return &vCopy
}
