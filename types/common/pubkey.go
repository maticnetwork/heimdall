package common

import (
	"encoding/json"

	"gopkg.in/yaml.v2"

	cosmossecp256k1 "github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmprotocrypto "github.com/tendermint/tendermint/proto/tendermint/crypto"

	"github.com/maticnetwork/bor/common"
	"github.com/maticnetwork/bor/common/hexutil"
)

// PubKey pubkey
type PubKey []byte

// ZeroPubKey represents empty pub key
var ZeroPubKey = PubKey{}

// NewPubKey from byte array
func NewPubKey(data []byte) PubKey {
	return data
}

// NewPubKeyFromHex from byte array
func NewPubKeyFromHex(pk string) PubKey {
	return NewPubKey(common.FromHex(pk))
}

// MarshalText returns the hex representation of a.
func (a PubKey) MarshalText() ([]byte, error) {
	return hexutil.Bytes(a[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (a *PubKey) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("PubKey", input, a.Bytes()[:])
}

// String returns string representations of key
func (a PubKey) String() string {
	return common.ToHex(a)
}

// Bytes returns bytes for pubkey
func (a PubKey) Bytes() []byte {
	return a[:]
}

// Address returns address
func (a PubKey) Address() common.Address {
	return common.BytesToAddress(a.CryptoPubKey().Address().Bytes())
}

// CryptoPubKey returns crypto pub key for crypto
func (a PubKey) CryptoPubKey() crypto.PubKey {
	pk := make(secp256k1.PubKey, secp256k1.PubKeySize)
	copy(pk, a[:])
	return pk
}

// TMProtoCryptoPubKey returns crypto pub key for tendermint
func (a PubKey) TMProtoCryptoPubKey() tmprotocrypto.PublicKey {
	return tmprotocrypto.PublicKey{
		Sum: &tmprotocrypto.PublicKey_Secp256K1{
			Secp256K1: a,
		},
	}
}

// CosmosCryptoPubKey returns crypto pub key for cosmos
func (a *PubKey) CosmosCryptoPubKey() cryptotypes.PubKey {
	return CosmosCryptoPubKey(a.Bytes())
}

// TODO: check if any interface is implementing
// ABCIPubKey returns abci pubkey for cosmos
// func (a PubKey) ABCIPubKey() abci.PubKey {
// 	return tmTypes.TM2PB.PubKey(a.CryptoPubKey())
// }

// Marshal returns the raw address bytes. It is needed for protobuf compatibility.
func (a PubKey) Marshal() ([]byte, error) {
	return a.Bytes(), nil
}

// Unmarshal sets the address to the given data. It is needed for protobuf
// compatibility.
func (a *PubKey) Unmarshal(data []byte) error {
	*a = data
	return nil
}

// MarshalJSON marshals to JSON using Bech32.
func (a PubKey) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

// MarshalYAML marshals to YAML using Bech32.
func (a PubKey) MarshalYAML() (interface{}, error) {
	return a.String(), nil
}

// UnmarshalJSON unmarshals from JSON assuming Bech32 encoding.
func (a *PubKey) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	*a = common.FromHex(s)
	return nil
}

// UnmarshalYAML unmarshals from JSON assuming Bech32 encoding.
func (a *PubKey) UnmarshalYAML(data []byte) error {
	var s string
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	*a = common.FromHex(s)
	return nil
}

//
// Utility methods
//

// CosmosCryptoPubKey returns crypto pub key for cosmos
func CosmosCryptoPubKey(pk []byte) cryptotypes.PubKey {
	return &cosmossecp256k1.PubKey{
		Key: pk,
	}
}
