package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/types"
)

type Validator struct {
	Address    common.Address
	StartEpoch int64
	EndEpoch   int64
	Pubkey     crypto.PubKey
	Power      int64 // aka Amount
	Signer     common.Address
}

func (validator Validator) IsCurrentValidator(ACKCount int) bool {
	// check if validator is current validator
	if validator.StartEpoch >= int64(ACKCount) && validator.EndEpoch <= int64(ACKCount) {
		return true
	}
	return false
}

// create empty validator without pubkey
func CreateEmptyValidator() Validator {
	validator := Validator{
		Address:    common.HexToAddress(""),
		StartEpoch: int64(0),
		EndEpoch:   int64(0),
		Power:      int64(0),
		Signer:     common.HexToAddress(""),
	}
	return validator
}

func CreateValidatorWithAddr(addr common.Address) Validator {
	validator := Validator{
		Address:    addr,
		StartEpoch: int64(0),
		EndEpoch:   int64(0),
		Power:      int64(0),
		Signer:     addr,
	}
	return validator
}

// todo add marshall and unmarshall methods here

// todo add human readable string

// convert heimdall validator to Tendermint validator
func (validator Validator) ToTmValidator() types.Validator {
	return types.Validator{
		Address:     validator.Address.Bytes(),
		PubKey:      validator.Pubkey,
		VotingPower: validator.Power,
	}
}
