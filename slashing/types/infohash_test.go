package types

import (
	"fmt"
	"testing"

	hmTypes "github.com/maticnetwork/heimdall/types"
	"github.com/stretchr/testify/require"
)

func TestSlashingInfoRLPEncoding(t *testing.T) {
	var slashingInfoList []*hmTypes.ValidatorSlashingInfo

	// Input data
	slashingInfo1 := hmTypes.NewValidatorSlashingInfo(1, "1000", false)
	slashingInfo2 := hmTypes.NewValidatorSlashingInfo(2, "234", true)
	slashingInfoList = append(slashingInfoList, &slashingInfo1)
	slashingInfoList = append(slashingInfoList, &slashingInfo2)

	// Encoding
	encodedSlashInfos, err := SortAndRLPEncodeSlashInfos(slashingInfoList)
	t.Log("RLP encoded", "encodedSlashInfos", encodedSlashInfos, "error", err)
	require.Empty(t, err)

	// Decoding
	decodedSlashInfoList, err := RLPDecodeSlashInfos(encodedSlashInfos)
	require.Empty(t, err)
	t.Log("RLP Decoded data", "valID", decodedSlashInfoList[0].ID, "amount", decodedSlashInfoList[0].SlashedAmount, "isJailed", decodedSlashInfoList[0].IsJailed)
	t.Log("RLP Decoded data", "valID", decodedSlashInfoList[1].ID, "amount", decodedSlashInfoList[1].SlashedAmount, "isJailed", decodedSlashInfoList[1].IsJailed)

	// Assertions
	for i := 0; i < len(slashingInfoList); i++ {
		require.Equal(t, slashingInfoList[i].ID, decodedSlashInfoList[i].ID, "ID mismatch between slashInfoList and decodedSlashInfoList")
		require.Equal(t, slashingInfoList[i].SlashedAmount, decodedSlashInfoList[i].SlashedAmount, "Amount mismatch between slashInfoList and decodedSlashInfoList")
		require.Equal(t, slashingInfoList[i].IsJailed, decodedSlashInfoList[i].IsJailed, "JailStatus mismatch between slashInfoList and decodedSlashInfoList")
	}
}

func TestSlashingInfoRLPDecoding(t *testing.T) {
	slashInfoEncodedBytesStr := "0x000000000000000000d6d5019233353030303030303030303030303030303080"

	slashInfoEncodedBytes := []byte(slashInfoEncodedBytesStr)
	fmt.Println(slashInfoEncodedBytes)
}
