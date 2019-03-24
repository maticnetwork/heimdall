package types

import (
	"encoding/hex"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmTypes "github.com/tendermint/tendermint/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// PubKey pubkey
type PubKey [65]byte

// ZeroPubKey represents empty pub key
var ZeroPubKey = PubKey{}

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

// String returns string representatin of key
func (a PubKey) String() string {
	return "0x" + hex.EncodeToString(a[:])
}

// Bytes returns bytes for pubkey
func (a PubKey) Bytes() []byte {
	return a[:]
}

// Address returns address
func (a PubKey) Address() common.Address {
	return common.BytesToAddress(a.CryptoPubKey().Address().Bytes())
}

// CryptoPubKey returns crypto pub key for tendermint
func (a PubKey) CryptoPubKey() crypto.PubKey {
	var pubkeyBytes secp256k1.PubKeySecp256k1
	copy(pubkeyBytes[:], a[:])
	return pubkeyBytes
}

// ABCIPubKey returns abci pubkey for cosmos
func (a PubKey) ABCIPubKey() abci.PubKey {
	return tmTypes.TM2PB.PubKey(a.CryptoPubKey())
}
