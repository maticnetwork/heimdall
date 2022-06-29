package types

import (
	"fmt"
	"sort"

	"github.com/cosmos/cosmos-sdk/codec"
)

// ValidatorSlashingInfo - contains ID, slashingAmount, isJailed
type ValidatorSlashingInfo struct {
	ID            ValidatorID `json:"ID"`
	SlashedAmount uint64      `json:"SlashedAmount"`
	IsJailed      bool        `json:"IsJailed"`
}

func NewValidatorSlashingInfo(id ValidatorID, slashedAmount uint64, isJailed bool) ValidatorSlashingInfo {
	return ValidatorSlashingInfo{
		ID:            id,
		SlashedAmount: slashedAmount,
		IsJailed:      isJailed,
	}
}

func (v ValidatorSlashingInfo) String() string {
	return fmt.Sprintf(`Validator Slashing Info:
	ID:               %d
	SlashedAmount:    %d
	IsJailed:         %v`,
		v.ID, v.SlashedAmount, v.IsJailed)
}

// SortValidatorSlashingInfoByID - Sorts ValidatorSlashingInfo By ID
func SortValidatorSlashingInfoByID(slashingInfos []*ValidatorSlashingInfo) []*ValidatorSlashingInfo {
	sort.Slice(slashingInfos, func(i, j int) bool { return slashingInfos[i].ID < slashingInfos[j].ID })
	return slashingInfos
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
	if err := cdc.UnmarshalBinaryBare(value, &valSlashingInfo); err != nil {
		return valSlashingInfo, err
	}

	return valSlashingInfo, nil
}
