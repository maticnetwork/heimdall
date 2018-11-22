package types

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/tendermint/tendermint/crypto"
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
