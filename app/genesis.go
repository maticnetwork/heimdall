package app

import (
	"encoding/json"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	tmTypes "github.com/tendermint/tendermint/types"
)

// GenesisAccount genesis account
type GenesisAccount struct {
	Address       common.Address `json:"address"`
	Coins         sdk.Coins      `json:"coins"`
	Sequence      int64          `json:"sequence_number"`
	AccountNumber int64          `json:"account_number"`
}

// GenesisPubKey pubkey
type GenesisPubKey [65]byte

// NewGenesisPubKey from byte array
func NewGenesisPubKey(data []byte) GenesisPubKey {
	var key GenesisPubKey
	copy(key[:], data[:])
	return key
}

// MarshalText returns the hex representation of a.
func (a GenesisPubKey) MarshalText() ([]byte, error) {
	return hexutil.Bytes(a[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (a *GenesisPubKey) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("GenesisPubKey", input, a[:])
}

// GenesisValidator genesis validator
type GenesisValidator struct {
	Address    common.Address `json:"address"`
	StartEpoch int64          `json:"start_epoch"`
	EndEpoch   int64          `json:"end_epoch"`
	Power      int64          `json:"power"` // aka Amount
	PubKey     GenesisPubKey  `json:"pub_key"`
	Signer     common.Address `json:"signer"`
}

// ToTmValidator converts genesis valdator validator to Tendermint validator
func (v *GenesisValidator) ToTmValidator() tmTypes.Validator {
	var pubkeyBytes secp256k1.PubKeySecp256k1
	copy(pubkeyBytes[:], v.PubKey[:])

	return tmTypes.Validator{
		Address:     v.Signer.Bytes(),
		PubKey:      pubkeyBytes,
		VotingPower: v.Power,
	}
}

// GenesisState to Unmarshal
type GenesisState struct {
	Accounts   []GenesisAccount   `json:"accounts"`
	Validators []GenesisValidator `json:"validators"`
	GenTxs     []json.RawMessage  `json:"gentxs"`
}
