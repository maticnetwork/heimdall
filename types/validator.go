package types

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"
	"strconv"

	"github.com/cosmos/cosmos-sdk/codec"
)

// Validator heimdall validator
type Validator struct {
	ID          ValidatorID     `json:"ID"`
	StartEpoch  uint64          `json:"startEpoch"`
	EndEpoch    uint64          `json:"endEpoch"`
	Power       uint64          `json:"power"` // TODO add 10^-18 here so that we dont overflow easily
	PubKey      PubKey          `json:"pubKey"`
	Signer      HeimdallAddress `json:"signer"`
	LastUpdated *big.Int        `json:"last_updated"`

	Accum int64 `json:"accum"`
}

// SortValidatorByAddress sorts a slice of validators by address
func SortValidatorByAddress(a []Validator) []Validator {
	sort.Slice(a, func(i, j int) bool {
		return bytes.Compare(a[i].Signer.Bytes(), a[j].Signer.Bytes()) < 0
	})
	return a
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

// Validates validator
func (v *Validator) ValidateBasic() bool {
	if v.StartEpoch < 0 || v.EndEpoch < 0 {
		return false
	}
	if bytes.Equal(v.PubKey.Bytes(), ZeroPubKey.Bytes()) {
		return false
	}
	if bytes.Equal(v.Signer.Bytes(), []byte("")) {
		return false
	}
	if v.ID < 0 {
		return false
	}
	return true
}

// amino marshall validator
func MarshallValidator(cdc *codec.Codec, validator Validator) (bz []byte, err error) {
	bz, err = cdc.MarshalBinaryBare(validator)
	if err != nil {
		return bz, err
	}
	return bz, nil
}

// amono unmarshall validator
func UnmarshallValidator(cdc *codec.Codec, value []byte) (Validator, error) {
	var validator Validator
	// unmarshall validator and return
	err := cdc.UnmarshalBinaryBare(value, &validator)
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
		result := bytes.Compare(v.Signer.Bytes(), other.Signer.Bytes())
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

	return fmt.Sprintf("Validator{%v::%v P:%v Start:%v End:%v A:%v}",
		v.ID,
		v.Signer.String(),
		v.Power,
		v.StartEpoch,
		v.EndEpoch,
		v.Accum,
	)
}

// returns block number of last validator update
func (v *Validator) UpdatedAt() *big.Int {
	return v.LastUpdated
}

// returns block number of last validator update
func (v *Validator) MinimalVal() MinimalVal {
	return MinimalVal{
		ID:     v.ID,
		Power:  v.Power,
		Signer: v.Signer,
	}
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

// --------

// validator ID and helper functions
type ValidatorID uint64

// generate new validator ID
func NewValidatorID(id uint64) ValidatorID {
	return ValidatorID(id)
}

// get bytes of validatorID
func (valID ValidatorID) Bytes() []byte {
	return []byte(strconv.Itoa(valID.Int()))
}

// convert validator ID to int
func (valID ValidatorID) Int() int {
	return int(valID)
}

// --------

// MinimalVal is the minimal validator representation
// Used to send validator information to bor validator contract
type MinimalVal struct {
	ID     ValidatorID     `json:"ID"`
	Power  uint64          `json:"power"` // TODO add 10^-18 here so that we dont overflow easily
	Signer HeimdallAddress `json:"signer"`
}

// ValToMinVal converts array of validators to minimal validators
func ValToMinVal(vals []Validator) (minVals []MinimalVal) {
	for _, val := range vals {
		minVals = append(minVals, val.MinimalVal())
	}
	return
}
