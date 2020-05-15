package types

import (
	"encoding/binary"
)

const (
	// ModuleName is the name of the module
	ModuleName = "slashing"

	// StoreKey is the store key string for slashing
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

	// QuerierRoute is the querier route for slashing
	QuerierRoute = ModuleName
)

// Keys for slashing store
// Items are stored with the following key: values
//
// - 0x01<consAddress_Bytes>: ValidatorSigningInfo
//
// - 0x02<consAddress_Bytes><period_Bytes>: bool
//
// - 0x03<accAddr_Bytes>: crypto.PubKey
var (
	DefaultValue = []byte{0x01} // Value to store for slashing sequence

	ValidatorSigningInfoKey         = []byte{0x01} // Prefix for signing info
	ValidatorMissedBlockBitArrayKey = []byte{0x02} // Prefix for missed block bit array
	AddrPubkeyRelationKey           = []byte{0x03} // Prefix for address-pubkey relation
	TotalSlashedAmountKey           = []byte{0x04} // Prefix for total slashed amount stored in buffer
	BufferValSlashingInfoKey        = []byte{0x05} // Prefix for Slashing Info stored in buffer
	TickValSlashingInfoKey          = []byte{0x06} // Prefix for Slashing Info stored after tick tx
	SlashingSequenceKey             = []byte{0x07} // prefix for each key for slashing sequence map
	TickCountKey                    = []byte{0x08} // key to store Tick counts
)

// GetValidatorSigningInfoKey - stored by *valID*
func GetValidatorSigningInfoKey(valID []byte) []byte {
	return append(ValidatorSigningInfoKey, valID...)
}

// GetValidatorMissedBlockBitArrayPrefixKey - stored by *Consensus* address (not operator address)
func GetValidatorMissedBlockBitArrayPrefixKey(valID []byte) []byte {
	return append(ValidatorMissedBlockBitArrayKey, valID...)
}

// GetValidatorMissedBlockBitArrayKey - stored by *Consensus* address (not operator address)
func GetValidatorMissedBlockBitArrayKey(valID []byte, i int64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(i))
	return append(GetValidatorMissedBlockBitArrayPrefixKey(valID), b...)
}

// GetBufferValSlashingInfoKey - gets buffer val slashing info key
func GetBufferValSlashingInfoKey(id []byte) []byte {
	return append(BufferValSlashingInfoKey, id...)
}

func GetTickValSlashingInfoKey(id []byte) []byte {
	return append(TickValSlashingInfoKey, id...)
}

// GetSlashingSequenceKey returns slashing sequence key
func GetSlashingSequenceKey(sequence string) []byte {
	return append(SlashingSequenceKey, []byte(sequence)...)
}
