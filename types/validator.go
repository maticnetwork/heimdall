package types

import (
	"bytes"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
)

// Validator heimdall validator
type Validator struct {
	ID          ValidatorID     `json:"ID"`
	StartEpoch  uint64          `json:"startEpoch"`
	EndEpoch    uint64          `json:"endEpoch"`
	Nonce       uint64          `json:"nonce"`
	VotingPower int64           `json:"power"` // TODO add 10^-18 here so that we dont overflow easily
	PubKey      PubKey          `json:"pubKey"`
	Signer      HeimdallAddress `json:"signer"`
	LastUpdated string          `json:"last_updated"`

	Jailed           bool  `json:"jailed"`
	ProposerPriority int64 `json:"accum"`
}

// NewValidator func creates a new validator,
// the HeimdallAddress field is generated using Address i.e. [20]byte
func NewValidator(
	id ValidatorID,
	startEpoch uint64,
	endEpoch uint64,
	nonce uint64,
	power int64,
	pubKey PubKey,
	signer HeimdallAddress,
) *Validator {
	return &Validator{
		ID:          id,
		StartEpoch:  startEpoch,
		EndEpoch:    endEpoch,
		Nonce:       nonce,
		VotingPower: power,
		PubKey:      pubKey,
		Signer:      signer,
	}
}

// SortValidatorByAddress sorts a slice of validators by address
// to sort it we compare the values of the Signer(HeimdallAddress i.e. [20]byte)
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
	if !v.Jailed && v.StartEpoch <= currentEpoch && (v.EndEpoch == 0 || v.EndEpoch > currentEpoch) && v.VotingPower > 0 {
		return true
	}

	return false
}

// Validates validator
func (v *Validator) ValidateBasic() bool {
	if bytes.Equal(v.PubKey.Bytes(), ZeroPubKey.Bytes()) {
		return false
	}
	if bytes.Equal(v.Signer.Bytes(), []byte("")) {
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

// Returns the one with higher ProposerPriority.
func (v *Validator) CompareProposerPriority(other *Validator) *Validator {
	if v == nil {
		return other
	}
	switch {
	case v.ProposerPriority > other.ProposerPriority:
		return v
	case v.ProposerPriority < other.ProposerPriority:
		return other
	default:
		result := bytes.Compare(v.Signer.Bytes(), other.Signer.Bytes())
		switch {
		case result < 0:
			return v
		case result > 0:
			return other
		default:
			panic("Cannot compare identical validators")
		}
	}
}

func (v *Validator) String() string {
	if v == nil {
		return "nil-Validator"
	}
	return fmt.Sprintf("Validator{%v %v %v VP:%v A:%v}",
		v.ID,
		v.Signer,
		v.PubKey,
		v.VotingPower,
		v.ProposerPriority)
}

// ValidatorListString returns a prettified validator list for logging purposes.
func ValidatorListString(vals []*Validator) string {
	chunks := make([]string, len(vals))
	for i, val := range vals {
		chunks[i] = fmt.Sprintf("%s:%d", val.Signer, val.VotingPower)
	}

	return strings.Join(chunks, ",")
}

// Bytes computes the unique encoding of a validator with a given voting power.
// These are the bytes that gets hashed in consensus. It excludes address
// as its redundant with the pubkey. This also excludes ProposerPriority
// which changes every round.
func (v *Validator) Bytes() []byte {
	result := make([]byte, 64)
	copy(result[12:], v.Signer.Bytes())
	copy(result[32:], new(big.Int).SetInt64(v.VotingPower).Bytes())
	return result
}

// UpdatedAt returns block number of last validator update
func (v *Validator) UpdatedAt() string {
	return v.LastUpdated
}

// MinimalVal returns block number of last validator update
func (v *Validator) MinimalVal() MinimalVal {
	return MinimalVal{
		ID:          v.ID,
		VotingPower: uint64(v.VotingPower),
		Signer:      v.Signer,
	}
}

// --------

// ValidatorID  validator ID and helper functions
type ValidatorID uint64

// NewValidatorID generate new validator ID
func NewValidatorID(id uint64) ValidatorID {
	return ValidatorID(id)
}

// Bytes get bytes of validatorID
func (valID ValidatorID) Bytes() []byte {
	return []byte(strconv.FormatUint(valID.Uint64(), 10))
}

// Int converts validator ID to int
func (valID ValidatorID) Int() int {
	return int(valID)
}

// Uint64 converts validator ID to int
func (valID ValidatorID) Uint64() uint64 {
	return uint64(valID)
}

// Uint64 converts validator ID to int
func (valID ValidatorID) String() string {
	return strconv.FormatUint(valID.Uint64(), 10)
}

// --------

// MinimalVal is the minimal validator representation
// Used to send validator information to bor validator contract
type MinimalVal struct {
	ID          ValidatorID     `json:"ID"`
	VotingPower uint64          `json:"power"` // TODO add 10^-18 here so that we dont overflow easily
	Signer      HeimdallAddress `json:"signer"`
}

// SortMinimalValByAddress sorts validators
func SortMinimalValByAddress(a []MinimalVal) []MinimalVal {
	sort.Slice(a, func(i, j int) bool {
		return bytes.Compare(a[i].Signer.Bytes(), a[j].Signer.Bytes()) < 0
	})
	return a
}

// ValToMinVal converts array of validators to minimal validators
func ValToMinVal(vals []Validator) (minVals []MinimalVal) {
	for _, val := range vals {
		minVals = append(minVals, val.MinimalVal())
	}
	return
}
