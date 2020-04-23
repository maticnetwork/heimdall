package types

import (
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// GetSlashingInfoHash returns hash of latest slashing info
func GetSlashingInfoHash(valSlashingInfos []*hmTypes.ValidatorSlashingInfo) ([]byte, error) {
	slashInfoHash, err := GenerateInfoHash(valSlashingInfos)
	if err != nil {
		return nil, err
	}

	return slashInfoHash, nil
}

// GetAccountTree returns roothash of Validator Account State Tree
func GenerateInfoHash(slashingInfos []*hmTypes.ValidatorSlashingInfo) ([]byte, error) {
	// Sort the dividendAccounts by ID
	slashingInfos = hmTypes.SortValidatorSlashingInfoByID(slashingInfos)

	// TODO - add logic to generate hash
	return nil, nil
}
