package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
)

// ValidatorAccount contains Rewards, Slashed Amount
type ValidatorAccount struct {
	ID            ValidatorID `json:"ID"`
	RewardAmount  []byte      `json:"rewardAmount"`
	SlashedAmount []byte      `json:"slashedAmount"`
}

func (va *ValidatorAccount) String() string {
	if va == nil {
		return "nil-ValidatorAccount"
	}

	return fmt.Sprintf("ValidatorAccount{%v %v %v}",
		va.ID,
		va.RewardAmount,
		va.SlashedAmount)
}

// MarshallValidatorAccount - amino Marshall ValidatorAccount
func MarshallValidatorAccount(cdc *codec.Codec, validatorAccount ValidatorAccount) (bz []byte, err error) {
	bz, err = cdc.MarshalBinaryBare(validatorAccount)
	if err != nil {
		return bz, err
	}

	return bz, nil
}

// UnMarshallValidatorAccount - amino Unmarshall ValidatorAccount
func UnMarshallValidatorAccount(cdc *codec.Codec, value []byte) (ValidatorAccount, error) {

	var validatorAccount ValidatorAccount
	err := cdc.UnmarshalBinaryBare(value, &validatorAccount)
	if err != nil {
		return validatorAccount, err
	}
	return validatorAccount, nil
}
