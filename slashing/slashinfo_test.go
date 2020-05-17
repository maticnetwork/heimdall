package slashing_test

import (
	"encoding/hex"
	"testing"

	slashingTypes "github.com/maticnetwork/heimdall/slashing/types"
	"github.com/stretchr/testify/require"

	hmTypes "github.com/maticnetwork/heimdall/types"
)

func TestSlashingInfoRLP(t *testing.T) {
	var slashingInfoList []*hmTypes.ValidatorSlashingInfo

	// Input data
	slashingInfo1 := hmTypes.NewValidatorSlashingInfo(1, uint64(1000), false) // on contract, "false" value decoded as "0"
	slashingInfo2 := hmTypes.NewValidatorSlashingInfo(2, uint64(234), true)   // on contract, "true" value decoded as "1"
	slashingInfoList = append(slashingInfoList, &slashingInfo1)
	slashingInfoList = append(slashingInfoList, &slashingInfo2)

	// Encoding
	encodedSlashInfos, err := slashingTypes.SortAndRLPEncodeSlashInfos(slashingInfoList)
	t.Log("RLP encoded", "encodedSlashInfos", hex.EncodeToString(encodedSlashInfos), "error", err)
	require.Empty(t, err)

	// Decoding
	decodedSlashInfoList, err := slashingTypes.RLPDecodeSlashInfos(encodedSlashInfos)
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

func TestSlashingInfoRLPEncoding(t *testing.T) {
	var slashingInfoList []*hmTypes.ValidatorSlashingInfo

	// Input data
	slashingInfo2 := hmTypes.NewValidatorSlashingInfo(2, uint64(120), false)
	slashingInfoList = append(slashingInfoList, &slashingInfo2)

	// Encoding
	encodedSlashInfos, err := slashingTypes.SortAndRLPEncodeSlashInfos(slashingInfoList)
	t.Log("RLP encoded", "encodedSlashInfos", hex.EncodeToString(encodedSlashInfos), "error", err)
	require.Empty(t, err)
}

func TestSlashingInfoRLPDecoding(t *testing.T) {
	// input data
	slashInfoEncodedBytesStr := "cdcc0289068155a43676e0000000"
	slashInfoEncodedBytes, err := hex.DecodeString(slashInfoEncodedBytesStr)
	require.Empty(t, err)

	// decoding input
	slashInfos, err := slashingTypes.RLPDecodeSlashInfos(slashInfoEncodedBytes)
	require.Empty(t, err)
	t.Log("RLP decoded data", "slashInfos - ", slashInfos)
}
