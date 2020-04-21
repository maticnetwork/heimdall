package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
)

// ValidatorSlashingInfo - contains ID, slashingAmount, isJailed
type ValidatorSlashingInfo struct {
	ID            ValidatorID `json:"ID"`
	SlashedAmount string      `json:"SlashedAmount"` // string representation of big.Int
	IsJailed      bool        `json:"IsJailed"`
}

func NewValidatorSlashingInfo(id ValidatorID, slashedAmount string, isJailed bool) ValidatorSlashingInfo {

	return ValidatorSlashingInfo{
		ID:            id,
		SlashedAmount: slashedAmount,
		IsJailed:      isJailed,
	}
}

func (v ValidatorSlashingInfo) String() string {
	return fmt.Sprintf(`Validator Signing Info:
	ID:               %d
	SlashedAmount:    %s
	IsJailed:         %v`,
		v.ID, v.SlashedAmount, v.IsJailed)
}

// amino marshall validator slashing info
func MarshallValSlashingInfo(cdc *codec.Codec, valSlashingInfo ValidatorSlashingInfo) (bz []byte, err error) {
	bz, err = cdc.MarshalBinaryBare(valSlashingInfo)
	if err != nil {
		return bz, err
	}
	return bz, nil
}

// amono unmarshall validator slashing info
func UnmarshallValSlashingInfo(cdc *codec.Codec, value []byte) (ValidatorSlashingInfo, error) {
	var valSlashingInfo ValidatorSlashingInfo
	// unmarshall validator and return
	err := cdc.UnmarshalBinaryBare(value, &valSlashingInfo)
	if err != nil {
		return valSlashingInfo, err
	}
	return valSlashingInfo, nil
}
