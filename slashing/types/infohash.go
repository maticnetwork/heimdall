package types

import (
	"github.com/maticnetwork/bor/rlp"
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
	encodedSlashInfo, err := SortAndRLPEncodeSlashInfos(slashingInfos)

	if err != nil {
		return nil, err
	}

	// TODO slashing - Generate Hash of encodedData
	return encodedSlashInfo, nil
}

// SortAndRLPEncodeSlashInfos  - RLP encoded slashing infos
func SortAndRLPEncodeSlashInfos(slashingInfos []*hmTypes.ValidatorSlashingInfo) ([]byte, error) {

	// Sort the slashingInfos by ID
	slashingInfos = hmTypes.SortValidatorSlashingInfoByID(slashingInfos)

	// Encode slashInfos
	encodedSlashInfos, err := rlp.EncodeToBytes(slashingInfos)

	return encodedSlashInfos, err
}

func RLPDecodeSlashInfos(encodedSlashInfo []byte) ([]hmTypes.ValidatorSlashingInfo, error) {
	var slashingInfoList []hmTypes.ValidatorSlashingInfo
	err := rlp.DecodeBytes(encodedSlashInfo, &slashingInfoList)
	return slashingInfoList, err
}
