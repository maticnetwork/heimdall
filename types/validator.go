package types

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/types"
)

var _ crypto.PubKey = secp256k1.PubKeySecp256k1{}

// Validator heimdall validator
type Validator struct {
	Address    common.Address
	StartEpoch int64
	EndEpoch   int64
	Power      int64 // aka Amount
	PubKey     crypto.PubKey
	Signer     common.Address
}

// IsCurrentValidator checks if validator is in current validator set
func (v *Validator) IsCurrentValidator(ackCount int) bool {
	// validator hasnt initialised unstake
	if v.StartEpoch >= int64(ackCount) && v.EndEpoch == int64(0) {
		return true
	}

	// validator has initialised unstake but is unbonding period
	if v.StartEpoch >= int64(ackCount) && v.EndEpoch <= int64(ackCount) {
		return true
	}

	return false
}

func (v *Validator) String() string {
	if v == nil {
		return "nil-Validator"
	}

	return fmt.Sprintf("Validator{%v :: %v %v P:%v}",
		v.Address,
		v.Signer,
		v.PubKey,
		v.Power,
	)
}

// ToTmValidator converts heimdall validator to Tendermint validator
func (v *Validator) ToTmValidator() types.Validator {
	return types.Validator{
		Address:     v.Signer.Bytes(),
		PubKey:      v.PubKey,
		VotingPower: v.Power,
	}
}

// // create empty validator without pubkey
// func CreateEmptyValidator() Validator {
// 	validator := Validator{
// 		Address:    common.HexToAddress(""),
// 		StartEpoch: int64(0),
// 		EndEpoch:   int64(0),
// 		Power:      int64(0),
// 		Signer:     common.HexToAddress(""),
// 	}
// 	return validator
// }

// func CreateValidatorWithAddr(addr common.Address) Validator {
// 	validator := Validator{
// 		Address:    addr,
// 		StartEpoch: int64(0),
// 		EndEpoch:   int64(0),
// 		Power:      int64(0),
// 		Signer:     addr,
// 	}
// 	return validator
// }

// todo add marshall and unmarshall methods here
