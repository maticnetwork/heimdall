package types

import (
	"bytes"
	"fmt"
	"math/big"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/ethereum/go-ethereum/common"
)

// Validator heimdall validator
type Validator struct {
	Address    common.Address `json:"address"`
	StartEpoch uint64         `json:"startEpoch"`
	EndEpoch   uint64         `json:"endEpoch"`
	Power      uint64         `json:"power"` // TODO add 10^-18 here so that we dont overflow easily
	PubKey     PubKey         `json:"pubKey"`
	Signer     common.Address `json:"signer"`

	Accum int64 `json:"accum"`
}

// IsCurrentValidator checks if validator is in current validator set
func (v *Validator) IsCurrentValidator(ackCount uint64) bool {
	// current epoch will be ack count + 1
	currentEpoch := ackCount + 1

	// validator hasnt initialised unstake
	if v.StartEpoch <= currentEpoch && (v.EndEpoch == 0 || v.EndEpoch >= currentEpoch) && v.Power > 0 {
		return true
	}

	return false
}

func MarshallValidator(cdc *codec.Codec, validator Validator) (bz []byte, err error) {
	bz, err = cdc.MarshalBinary(validator)
	if err != nil {
		return bz, err
	}
	return bz, nil
}

func UnmarshallValidator(cdc *codec.Codec, value []byte) (Validator, error) {
	var validator Validator
	// unmarshall validator and return
	err := cdc.UnmarshalBinary(value, &validator)
	if err != nil {
		return validator, err
	}
	return validator, nil
}

// Copy creates a new copy of the validator so we can mutate accum.
// Panics if the validator is nil.
func (v *Validator) Copy() *Validator {
	vCopy := *v
	return &vCopy
}

// CompareAccum returns the one with higher Accum.
func (v *Validator) CompareAccum(other *Validator) *Validator {
	if v == nil {
		return other
	}
	if v.Accum > other.Accum {
		return v
	} else if v.Accum < other.Accum {
		return other
	} else {
		result := bytes.Compare(v.Address.Bytes(), other.Address.Bytes())
		if result < 0 {
			return v
		} else if result > 0 {
			return other
		} else {
			return nil
		}
	}
}

func (v *Validator) String() string {
	if v == nil {
		return "nil-Validator"
	}

	return fmt.Sprintf("Validator{%v::%v P:%v}",
		v.Address.String(),
		v.Signer.String(),
		v.Power,
	)
}

// GetValidatorPower converts amount to power
func GetValidatorPower(amount string) uint64 {
	result := big.NewInt(0)
	result.SetString(amount, 10)
	if len(amount) >= 18 {
		t, _ := big.NewInt(0).SetString("1000000000000000000", 10)
		result.Div(result, t)
	}
	return result.Uint64()
}
