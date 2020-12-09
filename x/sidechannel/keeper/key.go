package keeper

import (
	"encoding/binary"
)

var (
	// TxsKeyPrefix prefix for txs
	TxsKeyPrefix = []byte{0x01}

	// ValidatorsKeyPrefix prefix for validators
	ValidatorsKeyPrefix = []byte{0x02}
)

// TxStoreKey returns key used to get tx from store
func TxStoreKey(height uint64, hash []byte) []byte {
	result := TxsStoreKey(height)
	result = append(result, hash...)
	return result
}

// TxsStoreKey returns key used to get txs from store
func TxsStoreKey(height uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, height)

	result := []byte{}
	result = append(result, TxsKeyPrefix...)
	result = append(result, b...)
	return result
}

// ValidatorsKey returns key used to get past-validators from store
func ValidatorsKey(height uint64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, height)

	result := []byte{}
	result = append(result, ValidatorsKeyPrefix...)
	result = append(result, b...)
	return result
}
