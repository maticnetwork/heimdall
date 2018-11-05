package helper

import (
	abci "github.com/tendermint/tendermint/abci/types"
)

func ValidatorsToString(validators []abci.Validator) (validatorStr []string) {
	for index, validator := range validators {
		validatorStr[index] = validator.String()
	}
	return
}
