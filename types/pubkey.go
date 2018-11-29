package types

import (
	"encoding/hex"

	"github.com/ethereum/go-ethereum/common/hexutil"
)

// PubKey pubkey
type PubKey [65]byte

// NewPubKey from byte array
func NewPubKey(data []byte) PubKey {
	var key PubKey
	copy(key[:], data[:])
	return key
}

// MarshalText returns the hex representation of a.
func (a PubKey) MarshalText() ([]byte, error) {
	return hexutil.Bytes(a[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (a *PubKey) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("PubKey", input, a[:])
}

// String return string representatin of key
func (a PubKey) String() string {
	return "0x" + hex.EncodeToString(a[:])
}
