package test

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"
	tmerkel "github.com/tendermint/tendermint/crypto/merkle"

	"github.com/maticnetwork/heimdall/helper"
)

func TestTxDecode(t *testing.T) {
	tx, err := helper.TendermintTxDecode("wWhvHPg6AQHY1wEBlP+zHe/ZNZTQii57ULFjrJulHewY2NcBAZT/sx3v2TWU0Ioue1CxY6ybpR3sGICEXTLzJQ==")
	if err != nil {
		t.Error(err)
	} else {
		expected := "c1686f1cf83a0101d8d7010194ffb31defd93594d08a2e7b50b163ac9ba51dec18d8d7010194ffb31defd93594d08a2e7b50b163ac9ba51dec1880845d32f325"
		require.Equal(t, expected, hex.EncodeToString(tx), "Tx encoding should match")
	}
}

func TestTxHash(t *testing.T) {
	// These allocations will be removed once Txs is switched to [][]byte,
	// ref #2603. This i[s because golang does not allow type casting slices without unsafe
	txs := []string{
		"CA7FE14F21B58259D87D6EBEC5E316865C100C22B7634B485AD5836AF40B37B9",
		"BDEF26BFF25D71CA8AB036581264DD12EDB0A183CF06D93EC1287AE8662F1BD6",
		"E9B8BC581775D36133CB5A443F844785A3DFF9F0DCC027A5E90636718405ACCB",
		"C6C0F589AA507BD0FAC79AF9C41D590C6E864BC4BE1B3EAC87C46F02CF21970E",
		"B1DB4473F55D259B6E0F8744D0AE3A146D82CB42CC587B9070E80C39E0DFF09A",
	}

	txBzs := make([][]byte, len(txs))
	for i := 0; i < len(txs); i++ {
		txBzs[i], _ = hex.DecodeString(txs[i])
	}
	expected := "22d3b708469ea7107c230a558d88ef82c18f1ce7f716e84f494c72edf50aeb0e"
	hash := hex.EncodeToString(tmerkel.SimpleHashFromByteSlices(txBzs))
	require.Equal(t, expected, hash, "Tx hash should match")
}

func TestTxMerkle(t *testing.T) {
	txs := []string{
		"CA7FE14F21B58259D87D6EBEC5E316865C100C22B7634B485AD5836AF40B37B9",
		"BDEF26BFF25D71CA8AB036581264DD12EDB0A183CF06D93EC1287AE8662F1BD6",
		"E9B8BC581775D36133CB5A443F844785A3DFF9F0DCC027A5E90636718405ACCB",
		"C6C0F589AA507BD0FAC79AF9C41D590C6E864BC4BE1B3EAC87C46F02CF21970E",
		"B1DB4473F55D259B6E0F8744D0AE3A146D82CB42CC587B9070E80C39E0DFF09A",
	}

	txBzs := make([][]byte, len(txs))
	for i := 0; i < len(txs); i++ {
		txBzs[i], _ = hex.DecodeString(txs[i])
	}

	rootHash, proof := tmerkel.SimpleProofsFromByteSlices(txBzs)
	result := helper.GetMerkleProofList(proof[2])

	expected := "22d3b708469ea7107c230a558d88ef82c18f1ce7f716e84f494c72edf50aeb0e"
	require.Equal(t, expected, hex.EncodeToString(rootHash), "Tx hash should match")
	require.Equal(t, 3, len(result), "Proof length should match")
}
