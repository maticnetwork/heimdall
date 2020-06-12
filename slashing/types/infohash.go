package types

import (
	"bytes"
	"math/big"
	"sort"

	"github.com/maticnetwork/bor/rlp"
	"github.com/maticnetwork/heimdall/helper"
	hmTypes "github.com/maticnetwork/heimdall/types"
)

type ModifiedSlashInfo struct {
	ID            hmTypes.ValidatorID `json:"ID"`
	SlashedAmount *big.Int            `json:"SlashedAmount"`
	IsJailed      []byte              `json:"IsJailed"`
}

// SortModifiedSlashInfoByID - Sorts ModifiedSlashInfo By ID
func sortModifiedSlashInfoByID(modifiedSlashInfos []*ModifiedSlashInfo) []*ModifiedSlashInfo {
	sort.Slice(modifiedSlashInfos, func(i, j int) bool { return modifiedSlashInfos[i].ID < modifiedSlashInfos[j].ID })
	return modifiedSlashInfos
}

// SortAndRLPEncodeSlashInfos  - RLP encoded slashing infos
func SortAndRLPEncodeSlashInfos(slashingInfos []*hmTypes.ValidatorSlashingInfo) ([]byte, error) {

	// convert slashingInfos to modifiedSlashingInfos
	var updatedslashInfos []*ModifiedSlashInfo
	for _, slashInfo := range slashingInfos {
		modifiedSlashInfo, err := slashInfoToModified(slashInfo)
		if err != nil {
			return nil, err
		}
		updatedslashInfos = append(updatedslashInfos, modifiedSlashInfo)
	}

	// Sort the slashingInfos by ID
	updatedslashInfos = sortModifiedSlashInfoByID(updatedslashInfos)

	// Encode encodedSlashInfos
	encodedSlashInfos, err := rlp.EncodeToBytes(updatedslashInfos)

	return encodedSlashInfos, err
}

func slashInfoToModified(slashInfo *hmTypes.ValidatorSlashingInfo) (modifiedSlashInfo *ModifiedSlashInfo, err error) {
	amount, err := helper.GetAmountFromPower(int64(slashInfo.SlashedAmount))
	if err != nil {
		return modifiedSlashInfo, err
	}

	// converting jailed from boolean to Byte. as boolean rlp is incompatible Issue - https://github.com/hamdiallam/Solidity-RLP/issues/5
	jailedByte := []byte{0x00}
	if slashInfo.IsJailed {
		jailedByte = []byte{0x01}
	}

	// convert slashing amount to 10^18. required for contracts.
	modifiedSlashInfo = &ModifiedSlashInfo{
		ID:            slashInfo.ID,
		SlashedAmount: amount,
		IsJailed:      jailedByte,
	}

	return modifiedSlashInfo, err
}

func RLPDecodeSlashInfos(encodedSlashInfo []byte) ([]*hmTypes.ValidatorSlashingInfo, error) {
	var modifiedSlashInfoList []*ModifiedSlashInfo
	err := rlp.DecodeBytes(encodedSlashInfo, &modifiedSlashInfoList)
	if err != nil {
		return nil, err
	}
	// convert modifiedSlashingInfos to slashingInfos
	var updatedslashInfos []*hmTypes.ValidatorSlashingInfo
	for _, modifiedSlashInfo := range modifiedSlashInfoList {
		slashInfo, err := modifiedToSlashInfo(modifiedSlashInfo)
		if err != nil {
			return nil, err
		}
		updatedslashInfos = append(updatedslashInfos, slashInfo)
	}

	return updatedslashInfos, err
}

func modifiedToSlashInfo(modifiedSlashInfo *ModifiedSlashInfo) (slashInfo *hmTypes.ValidatorSlashingInfo, err error) {
	amount := modifiedSlashInfo.SlashedAmount
	power, err := helper.GetPowerFromAmount(amount)
	if err != nil {
		return slashInfo, err
	}

	// converting jailed from boolean to bytes. as boolean RLP is incompatible Issue - https://github.com/hamdiallam/Solidity-RLP/issues/5
	var jailedBool bool
	if bytes.Equal(modifiedSlashInfo.IsJailed, []byte{0x01}) {
		jailedBool = true
	}

	// convert slashing amount to 10^18. required for contracts.
	slashInfo = &hmTypes.ValidatorSlashingInfo{
		ID:            modifiedSlashInfo.ID,
		SlashedAmount: power.Uint64(),
		IsJailed:      jailedBool,
	}

	return slashInfo, err
}
