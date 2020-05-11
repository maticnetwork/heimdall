package types

import "encoding/binary"

const (
	// ModuleName defines the name of the module
	ModuleName = "sidechannel"

	// StoreKey is the store key string for bor
	StoreKey = ModuleName

	// RouterKey is the message route for bor
	RouterKey = ModuleName

	// QuerierRoute is the querier route for bor
	QuerierRoute = ModuleName

	// DefaultParamspace default name for parameter store
	DefaultParamspace = ModuleName

	// TStoreKey is the string store key for the param transient store
	TStoreKey = "transient_params"
)

var (
	// TxsKeyPrefix prefix for txs
	TxsKeyPrefix = []byte{0x01}

	// ValidatorsKeyPrefix prefix for validators
	ValidatorsKeyPrefix = []byte{0x02}
)

// TxStoreKey returns key used to get tx from store
func TxStoreKey(height int64, hash []byte) []byte {
	result := TxsStoreKey(height)
	result = append(result, hash...)
	return result
}

// TxsStoreKey returns key used to get txs from store
func TxsStoreKey(height int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(height))

	result := []byte{}
	result = append(result, TxsKeyPrefix...)
	result = append(result, b...)
	return result
}

// ValidatorsKey returns key used to get past-validators from store
func ValidatorsKey(height int64) []byte {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(height))

	result := []byte{}
	result = append(result, ValidatorsKeyPrefix...)
	result = append(result, b...)
	return result
}
