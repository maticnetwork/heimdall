package types

import (
	"github.com/maticnetwork/bor/rlp"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

// SortAndRLPEncodeSlashInfos  - RLP encoded slashing infos
func SortAndRLPEncodeSlashInfos(slashingInfos []*hmTypes.ValidatorSlashingInfo) ([]byte, error) {

	// Sort the slashingInfos by ID
	slashingInfos = hmTypes.SortValidatorSlashingInfoByID(slashingInfos)

	// Encode slashInfos
	encodedSlashInfos, err := rlp.EncodeToBytes(slashingInfos)

	return encodedSlashInfos, err
}

func RLPDecodeSlashInfos(encodedSlashInfo []byte) ([]*hmTypes.ValidatorSlashingInfo, error) {
	var slashingInfoList []*hmTypes.ValidatorSlashingInfo
	err := rlp.DecodeBytes(encodedSlashInfo, &slashingInfoList)
	return slashingInfoList, err

}
